package voice

import (
	"context"
	"fmt"
)

type TTSService struct {
	providers map[string]TTSProvider
	defaultProvider string
}

func NewTTSService(providers map[string]TTSProvider, defaultProvider string) *TTSService {
	return &TTSService{
		providers: providers,
		defaultProvider: defaultProvider,
	}
}

func (s *TTSService) GenerateSpeech(ctx context.Context, req SpeechRequest) (SpeechResult, error) {
	providerName := req.Provider
	if providerName == "" {
		providerName = s.defaultProvider
	}
	provider, ok := s.providers[providerName]
	if !ok {
		return SpeechResult{}, fmt.Errorf("provider %s not found", providerName)
	}
	return provider.Synthesize(ctx, req)
}