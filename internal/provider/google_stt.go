package provider

import (
	"context"
	"errors"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"google.golang.org/api/option"

	"github.com/Empathify-FICPACT/Empathify-BE/config"
)

var ErrSTTFailed = errors.New("failed to transcribe audio")

type STTProvider struct{}

func NewSTTProvider() *STTProvider {
	return &STTProvider{}
}

func (s *STTProvider) Transcribe(ctx context.Context, audioBytes []byte) (string, error) {
	client, err := speech.NewClient(ctx, option.WithCredentialsFile(config.App.GoogleSTTCredentials))
	if err != nil {
		return "", ErrSTTFailed
	}
	defer client.Close()

	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_MP3,
			SampleRateHertz: 44100,
			LanguageCode:    "id-ID",
			Model:           "latest_long",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{
				Content: audioBytes,
			},
		},
	}

	resp, err := client.Recognize(ctx, req)
	if err != nil {
		return "", ErrSTTFailed
	}

	if len(resp.Results) == 0 || len(resp.Results[0].Alternatives) == 0 {
		return "", ErrSTTFailed
	}

	return resp.Results[0].Alternatives[0].Transcript, nil
}
