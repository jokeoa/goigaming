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

type sportEventRepository struct {
    pool *pgxpool.Pool
}

func NewSportEventRepository(pool *pgxpool.Pool) ports.SportEventRepository {
    return &sportEventRepository{pool: pool}
}

func (r *sportEventRepository) Create(ctx context.Context, event *domain.SportEvent) (*domain.SportEvent, error) {
    query := `
        INSERT INTO sport_events (
            id, sport_type, league, home_team, away_team, 
            home_odds, draw_odds, away_odds, event_time, 
            status, created_by, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
        )
        RETURNING id, created_at, updated_at
    `
    
    event.ID = uuid.New()
    event.CreatedAt = time.Now()
    event.UpdatedAt = time.Now()
    
    err := r.pool.QueryRow(ctx, query,
        event.ID, event.SportType, event.League, event.HomeTeam, event.AwayTeam,
        event.HomeOdds, event.DrawOdds, event.AwayOdds, event.EventTime,
        event.Status, event.CreatedBy, event.CreatedAt, event.UpdatedAt,
    ).Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)
    
    if err != nil {
        return nil, fmt.Errorf("failed to create sport event: %w", err)
    }
    
    return event, nil
}

func (r *sportEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.SportEvent, error) {
    query := `
        SELECT id, sport_type, league, home_team, away_team,
               home_odds, draw_odds, away_odds, event_time,
               status, home_score, away_score, created_by,
               created_at, updated_at, settled_at
        FROM sport_events
        WHERE id = $1
    `
    
    var event domain.SportEvent
    err := r.pool.QueryRow(ctx, query, id).Scan(
        &event.ID, &event.SportType, &event.League, &event.HomeTeam, &event.AwayTeam,
        &event.HomeOdds, &event.DrawOdds, &event.AwayOdds, &event.EventTime,
        &event.Status, &event.HomeScore, &event.AwayScore, &event.CreatedBy,
        &event.CreatedAt, &event.UpdatedAt, &event.SettledAt,
    )
    
    if err != nil {
        return nil, fmt.Errorf("failed to get sport event: %w", err)
    }
    
    return &event, nil
}

func (r *sportEventRepository) List(ctx context.Context, filter domain.SportEventFilter) ([]*domain.SportEvent, int64, error) {
    // Build dynamic query based on filter
    query := `
        SELECT id, sport_type, league, home_team, away_team,
               home_odds, draw_odds, away_odds, event_time,
               status, home_score, away_score, created_by,
               created_at, updated_at, settled_at
        FROM sport_events
        WHERE 1=1
    `
    countQuery := "SELECT COUNT(*) FROM sport_events WHERE 1=1"
    args := []interface{}{}
    argPos := 1
    
    if filter.SportType != nil {
        query += fmt.Sprintf(" AND sport_type = $%d", argPos)
        countQuery += fmt.Sprintf(" AND sport_type = $%d", argPos)
        args = append(args, *filter.SportType)
        argPos++
    }
    
    if filter.Status != nil {
        query += fmt.Sprintf(" AND status = $%d", argPos)
        countQuery += fmt.Sprintf(" AND status = $%d", argPos)
        args = append(args, *filter.Status)
        argPos++
    }
    
    if filter.FromDate != nil {
        query += fmt.Sprintf(" AND event_time >= $%d", argPos)
        args = append(args, *filter.FromDate)
        argPos++
    }
    
    if filter.ToDate != nil {
        query += fmt.Sprintf(" AND event_time <= $%d", argPos)
        args = append(args, *filter.ToDate)
        argPos++
    }
    
    query += " ORDER BY event_time ASC"
    
    if filter.Limit > 0 {
        query += fmt.Sprintf(" LIMIT $%d", argPos)
        args = append(args, filter.Limit)
        argPos++
    }
    
    if filter.Offset > 0 {
        query += fmt.Sprintf(" OFFSET $%d", argPos)
        args = append(args, filter.Offset)
    }
    
    // Get total count
    var total int64
    err := r.pool.QueryRow(ctx, countQuery, args[:argPos-1]...).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count events: %w", err)
    }
    
    // Get events
    rows, err := r.pool.Query(ctx, query, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to list events: %w", err)
    }
    defer rows.Close()
    
    events := make([]*domain.SportEvent, 0)
    for rows.Next() {
        var event domain.SportEvent
        err := rows.Scan(
            &event.ID, &event.SportType, &event.League, &event.HomeTeam, &event.AwayTeam,
            &event.HomeOdds, &event.DrawOdds, &event.AwayOdds, &event.EventTime,
            &event.Status, &event.HomeScore, &event.AwayScore, &event.CreatedBy,
            &event.CreatedAt, &event.UpdatedAt, &event.SettledAt,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan event: %w", err)
        }
        events = append(events, &event)
    }
    
    return events, total, nil
}

func (r *sportEventRepository) Update(ctx context.Context, event *domain.SportEvent) (*domain.SportEvent, error) {
    query := `
        UPDATE sport_events
        SET league = $2, home_team = $3, away_team = $4,
            home_odds = $5, draw_odds = $6, away_odds = $7,
            event_time = $8, updated_at = $9
        WHERE id = $1
        RETURNING updated_at
    `
    
    event.UpdatedAt = time.Now()
    err := r.pool.QueryRow(ctx, query,
        event.ID, event.League, event.HomeTeam, event.AwayTeam,
        event.HomeOdds, event.DrawOdds, event.AwayOdds,
        event.EventTime, event.UpdatedAt,
    ).Scan(&event.UpdatedAt)
    
    if err != nil {
        return nil, fmt.Errorf("failed to update event: %w", err)
    }
    
    return event, nil
}

func (r *sportEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
    query := "DELETE FROM sport_events WHERE id = $1 AND status = 'scheduled'"
    result, err := r.pool.Exec(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete event: %w", err)
    }
    
    if result.RowsAffected() == 0 {
        return fmt.Errorf("event not found or cannot be deleted")
    }
    
    return nil
}

func (r *sportEventRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
    query := "UPDATE sport_events SET status = $2, updated_at = $3 WHERE id = $1"
    _, err := r.pool.Exec(ctx, query, id, status, time.Now())
    return err
}

func (r *sportEventRepository) SetResult(ctx context.Context, id uuid.UUID, homeScore, awayScore int) error {
    query := `
        UPDATE sport_events
        SET home_score = $2, away_score = $3, status = 'finished', 
            settled_at = $4, updated_at = $5
        WHERE id = $1
    `
    now := time.Now()
    _, err := r.pool.Exec(ctx, query, id, homeScore, awayScore, now, now)
    return err
}
