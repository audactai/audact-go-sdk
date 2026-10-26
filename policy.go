package audact

import (
	"context"
	"net/http"
	"net/url"
)

// Evaluate evaluates call facts against a country policy. Calls POST /v1/evaluate.
func (s *PolicyService) Evaluate(ctx context.Context, req *EvaluateRequest) (*EvaluateResponse, error) {
	var out EvaluateResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/evaluate", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns the policy document of a country. Calls GET /v1/policies/{country}.
func (s *PolicyService) Get(ctx context.Context, country string) (*PolicyFileResponse, error) {
	var out PolicyFileResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/policies/"+url.PathEscape(country), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CheckCallingHours reports whether calling is currently permitted.
// Calls GET /v1/calling-hours/check.
func (s *PolicyService) CheckCallingHours(ctx context.Context, params CheckCallingHoursParams) (*CallingHoursCheckResponse, error) {
	q := url.Values{}
	q.Set("country", params.Country)
	addString(q, "call_time", params.CallTime)
	var out CallingHoursCheckResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/calling-hours/check", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
