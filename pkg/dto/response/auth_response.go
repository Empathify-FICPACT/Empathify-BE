package response

type AuthResponse struct {
	AccessToken string   `json:"access_token"`
	User        UserData `json:"user"`
}

type UserData struct {
	ID       string  `json:"id"`
	Name     *string `json:"name"`
	Gender   *string  `json:"gender"`
	AvatarID *int16   `json:"avatar_id"`
}