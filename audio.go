package audact

import (
	"context"
	"net/http"
	"net/url"
)

// Validate validates submitted audio against policy. Calls POST /v1/validate-audio.
func (s *AudioService) Validate(ctx context.Context, req ValidateAudioRequest) (*ValidateAudioResponse, error) {
	var out ValidateAudioResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/validate-audio", authAPIKey, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Reconcile reconciles audio against validated transcript segments. Calls
// POST /v1/audio-reconciliation.
func (s *AudioService) Reconcile(ctx context.Context, req AudioReconciliationRequest, params *ReconcileAudioParams) (*AudioReconciliationResponse, error) {
	p := orEmpty(params)
	q := url.Values{}
	addInt(q, "limit", p.Limit)
	addInt(q, "after_sequence", p.AfterSequence)
	extra := map[string]string{"Idempotency-Key": p.IdempotencyKey}
	var out AudioReconciliationResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/audio-reconciliation", authAPIKey, q, extra, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
