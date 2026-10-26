package audact

// NumberOrigin says where a phone number came from.
type NumberOrigin string

// Number origins.
const (
	NumberOriginAudact NumberOrigin = "audact"
	NumberOriginBYO    NumberOrigin = "byo"
)

// NumberStatus is the lifecycle state of a phone number.
type NumberStatus string

// Number lifecycle states.
const (
	NumberStatusAvailable NumberStatus = "available"
	NumberStatusAssigned  NumberStatus = "assigned"
	NumberStatusSuspended NumberStatus = "suspended"
)

// ConnectionHealth is the verification health of a BYO connection.
type ConnectionHealth string

// Connection health values.
const (
	ConnectionHealthUnverified  ConnectionHealth = "unverified"
	ConnectionHealthHealthy     ConnectionHealth = "healthy"
	ConnectionHealthUnreachable ConnectionHealth = "unreachable"
	ConnectionHealthAuthFailed  ConnectionHealth = "auth_failed"
)

// NumberProvider is the SIP provider behind a BYO connection.
type NumberProvider string

// Number providers.
const (
	NumberProviderTwilio     NumberProvider = "twilio"
	NumberProviderTelnyx     NumberProvider = "telnyx"
	NumberProviderGenericSIP NumberProvider = "generic_sip"
)

// NumberTrunkProvider is the carrier for an Audact-provisioned number
// (CreateNumberRequest.TrunkProvider); the spec only allows "telnyx".
type NumberTrunkProvider string

// Number trunk providers accepted by CreateNumberRequest.
const (
	NumberTrunkProviderTelnyx NumberTrunkProvider = "telnyx"
)

// TrunkProvider is the trunk provider of a number.
type TrunkProvider string

// Trunk providers.
const (
	TrunkProviderTelnyx TrunkProvider = "telnyx"
	TrunkProviderTwilio TrunkProvider = "twilio"
)

// NumberOrderStatus is the state of a number order.
type NumberOrderStatus string

// Number order states.
const (
	NumberOrderStatusPending           NumberOrderStatus = "pending"
	NumberOrderStatusRequirementReview NumberOrderStatus = "requirement_review"
	NumberOrderStatusOrdering          NumberOrderStatus = "ordering"
	NumberOrderStatusSuccess           NumberOrderStatus = "success"
	NumberOrderStatusFailure           NumberOrderStatus = "failure"
	NumberOrderStatusCancelled         NumberOrderStatus = "cancelled"
)

// SIPTransport is the SIP transport protocol of a BYO connection.
type SIPTransport string

// SIP transports.
const (
	SIPTransportUDP  SIPTransport = "udp"
	SIPTransportTCP  SIPTransport = "tcp"
	SIPTransportTLS  SIPTransport = "tls"
	SIPTransportAuto SIPTransport = "auto"
)

// BYOConnectionRequest describes a BYO SIP-trunk connection.
type BYOConnectionRequest struct {
	InboundAllowedIPs []string       `json:"inbound_allowed_ips,omitempty"`
	Name              string         `json:"name,omitempty"`
	Provider          NumberProvider `json:"provider,omitempty"`
	SIPHost           string         `json:"sip_host,omitempty"`
	SIPPassword       string         `json:"sip_password,omitempty"`
	SIPTransport      SIPTransport   `json:"sip_transport,omitempty"`
	SIPUsername       string         `json:"sip_username,omitempty"`
}

// PhoneNumber is a DID phone number in the tenant's pool.
type PhoneNumber struct {
	AssignedAgentID  string           `json:"assigned_agent_id,omitempty"`
	BYOConnectionID  string           `json:"byo_connection_id,omitempty"`
	ConnectionHealth ConnectionHealth `json:"connection_health,omitempty"`
	CountryCode      string           `json:"country_code,omitempty"`
	CreatedAt        string           `json:"created_at,omitempty"`
	DID              string           `json:"did,omitempty"`
	ID               string           `json:"id,omitempty"`
	Origin           NumberOrigin     `json:"origin,omitempty"`
	Provider         NumberProvider   `json:"provider,omitempty"`
	Status           NumberStatus     `json:"status,omitempty"`
	TrunkID          string           `json:"trunk_id,omitempty"`
	TrunkProvider    TrunkProvider    `json:"trunk_provider,omitempty"`
	UpdatedAt        string           `json:"updated_at,omitempty"`
}

// ListNumbersParams are the query parameters of Numbers.List.
type ListNumbersParams struct {
	Country string
	Status  NumberStatus
	Limit   int
	Cursor  string
}

// ListNumbersResponse is the response of Numbers.List.
type ListNumbersResponse struct {
	NextCursor string        `json:"next_cursor,omitempty"`
	Numbers    []PhoneNumber `json:"numbers,omitempty"`
	TotalCount int           `json:"total_count,omitempty"`
}

// CreateNumberRequest is the body of Numbers.Create.
type CreateNumberRequest struct {
	BYOConnectionID string                `json:"byo_connection_id,omitempty"`
	Connection      *BYOConnectionRequest `json:"connection,omitempty"`
	CountryCode     string                `json:"country_code,omitempty"`
	DID             string                `json:"did,omitempty"`
	Origin          NumberOrigin          `json:"origin,omitempty"`
	TrunkProvider   NumberTrunkProvider   `json:"trunk_provider,omitempty"`
}

// ListAvailableNumbersParams are the query parameters of Numbers.ListAvailable.
type ListAvailableNumbersParams struct {
	Country    string
	NumberType string
}

// AvailableNumber is a number offered by the provider's inventory.
type AvailableNumber struct {
	CountryCode     string `json:"country_code,omitempty"`
	PhoneNumber     string `json:"phone_number,omitempty"`
	PhoneNumberType string `json:"phone_number_type,omitempty"`
}

// ListAvailableNumbersResponse is the response of Numbers.ListAvailable.
type ListAvailableNumbersResponse struct {
	Numbers []AvailableNumber `json:"numbers,omitempty"`
}

// CreateNumberOrderRequest is the body of Numbers.CreateOrder.
type CreateNumberOrderRequest struct {
	AreaCode    string `json:"area_code,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	NumberType  string `json:"number_type,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	SetID       string `json:"set_id,omitempty"`
}

// NumberOrder is a number order.
type NumberOrder struct {
	CountryCode   string            `json:"country_code,omitempty"`
	CreatedAt     string            `json:"created_at,omitempty"`
	FailureReason string            `json:"failure_reason,omitempty"`
	ID            string            `json:"id,omitempty"`
	NumberType    string            `json:"number_type,omitempty"`
	PhoneNumber   string            `json:"phone_number,omitempty"`
	PhoneNumberID string            `json:"phone_number_id,omitempty"`
	SetID         string            `json:"set_id,omitempty"`
	Status        NumberOrderStatus `json:"status,omitempty"`
	TelnyxOrderID string            `json:"telnyx_order_id,omitempty"`
	UpdatedAt     string            `json:"updated_at,omitempty"`
}

// ListNumberOrdersResponse is the response of Numbers.ListOrders.
type ListNumberOrdersResponse struct {
	Orders []NumberOrder `json:"orders,omitempty"`
}

// AssignNumberRequest is the body of Numbers.Assign.
type AssignNumberRequest struct {
	AgentID string `json:"agent_id,omitempty"`
}

// SetNumberTrunkRequest is the body of Numbers.SetTrunk.
type SetNumberTrunkRequest struct {
	TrunkID string `json:"trunk_id,omitempty"`
}
