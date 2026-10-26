package audact

import (
	"net/http"
	"time"
)

// DefaultBaseURL is the production API base URL used when none is configured.
const DefaultBaseURL = "https://api.audact.ai"

// Option configures a Client.
type Option func(*Client) error

// WithBaseURL sets the API base URL. It takes precedence over AUDACT_BASE_URL.
func WithBaseURL(u string) Option {
	return func(c *Client) error {
		c.baseURLOpt = u
		c.baseURLSet = true
		return nil
	}
}

// WithHTTPClient sets the HTTP client. The default has a 30 second timeout.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) error {
		if hc != nil {
			c.httpClient = hc
		}
		return nil
	}
}

// WithOrgToken sets the organisation bearer token.
func WithOrgToken(t string) Option {
	return func(c *Client) error {
		c.orgToken = t
		return nil
	}
}

// WithPlatformToken sets the initial platform bearer token. Platform tokens
// are single-use and action-scoped, so a token set here serves one Platform
// call; use Client.SetPlatformToken to set a fresh one for each later call.
func WithPlatformToken(t string) Option {
	return func(c *Client) error {
		c.platformToken = t
		return nil
	}
}

// WithUserAgent appends a custom product token to the User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) error {
		c.userAgent = ua
		return nil
	}
}

func defaultHTTPClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}
