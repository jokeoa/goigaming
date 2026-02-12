package roulette

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jokeoa/goigaming/internal/core/domain"
	"github.com/jokeoa/goigaming/internal/core/ports"
	"github.com/jokeoa/goigaming/internal/repository/postgres"
	"github.com/shopspring/decimal"
)

type Service struct {
	pool      *pgxpool.Pool
	walletSvc ports.WalletService
	tableRepo ports.RouletteTableRepository
	roundRepo ports.RouletteRoundRepository
	betRepo   ports.RouletteBetRepository
	betFn     func(db postgres.DBTX) ports.RouletteBetRepository
}

func NewService(
	pool *pgxpool.Pool,
	walletSvc ports.WalletService,
	tableRepo ports.RouletteTableRepository,
	roundRepo ports.RouletteRoundRepository,
	betRepo ports.RouletteBetRepository,
	betFn func(db postgres.DBTX) ports.RouletteBetRepository,
) *Service {
	return &Service{
		pool:      pool,
		walletSvc: walletSvc,
		tableRepo: tableRepo,
		roundRepo: roundRepo,
		betRepo:   betRepo,
		betFn:     betFn,
	}
}

func (s *Service) GetTable(ctx context.Context, tableID uuid.UUID) (domain.RouletteTable, error) {
	table, err := s.tableRepo.FindByID(ctx, tableID)
	if err != nil {
		return domain.RouletteTable{}, fmt.Errorf("RouletteService.GetTable: %w", err)
	}
	return table, nil
}

func (s *Service) ListActiveTables(ctx context.Context) ([]domain.RouletteTable, error) {
	tables, err := s.tableRepo.FindActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("RouletteService.ListActiveTables: %w", err)
	}
	return tables, nil
}

func (s *Service) GetCurrentRound(ctx context.Context, tableID uuid.UUID) (domain.RouletteRound, error) {
	round, err := s.roundRepo.FindCurrentByTableID(ctx, tableID)
	if err != nil {
		return domain.RouletteRound{}, fmt.Errorf("RouletteService.GetCurrentRound: %w", err)
	}
	return round, nil
}

func (s *Service) GetRound(ctx context.Context, roundID uuid.UUID) (domain.RouletteRound, error) {
	round, err := s.roundRepo.FindByID(ctx, roundID)
	if err != nil {
		return domain.RouletteRound{}, fmt.Errorf("RouletteService.GetRound: %w", err)
	}
	return round, nil
}

func (s *Service) GetRoundHistory(ctx context.Context, tableID uuid.UUID, limit, offset int) ([]domain.RouletteRound, error) {
	rounds, err := s.roundRepo.FindSettledByTableID(ctx, tableID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("RouletteService.GetRoundHistory: %w", err)
	}
	return rounds, nil
}

func (s *Service) PlaceBet(ctx context.Context, userID, tableID, roundID uuid.UUID, betType, betValue, amount string) (domain.RouletteBet, error) {
	betAmount, err := decimal.NewFromString(amount)
	if err != nil || !betAmount.IsPositive() {
		return domain.RouletteBet{}, domain.ErrInvalidAmount
	}

	if !ValidateBetType(betType) {
		return domain.RouletteBet{}, domain.ErrInvalidBetType
	}

	if !ValidateBetValue(betType, betValue) {
		return domain.RouletteBet{}, domain.ErrInvalidBetType
	}

	round, err := s.roundRepo.FindByID(ctx, roundID)
	if err != nil {
		return domain.RouletteBet{}, fmt.Errorf("RouletteService.PlaceBet: %w", err)
	}

	if round.TableID != tableID {
		return domain.RouletteBet{}, domain.ErrRoundNotFound
	}

	if round.SettledAt != nil {
		return domain.RouletteBet{}, domain.ErrBettingClosed
	}
	if round.BettingEndsAt != nil && time.Now().After(*round.BettingEndsAt) {
		return domain.RouletteBet{}, domain.ErrBettingClosed
	}

	table, err := s.tableRepo.FindByID(ctx, round.TableID)
	if err != nil {
		return domain.RouletteBet{}, fmt.Errorf("RouletteService.PlaceBet get table: %w", err)
	}

	if table.Status != domain.RouletteTableStatusActive {
		return domain.RouletteBet{}, domain.ErrTableNotActive
	}

	if betAmount.LessThan(table.MinBet) || betAmount.GreaterThan(table.MaxBet) {
		return domain.RouletteBet{}, domain.ErrBetAmountOutOfRange
	}

	if _, err := s.walletSvc.Withdraw(ctx, userID, amount); err != nil {
		return domain.RouletteBet{}, fmt.Errorf("RouletteService.PlaceBet withdraw: %w", err)
	}

	var bet domain.RouletteBet

	err = postgres.RunInTx(ctx, s.pool, func(tx pgx.Tx) error {
		betRepo := s.betFn(tx)

		bet, err = betRepo.Create(ctx, domain.RouletteBet{
			RoundID:  roundID,
			UserID:   userID,
			BetType:  betType,
			BetValue: betValue,
			Amount:   betAmount,
			Payout:   decimal.Zero,
			Status:   domain.RouletteBetStatusPending,
		})
		return err
	})
	if err != nil {
		if _, depErr := s.walletSvc.Deposit(ctx, userID, amount); depErr != nil {
			return domain.RouletteBet{}, fmt.Errorf("RouletteService.PlaceBet refund failed: create=%w, refund=%v", err, depErr)
		}
		return domain.RouletteBet{}, fmt.Errorf("RouletteService.PlaceBet: %w", err)
	}

	return bet, nil
}

func (s *Service) GetUserBets(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.RouletteBet, error) {
	bets, err := s.betRepo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("RouletteService.GetUserBets: %w", err)
	}
	return bets, nil
}
