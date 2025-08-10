package main

import (
	"fmt"
	"log"
	"net/http"
	"polyglot/internal/adapters/elevenlabs"
	httpadapter "polyglot/internal/adapters/http"
	"polyglot/internal/config"
	"polyglot/internal/domain/voice"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadEnv()
	
	elevenlabsTTS := elevenlabs.NewTTSAdapter(cfg.Providers.ElevenLabs.APIKey)
	ttsProviders := map[string]voice.TTSProvider{
		"elevenlabs": elevenlabsTTS,
	}
	ttsService := voice.NewTTSService(ttsProviders, cfg.Defaults.TTSProvider)
	ttsHandler := httpadapter.NewTTSHandler(ttsService)

	elevenlabsSTT := elevenlabs.NewSTTAdapter(cfg.Providers.ElevenLabs.APIKey)
	sttProviders := map[string]voice.STTProvider{
		"elevenlabs": elevenlabsSTT,
	}
	sttService := voice.NewSTTService(sttProviders, cfg.Defaults.STTProvider)
	sttHandler := httpadapter.NewSTTHandler(sttService)

	router := chi.NewRouter()
	router.Post("/tts", ttsHandler.GenerateSpeech)
	router.Post("/stt", sttHandler.Transcribe)
	log.Printf("Server is running on port %d", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), router))
}