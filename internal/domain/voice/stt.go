package voice

import "context"

type TranscriptionRequest struct {
	Audio []byte
	Format string
	Language string
	Provider string
}

type TranscriptionResult struct {
	Text string
	Confidence float64
}

type STTProvider interface {
	Transcribe(ctx context.Context, req TranscriptionRequest) (TranscriptionResult, error)
}