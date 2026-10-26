package audact

import (
	"context"
	"net/http"
	"net/url"
)

// BYOConfig returns deployment-wide BYO trunk settings. Calls GET /v1/byo/config.
func (s *NumbersService) BYOConfig(ctx context.Context) (*BYOConfig, error) {
	var out BYOConfig
	if _, err := s.t.do(ctx, http.MethodGet, "/byo/config", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListBYOConnections lists BYO connections. Calls GET /v1/byo/connections.
func (s *NumbersService) ListBYOConnections(ctx context.Context, params *ListBYOConnectionsParams) (*ListBYOConnectionsResponse, error) {
	q := url.Values{}
	if params != nil {
		addInt(q, "limit", params.Limit)
		addString(q, "cursor", params.Cursor)
	}
	var out ListBYOConnectionsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/byo/connections", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateBYOConnection creates a BYO connection; a non-empty idempotencyKey is
// sent as the Idempotency-Key header. Calls POST /v1/byo/connections.
func (s *NumbersService) CreateBYOConnection(ctx context.Context, req *CreateBYOConnectionRequest, idempotencyKey string) (*BYOConnection, error) {
	var out BYOConnection
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/byo/connections", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetBYOConnection returns a BYO connection. Calls GET /v1/byo/connections/{id}.
func (s *NumbersService) GetBYOConnection(ctx context.Context, id string) (*BYOConnection, error) {
	var out BYOConnection
	if _, err := s.t.do(ctx, http.MethodGet, "/byo/connections/"+url.PathEscape(id), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateBYOConnection updates a BYO connection; a non-empty idempotencyKey is
// sent as the Idempotency-Key header. Calls PUT /v1/byo/connections/{id}.
func (s *NumbersService) UpdateBYOConnection(ctx context.Context, id string, req *UpdateBYOConnectionRequest, idempotencyKey string) (*BYOConnection, error) {
	var out BYOConnection
	if _, err := s.t.doHeaders(ctx, http.MethodPut, "/byo/connections/"+url.PathEscape(id), authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteBYOConnection deletes a BYO connection. Calls DELETE /v1/byo/connections/{id}.
func (s *NumbersService) DeleteBYOConnection(ctx context.Context, id string) error {
	_, err := s.t.do(ctx, http.MethodDelete, "/byo/connections/"+url.PathEscape(id), authAPIKey, nil, nil, nil)
	return err
}

// VerifyBYOConnection probes a BYO connection and returns its updated health.
// Calls POST /v1/byo/connections/{id}/verify.
func (s *NumbersService) VerifyBYOConnection(ctx context.Context, id string) (*BYOConnection, error) {
	var out BYOConnection
	if _, err := s.t.do(ctx, http.MethodPost, "/byo/connections/"+url.PathEscape(id)+"/verify", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
