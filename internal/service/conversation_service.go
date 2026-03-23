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
	ErrTopicNotFound      = errors.New("topic not found")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionNotOwned    = errors.New("session does not belong to user")
	ErrSessionCompleted   = errors.New("session already completed")
)

const conversationXP = 15

type ConversationService interface {
	GetTopics(ctx context.Context) ([]response.TopicResponse, error)
	StartSession(ctx context.Context, userID string, req *request.StartConversationRequest) (*response.SessionResponse, error)
	SendMessage(ctx context.Context, userID, sessionID string, req *request.SendMessageRequest) (*response.SendMessageResponse, error)
	CompleteSession(ctx context.Context, userID, sessionID string) (*response.CompleteSessionResponse, error)
}

type conversationService struct {
	convRepo repository.ConversationRepository
	userRepo repository.UserRepository
	gemini   *provider.GeminiProvider
}

func NewConversationService(
	convRepo repository.ConversationRepository,
	userRepo repository.UserRepository,
	gemini *provider.GeminiProvider,
) ConversationService {
	return &conversationService{
		convRepo: convRepo,
		userRepo: userRepo,
		gemini:   gemini,
	}
}

func (s *conversationService) GetTopics(ctx context.Context) ([]response.TopicResponse, error) {
	topics, err := s.convRepo.GetActiveTopics(ctx)
	if err != nil {
		return nil, err
	}

	var result []response.TopicResponse
	for _, t := range topics {
		desc := ""
		if t.Description != nil {
			desc = *t.Description
		}
		result = append(result, response.TopicResponse{
			ID:          t.ID,
			Title:       t.Title,
			Description: desc,
			Difficulty:  t.Difficulty,
		})
	}
	return result, nil
}

func (s *conversationService) StartSession(ctx context.Context, userID string, req *request.StartConversationRequest) (*response.SessionResponse, error) {
	topic, err := s.convRepo.GetTopicByID(ctx, req.TopicID)
	if err != nil {
		return nil, err
	}
	if topic == nil {
		return nil, ErrTopicNotFound
	}

	session := &domain.ConversationSession{
		ID:        uuid.NewString(),
		UserID:    userID,
		TopicID:   topic.ID,
		Status:    "in_progress",
		XPEarned:  0,
		StartedAt: time.Now(),
	}

	if err := s.convRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	// kirim pesan pembuka dari AI
	openingMsg, err := s.gemini.Chat(ctx, topic.SystemPrompt, nil, "Mulai percakapan")
	if err != nil {
		return nil, err
	}

	aiMsg := &domain.ConversationMessage{
		ID:        uuid.NewString(),
		SessionID: session.ID,
		Role:      "assistant",
		Content:   openingMsg,
		CreatedAt: time.Now(),
	}
	if err := s.convRepo.CreateMessage(ctx, aiMsg); err != nil {
		return nil, err
	}

	return &response.SessionResponse{
		ID:         session.ID,
		TopicID:    session.TopicID,
		TopicTitle: topic.Title,
		Status:     session.Status,
		XPEarned:   session.XPEarned,
		StartedAt:  session.StartedAt,
	}, nil
}

func (s *conversationService) SendMessage(ctx context.Context, userID, sessionID string, req *request.SendMessageRequest) (*response.SendMessageResponse, error) {
	session, err := s.convRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}
	if session.UserID != userID {
		return nil, ErrSessionNotOwned
	}
	if session.Status == "completed" {
		return nil, ErrSessionCompleted
	}

	// ambil topic untuk system prompt
	topic, err := s.convRepo.GetTopicByID(ctx, session.TopicID)
	if err != nil {
		return nil, err
	}

	// ambil history pesan
	history, err := s.convRepo.GetSessionMessages(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	// simpan pesan user
	userMsg := &domain.ConversationMessage{
		ID:        uuid.NewString(),
		SessionID: sessionID,
		Role:      "user",
		Content:   req.Content,
		CreatedAt: now,
	}
	if err := s.convRepo.CreateMessage(ctx, userMsg); err != nil {
		return nil, err
	}

	// kirim ke Gemini dengan history
	aiResponse, err := s.gemini.Chat(ctx, topic.SystemPrompt, history, req.Content)
	if err != nil {
		return nil, err
	}

	// simpan response AI
	aiMsg := &domain.ConversationMessage{
		ID:        uuid.NewString(),
		SessionID: sessionID,
		Role:      "assistant",
		Content:   aiResponse,
		CreatedAt: time.Now(),
	}
	if err := s.convRepo.CreateMessage(ctx, aiMsg); err != nil {
		return nil, err
	}

	return &response.SendMessageResponse{
		UserMessage: response.MessageResponse{
			ID:        userMsg.ID,
			Role:      userMsg.Role,
			Content:   userMsg.Content,
			CreatedAt: userMsg.CreatedAt,
		},
		AIMessage: response.MessageResponse{
			ID:        aiMsg.ID,
			Role:      aiMsg.Role,
			Content:   aiMsg.Content,
			CreatedAt: aiMsg.CreatedAt,
		},
	}, nil
}

func (s *conversationService) CompleteSession(ctx context.Context, userID, sessionID string) (*response.CompleteSessionResponse, error) {
	session, err := s.convRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}
	if session.UserID != userID {
		return nil, ErrSessionNotOwned
	}
	if session.Status == "completed" {
		return nil, ErrSessionCompleted
	}

	// complete session + kasih XP
	if err := s.convRepo.CompleteSession(ctx, sessionID, conversationXP); err != nil {
		return nil, err
	}

	// tambah total_xp user
	if err := s.userRepo.AddXP(ctx, userID, conversationXP); err != nil {
		return nil, err
	}

	// ambil total xp terbaru
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &response.CompleteSessionResponse{
		SessionID: sessionID,
		XPEarned:  conversationXP,
		TotalXP:   user.TotalXP,
	}, nil
}