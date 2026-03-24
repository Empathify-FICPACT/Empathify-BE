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
	ErrEmotionSessionNotFound  = errors.New("emotion session not found")
	ErrEmotionSessionNotOwned  = errors.New("emotion session does not belong to user")
	ErrEmotionSessionCompleted = errors.New("emotion session already completed")
	ErrEmotionQuestionNotFound = errors.New("emotion question not found")
	ErrEmotionAlreadyAnswered  = errors.New("all questions already answered")
	ErrInvalidAnswer           = errors.New("invalid answer option")
)

const emotionXP = 20

type EmotionService interface {
	StartSession(ctx context.Context, userID string, req *request.StartEmotionRequest) (*response.EmotionSessionResponse, error)
	SubmitAnswer(ctx context.Context, userID, sessionID string, req *request.SubmitEmotionAnswerRequest) (*response.EmotionAnswerResponse, error)
	CompleteSession(ctx context.Context, userID, sessionID string) (*response.CompleteEmotionResponse, error)
}

type emotionService struct {
	emotionRepo repository.EmotionRepository
	userRepo    repository.UserRepository
	missionSvc  MissionService
}

func NewEmotionService(
	emotionRepo repository.EmotionRepository,
	userRepo repository.UserRepository,
	missionSvc  MissionService,
) EmotionService {
	return &emotionService{
		emotionRepo: emotionRepo,
		userRepo:    userRepo,
		missionSvc: missionSvc,
	}
}

func (s *emotionService) StartSession(ctx context.Context, userID string, req *request.StartEmotionRequest) (*response.EmotionSessionResponse, error) {
	total := req.TotalQuestions
	if total < 5 || total > 8 {
		total = 5
	}

	questions, err := s.emotionRepo.GetRandomQuestions(ctx, total)
	if err != nil {
		return nil, err
	}
	if len(questions) < total {
		total = len(questions)
	}

	session := &domain.EmotionSession{
		ID:             uuid.NewString(),
		UserID:         userID,
		TotalQuestions: total,
		CorrectCount:   0,
		XPEarned:       0,
		Status:         "in_progress",
		StartedAt:      time.Now(),
	}

	if err := s.emotionRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	var questionResponses []response.EmotionQuestionResponse
	for i, q := range questions {
		questionResponses = append(questionResponses, response.EmotionQuestionResponse{
			ID:        q.ID,
			Situation: q.Situation,
			OptionA:   q.OptionA,
			OptionB:   q.OptionB,
			OptionC:   q.OptionC,
			OptionD:   q.OptionD,
			Order:     i + 1,
		})
	}

	return &response.EmotionSessionResponse{
		ID:             session.ID,
		Status:         session.Status,
		TotalQuestions: session.TotalQuestions,
		Questions:      questionResponses,
		StartedAt:      session.StartedAt,
	}, nil
}

func (s *emotionService) SubmitAnswer(ctx context.Context, userID, sessionID string, req *request.SubmitEmotionAnswerRequest) (*response.EmotionAnswerResponse, error) {
	session, err := s.emotionRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrEmotionSessionNotFound
	}
	if session.UserID != userID {
		return nil, ErrEmotionSessionNotOwned
	}
	if session.Status == "completed" {
		return nil, ErrEmotionSessionCompleted
	}

	// validasi jawaban
	validAnswers := map[string]bool{"a": true, "b": true, "c": true, "d": true}
	if !validAnswers[req.ChosenAnswer] {
		return nil, ErrInvalidAnswer
	}

	// cek sudah berapa jawaban yang masuk
	answers, err := s.emotionRepo.GetAnswersBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if len(answers) >= session.TotalQuestions {
		return nil, ErrEmotionAlreadyAnswered
	}

	// ambil soal
	question, err := s.emotionRepo.GetQuestionByID(ctx, req.QuestionID)
	if err != nil {
		return nil, err
	}
	if question == nil {
		return nil, ErrEmotionQuestionNotFound
	}

	isCorrect := req.ChosenAnswer == question.CorrectAnswer
	questionOrder := len(answers) + 1

	answer := &domain.EmotionAnswer{
		ID:            uuid.NewString(),
		SessionID:     sessionID,
		QuestionID:    req.QuestionID,
		ChosenAnswer:  req.ChosenAnswer,
		IsCorrect:     isCorrect,
		QuestionOrder: questionOrder,
		CreatedAt:     time.Now(),
	}

	if err := s.emotionRepo.CreateAnswer(ctx, answer); err != nil {
		return nil, err
	}

	return &response.EmotionAnswerResponse{
		QuestionID:    req.QuestionID,
		ChosenAnswer:  req.ChosenAnswer,
		CorrectAnswer: question.CorrectAnswer,
		IsCorrect:     isCorrect,
		QuestionOrder: questionOrder,
	}, nil
}

func (s *emotionService) CompleteSession(ctx context.Context, userID, sessionID string) (*response.CompleteEmotionResponse, error) {
	session, err := s.emotionRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrEmotionSessionNotFound
	}
	if session.UserID != userID {
		return nil, ErrEmotionSessionNotOwned
	}
	if session.Status == "completed" {
		return nil, ErrEmotionSessionCompleted
	}

	answers, err := s.emotionRepo.GetAnswersBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	correctCount := 0
	for _, a := range answers {
		if a.IsCorrect {
			correctCount++
		}
	}

	if err := s.emotionRepo.CompleteSession(ctx, sessionID, correctCount, emotionXP); err != nil {
		return nil, err
	}

	if err := s.userRepo.AddXP(ctx, userID, emotionXP); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	_ = s.missionSvc.UpdateProgressAfterSession(ctx, userID, "emotion", emotionXP)

	return &response.CompleteEmotionResponse{
		SessionID:    sessionID,
		CorrectCount: correctCount,
		Total:        session.TotalQuestions,
		XPEarned:     emotionXP,
		TotalXP:      user.TotalXP,
	}, nil
}