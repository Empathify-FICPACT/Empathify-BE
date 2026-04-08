package domain

import "time"

type EmotionQuestion struct {
	ID            string
	Situation     string
	OptionA       string
	OptionB       string
	OptionC       string
	OptionD       string
	CorrectAnswer string
	IsActive      bool
}

type EmotionSession struct {
	ID             string
	UserID         string
	TotalQuestions int
	CorrectCount   int
	XPEarned       int
	Status         string
	StartedAt      time.Time
	CompletedAt    *time.Time
}

type EmotionAnswer struct {
	ID            string
	SessionID     string
	QuestionID    string
	ChosenAnswer  string
	IsCorrect     bool
	QuestionOrder int
	CreatedAt     time.Time
}