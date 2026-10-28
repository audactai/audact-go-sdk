package audact

import "encoding/json"

// Tenant is a tenant record.
type Tenant struct {
	CreatedAt  string `json:"created_at,omitempty"`
	HMACKeyARN string `json:"hmac_key_arn,omitempty"`
	Name       string `json:"name,omitempty"`
	Plan       string `json:"plan,omitempty"`
	TenantID   string `json:"tenant_id,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

// CreateTenantRequest is the body of Tenants.Create.
type CreateTenantRequest struct {
	KeyName string `json:"key_name,omitempty"`
	Name    string `json:"name,omitempty"`
	Plan    string `json:"plan,omitempty"`
}

// CreateTenantAPIKeyInfo is the first API key provisioned with a new tenant;
// Key is returned only once.
type CreateTenantAPIKeyInfo struct {
	Key       string `json:"key,omitempty"`
	KeyID     string `json:"key_id,omitempty"`
	KeyPrefix string `json:"key_prefix,omitempty"`
	Name      string `json:"name,omitempty"`
	Role      string `json:"role,omitempty"`
}

// CreateTenantResponse is the response of Tenants.Create.
type CreateTenantResponse struct {
	APIKey     *CreateTenantAPIKeyInfo `json:"api_key,omitempty"`
	CreatedAt  string                  `json:"created_at,omitempty"`
	HMACKeyARN string                  `json:"hmac_key_arn,omitempty"`
	Name       string                  `json:"name,omitempty"`
	Plan       string                  `json:"plan,omitempty"`
	TenantID   string                  `json:"tenant_id,omitempty"`
	UpdatedAt  string                  `json:"updated_at,omitempty"`
}

// UpdateTenantRequest is the body of Tenants.Update.
type UpdateTenantRequest struct {
	Name string `json:"name,omitempty"`
	Plan string `json:"plan,omitempty"`
}

// GetTenantUsageParams are the optional RFC 3339 filters of Tenants.Usage.
type GetTenantUsageParams struct {
	PeriodStart string
	PeriodEnd   string
}

// UsageSource says where a UsageRecord's counts come from.
type UsageSource string

const (
	// UsageSourceStored means the counts come from a stored usage record.
	UsageSourceStored UsageSource = "stored"
	// UsageSourceLive means the counts were computed on request from the
	// tenant's calls and evidence events; UsageID is empty and
	// CreatedAt/UpdatedAt are the time it was computed.
	UsageSourceLive UsageSource = "live"
)

// UsageRecord is aggregated tenant usage for a period.
type UsageRecord struct {
	CallsInbound        int         `json:"calls_inbound,omitempty"`
	CallsOutbound       int         `json:"calls_outbound,omitempty"`
	CallsTotal          int         `json:"calls_total,omitempty"`
	CreatedAt           string      `json:"created_at,omitempty"`
	EvidenceEventsCount int         `json:"evidence_events_count,omitempty"`
	FailOpenCount       int         `json:"fail_open_count,omitempty"`
	LLMTokensIn         int         `json:"llm_tokens_in,omitempty"`
	LLMTokensOut        int         `json:"llm_tokens_out,omitempty"`
	PeriodEnd           string      `json:"period_end,omitempty"`
	PeriodStart         string      `json:"period_start,omitempty"`
	Source              UsageSource `json:"source,omitempty"`
	TenantID            string      `json:"tenant_id,omitempty"`
	UpdatedAt           string      `json:"updated_at,omitempty"`
	UsageID             string      `json:"usage_id,omitempty"`
	ValidatorChecks     int         `json:"validator_checks,omitempty"`
}

// TenantListParams are the pagination parameters of the tenant list calls.
type TenantListParams struct {
	Cursor string
	Limit  int
}

// ListTenantKeysParams are the optional parameters of Tenants.ListKeys.
type ListTenantKeysParams = TenantListParams

// ListConfigurationsParams are the optional parameters of Tenants.ListConfigurations.
type ListConfigurationsParams = TenantListParams

// TenantAPIKey is an API key of a tenant (without the secret).
type TenantAPIKey struct {
	APIKeyID  string   `json:"api_key_id,omitempty"`
	CreatedAt string   `json:"created_at,omitempty"`
	ExpiresAt string   `json:"expires_at,omitempty"`
	KeyPrefix string   `json:"key_prefix,omitempty"`
	Name      string   `json:"name,omitempty"`
	RevokedAt string   `json:"revoked_at,omitempty"`
	Role      string   `json:"role,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
	TenantID  string   `json:"tenant_id,omitempty"`
}

// ListKeysResponse is the response of Tenants.ListKeys.
type ListKeysResponse struct {
	HasMore    bool           `json:"has_more,omitempty"`
	Keys       []TenantAPIKey `json:"keys,omitempty"`
	NextCursor string         `json:"next_cursor,omitempty"`
}

// GenerateKeyRequest is the body of Tenants.GenerateKey.
type GenerateKeyRequest struct {
	Name string `json:"name,omitempty"`
	Role string `json:"role,omitempty"`
}

// GenerateKeyResponse is the response of Tenants.GenerateKey; Key is returned
// only once.
type GenerateKeyResponse struct {
	APIKeyID  string   `json:"api_key_id,omitempty"`
	CreatedAt string   `json:"created_at,omitempty"`
	ExpiresAt string   `json:"expires_at,omitempty"`
	Key       string   `json:"key,omitempty"`
	KeyPrefix string   `json:"key_prefix,omitempty"`
	Name      string   `json:"name,omitempty"`
	RevokedAt string   `json:"revoked_at,omitempty"`
	Role      string   `json:"role,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
	TenantID  string   `json:"tenant_id,omitempty"`
}

// BlocklistPatternRequest is one blocklist pattern of a configuration.
type BlocklistPatternRequest struct {
	Category        string   `json:"category,omitempty"`
	ID              string   `json:"id,omitempty"`
	ReplacementText string   `json:"replacement_text,omitempty"`
	Severity        string   `json:"severity,omitempty"`
	Terms           []string `json:"terms,omitempty"`
}

// Configuration is a tenant configuration. Free-form sub-objects are kept as
// raw JSON.
type Configuration struct {
	ASRConfig                   json.RawMessage           `json:"asr_config,omitempty"`
	AutoHandoverConfig          json.RawMessage           `json:"auto_handover_config,omitempty"`
	BeginMessage                string                    `json:"begin_message,omitempty"`
	BlocklistPatterns           json.RawMessage           `json:"blocklist_patterns,omitempty"`
	ConfigID                    string                    `json:"config_id,omitempty"`
	ConfigSealHash              string                    `json:"config_seal_hash,omitempty"`
	ConsentConfig               json.RawMessage           `json:"consent_config,omitempty"`
	Countries                   []string                  `json:"countries,omitempty"`
	CreatedAt                   string                    `json:"created_at,omitempty"`
	DeletedAt                   string                    `json:"deleted_at,omitempty"`
	DisclosureTexts             json.RawMessage           `json:"disclosure_texts,omitempty"`
	EndCall                     *EndCallConfig            `json:"end_call,omitempty"`
	FallbackResponse            json.RawMessage           `json:"fallback_response,omitempty"`
	LLMModel                    string                    `json:"llm_model,omitempty"`
	LLMProvider                 string                    `json:"llm_provider,omitempty"`
	MultiCandidate              json.RawMessage           `json:"multi_candidate,omitempty"`
	Name                        string                    `json:"name,omitempty"`
	OutputConstraint            json.RawMessage           `json:"output_constraint,omitempty"`
	PIIEncryptionEnabled        bool                      `json:"pii_encryption_enabled"`
	PostCallExtraction          *PostCallExtractionConfig `json:"post_call_extraction,omitempty"`
	PostCallSummary             *PostCallSummaryConfig    `json:"post_call_summary,omitempty"`
	PricingBands                json.RawMessage           `json:"pricing_bands,omitempty"`
	PromptConfidentialityConfig json.RawMessage           `json:"prompt_confidentiality_config,omitempty"`
	RecordCalls                 bool                      `json:"record_calls,omitempty"`
	RecordingRetentionDays      int                       `json:"recording_retention_days,omitempty"`
	StartSpeaker                string                    `json:"start_speaker,omitempty"`
	TenantID                    string                    `json:"tenant_id,omitempty"`
	UpdatedAt                   string                    `json:"updated_at,omitempty"`
	UpstreamLLMSecretARN        string                    `json:"upstream_llm_secret_arn,omitempty"`
	UpstreamLLMURL              string                    `json:"upstream_llm_url,omitempty"`
	ValidatorModes              json.RawMessage           `json:"validator_modes,omitempty"`
	Version                     int                       `json:"version,omitempty"`
	Voicemail                   *VoicemailConfig          `json:"voicemail,omitempty"`
}

// ListConfigurationsResponse is the response of Tenants.ListConfigurations.
type ListConfigurationsResponse struct {
	Configurations []Configuration `json:"configurations,omitempty"`
	HasMore        bool            `json:"has_more,omitempty"`
	NextCursor     string          `json:"next_cursor,omitempty"`
}

// CreateConfigurationRequest is the body of Tenants.CreateConfiguration.
type CreateConfigurationRequest struct {
	ASRConfig                   json.RawMessage           `json:"asr_config,omitempty"`
	AutoHandoverConfig          json.RawMessage           `json:"auto_handover_config,omitempty"`
	BlocklistPatterns           []BlocklistPatternRequest `json:"blocklist_patterns,omitempty"`
	ConsentConfig               json.RawMessage           `json:"consent_config,omitempty"`
	Countries                   []string                  `json:"countries,omitempty"`
	DisclosureTexts             json.RawMessage           `json:"disclosure_texts,omitempty"`
	EndCall                     *EndCallConfig            `json:"end_call,omitempty"`
	FallbackResponse            json.RawMessage           `json:"fallback_response,omitempty"`
	MultiCandidate              json.RawMessage           `json:"multi_candidate,omitempty"`
	Name                        string                    `json:"name,omitempty"`
	OutputConstraint            json.RawMessage           `json:"output_constraint,omitempty"`
	PricingBands                json.RawMessage           `json:"pricing_bands,omitempty"`
	PromptConfidentialityConfig json.RawMessage           `json:"prompt_confidentiality_config,omitempty"`
	RecordCalls                 *bool                     `json:"record_calls,omitempty"`
	RecordingRetentionDays      *int                      `json:"recording_retention_days,omitempty"`
	TenantID                    string                    `json:"tenant_id,omitempty"`
	UpstreamLLMKey              string                    `json:"upstream_llm_key,omitempty"`
	UpstreamLLMSecretARN        string                    `json:"upstream_llm_secret_arn,omitempty"`
	UpstreamLLMURL              string                    `json:"upstream_llm_url,omitempty"`
	ValidatorModes              json.RawMessage           `json:"validator_modes,omitempty"`
	Voicemail                   *VoicemailConfig          `json:"voicemail,omitempty"`
}

// UpdateConfigurationRequest is the body of Tenants.UpdateConfiguration.
type UpdateConfigurationRequest struct {
	ASRConfig                   json.RawMessage           `json:"asr_config,omitempty"`
	AutoHandoverConfig          json.RawMessage           `json:"auto_handover_config,omitempty"`
	BlocklistPatterns           []BlocklistPatternRequest `json:"blocklist_patterns,omitempty"`
	ConsentConfig               json.RawMessage           `json:"consent_config,omitempty"`
	Countries                   []string                  `json:"countries,omitempty"`
	DisclosureTexts             json.RawMessage           `json:"disclosure_texts,omitempty"`
	EndCall                     *EndCallConfig            `json:"end_call,omitempty"`
	FallbackResponse            json.RawMessage           `json:"fallback_response,omitempty"`
	MultiCandidate              json.RawMessage           `json:"multi_candidate,omitempty"`
	Name                        string                    `json:"name,omitempty"`
	OutputConstraint            json.RawMessage           `json:"output_constraint,omitempty"`
	PricingBands                json.RawMessage           `json:"pricing_bands,omitempty"`
	PromptConfidentialityConfig json.RawMessage           `json:"prompt_confidentiality_config,omitempty"`
	RecordCalls                 *bool                     `json:"record_calls,omitempty"`
	RecordingRetentionDays      *int                      `json:"recording_retention_days,omitempty"`
	UpstreamLLMKey              string                    `json:"upstream_llm_key,omitempty"`
	UpstreamLLMSecretARN        string                    `json:"upstream_llm_secret_arn,omitempty"`
	UpstreamLLMURL              string                    `json:"upstream_llm_url,omitempty"`
	ValidatorModes              json.RawMessage           `json:"validator_modes,omitempty"`
	Voicemail                   *VoicemailConfig          `json:"voicemail,omitempty"`
}
