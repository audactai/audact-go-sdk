package audact

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// List lists calls. Calls GET /v1/calls.
func (s *CallsService) List(ctx context.Context, params *ListCallsParams) (*ListCallsResponse, error) {
	q := url.Values{}
	if params != nil {
		addString(q, "date_from", params.DateFrom)
		addString(q, "date_to", params.DateTo)
		addString(q, "status", params.Status)
		addString(q, "verdict", params.Verdict)
		if params.Shadow != nil {
			q.Set("shadow", strconv.FormatBool(*params.Shadow))
		}
		addString(q, "country", params.Country)
		addString(q, "config_id", params.ConfigID)
		addString(q, "cursor", params.Cursor)
		addInt(q, "limit", params.Limit)
	}
	var out ListCallsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/calls", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns the detail of a call. Calls GET /v1/calls/{call_id}.
func (s *CallsService) Get(ctx context.Context, callID string) (*CallDetail, error) {
	var out CallDetail
	if _, err := s.t.do(ctx, http.MethodGet, "/calls/"+url.PathEscape(callID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Summary returns the composite summary of a call. Calls GET /v1/calls/{call_id}/summary.
func (s *CallsService) Summary(ctx context.Context, callID string) (*CallSummary, error) {
	var out CallSummary
	if _, err := s.t.do(ctx, http.MethodGet, "/calls/"+url.PathEscape(callID)+"/summary", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRecording returns a short-lived signed URL for the call recording.
// Calls GET /v1/calls/{call_id}/recording.
func (s *CallsService) GetRecording(ctx context.Context, callID string) (*RecordingResponse, error) {
	var out RecordingResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/calls/"+url.PathEscape(callID)+"/recording", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteRecording deletes the call recording. Calls DELETE /v1/calls/{call_id}/recording.
func (s *CallsService) DeleteRecording(ctx context.Context, callID string) error {
	_, err := s.t.do(ctx, http.MethodDelete, "/calls/"+url.PathEscape(callID)+"/recording", authAPIKey, nil, nil, nil)
	return err
}

func idemHeader(key string) map[string]string {
	if key == "" {
		return nil
	}
	return map[string]string{"Idempotency-Key": key}
}

// CreateOutbound places an outbound call; a non-empty idempotencyKey is sent
// as the Idempotency-Key header. Calls POST /v1/sip/calls/outbound.
// req.FromNumber and req.DestinationNumber must be strict E.164 with no
// spaces or separators (they are not normalised), and req.FromNumber must be
// one of your own numbers; otherwise the server answers 400 INVALID_REQUEST.
func (s *CallsService) CreateOutbound(ctx context.Context, req *OutboundCallRequest, idempotencyKey string) (*OutboundCallResponse, error) {
	var out OutboundCallResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/sip/calls/outbound", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateInbound registers an inbound call; a non-empty idempotencyKey is sent
// as the Idempotency-Key header. Calls POST /v1/sip/calls/inbound.
func (s *CallsService) CreateInbound(ctx context.Context, req *InboundCallRequest, idempotencyKey string) (*InboundCallResponse, error) {
	var out InboundCallResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/sip/calls/inbound", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateBrowser creates a browser (WebRTC) call; a non-empty idempotencyKey is
// sent as the Idempotency-Key header. Calls POST /v1/sip/calls/browser.
func (s *CallsService) CreateBrowser(ctx context.Context, req *BrowserCallRequest, idempotencyKey string) (*BrowserCallResponse, error) {
	var out BrowserCallResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/sip/calls/browser", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Status returns the lifecycle status of a call. Calls GET /v1/sip/calls/{call_id}/status.
func (s *CallsService) Status(ctx context.Context, callID string) (*CallStatusResponse, error) {
	var out CallStatusResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/sip/calls/"+url.PathEscape(callID)+"/status", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Hangup terminates an active call. Calls POST /v1/sip/calls/{call_id}/hangup.
func (s *CallsService) Hangup(ctx context.Context, callID string) (*HangupResponse, error) {
	var out HangupResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/sip/calls/"+url.PathEscape(callID)+"/hangup", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
