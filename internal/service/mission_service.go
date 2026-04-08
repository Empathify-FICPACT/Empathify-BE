package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/repository"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/response"
)

type MissionService interface {
	GetTodayMissions(ctx context.Context, userID string) (*response.DailyMissionsResponse, error)
	UpdateProgressAfterSession(ctx context.Context, userID, featureType string, xpEarned int) error
}

type missionService struct {
	missionRepo repository.MissionRepository
	userRepo    repository.UserRepository
}

func NewMissionService(
	missionRepo repository.MissionRepository,
	userRepo repository.UserRepository,
) MissionService {
	return &missionService{
		missionRepo: missionRepo,
		userRepo:    userRepo,
	}
}

func (s *missionService) GetTodayMissions(ctx context.Context, userID string) (*response.DailyMissionsResponse, error) {
	// pastikan misi hari ini sudah ada
	if err := s.missionRepo.EnsureTodayMissions(ctx); err != nil {
		return nil, err
	}

	missions, err := s.missionRepo.GetTodayMissions(ctx)
	if err != nil {
		return nil, err
	}

	// ambil progress user
	today := time.Now().Truncate(24 * time.Hour)
	progresses, err := s.missionRepo.GetUserTodayProgress(ctx, userID, today)
	if err != nil {
		return nil, err
	}

	// map progress by mission_id
	progressMap := make(map[string]domain.UserMissionProgress)
	for _, p := range progresses {
		progressMap[p.MissionID] = p
	}

	// ensure user punya progress row untuk tiap misi
	for _, m := range missions {
		if _, exists := progressMap[m.ID]; !exists {
			progress := &domain.UserMissionProgress{
				ID:           uuid.NewString(),
				UserID:       userID,
				MissionID:    m.ID,
				CurrentValue: 0,
				IsCompleted:  false,
			}
			_ = s.missionRepo.CreateUserMissionProgress(ctx, progress)
			progressMap[m.ID] = *progress
		}
	}

	var missionResponses []response.MissionResponse
	for _, m := range missions {
		p := progressMap[m.ID]
		missionResponses = append(missionResponses, response.MissionResponse{
			ID:           m.ID,
			MissionType:  m.MissionType,
			Description:  missionDescription(m.MissionType, m.TargetValue),
			TargetValue:  m.TargetValue,
			CurrentValue: p.CurrentValue,
			IsCompleted:  p.IsCompleted,
			XPReward:     m.XPReward,
			CompletedAt:  p.CompletedAt,
		})
	}

	return &response.DailyMissionsResponse{
		Date:     time.Now().Format("2006-01-02"),
		Missions: missionResponses,
	}, nil
}

func (s *missionService) UpdateProgressAfterSession(ctx context.Context, userID, featureType string, xpEarned int) error {
	if err := s.missionRepo.EnsureTodayMissions(ctx); err != nil {
		return err
	}

	missions, err := s.missionRepo.GetTodayMissions(ctx)
	if err != nil {
		return err
	}

	today := time.Now().Truncate(24 * time.Hour)
	progresses, err := s.missionRepo.GetUserTodayProgress(ctx, userID, today)
	if err != nil {
		return err
	}

	progressMap := make(map[string]domain.UserMissionProgress)
	for _, p := range progresses {
		progressMap[p.MissionID] = p
	}

	for _, m := range missions {
		p, exists := progressMap[m.ID]
		if !exists {
			newProgress := &domain.UserMissionProgress{
				ID:           uuid.NewString(),
				UserID:       userID,
				MissionID:    m.ID,
				CurrentValue: 0,
				IsCompleted:  false,
			}
			_ = s.missionRepo.CreateUserMissionProgress(ctx, newProgress)
			p = *newProgress
		}

		if p.IsCompleted {
			continue
		}

		newValue := p.CurrentValue
		completed := false

		switch m.MissionType {
		case "collect_xp":
			newValue += xpEarned
			if newValue >= m.TargetValue {
				completed = true
				newValue = m.TargetValue
			}

		case "complete_exercises":
			// semua fitur kecuali conversation dihitung sebagai exercise
			if featureType == "expression" || featureType == "emotion" || featureType == "story" {
				newValue++
				if newValue >= m.TargetValue {
					completed = true
				}
			}

		case "chat_pingo":
			if featureType == "conversation" {
				newValue++
				if newValue >= m.TargetValue {
					completed = true
				}
			}
		}

		if newValue != p.CurrentValue || completed != p.IsCompleted {
			_ = s.missionRepo.UpdateMissionProgress(ctx, userID, m.ID, newValue, completed)

			// kalau baru selesai, kasih XP bonus
			if completed && !p.IsCompleted {
				_ = s.userRepo.AddXP(ctx, userID, m.XPReward)
			}
		}
	}

	return nil
}

func missionDescription(missionType string, targetValue int) string {
	switch missionType {
	case "collect_xp":
		return "Kumpulkan 30 XP hari ini"
	case "complete_exercises":
		return "Selesaikan 2 latihan hari ini"
	case "chat_pingo":
		return "Lakukan percakapan dengan Pingo"
	default:
		return ""
	}
}