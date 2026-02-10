package repository

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jokeoa/igaming/models"
)

type RouletteBetRepository struct {
    db *pgxpool.Pool
}

func NewRouletteBetRepository(db *pgxpool.Pool) *RouletteBetRepository {
    return &RouletteBetRepository{db: db}
}

func (r *RouletteBetRepository) Create(ctx context.Context, bet models.RouletteBet) (models.RouletteBet, error) {
    query := `
        INSERT INTO roulette_bets (round_id, user_id, bet_type, bet_value, amount, status)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, round_id, user_id, bet_type, bet_value, amount, payout, status, created_at
    `
    var b models.RouletteBet
    err := r.db.QueryRow(ctx, query, bet.RoundID, bet.UserID, bet.BetType, bet.BetValue, bet.Amount, bet.Status).Scan(
        &b.ID, &b.RoundID, &b.UserID, &b.BetType, &b.BetValue, &b.Amount, &b.Payout, &b.Status, &b.CreatedAt,
    )
    if err != nil {
        return b, fmt.Errorf("RouletteBetRepository.Create: %w", err)
    }
    return b, nil
}

func (r *RouletteBetRepository) FindByRoundId(ctx context.Context, roundID uuid.UUID) ([]models.RouletteBet, error) {
    query := `
        SELECT id, round_id, user_id, bet_type, bet_value, amount, payout, status, created_at
        FROM roulette_bets
        WHERE round_id = $1
        ORDER BY created_at ASC
    `
    rows, err := r.db.Query(ctx, query, roundID)
    if err != nil {
        return nil, fmt.Errorf("RouletteBetRepository.FindByRoundId: %w", err)
    }
    defer rows.Close()

    var bets []models.RouletteBet
    for rows.Next() {
        var b models.RouletteBet
        if err := rows.Scan(&b.ID, &b.RoundID, &b.UserID, &b.BetType, &b.BetValue, &b.Amount, &b.Payout, &b.Status, &b.CreatedAt); err != nil {
            return nil, fmt.Errorf("RouletteBetRepository.FindByRoundId scan: %w", err)
        }
        bets = append(bets, b)
    }
    return bets, nil
}

func (r *RouletteBetRepository) FindByUserId(ctx context.Context, userID uuid.UUID) ([]models.RouletteBet, error) {
    query := `
        SELECT id, round_id, user_id, bet_type, bet_value, amount, payout, status, created_at
        FROM roulette_bets
        WHERE user_id = $1
        ORDER BY created_at DESC
    `
    rows, err := r.db.Query(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("RouletteBetRepository.FindByUserId: %w", err)
    }
    defer rows.Close()

    var bets []models.RouletteBet
    for rows.Next() {
        var b models.RouletteBet
        if err := rows.Scan(&b.ID, &b.RoundID, &b.UserID, &b.BetType, &b.BetValue, &b.Amount, &b.Payout, &b.Status, &b.CreatedAt); err != nil {
            return nil, fmt.Errorf("RouletteBetRepository.FindByUserId scan: %w", err)
        }
        bets = append(bets, b)
    }
    return bets, nil
}

func (r *RouletteBetRepository) Update(ctx context.Context, bet models.RouletteBet) (models.RouletteBet, error) {
    query := `
        UPDATE roulette_bets
        SET payout = $1, status = $2
        WHERE id = $3
        RETURNING id, round_id, user_id, bet_type, bet_value, amount, payout, status, created_at
    `
    var b models.RouletteBet
    err := r.db.QueryRow(ctx, query, bet.Payout, bet.Status, bet.ID).Scan(
        &b.ID, &b.RoundID, &b.UserID, &b.BetType, &b.BetValue, &b.Amount, &b.Payout, &b.Status, &b.CreatedAt,
    )
    if err != nil {
        return b, fmt.Errorf("RouletteBetRepository.Update: %w", err)
    }
    return b, nil
}
