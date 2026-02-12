package postgres

import (
    "context"
    "fmt"
    "time"
    
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jokeoa/goigaming/internal/core/domain"
    "github.com/jokeoa/goigaming/internal/core/ports"
)

type rouletteRoundRepository struct {
    pool *pgxpool.Pool
}

func NewRouletteRoundRepository(pool *pgxpool.Pool) ports.RouletteRoundRepository {
    return &rouletteRoundRepository{pool: pool}
}

func (r *rouletteRoundRepository) Create(ctx context.Context, round *domain.RouletteRound) (*domain.RouletteRound, error) {
    query := `
        INSERT INTO roulette_rounds (
            id, table_id, round_number, seed_hash, betting_ends_at, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at
    `
    
    round.ID = uuid.New()
    round.CreatedAt = time.Now()
    
    err := r.pool.QueryRow(ctx, query,
        round.ID, round.TableID, round.RoundNumber, round.SeedHash,
        round.BettingEndsAt, round.CreatedAt,
    ).Scan(&round.ID, &round.CreatedAt)
    
    if err != nil {
        return nil, fmt.Errorf("failed to create roulette round: %w", err)
    }
    
    return round, nil
}

func (r *rouletteRoundRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.RouletteRound, error) {
    query := `
        SELECT id, table_id, round_number, result, result_color,
               seed_hash, seed_revealed, betting_ends_at, created_at, settled_at
        FROM roulette_rounds
        WHERE id = $1
    `
    
    var round domain.RouletteRound
    err := r.pool.QueryRow(ctx, query, id).Scan(
        &round.ID, &round.TableID, &round.RoundNumber, &round.Result, &round.ResultColor,
        &round.SeedHash, &round.SeedRevealed, &round.BettingEndsAt,
        &round.CreatedAt, &round.SettledAt,
    )
    
    if err != nil {
        return nil, fmt.Errorf("failed to get roulette round: %w", err)
    }
    
    return &round, nil
}

func (r *rouletteRoundRepository) GetCurrent(ctx context.Context, tableID uuid.UUID) (*domain.RouletteRound, error) {
    query := `
        SELECT id, table_id, round_number, result, result_color,
               seed_hash, seed_revealed, betting_ends_at, created_at, settled_at
        FROM roulette_rounds
        WHERE table_id = $1 AND settled_at IS NULL
        ORDER BY created_at DESC
        LIMIT 1
    `
    
    var round domain.RouletteRound
    err := r.pool.QueryRow(ctx, query, tableID).Scan(
        &round.ID, &round.TableID, &round.RoundNumber, &round.Result, &round.ResultColor,
        &round.SeedHash, &round.SeedRevealed, &round.BettingEndsAt,
        &round.CreatedAt, &round.SettledAt,
    )
    
    if err != nil {
        return nil, fmt.Errorf("failed to get current round: %w", err)
    }
    
    return &round, nil
}

func (r *rouletteRoundRepository) SetResult(ctx context.Context, id uuid.UUID, result int, color string) error {
    query := `
        UPDATE roulette_rounds
        SET result = $2, result_color = $3, settled_at = $4
        WHERE id = $1
    `
    
    _, err := r.pool.Exec(ctx, query, id, result, color, time.Now())
    if err != nil {
        return fmt.Errorf("failed to set roulette result: %w", err)
    }
    
    return nil
}

func (r *rouletteRoundRepository) GetHistory(ctx context.Context, tableID uuid.UUID, limit int) ([]*domain.RouletteRound, error) {
    query := `
        SELECT id, table_id, round_number, result, result_color,
               seed_hash, seed_revealed, betting_ends_at, created_at, settled_at
        FROM roulette_rounds
        WHERE table_id = $1 AND settled_at IS NOT NULL
        ORDER BY settled_at DESC
        LIMIT $2
    `
    
    rows, err := r.pool.Query(ctx, query, tableID, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to get history: %w", err)
    }
    defer rows.Close()
    
    rounds := make([]*domain.RouletteRound, 0)
    for rows.Next() {
        var round domain.RouletteRound
        err := rows.Scan(
            &round.ID, &round.TableID, &round.RoundNumber, &round.Result, &round.ResultColor,
            &round.SeedHash, &round.SeedRevealed, &round.BettingEndsAt,
            &round.CreatedAt, &round.SettledAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan round: %w", err)
        }
        rounds = append(rounds, &round)
    }
    
    return rounds, nil
}
