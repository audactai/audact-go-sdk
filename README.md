# Audact Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/audactai/audact-go-sdk.svg)](https://pkg.go.dev/github.com/audactai/audact-go-sdk)

Hand-written, standard-library-only Go client for the Audact API. It covers the
full Audact API, with typed errors and `context.Context` on every call.

## Install

```sh
go get github.com/audactai/audact-go-sdk@latest
```

To pin this release:

```sh
go get github.com/audactai/audact-go-sdk@v0.4.0
```

Requires Go 1.24 or newer. The module has no third-party dependencies.

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	audact "github.com/audactai/audact-go-sdk"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c, err := audact.NewClient("") // reads AUDACT_API_KEY
	if err != nil {
		log.Fatal(err)
	}
	agents, err := c.Agents.List(ctx, &audact.ListAgentsParams{Limit: 10})
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range agents.Agents {
		fmt.Println(a.AgentID, a.Name)
	}
}
```

A fuller example lives in [`examples/quickstart`](examples/quickstart), and
[`examples/smoke`](examples/smoke) is an end-to-end staging smoke test.

## Base URL

The default is `https://api.audact.ai`. Resolution order:

1. `audact.WithBaseURL("https://staging.example")`
2. the `AUDACT_BASE_URL` environment variable
3. `audact.DefaultBaseURL`

A trailing `/` and a trailing `/v1` are stripped; every request path is
prefixed with `/v1`.

## Authentication

| Kind | How to set it | Header |
|---|---|---|
| Tenant API key | `NewClient(key)` or `AUDACT_API_KEY` | `Authorization: Bearer <key>` |
| Org token | `audact.WithOrgToken(t)`, or `c.Org.Login` / `Signup` (stored automatically). `c.Org.Refresh` reads the refresh cookie, so it needs an `http.Client` with a cookie jar (`WithHTTPClient`) | `Authorization: Bearer <token>` |
| Platform token | `audact.WithPlatformToken(t)`, then `c.SetPlatformToken(t)` before each later call | `Authorization: Bearer <token>` |

`NewClient` does not fail for a missing key because some operations are
public. It returns an error only for an invalid base URL or a nil option. If a
needed credential is empty the header is omitted and the server's 401 surfaces
as `audact.ErrUnauthorized`.

**Platform tokens are single-use and action-scoped.** The server burns each
token's `jti` on first use and checks its `action`: `org` for `c.Platform.Orgs`,
`admin_keys` for the admin-key methods, `export` / `import` for bundles. Mint
a fresh token for every Platform call (for example
`audact-admin auth --action admin_keys`) and set it with `c.SetPlatformToken`.
A Platform method called with no platform token returns
`audact.ErrPlatformTokenRequired` without sending a request; the tenant API
key is never sent to a platform route.

## Errors

Non-2xx responses return `*audact.APIError`, which carries `StatusCode`,
`Code`, `Message`, `Hint`, `Details` and `RequestID`.

```go
_, err := c.Calls.Get(ctx, "call_123")
var apiErr *audact.APIError
switch {
case errors.Is(err, audact.ErrNotFound):
	fmt.Println("no such call")
case errors.As(err, &apiErr):
	fmt.Println(apiErr.StatusCode, apiErr.Code, "request_id:", apiErr.RequestID)
case err != nil:
	log.Fatal(err) // transport error, wrapped with %w
}
```

Sentinels: `ErrUnauthorized` (401), `ErrForbidden` (403), `ErrNotFound` (404),
`ErrConflict` (409), `ErrRateLimited` (429). A 429 returns
`*audact.RateLimitError`; use `RetryAfter` for the server's requested delay:

```go
var rl *audact.RateLimitError
if errors.As(err, &rl) {
	time.Sleep(rl.RetryAfter)
}
```

## Retries and idempotency

The SDK never retries on its own. Writes that create evidence or other
permanent records accept an idempotency key: pass the same key when you retry
and the server replays the first response instead of acting twice. An empty
key sends no `Idempotency-Key` header.

```go
m, err := c.Evidence.GenerateC2PAManifest(ctx, callID, "c2pa-"+callID)
```

Methods with an idempotency key include `Evidence.PostEvent`,
`Evidence.ConfirmDelivery`, `Evidence.ReconcileTranscript`,
`Evidence.GenerateC2PAManifest` and `Platform.ImportTenant`
(`PlatformImportOptions.IdempotencyKey`). C2PA manifests and evidence events
are append-only, so a retry without a key stores a duplicate permanently.

## Context and timeouts

Every method takes a `context.Context` first; cancellation and deadlines are
honored, and `errors.Is(err, context.DeadlineExceeded)` works. The default HTTP
client has a 30 s timeout. Supply your own with `WithHTTPClient`:

```go
c, err := audact.NewClient("", audact.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}))
```

`WithUserAgent("myapp/1.0")` appends to the default `audact-go/0.3.0` agent.

## Streaming

`c.LLMProxy.Stream` returns a `*Stream` over server-sent events. The proxy
validates the reply one sentence at a time: `Current()` is the sentence text,
and `Chunk().AudactValidation` carries its verdict (`MODIFY` means the text
was rewritten; a `BLOCK` chunk carries the fallback line and ends the stream).
Tool-call deltas arrive as chunks with no `AudactValidation`.

```go
st, err := c.LLMProxy.Stream(ctx, nil, &audact.ChatCompletionRequest{
	Messages: []audact.ChatMessage{{Role: "user", Content: "Hello"}},
})
if err != nil {
	log.Fatal(err)
}
defer st.Close()
for st.Next() {
	if v := st.Chunk().AudactValidation; v != nil && v.Verdict == "BLOCK" {
		log.Println("reply blocked by compliance check")
	}
	fmt.Print(st.Current())
}
if err := st.Err(); err != nil {
	log.Fatal(err)
}
```

## PII helpers

`audact.NormalizeIdentifier(typ, value)` canonicalises a name, email or phone
identifier, and `audact.ClientHash(typ, value)` computes the client-side hash.
Both match the TypeScript and Python SDKs on the shared test vectors.

```go
h, err := audact.ClientHash("email", "Jane@Example.com")
```

## Choosing a TTS provider

Agents speak through Voxtral by default. To use ElevenLabs, set the provider,
optionally a model, and a voice from the provider's group in `ListVoices`:

```go
voices, _ := c.Agents.ListVoices(ctx)
voice := voices["elevenlabs"].Voices[0].VoiceID

_, err := c.Agents.Create(ctx, &audact.CreateAgentRequest{
	Name:        "Concierge",
	TTSProvider: audact.TTSProviderElevenLabs,
	TTSModel:    audact.TTSModelElevenFlashV25, // TTSModelDefault = provider default
	VoiceID:     voice,
})
```

A voice from another provider is rejected with a 400 naming `voice_id`.

## Services

| Accessor | Area |
|---|---|
| `c.Agents` | Agents and voices |
| `c.Calls` | Call records, SIP call control, transfers, routes |
| `c.Numbers` | Phone numbers, KYC, BYO connections |
| `c.Webhooks` | Webhooks and deliveries |
| `c.Tools` | Agent tools |
| `c.Integrations` | Third-party integrations |
| `c.KnowledgeBases` | Knowledge bases (RAG): sources, agent links, crawl exclusions, retrieval tests |
| `c.Evidence` | Evidence chains, receipts, compliance events, C2PA content credentials |
| `c.PII` | PII shredding, identities, exemptions |
| `c.Validator` | Validation and scoring |
| `c.LLMProxy` | LLM proxy (including streaming) |
| `c.Tenants` | Tenants, keys, configurations |
| `c.Org` | Organization accounts |
| `c.Platform` | Platform export/import, admin keys; `c.Platform.Orgs` for org administration |
| `c.Aggregation` | Aggregation proofs and metachain |
| `c.Policy` | Policy evaluation and calling hours |
| `c.Audio` | Audio validation and reconciliation |
| `c.Auth` | API key exchange and refresh |

## License

MIT

## Versioning

This module follows semantic versioning. Before 1.0 a minor release may contain
breaking changes; the [CHANGELOG](CHANGELOG.md) marks them with ⚠️. A future v2
will use the module path `github.com/audactai/audact-go-sdk/v2`.

## Contributing

This repository is published from Audact's private source repository; every commit here is an automated sync. Please open issues here. Pull requests can't be merged directly, but they are welcome as proposals and will be ported by the maintainers.
