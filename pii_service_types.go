package audact

// PIIIdentifierType is the kind of identifier a client hash was computed over.
type PIIIdentifierType string

// Identifier types accepted by the shred and link-identity routes.
const (
	// PIIIdentifierTypeName is a name hash.
	PIIIdentifierTypeName PIIIdentifierType = "name"
	// PIIIdentifierTypeEmail is an e-mail hash.
	PIIIdentifierTypeEmail PIIIdentifierType = "email"
	// PIIIdentifierTypePhone is a phone hash.
	PIIIdentifierTypePhone PIIIdentifierType = "phone"
)

// PIIExemptType is the kind of identifier an exemption covers.
type PIIExemptType string

// Exemption kinds.
const (
	// PIIExemptTypeEmail exempts an e-mail address.
	PIIExemptTypeEmail PIIExemptType = "email"
	// PIIExemptTypePhone exempts a phone number.
	PIIExemptTypePhone PIIExemptType = "phone"
)

// PIIIdentityType is the type of an identity linked to a call.
type PIIIdentityType string

// Identity types.
const (
	// PIIIdentityTypeName is a name identity.
	PIIIdentityTypeName PIIIdentityType = "name"
	// PIIIdentityTypeEmail is an e-mail identity.
	PIIIdentityTypeEmail PIIIdentityType = "email"
	// PIIIdentityTypePhone is a phone identity.
	PIIIdentityTypePhone PIIIdentityType = "phone"
	// PIIIdentityTypeCall is a call-scope identity.
	PIIIdentityTypeCall PIIIdentityType = "call"
)

// PIIIdentityRole is how an identity relates to a call.
type PIIIdentityRole string

// Identity roles.
const (
	// PIIIdentityRoleSubject marks the call's subject.
	PIIIdentityRoleSubject PIIIdentityRole = "subject"
	// PIIIdentityRoleMentioned marks a mentioned party.
	PIIIdentityRoleMentioned PIIIdentityRole = "mentioned"
)

// PIIIdentitySource is how an identity was linked to a call.
type PIIIdentitySource string

// Identity sources.
const (
	// PIIIdentitySourceCall is linked from the call itself.
	PIIIdentitySourceCall PIIIdentitySource = "call"
	// PIIIdentitySourcePSTN is linked from PSTN metadata.
	PIIIdentitySourcePSTN PIIIdentitySource = "pstn"
	// PIIIdentitySourceCallInit is linked at call initiation.
	PIIIdentitySourceCallInit PIIIdentitySource = "call_init"
	// PIIIdentitySourceAPI is linked through the API.
	PIIIdentitySourceAPI PIIIdentitySource = "api"
	// PIIIdentitySourceDetected is linked by detection.
	PIIIdentitySourceDetected PIIIdentitySource = "detected"
)

// PIIShredScope is the scope of an erasure.
type PIIShredScope string

// Shred scopes.
const (
	// PIIShredScopeIdentity is an identity-wide erasure.
	PIIShredScopeIdentity PIIShredScope = "identity"
	// PIIShredScopeCall is a single-call erasure.
	PIIShredScopeCall PIIShredScope = "call"
	// PIIShredScopeTenant is a tenant-wide erasure.
	PIIShredScopeTenant PIIShredScope = "tenant"
)

// PIIShredRequest is the body of POST /v1/pii/shred.
type PIIShredRequest struct {
	IdentifierType PIIIdentifierType `json:"identifier_type,omitempty"`
	IdentifierHash string            `json:"identifier_hash,omitempty"`
}

// PIIShredResponse is the 202 response of the shred routes.
type PIIShredResponse struct {
	CallsAffected      int    `json:"calls_affected,omitempty"`
	IdentitiesShredded int    `json:"identities_shredded,omitempty"`
	ShredID            string `json:"shred_id,omitempty"`
	ShreddedAt         string `json:"shredded_at,omitempty"`
}

// PIIShredEntry is one row of the PII-free erasure log.
type PIIShredEntry struct {
	ActorKeyID         string          `json:"actor_key_id,omitempty"`
	CallID             string          `json:"call_id,omitempty"`
	CallsAffected      int             `json:"calls_affected,omitempty"`
	IdentitiesShredded int             `json:"identities_shredded,omitempty"`
	IdentityType       PIIIdentityType `json:"identity_type,omitempty"`
	ImpersonatorID     string          `json:"impersonator_id,omitempty"`
	ImpersonatorOrgID  string          `json:"impersonator_org_id,omitempty"`
	Scope              PIIShredScope   `json:"scope,omitempty"`
	ShredID            string          `json:"shred_id,omitempty"`
	ShreddedAt         string          `json:"shredded_at,omitempty"`
}

// ListPIIShredsParams holds the query parameters of GET /v1/pii/shreds.
type ListPIIShredsParams struct {
	Cursor string
	Limit  int
}

// PIIShredsResponse is the response of GET /v1/pii/shreds.
type PIIShredsResponse struct {
	NextCursor string          `json:"next_cursor,omitempty"`
	Shreds     []PIIShredEntry `json:"shreds,omitempty"`
}

// PIIIdentity is one identity linked to a call; no hash or value is returned.
type PIIIdentity struct {
	IdentityID string            `json:"identity_id,omitempty"`
	LinkedAt   string            `json:"linked_at,omitempty"`
	Role       PIIIdentityRole   `json:"role,omitempty"`
	ShreddedAt string            `json:"shredded_at,omitempty"`
	Source     PIIIdentitySource `json:"source,omitempty"`
	Type       PIIIdentityType   `json:"type,omitempty"`
}

// PIIIdentitiesResponse is the response of GET /v1/calls/{call_id}/pii/identities.
type PIIIdentitiesResponse struct {
	Identities []PIIIdentity `json:"identities,omitempty"`
}

// PIILinkIdentityRequest is the body of POST /v1/calls/{call_id}/pii/identities.
type PIILinkIdentityRequest struct {
	IdentifierHash string            `json:"identifier_hash,omitempty"`
	IdentifierType PIIIdentifierType `json:"identifier_type,omitempty"`
}

// PIILinkIdentityResponse is the 201 response of the link-identity route.
type PIILinkIdentityResponse struct {
	IdentityID string `json:"identity_id,omitempty"`
}

// PIICallShredRequest is the body of POST /v1/calls/{call_id}/pii/shred.
// Omitting IdentityIDs shreds every identity of the call.
type PIICallShredRequest struct {
	IdentityIDs []string `json:"identity_ids,omitempty"`
}

// PIIDecryptRequest is the body of POST /v1/calls/{call_id}/pii/decrypt.
// Send exactly one of Text or Texts.
type PIIDecryptRequest struct {
	Text  string   `json:"text,omitempty"`
	Texts []string `json:"texts,omitempty"`
}

// PIIDecryptResponse is the response of the decrypt route.
type PIIDecryptResponse struct {
	Decrypted int      `json:"decrypted,omitempty"`
	Erased    int      `json:"erased,omitempty"`
	Text      string   `json:"text,omitempty"`
	Texts     []string `json:"texts,omitempty"`
}

// PIIErasedKeysetsResponse is the response of GET /v1/calls/{call_id}/pii/erased-keysets.
type PIIErasedKeysetsResponse struct {
	ErasedKeysetIDs []string `json:"erased_keyset_ids,omitempty"`
}

// PIIExemptionRequest is the body of POST and DELETE /v1/pii/exemptions.
type PIIExemptionRequest struct {
	ExemptType     PIIExemptType `json:"exempt_type,omitempty"`
	IdentifierHash string        `json:"identifier_hash,omitempty"`
}

// PIIExemption is one exemption; the value itself is never stored.
type PIIExemption struct {
	CreatedAt    string        `json:"created_at,omitempty"`
	ExemptType   PIIExemptType `json:"exempt_type,omitempty"`
	LookupPrefix string        `json:"lookup_prefix,omitempty"`
}

// PIIExemptionsResponse is the response of GET /v1/pii/exemptions.
type PIIExemptionsResponse struct {
	AutoExemptDIDs int            `json:"auto_exempt_dids,omitempty"`
	Exemptions     []PIIExemption `json:"exemptions,omitempty"`
}

// PIISettingsHistoryEntry is one entry of the encryption-switch history.
type PIISettingsHistoryEntry struct {
	ActorKeyID             string `json:"actor_key_id,omitempty"`
	ChangedAt              string `json:"changed_at,omitempty"`
	ConfigurationsResealed int    `json:"configurations_resealed,omitempty"`
	EncryptionEnabled      bool   `json:"encryption_enabled"`
	ImpersonatorID         string `json:"impersonator_id,omitempty"`
	ImpersonatorOrgID      string `json:"impersonator_org_id,omitempty"`
}

// PIISettingsPutRequest is the body of PUT /v1/pii/settings.
type PIISettingsPutRequest struct {
	// EncryptionEnabled true encrypts PII on new calls, false redacts it.
	EncryptionEnabled bool `json:"encryption_enabled"`
}

// PIISettingsResponse is the response of GET /v1/pii/settings.
type PIISettingsResponse struct {
	ChangedAt         string                    `json:"changed_at,omitempty"`
	ChangedBy         string                    `json:"changed_by,omitempty"`
	EncryptionEnabled bool                      `json:"encryption_enabled"`
	History           []PIISettingsHistoryEntry `json:"history,omitempty"`
}

// PIISettingsPutResponse is the response of PUT /v1/pii/settings.
type PIISettingsPutResponse struct {
	ChangedAt              string                    `json:"changed_at,omitempty"`
	ChangedBy              string                    `json:"changed_by,omitempty"`
	ConfigurationsResealed int                       `json:"configurations_resealed,omitempty"`
	EncryptionEnabled      bool                      `json:"encryption_enabled"`
	History                []PIISettingsHistoryEntry `json:"history,omitempty"`
}
