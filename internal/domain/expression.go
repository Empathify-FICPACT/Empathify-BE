package domain

import "time"

type ExpressionReference struct {
	ID       string
	Emotion  string
	ImageURL string
	IsActive bool
}

type ExpressionSession struct {
	ID               string
	UserID           string
	TotalExpressions int
	CorrectCount     int
	XPEarned         int
	Status           string
	StartedAt        time.Time
	CompletedAt      *time.Time
}

type ExpressionAttempt struct {
	ID              string
	SessionID       string
	ReferenceID     string
	UserPhotoURL    string
	SimilarityScore float64
	IsCorrect       bool
	AIFeedback      *string
	AttemptOrder    int
	CreatedAt       time.Time
}