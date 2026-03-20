package domain

import "time"

type User struct {
	ID             string
	Name           *string
	Gender         string
	AvatarID       int16
	CurrentStreak  int
	LongestStreak  int
	LastActiveDate *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserAuth struct {
	ID           string
	UserID       string
	Provider     string
	ProviderUID  *string
	Email        string
	PasswordHash *string
	IsVerified   bool
	CreatedAt    time.Time
}