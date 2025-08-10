package elevenlabs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"polyglot/internal/domain/voice"
)

type ElevenLabsSTT struct {
	apiKey string
}

func NewSTTAdapter(apiKey string) *ElevenLabsSTT {
	return &ElevenLabsSTT{
		apiKey: apiKey,
	}
}

func (e *ElevenLabsSTT) Transcribe(ctx context.Context, req voice.TranscriptionRequest) (voice.TranscriptionResult, error) {
	url := "https://api.elevenlabs.io/v1/speech-to-text"

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	fw, err := writer.CreateFormFile("file", fmt.Sprintf("audio.%s", req.Format))
	if err != nil {
		return voice.TranscriptionResult{}, err
	}

	_, err = fw.Write(req.Audio)
	if err != nil {
		return voice.TranscriptionResult{}, err
	}

	if err := writer.WriteField("model_id", "scribe_v1"); err != nil {
		return voice.TranscriptionResult{}, err
	}

	writer.Close()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, &buf)
	if err != nil {
		return voice.TranscriptionResult{}, err
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("xi-api-key", e.apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return voice.TranscriptionResult{}, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return voice.TranscriptionResult{}, fmt.Errorf("elevenlabs error: %s", string(body))
	}

	var result struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return voice.TranscriptionResult{}, err
	}

	return voice.TranscriptionResult{
		Text: result.Text,
	}, nil

}