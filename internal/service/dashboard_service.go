package service

import (
	"context"
	"math"
	"time"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/repository"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/response"
)

type DashboardService interface {
	GetDashboard(ctx context.Context, userID string) (*response.DashboardResponse, error)
}

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
	userRepo      repository.UserRepository
}

func NewDashboardService(
	dashboardRepo repository.DashboardRepository,
	userRepo repository.UserRepository,
) DashboardService {
	return &dashboardService{
		dashboardRepo: dashboardRepo,
		userRepo:      userRepo,
	}
}

func (s *dashboardService) GetDashboard(ctx context.Context, userID string) (*response.DashboardResponse, error) {
	// ambil data user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// conversation stats
	convTotal, err := s.dashboardRepo.GetConversationStats(ctx, userID)
	if err != nil {
		convTotal = 0
	}

	// expression stats
	exprTotal, exprAvg, err := s.dashboardRepo.GetExpressionStats(ctx, userID)
	if err != nil {
		exprTotal, exprAvg = 0, 0
	}

	// emotion stats
	emotionTotal, emotionAvg, err := s.dashboardRepo.GetEmotionStats(ctx, userID)
	if err != nil {
		emotionTotal, emotionAvg = 0, 0
	}

	// story stats
	storyTotal, storyAvg, err := s.dashboardRepo.GetStoryStats(ctx, userID)
	if err != nil {
		storyTotal, storyAvg = 0, 0
	}

	// misi hari ini
	today := time.Now().Truncate(24 * time.Hour)
	missionCompleted, missionTotal, err := s.dashboardRepo.GetTodayMissionStats(ctx, userID, today)
	if err != nil {
		missionCompleted, missionTotal = 0, 0
	}

	// badges
	unlockedBadges, _ := s.dashboardRepo.GetTotalUnlockedBadges(ctx, userID)
	totalBadges, _ := s.dashboardRepo.GetTotalBadges(ctx)
	latestBadges, _ := s.dashboardRepo.GetLatestBadges(ctx, userID, 3)

	// convert latest badges ke response
	var badgeResponses []response.BadgeResponse
	for _, b := range latestBadges {
		earnedAt := b.EarnedAt
		badgeResponses = append(badgeResponses, response.BadgeResponse{
			ID:          b.ID,
			Name:        b.Name,
			Description: b.Description,
			IconURL:     b.IconURL,
			IsUnlocked:  true,
			EarnedAt:    &earnedAt,
		})
	}

	return &response.DashboardResponse{
		User: response.DashboardUser{
			ID:       user.ID,
			Name:     user.Name,
			Gender:   user.Gender,
			AvatarID: user.AvatarID,
			Streak:   user.CurrentStreak,
		},
		XP: response.DashboardXP{
			Total: user.TotalXP,
		},
		Features: response.DashboardFeatures{
			Conversation: response.FeatureStats{
				TotalSessions: convTotal,
				AvgScore:      0, // conversation tidak ada score
			},
			Expression: response.FeatureStats{
				TotalSessions: exprTotal,
				AvgScore:      roundTo2(exprAvg * 100), // convert ke persentase
			},
			Emotion: response.FeatureStats{
				TotalSessions: emotionTotal,
				AvgScore:      roundTo2(emotionAvg * 100),
			},
			Story: response.FeatureStats{
				TotalSessions: storyTotal,
				AvgScore:      roundTo2(storyAvg * 100),
			},
		},
		Missions: response.DashboardMissions{
			CompletedToday: missionCompleted,
			TotalToday:     missionTotal,
		},
		Badges: response.DashboardBadges{
			Unlocked: unlockedBadges,
			Total:    totalBadges,
			Latest:   badgeResponses,
		},
	}, nil
}

func roundTo2(val float64) float64 {
	return math.Round(val*100) / 100
}