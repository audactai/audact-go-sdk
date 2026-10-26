package audact

// HandoverReason is why a warm handover was triggered.
type HandoverReason string

// Handover reasons accepted by the API.
const (
	HandoverReasonAgentRequested    HandoverReason = "agent_requested"
	HandoverReasonEscalation        HandoverReason = "escalation"
	HandoverReasonOperatorRequested HandoverReason = "operator_requested"
	HandoverReasonCallerRequested   HandoverReason = "caller_requested"
)

// WhisperMode controls how the human is briefed before the caller is bridged in.
type WhisperMode string

// Whisper modes accepted by the API.
const (
	WhisperModeSpoken WhisperMode = "spoken"
	WhisperModePortal WhisperMode = "portal"
	WhisperModeBoth   WhisperMode = "both"
)

// TransferOption is the SIP transfer type.
type TransferOption string

// Transfer options accepted by the API.
const (
	TransferOptionBlind TransferOption = "blind"
	TransferOptionWarm  TransferOption = "warm"
)

// HandoverRequest is the body of Calls.Handover.
type HandoverRequest struct {
	CallerReason   string         `json:"caller_reason,omitempty"`
	Destination    string         `json:"destination,omitempty"`
	HandoverReason HandoverReason `json:"handover_reason,omitempty"`
	WhisperMode    WhisperMode    `json:"whisper_mode,omitempty"`
}

// HandoverResponse is the response of Calls.Handover.
type HandoverResponse struct {
	CallID     string `json:"call_id,omitempty"`
	HandoverID string `json:"handover_id,omitempty"`
	State      string `json:"state,omitempty"`
}

// HandoverStateResponse is the response of Calls.GetHandover.
type HandoverStateResponse struct {
	CallID          string `json:"call_id,omitempty"`
	CompletedAt     string `json:"completed_at,omitempty"`
	DeliveryStatus  string `json:"delivery_status,omitempty"`
	DestinationHash string `json:"destination_hash,omitempty"`
	HandoverID      string `json:"handover_id,omitempty"`
	HandoverReason  string `json:"handover_reason,omitempty"`
	InitiatedAt     string `json:"initiated_at,omitempty"`
	State           string `json:"state,omitempty"`
	SummarySource   string `json:"summary_source,omitempty"`
	WhisperMode     string `json:"whisper_mode,omitempty"`
}

// TransferRequest is the body of Calls.Transfer.
type TransferRequest struct {
	Destination    string         `json:"destination,omitempty"`
	TransferOption TransferOption `json:"transfer_option,omitempty"`
}

// TransferResponse is the response of Calls.Transfer.
type TransferResponse struct {
	CallID         string `json:"call_id,omitempty"`
	Status         string `json:"status,omitempty"`
	TransferID     string `json:"transfer_id,omitempty"`
	TransferOption string `json:"transfer_option,omitempty"`
}

// TransferSessionResponse is the response of Calls.GetTransfer.
type TransferSessionResponse struct {
	AgentID         string `json:"agent_id,omitempty"`
	CallID          string `json:"call_id,omitempty"`
	CompletedAt     string `json:"completed_at,omitempty"`
	DestinationHash string `json:"destination_hash,omitempty"`
	InitiatedAt     string `json:"initiated_at,omitempty"`
	Status          string `json:"status,omitempty"`
	TransferID      string `json:"transfer_id,omitempty"`
	TransferOption  string `json:"transfer_option,omitempty"`
}

// TransferStateResponse is the response of the transfer bridge, cancel and end calls.
type TransferStateResponse struct {
	Status     string `json:"status,omitempty"`
	TransferID string `json:"transfer_id,omitempty"`
}

// ListPageParams are the optional limit and cursor query parameters of paged SIP lists.
type ListPageParams struct {
	Cursor string
	Limit  int
}

// RouteSummary is one inbound dispatch rule.
type RouteSummary struct {
	AgentID        string `json:"agent_id,omitempty"`
	DID            string `json:"did,omitempty"`
	DispatchRuleID string `json:"dispatch_rule_id,omitempty"`
	RoomPrefix     string `json:"room_prefix,omitempty"`
}

// ListRoutesResponse is the response of Calls.ListRoutes.
type ListRoutesResponse struct {
	HasMore    bool           `json:"has_more,omitempty"`
	NextCursor string         `json:"next_cursor,omitempty"`
	Routes     []RouteSummary `json:"routes,omitempty"`
}

// RouteInboundRequest is the body of Calls.CreateRoute.
type RouteInboundRequest struct {
	AgentID string `json:"agent_id,omitempty"`
	DID     string `json:"did,omitempty"`
}

// RouteInboundResponse is the response of Calls.CreateRoute.
type RouteInboundResponse struct {
	DispatchRuleID string `json:"dispatch_rule_id,omitempty"`
	Status         string `json:"status,omitempty"`
}

// InboundTrunkSummary is one inbound SIP trunk.
type InboundTrunkSummary struct {
	DID        string `json:"did,omitempty"`
	Provider   string `json:"provider,omitempty"`
	SIPTrunkID string `json:"sip_trunk_id,omitempty"`
}

// OutboundTrunkSummary is one outbound SIP trunk.
type OutboundTrunkSummary struct {
	Address    string `json:"address,omitempty"`
	Provider   string `json:"provider,omitempty"`
	SIPTrunkID string `json:"sip_trunk_id,omitempty"`
}

// ListTrunksResponse is the response of Calls.ListTrunks.
type ListTrunksResponse struct {
	HasMore        bool                   `json:"has_more,omitempty"`
	InboundTrunks  []InboundTrunkSummary  `json:"inbound_trunks,omitempty"`
	NextCursor     string                 `json:"next_cursor,omitempty"`
	OutboundTrunks []OutboundTrunkSummary `json:"outbound_trunks,omitempty"`
}

// RegisterInboundTrunkRequest is the body of Calls.CreateInboundTrunk.
type RegisterInboundTrunkRequest struct {
	AllowedAddresses []string `json:"allowed_addresses,omitempty"`
	AuthPassARN      string   `json:"auth_pass_arn,omitempty"`
	AuthPassword     string   `json:"auth_password,omitempty"`
	AuthUserARN      string   `json:"auth_user_arn,omitempty"`
	AuthUsername     string   `json:"auth_username,omitempty"`
	InboundNumbers   []string `json:"inbound_numbers,omitempty"`
	Provider         string   `json:"provider,omitempty"`
}

// RegisterInboundTrunkResponse is the response of Calls.CreateInboundTrunk.
type RegisterInboundTrunkResponse struct {
	InboundNumbers []string `json:"inbound_numbers,omitempty"`
	Provider       string   `json:"provider,omitempty"`
	SIPTrunkID     string   `json:"sip_trunk_id,omitempty"`
	Status         string   `json:"status,omitempty"`
}

// RegisterOutboundTrunkRequest is the body of Calls.CreateOutboundTrunk.
type RegisterOutboundTrunkRequest struct {
	Address     string   `json:"address,omitempty"`
	AuthPassARN string   `json:"auth_pass_arn,omitempty"`
	AuthUserARN string   `json:"auth_user_arn,omitempty"`
	Numbers     []string `json:"numbers,omitempty"`
	Provider    string   `json:"provider,omitempty"`
	Transport   string   `json:"transport,omitempty"`
}

// RegisterOutboundTrunkResponse is the response of Calls.CreateOutboundTrunk.
type RegisterOutboundTrunkResponse struct {
	Address    string `json:"address,omitempty"`
	Provider   string `json:"provider,omitempty"`
	SIPTrunkID string `json:"sip_trunk_id,omitempty"`
	Status     string `json:"status,omitempty"`
}

// AdmissionURLResponse is the response of Calls.GetAdmissionURL and Calls.SetAdmissionURL.
type AdmissionURLResponse struct {
	AdmissionURL string `json:"admission_url,omitempty"`
}

// SetAdmissionURLRequest is the body of Calls.SetAdmissionURL.
type SetAdmissionURLRequest struct {
	AdmissionURL string `json:"admission_url"`
}

// WSTicketResponse is the response of Calls.CreateWSTicket.
type WSTicketResponse struct {
	Ticket string `json:"ticket,omitempty"`
}
