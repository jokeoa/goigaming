package game

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jokeoa/goigaming/internal/core/domain"
	"github.com/jokeoa/goigaming/internal/core/ports"
	"github.com/jokeoa/goigaming/internal/repository/postgres"
	"github.com/shopspring/decimal"
)

type Service struct {
	pool       *pgxpool.Pool
	tableRepo  ports.PokerTableRepository
	playerRepo ports.PokerPlayerRepository
	handRepo   ports.PokerHandRepository
	walletSvc  ports.WalletService
	userSvc    ports.UserService
	hubManager *HubManager
	playerFn   func(db postgres.DBTX) ports.PokerPlayerRepository
}

func NewService(
	pool *pgxpool.Pool,
	tableRepo ports.PokerTableRepository,
	playerRepo ports.PokerPlayerRepository,
	handRepo ports.PokerHandRepository,
	walletSvc ports.WalletService,
	userSvc ports.UserService,
	hubManager *HubManager,
	playerFn func(db postgres.DBTX) ports.PokerPlayerRepository,
) *Service {
	return &Service{
		pool:       pool,
		tableRepo:  tableRepo,
		playerRepo: playerRepo,
		handRepo:   handRepo,
		walletSvc:  walletSvc,
		userSvc:    userSvc,
		hubManager: hubManager,
		playerFn:   playerFn,
	}
}

func (s *Service) CreateTable(ctx context.Context, table domain.PokerTable) (domain.PokerTable, error) {
	if table.SmallBlind.LessThanOrEqual(decimal.Zero) || table.BigBlind.LessThanOrEqual(decimal.Zero) {
		return domain.PokerTable{}, domain.ErrInvalidBetAmount
	}
	if table.MinBuyIn.LessThanOrEqual(decimal.Zero) || table.MaxBuyIn.LessThan(table.MinBuyIn) {
		return domain.PokerTable{}, domain.ErrInvalidBuyIn
	}
	if table.MaxPlayers < 2 || table.MaxPlayers > 9 {
		return domain.PokerTable{}, fmt.Errorf("max players must be between 2 and 9")
	}

	table.Status = domain.TableStatusWaiting

	created, err := s.tableRepo.Create(ctx, table)
	if err != nil {
		return domain.PokerTable{}, fmt.Errorf("PokerService.CreateTable: %w", err)
	}

	return created, nil
}

func (s *Service) GetTable(ctx context.Context, tableID uuid.UUID) (domain.PokerTable, error) {
	table, err := s.tableRepo.FindByID(ctx, tableID)
	if err != nil {
		return domain.PokerTable{}, fmt.Errorf("PokerService.GetTable: %w", err)
	}
	return table, nil
}

func (s *Service) ListTables(ctx context.Context) ([]domain.PokerTable, error) {
	tables, err := s.tableRepo.FindActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("PokerService.ListTables: %w", err)
	}
	return tables, nil
}

func (s *Service) JoinTable(ctx context.Context, tableID, userID uuid.UUID, seatNumber int, buyIn string) (domain.PokerPlayer, error) {
	buyInAmount, err := decimal.NewFromString(buyIn)
	if err != nil {
		return domain.PokerPlayer{}, domain.ErrInvalidBuyIn
	}

	table, err := s.tableRepo.FindByID(ctx, tableID)
	if err != nil {
		return domain.PokerPlayer{}, fmt.Errorf("PokerService.JoinTable: %w", err)
	}

	if buyInAmount.LessThan(table.MinBuyIn) || buyInAmount.GreaterThan(table.MaxBuyIn) {
		return domain.PokerPlayer{}, domain.ErrInvalidBuyIn
	}

	if seatNumber < 1 || seatNumber > table.MaxPlayers {
		return domain.PokerPlayer{}, domain.ErrSeatTaken
	}

	user, err := s.userSvc.GetByID(ctx, userID)
	if err != nil {
		return domain.PokerPlayer{}, fmt.Errorf("PokerService.JoinTable get user: %w", err)
	}

	if _, err := s.walletSvc.Withdraw(ctx, userID, buyIn); err != nil {
		return domain.PokerPlayer{}, fmt.Errorf("PokerService.JoinTable withdraw: %w", err)
	}

	var player domain.PokerPlayer

	err = postgres.RunInTx(ctx, s.pool, func(tx pgx.Tx) error {
		playerRepo := s.playerFn(tx)

		count, err := playerRepo.CountByTableID(ctx, tableID)
		if err != nil {
			return fmt.Errorf("count players: %w", err)
		}
		if count >= table.MaxPlayers {
			return domain.ErrTableFull
		}

		player, err = playerRepo.Create(ctx, domain.PokerPlayer{
			TableID:    tableID,
			UserID:     userID,
			Username:   user.Username,
			Stack:      buyInAmount,
			SeatNumber: seatNumber,
			Status:     domain.PlayerStatusActive,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if _, depErr := s.walletSvc.Deposit(ctx, userID, buyIn); depErr != nil {
			return domain.PokerPlayer{}, fmt.Errorf("PokerService.JoinTable refund failed: create=%w, refund=%v", err, depErr)
		}
		return domain.PokerPlayer{}, err
	}

	hub := s.hubManager.GetOrCreateHub(ctx, table)
	resultCh := make(chan HubResult, 1)
	if sendErr := hub.Send(HubEvent{
		Type:     EventPlayerJoin,
		UserID:   userID,
		PlayerID: player.ID,
		SeatNum:  seatNumber,
		BuyIn:    buyInAmount,
		Username: user.Username,
		ResultCh: resultCh,
	}); sendErr != nil {
		return domain.PokerPlayer{}, fmt.Errorf("PokerService.JoinTable hub send: %w", sendErr)
	}
	result := <-resultCh
	if result.Err != nil {
		return domain.PokerPlayer{}, result.Err
	}

	return player, nil
}

func (s *Service) LeaveTable(ctx context.Context, tableID, userID uuid.UUID) error {
	player, err := s.playerRepo.FindByTableAndUser(ctx, tableID, userID)
	if err != nil {
		return fmt.Errorf("PokerService.LeaveTable: %w", err)
	}

	hub := s.hubManager.GetHub(tableID)
	if hub != nil {
		resultCh := make(chan HubResult, 1)
		if err := hub.Send(HubEvent{
			Type:     EventPlayerLeave,
			UserID:   userID,
			PlayerID: player.ID,
			ResultCh: resultCh,
		}); err == nil {
			result := <-resultCh
			if result.Stack != nil {
				player.Stack = *result.Stack
			}
		}
	}

	if player.Stack.IsPositive() {
		if _, err := s.walletSvc.Deposit(ctx, userID, player.Stack.String()); err != nil {
			return fmt.Errorf("PokerService.LeaveTable deposit: %w", err)
		}
	}

	if err := s.playerRepo.Delete(ctx, player.ID); err != nil {
		return fmt.Errorf("PokerService.LeaveTable delete: %w", err)
	}

	return nil
}

func (s *Service) GetTableState(ctx context.Context, tableID uuid.UUID) (domain.WSTableState, error) {
	hub := s.hubManager.GetHub(tableID)
	if hub != nil {
		return hub.State().ToWSTableState(), nil
	}

	table, err := s.tableRepo.FindByID(ctx, tableID)
	if err != nil {
		return domain.WSTableState{}, fmt.Errorf("PokerService.GetTableState: %w", err)
	}

	players, err := s.playerRepo.FindByTableID(ctx, tableID)
	if err != nil {
		return domain.WSTableState{}, fmt.Errorf("PokerService.GetTableState players: %w", err)
	}

	wsPlayers := make([]domain.WSPlayerInfo, len(players))
	for i, p := range players {
		wsPlayers[i] = domain.WSPlayerInfo{
			UserID:     p.UserID,
			Username:   p.Username,
			Stack:      p.Stack,
			SeatNumber: p.SeatNumber,
			Status:     p.Status,
		}
	}

	return domain.WSTableState{
		TableID:    table.ID,
		Name:       table.Name,
		SmallBlind: table.SmallBlind,
		BigBlind:   table.BigBlind,
		Stage:      domain.StageWaiting,
		Players:    wsPlayers,
	}, nil
}
