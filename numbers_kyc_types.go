package audact

// KYCDocumentSetStatus is the review state of a KYC document set.
type KYCDocumentSetStatus string

// KYC document set statuses.
const (
	KYCDocumentSetStatusDraft     KYCDocumentSetStatus = "draft"
	KYCDocumentSetStatusSubmitted KYCDocumentSetStatus = "submitted"
	KYCDocumentSetStatusApproved  KYCDocumentSetStatus = "approved"
	KYCDocumentSetStatusRejected  KYCDocumentSetStatus = "rejected"
	KYCDocumentSetStatusExpired   KYCDocumentSetStatus = "expired"
)

// KYCRequirementsParams selects the regulatory requirements to look up.
type KYCRequirementsParams struct {
	Country    string
	NumberType string
}

// RegulatoryRequirement is one regulatory requirement for a country and number type.
type RegulatoryRequirement struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Required        bool   `json:"required,omitempty"`
	RequirementType string `json:"requirement_type,omitempty"`
}

// ListKYCRequirementsResponse lists regulatory requirements.
type ListKYCRequirementsResponse struct {
	Requirements []RegulatoryRequirement `json:"requirements,omitempty"`
}

// KYCDocument is the metadata of an uploaded KYC document.
type KYCDocument struct {
	ContentType      string `json:"content_type,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	DocType          string `json:"doc_type,omitempty"`
	Filename         string `json:"filename,omitempty"`
	ID               string `json:"id,omitempty"`
	SHA256           string `json:"sha256,omitempty"`
	SizeBytes        int    `json:"size_bytes,omitempty"`
	TelnyxDocumentID string `json:"telnyx_document_id,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}

// ListKYCDocumentsResponse lists KYC documents.
type ListKYCDocumentsResponse struct {
	Documents []KYCDocument `json:"documents,omitempty"`
}

// KYCUploadOptions are optional settings for KYC document uploads.
type KYCUploadOptions struct {
	// IdempotencyKey, when set, is sent as the Idempotency-Key header.
	IdempotencyKey string
}

// CreateKYCDocumentSetItemRequest answers one regulatory requirement with
// either a document or a text value.
type CreateKYCDocumentSetItemRequest struct {
	DocumentID      string `json:"document_id,omitempty"`
	RequirementID   string `json:"requirement_id,omitempty"`
	RequirementType string `json:"requirement_type,omitempty"`
	TextValue       string `json:"text_value,omitempty"`
}

// CreateKYCDocumentSetRequest creates a KYC document set.
type CreateKYCDocumentSetRequest struct {
	CountryCode string                            `json:"country_code,omitempty"`
	Items       []CreateKYCDocumentSetItemRequest `json:"items,omitempty"`
	Name        string                            `json:"name,omitempty"`
	NumberType  string                            `json:"number_type,omitempty"`
}

// KYCDocumentSetItem is one answered requirement of a KYC document set.
type KYCDocumentSetItem struct {
	DocType          string `json:"doc_type,omitempty"`
	DocumentFilename string `json:"document_filename,omitempty"`
	DocumentID       string `json:"document_id,omitempty"`
	RequirementID    string `json:"requirement_id,omitempty"`
	RequirementType  string `json:"requirement_type,omitempty"`
	TextValue        string `json:"text_value,omitempty"`
}

// KYCDocumentSet is a KYC document set submitted for regulatory approval.
type KYCDocumentSet struct {
	Action                   string               `json:"action,omitempty"`
	CountryCode              string               `json:"country_code,omitempty"`
	CreatedAt                string               `json:"created_at,omitempty"`
	ID                       string               `json:"id,omitempty"`
	Items                    []KYCDocumentSetItem `json:"items,omitempty"`
	Name                     string               `json:"name,omitempty"`
	NumberType               string               `json:"number_type,omitempty"`
	Status                   KYCDocumentSetStatus `json:"status,omitempty"`
	StatusDetail             string               `json:"status_detail,omitempty"`
	TelnyxRequirementGroupID string               `json:"telnyx_requirement_group_id,omitempty"`
	UpdatedAt                string               `json:"updated_at,omitempty"`
}

// ListKYCDocumentSetsResponse lists KYC document sets.
type ListKYCDocumentSetsResponse struct {
	DocumentSets []KYCDocumentSet `json:"document_sets,omitempty"`
}

// RotateKYCDEKResponse reports the new document encryption key version.
type RotateKYCDEKResponse struct {
	NewVersion int `json:"new_version,omitempty"`
}
