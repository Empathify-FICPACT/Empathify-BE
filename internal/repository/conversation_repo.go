package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
)

type ConversationRepository interface {
	GetActiveTopics(ctx context.Context) ([]domain.ConversationTopic, error)
	GetTopicByID(ctx context.Context, id string) (*domain.ConversationTopic, error)
	CreateSession(ctx context.Context, session *domain.ConversationSession) error
	GetSessionByID(ctx context.Context, id string) (*domain.ConversationSession, error)
	GetSessionMessages(ctx context.Context, sessionID string) ([]domain.ConversationMessage, error)
	CreateMessage(ctx context.Context, msg *domain.ConversationMessage) error
	CompleteSession(ctx context.Context, sessionID string, xpEarned int) error
}

type conversationRepository struct {
	db *pgxpool.Pool
}

func NewConversationRepository(db *pgxpool.Pool) ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) GetActiveTopics(ctx context.Context) ([]domain.ConversationTopic, error) {
	query := `
		SELECT id, title, description, system_prompt, difficulty, is_active
		FROM conversation_topics
		WHERE is_active = true
		ORDER BY difficulty
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []domain.ConversationTopic
	for rows.Next() {
		var t domain.ConversationTopic
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.SystemPrompt, &t.Difficulty, &t.IsActive); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, nil
}

func (r *conversationRepository) GetTopicByID(ctx context.Context, id string) (*domain.ConversationTopic, error) {
	query := `
		SELECT id, title, description, system_prompt, difficulty, is_active
		FROM conversation_topics
		WHERE id = $1
	`
	var t domain.ConversationTopic
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.SystemPrompt, &t.Difficulty, &t.IsActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *conversationRepository) CreateSession(ctx context.Context, session *domain.ConversationSession) error {
	query := `
		INSERT INTO conversation_sessions (id, user_id, topic_id, status, xp_earned, started_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.TopicID,
		session.Status,
		session.XPEarned,
		session.StartedAt,
	)
	return err
}

func (r *conversationRepository) GetSessionByID(ctx context.Context, id string) (*domain.ConversationSession, error) {
	query := `
		SELECT id, user_id, topic_id, status, xp_earned, started_at, completed_at
		FROM conversation_sessions
		WHERE id = $1
	`
	var s domain.ConversationSession
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.UserID, &s.TopicID, &s.Status, &s.XPEarned, &s.StartedAt, &s.CompletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *conversationRepository) GetSessionMessages(ctx context.Context, sessionID string) ([]domain.ConversationMessage, error) {
	query := `
		SELECT id, session_id, role, content, created_at
		FROM conversation_messages
		WHERE session_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.ConversationMessage
	for rows.Next() {
		var m domain.ConversationMessage
		if err := rows.Scan(&m.ID, &m.SessionID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (r *conversationRepository) CreateMessage(ctx context.Context, msg *domain.ConversationMessage) error {
	query := `
		INSERT INTO conversation_messages (id, session_id, role, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query,
		msg.ID,
		msg.SessionID,
		msg.Role,
		msg.Content,
		msg.CreatedAt,
	)
	return err
}

func (r *conversationRepository) CompleteSession(ctx context.Context, sessionID string, xpEarned int) error {
	query := `
		UPDATE conversation_sessions
		SET status = 'completed', xp_earned = $1, completed_at = now()
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, xpEarned, sessionID)
	return err
}