package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/provider"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/repository"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/request"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/response"
)

var (
	ErrExpressionSessionNotFound  = errors.New("expression session not found")
	ErrExpressionSessionNotOwned  = errors.New("expression session does not belong to user")
	ErrExpressionSessionCompleted = errors.New("expression session already completed")
	ErrReferenceNotFound          = errors.New("expression reference not found")
	ErrAllAttemptsSubmitted       = errors.New("all expressions already submitted")
)

const (
	expressionXP        = 15
	similarityThreshold = 0.6 // score >= 0.6 dianggap benar
)

type ExpressionService interface {
	StartSession(ctx context.Context, userID string, req *request.StartExpressionRequest) (*response.ExpressionSessionResponse, error)
	SubmitAttempt(ctx context.Context, userID, sessionID, referenceID string, photoBytes []byte) (*response.ExpressionAttemptResponse, error)
	CompleteSession(ctx context.Context, userID, sessionID string) (*response.CompleteExpressionResponse, error)
}

type expressionService struct {
	exprRepo repository.ExpressionRepository
	userRepo repository.UserRepository
	gemini   *provider.GeminiProvider
}

func NewExpressionService(
	exprRepo repository.ExpressionRepository,
	userRepo repository.UserRepository,
	gemini *provider.GeminiProvider,
) ExpressionService {
	return &expressionService{
		exprRepo: exprRepo,
		userRepo: userRepo,
		gemini:   gemini,
	}
}

func (s *expressionService) StartSession(ctx context.Context, userID string, req *request.StartExpressionRequest) (*response.ExpressionSessionResponse, error) {
	total := req.TotalExpressions
	if total < 3 || total > 5 {
		total = 3 // default
	}

	// random ekspresi dari DB
	refs, err := s.exprRepo.GetRandomReferences(ctx, total)
	if err != nil {
		return nil, err
	}
	if len(refs) < total {
		total = len(refs) // kalau data kurang, sesuaikan
	}

	session := &domain.ExpressionSession{
		ID:               uuid.NewString(),
		UserID:           userID,
		TotalExpressions: total,
		CorrectCount:     0,
		XPEarned:         0,
		Status:           "in_progress",
		StartedAt:        time.Now(),
	}

	if err := s.exprRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	// build response expressions
	var expressions []response.ExpressionReferenceResponse
	for i, ref := range refs {
		expressions = append(expressions, response.ExpressionReferenceResponse{
			ID:       ref.ID,
			Emotion:  ref.Emotion,
			ImageURL: ref.ImageURL,
			Order:    i + 1,
		})
	}

	return &response.ExpressionSessionResponse{
		ID:               session.ID,
		Status:           session.Status,
		TotalExpressions: session.TotalExpressions,
		Expressions:      expressions,
		StartedAt:        session.StartedAt,
	}, nil
}

func (s *expressionService) SubmitAttempt(ctx context.Context, userID, sessionID, referenceID string, photoBytes []byte) (*response.ExpressionAttemptResponse, error) {
	session, err := s.exprRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrExpressionSessionNotFound
	}
	if session.UserID != userID {
		return nil, ErrExpressionSessionNotOwned
	}
	if session.Status == "completed" {
		return nil, ErrExpressionSessionCompleted
	}

	// cek sudah berapa attempt yang masuk
	attempts, err := s.exprRepo.GetAttemptsBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if len(attempts) >= session.TotalExpressions {
		return nil, ErrAllAttemptsSubmitted
	}

	// ambil referensi
	ref, err := s.exprRepo.GetReferenceByID(ctx, referenceID)
	if err != nil {
		return nil, err
	}
	if ref == nil {
		return nil, ErrReferenceNotFound
	}

	// analisis dengan Gemini
	score, feedback, err := s.gemini.AnalyzeExpression(ctx, ref.ImageURL, photoBytes, ref.Emotion)
	if err != nil {
		return nil, err
	}

	isCorrect := score >= similarityThreshold
	attemptOrder := len(attempts) + 1

	attempt := &domain.ExpressionAttempt{
		ID:              uuid.NewString(),
		SessionID:       sessionID,
		ReferenceID:     referenceID,
		UserPhotoURL:    "", // opsional: simpan ke storage
		SimilarityScore: score,
		IsCorrect:       isCorrect,
		AIFeedback:      &feedback,
		AttemptOrder:    attemptOrder,
		CreatedAt:       time.Now(),
	}

	if err := s.exprRepo.CreateAttempt(ctx, attempt); err != nil {
		return nil, err
	}

	return &response.ExpressionAttemptResponse{
		AttemptID:       attempt.ID,
		ReferenceID:     referenceID,
		Emotion:         ref.Emotion,
		SimilarityScore: score,
		IsCorrect:       isCorrect,
		AIFeedback:      feedback,
		AttemptOrder:    attemptOrder,
	}, nil
}

func (s *expressionService) CompleteSession(ctx context.Context, userID, sessionID string) (*response.CompleteExpressionResponse, error) {
	session, err := s.exprRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrExpressionSessionNotFound
	}
	if session.UserID != userID {
		return nil, ErrExpressionSessionNotOwned
	}
	if session.Status == "completed" {
		return nil, ErrExpressionSessionCompleted
	}

	// hitung correct count dari attempts
	attempts, err := s.exprRepo.GetAttemptsBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	correctCount := 0
	for _, a := range attempts {
		if a.IsCorrect {
			correctCount++
		}
	}

	// complete session
	if err := s.exprRepo.CompleteSession(ctx, sessionID, correctCount, expressionXP); err != nil {
		return nil, err
	}

	// tambah XP user
	if err := s.userRepo.AddXP(ctx, userID, expressionXP); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &response.CompleteExpressionResponse{
		SessionID:    sessionID,
		CorrectCount: correctCount,
		Total:        session.TotalExpressions,
		XPEarned:     expressionXP,
		TotalXP:      user.TotalXP,
	}, nil
}