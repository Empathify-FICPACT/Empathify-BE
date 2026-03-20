package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	UpdateOnboarding(ctx context.Context, id string, gender string, avatarID int16) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
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
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateOnboarding(ctx context.Context, id string, gender string, avatarID int16) error {
	query := `
		UPDATE users
		SET gender = $1, avatar_id = $2, updated_at = now()
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, gender, avatarID, id)
	return err
}