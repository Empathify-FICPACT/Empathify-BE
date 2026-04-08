package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
)

type BadgeRepository interface {
	GetAllBadges(ctx context.Context) ([]domain.Badge, error)
	GetUserBadges(ctx context.Context, userID string) ([]domain.UserBadge, error)
	AwardBadge(ctx context.Context, userBadge *domain.UserBadge) error
	IsUserHasBadge(ctx context.Context, userID, badgeID string) (bool, error)
}

type badgeRepository struct {
	db *pgxpool.Pool
}

func NewBadgeRepository(db *pgxpool.Pool) BadgeRepository {
	return &badgeRepository{db: db}
}

func (r *badgeRepository) GetAllBadges(ctx context.Context) ([]domain.Badge, error) {
	query := `
		SELECT id, name, description, icon_url, trigger_type, trigger_value
		FROM badges
		ORDER BY trigger_value ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []domain.Badge
	for rows.Next() {
		var b domain.Badge
		if err := rows.Scan(
			&b.ID, &b.Name, &b.Description,
			&b.IconURL, &b.TriggerType, &b.TriggerValue,
		); err != nil {
			return nil, err
		}
		badges = append(badges, b)
	}
	return badges, nil
}

func (r *badgeRepository) GetUserBadges(ctx context.Context, userID string) ([]domain.UserBadge, error) {
	query := `
		SELECT id, user_id, badge_id, earned_at
		FROM user_badges
		WHERE user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userBadges []domain.UserBadge
	for rows.Next() {
		var ub domain.UserBadge
		if err := rows.Scan(&ub.ID, &ub.UserID, &ub.BadgeID, &ub.EarnedAt); err != nil {
			return nil, err
		}
		userBadges = append(userBadges, ub)
	}
	return userBadges, nil
}

func (r *badgeRepository) AwardBadge(ctx context.Context, userBadge *domain.UserBadge) error {
	query := `
		INSERT INTO user_badges (id, user_id, badge_id, earned_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, badge_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query,
		userBadge.ID,
		userBadge.UserID,
		userBadge.BadgeID,
		userBadge.EarnedAt,
	)
	return err
}

func (r *badgeRepository) IsUserHasBadge(ctx context.Context, userID, badgeID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_badges WHERE user_id = $1 AND badge_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, badgeID).Scan(&exists)
	return exists, err
}