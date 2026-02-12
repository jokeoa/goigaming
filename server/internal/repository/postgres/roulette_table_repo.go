package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jokeoa/goigaming/internal/core/domain"
)

type RouletteTableRepo struct {
	db DBTX
}

func NewRouletteTableRepo(db DBTX) *RouletteTableRepo {
	return &RouletteTableRepo{db: db}
}

func (r *RouletteTableRepo) FindByID(ctx context.Context, id uuid.UUID) (domain.RouletteTable, error) {
	query := `
		SELECT id, name, min_bet, max_bet, status, created_at
		FROM roulette_tables
		WHERE id = $1
	`

	var t domain.RouletteTable
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.Name, &t.MinBet, &t.MaxBet, &t.Status, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return t, domain.ErrTableNotFound
		}
		return t, fmt.Errorf("RouletteTableRepo.FindByID: %w", err)
	}

	return t, nil
}

func (r *RouletteTableRepo) FindActive(ctx context.Context) ([]domain.RouletteTable, error) {
	query := `
		SELECT id, name, min_bet, max_bet, status, created_at
		FROM roulette_tables
		WHERE status = 'active'
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("RouletteTableRepo.FindActive: %w", err)
	}
	defer rows.Close()

	var tables []domain.RouletteTable
	for rows.Next() {
		var t domain.RouletteTable
		if err := rows.Scan(
			&t.ID, &t.Name, &t.MinBet, &t.MaxBet, &t.Status, &t.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("RouletteTableRepo.FindActive scan: %w", err)
		}
		tables = append(tables, t)
	}

	return tables, rows.Err()
}
