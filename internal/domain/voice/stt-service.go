package voice

import (
	"context"
	"fmt"
)

type STTService struct {
	providers map[string]STTProvider
	defaultProvider string
}

func NewSTTService(providers map[string]STTProvider, defaultProvider string) *STTService {
	return &STTService{
		providers: providers,
		defaultProvider: defaultProvider,
	}
}

func (s *STTService) Transcribe(ctx context.Context, req TranscriptionRequest) (TranscriptionResult, error) {
	providerName := req.Provider
	if providerName == "" {
		providerName = s.defaultProvider
	}
	provider, ok := s.providers[providerName]
	if !ok {
		return TranscriptionResult{}, fmt.Errorf("provider %s not found", providerName)
	}
	return provider.Transcribe(ctx, req)
}