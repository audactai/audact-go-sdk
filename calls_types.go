package audact

import "encoding/json"

// CallRecord is one row of the call list.
type CallRecord struct {
	CallID    string `json:"call_id,omitempty"`
	Channel   string `json:"channel,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Direction string `json:"direction,omitempty"`
	EndedAt   string `json:"ended_at,omitempty"`
	Metadata  string `json:"metadata,omitempty"`
	StartedAt string `json:"started_at,omitempty"`
	Status    string `json:"status,omitempty"`
	TenantID  string `json:"tenant_id,omitempty"`
}

// ListCallsParams are the optional query parameters of Calls.List.
type ListCallsParams struct {
	DateFrom string
	DateTo   string
	Status   string
	Verdict  string
	// Shadow filters on shadow-mode calls when non-nil.
	Shadow   *bool
	Country  string
	ConfigID string
	Cursor   string
	Limit    int
}

// ListCallsResponse is the response of Calls.List.
type ListCallsResponse struct {
	Calls      []CallRecord `json:"calls,omitempty"`
	HasMore    bool         `json:"has_more,omitempty"`
	NextCursor string       `json:"next_cursor,omitempty"`
}

// CallComplianceScore is the finalized compliance score of a call.
type CallComplianceScore struct {
	Composite           float64 `json:"composite,omitempty"`
	FailOpenCount       int     `json:"fail_open_count,omitempty"`
	TotalTurns          int     `json:"total_turns,omitempty"`
	TurnsWithViolations int     `json:"turns_with_violations,omitempty"`
}

// TranscriptTurn is one turn of a call transcript.
type TranscriptTurn struct {
	Content string `json:"content,omitempty"`
	Role    string `json:"role,omitempty"`
}

// CallDetail is the full detail of one call.
type CallDetail struct {
	AgentID            string               `json:"agent_id,omitempty"`
	CallID             string               `json:"call_id,omitempty"`
	ComplianceScore    *CallComplianceScore `json:"compliance_score,omitempty"`
	Connected          *bool                `json:"connected,omitempty"`
	Direction          string               `json:"direction,omitempty"`
	DurationMs         *int                 `json:"duration_ms,omitempty"`
	EndedAt            string               `json:"ended_at,omitempty"`
	EvidenceLink       string               `json:"evidence_link,omitempty"`
	Metadata           string               `json:"metadata,omitempty"`
	PIIMode            string               `json:"pii_mode,omitempty"`
	RecordingAvailable bool                 `json:"recording_available,omitempty"`
	RecordingRef       string               `json:"recording_ref,omitempty"`
	SIPResponseCode    *int                 `json:"sip_response_code,omitempty"`
	StartedAt          string               `json:"started_at,omitempty"`
	Status             string               `json:"status,omitempty"`
	TenantID           string               `json:"tenant_id,omitempty"`
	TerminationClass   string               `json:"termination_class,omitempty"`
	Transcript         []TranscriptTurn     `json:"transcript,omitempty"`
	TranscriptStatus   string               `json:"transcript_status,omitempty"`
}

// CallEndSummary is the end-of-call and voicemail detection context.
type CallEndSummary struct {
	DetectedAt                  string  `json:"detected_at,omitempty"`
	DisclosureSkipReason        string  `json:"disclosure_skip_reason,omitempty"`
	DisclosureSkippedVoicemail  bool    `json:"disclosure_skipped_voicemail,omitempty"`
	TerminationConfidence       float64 `json:"termination_confidence,omitempty"`
	TerminationMatchedPhraseKey string  `json:"termination_matched_phrase_key,omitempty"`
	TerminationReason           string  `json:"termination_reason,omitempty"`
	VoicemailConfidence         float64 `json:"voicemail_confidence,omitempty"`
	VoicemailDetected           bool    `json:"voicemail_detected,omitempty"`
	VoicemailMatchedPhraseKey   string  `json:"voicemail_matched_phrase_key,omitempty"`
	VoicemailReason             string  `json:"voicemail_reason,omitempty"`
}

// HandoverSummary is the latest warm-handover context of a call.
type HandoverSummary struct {
	CompletedAt     string `json:"completed_at,omitempty"`
	DestinationHash string `json:"destination_hash,omitempty"`
	HandoverID      string `json:"handover_id,omitempty"`
	HandoverReason  string `json:"handover_reason,omitempty"`
	Outcome         string `json:"outcome,omitempty"`
	SummarySource   string `json:"summary_source,omitempty"`
	SummaryText     string `json:"summary_text,omitempty"`
	WhisperMode     bool   `json:"whisper_mode,omitempty"`
}

// CallSummary is the composite view of a call.
type CallSummary struct {
	AllowCount     int              `json:"allow_count,omitempty"`
	AnalysisStatus string           `json:"analysis_status,omitempty"`
	BlockCount     int              `json:"block_count,omitempty"`
	CallEnd        *CallEndSummary  `json:"call_end,omitempty"`
	CallID         string           `json:"call_id,omitempty"`
	Disposition    string           `json:"disposition,omitempty"`
	EndedAt        string           `json:"ended_at,omitempty"`
	EventCount     int              `json:"event_count,omitempty"`
	Extraction     json.RawMessage  `json:"extraction,omitempty"`
	Handover       *HandoverSummary `json:"handover,omitempty"`
	KeyFacts       []string         `json:"key_facts,omitempty"`
	ModifyCount    int              `json:"modify_count,omitempty"`
	Sentiment      string           `json:"sentiment,omitempty"`
	StartedAt      string           `json:"started_at,omitempty"`
	Status         string           `json:"status,omitempty"`
	SummaryStatus  string           `json:"summary_status,omitempty"`
	TenantID       string           `json:"tenant_id,omitempty"`
}

// RecordingResponse is a short-lived signed URL for a call recording.
type RecordingResponse struct {
	CallID    string `json:"call_id,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
	URL       string `json:"url,omitempty"`
}

// IdentifierHash is a client-side hashed identifier linked to a call.
type IdentifierHash struct {
	Hash string `json:"hash,omitempty"`
	Type string `json:"type,omitempty"`
}

// OutboundCallRequest is the body of Calls.CreateOutbound.
//
// FromNumber and DestinationNumber must be strict E.164: a leading '+', the
// country code, then the subscriber number — digits only, 8 to 16 characters
// in total, no spaces, dashes, dots or parentheses. Numbers are not
// normalised; "+44 20 7946 0958" is rejected with 400 INVALID_REQUEST.
// Example: "+442079460958". FromNumber must also be one of your own numbers
// (Numbers.List, GET /v1/numbers).
type OutboundCallRequest struct {
	AgentID            string           `json:"agent_id,omitempty"`
	DestinationNumber  string           `json:"destination_number,omitempty"`
	DynamicVariables   map[string]any   `json:"dynamic_variables,omitempty"`
	FromNumber         string           `json:"from_number,omitempty"`
	Identifiers        []IdentifierHash `json:"identifiers,omitempty"`
	MaxDurationSeconds int              `json:"max_duration_seconds,omitempty"`
	Metadata           map[string]any   `json:"metadata,omitempty"`
	RequestedOutcome   string           `json:"requested_outcome,omitempty"`
}

// OutboundCallResponse is the response of Calls.CreateOutbound.
type OutboundCallResponse struct {
	CallID        string `json:"call_id,omitempty"`
	ParticipantID string `json:"participant_id,omitempty"`
	RoomName      string `json:"room_name,omitempty"`
	SIPCallID     string `json:"sip_call_id,omitempty"`
}

// InboundCallRequest is the body of Calls.CreateInbound.
type InboundCallRequest struct {
	CalledDID        string           `json:"called_did,omitempty"`
	CallerNumber     string           `json:"caller_number,omitempty"`
	DynamicVariables map[string]any   `json:"dynamic_variables,omitempty"`
	Identifiers      []IdentifierHash `json:"identifiers,omitempty"`
	Metadata         map[string]any   `json:"metadata,omitempty"`
}

// InboundCallResponse is the response of Calls.CreateInbound.
type InboundCallResponse struct {
	CallID   string `json:"call_id,omitempty"`
	RoomName string `json:"room_name,omitempty"`
	TrunkID  string `json:"trunk_id,omitempty"`
}

// BrowserCallRequest is the body of Calls.CreateBrowser.
type BrowserCallRequest struct {
	AgentID          string           `json:"agent_id,omitempty"`
	DynamicVariables map[string]any   `json:"dynamic_variables,omitempty"`
	Identifiers      []IdentifierHash `json:"identifiers,omitempty"`
	Metadata         map[string]any   `json:"metadata,omitempty"`
}

// BrowserCallResponse is the response of Calls.CreateBrowser.
type BrowserCallResponse struct {
	CallID       string `json:"call_id,omitempty"`
	LiveKitToken string `json:"livekit_token,omitempty"`
	LiveKitURL   string `json:"livekit_url,omitempty"`
	RoomName     string `json:"room_name,omitempty"`
}

// HangupResponse is the response of Calls.Hangup.
type HangupResponse struct {
	CallID   string `json:"call_id,omitempty"`
	RoomName string `json:"room_name,omitempty"`
	Status   string `json:"status,omitempty"`
}

// CallStatusResponse is the response of Calls.Status.
type CallStatusResponse struct {
	CallID    string `json:"call_id,omitempty"`
	RoomName  string `json:"room_name,omitempty"`
	StartedAt string `json:"started_at,omitempty"`
	Status    string `json:"status,omitempty"`
}
