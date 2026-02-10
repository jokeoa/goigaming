package postgres

import (
	"context"
	"fmt"

	"github.com/jokeoa/goigaming/internal/core/domain"
)

type TransactionRepository struct {
	db DBTX
}

func NewTransactionRepository(db DBTX) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, t domain.Transaction) (domain.Transaction, error) {
	query := `
		INSERT INTO transactions (wallet_id, amount, balance_after, reference_type, reference_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, wallet_id, amount, balance_after, reference_type, reference_id, created_at
	`

	var tx domain.Transaction
	err := r.db.QueryRow(ctx, query, t.WalletID, t.Amount, t.BalanceAfter, t.ReferenceType, t.ReferenceID).Scan(
		&tx.ID, &tx.WalletID, &tx.Amount, &tx.BalanceAfter, &tx.ReferenceType, &tx.ReferenceID, &tx.CreatedAt,
	)
	if err != nil {
		return tx, fmt.Errorf("TransactionRepository.Create: %w", err)
	}

	return tx, nil
}

func (r *TransactionRepository) FindByWalletID(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, error) {
	query := `
		SELECT id, wallet_id, amount, balance_after, reference_type, reference_id, created_at
		FROM transactions
		WHERE wallet_id = $1
	`
	args := []any{filter.WalletID}
	argIdx := 2

	if filter.ReferenceType != "" {
		query += fmt.Sprintf(" AND reference_type = $%d", argIdx)
		args = append(args, filter.ReferenceType)
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
		argIdx++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("TransactionRepository.FindByWalletID: %w", err)
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(&t.ID, &t.WalletID, &t.Amount, &t.BalanceAfter, &t.ReferenceType, &t.ReferenceID, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("TransactionRepository.FindByWalletID scan: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("TransactionRepository.FindByWalletID rows: %w", err)
	}

	return transactions, nil
}
