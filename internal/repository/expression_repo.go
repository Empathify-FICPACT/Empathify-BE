package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
)

type ExpressionRepository interface {
	GetRandomReferences(ctx context.Context, limit int) ([]domain.ExpressionReference, error)
	GetReferenceByID(ctx context.Context, id string) (*domain.ExpressionReference, error)
	CreateSession(ctx context.Context, session *domain.ExpressionSession) error
	GetSessionByID(ctx context.Context, id string) (*domain.ExpressionSession, error)
	CreateAttempt(ctx context.Context, attempt *domain.ExpressionAttempt) error
	GetAttemptsBySessionID(ctx context.Context, sessionID string) ([]domain.ExpressionAttempt, error)
	CompleteSession(ctx context.Context, sessionID string, correctCount, xpEarned int) error
}

type expressionRepository struct {
	db *pgxpool.Pool
}

func NewExpressionRepository(db *pgxpool.Pool) ExpressionRepository {
	return &expressionRepository{db: db}
}

func (r *expressionRepository) GetRandomReferences(ctx context.Context, limit int) ([]domain.ExpressionReference, error) {
	query := `
		SELECT id, emotion, image_url, is_active
		FROM expression_references
		WHERE is_active = true
		ORDER BY RANDOM()
		LIMIT $1
	`
	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []domain.ExpressionReference
	for rows.Next() {
		var ref domain.ExpressionReference
		if err := rows.Scan(&ref.ID, &ref.Emotion, &ref.ImageURL, &ref.IsActive); err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

func (r *expressionRepository) GetReferenceByID(ctx context.Context, id string) (*domain.ExpressionReference, error) {
	query := `
		SELECT id, emotion, image_url, is_active
		FROM expression_references
		WHERE id = $1
	`
	var ref domain.ExpressionReference
	err := r.db.QueryRow(ctx, query, id).Scan(
		&ref.ID, &ref.Emotion, &ref.ImageURL, &ref.IsActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ref, nil
}

func (r *expressionRepository) CreateSession(ctx context.Context, session *domain.ExpressionSession) error {
	query := `
		INSERT INTO expression_sessions (id, user_id, total_expressions, correct_count, xp_earned, status, started_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.TotalExpressions,
		session.CorrectCount,
		session.XPEarned,
		session.Status,
		session.StartedAt,
	)
	return err
}

func (r *expressionRepository) GetSessionByID(ctx context.Context, id string) (*domain.ExpressionSession, error) {
	query := `
		SELECT id, user_id, total_expressions, correct_count, xp_earned, status, started_at, completed_at
		FROM expression_sessions
		WHERE id = $1
	`
	var s domain.ExpressionSession
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.UserID, &s.TotalExpressions, &s.CorrectCount,
		&s.XPEarned, &s.Status, &s.StartedAt, &s.CompletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *expressionRepository) CreateAttempt(ctx context.Context, attempt *domain.ExpressionAttempt) error {
	query := `
		INSERT INTO expression_attempts (id, session_id, reference_id, user_photo_url, similarity_score, is_correct, ai_feedback, attempt_order, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		attempt.ID,
		attempt.SessionID,
		attempt.ReferenceID,
		attempt.UserPhotoURL,
		attempt.SimilarityScore,
		attempt.IsCorrect,
		attempt.AIFeedback,
		attempt.AttemptOrder,
		attempt.CreatedAt,
	)
	return err
}

func (r *expressionRepository) GetAttemptsBySessionID(ctx context.Context, sessionID string) ([]domain.ExpressionAttempt, error) {
	query := `
		SELECT id, session_id, reference_id, user_photo_url, similarity_score, is_correct, ai_feedback, attempt_order, created_at
		FROM expression_attempts
		WHERE session_id = $1
		ORDER BY attempt_order ASC
	`
	rows, err := r.db.Query(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []domain.ExpressionAttempt
	for rows.Next() {
		var a domain.ExpressionAttempt
		if err := rows.Scan(
			&a.ID, &a.SessionID, &a.ReferenceID, &a.UserPhotoURL,
			&a.SimilarityScore, &a.IsCorrect, &a.AIFeedback,
			&a.AttemptOrder, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}
	return attempts, nil
}

func (r *expressionRepository) CompleteSession(ctx context.Context, sessionID string, correctCount, xpEarned int) error {
	query := `
		UPDATE expression_sessions
		SET status = 'completed', correct_count = $1, xp_earned = $2, completed_at = now()
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, correctCount, xpEarned, sessionID)
	return err
}