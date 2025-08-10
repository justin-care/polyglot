## polyglot

A small Go service to experiment with hexagonal architecture (ports and adapters) for Text-to-Speech (TTS) and Speech-to-Text (STT), currently using ElevenLabs as the external provider.

### Why

This repo was created to explore a clear separation of concerns via hexagonal architecture:

- **Domain (core)**: business logic and ports (interfaces)
- **Adapters**: HTTP inbounds and ElevenLabs outbounds
- **Composition**: wiring in `main` without leaking implementation details into the core

## Architecture

- **Domain layer** (`internal/domain/voice`)
  - Ports: `TTSProvider`, `STTProvider`
  - Use cases/services: `TTSService`, `STTService` route requests to a chosen provider
  - Data types: `SpeechRequest`, `SpeechResult`, `TranscriptionRequest`, `TranscriptionResult`

- **Inbound adapter** (`internal/adapters/http`)
  - `TTSHandler` exposes `POST /tts` (JSON in, audio out)
  - `STTHandler` exposes `POST /stt` (audio in, text out)

- **Outbound adapter** (`internal/adapters/elevenlabs`)
  - Implements domain ports to call ElevenLabs REST APIs for TTS and STT

- **Configuration** (`internal/config`)
  - Loads env vars, validates required values, provides defaults

- **Composition root** (`cmd/server/main.go`)
  - Builds providers map, initializes services, registers HTTP routes

### Request flow (example: TTS)

HTTP `POST /tts` → `TTSHandler` → `TTSService` → `TTSProvider` port → `ElevenLabsTTS` adapter → ElevenLabs API → audio response.

## Project layout

```
cmd/
  server/
    main.go
internal/
  adapters/
    http/           # HTTP handlers (inbound)
      tts-handler.go
      stt-handler.go
    elevenlabs/     # External provider (outbound)
      tts.go
      stt.go
  config/           # Env loading and validation
    config.go
  domain/
    voice/          # Core domain: ports and services
      tts.go
      tts-service.go
      stt.go
      stt-service.go
```

## Requirements

- Go 1.23+
- ElevenLabs API key

## Configuration

Set environment variables (a `.env` file is supported via `github.com/joho/godotenv`):

- **PORT**: HTTP port (default: `8080`)
- **DEFAULT_TTS_PROVIDER**: default TTS provider key (e.g. `elevenlabs`)
- **DEFAULT_STT_PROVIDER**: default STT provider key (e.g. `elevenlabs`)
- **ELEVEN_LABS_API_KEY**: your ElevenLabs API key

Example `.env`:

```env
PORT=8080
DEFAULT_TTS_PROVIDER=elevenlabs
DEFAULT_STT_PROVIDER=elevenlabs
ELEVEN_LABS_API_KEY=your_api_key_here
```

## Run

```bash
go run ./cmd/server
```

Server starts on `:${PORT}` (default `:8080`).

## HTTP API

### POST /tts

- **Body (JSON)**:
  - `text` (string, required): text to synthesize
  - `voice` (string, required): ElevenLabs voice ID
  - `format` (string, optional): desired audio format (adapter returns `mp3`)
  - `provider` (string, optional): provider key; defaults to `DEFAULT_TTS_PROVIDER`
- **Response**: `audio/mpeg` bytes

Example:

```bash
curl -X POST http://localhost:8080/tts \
  -H "Content-Type: application/json" \
  --data '{
    "text": "Hello from polyglot",
    "voice": "YOUR_ELEVENLABS_VOICE_ID"
  }' \
  --output out.mp3
```

### POST /stt?format=mp3

- **Query**: `format` (optional; default `mp3`)
- **Body**: raw audio bytes
- **Response**: `text/plain` with the transcription
  - Note: provider is currently fixed to ElevenLabs in the HTTP adapter

Example (transcribe an mp3 file):

```bash
curl -X POST "http://localhost:8080/stt?format=mp3" \
  --data-binary @sample.mp3
```

## Extending the system

- Add a new outbound adapter that implements the domain port(s):
  - `TTSProvider.Synthesize(ctx, req)`
  - `STTProvider.Transcribe(ctx, req)`
- Register the adapter in `cmd/server/main.go` by adding it to the providers map(s)
- Use its key via `DEFAULT_TTS_PROVIDER` / `DEFAULT_STT_PROVIDER` or per-request `provider` (currently per-request is supported by TTS handler)

This keeps the domain agnostic of specific vendors while allowing easy swaps and multi-provider support.

## Notes

- Responses and error handling are intentionally minimal for clarity
- The TTS adapter accepts `voice` as the ElevenLabs voice ID and requests `audio/mpeg`
- The STT adapter uses ElevenLabs `scribe_v1` model


