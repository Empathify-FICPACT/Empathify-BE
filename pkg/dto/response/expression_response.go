package response

import "time"

type ExpressionReferenceResponse struct {
	ID       string `json:"id"`
	Emotion  string `json:"emotion"`
	ImageURL string `json:"image_url"`
	Order    int    `json:"order"`
}

type ExpressionSessionResponse struct {
	ID               string                        `json:"id"`
	Status           string                        `json:"status"`
	TotalExpressions int                           `json:"total_expressions"`
	Expressions      []ExpressionReferenceResponse `json:"expressions"`
	StartedAt        time.Time                     `json:"started_at"`
}

type ExpressionAttemptResponse struct {
	AttemptID       string  `json:"attempt_id"`
	ReferenceID     string  `json:"reference_id"`
	Emotion         string  `json:"emotion"`
	SimilarityScore float64 `json:"similarity_score"`
	IsCorrect       bool    `json:"is_correct"`
	AIFeedback      string  `json:"ai_feedback"`
	AttemptOrder    int     `json:"attempt_order"`
}

type CompleteExpressionResponse struct {
	SessionID    string `json:"session_id"`
	CorrectCount int    `json:"correct_count"`
	Total        int    `json:"total"`
	XPEarned     int    `json:"xp_earned"`
	TotalXP      int    `json:"total_xp"`
}