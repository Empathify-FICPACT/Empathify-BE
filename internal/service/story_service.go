package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/repository"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/request"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/response"
)

var (
	ErrStorySessionNotFound  = errors.New("story session not found")
	ErrStorySessionNotOwned  = errors.New("story session does not belong to user")
	ErrStorySessionCompleted = errors.New("story session already completed")
	ErrStoryScenarioNotFound = errors.New("story scenario not found")
	ErrStoryAlreadyAnswered  = errors.New("all scenarios already answered")
	ErrInvalidStoryAnswer    = errors.New("invalid answer option")
)

const storyXP = 25

type StoryService interface {
	StartSession(ctx context.Context, userID string, req *request.StartStoryRequest) (*response.StorySessionResponse, error)
	SubmitAnswer(ctx context.Context, userID, sessionID string, req *request.SubmitStoryAnswerRequest) (*response.StoryAnswerResponse, error)
	CompleteSession(ctx context.Context, userID, sessionID string) (*response.CompleteStoryResponse, error)
}

type storyService struct {
	storyRepo repository.StoryRepository
	userRepo  repository.UserRepository
	missionSvc  MissionService
}

func NewStoryService(
	storyRepo repository.StoryRepository,
	userRepo repository.UserRepository,
	missionSvc  MissionService,
) StoryService {
	return &storyService{
		storyRepo: storyRepo,
		userRepo:  userRepo,
		missionSvc: missionSvc,
	}
}

func (s *storyService) StartSession(ctx context.Context, userID string, req *request.StartStoryRequest) (*response.StorySessionResponse, error) {
	total := req.TotalScenarios
	if total < 5 || total > 8 {
		total = 5
	}

	scenarios, err := s.storyRepo.GetRandomScenarios(ctx, total)
	if err != nil {
		return nil, err
	}
	if len(scenarios) < total {
		total = len(scenarios)
	}

	session := &domain.StorySession{
		ID:             uuid.NewString(),
		UserID:         userID,
		TotalScenarios: total,
		CorrectCount:   0,
		XPEarned:       0,
		Status:         "in_progress",
		StartedAt:      time.Now(),
	}

	if err := s.storyRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	var scenarioResponses []response.StoryScenarioResponse
	for i, sc := range scenarios {
		scenarioResponses = append(scenarioResponses, response.StoryScenarioResponse{
			ID:        sc.ID,
			Situation: sc.Situation,
			OptionA:   sc.OptionA,
			OptionB:   sc.OptionB,
			OptionC:   sc.OptionC,
			OptionD:   sc.OptionD,
			Order:     i + 1,
		})
	}

	return &response.StorySessionResponse{
		ID:             session.ID,
		Status:         session.Status,
		TotalScenarios: session.TotalScenarios,
		Scenarios:      scenarioResponses,
		StartedAt:      session.StartedAt,
	}, nil
}

func (s *storyService) SubmitAnswer(ctx context.Context, userID, sessionID string, req *request.SubmitStoryAnswerRequest) (*response.StoryAnswerResponse, error) {
	session, err := s.storyRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrStorySessionNotFound
	}
	if session.UserID != userID {
		return nil, ErrStorySessionNotOwned
	}
	if session.Status == "completed" {
		return nil, ErrStorySessionCompleted
	}

	validAnswers := map[string]bool{"a": true, "b": true, "c": true, "d": true}
	if !validAnswers[req.ChosenAnswer] {
		return nil, ErrInvalidStoryAnswer
	}

	answers, err := s.storyRepo.GetAnswersBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if len(answers) >= session.TotalScenarios {
		return nil, ErrStoryAlreadyAnswered
	}

	scenario, err := s.storyRepo.GetScenarioByID(ctx, req.ScenarioID)
	if err != nil {
		return nil, err
	}
	if scenario == nil {
		return nil, ErrStoryScenarioNotFound
	}

	isCorrect := req.ChosenAnswer == scenario.CorrectAnswer
	scenarioOrder := len(answers) + 1

	answer := &domain.StoryAnswer{
		ID:            uuid.NewString(),
		SessionID:     sessionID,
		ScenarioID:    req.ScenarioID,
		ChosenAnswer:  req.ChosenAnswer,
		IsCorrect:     isCorrect,
		ScenarioOrder: scenarioOrder,
		CreatedAt:     time.Now(),
	}

	if err := s.storyRepo.CreateAnswer(ctx, answer); err != nil {
		return nil, err
	}

	return &response.StoryAnswerResponse{
		ScenarioID:    req.ScenarioID,
		ChosenAnswer:  req.ChosenAnswer,
		CorrectAnswer: scenario.CorrectAnswer,
		IsCorrect:     isCorrect,
		Explanation:   scenario.Explanation,
		ScenarioOrder: scenarioOrder,
	}, nil
}

func (s *storyService) CompleteSession(ctx context.Context, userID, sessionID string) (*response.CompleteStoryResponse, error) {
	session, err := s.storyRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrStorySessionNotFound
	}
	if session.UserID != userID {
		return nil, ErrStorySessionNotOwned
	}
	if session.Status == "completed" {
		return nil, ErrStorySessionCompleted
	}

	answers, err := s.storyRepo.GetAnswersBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	correctCount := 0
	for _, a := range answers {
		if a.IsCorrect {
			correctCount++
		}
	}

	if err := s.storyRepo.CompleteSession(ctx, sessionID, correctCount, storyXP); err != nil {
		return nil, err
	}

	if err := s.userRepo.AddXP(ctx, userID, storyXP); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	_ = s.missionSvc.UpdateProgressAfterSession(ctx, userID, "story", storyXP)

	return &response.CompleteStoryResponse{
		SessionID:    sessionID,
		CorrectCount: correctCount,
		Total:        session.TotalScenarios,
		XPEarned:     storyXP,
		TotalXP:      user.TotalXP,
	}, nil
}