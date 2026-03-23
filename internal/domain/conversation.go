package domain

import "time"

type ConversationTopic struct {
	ID           string
	Title        string
	Description  *string
	SystemPrompt string
	Difficulty   string
	IsActive     bool
}

type ConversationSession struct {
	ID          string
	UserID      string
	TopicID     string
	Status      string
	XPEarned    int
	StartedAt   time.Time
	CompletedAt *time.Time
}

type ConversationMessage struct {
	ID        string
	SessionID string
	Role      string
	Content   string
	CreatedAt time.Time
}