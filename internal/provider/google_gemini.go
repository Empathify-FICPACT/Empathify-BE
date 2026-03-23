package provider

import (
	"context"
	"errors"

	"google.golang.org/genai"

	"github.com/Empathify-FICPACT/Empathify-BE/config"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
)

var ErrGeminiFailed = errors.New("failed to get response from Gemini")

type GeminiProvider struct{}

func NewGeminiProvider() *GeminiProvider {
	return &GeminiProvider{}
}

func (g *GeminiProvider) Chat(ctx context.Context, systemPrompt string, history []domain.ConversationMessage, userMessage string) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.App.GeminiAPIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", ErrGeminiFailed
	}

	// convert history ke format genai
	var contents []*genai.Content
	for _, msg := range history {
		role := "user"
		if msg.Role == "assistant" {
			role = "model"
		}
		contents = append(contents, &genai.Content{
			Role:  role,
			Parts: []*genai.Part{{Text: msg.Content}},
		})
	}

	// tambah pesan user sekarang
	contents = append(contents, &genai.Content{
		Role:  "user",
		Parts: []*genai.Part{{Text: userMessage}},
	})

	resp, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash", contents, &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: systemPrompt}},
		},
	})
	if err != nil {
		return "", ErrGeminiFailed
	}

	if resp == nil || len(resp.Candidates) == 0 {
		return "", ErrGeminiFailed
	}

	return resp.Text(), nil
}