package audact

// IntegrationHealth is the last-observed reachability or token state of a connection.
type IntegrationHealth string

// Integration health values.
const (
	IntegrationHealthConnected    IntegrationHealth = "connected"
	IntegrationHealthTokenExpired IntegrationHealth = "token_expired"
	IntegrationHealthRevoked      IntegrationHealth = "revoked"
	IntegrationHealthError        IntegrationHealth = "error"
	IntegrationHealthUnknown      IntegrationHealth = "unknown"
)

// IntegrationStatus is the lifecycle state of a connection.
type IntegrationStatus string

// Integration status values.
const (
	IntegrationStatusActive  IntegrationStatus = "active"
	IntegrationStatusRevoked IntegrationStatus = "revoked"
	IntegrationStatusError   IntegrationStatus = "error"
)

// IntegrationAuthMethod is how a connection authenticates with its provider.
type IntegrationAuthMethod string

// Integration auth method values.
const (
	IntegrationAuthMethodAPIKey IntegrationAuthMethod = "api_key"
	IntegrationAuthMethodOAuth  IntegrationAuthMethod = "oauth"
)

// Connection is a third-party integration connection.
type Connection struct {
	AuthMethod     IntegrationAuthMethod `json:"auth_method,omitempty"`
	CalUsername    string                `json:"cal_username,omitempty"`
	ConnectionID   string                `json:"connection_id,omitempty"`
	CreatedAt      string                `json:"created_at,omitempty"`
	DeletedAt      string                `json:"deleted_at,omitempty"`
	HasAPIKey      bool                  `json:"has_api_key,omitempty"`
	HasOAuthToken  bool                  `json:"has_oauth_token,omitempty"`
	Health         IntegrationHealth     `json:"health,omitempty"`
	LastCheckedAt  string                `json:"last_checked_at,omitempty"`
	LastError      string                `json:"last_error,omitempty"`
	Name           string                `json:"name,omitempty"`
	OAuthExpiresAt string                `json:"oauth_expires_at,omitempty"`
	OAuthScope     string                `json:"oauth_scope,omitempty"`
	Provider       string                `json:"provider,omitempty"`
	Status         IntegrationStatus     `json:"status,omitempty"`
	TenantID       string                `json:"tenant_id,omitempty"`
	UpdatedAt      string                `json:"updated_at,omitempty"`
}

// ListConnectionsParams holds optional query parameters for listing connections.
type ListConnectionsParams struct {
	// Provider filters by provider, for example "cal_com".
	Provider string
}

// ListConnectionsResponse is the response of listing connections.
type ListConnectionsResponse struct {
	Connections []Connection `json:"connections,omitempty"`
}

// CreateConnectionRequest is the body for creating a connection.
type CreateConnectionRequest struct {
	APIKey      string `json:"api_key,omitempty"`
	AuthMethod  string `json:"auth_method,omitempty"`
	CalUsername string `json:"cal_username,omitempty"`
	Name        string `json:"name,omitempty"`
	Provider    string `json:"provider,omitempty"`
}

// DeleteConnectionResponse is the response of deleting a connection.
type DeleteConnectionResponse struct {
	Deleted bool `json:"deleted,omitempty"`
}

// CalcomEventType is a Cal.com event type.
type CalcomEventType struct {
	Hidden          bool   `json:"hidden,omitempty"`
	ID              int    `json:"id,omitempty"`
	LengthInMinutes int    `json:"lengthInMinutes,omitempty"`
	Slug            string `json:"slug,omitempty"`
	Title           string `json:"title,omitempty"`
}

// ListEventTypesResponse is the response of listing a connection's event types.
type ListEventTypesResponse struct {
	EventTypes []CalcomEventType `json:"event_types,omitempty"`
}

// TestConnectionResponse is the response of testing a connection.
type TestConnectionResponse struct {
	ConnectionID   string `json:"connection_id,omitempty"`
	Error          string `json:"error,omitempty"`
	EventTypeCount int    `json:"event_type_count,omitempty"`
	Health         string `json:"health,omitempty"`
	Healthy        bool   `json:"healthy,omitempty"`
}

// CalcomOAuthStartParams holds optional query parameters for starting Cal.com OAuth.
type CalcomOAuthStartParams struct {
	Name        string
	CalUsername string
}

// CalcomOAuthStartResponse is the response of starting Cal.com OAuth.
type CalcomOAuthStartResponse struct {
	AuthorizeURL string `json:"authorize_url,omitempty"`
}

// CalcomOAuthCallbackParams holds the query parameters of the Cal.com OAuth callback.
type CalcomOAuthCallbackParams struct {
	// State is the single-use state token from CalcomOAuthStart (required).
	State string
	// Code is the authorization code; absent when Error is set.
	Code string
	// Error is Cal.com's error code when the user declined consent.
	Error string
}

// CalcomOAuthCallbackResponse is the unfollowed redirect returned by the callback.
type CalcomOAuthCallbackResponse struct {
	// StatusCode is the HTTP status, normally 302.
	StatusCode int
	// Location is the redirect target.
	Location string
}
