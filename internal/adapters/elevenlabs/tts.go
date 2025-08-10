package elevenlabs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"polyglot/internal/domain/voice"
)

type ElevenLabsTTS struct {
	apiKey string
}

func NewTTSAdapter(apiKey string) *ElevenLabsTTS {
	return &ElevenLabsTTS{
		apiKey: apiKey,
	}
}

func (e *ElevenLabsTTS) Synthesize(ctx context.Context, req voice.SpeechRequest) (voice.SpeechResult, error) {
    url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", req.Voice)

    body := map[string]interface{}{
        "text":     req.Text,
        "model_id": "eleven_turbo_v2",
        "voice_settings": map[string]float64{
            "stability":        0.75,
            "similarity_boost": 0.75,
        },
    }

    payload, err := json.Marshal(body)
    if err != nil {
        return voice.SpeechResult{}, err
    }

	fmt.Printf("DEBUG URL: %q\n", url)
	fmt.Printf("DEBUG Payload: %s\n", payload)


    httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
    if err != nil {
        return voice.SpeechResult{}, err
    }

    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Accept", "audio/mpeg")
    httpReq.Header.Set("xi-api-key", e.apiKey)

    resp, err := http.DefaultClient.Do(httpReq)
    if err != nil {
        return voice.SpeechResult{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return voice.SpeechResult{}, fmt.Errorf("elevenlabs error: %s", string(body))
    }

    audio, err := io.ReadAll(resp.Body)
    if err != nil {
        return voice.SpeechResult{}, err
    }

    return voice.SpeechResult{
        Audio:  audio,
        Format: "mp3",
    }, nil
}
