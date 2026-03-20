package request

type RegisterRequest struct {
	Name     *string `json:"name"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8"`
	Gender   string  `json:"gender" validate:"required,oneof=male female"`
	AvatarID int16   `json:"avatar_id" validate:"required,min=1,max=4"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}