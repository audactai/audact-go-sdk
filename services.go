package audact

// AgentsService groups agent and voice operations.
type AgentsService struct{ t *transport }

// CallsService groups call record and call control operations.
type CallsService struct{ t *transport }

// NumbersService groups phone number, KYC and BYO operations.
type NumbersService struct{ t *transport }

// WebhooksService groups webhook operations.
type WebhooksService struct{ t *transport }

// ToolsService groups tool operations.
type ToolsService struct{ t *transport }

// IntegrationsService groups integration operations.
type IntegrationsService struct{ t *transport }

// KnowledgeBasesService groups knowledge-base (RAG) operations: knowledge
// bases, sources, agent links, crawl exclusions and retrieval tests.
type KnowledgeBasesService struct{ t *transport }

// EvidenceService groups evidence operations.
type EvidenceService struct{ t *transport }

// PIIService groups PII operations.
type PIIService struct{ t *transport }

// ValidatorService groups validator operations.
type ValidatorService struct{ t *transport }

// LLMProxyService groups LLM proxy operations.
type LLMProxyService struct{ t *transport }

// TenantsService groups tenant operations.
type TenantsService struct{ t *transport }

// OrgService groups organisation operations.
type OrgService struct{ t *transport }

// PlatformOrgsService groups platform organisation administration operations.
type PlatformOrgsService struct{ t *transport }

// PlatformService groups platform operations.
type PlatformService struct {
	t *transport

	// Orgs groups platform organisation administration operations.
	Orgs *PlatformOrgsService
}

// AggregationService groups aggregation operations.
type AggregationService struct{ t *transport }

// PolicyService groups policy operations.
type PolicyService struct{ t *transport }

// AudioService groups audio validation operations.
type AudioService struct{ t *transport }

// AuthService groups authentication operations.
type AuthService struct{ t *transport }

func (c *Client) initServices() {
	t := c.t
	c.Agents = &AgentsService{t}
	c.Calls = &CallsService{t}
	c.Numbers = &NumbersService{t}
	c.Webhooks = &WebhooksService{t}
	c.Tools = &ToolsService{t}
	c.Integrations = &IntegrationsService{t}
	c.KnowledgeBases = &KnowledgeBasesService{t}
	c.Evidence = &EvidenceService{t}
	c.PII = &PIIService{t}
	c.Validator = &ValidatorService{t}
	c.LLMProxy = &LLMProxyService{t}
	c.Tenants = &TenantsService{t}
	c.Org = &OrgService{t}
	c.Platform = &PlatformService{t: t, Orgs: &PlatformOrgsService{t}}
	c.Aggregation = &AggregationService{t}
	c.Policy = &PolicyService{t}
	c.Audio = &AudioService{t}
	c.Auth = &AuthService{t}
}
