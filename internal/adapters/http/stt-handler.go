package httpadapter

import (
	"io"
	"net/http"
	"polyglot/internal/domain/voice"
)

type STTHandler struct {
	sttService *voice.STTService
}

func NewSTTHandler(sttService *voice.STTService) *STTHandler {
	return &STTHandler{sttService: sttService}
}

func (h *STTHandler) Transcribe(w http.ResponseWriter, r *http.Request) {
	audioData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "mp3"
	}

	result, err := h.sttService.Transcribe(r.Context(), voice.TranscriptionRequest{
		Audio: audioData,
		Format: format,
		Language: "en",
		Provider: "elevenlabs",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result.Text))
}