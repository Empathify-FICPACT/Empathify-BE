package domain

import "time"

type DailyMission struct {
	ID          string
	MissionDate time.Time
	MissionType string
	TargetValue int
	XPReward    int
}

type UserMissionProgress struct {
	ID           string
	UserID       string
	MissionID    string
	CurrentValue int
	IsCompleted  bool
	CompletedAt  *time.Time
}