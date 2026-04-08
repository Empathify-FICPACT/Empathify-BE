package response

import "time"

type MissionResponse struct {
	ID           string    `json:"id"`
	MissionType  string    `json:"mission_type"`
	Description  string    `json:"description"`
	TargetValue  int       `json:"target_value"`
	CurrentValue int       `json:"current_value"`
	IsCompleted  bool      `json:"is_completed"`
	XPReward     int       `json:"xp_reward"`
	CompletedAt  *time.Time `json:"completed_at"`
}

type DailyMissionsResponse struct {
	Date     string            `json:"date"`
	Missions []MissionResponse `json:"missions"`
}