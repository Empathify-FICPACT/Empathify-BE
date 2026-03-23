package response

import "time"

type TopicResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
}

type SessionResponse struct {
	ID        string    `json:"id"`
	TopicID   string    `json:"topic_id"`
	TopicTitle string   `json:"topic_title"`
	Status    string    `json:"status"`
	XPEarned  int       `json:"xp_earned"`
	StartedAt time.Time `json:"started_at"`
}

type MessageResponse struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type SendMessageResponse struct {
	UserMessage MessageResponse `json:"user_message"`
	AIMessage   MessageResponse `json:"ai_message"`
}

type CompleteSessionResponse struct {
	SessionID string `json:"session_id"`
	XPEarned  int    `json:"xp_earned"`
	TotalXP   int    `json:"total_xp"`
}