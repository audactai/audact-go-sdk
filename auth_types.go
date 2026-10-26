package audact

// TokenResponse is a short-lived access token.
type TokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
}

// ExchangeKeyResponse is the result of Auth.ExchangeKey.
type ExchangeKeyResponse = TokenResponse

// RefreshTokenResponse is the result of Auth.Refresh.
type RefreshTokenResponse = TokenResponse
