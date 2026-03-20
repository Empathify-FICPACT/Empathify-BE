package service

import (
	"context"
	"errors"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/repository"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/request"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/response"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrAlreadyOnboarded    = errors.New("user already onboarded")
)

type UserService interface {
	Onboarding(ctx context.Context, userID string, req *request.OnboardingRequest) (*response.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Onboarding(ctx context.Context, userID string, req *request.OnboardingRequest) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// cek apakah sudah pernah onboarding
	if user.Gender != nil && user.AvatarID != nil {
		return nil, ErrAlreadyOnboarded
	}

	if err := s.userRepo.UpdateOnboarding(ctx, userID, req.Gender, req.AvatarID); err != nil {
		return nil, err
	}

	// ambil data terbaru
	updated, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:            updated.ID,
		Name:          updated.Name,
		Gender:        updated.Gender,
		AvatarID:      updated.AvatarID,
		CurrentStreak: updated.CurrentStreak,
		LongestStreak: updated.LongestStreak,
	}, nil
}