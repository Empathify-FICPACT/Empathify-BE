package provider

import (
	"context"
	"errors"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"google.golang.org/api/option"

	"github.com/Empathify-FICPACT/Empathify-BE/config"
)

var ErrSTTFailed = errors.New("failed to transcribe audio")

type STTProvider struct {
	client *speech.Client
}

func NewSTTProvider() *STTProvider {
	ctx := context.Background()

	client, err := speech.NewClient(ctx, option.WithCredentialsFile(config.App.GoogleSTTCredentials))
	if err != nil {
		panic(err)
	}

	return &STTProvider{client: client}
}

func (s *STTProvider) Transcribe(ctx context.Context, audioBytes []byte) (string, error) {
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_MP3,
			SampleRateHertz: 44100,
			LanguageCode:    "id-ID",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{
				Content: audioBytes,
			},
		},
	}

	resp, err := s.client.Recognize(ctx, req)
	if err != nil {
		log.Println("GOOGLE STT ERROR:", err)
		return "", err
	}

	if len(resp.Results) == 0 {
		log.Println("STT EMPTY RESPONSE")
		return "", errors.New("empty transcription")
	}

	return resp.Results[0].Alternatives[0].Transcript, nil
}
