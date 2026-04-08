package response

type AuthResponse struct {
	AccessToken string   `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User        UserData `json:"user"`
}

type UserData struct {
	ID       string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name     *string `json:"name" example:"Gantang Satria"`
	Gender   *string  `json:"gender" example:"male"`
	AvatarID *int16   `json:"avatar_id" example:"1"`
}