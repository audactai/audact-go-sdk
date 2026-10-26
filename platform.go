package audact

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// Platform routes send the platform token (WithPlatformToken or
// Client.SetPlatformToken). Without one they fail locally with
// ErrPlatformTokenRequired; the API key is never sent to a platform route.
//
// Platform tokens are single-use (the server burns each token's jti) and
// action-scoped ("org" for organisations and org admins, "admin_keys" for
// admin keys, "export"/"import" for bundles), so set a freshly minted token
// with SetPlatformToken before every Platform call.

func acceptFor(f PlatformBundleFormat) string {
	if f == PlatformBundleJSON {
		return "application/json"
	}
	return "application/yaml"
}

func (s *PlatformService) export(ctx context.Context, path string, q url.Values, f PlatformBundleFormat) (*PlatformBundle, error) {
	var raw rawResponse
	resp, err := s.t.doHeaders(ctx, http.MethodGet, path, authPlatform, q, map[string]string{"Accept": acceptFor(f)}, nil, &raw)
	if err != nil {
		return nil, err
	}
	return &PlatformBundle{Data: raw.data, ContentType: resp.Header.Get("Content-Type")}, nil
}

// ExportAgent exports one agent as a bundle. Calls GET /v1/platform/export/agents/{agent_id}.
func (s *PlatformService) ExportAgent(ctx context.Context, agentID string, opts *PlatformExportAgentOptions) (*PlatformBundle, error) {
	o := orEmpty(opts)
	q := url.Values{}
	addString(q, "tenant_id", o.TenantID)
	return s.export(ctx, "/platform/export/agents/"+url.PathEscape(agentID), q, o.Format)
}

// ExportTenant exports a whole tenant as a bundle. Calls GET /v1/platform/export/tenants/{tenant_id}.
func (s *PlatformService) ExportTenant(ctx context.Context, tenantID string, opts *PlatformExportOptions) (*PlatformBundle, error) {
	o := orEmpty(opts)
	return s.export(ctx, "/platform/export/tenants/"+url.PathEscape(tenantID), nil, o.Format)
}

// ImportTenant imports a raw YAML (or JSON) bundle into a tenant. The bundle is
// sent verbatim as application/yaml. Calls POST /v1/platform/import/tenants/{target_tenant_id}.
func (s *PlatformService) ImportTenant(ctx context.Context, targetTenantID string, bundle []byte, opts *PlatformImportOptions) (*PlatformImportResponse, error) {
	o := orEmpty(opts)
	q := url.Values{}
	addBool(q, "dryRun", o.DryRun)
	addString(q, "mode", string(o.Mode))
	hdr := http.Header{}
	for _, v := range o.SetSecret {
		hdr.Add("X-Audact-Set-Secret", v)
	}
	extra := map[string]string{"Idempotency-Key": o.IdempotencyKey}
	in := &rawBody{contentType: "application/yaml", data: bundle, header: hdr}
	var out PlatformImportResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/platform/import/tenants/"+url.PathEscape(targetTenantID), authPlatform, q, extra, in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAdminKeys lists platform admin keys. Calls GET /v1/platform/admin-keys.
func (s *PlatformService) ListAdminKeys(ctx context.Context) (*ListAdminKeysResponse, error) {
	var out ListAdminKeysResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/platform/admin-keys", authPlatform, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddAdminKey registers a platform admin key. Calls POST /v1/platform/admin-keys.
func (s *PlatformService) AddAdminKey(ctx context.Context, req AddAdminKeyRequest) (*AdminKeyResponse, error) {
	var out AdminKeyResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/platform/admin-keys", authPlatform, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DisableAdminKey disables a platform admin key. Calls POST /v1/platform/admin-keys/{kid}/disable.
func (s *PlatformService) DisableAdminKey(ctx context.Context, kid string) (*AdminKeyResponse, error) {
	var out AdminKeyResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/platform/admin-keys/"+url.PathEscape(kid)+"/disable", authPlatform, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List lists organisations. Calls GET /v1/platform/organizations.
func (s *PlatformOrgsService) List(ctx context.Context, params *PlatformListOrganizationsParams) (*PlatformListOrganizationsResponse, error) {
	p := orEmpty(params)
	q := url.Values{}
	addString(q, "status", p.Status)
	addString(q, "email", p.Email)
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	addString(q, "cursor", p.Cursor)
	var out PlatformListOrganizationsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/platform/organizations", authPlatform, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Verify verifies an organisation. Calls POST /v1/platform/organizations/{org_id}/verify.
func (s *PlatformOrgsService) Verify(ctx context.Context, orgID string) (*PlatformOrganizationResponse, error) {
	var out PlatformOrganizationResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/platform/organizations/"+url.PathEscape(orgID)+"/verify", authPlatform, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Reject rejects an organisation. Calls POST /v1/platform/organizations/{org_id}/reject.
func (s *PlatformOrgsService) Reject(ctx context.Context, orgID string, req PlatformRejectOrganizationRequest) (*PlatformOrganizationResponse, error) {
	var out PlatformOrganizationResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/platform/organizations/"+url.PathEscape(orgID)+"/reject", authPlatform, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FindAdmins finds org admins by email or org. Calls GET /v1/platform/org-admins.
func (s *PlatformOrgsService) FindAdmins(ctx context.Context, params *PlatformFindOrgAdminsParams) (*PlatformFindOrgAdminsResponse, error) {
	p := orEmpty(params)
	q := url.Values{}
	addString(q, "email", p.Email)
	addString(q, "org_id", p.OrgID)
	var out PlatformFindOrgAdminsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/platform/org-admins", authPlatform, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ChangeAdminEmail changes an org admin's email. Calls PATCH /v1/platform/org-admins/{admin_id}.
func (s *PlatformOrgsService) ChangeAdminEmail(ctx context.Context, adminID string, req PlatformChangeOrgAdminEmailRequest) (*PlatformOrgAdminResponse, error) {
	var out PlatformOrgAdminResponse
	if _, err := s.t.do(ctx, http.MethodPatch, "/platform/org-admins/"+url.PathEscape(adminID), authPlatform, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ResetAdminPassword resets an org admin's password. Calls POST /v1/platform/org-admins/{admin_id}/reset-password.
func (s *PlatformOrgsService) ResetAdminPassword(ctx context.Context, adminID string, req *PlatformResetOrgAdminPasswordRequest) (*PlatformResetOrgAdminPasswordResponse, error) {
	var out PlatformResetOrgAdminPasswordResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/platform/org-admins/"+url.PathEscape(adminID)+"/reset-password", authPlatform, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}
