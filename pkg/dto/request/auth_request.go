package request

type RegisterRequest struct {
	Name     *string `json:"name" example:"Gantang Satria"`
	Email    string  `json:"email" validate:"required,email" example:"satria@gmail.com"`
	Password string  `json:"password" validate:"required,min=8" example:"satria123"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"satria@gmail.com"`
	Password string `json:"password" validate:"required" example:"satria123"`
}