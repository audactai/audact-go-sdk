package audact

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// List lists agents. Calls GET /v1/agents.
func (s *AgentsService) List(ctx context.Context, params *ListAgentsParams) (*ListAgentsResponse, error) {
	q := url.Values{}
	if params != nil {
		addInt(q, "limit", params.Limit)
		addString(q, "cursor", params.Cursor)
		addString(q, "configuration_id", params.ConfigurationID)
	}
	var out ListAgentsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/agents", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates an agent. Calls POST /v1/agents.
func (s *AgentsService) Create(ctx context.Context, req *CreateAgentRequest) (*CreateAgentResponse, error) {
	var out CreateAgentResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/agents", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns an agent. Calls GET /v1/agents/{agent_id}.
func (s *AgentsService) Get(ctx context.Context, agentID string) (*Agent, error) {
	var out Agent
	if _, err := s.t.do(ctx, http.MethodGet, "/agents/"+url.PathEscape(agentID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update replaces an agent's settings. Calls PUT /v1/agents/{agent_id}.
func (s *AgentsService) Update(ctx context.Context, agentID string, req *UpdateAgentRequest) (*Agent, error) {
	var out Agent
	if _, err := s.t.do(ctx, http.MethodPut, "/agents/"+url.PathEscape(agentID), authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Patch partially updates an agent; unset fields are not sent. Calls PATCH /v1/agents/{agent_id}.
func (s *AgentsService) Patch(ctx context.Context, agentID string, req *PatchAgentRequest) (*Agent, error) {
	var out Agent
	if _, err := s.t.do(ctx, http.MethodPatch, "/agents/"+url.PathEscape(agentID), authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes an agent. Calls DELETE /v1/agents/{agent_id}.
func (s *AgentsService) Delete(ctx context.Context, agentID string) (*DeleteAgentResponse, error) {
	var out DeleteAgentResponse
	if _, err := s.t.do(ctx, http.MethodDelete, "/agents/"+url.PathEscape(agentID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Duplicate duplicates an agent. Calls POST /v1/agents/{agent_id}/duplicate.
func (s *AgentsService) Duplicate(ctx context.Context, agentID string) (*Agent, error) {
	var out Agent
	if _, err := s.t.do(ctx, http.MethodPost, "/agents/"+url.PathEscape(agentID)+"/duplicate", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListVoices lists TTS voices grouped by provider. Calls GET /v1/voices.
func (s *AgentsService) ListVoices(ctx context.Context) (ListVoicesResponse, error) {
	return s.ListVoicesFor(ctx, nil)
}

// ListVoicesFor lists TTS voices with optional filters: Language keeps voices
// (and models) that speak it, native speakers first; Provider keeps one TTS
// provider. Calls GET /v1/voices.
func (s *AgentsService) ListVoicesFor(ctx context.Context, params *ListVoicesParams) (ListVoicesResponse, error) {
	q := url.Values{}
	if params != nil {
		addString(q, "language", params.Language)
		addString(q, "provider", string(params.Provider))
		addInt(q, "limit", params.Limit)
	}
	var out ListVoicesResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/voices", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateVoice clones a voice from a reference clip using multipart/form-data.
// Calls POST /v1/voices.
func (s *AgentsService) CreateVoice(ctx context.Context, req *CreateVoiceRequest) (*CreateVoiceResponse, error) {
	if req == nil || req.Audio == nil {
		return nil, errors.New("audact: CreateVoice: audio is required")
	}
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	if err := mw.WriteField("name", req.Name); err != nil {
		return nil, fmt.Errorf("audact: CreateVoice: %w", err)
	}
	filename := req.Filename
	if filename == "" {
		filename = "audio"
	}
	fw, err := mw.CreateFormFile("audio", filename)
	if err != nil {
		return nil, fmt.Errorf("audact: CreateVoice: %w", err)
	}
	if _, err := io.Copy(fw, req.Audio); err != nil {
		return nil, fmt.Errorf("audact: CreateVoice: read audio: %w", err)
	}
	consent := "false"
	if req.ConsentConfirmed {
		consent = "true"
	}
	fields := []struct{ k, v string }{
		{"consent_confirmed", consent}, {"language", req.Language}, {"accent", req.Accent},
		{"gender", req.Gender}, {"style", req.Style},
	}
	for _, f := range fields {
		if f.v == "" {
			continue
		}
		if err := mw.WriteField(f.k, f.v); err != nil {
			return nil, fmt.Errorf("audact: CreateVoice: %w", err)
		}
	}
	if err := mw.Close(); err != nil {
		return nil, fmt.Errorf("audact: CreateVoice: %w", err)
	}
	var out CreateVoiceResponse
	body := &rawBody{contentType: mw.FormDataContentType(), data: buf.Bytes()}
	if _, err := s.t.do(ctx, http.MethodPost, "/voices", authAPIKey, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetVoice returns a voice. Calls GET /v1/voices/{voice_id}.
func (s *AgentsService) GetVoice(ctx context.Context, voiceID string) (*Voice, error) {
	var out Voice
	if _, err := s.t.do(ctx, http.MethodGet, "/voices/"+url.PathEscape(voiceID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteVoice deletes a cloned voice. Calls DELETE /v1/voices/{voice_id}.
func (s *AgentsService) DeleteVoice(ctx context.Context, voiceID string) error {
	_, err := s.t.do(ctx, http.MethodDelete, "/voices/"+url.PathEscape(voiceID), authAPIKey, nil, nil, nil)
	return err
}

// PipelineStatus returns the number of active voice pipelines. This endpoint
// is public in the spec, but like the TS SDK the client sends the API key when
// one is configured. Calls GET /v1/pipeline/status.
func (s *AgentsService) PipelineStatus(ctx context.Context) (*PipelineStatusResponse, error) {
	var out PipelineStatusResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/pipeline/status", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// orEmpty returns an empty struct for a nil request so the body is "{}".
func orEmpty[T any](req *T) *T {
	if req == nil {
		return new(T)
	}
	return req
}
