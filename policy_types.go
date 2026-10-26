package audact

// EvaluationContext holds the call-level facts referenced by policy rules.
type EvaluationContext struct {
	CallTime                 string `json:"call_time,omitempty"`
	CallWithinPermittedHours *bool  `json:"call_within_permitted_hours,omitempty"`
	ConsentStatus            string `json:"consent_status,omitempty"`
	ContactType              string `json:"contact_type,omitempty"`
	Direction                string `json:"direction,omitempty"`
	DNCRegistered            *bool  `json:"dnc_registered,omitempty"`
	ExistingRelationship     *bool  `json:"existing_relationship,omitempty"`
}

// EvaluateRequest is the body for evaluating a policy.
type EvaluateRequest struct {
	AgentID  string            `json:"agent_id,omitempty"`
	CallID   string            `json:"call_id,omitempty"`
	Context  EvaluationContext `json:"context"`
	Country  string            `json:"country"`
	Language string            `json:"language,omitempty"`
	TenantID string            `json:"tenant_id,omitempty"`
}

// EvaluateResponse is the result of a policy evaluation.
type EvaluateResponse struct {
	ApplicableRules    []string `json:"applicable_rules,omitempty"`
	DisclosureLanguage string   `json:"disclosure_language,omitempty"`
	DisclosureText     string   `json:"disclosure_text,omitempty"`
	NextPermittedAt    string   `json:"next_permitted_at,omitempty"`
	PolicyVersion      string   `json:"policy_version,omitempty"`
	Reasons            []string `json:"reasons,omitempty"`
	RequireDisclosure  bool     `json:"require_disclosure,omitempty"`
	Verdict            Verdict  `json:"verdict,omitempty"`
}

// CheckCallingHoursParams holds the query parameters for a calling-hours check.
type CheckCallingHoursParams struct {
	// Country is the jurisdiction code (required).
	Country string
	// CallTime is an optional RFC 3339 time to check; the server uses now when empty.
	CallTime string
}

// CallingHoursCheckResponse is the result of a calling-hours check.
type CallingHoursCheckResponse struct {
	Allowed         bool   `json:"allowed,omitempty"`
	NextPermittedAt string `json:"next_permitted_at,omitempty"`
	PolicyVersion   string `json:"policy_version,omitempty"`
	Reason          string `json:"reason,omitempty"`
}

// PolicyDNCAdapter describes a do-not-call registry adapter.
type PolicyDNCAdapter struct {
	Description string `json:"description,omitempty"`
	ID          string `json:"id,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Type        string `json:"type,omitempty"`
}

// PolicyHoliday is a holiday in a policy file.
type PolicyHoliday struct {
	Date string `json:"date,omitempty"`
	Name string `json:"name,omitempty"`
}

// PolicyJurisdictionCascade describes a parent jurisdiction and its children.
type PolicyJurisdictionCascade struct {
	Children []string `json:"children,omitempty"`
	Parent   string   `json:"parent,omitempty"`
}

// PolicyQuietWindow is a quiet-hours window.
type PolicyQuietWindow struct {
	Days     []string `json:"days,omitempty"`
	End      string   `json:"end,omitempty"`
	Label    string   `json:"label,omitempty"`
	Start    string   `json:"start,omitempty"`
	Timezone string   `json:"timezone,omitempty"`
}

// PolicyRuleCondition is a condition of a policy use case.
type PolicyRuleCondition struct {
	Field    string `json:"field,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}

// PolicyUseCase is a use case in a policy file.
type PolicyUseCase struct {
	BlockReasons   []string              `json:"block_reasons,omitempty"`
	Conditions     []PolicyRuleCondition `json:"conditions,omitempty"`
	DefaultVerdict string                `json:"default_verdict,omitempty"`
	Description    string                `json:"description,omitempty"`
	ID             string                `json:"id,omitempty"`
}

// PolicyFileResponse is a country policy document.
type PolicyFileResponse struct {
	Country             string                      `json:"country,omitempty"`
	DNCAdapters         []PolicyDNCAdapter          `json:"dnc_adapters,omitempty"`
	Holidays            []PolicyHoliday             `json:"holidays,omitempty"`
	JurisdictionCascade []PolicyJurisdictionCascade `json:"jurisdiction_cascade,omitempty"`
	QuietHours          []PolicyQuietWindow         `json:"quiet_hours,omitempty"`
	UpdatedAt           string                      `json:"updated_at,omitempty"`
	UseCases            []PolicyUseCase             `json:"use_cases,omitempty"`
	Version             string                      `json:"version,omitempty"`
}
