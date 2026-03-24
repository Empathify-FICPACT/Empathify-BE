package domain

import "time"

type StoryScenario struct {
	ID            string
	Situation     string
	OptionA       string
	OptionB       string
	OptionC       string
	OptionD       string
	CorrectAnswer string
	Explanation   string
	IsActive      bool
}

type StorySession struct {
	ID             string
	UserID         string
	TotalScenarios int
	CorrectCount   int
	XPEarned       int
	Status         string
	StartedAt      time.Time
	CompletedAt    *time.Time
}

type StoryAnswer struct {
	ID             string
	SessionID      string
	ScenarioID     string
	ChosenAnswer   string
	IsCorrect      bool
	ScenarioOrder  int
	CreatedAt      time.Time
}