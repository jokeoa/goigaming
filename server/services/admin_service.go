package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminService struct {
	db *pgxpool.Pool
}

func NewAdminService(db *pgxpool.Pool) *AdminService {
	return &AdminService{db: db}
}

type AnalyticsData struct {
	TotalUsers    int     `json:"total_users"`
	ActiveTables  int     `json:"active_tables"`
	TotalBets     int     `json:"total_bets"`
	TotalRevenue  float64 `json:"total_revenue"`
	TotalPayouts  float64 `json:"total_payouts"`
	NetRevenue    float64 `json:"net_revenue"`
	ActiveRounds  int     `json:"active_rounds"`
	SettledRounds int     `json:"settled_rounds"`
}

func (s *AdminService) GetAnalytics(ctx context.Context) (AnalyticsData, error) {
	var analytics AnalyticsData

	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&analytics.TotalUsers)
	if err != nil {
		return analytics, fmt.Errorf("failed to count users: %w", err)
	}

	err = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM roulette_tables WHERE status = 'active'`).Scan(&analytics.ActiveTables)
	if err != nil {
		return analytics, fmt.Errorf("failed to count active tables: %w", err)
	}

	err = s.db.QueryRow(ctx, `
		SELECT 
			COUNT(*), 
			COALESCE(SUM(amount), 0), 
			COALESCE(SUM(payout), 0)
		FROM roulette_bets
	`).Scan(&analytics.TotalBets, &analytics.TotalRevenue, &analytics.TotalPayouts)
	if err != nil {
		return analytics, fmt.Errorf("failed to get bet stats: %w", err)
	}

	analytics.NetRevenue = analytics.TotalRevenue - analytics.TotalPayouts

	err = s.db.QueryRow(ctx, `
		SELECT 
			COUNT(CASE WHEN settled_at IS NULL THEN 1 END) as active,
			COUNT(CASE WHEN settled_at IS NOT NULL THEN 1 END) as settled
		FROM roulette_rounds
	`).Scan(&analytics.ActiveRounds, &analytics.SettledRounds)
	if err != nil {
		return analytics, fmt.Errorf("failed to count rounds: %w", err)
	}

	return analytics, nil
}

type GameSession struct {
	ID         uuid.UUID `json:"id"`
	TableID    uuid.UUID `json:"table_id"`
	TableName  string    `json:"table_name"`
	RoundNum   int       `json:"round_number"`
	BetsCount  int       `json:"bets_count"`
	TotalStake float64   `json:"total_stake"`
}

func (s *AdminService) GetActiveSessions(ctx context.Context) ([]GameSession, error) {
	query := `
		SELECT 
			rr.id,
			rr.table_id,
			rt.name,
			rr.round_number,
			COUNT(rb.id) as bets_count,
			COALESCE(SUM(rb.amount), 0) as total_stake
		FROM roulette_rounds rr
		JOIN roulette_tables rt ON rr.table_id = rt.id
		LEFT JOIN roulette_bets rb ON rr.id = rb.round_id
		WHERE rr.settled_at IS NULL
		GROUP BY rr.id, rr.table_id, rt.name, rr.round_number
		ORDER BY rr.created_at DESC
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}
	defer rows.Close()

	var sessions []GameSession
	for rows.Next() {
		var session GameSession
		err := rows.Scan(
			&session.ID,
			&session.TableID,
			&session.TableName,
			&session.RoundNum,
			&session.BetsCount,
			&session.TotalStake,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *AdminService) ForceCloseSession(ctx context.Context, roundID uuid.UUID) error {
	query := `
		UPDATE roulette_rounds
		SET 
			result = 0,
			result_color = 'green',
			settled_at = NOW()
		WHERE id = $1 AND settled_at IS NULL
	`

	tag, err := s.db.Exec(ctx, query, roundID)
	if err != nil {
		return fmt.Errorf("failed to close session: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("session not found or already closed")
	}

	_, err = s.db.Exec(ctx, `
		UPDATE roulette_bets
		SET status = 'cancelled', payout = amount
		WHERE round_id = $1 AND status = 'pending'
	`, roundID)
	if err != nil {
		return fmt.Errorf("failed to cancel bets: %w", err)
	}

	return nil
}

type UserStats struct {
	UserID       uuid.UUID `json:"user_id"`
	TotalBets    int       `json:"total_bets"`
	TotalWagered float64   `json:"total_wagered"`
	TotalWon     float64   `json:"total_won"`
	NetProfit    float64   `json:"net_profit"`
}

func (s *AdminService) GetUserStats(ctx context.Context, userID uuid.UUID) (UserStats, error) {
	var stats UserStats
	stats.UserID = userID

	query := `
		SELECT 
			COUNT(*),
			COALESCE(SUM(amount), 0),
			COALESCE(SUM(payout), 0)
		FROM roulette_bets
		WHERE user_id = $1
	`

	err := s.db.QueryRow(ctx, query, userID).Scan(
		&stats.TotalBets,
		&stats.TotalWagered,
		&stats.TotalWon,
	)
	if err != nil {
		return stats, fmt.Errorf("failed to get user stats: %w", err)
	}

	stats.NetProfit = stats.TotalWon - stats.TotalWagered

	return stats, nil
}

type TopWinner struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	TotalWon  float64   `json:"total_won"`
	BetsCount int       `json:"bets_count"`
}

func (s *AdminService) GetTopWinners(ctx context.Context, limit int) ([]TopWinner, error) {
	query := `
		SELECT 
			u.id,
			u.username,
			COALESCE(SUM(rb.payout), 0) as total_won,
			COUNT(rb.id) as bets_count
		FROM users u
		JOIN roulette_bets rb ON u.id = rb.user_id
		WHERE rb.status = 'won'
		GROUP BY u.id, u.username
		ORDER BY total_won DESC
		LIMIT $1
	`

	rows, err := s.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top winners: %w", err)
	}
	defer rows.Close()

	var winners []TopWinner
	for rows.Next() {
		var winner TopWinner
		err := rows.Scan(&winner.UserID, &winner.Username, &winner.TotalWon, &winner.BetsCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan winner: %w", err)
		}
		winners = append(winners, winner)
	}

	return winners, nil
}
