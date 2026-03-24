package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
)

type StoryRepository interface {
	GetRandomScenarios(ctx context.Context, limit int) ([]domain.StoryScenario, error)
	GetScenarioByID(ctx context.Context, id string) (*domain.StoryScenario, error)
	CreateSession(ctx context.Context, session *domain.StorySession) error
	GetSessionByID(ctx context.Context, id string) (*domain.StorySession, error)
	CreateAnswer(ctx context.Context, answer *domain.StoryAnswer) error
	GetAnswersBySessionID(ctx context.Context, sessionID string) ([]domain.StoryAnswer, error)
	CompleteSession(ctx context.Context, sessionID string, correctCount, xpEarned int) error
}

type storyRepository struct {
	db *pgxpool.Pool
}

func NewStoryRepository(db *pgxpool.Pool) StoryRepository {
	return &storyRepository{db: db}
}

func (r *storyRepository) GetRandomScenarios(ctx context.Context, limit int) ([]domain.StoryScenario, error) {
	query := `
		SELECT id, situation, option_a, option_b, option_c, option_d, correct_answer, explanation, is_active
		FROM story_scenarios
		WHERE is_active = true
		ORDER BY RANDOM()
		LIMIT $1
	`
	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scenarios []domain.StoryScenario
	for rows.Next() {
		var s domain.StoryScenario
		if err := rows.Scan(
			&s.ID, &s.Situation,
			&s.OptionA, &s.OptionB, &s.OptionC, &s.OptionD,
			&s.CorrectAnswer, &s.Explanation, &s.IsActive,
		); err != nil {
			return nil, err
		}
		scenarios = append(scenarios, s)
	}
	return scenarios, nil
}

func (r *storyRepository) GetScenarioByID(ctx context.Context, id string) (*domain.StoryScenario, error) {
	query := `
		SELECT id, situation, option_a, option_b, option_c, option_d, correct_answer, explanation, is_active
		FROM story_scenarios
		WHERE id = $1
	`
	var s domain.StoryScenario
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.Situation,
		&s.OptionA, &s.OptionB, &s.OptionC, &s.OptionD,
		&s.CorrectAnswer, &s.Explanation, &s.IsActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *storyRepository) CreateSession(ctx context.Context, session *domain.StorySession) error {
	query := `
		INSERT INTO story_sessions (id, user_id, total_scenarios, correct_count, xp_earned, status, started_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.TotalScenarios,
		session.CorrectCount,
		session.XPEarned,
		session.Status,
		session.StartedAt,
	)
	return err
}

func (r *storyRepository) GetSessionByID(ctx context.Context, id string) (*domain.StorySession, error) {
	query := `
		SELECT id, user_id, total_scenarios, correct_count, xp_earned, status, started_at, completed_at
		FROM story_sessions
		WHERE id = $1
	`
	var s domain.StorySession
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.UserID, &s.TotalScenarios,
		&s.CorrectCount, &s.XPEarned, &s.Status,
		&s.StartedAt, &s.CompletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *storyRepository) CreateAnswer(ctx context.Context, answer *domain.StoryAnswer) error {
	query := `
		INSERT INTO story_answers (id, session_id, scenario_id, chosen_answer, is_correct, scenario_order, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		answer.ID,
		answer.SessionID,
		answer.ScenarioID,
		answer.ChosenAnswer,
		answer.IsCorrect,
		answer.ScenarioOrder,
		answer.CreatedAt,
	)
	return err
}

func (r *storyRepository) GetAnswersBySessionID(ctx context.Context, sessionID string) ([]domain.StoryAnswer, error) {
	query := `
		SELECT id, session_id, scenario_id, chosen_answer, is_correct, scenario_order, created_at
		FROM story_answers
		WHERE session_id = $1
		ORDER BY scenario_order ASC
	`
	rows, err := r.db.Query(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []domain.StoryAnswer
	for rows.Next() {
		var a domain.StoryAnswer
		if err := rows.Scan(
			&a.ID, &a.SessionID, &a.ScenarioID,
			&a.ChosenAnswer, &a.IsCorrect,
			&a.ScenarioOrder, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		answers = append(answers, a)
	}
	return answers, nil
}

func (r *storyRepository) CompleteSession(ctx context.Context, sessionID string, correctCount, xpEarned int) error {
	query := `
		UPDATE story_sessions
		SET status = 'completed', correct_count = $1, xp_earned = $2, completed_at = now()
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, correctCount, xpEarned, sessionID)
	return err
}