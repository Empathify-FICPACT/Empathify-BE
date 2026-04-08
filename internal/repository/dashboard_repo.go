package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository interface {
	GetConversationStats(ctx context.Context, userID string) (totalSessions int, err error)
	GetExpressionStats(ctx context.Context, userID string) (totalSessions int, avgScore float64, err error)
	GetEmotionStats(ctx context.Context, userID string) (totalSessions int, avgScore float64, err error)
	GetStoryStats(ctx context.Context, userID string) (totalSessions int, avgScore float64, err error)
	GetTodayMissionStats(ctx context.Context, userID string, date time.Time) (completed, total int, err error)
	GetLatestBadges(ctx context.Context, userID string, limit int) ([]LatestBadge, error)
	GetTotalUnlockedBadges(ctx context.Context, userID string) (int, error)
	GetTotalBadges(ctx context.Context) (int, error)
}

type LatestBadge struct {
	ID          string
	Name        string
	Description string
	IconURL     string
	EarnedAt    time.Time
}

type dashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetConversationStats(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM conversation_sessions
		WHERE user_id = $1 AND status = 'completed'
	`
	var total int
	err := r.db.QueryRow(ctx, query, userID).Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetExpressionStats(ctx context.Context, userID string) (int, float64, error) {
	query := `
		SELECT
			COUNT(*),
			COALESCE(AVG(
				CASE WHEN total_expressions > 0
				THEN correct_count::float / total_expressions::float
				ELSE 0 END
			), 0)
		FROM expression_sessions
		WHERE user_id = $1 AND status = 'completed'
	`
	var total int
	var avg float64
	err := r.db.QueryRow(ctx, query, userID).Scan(&total, &avg)
	return total, avg, err
}

func (r *dashboardRepository) GetEmotionStats(ctx context.Context, userID string) (int, float64, error) {
	query := `
		SELECT
			COUNT(*),
			COALESCE(AVG(
				CASE WHEN total_questions > 0
				THEN correct_count::float / total_questions::float
				ELSE 0 END
			), 0)
		FROM emotion_sessions
		WHERE user_id = $1 AND status = 'completed'
	`
	var total int
	var avg float64
	err := r.db.QueryRow(ctx, query, userID).Scan(&total, &avg)
	return total, avg, err
}

func (r *dashboardRepository) GetStoryStats(ctx context.Context, userID string) (int, float64, error) {
	query := `
		SELECT
			COUNT(*),
			COALESCE(AVG(
				CASE WHEN total_scenarios > 0
				THEN correct_count::float / total_scenarios::float
				ELSE 0 END
			), 0)
		FROM story_sessions
		WHERE user_id = $1 AND status = 'completed'
	`
	var total int
	var avg float64
	err := r.db.QueryRow(ctx, query, userID).Scan(&total, &avg)
	return total, avg, err
}

func (r *dashboardRepository) GetTodayMissionStats(ctx context.Context, userID string, date time.Time) (int, int, error) {
	query := `
		SELECT
			COUNT(CASE WHEN ump.is_completed = true THEN 1 END),
			COUNT(*)
		FROM daily_missions dm
		LEFT JOIN user_mission_progress ump
			ON ump.mission_id = dm.id AND ump.user_id = $1
		WHERE dm.mission_date = $2
	`
	var completed, total int
	err := r.db.QueryRow(ctx, query, userID, date).Scan(&completed, &total)
	return completed, total, err
}

func (r *dashboardRepository) GetLatestBadges(ctx context.Context, userID string, limit int) ([]LatestBadge, error) {
	query := `
		SELECT b.id, b.name, b.description, b.icon_url, ub.earned_at
		FROM user_badges ub
		JOIN badges b ON b.id = ub.badge_id
		WHERE ub.user_id = $1
		ORDER BY ub.earned_at DESC
		LIMIT $2
	`
	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []LatestBadge
	for rows.Next() {
		var b LatestBadge
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.IconURL, &b.EarnedAt); err != nil {
			return nil, err
		}
		badges = append(badges, b)
	}
	return badges, nil
}

func (r *dashboardRepository) GetTotalUnlockedBadges(ctx context.Context, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM user_badges WHERE user_id = $1`
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *dashboardRepository) GetTotalBadges(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM badges`
	var count int
	err := r.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}