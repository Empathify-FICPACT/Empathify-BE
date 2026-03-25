package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/repository"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/response"
)

type BadgeService interface {
	GetUserBadges(ctx context.Context, userID string) (*response.BadgeListResponse, error)
	CheckAndAwardBadges(ctx context.Context, userID string) error
}

type badgeService struct {
	badgeRepo   repository.BadgeRepository
	userRepo    repository.UserRepository
	exprRepo    repository.ExpressionRepository
	emotionRepo repository.EmotionRepository
	storyRepo   repository.StoryRepository
	convRepo    repository.ConversationRepository
}

func NewBadgeService(
	badgeRepo repository.BadgeRepository,
	userRepo repository.UserRepository,
	exprRepo repository.ExpressionRepository,
	emotionRepo repository.EmotionRepository,
	storyRepo repository.StoryRepository,
	convRepo repository.ConversationRepository,
) BadgeService {
	return &badgeService{
		badgeRepo:   badgeRepo,
		userRepo:    userRepo,
		exprRepo:    exprRepo,
		emotionRepo: emotionRepo,
		storyRepo:   storyRepo,
		convRepo:    convRepo,
	}
}

func (s *badgeService) GetUserBadges(ctx context.Context, userID string) (*response.BadgeListResponse, error) {
	allBadges, err := s.badgeRepo.GetAllBadges(ctx)
	if err != nil {
		return nil, err
	}

	userBadges, err := s.badgeRepo.GetUserBadges(ctx, userID)
	if err != nil {
		return nil, err
	}

	// map earned badges
	earnedMap := make(map[string]time.Time)
	for _, ub := range userBadges {
		earnedMap[ub.BadgeID] = ub.EarnedAt
	}

	var badgeResponses []response.BadgeResponse
	for _, b := range allBadges {
		earnedAt, isUnlocked := earnedMap[b.ID]
		br := response.BadgeResponse{
			ID:          b.ID,
			Name:        b.Name,
			Description: b.Description,
			IconURL:     b.IconURL,
			IsUnlocked:  isUnlocked,
		}
		if isUnlocked {
			br.EarnedAt = &earnedAt
		}
		badgeResponses = append(badgeResponses, br)
	}

	return &response.BadgeListResponse{
		Unlocked: len(userBadges),
		Total:    len(allBadges),
		Badges:   badgeResponses,
	}, nil
}

func (s *badgeService) CheckAndAwardBadges(ctx context.Context, userID string) error {
	allBadges, err := s.badgeRepo.GetAllBadges(ctx)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// hitung total exercise sessions yang completed
	totalExercises, err := s.countCompletedExercises(ctx, userID)
	if err != nil {
		return err
	}

	for _, badge := range allBadges {
		// cek apakah sudah punya badge ini
		has, err := s.badgeRepo.IsUserHasBadge(ctx, userID, badge.ID)
		if err != nil {
			continue
		}
		if has {
			continue
		}

		// cek apakah memenuhi syarat
		qualified := false
		switch badge.TriggerType {
		case "xp_total":
			qualified = user.TotalXP >= badge.TriggerValue
		case "exercise_count":
			qualified = totalExercises >= badge.TriggerValue
		}

		if qualified {
			userBadge := &domain.UserBadge{
				ID:       uuid.NewString(),
				UserID:   userID,
				BadgeID:  badge.ID,
				EarnedAt: time.Now(),
			}
			_ = s.badgeRepo.AwardBadge(ctx, userBadge)
		}
	}

	return nil
}

// func (s *badgeService) countCompletedExercises(ctx context.Context, userID string) (int, error) {
// 	// hitung dari semua fitur latihan
// 	count := 0

// 	// expression
// 	exprSessions, err := s.countCompletedExpressionSessions(ctx, userID)
// 	if err == nil {
// 		count += exprSessions
// 	}

// 	// emotion
// 	emotionSessions, err := s.countCompletedEmotionSessions(ctx, userID)
// 	if err == nil {
// 		count += emotionSessions
// 	}

// 	// story
// 	storySessions, err := s.countCompletedStorySessions(ctx, userID)
// 	if err == nil {
// 		count += storySessions
// 	}

// 	// conversation
// 	convSessions, err := s.countCompletedConvSessions(ctx, userID)
// 	if err == nil {
// 		count += convSessions
// 	}

// 	return count, nil
// }

func (s *badgeService) countCompletedExercises(ctx context.Context, userID string) (int, error) {
	count := 0

	n, _ := s.exprRepo.CountCompletedSessions(ctx, userID)
	count += n

	n, _ = s.emotionRepo.CountCompletedSessions(ctx, userID)
	count += n

	n, _ = s.storyRepo.CountCompletedSessions(ctx, userID)
	count += n

	n, _ = s.convRepo.CountCompletedSessions(ctx, userID)
	count += n

	return count, nil
}

// func (s *badgeService) countCompletedExpressionSessions(ctx context.Context, userID string) (int, error) {
// 	query := `SELECT COUNT(*) FROM expression_sessions WHERE user_id = $1 AND status = 'completed'`
// 	var count int
// 	err := s.exprRepo.(*expressionRepository).db.QueryRow(ctx, query, userID).Scan(&count)
// 	return count, err
// }

// func (s *badgeService) countCompletedEmotionSessions(ctx context.Context, userID string) (int, error) {
// 	query := `SELECT COUNT(*) FROM emotion_sessions WHERE user_id = $1 AND status = 'completed'`
// 	var count int
// 	err := s.emotionRepo.(*emotionRepository).db.QueryRow(ctx, query, userID).Scan(&count)
// 	return count, err
// }

// func (s *badgeService) countCompletedStorySessions(ctx context.Context, userID string) (int, error) {
// 	query := `SELECT COUNT(*) FROM story_sessions WHERE user_id = $1 AND status = 'completed'`
// 	var count int
// 	err := s.storyRepo.(*storyRepository).db.QueryRow(ctx, query, userID).Scan(&count)
// 	return count, err
// }

// func (s *badgeService) countCompletedConvSessions(ctx context.Context, userID string) (int, error) {
// 	query := `SELECT COUNT(*) FROM conversation_sessions WHERE user_id = $1 AND status = 'completed'`
// 	var count int
// 	err := s.convRepo.(*conversationRepository).db.QueryRow(ctx, query, userID).Scan(&count)
// 	return count, err
// }