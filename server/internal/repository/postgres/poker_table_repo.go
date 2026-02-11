package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jokeoa/goigaming/internal/core/domain"
)

type PokerTableRepository struct {
	db DBTX
}

func NewPokerTableRepository(db DBTX) *PokerTableRepository {
	return &PokerTableRepository{db: db}
}

func (r *PokerTableRepository) Create(ctx context.Context, table domain.PokerTable) (domain.PokerTable, error) {
	query := `
		INSERT INTO poker_tables (name, small_blind, big_blind, min_buy_in, max_buy_in, max_players, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, small_blind, big_blind, min_buy_in, max_buy_in, max_players, status, created_at
	`

	var t domain.PokerTable
	err := r.db.QueryRow(ctx, query,
		table.Name, table.SmallBlind, table.BigBlind,
		table.MinBuyIn, table.MaxBuyIn, table.MaxPlayers, table.Status,
	).Scan(
		&t.ID, &t.Name, &t.SmallBlind, &t.BigBlind,
		&t.MinBuyIn, &t.MaxBuyIn, &t.MaxPlayers, &t.Status, &t.CreatedAt,
	)
	if err != nil {
		return t, fmt.Errorf("PokerTableRepository.Create: %w", err)
	}

	return t, nil
}

func (r *PokerTableRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.PokerTable, error) {
	query := `
		SELECT id, name, small_blind, big_blind, min_buy_in, max_buy_in, max_players, status, created_at
		FROM poker_tables
		WHERE id = $1
	`

	var t domain.PokerTable
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.Name, &t.SmallBlind, &t.BigBlind,
		&t.MinBuyIn, &t.MaxBuyIn, &t.MaxPlayers, &t.Status, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return t, domain.ErrTableNotFound
		}
		return t, fmt.Errorf("PokerTableRepository.FindByID: %w", err)
	}

	return t, nil
}

func (r *PokerTableRepository) FindActive(ctx context.Context) ([]domain.PokerTable, error) {
	query := `
		SELECT id, name, small_blind, big_blind, min_buy_in, max_buy_in, max_players, status, created_at
		FROM poker_tables
		WHERE status != 'closed'
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PokerTableRepository.FindActive: %w", err)
	}
	defer rows.Close()

	var tables []domain.PokerTable
	for rows.Next() {
		var t domain.PokerTable
		if err := rows.Scan(
			&t.ID, &t.Name, &t.SmallBlind, &t.BigBlind,
			&t.MinBuyIn, &t.MaxBuyIn, &t.MaxPlayers, &t.Status, &t.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("PokerTableRepository.FindActive scan: %w", err)
		}
		tables = append(tables, t)
	}

	return tables, rows.Err()
}

func (r *PokerTableRepository) Update(ctx context.Context, table domain.PokerTable) (domain.PokerTable, error) {
	query := `
		UPDATE poker_tables
		SET name = $1, small_blind = $2, big_blind = $3, min_buy_in = $4, max_buy_in = $5,
		    max_players = $6, status = $7
		WHERE id = $8
		RETURNING id, name, small_blind, big_blind, min_buy_in, max_buy_in, max_players, status, created_at
	`

	var t domain.PokerTable
	err := r.db.QueryRow(ctx, query,
		table.Name, table.SmallBlind, table.BigBlind,
		table.MinBuyIn, table.MaxBuyIn, table.MaxPlayers, table.Status, table.ID,
	).Scan(
		&t.ID, &t.Name, &t.SmallBlind, &t.BigBlind,
		&t.MinBuyIn, &t.MaxBuyIn, &t.MaxPlayers, &t.Status, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return t, domain.ErrTableNotFound
		}
		return t, fmt.Errorf("PokerTableRepository.Update: %w", err)
	}

	return t, nil
}
