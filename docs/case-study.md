## Case Study: Polyglot — Applying Hexagonal Architecture to a TTS/STT Service

### Background
An internal proof-of-concept AI assistant platform had evolved into a monolith. Before adding more sophisticated capabilities, a focused refactor of Text-to-Speech (TTS) and Speech-to-Text (STT) was undertaken to explore a maintainable architecture pattern. Hexagonal architecture (ports and adapters) was selected to improve separation of concerns and enable straightforward provider substitution. The initial audience was a single developer, with an eye toward future production readiness. There were no strict time or budget constraints; the emphasis was on learning and correctness of the pattern.

### Objectives
- **Separation of concerns** and clear boundaries between domain logic and infrastructure
- **Provider swapability** so additional vendors can be added with minimal changes
- **Testability and maintainability** as secondary objectives

Success was defined as the ability to add new providers without rewiring major components or navigating large, tightly coupled files.

### Scope
- **In scope**: minimal TTS and STT functionality with simple, predictable inputs/outputs
- **Out of scope**: broader orchestration and assistant features, which belong in higher-level services

### Architecture Overview
- **Domain core** (`internal/domain/voice`)
  - Ports (interfaces): `TTSProvider`, `STTProvider`
  - Services: `TTSService`, `STTService` select a provider by key with a configured default
  - Types: `SpeechRequest/Result`, `TranscriptionRequest/Result`
- **Inbound adapter** (`internal/adapters/http`)
  - `POST /tts`: JSON in, `audio/mpeg` out
  - `POST /stt`: binary audio in, plain text out
- **Outbound adapter** (`internal/adapters/elevenlabs`)
  - `ElevenLabsTTS` and `ElevenLabsSTT` invoke ElevenLabs REST APIs
- **Composition root** (`cmd/server/main.go`)
  - Wires providers, reads configuration, and mounts routes with `chi`
- **Configuration** (`internal/config/config.go`)
  - Loads `.env`, validates defaults and API keys

### Notable Design Choices
- Providers are injected as maps keyed by name; defaults are set via environment configuration
- HTTP handlers are intentionally thin and delegate to services; services encapsulate provider dispatch
- TTS supports per-request provider override; STT currently fixes the provider to ElevenLabs within the handler for simplicity
- Error and response handling are kept minimal to emphasize architectural structure

### Provider Selection
ElevenLabs was selected initially for both TTS and STT due to available credentials, allowing rapid validation of the architectural approach.

### Implementation Highlights
- Clear port/adapter boundaries enable adding a new provider by implementing the domain interfaces and registering the adapter at the composition root
- A simple HTTP surface reduces cognitive overhead and eases manual and automated testing
- The domain remains vendor-agnostic; provider-specific details are isolated within adapters

### Operational Approach
Local development relies on a `.env` file for configuration and the Air tool for hot reloading.

### Security & Privacy
Secrets are sourced from a local `.env` file during development; no production-grade secret management is included in this exploration.

### Outcomes
The refactor validated hexagonal architecture for this problem space, delivering clean separation and straightforward extensibility. The resulting service behaves as intended and provides a solid foundation for reintegration into a broader assistant platform.

### Future Work
- Reintegrate the service into the main platform following wider refactors
- Add multiple providers, streaming support, authentication, and improved error handling
- Consider infrastructure and developer-experience improvements during reintegration

### Team & Timeline
This was a solo effort completed in approximately 6–8 hours.