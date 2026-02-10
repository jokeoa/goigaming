package repository

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jokeoa/goigaming/models"
)

type RouletteRoundRepository struct {
    db *pgxpool.Pool
}

func NewRouletteRoundRepository(db *pgxpool.Pool) *RouletteRoundRepository {
    return &RouletteRoundRepository{db: db}
}

func (r *RouletteRoundRepository) Create(ctx context.Context, round models.RouletteRound) (models.RouletteRound, error) {
    query := `
        INSERT INTO roulette_rounds (table_id, round_number, seed_hash, betting_ends_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id, table_id, round_number, result, result_color, seed_hash, seed_revealed, betting_ends_at, created_at, settled_at
    `
    var rd models.RouletteRound
    err := r.db.QueryRow(ctx, query, round.TableID, round.RoundNumber, round.SeedHash, round.BettingEndsAt).Scan(
        &rd.ID, &rd.TableID, &rd.RoundNumber, &rd.Result, &rd.ResultColor,
        &rd.SeedHash, &rd.SeedRevealed, &rd.BettingEndsAt, &rd.CreatedAt, &rd.SettledAt,
    )
    if err != nil {
        return rd, fmt.Errorf("RouletteRoundRepository.Create: %w", err)
    }
    return rd, nil
}

func (r *RouletteRoundRepository) FindById(ctx context.Context, id uuid.UUID) (models.RouletteRound, error) {
    query := `
        SELECT id, table_id, round_number, result, result_color, seed_hash, seed_revealed, betting_ends_at, created_at, settled_at
        FROM roulette_rounds
        WHERE id = $1
    `
    var rd models.RouletteRound
    err := r.db.QueryRow(ctx, query, id).Scan(
        &rd.ID, &rd.TableID, &rd.RoundNumber, &rd.Result, &rd.ResultColor,
        &rd.SeedHash, &rd.SeedRevealed, &rd.BettingEndsAt, &rd.CreatedAt, &rd.SettledAt,
    )
    if err != nil {
        return rd, fmt.Errorf("RouletteRoundRepository.FindById: %w", err)
    }
    return rd, nil
}

func (r *RouletteRoundRepository) FindCurrent(ctx context.Context, tableID uuid.UUID) (models.RouletteRound, error) {
    query := `
        SELECT id, table_id, round_number, result, result_color, seed_hash, seed_revealed, betting_ends_at, created_at, settled_at
        FROM roulette_rounds
        WHERE table_id = $1 AND settled_at IS NULL
        ORDER BY created_at DESC
        LIMIT 1
    `
    var rd models.RouletteRound
    err := r.db.QueryRow(ctx, query, tableID).Scan(
        &rd.ID, &rd.TableID, &rd.RoundNumber, &rd.Result, &rd.ResultColor,
        &rd.SeedHash, &rd.SeedRevealed, &rd.BettingEndsAt, &rd.CreatedAt, &rd.SettledAt,
    )
    if err != nil {
        return rd, fmt.Errorf("RouletteRoundRepository.FindCurrent: %w", err)
    }
    return rd, nil
}

func (r *RouletteRoundRepository) Update(ctx context.Context, round models.RouletteRound) (models.RouletteRound, error) {
    query := `
        UPDATE roulette_rounds
        SET result = $1, result_color = $2, seed_revealed = $3, settled_at = $4
        WHERE id = $5
        RETURNING id, table_id, round_number, result, result_color, seed_hash, seed_revealed, betting_ends_at, created_at, settled_at
    `
    var rd models.RouletteRound
    err := r.db.QueryRow(ctx, query, round.Result, round.ResultColor, round.SeedRevealed, round.SettledAt, round.ID).Scan(
        &rd.ID, &rd.TableID, &rd.RoundNumber, &rd.Result, &rd.ResultColor,
        &rd.SeedHash, &rd.SeedRevealed, &rd.BettingEndsAt, &rd.CreatedAt, &rd.SettledAt,
    )
    if err != nil {
        return rd, fmt.Errorf("RouletteRoundRepository.Update: %w", err)
    }
    return rd, nil
}
