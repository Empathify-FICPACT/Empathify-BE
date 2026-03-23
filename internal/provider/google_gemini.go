package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

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

	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", contents, &genai.GenerateContentConfig{
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

func (g *GeminiProvider) AnalyzeExpression(ctx context.Context, referenceImageURL string, userPhotoBytes []byte, emotion string) (float64, string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.App.GeminiAPIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return 0, "", ErrGeminiFailed
	}

	// fetch gambar referensi dari URL
	refImageBytes, err := fetchImageFromURL(referenceImageURL)
	if err != nil {
		return 0, "", ErrGeminiFailed
	}

	prompt := fmt.Sprintf(`Bandingkan dua gambar ekspresi wajah berikut.
Gambar pertama adalah referensi ekspresi "%s".
Gambar kedua adalah foto ekspresi dari pengguna.

Berikan:
1. Skor kemiripan ekspresi antara 0.0 hingga 1.0 (1.0 = sangat mirip)
2. Feedback singkat dalam bahasa Indonesia yang ramah untuk anak-anak (max 2 kalimat)

Jawab HANYA dalam format JSON berikut tanpa penjelasan tambahan:
{"score": 0.0, "feedback": "..."}`, emotion)

	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash",
		[]*genai.Content{
			{
				Role: "user",
				Parts: []*genai.Part{
					{
						InlineData: &genai.Blob{
							MIMEType: "image/jpeg",
							Data:     refImageBytes,
						},
					},
					{
						InlineData: &genai.Blob{
							MIMEType: "image/jpeg",
							Data:     userPhotoBytes,
						},
					},
					{Text: prompt},
				},
			},
		},
		nil,
	)
	if err != nil {
		return 0, "", ErrGeminiFailed
	}

	if resp == nil || len(resp.Candidates) == 0 {
		return 0, "", ErrGeminiFailed
	}

	// parse JSON response
	text := resp.Text()
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var result struct {
		Score    float64 `json:"score"`
		Feedback string  `json:"feedback"`
	}
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return 0, "", ErrGeminiFailed
	}

	return result.Score, result.Feedback, nil
}

func fetchImageFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}