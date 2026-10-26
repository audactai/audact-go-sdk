package audact

// Verdict is a compliance verdict.
type Verdict string

// Verdict values.
const (
	VerdictAllow Verdict = "ALLOW"
	VerdictBlock Verdict = "BLOCK"
)

// ValidateCallState is the call state sent with a validation request.
type ValidateCallState struct {
	DisclosureDelivered bool `json:"disclosure_delivered,omitempty"`
}

// ValidateTurn is a previous conversation turn.
type ValidateTurn struct {
	Content string `json:"content,omitempty"`
	Turn    int    `json:"turn,omitempty"`
	Verdict string `json:"verdict,omitempty"`
}

// ValidateContext carries conversation history and call state.
type ValidateContext struct {
	CallState     *ValidateCallState `json:"call_state,omitempty"`
	PreviousTurns []ValidateTurn     `json:"previous_turns,omitempty"`
}

// ValidateRequest is the body for validating a turn.
type ValidateRequest struct {
	AgentID        string           `json:"agent_id,omitempty"`
	CallID         string           `json:"call_id,omitempty"`
	ConfigID       string           `json:"config_id,omitempty"`
	ConfigSealHash string           `json:"config_seal_hash,omitempty"`
	Content        string           `json:"content,omitempty"`
	Context        *ValidateContext `json:"context,omitempty"`
	StreamSegment  bool             `json:"stream_segment,omitempty"`
	TenantID       string           `json:"tenant_id,omitempty"`
	Turn           int              `json:"turn,omitempty"`
}

// ValidateCheckResult is the outcome of one validation check.
type ValidateCheckResult struct {
	CheckID                string   `json:"check_id,omitempty"`
	DetectedLanguage       string   `json:"detected_language,omitempty"`
	DurationMs             float64  `json:"duration_ms,omitempty"`
	Layers                 []string `json:"layers,omitempty"`
	Mode                   string   `json:"mode,omitempty"`
	ModifiedContent        string   `json:"modified_content,omitempty"`
	NegationLexiconVersion string   `json:"negation_lexicon_version,omitempty"`
	NotEvaluated           bool     `json:"not_evaluated,omitempty"`
	NotRun                 bool     `json:"not_run,omitempty"`
	PromptFingerprint      string   `json:"prompt_fingerprint,omitempty"`
	PromptScope            string   `json:"prompt_scope,omitempty"`
	Reason                 string   `json:"reason,omitempty"`
	ShadowDetection        bool     `json:"shadow_detection,omitempty"`
	ShadowFindings         int      `json:"shadow_findings,omitempty"`
	ShadowVerdict          string   `json:"shadow_verdict,omitempty"`
	Verdict                string   `json:"verdict,omitempty"`
}

// ValidateResponse is the result of validating a turn.
type ValidateResponse struct {
	CallID   string                `json:"call_id,omitempty"`
	Checks   []ValidateCheckResult `json:"checks,omitempty"`
	Content  string                `json:"content,omitempty"`
	FailOpen bool                  `json:"fail_open,omitempty"`
	Turn     int                   `json:"turn,omitempty"`
	Verdict  string                `json:"verdict,omitempty"`
}

// FinalizeScoreRequest is the body for finalizing a call score.
type FinalizeScoreRequest struct {
	AgentID string `json:"agent_id,omitempty"`
}

// FinalizeScoreResponse is the result of finalizing a call score. CompositeScore
// is nil when the call had no recorded turns (status "skipped").
type FinalizeScoreResponse struct {
	CompositeScore      *float64 `json:"composite_score,omitempty"`
	FailOpenCount       int      `json:"fail_open_count,omitempty"`
	Status              string   `json:"status,omitempty"`
	TotalTurns          int      `json:"total_turns,omitempty"`
	TurnsWithViolations int      `json:"turns_with_violations,omitempty"`
}

// ComponentStatus is the health of one service component.
type ComponentStatus struct {
	Error   string `json:"error,omitempty"`
	Latency string `json:"latency,omitempty"`
	Status  string `json:"status,omitempty"`
}

// HealthResponse is a service health report.
type HealthResponse struct {
	Components map[string]ComponentStatus `json:"components,omitempty"`
	Status     string                     `json:"status,omitempty"`
	Timestamp  string                     `json:"timestamp,omitempty"`
}
