package response

type UserResponse struct {
	ID             string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name           *string `json:"name" example:"Gantang Satria"`
	Gender         *string `json:"gender" example:"male"`
	AvatarID       *int16  `json:"avatar_id" example:"1"`
	TotalXP       int     `json:"total_xp"`
	CurrentStreak  int     `json:"current_streak" example:"3"`
	LongestStreak  int     `json:"longest_streak" example:"7"`
}