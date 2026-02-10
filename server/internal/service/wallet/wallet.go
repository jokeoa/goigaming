package wallet

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jokeoa/igaming/internal/core/domain"
	"github.com/jokeoa/igaming/internal/core/ports"
	"github.com/jokeoa/igaming/internal/repository/postgres"
	"github.com/shopspring/decimal"
)

const maxRetries = 3

type Service struct {
	pool       *pgxpool.Pool
	walletFn   func(db postgres.DBTX) ports.WalletRepository
	txFn       func(db postgres.DBTX) ports.TransactionRepository
	walletRepo ports.WalletRepository
	txRepo     ports.TransactionRepository
}

func NewService(
	pool *pgxpool.Pool,
	walletRepo ports.WalletRepository,
	txRepo ports.TransactionRepository,
	walletFn func(db postgres.DBTX) ports.WalletRepository,
	txFn func(db postgres.DBTX) ports.TransactionRepository,
) *Service {
	return &Service{
		pool:       pool,
		walletFn:   walletFn,
		txFn:       txFn,
		walletRepo: walletRepo,
		txRepo:     txRepo,
	}
}

func (s *Service) CreateWallet(ctx context.Context, userID uuid.UUID) (domain.Wallet, error) {
	w, err := s.walletRepo.Create(ctx, domain.Wallet{
		UserID:  userID,
		Balance: decimal.Zero,
	})
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("WalletService.CreateWallet: %w", err)
	}
	return w, nil
}

func (s *Service) GetBalance(ctx context.Context, userID uuid.UUID) (domain.Wallet, error) {
	w, err := s.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("WalletService.GetBalance: %w", err)
	}
	return w, nil
}

func (s *Service) Deposit(ctx context.Context, userID uuid.UUID, amount string) (domain.Wallet, error) {
	amt, err := decimal.NewFromString(amount)
	if err != nil {
		return domain.Wallet{}, domain.ErrInvalidAmount
	}
	if !amt.IsPositive() {
		return domain.Wallet{}, domain.ErrInvalidAmount
	}

	var result domain.Wallet

	for attempt := range maxRetries {
		result, err = s.executeDeposit(ctx, userID, amt)
		if err == nil {
			return result, nil
		}
		if !errors.Is(err, domain.ErrOptimisticLock) {
			return domain.Wallet{}, err
		}
		if attempt == maxRetries-1 {
			return domain.Wallet{}, fmt.Errorf("WalletService.Deposit: max retries exceeded: %w", err)
		}
	}

	return result, nil
}

func (s *Service) executeDeposit(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) (domain.Wallet, error) {
	var result domain.Wallet

	err := postgres.RunInTx(ctx, s.pool, func(tx pgx.Tx) error {
		walletRepo := s.walletFn(tx)
		txRepo := s.txFn(tx)

		w, err := walletRepo.FindByUserID(ctx, userID)
		if err != nil {
			return err
		}

		newBalance := w.Balance.Add(amount)

		updated, err := walletRepo.UpdateBalance(ctx, userID, newBalance, w.Version)
		if err != nil {
			return err
		}

		_, err = txRepo.Create(ctx, domain.Transaction{
			WalletID:      userID,
			Amount:        amount,
			BalanceAfter:  newBalance,
			ReferenceType: "deposit",
		})
		if err != nil {
			return fmt.Errorf("WalletService.Deposit create transaction: %w", err)
		}

		result = updated
		return nil
	})

	return result, err
}

func (s *Service) Withdraw(ctx context.Context, userID uuid.UUID, amount string) (domain.Wallet, error) {
	amt, err := decimal.NewFromString(amount)
	if err != nil {
		return domain.Wallet{}, domain.ErrInvalidAmount
	}
	if !amt.IsPositive() {
		return domain.Wallet{}, domain.ErrInvalidAmount
	}

	var result domain.Wallet

	for attempt := range maxRetries {
		result, err = s.executeWithdraw(ctx, userID, amt)
		if err == nil {
			return result, nil
		}
		if !errors.Is(err, domain.ErrOptimisticLock) {
			return domain.Wallet{}, err
		}
		if attempt == maxRetries-1 {
			return domain.Wallet{}, fmt.Errorf("WalletService.Withdraw: max retries exceeded: %w", err)
		}
	}

	return result, nil
}

func (s *Service) executeWithdraw(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) (domain.Wallet, error) {
	var result domain.Wallet

	err := postgres.RunInTx(ctx, s.pool, func(tx pgx.Tx) error {
		walletRepo := s.walletFn(tx)
		txRepo := s.txFn(tx)

		w, err := walletRepo.FindByUserID(ctx, userID)
		if err != nil {
			return err
		}

		if w.Balance.LessThan(amount) {
			return domain.ErrInsufficientFunds
		}

		newBalance := w.Balance.Sub(amount)

		updated, err := walletRepo.UpdateBalance(ctx, userID, newBalance, w.Version)
		if err != nil {
			return err
		}

		_, err = txRepo.Create(ctx, domain.Transaction{
			WalletID:      userID,
			Amount:        amount.Neg(),
			BalanceAfter:  newBalance,
			ReferenceType: "withdrawal",
		})
		if err != nil {
			return fmt.Errorf("WalletService.Withdraw create transaction: %w", err)
		}

		result = updated
		return nil
	})

	return result, err
}

func (s *Service) GetTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Transaction, error) {
	w, err := s.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("WalletService.GetTransactions: %w", err)
	}

	txs, err := s.txRepo.FindByWalletID(ctx, domain.TransactionFilter{
		WalletID: w.UserID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("WalletService.GetTransactions: %w", err)
	}

	return txs, nil
}
