package audact

import (
	"context"
	"net/http"
	"net/url"
)

// List lists phone numbers in the pool. Calls GET /v1/numbers.
func (s *NumbersService) List(ctx context.Context, params *ListNumbersParams) (*ListNumbersResponse, error) {
	q := url.Values{}
	if params != nil {
		addString(q, "country", params.Country)
		addString(q, "status", string(params.Status))
		addInt(q, "limit", params.Limit)
		addString(q, "cursor", params.Cursor)
	}
	var out ListNumbersResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/numbers", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create adds a phone number to the pool; a non-empty idempotencyKey is sent
// as the Idempotency-Key header. Calls POST /v1/numbers.
func (s *NumbersService) Create(ctx context.Context, req *CreateNumberRequest, idempotencyKey string) (*PhoneNumber, error) {
	var out PhoneNumber
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/numbers", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAvailable searches the provider's live inventory. Calls GET /v1/numbers/available.
func (s *NumbersService) ListAvailable(ctx context.Context, params *ListAvailableNumbersParams) (*ListAvailableNumbersResponse, error) {
	q := url.Values{}
	if params != nil {
		addString(q, "country", params.Country)
		addString(q, "number_type", params.NumberType)
	}
	var out ListAvailableNumbersResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/numbers/available", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Lookup finds a phone number by DID. Calls GET /v1/numbers/lookup.
func (s *NumbersService) Lookup(ctx context.Context, did string) (*PhoneNumber, error) {
	q := url.Values{}
	q.Set("did", did)
	var out PhoneNumber
	if _, err := s.t.do(ctx, http.MethodGet, "/numbers/lookup", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListOrders lists number orders. Calls GET /v1/numbers/orders.
func (s *NumbersService) ListOrders(ctx context.Context) (*ListNumberOrdersResponse, error) {
	var out ListNumberOrdersResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/numbers/orders", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateOrder orders a number; a non-empty idempotencyKey is sent as the
// Idempotency-Key header. Calls POST /v1/numbers/orders.
func (s *NumbersService) CreateOrder(ctx context.Context, req *CreateNumberOrderRequest, idempotencyKey string) (*NumberOrder, error) {
	var out NumberOrder
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/numbers/orders", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetOrder returns a number order. Calls GET /v1/numbers/orders/{id}.
func (s *NumbersService) GetOrder(ctx context.Context, id string) (*NumberOrder, error) {
	var out NumberOrder
	if _, err := s.t.do(ctx, http.MethodGet, "/numbers/orders/"+url.PathEscape(id), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes an available phone number. Calls DELETE /v1/numbers/{id}.
func (s *NumbersService) Delete(ctx context.Context, id string) error {
	_, err := s.t.do(ctx, http.MethodDelete, "/numbers/"+url.PathEscape(id), authAPIKey, nil, nil, nil)
	return err
}

// Assign assigns a number to an agent; a non-empty idempotencyKey is sent as
// the Idempotency-Key header. Calls POST /v1/numbers/{id}/assign.
func (s *NumbersService) Assign(ctx context.Context, id string, req *AssignNumberRequest, idempotencyKey string) error {
	_, err := s.t.doHeaders(ctx, http.MethodPost, "/numbers/"+url.PathEscape(id)+"/assign", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), nil)
	return err
}

// Release releases a number; a non-empty idempotencyKey is sent as the
// Idempotency-Key header. Calls DELETE /v1/numbers/{id}/release.
func (s *NumbersService) Release(ctx context.Context, id, idempotencyKey string) error {
	_, err := s.t.doHeaders(ctx, http.MethodDelete, "/numbers/"+url.PathEscape(id)+"/release", authAPIKey, nil, idemHeader(idempotencyKey), nil, nil)
	return err
}

// Resume resumes a suspended number; a non-empty idempotencyKey is sent as the
// Idempotency-Key header. Calls POST /v1/numbers/{id}/resume.
func (s *NumbersService) Resume(ctx context.Context, id, idempotencyKey string) error {
	_, err := s.t.doHeaders(ctx, http.MethodPost, "/numbers/"+url.PathEscape(id)+"/resume", authAPIKey, nil, idemHeader(idempotencyKey), nil, nil)
	return err
}

// Suspend suspends a number; a non-empty idempotencyKey is sent as the
// Idempotency-Key header. Calls POST /v1/numbers/{id}/suspend.
func (s *NumbersService) Suspend(ctx context.Context, id, idempotencyKey string) error {
	_, err := s.t.doHeaders(ctx, http.MethodPost, "/numbers/"+url.PathEscape(id)+"/suspend", authAPIKey, nil, idemHeader(idempotencyKey), nil, nil)
	return err
}

// SetTrunk sets the trunk of a number. Calls PUT /v1/numbers/{id}/trunk.
func (s *NumbersService) SetTrunk(ctx context.Context, id string, req *SetNumberTrunkRequest) (*PhoneNumber, error) {
	var out PhoneNumber
	if _, err := s.t.do(ctx, http.MethodPut, "/numbers/"+url.PathEscape(id)+"/trunk", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}
