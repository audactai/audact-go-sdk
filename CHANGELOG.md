# Changelog

## [Unreleased]

## [0.4.0] - 2026-10-28

### ⚠️ Breaking

- `CallRecord.Metadata` and `CallDetail.Metadata`: `string` → `map[string]any`.
  Read it as a map; no more `json.Unmarshal([]byte(rec.Metadata), …)`.
- `ReconstructedEvent.Data`: `string` → `json.RawMessage`. Unmarshal it
  directly with `json.Unmarshal(ev.Data, &v)`.
- `ViolationRecord.Data`: `string` → `json.RawMessage`. Unmarshal it directly.
- `PostEventRequest.Data`: `string` → `json.RawMessage`. Pass valid JSON, e.g.
  `json.RawMessage(`{"k":1}`)`; it is sent as an object, not a quoted string.
- `PolicyRuleCondition.Value`: `string` → `any`. After decoding it is a
  `string`, `float64`, `bool` or `[]any` of strings; use a type switch.
- `CreateWebhookRequest.Scope`, `WebhookEndpoint.Scope` and the other scope
  fields: `string` → `WebhookScope`. Untyped string constants still compile;
  convert variables with `WebhookScope(s)`. `Scope` is now required on create.
- Response booleans (`IsActive`, `Active`, `Enabled`, `EncryptionEnabled`,
  `PIIEncryptionEnabled`) no longer use `omitempty`: `false` is now always
  encoded when a response value is marshalled again. The request field
  `UpdateWebhookRequest.IsActive` is unchanged: when nil it is not sent.

### Fixed

- `Calls.List`, `Calls.Get`, `Evidence.Reconstruct` and `Policy.Get` failed to
  decode real API responses because fields typed as strings are JSON objects
  or booleans on the wire.
- `Evidence.ListViolations` would fail to decode as soon as a violation existed.
- `Evidence.PostEvent` sent `data` as a quoted string instead of a JSON value.
- `Webhooks.Create`: `scope` is now a required, typed field.
- `Policy.CheckCallingHours` now works because the API routes the endpoint
  (no SDK change).

### Added

- `WebhookScope` with `WebhookScopeAccount`, `WebhookScopeAgent` and
  `WebhookScopeCompliance`.

### Deprecated

- `Validator.FinalizeScore`: the endpoint is internal and is not routed on the
  public API, so it returns 404. It will be removed in a future release.

## [0.3.0] - 2026-10-23

### Added

- `TTSModelElevenV4Turbo` (`eleven_v4_turbo`) and
  `TTSModelElevenV3Conversational` (`eleven_v3_conversational`).

- `KBSource.ChunkCount`, `KBSource.CrawlCoverage` (new `KBCrawlCoverage`:
  discovered / stored / excluded / failed pages) and `KBSource.LastIngestedAt`:
  the outcome of the source's latest successful ingestion. Nil /
  empty until a first ingestion has finished.

- `ListVoicesParams.Limit` (1-1000; 0 = server default 100 per provider) and
  `VoiceProvider.Total` / `VoiceProvider.Truncated`.

- ElevenLabs languages: `Voice.Languages` (what a voice speaks; `Language` is
  its native one), `VoiceProvider.Models` / `TTSModelInfo` (each model and its
  languages), and `Agents.ListVoicesFor(ctx, &ListVoicesParams{Language,
  Provider})` — native speakers first. `ListVoices` is unchanged.
- ElevenLabs TTS: `TTSProviderElevenLabs`; new `TTSModel` type with
  `TTSModelDefault` (`""`, provider default), `TTSModelElevenFlashV25`,
  `TTSModelElevenTurboV25`, `TTSModelElevenMultilingualV2` and
  `TTSModelVoxtralMiniTTS2603`; `tts_model` on `CreateAgentRequest`
  (and Update/Patch), `Agent` and `CreateAgentResponse`. `ListVoices` returns an
  `elevenlabs` group. The voice must belong to the chosen provider (400,
  `field: voice_id`); `voice_id` stays a plain string because tenant-cloned
  voices are valid.
- `Client.SetPlatformToken` replaces the platform token. Platform tokens are
  single-use and action-scoped, so set a fresh one before each Platform call.
- `ErrPlatformTokenRequired`.
- `ValidateAudioRequest.ConfigSealHash` and
  `TurnDeliveryStatus.ConfirmationCount` close schema drift.
- `ChatCompletionChunk` and `Stream.Chunk()` / `Stream.Raw()`.
- `Tenants.GetConfigurationVersion(ctx, configID, version)` for
  `GET /v1/configurations/{config_id}/versions/{version}`: one immutable
  earlier version of a configuration.
- `AudioComplianceReceipt` gains `CallID`, `TurnNumber`, `Timestamp`,
  `ConfigSealHash`, `AuthMechanism` and `ReceiptFormatVersion`.
- `ResolvedConfigVersion` carries the version's settings: `BeginMessage`,
  `BlocklistPatterns`, `ConsentConfig`, `Countries`, `DisclosureTexts`,
  `EndCall`, `FallbackResponse`, `LLMModel`, `LLMProvider`,
  `OutputConstraint`, `PricingBands`, `PromptConfidentialityConfig`,
  `ValidatorModes`, `Voicemail`.
- `UsageRecord.Source` (`UsageSourceStored` / `UsageSourceLive`).

- `AudioSegmentResult.TranscriptionStatus` with the `TranscriptionStatus` type
  and constants `TranscriptionStatusTranscribed` and `TranscriptionStatusNoSpeech`.
- `Evidence.ListC2PAManifests` (`GET /v1/calls/{call_id}/c2pa`),
  `Evidence.DownloadC2PAManifestStore` (`application/c2pa` bytes),
  `Evidence.DownloadC2PAAsset` (follows the 302 to the presigned recording URL
  without sending the API key) and `Evidence.C2PATrustAnchor` (PEM,
  unauthenticated).
- Types: `C2PAAssetRef`, `C2PAManifestStoreRef`, `C2PASigningCredential`,
  `C2PAListSigner`, `C2PAManifestListItem`, `C2PAManifestListResponse`,
  `ListC2PAManifestsParams`.
- `c.KnowledgeBases`, covering the 20 `/v1/knowledge-bases` operations:
  `List`, `Create`, `Get`, `Update`, `Delete`; sources `ListSources`,
  `AddSource` (JSON url/text, 202), `GetSource`, `UpdateSource`,
  `ReingestSource` (202), `DeleteSource`; agent links `ListAgentLinks`,
  `ListLinksByAgent`, `LinkAgent`, `UpdateAgentLink`, `UnlinkAgent`; crawl
  exclusions `ListExclusions`, `AddExclusion`, `DeleteExclusion`; and
  `TestRetrieval`. `AddFileSource` uploads a file source as
  multipart/form-data on the same endpoint.
- Types: `KnowledgeBase`, `KBSource`, `KBSourceType`, `KBSourceStatus`,
  `KBAgentLink`, `KBCrawlExclusion`, `KBRetrievalTurn`, `KBRetrievalHit` and
  the request/response types of the methods above.

### ⚠️ Breaking

- `Evidence.Export` and `Evidence.Article50Report` return `*Article50Export`,
  which holds either `Report` (JSON) or `PDF` (raw `application/pdf` bytes).
  Platform-minted evidence links default to PDF, which the old JSON-only
  return type could not decode.
- `Evidence.GenerateC2PAManifest(ctx, callID, idempotencyKey)` takes an
  idempotency key. Pass `""` to keep the old behaviour. Every manifest is
  stored permanently, so retries should carry a key.
- Platform methods no longer fall back to the tenant API key. Without a
  platform token they return `ErrPlatformTokenRequired` and send nothing.

- C2PA manifests are really signed (generator 2.0.0). `C2PAManifest.Signature` (HMAC hex) is removed; `SignatureAlgorithm` is now
  `"ES256"`. Verify the downloaded manifest store with `c2patool` instead.
- `C2PAManifest` gains `ManifestID`, `CallID`, `Asset`, `ManifestStore` and
  `SigningCredential`.

### Changed

- `KnowledgeBase.RobotsPolicy` and `CreateKnowledgeBaseRequest.RobotsPolicy` are now
  `KBRobotsPolicy` (`KBRobotsPolicyHonor`, `KBRobotsPolicyOverride`), and
  `KBRetrievalTurn.Role` is now `KBRetrievalTurnRole` (`KBRetrievalTurnRoleUser`,
  `KBRetrievalTurnRoleAssistant`, `KBRetrievalTurnRoleSystem`). Untyped string
  constants still compile; typed `string` variables need a conversion.

### Fixed

- `LLMProxy.Stream` decodes the proxy's real `chat.completion.chunk` events.
  `Current()` returns `choices[0].delta.content`, and the per-sentence
  `audact_validation` verdict is available on `Chunk()`.
- Transport errors no longer include request query strings, so an evidence
  capability token cannot leak into a logged error.
- The smoke example treats `terminated` as a terminal status, instead of
  polling until its deadline.

### Overview

Initial release: hand-written, standard library only Go SDK covering all 184
non-ignored operations of the Audact API, with typed errors (`APIError`,
`RateLimitError`, sentinel errors), `context.Context` on every call, a single
configurable base URL, SSE streaming for the LLM proxy, and PII normalisation
helpers.
