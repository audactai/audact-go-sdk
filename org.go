package audact

import (
	"context"
	"net/http"
	"net/url"
)

func orgTenantPath(id, suffix string) string {
	return "/org/tenants/" + url.PathEscape(id) + suffix
}

// Signup creates an organisation and its first admin, then stores the returned
// access token on the client. No authentication is sent. Calls POST /v1/org/signup.
func (s *OrgService) Signup(ctx context.Context, req *OrgSignupRequest) (*OrgSignupResponse, error) {
	var out OrgSignupResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/org/signup", authNone, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	s.t.setOrgToken(out.Token)
	return &out, nil
}

// Login authenticates an org admin and stores the returned access token on the
// client. No authentication is sent. Calls POST /v1/org/login.
func (s *OrgService) Login(ctx context.Context, req *OrgLoginRequest) (*OrgLoginResponse, error) {
	var out OrgLoginResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/org/login", authNone, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	s.t.setOrgToken(out.Token)
	return &out, nil
}

// Refresh rotates the session using the refresh cookie (configure a cookie jar
// via WithHTTPClient) and stores the new access token. Calls POST /v1/org/refresh.
func (s *OrgService) Refresh(ctx context.Context) (*OrgRefreshResponse, error) {
	var out OrgRefreshResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/org/refresh", authNone, nil, nil, &out); err != nil {
		return nil, err
	}
	s.t.setOrgToken(out.Token)
	return &out, nil
}

// Logout revokes the session and clears the stored org token. Calls POST /v1/org/logout.
func (s *OrgService) Logout(ctx context.Context) error {
	if _, err := s.t.do(ctx, http.MethodPost, "/org/logout", authOrg, nil, nil, nil); err != nil {
		return err
	}
	s.t.setOrgToken("")
	return nil
}

// Profile returns the calling organisation and admin. Calls GET /v1/org/profile.
func (s *OrgService) Profile(ctx context.Context) (*OrgProfileResponse, error) {
	var out OrgProfileResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/org/profile", authOrg, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateProfile updates the organisation and/or admin display name. Calls
// PATCH /v1/org/profile.
func (s *OrgService) UpdateProfile(ctx context.Context, req *OrgUpdateProfileRequest) (*OrgProfileResponse, error) {
	var out OrgProfileResponse
	if _, err := s.t.do(ctx, http.MethodPatch, "/org/profile", authOrg, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// KYC returns the organisation's KYC record. Calls GET /v1/org/kyc.
func (s *OrgService) KYC(ctx context.Context) (*OrgKYCResponse, error) {
	var out OrgKYCResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/org/kyc", authOrg, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateKYC merges the given fields onto the KYC record. Calls PATCH /v1/org/kyc.
func (s *OrgService) UpdateKYC(ctx context.Context, req *OrgKYCPatchRequest) (*OrgKYCResponse, error) {
	var out OrgKYCResponse
	if _, err := s.t.do(ctx, http.MethodPatch, "/org/kyc", authOrg, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListTenants lists the organisation's tenants; params may be nil. Calls
// GET /v1/org/tenants.
func (s *OrgService) ListTenants(ctx context.Context, params *ListOrgTenantsParams) (*OrgListTenantsResponse, error) {
	q := url.Values{}
	if params != nil {
		addBool(q, "include_deleted", params.IncludeDeleted)
		addInt(q, "limit", params.Limit)
		addString(q, "cursor", params.Cursor)
	}
	var out OrgListTenantsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/org/tenants", authOrg, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateTenant creates a tenant under the organisation. Calls POST /v1/org/tenants.
func (s *OrgService) CreateTenant(ctx context.Context, req *OrgCreateTenantRequest) (*OrgCreateTenantResponse, error) {
	var out OrgCreateTenantResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/org/tenants", authOrg, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTenant retrieves an organisation tenant. Calls GET /v1/org/tenants/{id}.
func (s *OrgService) GetTenant(ctx context.Context, id string) (*OrgGetTenantResponse, error) {
	var out OrgGetTenantResponse
	if _, err := s.t.do(ctx, http.MethodGet, orgTenantPath(id, ""), authOrg, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *OrgService) tenantAction(ctx context.Context, method, id, suffix string) (*OrgTenantActionResponse, error) {
	var out OrgTenantActionResponse
	if _, err := s.t.do(ctx, method, orgTenantPath(id, suffix), authOrg, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteTenant soft-deletes a tenant. Calls DELETE /v1/org/tenants/{id}.
func (s *OrgService) DeleteTenant(ctx context.Context, id string) (*OrgTenantActionResponse, error) {
	return s.tenantAction(ctx, http.MethodDelete, id, "")
}

// DisableTenant disables a tenant. Calls POST /v1/org/tenants/{id}/disable.
func (s *OrgService) DisableTenant(ctx context.Context, id string) (*OrgTenantActionResponse, error) {
	return s.tenantAction(ctx, http.MethodPost, id, "/disable")
}

// EnableTenant enables a tenant. Calls POST /v1/org/tenants/{id}/enable.
func (s *OrgService) EnableTenant(ctx context.Context, id string) (*OrgTenantActionResponse, error) {
	return s.tenantAction(ctx, http.MethodPost, id, "/enable")
}

// RestoreTenant restores a soft-deleted tenant. Calls POST /v1/org/tenants/{id}/restore.
func (s *OrgService) RestoreTenant(ctx context.Context, id string) (*OrgTenantActionResponse, error) {
	return s.tenantAction(ctx, http.MethodPost, id, "/restore")
}

// ImpersonateTenant mints a short-lived tenant token. Calls
// POST /v1/org/tenants/{id}/impersonate.
func (s *OrgService) ImpersonateTenant(ctx context.Context, id string) (*OrgImpersonateTenantResponse, error) {
	var out OrgImpersonateTenantResponse
	if _, err := s.t.do(ctx, http.MethodPost, orgTenantPath(id, "/impersonate"), authOrg, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// InitializeTenant provisions the tenant's default configuration and first key;
// a non-empty idempotencyKey is sent as Idempotency-Key. Calls
// POST /v1/org/tenants/{id}/initialize.
func (s *OrgService) InitializeTenant(ctx context.Context, id, idempotencyKey string) (*OrgInitializeTenantResponse, error) {
	var out OrgInitializeTenantResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, orgTenantPath(id, "/initialize"), authOrg, nil, idemHeader(idempotencyKey), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// IssueTenantKey issues an additional tenant API key (returned once); a
// non-empty idempotencyKey is sent as Idempotency-Key. Calls
// POST /v1/org/tenants/{id}/keys.
func (s *OrgService) IssueTenantKey(ctx context.Context, id string, req *OrgIssueTenantKeyRequest, idempotencyKey string) (*OrgIssueTenantKeyResponse, error) {
	var out OrgIssueTenantKeyResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, orgTenantPath(id, "/keys"), authOrg, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}
