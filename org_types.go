package audact

// OrgAdmin is an organisation administrator account.
type OrgAdmin struct {
	AdminID         string `json:"admin_id,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	DisabledAt      string `json:"disabled_at,omitempty"`
	Email           string `json:"email,omitempty"`
	EmailVerifiedAt string `json:"email_verified_at,omitempty"`
	LastLoginAt     string `json:"last_login_at,omitempty"`
	Name            string `json:"name,omitempty"`
	OrgID           string `json:"org_id,omitempty"`
	Role            string `json:"role,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
}

// Organization is an Audact organisation.
type Organization struct {
	Country            string `json:"country,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	DefaultBillingRef  string `json:"default_billing_ref,omitempty"`
	EntityType         string `json:"entity_type,omitempty"`
	LegalName          string `json:"legal_name,omitempty"`
	Name               string `json:"name,omitempty"`
	OrgID              string `json:"org_id,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
	VerificationStatus string `json:"verification_status,omitempty"`
	VerifiedAt         string `json:"verified_at,omitempty"`
	VerifiedByKID      string `json:"verified_by_kid,omitempty"`
}

// OrgKYC is an organisation's KYC record.
type OrgKYC struct {
	Additional           map[string]any `json:"additional,omitempty"`
	AddressCountry       string         `json:"address_country,omitempty"`
	AddressLine1         string         `json:"address_line1,omitempty"`
	AddressLine2         string         `json:"address_line2,omitempty"`
	City                 string         `json:"city,omitempty"`
	ContactEmail         string         `json:"contact_email,omitempty"`
	ContactPhone         string         `json:"contact_phone,omitempty"`
	CreatedAt            string         `json:"created_at,omitempty"`
	DateOfBirth          string         `json:"date_of_birth,omitempty"`
	IncorporationCountry string         `json:"incorporation_country,omitempty"`
	IncorporationDate    string         `json:"incorporation_date,omitempty"`
	IncorporationRegion  string         `json:"incorporation_region,omitempty"`
	LEI                  string         `json:"lei,omitempty"`
	Nationality          string         `json:"nationality,omitempty"`
	OrgID                string         `json:"org_id,omitempty"`
	PostalCode           string         `json:"postal_code,omitempty"`
	Region               string         `json:"region,omitempty"`
	RegistrationNumber   string         `json:"registration_number,omitempty"`
	RepresentativeName   string         `json:"representative_name,omitempty"`
	RepresentativeRole   string         `json:"representative_role,omitempty"`
	TaxID                string         `json:"tax_id,omitempty"`
	TaxResidenceCountry  string         `json:"tax_residence_country,omitempty"`
	UpdatedAt            string         `json:"updated_at,omitempty"`
	VATID                string         `json:"vat_id,omitempty"`
	Website              string         `json:"website,omitempty"`
}

// OrgKYCPatchRequest holds the KYC fields to merge; unset fields are not sent.
type OrgKYCPatchRequest struct {
	Additional           map[string]any `json:"additional,omitempty"`
	AddressCountry       string         `json:"address_country,omitempty"`
	AddressLine1         string         `json:"address_line1,omitempty"`
	AddressLine2         string         `json:"address_line2,omitempty"`
	City                 string         `json:"city,omitempty"`
	ContactEmail         string         `json:"contact_email,omitempty"`
	ContactPhone         string         `json:"contact_phone,omitempty"`
	DateOfBirth          string         `json:"date_of_birth,omitempty"`
	IncorporationCountry string         `json:"incorporation_country,omitempty"`
	IncorporationDate    string         `json:"incorporation_date,omitempty"`
	IncorporationRegion  string         `json:"incorporation_region,omitempty"`
	LEI                  string         `json:"lei,omitempty"`
	Nationality          string         `json:"nationality,omitempty"`
	PostalCode           string         `json:"postal_code,omitempty"`
	Region               string         `json:"region,omitempty"`
	RegistrationNumber   string         `json:"registration_number,omitempty"`
	RepresentativeName   string         `json:"representative_name,omitempty"`
	RepresentativeRole   string         `json:"representative_role,omitempty"`
	TaxID                string         `json:"tax_id,omitempty"`
	TaxResidenceCountry  string         `json:"tax_residence_country,omitempty"`
	VATID                string         `json:"vat_id,omitempty"`
	Website              string         `json:"website,omitempty"`
}

// OrgKYCResponse wraps an organisation's KYC record.
type OrgKYCResponse struct {
	KYC *OrgKYC `json:"kyc,omitempty"`
}

// OrgEntityType is the kind of legal entity behind an organisation.
type OrgEntityType string

// Organisation entity types.
const (
	OrgEntityIndividual OrgEntityType = "individual"
	OrgEntityCompany    OrgEntityType = "company"
)

// OrgSignupRequest creates an organisation and its first admin.
type OrgSignupRequest struct {
	Country    string        `json:"country,omitempty"`
	Email      string        `json:"email,omitempty"`
	EntityType OrgEntityType `json:"entity_type,omitempty"`
	KYC        *OrgKYC       `json:"kyc,omitempty"`
	LegalName  string        `json:"legal_name,omitempty"`
	Name       string        `json:"name,omitempty"`
	OrgName    string        `json:"org_name,omitempty"`
	Password   string        `json:"password,omitempty"`
}

// OrgSignupResponse is returned by signup.
type OrgSignupResponse struct {
	Admin     *OrgAdmin     `json:"admin,omitempty"`
	ExpiresIn int           `json:"expires_in,omitempty"`
	Org       *Organization `json:"org,omitempty"`
	Token     string        `json:"token,omitempty"`
}

// OrgLoginRequest authenticates an organisation admin.
type OrgLoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// OrgLoginResponse is returned by login.
type OrgLoginResponse struct {
	ExpiresIn int    `json:"expires_in,omitempty"`
	Token     string `json:"token,omitempty"`
	Ver       bool   `json:"ver,omitempty"`
}

// OrgRefreshResponse is returned by refresh.
type OrgRefreshResponse struct {
	ExpiresIn int    `json:"expires_in,omitempty"`
	Token     string `json:"token,omitempty"`
}

// OrgProfileResponse is the calling organisation and admin.
type OrgProfileResponse struct {
	Admin *OrgAdmin     `json:"admin,omitempty"`
	Org   *Organization `json:"org,omitempty"`
}

// OrgUpdateProfileRequest updates display names; empty fields are unchanged.
type OrgUpdateProfileRequest struct {
	AdminName string `json:"admin_name,omitempty"`
	OrgName   string `json:"org_name,omitempty"`
}

// OrgTenantSummary describes a tenant owned by an organisation.
type OrgTenantSummary struct {
	CreatedAt     string `json:"created_at,omitempty"`
	DeletedAt     string `json:"deleted_at,omitempty"`
	InitializedAt string `json:"initialized_at,omitempty"`
	Name          string `json:"name,omitempty"`
	Status        string `json:"status,omitempty"`
	TenantID      string `json:"tenant_id,omitempty"`
}

// ListOrgTenantsParams are optional query parameters for Org.ListTenants.
type ListOrgTenantsParams struct {
	IncludeDeleted bool
	Limit          int
	Cursor         string
}

// OrgListTenantsResponse is a page of organisation tenants.
type OrgListTenantsResponse struct {
	NextCursor string             `json:"next_cursor,omitempty"`
	Tenants    []OrgTenantSummary `json:"tenants,omitempty"`
}

// OrgCreateTenantRequest creates a tenant under the organisation.
type OrgCreateTenantRequest struct {
	Name string `json:"name,omitempty"`
}

// OrgCreateTenantResponse is returned by CreateTenant.
type OrgCreateTenantResponse struct {
	Tenant *OrgTenantSummary `json:"tenant,omitempty"`
}

// OrgGetTenantResponse is returned by GetTenant.
type OrgGetTenantResponse struct {
	Tenant *OrgTenantSummary `json:"tenant,omitempty"`
}

// OrgTenantActionResponse is returned by tenant lifecycle actions.
type OrgTenantActionResponse struct {
	Tenant *OrgTenantSummary `json:"tenant,omitempty"`
}

// OrgImpersonateTenantResponse carries a short-lived tenant token.
type OrgImpersonateTenantResponse struct {
	ExpiresIn   int    `json:"expires_in,omitempty"`
	TenantToken string `json:"tenant_token,omitempty"`
}

// OrgAPIKeyView is a newly minted API key; Key is returned only once.
type OrgAPIKeyView struct {
	APIKeyID  string `json:"api_key_id,omitempty"`
	Key       string `json:"key,omitempty"`
	KeyPrefix string `json:"key_prefix,omitempty"`
	Role      string `json:"role,omitempty"`
}

// OrgInitializeTenantResponse is returned by InitializeTenant.
type OrgInitializeTenantResponse struct {
	APIKey *OrgAPIKeyView    `json:"api_key,omitempty"`
	Tenant *OrgTenantSummary `json:"tenant,omitempty"`
}

// OrgIssueTenantKeyRequest names a new tenant key.
type OrgIssueTenantKeyRequest struct {
	Name string `json:"name,omitempty"`
}

// OrgIssueTenantKeyResponse is returned by IssueTenantKey.
type OrgIssueTenantKeyResponse struct {
	APIKey *OrgAPIKeyView `json:"api_key,omitempty"`
}
