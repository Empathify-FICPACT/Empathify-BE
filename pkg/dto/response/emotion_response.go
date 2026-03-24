package response

import "time"

type EmotionQuestionResponse struct {
	ID        string `json:"id"`
	Situation string `json:"situation"`
	OptionA   string `json:"option_a"`
	OptionB   string `json:"option_b"`
	OptionC   string `json:"option_c"`
	OptionD   string `json:"option_d"`
	Order     int    `json:"order"`
}

type EmotionSessionResponse struct {
	ID             string                    `json:"id"`
	Status         string                    `json:"status"`
	TotalQuestions int                       `json:"total_questions"`
	Questions      []EmotionQuestionResponse `json:"questions"`
	StartedAt      time.Time                 `json:"started_at"`
}

type EmotionAnswerResponse struct {
	QuestionID    string `json:"question_id"`
	ChosenAnswer  string `json:"chosen_answer"`
	CorrectAnswer string `json:"correct_answer"`
	IsCorrect     bool   `json:"is_correct"`
	QuestionOrder int    `json:"question_order"`
}

type CompleteEmotionResponse struct {
	SessionID    string `json:"session_id"`
	CorrectCount int    `json:"correct_count"`
	Total        int    `json:"total"`
	XPEarned     int    `json:"xp_earned"`
	TotalXP      int    `json:"total_xp"`
}