package audact

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// List lists integration connections. Calls GET /v1/integrations.
func (s *IntegrationsService) List(ctx context.Context, params *ListConnectionsParams) (*ListConnectionsResponse, error) {
	q := url.Values{}
	if params != nil {
		addString(q, "provider", params.Provider)
	}
	var out ListConnectionsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/integrations", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates an integration connection. Calls POST /v1/integrations.
func (s *IntegrationsService) Create(ctx context.Context, req *CreateConnectionRequest) (*Connection, error) {
	var out Connection
	if _, err := s.t.do(ctx, http.MethodPost, "/integrations", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns a connection. Calls GET /v1/integrations/{connection_id}.
func (s *IntegrationsService) Get(ctx context.Context, connectionID string) (*Connection, error) {
	var out Connection
	if _, err := s.t.do(ctx, http.MethodGet, "/integrations/"+url.PathEscape(connectionID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a connection; force deletes it even when it is in use.
// Calls DELETE /v1/integrations/{connection_id}.
func (s *IntegrationsService) Delete(ctx context.Context, connectionID string, force bool) (*DeleteConnectionResponse, error) {
	q := url.Values{}
	addBool(q, "force", force)
	var out DeleteConnectionResponse
	if _, err := s.t.do(ctx, http.MethodDelete, "/integrations/"+url.PathEscape(connectionID), authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListEventTypes lists the event types of a connection.
// Calls GET /v1/integrations/{connection_id}/event-types.
func (s *IntegrationsService) ListEventTypes(ctx context.Context, connectionID string) (*ListEventTypesResponse, error) {
	var out ListEventTypesResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/integrations/"+url.PathEscape(connectionID)+"/event-types", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Test checks a connection's health. Calls POST /v1/integrations/{connection_id}/test.
func (s *IntegrationsService) Test(ctx context.Context, connectionID string) (*TestConnectionResponse, error) {
	var out TestConnectionResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/integrations/"+url.PathEscape(connectionID)+"/test", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CalcomOAuthStart mints a state token and returns the Cal.com authorize URL.
// Calls GET /v1/integrations/cal_com/oauth/start.
func (s *IntegrationsService) CalcomOAuthStart(ctx context.Context, params *CalcomOAuthStartParams) (*CalcomOAuthStartResponse, error) {
	q := url.Values{}
	if params != nil {
		addString(q, "name", params.Name)
		addString(q, "cal_username", params.CalUsername)
	}
	var out CalcomOAuthStartResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/integrations/cal_com/oauth/start", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CalcomOAuthCallback calls the Cal.com OAuth callback without following the
// redirect and returns the status and Location. It sends no Authorization
// header: the single-use state token is the credential. Calls
// GET /v1/integrations/cal_com/oauth/callback.
func (s *IntegrationsService) CalcomOAuthCallback(ctx context.Context, params CalcomOAuthCallbackParams) (*CalcomOAuthCallbackResponse, error) {
	q := url.Values{}
	q.Set("state", params.State)
	addString(q, "code", params.Code)
	addString(q, "error", params.Error)
	const path = "/integrations/cal_com/oauth/callback"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.t.baseURL+"/v1"+path+"?"+q.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("audact: GET %s: %w", path, err)
	}
	req.Header.Set("User-Agent", s.t.userAgent)
	req.Header.Set("Accept", "application/json")
	hc := *s.t.hc
	hc.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	resp, err := hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("audact: GET %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
		return nil, parseError(resp, b)
	}
	return &CalcomOAuthCallbackResponse{StatusCode: resp.StatusCode, Location: resp.Header.Get("Location")}, nil
}
