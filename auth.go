package audact

import (
	"context"
	"net/http"
)

// ExchangeKey exchanges the API key for a short-lived access token. Calls
// POST /v1/auth/exchange-key.
func (s *AuthService) ExchangeKey(ctx context.Context) (*ExchangeKeyResponse, error) {
	var out ExchangeKeyResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/auth/exchange-key", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Refresh refreshes an access token. Calls POST /v1/auth/refresh.
func (s *AuthService) Refresh(ctx context.Context) (*RefreshTokenResponse, error) {
	var out RefreshTokenResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/auth/refresh", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
