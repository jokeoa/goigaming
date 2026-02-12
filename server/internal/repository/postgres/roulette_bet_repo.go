package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/domain"
)

type RouletteBetRepo struct {
	db DBTX
}

func NewRouletteBetRepo(db DBTX) *RouletteBetRepo {
	return &RouletteBetRepo{db: db}
}

func (r *RouletteBetRepo) Create(ctx context.Context, bet domain.RouletteBet) (domain.RouletteBet, error) {
	query := `
		INSERT INTO roulette_bets (round_id, user_id, bet_type, bet_value, amount, payout, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, round_id, user_id, bet_type, bet_value, amount, payout, status, created_at
	`

	var b domain.RouletteBet
	err := r.db.QueryRow(ctx, query,
		bet.RoundID, bet.UserID, bet.BetType, bet.BetValue,
		bet.Amount, bet.Payout, bet.Status,
	).Scan(
		&b.ID, &b.RoundID, &b.UserID, &b.BetType, &b.BetValue,
		&b.Amount, &b.Payout, &b.Status, &b.CreatedAt,
	)
	if err != nil {
		return b, fmt.Errorf("RouletteBetRepo.Create: %w", err)
	}

	return b, nil
}

func (r *RouletteBetRepo) FindByRoundID(ctx context.Context, roundID uuid.UUID) ([]domain.RouletteBet, error) {
	query := `
		SELECT id, round_id, user_id, bet_type, bet_value, amount, payout, status, created_at
		FROM roulette_bets
		WHERE round_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, roundID)
	if err != nil {
		return nil, fmt.Errorf("RouletteBetRepo.FindByRoundID: %w", err)
	}
	defer rows.Close()

	var bets []domain.RouletteBet
	for rows.Next() {
		var b domain.RouletteBet
		if err := rows.Scan(
			&b.ID, &b.RoundID, &b.UserID, &b.BetType, &b.BetValue,
			&b.Amount, &b.Payout, &b.Status, &b.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("RouletteBetRepo.FindByRoundID scan: %w", err)
		}
		bets = append(bets, b)
	}

	return bets, rows.Err()
}

func (r *RouletteBetRepo) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.RouletteBet, error) {
	query := `
		SELECT id, round_id, user_id, bet_type, bet_value, amount, payout, status, created_at
		FROM roulette_bets
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("RouletteBetRepo.FindByUserID: %w", err)
	}
	defer rows.Close()

	var bets []domain.RouletteBet
	for rows.Next() {
		var b domain.RouletteBet
		if err := rows.Scan(
			&b.ID, &b.RoundID, &b.UserID, &b.BetType, &b.BetValue,
			&b.Amount, &b.Payout, &b.Status, &b.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("RouletteBetRepo.FindByUserID scan: %w", err)
		}
		bets = append(bets, b)
	}

	return bets, rows.Err()
}
