package request

type StartEmotionRequest struct {
	TotalQuestions int `json:"total_questions"` // 5-8, default 5
}

type SubmitEmotionAnswerRequest struct {
	QuestionID   string `json:"question_id"`
	ChosenAnswer string `json:"chosen_answer"` // a, b, c, d
}