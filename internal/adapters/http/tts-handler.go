package httpadapter

import (
	"encoding/json"
	"net/http"
	"polyglot/internal/domain/voice"
)

type TTSHandler struct {
	ttsService *voice.TTSService
}

func NewTTSHandler(ttsService *voice.TTSService) *TTSHandler {
	return &TTSHandler{ttsService: ttsService}
}

func (h *TTSHandler) GenerateSpeech(w http.ResponseWriter, r *http.Request) {
	var req voice.SpeechRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.ttsService.GenerateSpeech(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(result.Audio)
}