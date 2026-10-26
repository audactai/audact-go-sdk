package audact

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

func tenantPath(tenantID, suffix string) string {
	return "/tenants/" + url.PathEscape(tenantID) + suffix
}

func tenantListQuery(p *TenantListParams) url.Values {
	q := url.Values{}
	if p != nil {
		addString(q, "cursor", p.Cursor)
		addInt(q, "limit", p.Limit)
	}
	return q
}

// Create creates a tenant and provisions its first API key (returned once);
// no authentication is sent. A non-empty idempotencyKey is sent as the
// Idempotency-Key header. Calls POST /v1/tenants.
func (s *TenantsService) Create(ctx context.Context, req *CreateTenantRequest, idempotencyKey string) (*CreateTenantResponse, error) {
	var out CreateTenantResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/tenants", authNone, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get retrieves a tenant. Calls GET /v1/tenants/{tenant_id}.
func (s *TenantsService) Get(ctx context.Context, tenantID string) (*Tenant, error) {
	var out Tenant
	if _, err := s.t.do(ctx, http.MethodGet, tenantPath(tenantID, ""), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates a tenant's name and/or plan. Calls PUT /v1/tenants/{tenant_id}.
func (s *TenantsService) Update(ctx context.Context, tenantID string, req *UpdateTenantRequest) (*Tenant, error) {
	var out Tenant
	if _, err := s.t.do(ctx, http.MethodPut, tenantPath(tenantID, ""), authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Usage returns aggregated usage for a tenant. params may be nil. Calls
// GET /v1/tenants/{tenant_id}/usage.
func (s *TenantsService) Usage(ctx context.Context, tenantID string, params *GetTenantUsageParams) (*UsageRecord, error) {
	q := url.Values{}
	if params != nil {
		addString(q, "period_start", params.PeriodStart)
		addString(q, "period_end", params.PeriodEnd)
	}
	var out UsageRecord
	if _, err := s.t.do(ctx, http.MethodGet, tenantPath(tenantID, "/usage"), authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListKeys lists a tenant's API keys. Calls GET /v1/tenants/{tenant_id}/keys.
func (s *TenantsService) ListKeys(ctx context.Context, tenantID string, params *ListTenantKeysParams) (*ListKeysResponse, error) {
	var out ListKeysResponse
	if _, err := s.t.do(ctx, http.MethodGet, tenantPath(tenantID, "/keys"), authAPIKey, tenantListQuery(params), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GenerateKey creates an API key; the secret is returned only once. Calls
// POST /v1/tenants/{tenant_id}/keys.
func (s *TenantsService) GenerateKey(ctx context.Context, tenantID string, req *GenerateKeyRequest) (*GenerateKeyResponse, error) {
	var out GenerateKeyResponse
	if _, err := s.t.do(ctx, http.MethodPost, tenantPath(tenantID, "/keys"), authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RevokeKey revokes an API key. Calls DELETE /v1/tenants/{tenant_id}/keys/{key_id}.
func (s *TenantsService) RevokeKey(ctx context.Context, tenantID, keyID string) error {
	_, err := s.t.do(ctx, http.MethodDelete, tenantPath(tenantID, "/keys/"+url.PathEscape(keyID)), authAPIKey, nil, nil, nil)
	return err
}

// ListConfigurations lists a tenant's configurations. Calls
// GET /v1/tenants/{tenant_id}/configurations.
func (s *TenantsService) ListConfigurations(ctx context.Context, tenantID string, params *ListConfigurationsParams) (*ListConfigurationsResponse, error) {
	var out ListConfigurationsResponse
	if _, err := s.t.do(ctx, http.MethodGet, tenantPath(tenantID, "/configurations"), authAPIKey, tenantListQuery(params), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateConfiguration creates a configuration. Calls POST /v1/configurations.
func (s *TenantsService) CreateConfiguration(ctx context.Context, req *CreateConfigurationRequest) (*Configuration, error) {
	var out Configuration
	if _, err := s.t.do(ctx, http.MethodPost, "/configurations", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetConfiguration retrieves a configuration. Calls GET /v1/configurations/{config_id}.
func (s *TenantsService) GetConfiguration(ctx context.Context, configID string) (*Configuration, error) {
	var out Configuration
	if _, err := s.t.do(ctx, http.MethodGet, "/configurations/"+url.PathEscape(configID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetConfigurationVersion retrieves one immutable earlier version of a
// configuration: the settings calls pinned to that version ran with. version
// is a receipt's config_version (>= 1). Calls
// GET /v1/configurations/{config_id}/versions/{version}.
func (s *TenantsService) GetConfigurationVersion(ctx context.Context, configID string, version int) (*Configuration, error) {
	var out Configuration
	path := "/configurations/" + url.PathEscape(configID) + "/versions/" + strconv.Itoa(version)
	if _, err := s.t.do(ctx, http.MethodGet, path, authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateConfiguration updates a configuration. Calls PUT /v1/configurations/{config_id}.
func (s *TenantsService) UpdateConfiguration(ctx context.Context, configID string, req *UpdateConfigurationRequest) (*Configuration, error) {
	var out Configuration
	if _, err := s.t.do(ctx, http.MethodPut, "/configurations/"+url.PathEscape(configID), authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteConfiguration deletes a configuration. Calls DELETE /v1/configurations/{config_id}.
func (s *TenantsService) DeleteConfiguration(ctx context.Context, configID string) error {
	_, err := s.t.do(ctx, http.MethodDelete, "/configurations/"+url.PathEscape(configID), authAPIKey, nil, nil, nil)
	return err
}
