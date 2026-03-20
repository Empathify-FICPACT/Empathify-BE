package request

type OnboardingRequest struct {
	Gender   string `json:"gender" validate:"required,oneof=male female"`
	AvatarID int16  `json:"avatar_id" validate:"required,min=1,max=4"`
}