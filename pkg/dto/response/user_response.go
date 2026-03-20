package response

type UserResponse struct {
	ID             string  `json:"id"`
	Name           *string `json:"name"`
	Gender         *string `json:"gender"`
	AvatarID       *int16  `json:"avatar_id"`
	CurrentStreak  int     `json:"current_streak"`
	LongestStreak  int     `json:"longest_streak"`
}