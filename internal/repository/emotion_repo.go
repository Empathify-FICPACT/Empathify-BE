package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
)

type EmotionRepository interface {
	GetRandomQuestions(ctx context.Context, limit int) ([]domain.EmotionQuestion, error)
	GetQuestionByID(ctx context.Context, id string) (*domain.EmotionQuestion, error)
	CreateSession(ctx context.Context, session *domain.EmotionSession) error
	GetSessionByID(ctx context.Context, id string) (*domain.EmotionSession, error)
	CreateAnswer(ctx context.Context, answer *domain.EmotionAnswer) error
	GetAnswersBySessionID(ctx context.Context, sessionID string) ([]domain.EmotionAnswer, error)
	CompleteSession(ctx context.Context, sessionID string, correctCount, xpEarned int) error
	CountCompletedSessions(ctx context.Context, userID string) (int, error)
}

type emotionRepository struct {
	db *pgxpool.Pool
}

func NewEmotionRepository(db *pgxpool.Pool) EmotionRepository {
	return &emotionRepository{db: db}
}

func (r *emotionRepository) GetRandomQuestions(ctx context.Context, limit int) ([]domain.EmotionQuestion, error) {
	query := `
		SELECT id, situation, option_a, option_b, option_c, option_d, correct_answer, is_active
		FROM emotion_questions
		WHERE is_active = true
		ORDER BY RANDOM()
		LIMIT $1
	`
	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []domain.EmotionQuestion
	for rows.Next() {
		var q domain.EmotionQuestion
		if err := rows.Scan(
			&q.ID, &q.Situation,
			&q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD,
			&q.CorrectAnswer, &q.IsActive,
		); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, nil
}

func (r *emotionRepository) GetQuestionByID(ctx context.Context, id string) (*domain.EmotionQuestion, error) {
	query := `
		SELECT id, situation, option_a, option_b, option_c, option_d, correct_answer, is_active
		FROM emotion_questions
		WHERE id = $1
	`
	var q domain.EmotionQuestion
	err := r.db.QueryRow(ctx, query, id).Scan(
		&q.ID, &q.Situation,
		&q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD,
		&q.CorrectAnswer, &q.IsActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &q, nil
}

func (r *emotionRepository) CreateSession(ctx context.Context, session *domain.EmotionSession) error {
	query := `
		INSERT INTO emotion_sessions (id, user_id, total_questions, correct_count, xp_earned, status, started_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.TotalQuestions,
		session.CorrectCount,
		session.XPEarned,
		session.Status,
		session.StartedAt,
	)
	return err
}

func (r *emotionRepository) GetSessionByID(ctx context.Context, id string) (*domain.EmotionSession, error) {
	query := `
		SELECT id, user_id, total_questions, correct_count, xp_earned, status, started_at, completed_at
		FROM emotion_sessions
		WHERE id = $1
	`
	var s domain.EmotionSession
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.UserID, &s.TotalQuestions,
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

func (r *emotionRepository) CreateAnswer(ctx context.Context, answer *domain.EmotionAnswer) error {
	query := `
		INSERT INTO emotion_answers (id, session_id, question_id, chosen_answer, is_correct, question_order, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		answer.ID,
		answer.SessionID,
		answer.QuestionID,
		answer.ChosenAnswer,
		answer.IsCorrect,
		answer.QuestionOrder,
		answer.CreatedAt,
	)
	return err
}

func (r *emotionRepository) GetAnswersBySessionID(ctx context.Context, sessionID string) ([]domain.EmotionAnswer, error) {
	query := `
		SELECT id, session_id, question_id, chosen_answer, is_correct, question_order, created_at
		FROM emotion_answers
		WHERE session_id = $1
		ORDER BY question_order ASC
	`
	rows, err := r.db.Query(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []domain.EmotionAnswer
	for rows.Next() {
		var a domain.EmotionAnswer
		if err := rows.Scan(
			&a.ID, &a.SessionID, &a.QuestionID,
			&a.ChosenAnswer, &a.IsCorrect,
			&a.QuestionOrder, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		answers = append(answers, a)
	}
	return answers, nil
}

func (r *emotionRepository) CompleteSession(ctx context.Context, sessionID string, correctCount, xpEarned int) error {
	query := `
		UPDATE emotion_sessions
		SET status = 'completed', correct_count = $1, xp_earned = $2, completed_at = now()
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, correctCount, xpEarned, sessionID)
	return err
}

func (r *emotionRepository) CountCompletedSessions(ctx context.Context, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM emotion_sessions WHERE user_id = $1 AND status = 'completed'`
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}