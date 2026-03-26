package response

type FeatureStats struct {
	TotalSessions int     `json:"total_sessions"`
	AvgScore      float64 `json:"avg_score"`
}

type DashboardResponse struct {
	User        DashboardUser        `json:"user"`
	XP          DashboardXP          `json:"xp"`
	Features    DashboardFeatures    `json:"features"`
	Missions    DashboardMissions    `json:"missions"`
	Badges      DashboardBadges      `json:"badges"`
}

type DashboardUser struct {
	ID       string  `json:"id"`
	Name     *string `json:"name"`
	Gender   *string `json:"gender"`
	AvatarID *int16  `json:"avatar_id"`
	Streak   int     `json:"streak"`
}

type DashboardXP struct {
	Total int `json:"total"`
}

type DashboardFeatures struct {
	Conversation FeatureStats `json:"conversation"`
	Expression   FeatureStats `json:"expression"`
	Emotion      FeatureStats `json:"emotion"`
	Story        FeatureStats `json:"story"`
}

type DashboardMissions struct {
	CompletedToday int `json:"completed_today"`
	TotalToday     int `json:"total_today"`
}

type DashboardBadges struct {
	Unlocked int             `json:"unlocked"`
	Total    int             `json:"total"`
	Latest   []BadgeResponse `json:"latest"` // 3 badge terakhir yang diraih
}