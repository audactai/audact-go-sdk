package audact

import (
	"context"
	"net/http"
	"net/url"
)

// Trigger triggers an aggregation run. Calls POST /v1/admin/aggregate.
func (s *AggregationService) Trigger(ctx context.Context, req *AdminAggregateRequest) (*AdminAggregateResponse, error) {
	var out AdminAggregateResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/admin/aggregate", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListProofs lists the tenant's aggregation proofs. Calls GET /v1/proofs.
func (s *AggregationService) ListProofs(ctx context.Context, params *ListProofsParams) (*ProofListResponse, error) {
	p := orEmpty(params)
	q := url.Values{}
	addInt(q, "limit", p.Limit)
	addString(q, "cursor", p.Cursor)
	var out ProofListResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/proofs", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Proof returns the proof of one aggregation window. Calls GET /v1/proofs/{window_id}.
func (s *AggregationService) Proof(ctx context.Context, windowID string) (*ProofResponse, error) {
	var out ProofResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/proofs/"+url.PathEscape(windowID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Metachain returns a page of the public meta-chain. Calls GET /v1/metachain.
func (s *AggregationService) Metachain(ctx context.Context, params *MetachainParams) (*MetaChainResponse, error) {
	p := orEmpty(params)
	q := url.Values{}
	addInt(q, "limit", p.Limit)
	addString(q, "cursor", p.Cursor)
	addString(q, "order", p.Order)
	var out MetaChainResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/metachain", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// VerifyInclusion verifies a Merkle inclusion proof. Calls POST /v1/verify-inclusion.
func (s *AggregationService) VerifyInclusion(ctx context.Context, req VerifyInclusionRequest) (*VerifyInclusionResponse, error) {
	var out VerifyInclusionResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/verify-inclusion", authAPIKey, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
