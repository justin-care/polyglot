package voice

import "context"

type SpeechRequest struct {
	Text string
	Voice string
	Format string
	Provider string
}

type SpeechResult struct {
	Audio []byte
	Format string
}

type TTSProvider interface {
	Synthesize(ctx context.Context, req SpeechRequest) (SpeechResult, error)
}
