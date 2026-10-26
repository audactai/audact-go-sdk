package audact

import (
	"context"
	"net/http"
	"net/url"
)

func pageQuery(p *ListPageParams) url.Values {
	q := url.Values{}
	if p != nil {
		addString(q, "cursor", p.Cursor)
		addInt(q, "limit", p.Limit)
	}
	return q
}

// Handover starts a warm handover for an active call.
// Calls POST /v1/sip/calls/{call_id}/handover.
func (s *CallsService) Handover(ctx context.Context, callID string, req *HandoverRequest) (*HandoverResponse, error) {
	var out HandoverResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/sip/calls/"+url.PathEscape(callID)+"/handover", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHandover returns the state of a handover session.
// Calls GET /v1/sip/handovers/{handover_id}.
func (s *CallsService) GetHandover(ctx context.Context, handoverID string) (*HandoverStateResponse, error) {
	var out HandoverStateResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/sip/handovers/"+url.PathEscape(handoverID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Transfer starts a SIP transfer of an active call; a non-empty idempotencyKey
// is sent as the Idempotency-Key header. Calls POST /v1/sip/calls/{call_id}/transfer.
func (s *CallsService) Transfer(ctx context.Context, callID string, req *TransferRequest, idempotencyKey string) (*TransferResponse, error) {
	var out TransferResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/sip/calls/"+url.PathEscape(callID)+"/transfer", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTransfer returns the state of a transfer session.
// Calls GET /v1/sip/transfers/{transfer_id}.
func (s *CallsService) GetTransfer(ctx context.Context, transferID string) (*TransferSessionResponse, error) {
	var out TransferSessionResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/sip/transfers/"+url.PathEscape(transferID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *CallsService) transferAction(ctx context.Context, transferID, action string) (*TransferStateResponse, error) {
	var out TransferStateResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/sip/transfers/"+url.PathEscape(transferID)+"/"+action, authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// BridgeTransfer bridges the caller to the transfer target.
// Calls POST /v1/sip/transfers/{transfer_id}/bridge.
func (s *CallsService) BridgeTransfer(ctx context.Context, transferID string) (*TransferStateResponse, error) {
	return s.transferAction(ctx, transferID, "bridge")
}

// CancelTransfer cancels a transfer. Calls POST /v1/sip/transfers/{transfer_id}/cancel.
func (s *CallsService) CancelTransfer(ctx context.Context, transferID string) (*TransferStateResponse, error) {
	return s.transferAction(ctx, transferID, "cancel")
}

// EndTransfer ends a transfer. Calls POST /v1/sip/transfers/{transfer_id}/end.
func (s *CallsService) EndTransfer(ctx context.Context, transferID string) (*TransferStateResponse, error) {
	return s.transferAction(ctx, transferID, "end")
}

// ListRoutes lists inbound dispatch rules. Calls GET /v1/sip/routes.
func (s *CallsService) ListRoutes(ctx context.Context, params *ListPageParams) (*ListRoutesResponse, error) {
	var out ListRoutesResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/sip/routes", authAPIKey, pageQuery(params), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateRoute routes an inbound DID to an agent; a non-empty idempotencyKey is
// sent as the Idempotency-Key header. Calls POST /v1/sip/routes.
func (s *CallsService) CreateRoute(ctx context.Context, req *RouteInboundRequest, idempotencyKey string) (*RouteInboundResponse, error) {
	var out RouteInboundResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/sip/routes", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListTrunks lists SIP trunks. Calls GET /v1/sip/trunks.
func (s *CallsService) ListTrunks(ctx context.Context, params *ListPageParams) (*ListTrunksResponse, error) {
	var out ListTrunksResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/sip/trunks", authAPIKey, pageQuery(params), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateInboundTrunk registers an inbound SIP trunk; a non-empty idempotencyKey
// is sent as the Idempotency-Key header. Calls POST /v1/sip/trunks/inbound.
func (s *CallsService) CreateInboundTrunk(ctx context.Context, req *RegisterInboundTrunkRequest, idempotencyKey string) (*RegisterInboundTrunkResponse, error) {
	var out RegisterInboundTrunkResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/sip/trunks/inbound", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateOutboundTrunk registers an outbound SIP trunk; a non-empty idempotencyKey
// is sent as the Idempotency-Key header. Calls POST /v1/sip/trunks/outbound.
func (s *CallsService) CreateOutboundTrunk(ctx context.Context, req *RegisterOutboundTrunkRequest, idempotencyKey string) (*RegisterOutboundTrunkResponse, error) {
	var out RegisterOutboundTrunkResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/sip/trunks/outbound", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAdmissionURL returns the pre-connect admission hook URL.
// Calls GET /v1/sip/admission-url.
func (s *CallsService) GetAdmissionURL(ctx context.Context) (*AdmissionURLResponse, error) {
	var out AdmissionURLResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/sip/admission-url", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetAdmissionURL sets the pre-connect admission hook URL.
// Calls PUT /v1/sip/admission-url.
func (s *CallsService) SetAdmissionURL(ctx context.Context, req *SetAdmissionURLRequest) (*AdmissionURLResponse, error) {
	var out AdmissionURLResponse
	if _, err := s.t.do(ctx, http.MethodPut, "/sip/admission-url", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateWSTicket creates a short-lived WebSocket authentication ticket.
// Calls POST /v1/ws/ticket.
func (s *CallsService) CreateWSTicket(ctx context.Context) (*WSTicketResponse, error) {
	var out WSTicketResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/ws/ticket", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
