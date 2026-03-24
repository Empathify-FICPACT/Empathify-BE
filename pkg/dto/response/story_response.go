package response

import "time"

type StoryScenarioResponse struct {
	ID        string `json:"id"`
	Situation string `json:"situation"`
	OptionA   string `json:"option_a"`
	OptionB   string `json:"option_b"`
	OptionC   string `json:"option_c"`
	OptionD   string `json:"option_d"`
	Order     int    `json:"order"`
}

type StorySessionResponse struct {
	ID             string                  `json:"id"`
	Status         string                  `json:"status"`
	TotalScenarios int                     `json:"total_scenarios"`
	Scenarios      []StoryScenarioResponse `json:"scenarios"`
	StartedAt      time.Time               `json:"started_at"`
}

type StoryAnswerResponse struct {
	ScenarioID    string `json:"scenario_id"`
	ChosenAnswer  string `json:"chosen_answer"`
	CorrectAnswer string `json:"correct_answer"`
	IsCorrect     bool   `json:"is_correct"`
	Explanation   string `json:"explanation"`
	ScenarioOrder int    `json:"scenario_order"`
}

type CompleteStoryResponse struct {
	SessionID    string `json:"session_id"`
	CorrectCount int    `json:"correct_count"`
	Total        int    `json:"total"`
	XPEarned     int    `json:"xp_earned"`
	TotalXP      int    `json:"total_xp"`
}