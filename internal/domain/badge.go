package domain

import "time"

type Badge struct {
	ID           string
	Name         string
	Description  string
	IconURL      string
	TriggerType  string
	TriggerValue int
}

type UserBadge struct {
	ID       string
	UserID   string
	BadgeID  string
	EarnedAt time.Time
}