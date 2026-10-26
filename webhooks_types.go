package audact

import "encoding/json"

// WebhookEventType is a webhook event name such as call_ended.
type WebhookEventType string

// Webhook event types accepted when creating or updating an endpoint.
const (
	WebhookEventCallAnalyzed                          WebhookEventType = "call_analyzed"
	WebhookEventCallEndDetected                       WebhookEventType = "call_end_detected"
	WebhookEventCallEnded                             WebhookEventType = "call_ended"
	WebhookEventCallStarted                           WebhookEventType = "call_started"
	WebhookEventCallSummary                           WebhookEventType = "call_summary"
	WebhookEventChatAnalyzed                          WebhookEventType = "chat_analyzed"
	WebhookEventChatEnded                             WebhookEventType = "chat_ended"
	WebhookEventChatStarted                           WebhookEventType = "chat_started"
	WebhookEventComplianceBYOConnectionCreated        WebhookEventType = "compliance.byo_connection_created"
	WebhookEventComplianceBYOConnectionDeleted        WebhookEventType = "compliance.byo_connection_deleted"
	WebhookEventComplianceBYOConnectionHealthChanged  WebhookEventType = "compliance.byo_connection_health_changed"
	WebhookEventComplianceBYOConnectionVerified       WebhookEventType = "compliance.byo_connection_verified"
	WebhookEventComplianceBYONumberAttached           WebhookEventType = "compliance.byo_number_attached"
	WebhookEventComplianceBYONumberDetached           WebhookEventType = "compliance.byo_number_detached"
	WebhookEventComplianceCallBlocked                 WebhookEventType = "compliance.call_blocked"
	WebhookEventComplianceCallScoreFinalized          WebhookEventType = "compliance.call_score_finalized"
	WebhookEventComplianceChainTampered               WebhookEventType = "compliance.chain_tampered"
	WebhookEventComplianceChainVerified               WebhookEventType = "compliance.chain_verified"
	WebhookEventComplianceFailOpen                    WebhookEventType = "compliance.fail_open"
	WebhookEventComplianceKYCSetApproved              WebhookEventType = "compliance.kyc_set_approved"
	WebhookEventComplianceKYCSetRejected              WebhookEventType = "compliance.kyc_set_rejected"
	WebhookEventComplianceKYCSetSubmitted             WebhookEventType = "compliance.kyc_set_submitted"
	WebhookEventComplianceNumberProvisioningApproved  WebhookEventType = "compliance.number_provisioning_approved"
	WebhookEventComplianceNumberProvisioningFailed    WebhookEventType = "compliance.number_provisioning_failed"
	WebhookEventComplianceNumberProvisioningRequested WebhookEventType = "compliance.number_provisioning_requested"
	WebhookEventCompliancePolicyEvaluated             WebhookEventType = "compliance.policy_evaluated"
	WebhookEventComplianceReceiptIssued               WebhookEventType = "compliance.receipt_issued"
	WebhookEventComplianceValidationCompleted         WebhookEventType = "compliance.validation_completed"
	WebhookEventComplianceViolationDetected           WebhookEventType = "compliance.violation_detected"
	WebhookEventHandoverCompleted                     WebhookEventType = "handover_completed"
	WebhookEventHandoverInitiated                     WebhookEventType = "handover_initiated"
	WebhookEventHandoverSummaryReady                  WebhookEventType = "handover_summary_ready"
	WebhookEventTranscriptUpdated                     WebhookEventType = "transcript_updated"
	WebhookEventTransferBridged                       WebhookEventType = "transfer_bridged"
	WebhookEventTransferCancelled                     WebhookEventType = "transfer_cancelled"
	WebhookEventTransferEnded                         WebhookEventType = "transfer_ended"
	WebhookEventTransferStarted                       WebhookEventType = "transfer_started"
	WebhookEventVoicemailDetected                     WebhookEventType = "voicemail_detected"
)

// WebhookDeliveryStatus is the lifecycle state of a webhook delivery.
type WebhookDeliveryStatus string

// Webhook delivery statuses.
const (
	WebhookDeliveryPending    WebhookDeliveryStatus = "pending"
	WebhookDeliverySuccess    WebhookDeliveryStatus = "success"
	WebhookDeliveryFailed     WebhookDeliveryStatus = "failed"
	WebhookDeliveryDeadLetter WebhookDeliveryStatus = "dead_letter"
)

// WebhookEndpoint is a registered webhook endpoint.
type WebhookEndpoint struct {
	AgentID  string   `json:"agent_id,omitempty"`
	Events   []string `json:"events,omitempty"`
	ID       string   `json:"id,omitempty"`
	IsActive bool     `json:"is_active,omitempty"`
	Scope    string   `json:"scope,omitempty"`
	TenantID string   `json:"tenant_id,omitempty"`
	URL      string   `json:"url,omitempty"`
}

// ListWebhooksParams holds optional query parameters for Webhooks.List.
type ListWebhooksParams struct {
	Limit  int
	Cursor string
}

// ListWebhooksResponse is a page of webhook endpoints.
type ListWebhooksResponse struct {
	HasMore    bool              `json:"has_more,omitempty"`
	Items      []WebhookEndpoint `json:"items,omitempty"`
	NextCursor string            `json:"next_cursor,omitempty"`
}

// CreateWebhookRequest is the body of Webhooks.Create.
type CreateWebhookRequest struct {
	AgentID string             `json:"agent_id,omitempty"`
	Events  []WebhookEventType `json:"events,omitempty"`
	Scope   string             `json:"scope,omitempty"`
	URL     string             `json:"url,omitempty"`
}

// CreateWebhookResponse is the created endpoint plus its one-time signing secret.
type CreateWebhookResponse struct {
	AgentID       string   `json:"agent_id,omitempty"`
	Events        []string `json:"events,omitempty"`
	ID            string   `json:"id,omitempty"`
	IsActive      bool     `json:"is_active,omitempty"`
	Scope         string   `json:"scope,omitempty"`
	SigningSecret string   `json:"signing_secret,omitempty"`
	TenantID      string   `json:"tenant_id,omitempty"`
	URL           string   `json:"url,omitempty"`
}

// UpdateWebhookRequest is the body of Webhooks.Update; unset fields are not sent.
type UpdateWebhookRequest struct {
	Events   []WebhookEventType `json:"events,omitempty"`
	IsActive *bool              `json:"is_active,omitempty"`
	URL      string             `json:"url,omitempty"`
}

// ListWebhookDeliveriesParams holds optional query parameters for Webhooks.ListDeliveries.
type ListWebhookDeliveriesParams struct {
	Limit  int
	Offset int
}

// WebhookDeliveryRecord is one delivery attempt log entry.
type WebhookDeliveryRecord struct {
	AttemptCount int                   `json:"attempt_count,omitempty"`
	CompletedAt  string                `json:"completed_at,omitempty"`
	CreatedAt    string                `json:"created_at,omitempty"`
	EndpointID   string                `json:"endpoint_id,omitempty"`
	EventID      string                `json:"event_id,omitempty"`
	EventType    string                `json:"event_type,omitempty"`
	HTTPStatus   int                   `json:"http_status,omitempty"`
	ID           string                `json:"id,omitempty"`
	LastError    string                `json:"last_error,omitempty"`
	LatencyMs    float64               `json:"latency_ms,omitempty"`
	Status       WebhookDeliveryStatus `json:"status,omitempty"`
}

// ListWebhookDeliveriesResponse is a page of delivery records.
type ListWebhookDeliveriesResponse struct {
	Deliveries []WebhookDeliveryRecord `json:"deliveries,omitempty"`
	Total      int                     `json:"total,omitempty"`
}

// WebhookTestResponse is the result of a test ping.
type WebhookTestResponse struct {
	Error      string  `json:"error,omitempty"`
	LatencyMs  float64 `json:"latency_ms,omitempty"`
	StatusCode int     `json:"status_code,omitempty"`
	Success    bool    `json:"success,omitempty"`
}

// ResendWebhookDeliveryResponse reports a requeued delivery.
type ResendWebhookDeliveryResponse struct {
	DeliveryID string `json:"delivery_id,omitempty"`
	Requeued   bool   `json:"requeued,omitempty"`
}

// WebhookEventTypeInfo describes one subscribable event type.
type WebhookEventTypeInfo struct {
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	SchemaURL   string `json:"schema_url,omitempty"`
}

// WebhookEventTypesResponse lists the subscribable event types.
type WebhookEventTypesResponse struct {
	EventTypes []WebhookEventTypeInfo `json:"event_types,omitempty"`
}

// WebhookEventSchema is the schema and example payload of an event type.
type WebhookEventSchema struct {
	Category        string          `json:"category,omitempty"`
	Description     string          `json:"description,omitempty"`
	EventType       string          `json:"event_type,omitempty"`
	ExamplePayload  json.RawMessage `json:"example_payload,omitempty"`
	PIIHandlingNote string          `json:"pii_handling_note,omitempty"`
	SchemaVersion   string          `json:"schema_version,omitempty"`
}

// RotateWebhookSecretResponse carries the new signing secret.
type RotateWebhookSecretResponse struct {
	GracePeriodEnds string `json:"grace_period_ends,omitempty"`
	Rotated         bool   `json:"rotated,omitempty"`
	SigningSecret   string `json:"signing_secret,omitempty"`
}

// WebhookSecretStatusResponse is the signing-secret rotation status.
type WebhookSecretStatusResponse struct {
	CreatedAt       string `json:"created_at,omitempty"`
	GracePeriodEnds string `json:"grace_period_ends,omitempty"`
	InRotation      bool   `json:"in_rotation,omitempty"`
}
