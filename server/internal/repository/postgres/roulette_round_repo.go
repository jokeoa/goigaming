package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jokeoa/goigaming/internal/core/domain"
)

type RouletteRoundRepo struct {
	db DBTX
}

func NewRouletteRoundRepo(db DBTX) *RouletteRoundRepo {
	return &RouletteRoundRepo{db: db}
}

func (r *RouletteRoundRepo) FindByID(ctx context.Context, id uuid.UUID) (domain.RouletteRound, error) {
	query := `
		SELECT id, table_id, round_number, result, result_color, seed_hash, seed_revealed,
		       betting_ends_at, created_at, settled_at
		FROM roulette_rounds
		WHERE id = $1
	`

	var round domain.RouletteRound
	err := r.db.QueryRow(ctx, query, id).Scan(
		&round.ID, &round.TableID, &round.RoundNumber,
		&round.Result, &round.ResultColor,
		&round.SeedHash, &round.SeedRevealed,
		&round.BettingEndsAt, &round.CreatedAt, &round.SettledAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return round, domain.ErrRoundNotFound
		}
		return round, fmt.Errorf("RouletteRoundRepo.FindByID: %w", err)
	}

	return round, nil
}

func (r *RouletteRoundRepo) FindCurrentByTableID(ctx context.Context, tableID uuid.UUID) (domain.RouletteRound, error) {
	query := `
		SELECT id, table_id, round_number, result, result_color, seed_hash, seed_revealed,
		       betting_ends_at, created_at, settled_at
		FROM roulette_rounds
		WHERE table_id = $1 AND settled_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`

	var round domain.RouletteRound
	err := r.db.QueryRow(ctx, query, tableID).Scan(
		&round.ID, &round.TableID, &round.RoundNumber,
		&round.Result, &round.ResultColor,
		&round.SeedHash, &round.SeedRevealed,
		&round.BettingEndsAt, &round.CreatedAt, &round.SettledAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return round, domain.ErrRoundNotFound
		}
		return round, fmt.Errorf("RouletteRoundRepo.FindCurrentByTableID: %w", err)
	}

	return round, nil
}

func (r *RouletteRoundRepo) FindSettledByTableID(ctx context.Context, tableID uuid.UUID, limit, offset int) ([]domain.RouletteRound, error) {
	query := `
		SELECT id, table_id, round_number, result, result_color, seed_hash, seed_revealed,
		       betting_ends_at, created_at, settled_at
		FROM roulette_rounds
		WHERE table_id = $1 AND settled_at IS NOT NULL
		ORDER BY settled_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, tableID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("RouletteRoundRepo.FindSettledByTableID: %w", err)
	}
	defer rows.Close()

	var rounds []domain.RouletteRound
	for rows.Next() {
		var round domain.RouletteRound
		if err := rows.Scan(
			&round.ID, &round.TableID, &round.RoundNumber,
			&round.Result, &round.ResultColor,
			&round.SeedHash, &round.SeedRevealed,
			&round.BettingEndsAt, &round.CreatedAt, &round.SettledAt,
		); err != nil {
			return nil, fmt.Errorf("RouletteRoundRepo.FindSettledByTableID scan: %w", err)
		}
		rounds = append(rounds, round)
	}

	return rounds, rows.Err()
}
