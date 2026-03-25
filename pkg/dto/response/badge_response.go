package response

import "time"

type BadgeResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IconURL     string     `json:"icon_url"`
	IsUnlocked  bool       `json:"is_unlocked"`
	EarnedAt    *time.Time `json:"earned_at"`
}

type BadgeListResponse struct {
	Unlocked int             `json:"unlocked"`
	Total    int             `json:"total"`
	Badges   []BadgeResponse `json:"badges"`
}