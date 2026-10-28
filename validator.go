package audact

import (
	"context"
	"net/http"
	"net/url"
)

// Validate validates one turn of agent output. Calls POST /v1/validate.
func (s *ValidatorService) Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error) {
	var out ValidateResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/validate", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FinalizeScore computes the aggregate compliance score for a completed call.
// Calls POST /v1/calls/{call_id}/finalize-score.
//
// Deprecated: POST /v1/calls/{call_id}/finalize-score is an internal endpoint
// and is not routed on the public API, so this method returns 404 against
// api.audact.ai. It will be removed in a future release.
func (s *ValidatorService) FinalizeScore(ctx context.Context, callID string, req *FinalizeScoreRequest) (*FinalizeScoreResponse, error) {
	var out FinalizeScoreResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/calls/"+url.PathEscape(callID)+"/finalize-score", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Health returns the Validator health. Calls GET /v1/health.
func (s *ValidatorService) Health(ctx context.Context) (*HealthResponse, error) {
	var out HealthResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/health", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
