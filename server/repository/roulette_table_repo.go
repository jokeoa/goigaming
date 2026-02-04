package repository

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
    "goigaming/models"
)

type RouletteTableRepository struct {
    db *pgxpool.Pool
}

func NewRouletteTableRepository(db *pgxpool.Pool) *RouletteTableRepository {
    return &RouletteTableRepository{db: db}
}

func (r *RouletteTableRepository) Create(ctx context.Context, table models.RouletteTable) (models.RouletteTable, error) {
    query := `
        INSERT INTO roulette_tables (name, min_bet, max_bet, status)
        VALUES ($1, $2, $3, $4)
        RETURNING id, name, min_bet, max_bet, status, created_at
    `
    var t models.RouletteTable
    err := r.db.QueryRow(ctx, query, table.Name, table.MinBet, table.MaxBet, table.Status).Scan(
        &t.ID, &t.Name, &t.MinBet, &t.MaxBet, &t.Status, &t.CreatedAt,
    )
    if err != nil {
        return t, fmt.Errorf("RouletteTableRepository.Create: %w", err)
    }
    return t, nil
}

func (r *RouletteTableRepository) FindActive(ctx context.Context) ([]models.RouletteTable, error) {
    query := `
        SELECT id, name, min_bet, max_bet, status, created_at
        FROM roulette_tables
        WHERE status = 'active'
        ORDER BY created_at ASC
    `
    rows, err := r.db.Query(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("RouletteTableRepository.FindActive: %w", err)
    }
    defer rows.Close()

    var tables []models.RouletteTable
    for rows.Next() {
        var t models.RouletteTable
        if err := rows.Scan(&t.ID, &t.Name, &t.MinBet, &t.MaxBet, &t.Status, &t.CreatedAt); err != nil {
            return nil, fmt.Errorf("RouletteTableRepository.FindActive scan: %w", err)
        }
        tables = append(tables, t)
    }
    return tables, nil
}

func (r *RouletteTableRepository) Update(ctx context.Context, table models.RouletteTable) (models.RouletteTable, error) {
    query := `
        UPDATE roulette_tables
        SET name = $1, min_bet = $2, max_bet = $3, status = $4
        WHERE id = $5
        RETURNING id, name, min_bet, max_bet, status, created_at
    `
    var t models.RouletteTable
    err := r.db.QueryRow(ctx, query, table.Name, table.MinBet, table.MaxBet, table.Status, table.ID).Scan(
        &t.ID, &t.Name, &t.MinBet, &t.MaxBet, &t.Status, &t.CreatedAt,
    )
    if err != nil {
        return t, fmt.Errorf("RouletteTableRepository.Update: %w", err)
    }
    return t, nil
}
