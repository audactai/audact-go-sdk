package audact

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// List lists tools. Calls GET /v1/tools.
func (s *ToolsService) List(ctx context.Context, params *ListToolsParams) (*ListToolsResponse, error) {
	q := url.Values{}
	if params != nil {
		addString(q, "name", params.Name)
		addString(q, "cursor", params.Cursor)
		addInt(q, "limit", params.Limit)
	}
	var out ListToolsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/tools", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a tool; a non-empty idempotencyKey is sent as the
// Idempotency-Key header. Calls POST /v1/tools.
func (s *ToolsService) Create(ctx context.Context, req *CreateToolRequest, idempotencyKey string) (*Tool, error) {
	var out Tool
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/tools", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns a tool. Calls GET /v1/tools/{tool_id}.
func (s *ToolsService) Get(ctx context.Context, toolID string) (*Tool, error) {
	var out Tool
	if _, err := s.t.do(ctx, http.MethodGet, "/tools/"+url.PathEscape(toolID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Replace replaces a tool. Calls PUT /v1/tools/{tool_id}.
func (s *ToolsService) Replace(ctx context.Context, toolID string, req *CreateToolRequest) (*Tool, error) {
	var out Tool
	if _, err := s.t.do(ctx, http.MethodPut, "/tools/"+url.PathEscape(toolID), authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update partially updates a tool; unset fields are not sent. Calls PATCH /v1/tools/{tool_id}.
func (s *ToolsService) Update(ctx context.Context, toolID string, req *UpdateToolRequest) (*Tool, error) {
	var out Tool
	if _, err := s.t.do(ctx, http.MethodPatch, "/tools/"+url.PathEscape(toolID), authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a tool; force detaches it from agents first. Calls DELETE /v1/tools/{tool_id}.
func (s *ToolsService) Delete(ctx context.Context, toolID string, force bool) (*DeleteToolResponse, error) {
	q := url.Values{}
	if force {
		q.Set("force", strconv.FormatBool(force))
	}
	var out DeleteToolResponse
	if _, err := s.t.do(ctx, http.MethodDelete, "/tools/"+url.PathEscape(toolID), authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Test runs a health check on a tool. Calls POST /v1/tools/{tool_id}/test.
func (s *ToolsService) Test(ctx context.Context, toolID string) (*TestToolResponse, error) {
	var out TestToolResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/tools/"+url.PathEscape(toolID)+"/test", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListForAgent lists tools attached to an agent. Calls GET /v1/agents/{agent_id}/tools.
func (s *ToolsService) ListForAgent(ctx context.Context, agentID string, params *ListAgentToolsParams) (*ListAgentToolsResponse, error) {
	q := url.Values{}
	if params != nil {
		addString(q, "cursor", params.Cursor)
		addInt(q, "limit", params.Limit)
	}
	var out ListAgentToolsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/agents/"+url.PathEscape(agentID)+"/tools", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AttachToAgent attaches a tool to an agent. Calls POST /v1/agents/{agent_id}/tools.
func (s *ToolsService) AttachToAgent(ctx context.Context, agentID string, req *AttachToolRequest) (*AttachToolResponse, error) {
	var out AttachToolResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/agents/"+url.PathEscape(agentID)+"/tools", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DetachFromAgent detaches a tool from an agent. Calls DELETE /v1/agents/{agent_id}/tools/{tool_id}.
func (s *ToolsService) DetachFromAgent(ctx context.Context, agentID, toolID string) (*DetachToolResponse, error) {
	var out DetachToolResponse
	if _, err := s.t.do(ctx, http.MethodDelete, "/agents/"+url.PathEscape(agentID)+"/tools/"+url.PathEscape(toolID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
