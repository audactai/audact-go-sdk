package audact

// PlatformBundleFormat selects the serialisation of an exported bundle.
type PlatformBundleFormat string

// Bundle formats accepted by the export routes.
const (
	// PlatformBundleYAML requests application/yaml (the default).
	PlatformBundleYAML PlatformBundleFormat = "yaml"
	// PlatformBundleJSON requests application/json.
	PlatformBundleJSON PlatformBundleFormat = "json"
)

// PlatformImportMode is the conflict mode of a bundle import.
type PlatformImportMode string

// PlatformBundle is a raw exported bundle. The Go SDK uses only the standard
// library, so the YAML/JSON document is returned undecoded.
type PlatformBundle struct {
	// Data is the bundle document exactly as the server returned it.
	Data []byte
	// ContentType is the response Content-Type (application/yaml or application/json).
	ContentType string
}

// PlatformExportOptions are the options of ExportTenant.
type PlatformExportOptions struct {
	// Format is the requested bundle format; empty means YAML.
	Format PlatformBundleFormat
}

// PlatformExportAgentOptions are the options of ExportAgent.
type PlatformExportAgentOptions struct {
	// TenantID is the tenant owning the agent (the tenant_id query parameter).
	TenantID string
	// Format is the requested bundle format; empty means YAML.
	Format PlatformBundleFormat
}

// PlatformImportOptions are the options of ImportTenant.
type PlatformImportOptions struct {
	// DryRun validates the bundle without writing (dryRun=true).
	DryRun bool
	// Mode is the import conflict mode (mode query parameter).
	Mode PlatformImportMode
	// SetSecret holds NAME=value assignments sent as repeated
	// X-Audact-Set-Secret headers (never as query parameters).
	SetSecret []string
	// IdempotencyKey is sent as the Idempotency-Key header when non-empty.
	IdempotencyKey string
}

// PlatformImportEntityResult is the outcome for one imported entity.
type PlatformImportEntityResult struct {
	Kind        string `json:"kind,omitempty"`
	Name        string `json:"name,omitempty"`
	NeedsSecret bool   `json:"needs_secret,omitempty"`
	NewID       string `json:"new_id,omitempty"`
	Ref         string `json:"ref,omitempty"`
	Status      string `json:"status,omitempty"`
}

// PlatformImportResponse is the result of a bundle import.
type PlatformImportResponse struct {
	BundleSHA256 string                       `json:"bundle_sha256,omitempty"`
	Counts       map[string]int               `json:"counts,omitempty"`
	DryRun       bool                         `json:"dry_run,omitempty"`
	Entities     []PlatformImportEntityResult `json:"entities,omitempty"`
	Mode         string                       `json:"mode,omitempty"`
	Warnings     []string                     `json:"warnings,omitempty"`
}

// AdminKey is a platform admin signing key.
type AdminKey struct {
	Active     bool   `json:"active,omitempty"`
	AddedAt    string `json:"added_at,omitempty"`
	DisabledAt string `json:"disabled_at,omitempty"`
	KID        string `json:"kid,omitempty"`
	Operator   string `json:"operator,omitempty"`
}

// ListAdminKeysResponse lists platform admin keys.
type ListAdminKeysResponse struct {
	AdminKeys []AdminKey `json:"admin_keys,omitempty"`
}

// AddAdminKeyRequest registers a platform admin public key.
type AddAdminKeyRequest struct {
	KID          string `json:"kid,omitempty"`
	Operator     string `json:"operator,omitempty"`
	PublicKeyPEM string `json:"public_key_pem,omitempty"`
}

// AdminKeyResponse wraps one platform admin key.
type AdminKeyResponse struct {
	AdminKey *AdminKey `json:"admin_key,omitempty"`
}

// PlatformListOrganizationsParams are the query parameters of Orgs.List.
type PlatformListOrganizationsParams struct {
	Status string
	Email  string
	Limit  int
	Cursor string
}

// PlatformListOrganizationsResponse is a page of organisations.
type PlatformListOrganizationsResponse struct {
	NextCursor    string         `json:"next_cursor,omitempty"`
	Organizations []Organization `json:"organizations,omitempty"`
}

// PlatformOrganizationResponse wraps one organisation.
type PlatformOrganizationResponse struct {
	Organization *Organization `json:"organization,omitempty"`
}

// PlatformRejectOrganizationRequest rejects an organisation.
type PlatformRejectOrganizationRequest struct {
	Reason string `json:"reason"`
}

// PlatformFindOrgAdminsParams are the query parameters of Orgs.FindAdmins.
type PlatformFindOrgAdminsParams struct {
	Email string
	OrgID string
}

// PlatformFindOrgAdminsResponse lists matching org admins.
type PlatformFindOrgAdminsResponse struct {
	OrgAdmins []OrgAdmin `json:"org_admins,omitempty"`
}

// PlatformChangeOrgAdminEmailRequest changes an org admin's email.
type PlatformChangeOrgAdminEmailRequest struct {
	Email string `json:"email"`
}

// PlatformOrgAdminResponse wraps one org admin.
type PlatformOrgAdminResponse struct {
	OrgAdmin *OrgAdmin `json:"org_admin,omitempty"`
}

// PlatformResetOrgAdminPasswordRequest resets an org admin's password.
type PlatformResetOrgAdminPasswordRequest struct {
	NewPassword string `json:"new_password,omitempty"`
}

// PlatformResetOrgAdminPasswordResponse is the result of a password reset.
type PlatformResetOrgAdminPasswordResponse struct {
	ExpiresIn int    `json:"expires_in,omitempty"`
	Reset     string `json:"reset,omitempty"`
	Token     string `json:"token,omitempty"`
}
