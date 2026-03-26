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
	ErrInvalidGender    = errors.New("gender harus male atau female")
	ErrInvalidAvatarID  = errors.New("avatar_id harus antara 1 sampai 4")
)

type UserService interface {
	Onboarding(ctx context.Context, userID string, req *request.OnboardingRequest) (*response.UserResponse, error)
	EditProfile(ctx context.Context, userID string, req *request.EditProfileRequest) (*response.UserResponse, error)
	GetProfile(ctx context.Context, userID string) (*response.UserResponse, error)
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

func (s *userService) GetProfile(ctx context.Context, userID string) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &response.UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Gender:        user.Gender,
		AvatarID:      user.AvatarID,
		TotalXP:       user.TotalXP,
		CurrentStreak: user.CurrentStreak,
		LongestStreak: user.LongestStreak,
	}, nil
}

func (s *userService) EditProfile(ctx context.Context, userID string, req *request.EditProfileRequest) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	// validasi gender kalau diisi
	if req.Gender != nil {
		if *req.Gender != "male" && *req.Gender != "female" {
			return nil, ErrInvalidGender
		}
	}

	// validasi avatar_id kalau diisi
	if req.AvatarID != nil {
		if *req.AvatarID < 1 || *req.AvatarID > 4 {
			return nil, ErrInvalidAvatarID
		}
	}

	if err := s.userRepo.UpdateProfile(ctx, userID, req.Name, req.Gender, req.AvatarID); err != nil {
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
		TotalXP:       updated.TotalXP,
		CurrentStreak: updated.CurrentStreak,
		LongestStreak: updated.LongestStreak,
	}, nil
}