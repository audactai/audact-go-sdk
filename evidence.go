package audact

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strconv"
)

func evidencePath(prefix, id string) string { return prefix + url.PathEscape(id) }

// PostEvent writes an evidence event and returns its hash receipt. An empty
// idempotencyKey sends no Idempotency-Key header. Calls POST /v1/events.
func (s *EvidenceService) PostEvent(ctx context.Context, req *PostEventRequest, idempotencyKey string) (*PostEventResponse, error) {
	var out PostEventResponse
	h := map[string]string{"Idempotency-Key": idempotencyKey}
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/events", authAPIKey, nil, h, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Chain returns a page of the evidence chain for a call. Calls GET /v1/chains/{call_id}.
func (s *EvidenceService) Chain(ctx context.Context, callID string, params *ChainParams) (*GetChainResponse, error) {
	q := url.Values{}
	if params != nil {
		if params.AfterSequence != 0 {
			q.Set("after_sequence", strconv.Itoa(params.AfterSequence))
		}
		addInt(q, "limit", params.Limit)
	}
	var out GetChainResponse
	if _, err := s.t.do(ctx, http.MethodGet, evidencePath("/chains/", callID), authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// VerifyChain verifies the evidence hash chain of a call. Calls GET /v1/verify/{call_id}.
func (s *EvidenceService) VerifyChain(ctx context.Context, callID string) (*VerifyChainResponse, error) {
	var out VerifyChainResponse
	if _, err := s.t.do(ctx, http.MethodGet, evidencePath("/verify/", callID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Reconstruct reconstructs the event chain of a call. Calls GET /v1/reconstruct/{call_id}.
func (s *EvidenceService) Reconstruct(ctx context.Context, callID string) (*ReconstructResponse, error) {
	var out ReconstructResponse
	if _, err := s.t.do(ctx, http.MethodGet, evidencePath("/reconstruct/", callID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListViolations lists non-ALLOW check results. Calls GET /v1/violations.
func (s *EvidenceService) ListViolations(ctx context.Context, params *ListViolationsParams) (*GetViolationsResponse, error) {
	q := url.Values{}
	if params != nil {
		addString(q, "call_id", params.CallID)
		addString(q, "date_from", params.DateFrom)
		addString(q, "date_to", params.DateTo)
		addString(q, "verdict", params.Verdict)
		addString(q, "cursor", params.Cursor)
		addInt(q, "limit", params.Limit)
	}
	var out GetViolationsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/violations", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Receipt returns the compliance receipt of a call. Calls GET /v1/receipts/{call_id}.
func (s *EvidenceService) Receipt(ctx context.Context, callID string) (*GetReceiptResponse, error) {
	var out GetReceiptResponse
	if _, err := s.t.do(ctx, http.MethodGet, evidencePath("/receipts/", callID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeliveryStatus returns per-turn delivery status of a call. Calls GET /v1/calls/{call_id}/delivery.
func (s *EvidenceService) DeliveryStatus(ctx context.Context, callID string) ([]TurnDeliveryStatus, error) {
	var out []TurnDeliveryStatus
	if _, err := s.t.do(ctx, http.MethodGet, "/calls/"+url.PathEscape(callID)+"/delivery", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ConfirmDelivery records a delivery confirmation. An empty idempotencyKey
// sends no header. Calls POST /v1/delivery-confirmation.
func (s *EvidenceService) ConfirmDelivery(ctx context.Context, req *ConfirmDeliveryRequest, idempotencyKey string) (*ConfirmDeliveryResponse, error) {
	var out ConfirmDeliveryResponse
	h := map[string]string{"Idempotency-Key": idempotencyKey}
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/delivery-confirmation", authAPIKey, nil, h, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ReconcileTranscript reconciles a transcript (an object with call_id,
// segments and source) against validated segments. An empty idempotencyKey
// sends no header. Calls POST /v1/transcript-reconciliation.
func (s *EvidenceService) ReconcileTranscript(ctx context.Context, req map[string]any, idempotencyKey string) (*ReconcileTranscriptResponse, error) {
	if req == nil {
		req = map[string]any{}
	}
	var out ReconcileTranscriptResponse
	h := map[string]string{"Idempotency-Key": idempotencyKey}
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/transcript-reconciliation", authAPIKey, nil, h, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Export downloads an Article 50 evidence export using a capability token
// from an evidence link; no API key is sent. The format defaults to the
// token's (platform-minted links are "pdf"), so the result carries whichever
// the server served. Calls GET /v1/evidence-exports/{call_id}.
func (s *EvidenceService) Export(ctx context.Context, callID, token string, params *Article50Params) (*Article50Export, error) {
	q := url.Values{}
	q.Set("token", token)
	if params != nil {
		addString(q, "format", params.Format)
	}
	return s.article50(ctx, evidencePath("/evidence-exports/", callID), authNone, q)
}

// Article50Report assembles and signs the Article 50 report of a call: JSON
// by default, or a PDF with Format "pdf".
// Calls GET /v1/calls/{call_id}/article50-report.
func (s *EvidenceService) Article50Report(ctx context.Context, callID string, params *Article50Params) (*Article50Export, error) {
	q := url.Values{}
	if params != nil {
		addString(q, "format", params.Format)
	}
	return s.article50(ctx, "/calls/"+url.PathEscape(callID)+"/article50-report", authAPIKey, q)
}

// article50 fetches an Article 50 export and decodes it by Content-Type.
func (s *EvidenceService) article50(ctx context.Context, path string, auth authKind, q url.Values) (*Article50Export, error) {
	var raw rawResponse
	h := map[string]string{"Accept": "application/json, application/pdf"}
	if _, err := s.t.doHeaders(ctx, http.MethodGet, path, auth, q, h, nil, &raw); err != nil {
		return nil, err
	}
	if mt, _, _ := mime.ParseMediaType(raw.contentType); mt == "application/pdf" {
		return &Article50Export{PDF: raw.data}, nil
	}
	var rep Article50Report
	if err := json.Unmarshal(raw.data, &rep); err != nil {
		return nil, fmt.Errorf("audact: GET %s: decode response: %w", path, err)
	}
	return &Article50Export{Report: &rep}, nil
}

// GenerateC2PAManifest generates a signed C2PA manifest for a call. Every
// successful call issues and permanently stores a new manifest (evidence is
// append-only), so pass an idempotencyKey to make retries safe; an empty
// idempotencyKey sends no Idempotency-Key header.
// Calls POST /v1/calls/{call_id}/c2pa.
func (s *EvidenceService) GenerateC2PAManifest(ctx context.Context, callID, idempotencyKey string) (*C2PAManifest, error) {
	var out C2PAManifest
	h := map[string]string{"Idempotency-Key": idempotencyKey}
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/calls/"+url.PathEscape(callID)+"/c2pa", authAPIKey, nil, h, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListC2PAManifests lists the C2PA manifests issued for a call, newest first.
// Calls GET /v1/calls/{call_id}/c2pa.
func (s *EvidenceService) ListC2PAManifests(ctx context.Context, callID string, params *ListC2PAManifestsParams) (*C2PAManifestListResponse, error) {
	q := url.Values{}
	if params != nil {
		addInt(q, "limit", params.Limit)
		addString(q, "cursor", params.Cursor)
	}
	var out C2PAManifestListResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/calls/"+url.PathEscape(callID)+"/c2pa", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadC2PAManifestStore downloads the stored C2PA manifest store
// (application/c2pa). Verify it with
// `c2patool <asset> --external-manifest <file>.c2pa`.
// Calls GET /v1/calls/{call_id}/c2pa/{manifest_id}/manifest.c2pa.
func (s *EvidenceService) DownloadC2PAManifestStore(ctx context.Context, callID, manifestID string) ([]byte, error) {
	var raw rawResponse
	h := map[string]string{"Accept": "application/c2pa"}
	if _, err := s.t.doHeaders(ctx, http.MethodGet, c2paManifestPath(callID, manifestID, "/manifest.c2pa"), authAPIKey, nil, h, nil, &raw); err != nil {
		return nil, err
	}
	return raw.data, nil
}

// DownloadC2PAAsset downloads the asset a manifest is hard-bound to. For a
// recording the server answers 302 to a presigned URL; that URL is fetched
// without the Authorization header, because it authenticates itself and the
// API key must not reach the storage provider.
// Calls GET /v1/calls/{call_id}/c2pa/{manifest_id}/asset.
func (s *EvidenceService) DownloadC2PAAsset(ctx context.Context, callID, manifestID string) ([]byte, error) {
	path := c2paManifestPath(callID, manifestID, "/asset")
	hc := *s.t.hc
	hc.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }

	tok, _ := s.t.credential(authAPIKey)
	resp, err := s.fetchAsset(ctx, &hc, path, s.t.baseURL+"/v1"+path, tok)
	if err != nil {
		return nil, err
	}
	if resp.status >= 300 && resp.status < 400 && resp.location != "" {
		loc, err := resp.url.Parse(resp.location)
		if err != nil {
			return nil, fmt.Errorf("audact: GET %s: redirect location: %w", path, redactURLError(err))
		}
		if resp, err = s.fetchAsset(ctx, &hc, path, loc.String(), ""); err != nil {
			return nil, err
		}
	}
	return resp.body, resp.err
}

type assetResponse struct {
	status   int
	location string
	url      *url.URL
	body     []byte
	err      error
}

// fetchAsset performs one GET without following redirects. tok, when set, is
// sent as the Bearer credential. Errors name only path: rawURL may be a
// presigned URL whose query is a credential.
func (s *EvidenceService) fetchAsset(ctx context.Context, hc *http.Client, path, rawURL, tok string) (*assetResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("audact: GET %s: %w", path, redactURLError(err))
	}
	req.Header.Set("User-Agent", s.t.userAgent)
	if tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("audact: GET %s: %w", path, redactURLError(err))
	}
	defer func() { _ = resp.Body.Close() }()
	out := &assetResponse{status: resp.StatusCode, location: resp.Header.Get("Location"), url: req.URL}
	switch {
	case resp.StatusCode >= 300 && resp.StatusCode < 400 && out.location != "":
		return out, nil
	case resp.StatusCode < 200 || resp.StatusCode > 299:
		b, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
		out.err = parseError(resp, b)
		return out, nil
	}
	if out.body, err = io.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("audact: GET %s: read response: %w", path, err)
	}
	return out, nil
}

// C2PATrustAnchor returns the PEM trust anchor (top certificate of the
// signing chain) for `c2patool --trust_anchors`. No API key is sent.
// Calls GET /v1/c2pa/trust-anchor.
func (s *EvidenceService) C2PATrustAnchor(ctx context.Context) (string, error) {
	var raw rawResponse
	h := map[string]string{"Accept": "application/x-pem-file"}
	if _, err := s.t.doHeaders(ctx, http.MethodGet, "/c2pa/trust-anchor", authNone, nil, h, nil, &raw); err != nil {
		return "", err
	}
	return string(raw.data), nil
}

func c2paManifestPath(callID, manifestID, suffix string) string {
	return "/calls/" + url.PathEscape(callID) + "/c2pa/" + url.PathEscape(manifestID) + suffix
}

// Article50PublicKey returns the public key (as served, typically PEM) that
// verifies Article 50 report signatures. Calls GET /v1/keys/{kid}.
func (s *EvidenceService) Article50PublicKey(ctx context.Context, kid string) (string, error) {
	var raw rawResponse
	if _, err := s.t.do(ctx, http.MethodGet, evidencePath("/keys/", kid), authNone, nil, nil, &raw); err != nil {
		return "", err
	}
	return string(raw.data), nil
}

// ListComplianceEvents lists compliance events with keyset pagination.
// Calls GET /v1/compliance/events.
func (s *EvidenceService) ListComplianceEvents(ctx context.Context, params *ListComplianceEventsParams) (*ListComplianceEventsResponse, error) {
	q := url.Values{}
	if params != nil {
		addString(q, "event_subtype", string(params.EventSubtype))
		addString(q, "verdict", string(params.Verdict))
		addString(q, "severity", string(params.Severity))
		addString(q, "call_id", params.CallID)
		addString(q, "cursor", params.Cursor)
		addInt(q, "limit", params.Limit)
	}
	var out ListComplianceEventsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/compliance/events", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
