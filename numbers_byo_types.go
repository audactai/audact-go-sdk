package audact

// BYOConfig holds deployment-wide values needed to configure a BYO trunk.
type BYOConfig struct {
	SignalingIPs   []string `json:"signaling_ips,omitempty"`
	SIPIngressHost string   `json:"sip_ingress_host,omitempty"`
}

// ListBYOConnectionsParams are the optional query parameters of ListBYOConnections.
type ListBYOConnectionsParams struct {
	Limit  int
	Cursor string
}

// CreateBYOConnectionRequest creates a BYO SIP-trunk connection.
type CreateBYOConnectionRequest struct {
	InboundAllowedIPs []string       `json:"inbound_allowed_ips,omitempty"`
	Name              string         `json:"name"`
	Provider          NumberProvider `json:"provider"`
	SIPHost           string         `json:"sip_host"`
	SIPPassword       string         `json:"sip_password,omitempty"`
	SIPTransport      SIPTransport   `json:"sip_transport,omitempty"`
	SIPUsername       string         `json:"sip_username,omitempty"`
}

// UpdateBYOConnectionRequest partially updates a BYO connection; unset fields
// are not sent. A non-nil InboundAllowedIPs (even empty) replaces the stored allowlist.
type UpdateBYOConnectionRequest struct {
	InboundAllowedIPs *[]string      `json:"inbound_allowed_ips,omitempty"`
	Name              string         `json:"name,omitempty"`
	Provider          NumberProvider `json:"provider,omitempty"`
	SIPHost           string         `json:"sip_host,omitempty"`
	SIPPassword       string         `json:"sip_password,omitempty"`
	SIPTransport      SIPTransport   `json:"sip_transport,omitempty"`
	SIPUsername       string         `json:"sip_username,omitempty"`
}

// BYOConnection is a BYO SIP-trunk connection.
type BYOConnection struct {
	ConnectionHealth  ConnectionHealth `json:"connection_health,omitempty"`
	CreatedAt         string           `json:"created_at,omitempty"`
	HasCredentials    bool             `json:"has_credentials,omitempty"`
	HealthDetail      string           `json:"health_detail,omitempty"`
	ID                string           `json:"id,omitempty"`
	InboundAllowedIPs []string         `json:"inbound_allowed_ips,omitempty"`
	LastErrorAt       string           `json:"last_error_at,omitempty"`
	LastErrorCode     int              `json:"last_error_code,omitempty"`
	LastVerifiedAt    string           `json:"last_verified_at,omitempty"`
	Name              string           `json:"name,omitempty"`
	Provider          NumberProvider   `json:"provider,omitempty"`
	SIPHost           string           `json:"sip_host,omitempty"`
	SIPTransport      SIPTransport     `json:"sip_transport,omitempty"`
	UpdatedAt         string           `json:"updated_at,omitempty"`
}

// ListBYOConnectionsResponse is a page of BYO connections.
type ListBYOConnectionsResponse struct {
	Connections []BYOConnection `json:"connections,omitempty"`
	NextCursor  string          `json:"next_cursor,omitempty"`
}
