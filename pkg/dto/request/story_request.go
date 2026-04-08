package request

type StartStoryRequest struct {
	TotalScenarios int `json:"total_scenarios"` // 5-8, default 5
}

type SubmitStoryAnswerRequest struct {
	ScenarioID   string `json:"scenario_id"`
	ChosenAnswer string `json:"chosen_answer"` // a, b, c, d
}