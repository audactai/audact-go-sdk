package audact

import "encoding/json"

// ComplianceEventSubtype is the subtype of a compliance event.
type ComplianceEventSubtype string

// ComplianceEventSubtype values. The last three are webhook-only and always
// yield an empty page when used as a filter.
const (
	ComplianceEventValidationCompleted ComplianceEventSubtype = "validation_completed"
	ComplianceEventViolationDetected   ComplianceEventSubtype = "violation_detected"
	ComplianceEventFailOpen            ComplianceEventSubtype = "fail_open"
	ComplianceEventPolicyTriggered     ComplianceEventSubtype = "policy_triggered"
	ComplianceEventAuditLogged         ComplianceEventSubtype = "audit_logged"
	ComplianceEventReceiptIssued       ComplianceEventSubtype = "receipt_issued"
)

// ComplianceVerdict is the verdict of a compliance event.
type ComplianceVerdict string

// ComplianceVerdict values.
const (
	ComplianceVerdictAllow  ComplianceVerdict = "ALLOW"
	ComplianceVerdictModify ComplianceVerdict = "MODIFY"
	ComplianceVerdictBlock  ComplianceVerdict = "BLOCK"
)

// ComplianceSeverity is the severity band of a compliance event.
type ComplianceSeverity string

// ComplianceSeverity values.
const (
	ComplianceSeverityCritical ComplianceSeverity = "critical"
	ComplianceSeverityHigh     ComplianceSeverity = "high"
	ComplianceSeverityMedium   ComplianceSeverity = "medium"
	ComplianceSeverityLow      ComplianceSeverity = "low"
)

// ListComplianceEventsParams filters GET /v1/compliance/events. Zero values are omitted.
type ListComplianceEventsParams struct {
	EventSubtype ComplianceEventSubtype
	Verdict      ComplianceVerdict
	Severity     ComplianceSeverity
	CallID       string
	Cursor       string
	Limit        int
}

// ComplianceEvent is one compliance event.
type ComplianceEvent struct {
	AgentID         string                 `json:"agent_id,omitempty"`
	CallID          string                 `json:"call_id,omitempty"`
	CheckID         string                 `json:"check_id,omitempty"`
	ComplianceScore float64                `json:"compliance_score,omitempty"`
	Detail          string                 `json:"detail,omitempty"`
	EventSubtype    ComplianceEventSubtype `json:"event_subtype,omitempty"`
	Severity        ComplianceSeverity     `json:"severity,omitempty"`
	Timestamp       string                 `json:"timestamp,omitempty"`
	TurnSequence    int                    `json:"turn_sequence,omitempty"`
	Verdict         ComplianceVerdict      `json:"verdict,omitempty"`
}

// ListComplianceEventsResponse is one page of compliance events.
type ListComplianceEventsResponse struct {
	Events     []ComplianceEvent `json:"events,omitempty"`
	HasMore    bool              `json:"has_more,omitempty"`
	NextCursor string            `json:"next_cursor,omitempty"`
	Total      int               `json:"total,omitempty"`
}

// PostEventRequest is the body for writing an evidence event.
type PostEventRequest struct {
	CallID    string `json:"call_id,omitempty"`
	Data      string `json:"data,omitempty"`
	EventType string `json:"event_type,omitempty"`
	TenantID  string `json:"tenant_id,omitempty"`
}

// PostEventResponse is the hash receipt for a written evidence event.
type PostEventResponse struct {
	CallID       string `json:"call_id,omitempty"`
	EventID      string `json:"event_id,omitempty"`
	Hash         string `json:"hash,omitempty"`
	PreviousHash string `json:"previous_hash,omitempty"`
	Sequence     int    `json:"sequence,omitempty"`
	TenantID     string `json:"tenant_id,omitempty"`
	Timestamp    string `json:"timestamp,omitempty"`
}

// ChainParams are the query parameters of GET /v1/chains/{call_id}.
type ChainParams struct {
	AfterSequence int
	Limit         int
}

// GetChainResponse is one page of an evidence chain. Events are free-form objects.
type GetChainResponse struct {
	CallID     string            `json:"call_id,omitempty"`
	Events     []json.RawMessage `json:"events,omitempty"`
	HasMore    bool              `json:"has_more,omitempty"`
	NextCursor int               `json:"next_cursor,omitempty"`
}

// VerifyChainResponse is the result of verifying an evidence chain.
type VerifyChainResponse struct {
	BrokenAtSequence int    `json:"broken_at_sequence,omitempty"`
	CallID           string `json:"call_id,omitempty"`
	EventCount       int    `json:"event_count,omitempty"`
	Status           string `json:"status,omitempty"`
}

// EventComplianceReceipt is the compliance receipt attached to a reconstructed event.
type EventComplianceReceipt struct {
	ConfigSealHash string `json:"config_seal_hash,omitempty"`
	EventID        string `json:"event_id,omitempty"`
	IssuedAt       string `json:"issued_at,omitempty"`
	Verdict        string `json:"verdict,omitempty"`
}

// ReconstructedEvent is an evidence event in a reconstruction.
type ReconstructedEvent struct {
	CallID            string                  `json:"call_id,omitempty"`
	ComplianceReceipt *EventComplianceReceipt `json:"compliance_receipt,omitempty"`
	ConfigSealHash    string                  `json:"config_seal_hash,omitempty"`
	Data              string                  `json:"data,omitempty"`
	EventID           string                  `json:"event_id,omitempty"`
	EventType         string                  `json:"event_type,omitempty"`
	Hash              string                  `json:"hash,omitempty"`
	PreviousHash      string                  `json:"previous_hash,omitempty"`
	Sequence          int                     `json:"sequence,omitempty"`
	TenantID          string                  `json:"tenant_id,omitempty"`
	Timestamp         string                  `json:"timestamp,omitempty"`
}

// ResolvedConfigVersion is a tenant configuration version in effect during a call.
//
// The settings fields (BeginMessage through Voicemail) are the version's own
// settings, the same shape as on Configuration.
type ResolvedConfigVersion struct {
	AutoHandoverConfig          json.RawMessage  `json:"auto_handover_config,omitempty"`
	BeginMessage                string           `json:"begin_message,omitempty"`
	BlocklistPatterns           json.RawMessage  `json:"blocklist_patterns,omitempty"`
	ConfigID                    string           `json:"config_id,omitempty"`
	ConfigSealHash              string           `json:"config_seal_hash,omitempty"`
	ConsentConfig               json.RawMessage  `json:"consent_config,omitempty"`
	Countries                   []string         `json:"countries,omitempty"`
	CreatedAt                   string           `json:"created_at,omitempty"`
	DeletedAt                   string           `json:"deleted_at,omitempty"`
	DisclosureTexts             json.RawMessage  `json:"disclosure_texts,omitempty"`
	EndCall                     *EndCallConfig   `json:"end_call,omitempty"`
	FallbackResponse            json.RawMessage  `json:"fallback_response,omitempty"`
	LLMModel                    string           `json:"llm_model,omitempty"`
	LLMProvider                 string           `json:"llm_provider,omitempty"`
	Name                        string           `json:"name,omitempty"`
	OutputConstraint            json.RawMessage  `json:"output_constraint,omitempty"`
	PIIEncryption               bool             `json:"pii_encryption"`
	PricingBands                json.RawMessage  `json:"pricing_bands,omitempty"`
	PromptConfidentialityConfig json.RawMessage  `json:"prompt_confidentiality_config,omitempty"`
	TenantID                    string           `json:"tenant_id,omitempty"`
	UpdatedAt                   string           `json:"updated_at,omitempty"`
	ValidatorModes              json.RawMessage  `json:"validator_modes,omitempty"`
	Version                     int              `json:"version,omitempty"`
	Voicemail                   *VoicemailConfig `json:"voicemail,omitempty"`
}

// ReconstructResponse is a reconstructed evidence chain.
type ReconstructResponse struct {
	CallID         string                           `json:"call_id,omitempty"`
	ChainStatus    string                           `json:"chain_status,omitempty"`
	ConfigVersions map[string]ResolvedConfigVersion `json:"config_versions,omitempty"`
	EventCount     int                              `json:"event_count,omitempty"`
	Events         []ReconstructedEvent             `json:"events,omitempty"`
	PolicyVersion  string                           `json:"policy_version,omitempty"`
}

// ListViolationsParams filters GET /v1/violations. Zero values are omitted.
type ListViolationsParams struct {
	CallID   string
	DateFrom string
	DateTo   string
	Verdict  string
	Cursor   string
	Limit    int
}

// ViolationRecord is one violation evidence event.
type ViolationRecord struct {
	CallID    string `json:"call_id,omitempty"`
	Data      string `json:"data,omitempty"`
	EventID   string `json:"event_id,omitempty"`
	EventType string `json:"event_type,omitempty"`
	Hash      string `json:"hash,omitempty"`
	Sequence  int    `json:"sequence,omitempty"`
	TenantID  string `json:"tenant_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// GetViolationsResponse is one page of violations.
type GetViolationsResponse struct {
	HasMore    bool              `json:"has_more,omitempty"`
	NextCursor string            `json:"next_cursor,omitempty"`
	Violations []ViolationRecord `json:"violations,omitempty"`
}

// ReceiptVerdict is the verdict recorded on a compliance receipt.
type ReceiptVerdict string

// ReceiptVerdict values.
const (
	ReceiptVerdictAllow  ReceiptVerdict = "ALLOW"
	ReceiptVerdictBlock  ReceiptVerdict = "BLOCK"
	ReceiptVerdictModify ReceiptVerdict = "MODIFY"
)

// ComplianceReceipt is a compliance receipt carrying an HMAC integrity tag.
type ComplianceReceipt struct {
	AuthMechanism           string         `json:"auth_mechanism,omitempty"`
	CallID                  string         `json:"call_id,omitempty"`
	CandidateCount          int            `json:"candidate_count,omitempty"`
	CandidateSetHash        string         `json:"candidate_set_hash,omitempty"`
	ConfigSealHash          string         `json:"config_seal_hash,omitempty"`
	ConfigVersion           int            `json:"config_version,omitempty"`
	ConfirmationGrade       string         `json:"confirmation_grade,omitempty"`
	DeliveryAttestationHash string         `json:"delivery_attestation_hash,omitempty"`
	DeliveryStatus          string         `json:"delivery_status,omitempty"`
	EarlyTerminated         bool           `json:"early_terminated,omitempty"`
	EvidenceHash            string         `json:"evidence_hash,omitempty"`
	InputHash               string         `json:"input_hash,omitempty"`
	OutputHash              string         `json:"output_hash,omitempty"`
	PolicyVersion           string         `json:"policy_version,omitempty"`
	ReceiptFormatVersion    string         `json:"receipt_format_version,omitempty"`
	SelectedCandidateID     string         `json:"selected_candidate_id,omitempty"`
	Signature               string         `json:"signature,omitempty"`
	Timestamp               string         `json:"timestamp,omitempty"`
	Verdict                 ReceiptVerdict `json:"verdict,omitempty"`
}

// GetReceiptResponse is a receipt plus the result of verifying its signature.
type GetReceiptResponse struct {
	HMACValid     bool               `json:"hmac_valid,omitempty"`
	Receipt       *ComplianceReceipt `json:"receipt,omitempty"`
	Verifiability string             `json:"verifiability,omitempty"`
}

// TurnDeliveryStatus is the delivery confirmation status of one turn.
type TurnDeliveryStatus struct {
	// ConfirmationCount is how many confirmations exist for the turn. When it
	// is more than 1 the other fields describe the worst of them. Zero while
	// pending.
	ConfirmationCount int    `json:"confirmation_count,omitempty"`
	ConfirmationGrade string `json:"confirmation_grade,omitempty"`
	ConfirmedAt       string `json:"confirmed_at,omitempty"`
	ContentMatch      bool   `json:"content_match,omitempty"`
	LinkingHash       string `json:"linking_hash,omitempty"`
	Mechanism         string `json:"mechanism,omitempty"`
	MismatchSeverity  string `json:"mismatch_severity,omitempty"`
	MismatchType      string `json:"mismatch_type,omitempty"`
	Status            string `json:"status,omitempty"`
	TurnNumber        int    `json:"turn_number,omitempty"`
}

// ConfirmDeliveryRequest is the body for confirming delivery of a turn.
type ConfirmDeliveryRequest struct {
	ASRConfidence     float64  `json:"asr_confidence,omitempty"`
	ASRTranscription  string   `json:"asr_transcription,omitempty"`
	AudioHash         string   `json:"audio_hash,omitempty"`
	CallID            string   `json:"call_id,omitempty"`
	DeliveredSegments int      `json:"delivered_segments,omitempty"`
	DeliveredText     string   `json:"delivered_text,omitempty"`
	Interrupted       bool     `json:"interrupted,omitempty"`
	Mechanism         string   `json:"mechanism,omitempty"`
	PlatformTimestamp string   `json:"platform_timestamp,omitempty"`
	TurnNumber        int      `json:"turn_number,omitempty"`
	ValidatedSegments []string `json:"validated_segments,omitempty"`
}

// ConfirmDeliveryResponse is the result of processing a delivery confirmation.
type ConfirmDeliveryResponse struct {
	AttestationEventID string `json:"attestation_event_id,omitempty"`
	ConfirmationGrade  string `json:"confirmation_grade,omitempty"`
	ContentMatch       bool   `json:"content_match,omitempty"`
	LinkingHash        string `json:"linking_hash,omitempty"`
	MismatchSeverity   string `json:"mismatch_severity,omitempty"`
	MismatchType       string `json:"mismatch_type,omitempty"`
}

// CrossTurnFinding is a cross-turn inconsistency found by reconciliation.
type CrossTurnFinding struct {
	Confidence string `json:"confidence,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Summary    string `json:"summary,omitempty"`
	TurnA      int    `json:"turn_a,omitempty"`
	TurnB      int    `json:"turn_b,omitempty"`
}

// ReconciliationSegmentResult is the reconciliation outcome of one segment.
type ReconciliationSegmentResult struct {
	MismatchClassification string `json:"mismatch_classification,omitempty"`
	SequenceNumber         int    `json:"sequence_number,omitempty"`
	Status                 string `json:"status,omitempty"`
	TranscriptHash         string `json:"transcript_hash,omitempty"`
	ValidatedHash          string `json:"validated_hash,omitempty"`
}

// ReconcileTranscriptResponse is the result of reconciling a transcript.
type ReconcileTranscriptResponse struct {
	CallID                string                        `json:"call_id,omitempty"`
	CrossTurnFindingCount int                           `json:"cross_turn_finding_count,omitempty"`
	CrossTurnFindings     []CrossTurnFinding            `json:"cross_turn_findings,omitempty"`
	CrossTurnModel        string                        `json:"cross_turn_model,omitempty"`
	CrossTurnProvider     string                        `json:"cross_turn_provider,omitempty"`
	CrossTurnStatus       string                        `json:"cross_turn_status,omitempty"`
	MatchedSegments       int                           `json:"matched_segments,omitempty"`
	MismatchedSegments    int                           `json:"mismatched_segments,omitempty"`
	SegmentResults        []ReconciliationSegmentResult `json:"segment_results,omitempty"`
	TotalSegments         int                           `json:"total_segments,omitempty"`
	UnvalidatedSegments   int                           `json:"unvalidated_segments,omitempty"`
}

// Article50Params are the query parameters of the Article 50 report endpoints.
type Article50Params struct {
	Format string
}

// Article50Export is an Article 50 export in the format the server served:
// exactly one of Report (application/json) and PDF (application/pdf) is set.
type Article50Export struct {
	Report *Article50Report
	PDF    []byte
}

// Article50Report is the signed Article 50 evidence export. Its nested
// sections are kept as raw JSON for wire fidelity.
type Article50Report struct {
	CallID              string          `json:"call_id,omitempty"`
	Channel             string          `json:"channel,omitempty"`
	ConfigurationSeal   json.RawMessage `json:"configuration_seal,omitempty"`
	Direction           string          `json:"direction,omitempty"`
	Disclosure          json.RawMessage `json:"disclosure,omitempty"`
	EvidenceChain       json.RawMessage `json:"evidence_chain,omitempty"`
	GeneratedAt         string          `json:"generated_at,omitempty"`
	MerkleProof         json.RawMessage `json:"merkle_proof,omitempty"`
	Receipts            json.RawMessage `json:"receipts,omitempty"`
	ReportFormatVersion string          `json:"report_format_version,omitempty"`
	ReportType          string          `json:"report_type,omitempty"`
	Signature           json.RawMessage `json:"signature,omitempty"`
	TenantID            string          `json:"tenant_id,omitempty"`
	TSATimestamp        json.RawMessage `json:"tsa_timestamp,omitempty"`
}

// C2PAManifest is the summary of an issued C2PA content credential. The
// claim is signed ES256 with an X.509 chain; verify the downloaded manifest
// store with c2patool. Assertions are raw JSON.
type C2PAManifest struct {
	Assertions         []json.RawMessage      `json:"assertions,omitempty"`
	Asset              *C2PAAssetRef          `json:"asset,omitempty"`
	CallID             string                 `json:"call_id,omitempty"`
	Claim              json.RawMessage        `json:"claim,omitempty"`
	GeneratedAt        string                 `json:"generated_at,omitempty"`
	ManifestID         string                 `json:"manifest_id,omitempty"`
	ManifestStore      *C2PAManifestStoreRef  `json:"manifest_store,omitempty"`
	SignatureAlgorithm string                 `json:"signature_algorithm,omitempty"`
	SigningCredential  *C2PASigningCredential `json:"signing_credential,omitempty"`
}

// C2PAAssetRef is the asset a manifest is hard-bound to. Kind is
// "evidence_export" or "recording"; URL is the relative API path that
// DownloadC2PAAsset fetches.
type C2PAAssetRef struct {
	Kind      string `json:"kind,omitempty"`
	MediaType string `json:"media_type,omitempty"`
	SHA256    string `json:"sha256,omitempty"`
	URL       string `json:"url,omitempty"`
}

// C2PAManifestStoreRef is the stored sidecar manifest store
// (application/c2pa) that DownloadC2PAManifestStore fetches.
type C2PAManifestStoreRef struct {
	MediaType string `json:"media_type,omitempty"`
	SHA256    string `json:"sha256,omitempty"`
	URL       string `json:"url,omitempty"`
}

// C2PASigningCredential holds public facts about the certificate that signed
// the claim; never key material. Trust is "self_signed" or "ca_issued".
type C2PASigningCredential struct {
	FingerprintSHA256 string `json:"fingerprint_sha256,omitempty"`
	Issuer            string `json:"issuer,omitempty"`
	NotAfter          string `json:"not_after,omitempty"`
	NotBefore         string `json:"not_before,omitempty"`
	Serial            string `json:"serial,omitempty"`
	Subject           string `json:"subject,omitempty"`
	Trust             string `json:"trust,omitempty"`
}

// C2PAListSigner is the stored subset of the signing credential carried by a
// manifest list item.
type C2PAListSigner struct {
	FingerprintSHA256 string `json:"fingerprint_sha256,omitempty"`
	Trust             string `json:"trust,omitempty"`
}

// C2PAManifestListItem is one manifest of a call: a subset of C2PAManifest
// with no claim or assertions.
type C2PAManifestListItem struct {
	Asset              *C2PAAssetRef         `json:"asset,omitempty"`
	CallID             string                `json:"call_id,omitempty"`
	GeneratedAt        string                `json:"generated_at,omitempty"`
	ManifestID         string                `json:"manifest_id,omitempty"`
	ManifestStore      *C2PAManifestStoreRef `json:"manifest_store,omitempty"`
	SignatureAlgorithm string                `json:"signature_algorithm,omitempty"`
	SigningCredential  *C2PAListSigner       `json:"signing_credential,omitempty"`
}

// C2PAManifestListResponse is a page of a call's manifests, newest first.
// NextCursor is empty on the last page.
type C2PAManifestListResponse struct {
	Manifests  []C2PAManifestListItem `json:"manifests,omitempty"`
	NextCursor string                 `json:"next_cursor,omitempty"`
}

// ListC2PAManifestsParams are the query parameters of ListC2PAManifests.
type ListC2PAManifestsParams struct {
	Limit  int
	Cursor string
}
