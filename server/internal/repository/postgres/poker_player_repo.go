package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jokeoa/goigaming/internal/core/domain"
	"github.com/shopspring/decimal"
)

type PokerPlayerRepository struct {
	db DBTX
}

func NewPokerPlayerRepository(db DBTX) *PokerPlayerRepository {
	return &PokerPlayerRepository{db: db}
}

func (r *PokerPlayerRepository) Create(ctx context.Context, player domain.PokerPlayer) (domain.PokerPlayer, error) {
	query := `
		INSERT INTO poker_players (table_id, user_id, username, stack, seat_number, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, table_id, user_id, username, stack, seat_number, status, joined_at
	`

	var p domain.PokerPlayer
	err := r.db.QueryRow(ctx, query,
		player.TableID, player.UserID, player.Username,
		player.Stack, player.SeatNumber, player.Status,
	).Scan(
		&p.ID, &p.TableID, &p.UserID, &p.Username,
		&p.Stack, &p.SeatNumber, &p.Status, &p.JoinedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "poker_players_table_id_seat_number_key" {
				return p, domain.ErrSeatTaken
			}
			return p, domain.ErrPlayerAlreadySeated
		}
		return p, fmt.Errorf("PokerPlayerRepository.Create: %w", err)
	}

	return p, nil
}

func (r *PokerPlayerRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.PokerPlayer, error) {
	query := `
		SELECT id, table_id, user_id, username, stack, seat_number, status, joined_at
		FROM poker_players
		WHERE id = $1
	`

	var p domain.PokerPlayer
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.TableID, &p.UserID, &p.Username,
		&p.Stack, &p.SeatNumber, &p.Status, &p.JoinedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p, domain.ErrPlayerNotFound
		}
		return p, fmt.Errorf("PokerPlayerRepository.FindByID: %w", err)
	}

	return p, nil
}

func (r *PokerPlayerRepository) FindByTableID(ctx context.Context, tableID uuid.UUID) ([]domain.PokerPlayer, error) {
	query := `
		SELECT id, table_id, user_id, username, stack, seat_number, status, joined_at
		FROM poker_players
		WHERE table_id = $1
		ORDER BY seat_number
	`

	rows, err := r.db.Query(ctx, query, tableID)
	if err != nil {
		return nil, fmt.Errorf("PokerPlayerRepository.FindByTableID: %w", err)
	}
	defer rows.Close()

	var players []domain.PokerPlayer
	for rows.Next() {
		var p domain.PokerPlayer
		if err := rows.Scan(
			&p.ID, &p.TableID, &p.UserID, &p.Username,
			&p.Stack, &p.SeatNumber, &p.Status, &p.JoinedAt,
		); err != nil {
			return nil, fmt.Errorf("PokerPlayerRepository.FindByTableID scan: %w", err)
		}
		players = append(players, p)
	}

	return players, rows.Err()
}

func (r *PokerPlayerRepository) FindByTableAndUser(ctx context.Context, tableID, userID uuid.UUID) (domain.PokerPlayer, error) {
	query := `
		SELECT id, table_id, user_id, username, stack, seat_number, status, joined_at
		FROM poker_players
		WHERE table_id = $1 AND user_id = $2
	`

	var p domain.PokerPlayer
	err := r.db.QueryRow(ctx, query, tableID, userID).Scan(
		&p.ID, &p.TableID, &p.UserID, &p.Username,
		&p.Stack, &p.SeatNumber, &p.Status, &p.JoinedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p, domain.ErrPlayerNotFound
		}
		return p, fmt.Errorf("PokerPlayerRepository.FindByTableAndUser: %w", err)
	}

	return p, nil
}

func (r *PokerPlayerRepository) UpdateStack(ctx context.Context, playerID uuid.UUID, stack decimal.Decimal) error {
	query := `UPDATE poker_players SET stack = $1 WHERE id = $2`

	tag, err := r.db.Exec(ctx, query, stack, playerID)
	if err != nil {
		return fmt.Errorf("PokerPlayerRepository.UpdateStack: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrPlayerNotFound
	}

	return nil
}

func (r *PokerPlayerRepository) UpdateStatus(ctx context.Context, playerID uuid.UUID, status domain.PlayerStatus) error {
	query := `UPDATE poker_players SET status = $1 WHERE id = $2`

	tag, err := r.db.Exec(ctx, query, status, playerID)
	if err != nil {
		return fmt.Errorf("PokerPlayerRepository.UpdateStatus: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrPlayerNotFound
	}

	return nil
}

func (r *PokerPlayerRepository) Delete(ctx context.Context, playerID uuid.UUID) error {
	query := `DELETE FROM poker_players WHERE id = $1`

	tag, err := r.db.Exec(ctx, query, playerID)
	if err != nil {
		return fmt.Errorf("PokerPlayerRepository.Delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrPlayerNotFound
	}

	return nil
}

func (r *PokerPlayerRepository) CountByTableID(ctx context.Context, tableID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM poker_players WHERE table_id = $1`

	var count int
	err := r.db.QueryRow(ctx, query, tableID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("PokerPlayerRepository.CountByTableID: %w", err)
	}

	return count, nil
}
