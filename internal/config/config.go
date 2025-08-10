package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Server struct {
	Port         int
}

type Defaults struct {
	TTSProvider string
	STTProvider string
}

type Providers struct {
	ElevenLabs struct {
		APIKey string
	}
}

type Config struct {
	Server Server
	Defaults Defaults
	Providers Providers
}

func mustString(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("%s is required", key))
	}
	return v
}

func mustInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Sprintf("%s must be an integer", key))
	}
	return i
}

func (c Config) Validate() error {
	if c.Defaults.TTSProvider == "" {
		return errors.New("DEFAULT_TTS_PROVIDER is required")
	}
	if c.Defaults.STTProvider == "" {
		return errors.New("DEFAULT_STT_PROVIDER is required")
	}
	return nil
}

func LoadEnv() Config {
	godotenv.Load()
	cfg := Config{
		Server: Server{
			Port: mustInt("PORT", 8080),
		},
		Defaults: Defaults{
			TTSProvider: mustString("DEFAULT_TTS_PROVIDER"),
			STTProvider: mustString("DEFAULT_STT_PROVIDER"),
		},
		Providers: Providers{
			ElevenLabs: struct {
				APIKey string
			}{
				APIKey: mustString("ELEVEN_LABS_API_KEY"),
			},

		},
	}
	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	return cfg
}