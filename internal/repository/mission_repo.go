package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
)

type MissionRepository interface {
	GetTodayMissions(ctx context.Context) ([]domain.DailyMission, error)
	GetUserMissionProgress(ctx context.Context, userID string, missionID string) (*domain.UserMissionProgress, error)
	GetUserTodayProgress(ctx context.Context, userID string, date time.Time) ([]domain.UserMissionProgress, error)
	CreateUserMissionProgress(ctx context.Context, progress *domain.UserMissionProgress) error
	UpdateMissionProgress(ctx context.Context, userID, missionID string, currentValue int, isCompleted bool) error
	EnsureTodayMissions(ctx context.Context) error
}

type missionRepository struct {
	db *pgxpool.Pool
}

func NewMissionRepository(db *pgxpool.Pool) MissionRepository {
	return &missionRepository{db: db}
}

func (r *missionRepository) EnsureTodayMissions(ctx context.Context) error {
	query := `
		INSERT INTO daily_missions (mission_date, mission_type, target_value, xp_reward)
		VALUES
			(CURRENT_DATE, 'collect_xp', 30, 10),
			(CURRENT_DATE, 'complete_exercises', 2, 10),
			(CURRENT_DATE, 'chat_pingo', 1, 10)
		ON CONFLICT (mission_date, mission_type) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query)
	return err
}

func (r *missionRepository) GetTodayMissions(ctx context.Context) ([]domain.DailyMission, error) {
	query := `
		SELECT id, mission_date, mission_type, target_value, xp_reward
		FROM daily_missions
		WHERE mission_date = CURRENT_DATE
		ORDER BY mission_type
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []domain.DailyMission
	for rows.Next() {
		var m domain.DailyMission
		if err := rows.Scan(&m.ID, &m.MissionDate, &m.MissionType, &m.TargetValue, &m.XPReward); err != nil {
			return nil, err
		}
		missions = append(missions, m)
	}
	return missions, nil
}

func (r *missionRepository) GetUserMissionProgress(ctx context.Context, userID, missionID string) (*domain.UserMissionProgress, error) {
	query := `
		SELECT id, user_id, mission_id, current_value, is_completed, completed_at
		FROM user_mission_progress
		WHERE user_id = $1 AND mission_id = $2
	`
	var p domain.UserMissionProgress
	err := r.db.QueryRow(ctx, query, userID, missionID).Scan(
		&p.ID, &p.UserID, &p.MissionID,
		&p.CurrentValue, &p.IsCompleted, &p.CompletedAt,
	)
	if err != nil {
		return nil, nil // belum ada progress
	}
	return &p, nil
}

func (r *missionRepository) GetUserTodayProgress(ctx context.Context, userID string, date time.Time) ([]domain.UserMissionProgress, error) {
	query := `
		SELECT ump.id, ump.user_id, ump.mission_id, ump.current_value, ump.is_completed, ump.completed_at
		FROM user_mission_progress ump
		JOIN daily_missions dm ON dm.id = ump.mission_id
		WHERE ump.user_id = $1 AND dm.mission_date = $2
	`
	rows, err := r.db.Query(ctx, query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progresses []domain.UserMissionProgress
	for rows.Next() {
		var p domain.UserMissionProgress
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.MissionID,
			&p.CurrentValue, &p.IsCompleted, &p.CompletedAt,
		); err != nil {
			return nil, err
		}
		progresses = append(progresses, p)
	}
	return progresses, nil
}

func (r *missionRepository) CreateUserMissionProgress(ctx context.Context, progress *domain.UserMissionProgress) error {
	query := `
		INSERT INTO user_mission_progress (id, user_id, mission_id, current_value, is_completed)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, mission_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query,
		progress.ID,
		progress.UserID,
		progress.MissionID,
		progress.CurrentValue,
		progress.IsCompleted,
	)
	return err
}

func (r *missionRepository) UpdateMissionProgress(ctx context.Context, userID, missionID string, currentValue int, isCompleted bool) error {
	query := `
		UPDATE user_mission_progress
		SET current_value = $1,
		    is_completed = $2,
		    completed_at = CASE WHEN $2 = true AND completed_at IS NULL THEN now() ELSE completed_at END
		WHERE user_id = $3 AND mission_id = $4
	`
	_, err := r.db.Exec(ctx, query, currentValue, isCompleted, userID, missionID)
	return err
}