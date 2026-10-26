package audact

import (
	"context"
	"net/http"
	"net/url"
)

// List lists webhook endpoints. Calls GET /v1/webhooks.
func (s *WebhooksService) List(ctx context.Context, params *ListWebhooksParams) (*ListWebhooksResponse, error) {
	q := url.Values{}
	if params != nil {
		addInt(q, "limit", params.Limit)
		addString(q, "cursor", params.Cursor)
	}
	var out ListWebhooksResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/webhooks", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create registers a webhook endpoint; a non-empty idempotencyKey is sent as
// the Idempotency-Key header. Calls POST /v1/webhooks.
func (s *WebhooksService) Create(ctx context.Context, req *CreateWebhookRequest, idempotencyKey string) (*CreateWebhookResponse, error) {
	var out CreateWebhookResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/webhooks", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns a webhook endpoint. Calls GET /v1/webhooks/{id}.
func (s *WebhooksService) Get(ctx context.Context, id string) (*WebhookEndpoint, error) {
	var out WebhookEndpoint
	if _, err := s.t.do(ctx, http.MethodGet, "/webhooks/"+url.PathEscape(id), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates a webhook endpoint. Calls PUT /v1/webhooks/{id}.
func (s *WebhooksService) Update(ctx context.Context, id string, req *UpdateWebhookRequest) (*WebhookEndpoint, error) {
	var out WebhookEndpoint
	if _, err := s.t.do(ctx, http.MethodPut, "/webhooks/"+url.PathEscape(id), authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a webhook endpoint. Calls DELETE /v1/webhooks/{id}.
func (s *WebhooksService) Delete(ctx context.Context, id string) error {
	_, err := s.t.do(ctx, http.MethodDelete, "/webhooks/"+url.PathEscape(id), authAPIKey, nil, nil, nil)
	return err
}

// ListDeliveries lists delivery attempts of an endpoint. Calls GET /v1/webhooks/{id}/deliveries.
func (s *WebhooksService) ListDeliveries(ctx context.Context, id string, params *ListWebhookDeliveriesParams) (*ListWebhookDeliveriesResponse, error) {
	q := url.Values{}
	if params != nil {
		addInt(q, "limit", params.Limit)
		addInt(q, "offset", params.Offset)
	}
	var out ListWebhookDeliveriesResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/webhooks/"+url.PathEscape(id)+"/deliveries", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Test sends a test ping to an endpoint. Calls POST /v1/webhooks/{id}/test.
func (s *WebhooksService) Test(ctx context.Context, id string) (*WebhookTestResponse, error) {
	var out WebhookTestResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/webhooks/"+url.PathEscape(id)+"/test", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ResendDelivery requeues a delivery. Calls POST /v1/webhooks/deliveries/{id}/resend.
func (s *WebhooksService) ResendDelivery(ctx context.Context, deliveryID string) (*ResendWebhookDeliveryResponse, error) {
	var out ResendWebhookDeliveryResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/webhooks/deliveries/"+url.PathEscape(deliveryID)+"/resend", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListEventTypes lists subscribable event types. Calls GET /v1/webhooks/event-types.
func (s *WebhooksService) ListEventTypes(ctx context.Context) (*WebhookEventTypesResponse, error) {
	var out WebhookEventTypesResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/webhooks/event-types", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EventSchemas returns the schema of one event type. Calls GET /v1/webhooks/event-schemas.
func (s *WebhooksService) EventSchemas(ctx context.Context, eventType string) (*WebhookEventSchema, error) {
	q := url.Values{}
	addString(q, "type", eventType)
	var out WebhookEventSchema
	if _, err := s.t.do(ctx, http.MethodGet, "/webhooks/event-schemas", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RotateSecret rotates the webhook signing secret. Calls POST /v1/webhooks/secrets/rotate.
func (s *WebhooksService) RotateSecret(ctx context.Context) (*RotateWebhookSecretResponse, error) {
	var out RotateWebhookSecretResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/webhooks/secrets/rotate", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SecretStatus returns the signing-secret rotation status. Calls GET /v1/webhooks/secrets/status.
func (s *WebhooksService) SecretStatus(ctx context.Context) (*WebhookSecretStatusResponse, error) {
	var out WebhookSecretStatusResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/webhooks/secrets/status", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
