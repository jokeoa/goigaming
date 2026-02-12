package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jokeoa/goigaming/internal/core/domain"
)

type PokerHandRepository struct {
	db DBTX
}

func NewPokerHandRepository(db DBTX) *PokerHandRepository {
	return &PokerHandRepository{db: db}
}

func (r *PokerHandRepository) Create(ctx context.Context, hand domain.PokerHand) (domain.PokerHand, error) {
	query := `
		INSERT INTO poker_hands (table_id, hand_number, pot, community_cards, stage)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, table_id, hand_number, pot, community_cards, stage, winner_id, started_at, ended_at
	`

	var h domain.PokerHand
	err := r.db.QueryRow(ctx, query,
		hand.TableID, hand.HandNumber, hand.Pot, hand.CommunityCards, hand.Stage,
	).Scan(
		&h.ID, &h.TableID, &h.HandNumber, &h.Pot, &h.CommunityCards,
		&h.Stage, &h.WinnerID, &h.StartedAt, &h.EndedAt,
	)
	if err != nil {
		return h, fmt.Errorf("PokerHandRepository.Create: %w", err)
	}

	return h, nil
}

func (r *PokerHandRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.PokerHand, error) {
	query := `
		SELECT id, table_id, hand_number, pot, community_cards, stage, winner_id, started_at, ended_at
		FROM poker_hands
		WHERE id = $1
	`

	var h domain.PokerHand
	err := r.db.QueryRow(ctx, query, id).Scan(
		&h.ID, &h.TableID, &h.HandNumber, &h.Pot, &h.CommunityCards,
		&h.Stage, &h.WinnerID, &h.StartedAt, &h.EndedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return h, domain.ErrHandNotFound
		}
		return h, fmt.Errorf("PokerHandRepository.FindByID: %w", err)
	}

	return h, nil
}

func (r *PokerHandRepository) Update(ctx context.Context, hand domain.PokerHand) (domain.PokerHand, error) {
	query := `
		UPDATE poker_hands
		SET pot = $1, community_cards = $2, stage = $3, winner_id = $4, ended_at = $5
		WHERE id = $6
		RETURNING id, table_id, hand_number, pot, community_cards, stage, winner_id, started_at, ended_at
	`

	var h domain.PokerHand
	err := r.db.QueryRow(ctx, query,
		hand.Pot, hand.CommunityCards, hand.Stage, hand.WinnerID, hand.EndedAt, hand.ID,
	).Scan(
		&h.ID, &h.TableID, &h.HandNumber, &h.Pot, &h.CommunityCards,
		&h.Stage, &h.WinnerID, &h.StartedAt, &h.EndedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return h, domain.ErrHandNotFound
		}
		return h, fmt.Errorf("PokerHandRepository.Update: %w", err)
	}

	return h, nil
}

func (r *PokerHandRepository) FindLatestByTableID(ctx context.Context, tableID uuid.UUID) (domain.PokerHand, error) {
	query := `
		SELECT id, table_id, hand_number, pot, community_cards, stage, winner_id, started_at, ended_at
		FROM poker_hands
		WHERE table_id = $1
		ORDER BY hand_number DESC
		LIMIT 1
	`

	var h domain.PokerHand
	err := r.db.QueryRow(ctx, query, tableID).Scan(
		&h.ID, &h.TableID, &h.HandNumber, &h.Pot, &h.CommunityCards,
		&h.Stage, &h.WinnerID, &h.StartedAt, &h.EndedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return h, domain.ErrHandNotFound
		}
		return h, fmt.Errorf("PokerHandRepository.FindLatestByTableID: %w", err)
	}

	return h, nil
}

func (r *PokerHandRepository) CreateHandPlayer(ctx context.Context, hp domain.PokerHandPlayer) (domain.PokerHandPlayer, error) {
	query := `
		INSERT INTO poker_hand_players (hand_id, player_id, hole_cards, bet_amount, last_action, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, hand_id, player_id, hole_cards, bet_amount, last_action, is_active
	`

	var result domain.PokerHandPlayer
	err := r.db.QueryRow(ctx, query,
		hp.HandID, hp.PlayerID, hp.HoleCards, hp.BetAmount, hp.LastAction, hp.IsActive,
	).Scan(
		&result.ID, &result.HandID, &result.PlayerID,
		&result.HoleCards, &result.BetAmount, &result.LastAction, &result.IsActive,
	)
	if err != nil {
		return result, fmt.Errorf("PokerHandRepository.CreateHandPlayer: %w", err)
	}

	return result, nil
}

func (r *PokerHandRepository) FindHandPlayers(ctx context.Context, handID uuid.UUID) ([]domain.PokerHandPlayer, error) {
	query := `
		SELECT id, hand_id, player_id, hole_cards, bet_amount, last_action, is_active
		FROM poker_hand_players
		WHERE hand_id = $1
	`

	rows, err := r.db.Query(ctx, query, handID)
	if err != nil {
		return nil, fmt.Errorf("PokerHandRepository.FindHandPlayers: %w", err)
	}
	defer rows.Close()

	var players []domain.PokerHandPlayer
	for rows.Next() {
		var hp domain.PokerHandPlayer
		if err := rows.Scan(
			&hp.ID, &hp.HandID, &hp.PlayerID,
			&hp.HoleCards, &hp.BetAmount, &hp.LastAction, &hp.IsActive,
		); err != nil {
			return nil, fmt.Errorf("PokerHandRepository.FindHandPlayers scan: %w", err)
		}
		players = append(players, hp)
	}

	return players, rows.Err()
}

func (r *PokerHandRepository) UpdateHandPlayer(ctx context.Context, hp domain.PokerHandPlayer) error {
	query := `
		UPDATE poker_hand_players
		SET hole_cards = $1, bet_amount = $2, last_action = $3, is_active = $4
		WHERE id = $5
	`

	tag, err := r.db.Exec(ctx, query, hp.HoleCards, hp.BetAmount, hp.LastAction, hp.IsActive, hp.ID)
	if err != nil {
		return fmt.Errorf("PokerHandRepository.UpdateHandPlayer: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrPlayerNotFound
	}

	return nil
}

func (r *PokerHandRepository) CreateAction(ctx context.Context, action domain.PokerAction) (domain.PokerAction, error) {
	query := `
		INSERT INTO poker_actions (hand_id, player_id, action, amount, stage, action_order)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, hand_id, player_id, action, amount, stage, action_order, created_at
	`

	var a domain.PokerAction
	err := r.db.QueryRow(ctx, query,
		action.HandID, action.PlayerID, action.Action,
		action.Amount, action.Stage, action.ActionOrder,
	).Scan(
		&a.ID, &a.HandID, &a.PlayerID, &a.Action,
		&a.Amount, &a.Stage, &a.ActionOrder, &a.CreatedAt,
	)
	if err != nil {
		return a, fmt.Errorf("PokerHandRepository.CreateAction: %w", err)
	}

	return a, nil
}

func (r *PokerHandRepository) FindActionsByHandID(ctx context.Context, handID uuid.UUID) ([]domain.PokerAction, error) {
	query := `
		SELECT id, hand_id, player_id, action, amount, stage, action_order, created_at
		FROM poker_actions
		WHERE hand_id = $1
		ORDER BY action_order
	`

	rows, err := r.db.Query(ctx, query, handID)
	if err != nil {
		return nil, fmt.Errorf("PokerHandRepository.FindActionsByHandID: %w", err)
	}
	defer rows.Close()

	var actions []domain.PokerAction
	for rows.Next() {
		var a domain.PokerAction
		if err := rows.Scan(
			&a.ID, &a.HandID, &a.PlayerID, &a.Action,
			&a.Amount, &a.Stage, &a.ActionOrder, &a.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("PokerHandRepository.FindActionsByHandID scan: %w", err)
		}
		actions = append(actions, a)
	}

	return actions, rows.Err()
}
