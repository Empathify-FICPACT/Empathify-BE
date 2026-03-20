package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	CreateUserAuth(ctx context.Context, auth *domain.UserAuth) error
	FindUserByID(ctx context.Context, id string) (*domain.User, error)
	FindAuthByEmail(ctx context.Context, email string) (*domain.UserAuth, error)
	FindAuthByProviderUID(ctx context.Context, provider, providerUID string) (*domain.UserAuth, error)
	IsEmailExists(ctx context.Context, email string) (bool, error)
}

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, name, gender, avatar_id, current_streak, longest_streak, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Gender,
		user.AvatarID,
		user.CurrentStreak,
		user.LongestStreak,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *authRepository) CreateUserAuth(ctx context.Context, auth *domain.UserAuth) error {
	query := `
		INSERT INTO user_auth (id, user_id, provider, provider_uid, email, password_hash, is_verified, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		auth.ID,
		auth.UserID,
		auth.Provider,
		auth.ProviderUID,
		auth.Email,
		auth.PasswordHash,
		auth.IsVerified,
		auth.CreatedAt,
	)
	return err
}

func (r *authRepository) FindUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, name, gender, avatar_id, current_streak, longest_streak, last_active_date, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Gender,
		&user.AvatarID,
		&user.CurrentStreak,
		&user.LongestStreak,
		&user.LastActiveDate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *authRepository) FindAuthByEmail(ctx context.Context, email string) (*domain.UserAuth, error) {
	query := `
		SELECT id, user_id, provider, provider_uid, email, password_hash, is_verified, created_at
		FROM user_auth
		WHERE email = $1 AND provider = 'email'
	`
	auth := &domain.UserAuth{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&auth.ID,
		&auth.UserID,
		&auth.Provider,
		&auth.ProviderUID,
		&auth.Email,
		&auth.PasswordHash,
		&auth.IsVerified,
		&auth.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return auth, nil
}

func (r *authRepository) FindAuthByProviderUID(ctx context.Context, provider, providerUID string) (*domain.UserAuth, error) {
	query := `
		SELECT id, user_id, provider, provider_uid, email, password_hash, is_verified, created_at
		FROM user_auth
		WHERE provider = $1 AND provider_uid = $2
	`
	auth := &domain.UserAuth{}
	err := r.db.QueryRow(ctx, query, provider, providerUID).Scan(
		&auth.ID,
		&auth.UserID,
		&auth.Provider,
		&auth.ProviderUID,
		&auth.Email,
		&auth.PasswordHash,
		&auth.IsVerified,
		&auth.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return auth, nil
}

func (r *authRepository) IsEmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_auth WHERE email = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}