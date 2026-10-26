package audact

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Client is the Audact API client. It is safe for concurrent use.
type Client struct {
	Agents         *AgentsService
	Calls          *CallsService
	Numbers        *NumbersService
	Webhooks       *WebhooksService
	Tools          *ToolsService
	Integrations   *IntegrationsService
	KnowledgeBases *KnowledgeBasesService
	Evidence       *EvidenceService
	PII            *PIIService
	Validator      *ValidatorService
	LLMProxy       *LLMProxyService
	Tenants        *TenantsService
	Org            *OrgService
	Platform       *PlatformService
	Aggregation    *AggregationService
	Policy         *PolicyService
	Audio          *AudioService
	Auth           *AuthService

	t *transport

	baseURLOpt    string
	baseURLSet    bool
	httpClient    *http.Client
	orgToken      string
	platformToken string
	userAgent     string
}

// NewClient creates a Client. An empty apiKey falls back to AUDACT_API_KEY.
// A missing key is not an error (some operations are public; the server
// answers 401 otherwise). NewClient returns an error only for a nil option, an
// option that fails, or an invalid base URL (unparseable, or a scheme other
// than http or https).
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	c := &Client{httpClient: defaultHTTPClient()}
	for _, o := range opts {
		if o == nil {
			return nil, errors.New("audact: nil option")
		}
		if err := o(c); err != nil {
			return nil, fmt.Errorf("audact: apply option: %w", err)
		}
	}
	if apiKey == "" {
		apiKey = os.Getenv("AUDACT_API_KEY")
	}
	raw := DefaultBaseURL
	switch {
	case c.baseURLSet:
		raw = c.baseURLOpt
	case os.Getenv("AUDACT_BASE_URL") != "":
		raw = os.Getenv("AUDACT_BASE_URL")
	}
	base, err := normalizeBaseURL(raw)
	if err != nil {
		return nil, err
	}
	ua := "audact-go/" + Version
	if c.userAgent != "" {
		ua += " " + c.userAgent
	}
	c.t = &transport{
		baseURL:       base,
		hc:            c.httpClient,
		userAgent:     ua,
		apiKey:        apiKey,
		orgToken:      c.orgToken,
		platformToken: c.platformToken,
	}
	c.initServices()
	return c, nil
}

// SetPlatformToken replaces the platform bearer token used by Platform
// methods. Platform tokens are single-use (the server rejects a replayed jti)
// and scoped to one action, so mint and set a fresh token before each Platform
// call. An empty string clears it. Safe for concurrent use.
func (c *Client) SetPlatformToken(tok string) { c.t.setPlatformToken(tok) }

// BaseURL returns the normalised base URL (no trailing slash and no /v1).
func (c *Client) BaseURL() string { return c.t.baseURL }

func normalizeBaseURL(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("audact: invalid base URL: %w", err)
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", fmt.Errorf("audact: invalid base URL %q: want http or https URL with host", raw)
	}
	s := strings.TrimRight(u.String(), "/")
	s = strings.TrimSuffix(s, "/v1")
	s = strings.TrimRight(s, "/")
	return s, nil
}
