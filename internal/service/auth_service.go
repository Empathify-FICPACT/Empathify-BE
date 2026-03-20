package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/repository"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/request"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/response"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/jwt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountNotVerified = errors.New("account not verified")
)

type AuthService interface {
	RegisterEmail(ctx context.Context, req *request.RegisterRequest) (*response.AuthResponse, error)
	LoginEmail(ctx context.Context, req *request.LoginRequest) (*response.AuthResponse, error)
	LoginGoogle(ctx context.Context, googleUser *domain.GoogleUser) (*response.AuthResponse, error)
}

type authService struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &authService{authRepo: authRepo}
}

func (s *authService) RegisterEmail(ctx context.Context, req *request.RegisterRequest) (*response.AuthResponse, error) {
	// cek email sudah dipakai belum
	exists, err := s.authRepo.IsEmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashedStr := string(hashedPassword)

	now := time.Now()
	userID := uuid.NewString()
	authID := uuid.NewString()

	// buat user
	user := &domain.User{
		ID:            userID,
		Name:          req.Name,
		Gender:        req.Gender,
		AvatarID:      req.AvatarID,
		CurrentStreak: 0,
		LongestStreak: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.authRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// buat user_auth
	userAuth := &domain.UserAuth{
		ID:           authID,
		UserID:       userID,
		Provider:     "email",
		ProviderUID:  nil,
		Email:        req.Email,
		PasswordHash: &hashedStr,
		IsVerified:   false,
		CreatedAt:    now,
	}

	if err := s.authRepo.CreateUserAuth(ctx, userAuth); err != nil {
		return nil, err
	}

	// generate token
	token, err := jwt.GenerateToken(userID)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		AccessToken: token,
		User: response.UserData{
			ID:       user.ID,
			Name:     user.Name,
			Gender:   user.Gender,
			AvatarID: user.AvatarID,
		},
	}, nil
}

func (s *authService) LoginEmail(ctx context.Context, req *request.LoginRequest) (*response.AuthResponse, error) {
	// cari auth by email
	userAuth, err := s.authRepo.FindAuthByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if userAuth == nil {
		return nil, ErrInvalidCredentials
	}

	// cek password
	if err := bcrypt.CompareHashAndPassword([]byte(*userAuth.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// ambil data user
	user, err := s.authRepo.FindUserByID(ctx, userAuth.UserID)
	if err != nil {
		return nil, err
	}

	// generate token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		AccessToken: token,
		User: response.UserData{
			ID:       user.ID,
			Name:     user.Name,
			Gender:   user.Gender,
			AvatarID: user.AvatarID,
		},
	}, nil
}

func (s *authService) LoginGoogle(ctx context.Context, googleUser *domain.GoogleUser) (*response.AuthResponse, error) {
	// cek apakah sudah pernah login google
	userAuth, err := s.authRepo.FindAuthByProviderUID(ctx, "google", googleUser.ID)
	if err != nil {
		return nil, err
	}

	var userID string

	if userAuth == nil {
		// user baru — buat user + user_auth
		now := time.Now()
		userID = uuid.NewString()
		authID := uuid.NewString()

		user := &domain.User{
			ID:            userID,
			Name:          &googleUser.Name,
			Gender:        googleUser.Gender,
			AvatarID:      1, // default avatar
			CurrentStreak: 0,
			LongestStreak: 0,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := s.authRepo.CreateUser(ctx, user); err != nil {
			return nil, err
		}

		providerUID := googleUser.ID
		newAuth := &domain.UserAuth{
			ID:           authID,
			UserID:       userID,
			Provider:     "google",
			ProviderUID:  &providerUID,
			Email:        googleUser.Email,
			PasswordHash: nil,
			IsVerified:   true, // google sudah verified
			CreatedAt:    now,
		}

		if err := s.authRepo.CreateUserAuth(ctx, newAuth); err != nil {
			return nil, err
		}
	} else {
		// user lama — langsung login
		userID = userAuth.UserID
	}

	// ambil data user
	user, err := s.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// generate token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		AccessToken: token,
		User: response.UserData{
			ID:       user.ID,
			Name:     user.Name,
			Gender:   user.Gender,
			AvatarID: user.AvatarID,
		},
	}, nil
}