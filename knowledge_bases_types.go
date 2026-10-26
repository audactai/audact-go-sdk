package audact

import "io"

// KBSourceType is the origin of a knowledge-base source.
type KBSourceType string

// Knowledge-base source type values.
const (
	KBSourceTypeURL  KBSourceType = "url"
	KBSourceTypeFile KBSourceType = "file"
	KBSourceTypeText KBSourceType = "text"
)

// KBSourceStatus is the ingestion lifecycle state of a knowledge-base source.
type KBSourceStatus string

// Knowledge-base source status values.
const (
	KBSourceStatusPending    KBSourceStatus = "pending"
	KBSourceStatusProcessing KBSourceStatus = "processing"
	KBSourceStatusReady      KBSourceStatus = "ready"
	KBSourceStatusFailed     KBSourceStatus = "failed"
)

// KBRobotsPolicy controls how the crawler treats robots.txt for URL sources.
type KBRobotsPolicy string

const (
	// KBRobotsPolicyHonor obeys robots.txt (the server default).
	KBRobotsPolicyHonor KBRobotsPolicy = "honor"
	// KBRobotsPolicyOverride bypasses robots.txt; only for domains the tenant owns or has attested.
	KBRobotsPolicyOverride KBRobotsPolicy = "override"
)

// KBRetrievalTurnRole is the speaker of one turn in a retrieval test.
type KBRetrievalTurnRole string

const (
	// KBRetrievalTurnRoleUser is a caller turn.
	KBRetrievalTurnRoleUser KBRetrievalTurnRole = "user"
	// KBRetrievalTurnRoleAssistant is an agent turn.
	KBRetrievalTurnRoleAssistant KBRetrievalTurnRole = "assistant"
	// KBRetrievalTurnRoleSystem is a system turn; it is excluded from query condensation.
	KBRetrievalTurnRoleSystem KBRetrievalTurnRole = "system"
)

// KnowledgeBase is a tenant knowledge base used for retrieval-augmented generation.
type KnowledgeBase struct {
	AutoRefresh  bool           `json:"auto_refresh,omitempty"`
	CreatedAt    string         `json:"created_at,omitempty"`
	KBID         string         `json:"kb_id,omitempty"`
	Name         string         `json:"name,omitempty"`
	RobotsPolicy KBRobotsPolicy `json:"robots_policy,omitempty"`
	TenantID     string         `json:"tenant_id,omitempty"`
	UpdatedAt    string         `json:"updated_at,omitempty"`
}

// CreateKnowledgeBaseRequest is the body of KnowledgeBases.Create. Name is
// required (1-200 characters); RobotsPolicy is KBRobotsPolicyHonor (default) or KBRobotsPolicyOverride.
type CreateKnowledgeBaseRequest struct {
	AutoRefresh  bool           `json:"auto_refresh,omitempty"`
	Name         string         `json:"name"`
	RobotsPolicy KBRobotsPolicy `json:"robots_policy,omitempty"`
}

// UpdateKnowledgeBaseRequest is the body of KnowledgeBases.Update; nil fields are not sent.
type UpdateKnowledgeBaseRequest struct {
	AutoRefresh *bool   `json:"auto_refresh,omitempty"`
	Name        *string `json:"name,omitempty"`
}

// ListKnowledgeBasesParams holds optional query parameters for KnowledgeBases.List.
type ListKnowledgeBasesParams struct {
	Limit  int
	Cursor string
}

// ListKnowledgeBasesResponse is a page of knowledge bases, ordered by kb_id.
// Pass NextCursor as ListKnowledgeBasesParams.Cursor while HasMore is true.
type ListKnowledgeBasesResponse struct {
	HasMore    bool            `json:"has_more,omitempty"`
	Items      []KnowledgeBase `json:"items,omitempty"`
	NextCursor string          `json:"next_cursor,omitempty"`
}

// KBSource is a URL, file or text source of a knowledge base.
//
// ChunkCount, CrawlCoverage and LastIngestedAt describe the source's latest
// successful ingestion; they are nil/empty until a first ingestion has
// finished, and on sources returned by older servers.
type KBSource struct {
	ChunkCount     *int             `json:"chunk_count,omitempty"`
	Config         map[string]any   `json:"config,omitempty"`
	ContentPreview string           `json:"content_preview,omitempty"`
	CrawlCoverage  *KBCrawlCoverage `json:"crawl_coverage,omitempty"`
	CreatedAt      string           `json:"created_at,omitempty"`
	Error          string           `json:"error,omitempty"`
	IngestVersion  int              `json:"ingest_version,omitempty"`
	KBID           string           `json:"kb_id,omitempty"`
	LastIngestedAt string           `json:"last_ingested_at,omitempty"`
	SourceID       string           `json:"source_id,omitempty"`
	Status         KBSourceStatus   `json:"status,omitempty"`
	TenantID       string           `json:"tenant_id,omitempty"`
	Title          string           `json:"title,omitempty"`
	Type           KBSourceType     `json:"type,omitempty"`
	UpdatedAt      string           `json:"updated_at,omitempty"`
	URI            string           `json:"uri,omitempty"`
}

// KBCrawlCoverage is the page coverage of a source's latest successful
// ingestion: pages discovered, stored, excluded by the knowledge base's
// exclusion list, and failed. A single page, file or text source reports
// 1 discovered / 1 stored.
type KBCrawlCoverage struct {
	Discovered int `json:"discovered"`
	Excluded   int `json:"excluded"`
	Failed     int `json:"failed"`
	Stored     int `json:"stored"`
}

// AddKBSourceRequest is the JSON body of KnowledgeBases.AddSource. Type is
// KBSourceTypeURL (URL required) or KBSourceTypeText (Content required); use
// KnowledgeBases.AddFileSource for file uploads.
type AddKBSourceRequest struct {
	Config  map[string]any `json:"config,omitempty"`
	Content string         `json:"content,omitempty"`
	Title   string         `json:"title,omitempty"`
	Type    KBSourceType   `json:"type"`
	URL     string         `json:"url,omitempty"`
}

// AddKBFileSourceRequest is the multipart upload of KnowledgeBases.AddFileSource
// (at most 50 MiB). Filename defaults to "upload" server-side; Config, when
// set, is sent JSON-encoded in the "config" form field.
type AddKBFileSourceRequest struct {
	File     io.Reader
	Filename string
	Config   map[string]any
}

// UpdateKBSourceRequest is the body of KnowledgeBases.UpdateSource; nil
// fields are not sent. Content is only editable on text sources.
type UpdateKBSourceRequest struct {
	Content *string `json:"content,omitempty"`
	Title   *string `json:"title,omitempty"`
}

// ListKBSourcesResponse lists the sources of a knowledge base.
type ListKBSourcesResponse struct {
	Items []KBSource `json:"items,omitempty"`
	Total int        `json:"total,omitempty"`
}

// KBAgentLink links a knowledge base to an agent with retrieval tuning.
type KBAgentLink struct {
	AgentID              string  `json:"agent_id,omitempty"`
	CreatedAt            string  `json:"created_at,omitempty"`
	KBID                 string  `json:"kb_id,omitempty"`
	RetrievalInstruction string  `json:"retrieval_instruction,omitempty"`
	SimilarityThreshold  float64 `json:"similarity_threshold,omitempty"`
	TenantID             string  `json:"tenant_id,omitempty"`
	TopK                 int     `json:"top_k,omitempty"`
}

// LinkKBAgentRequest is the optional body of KnowledgeBases.LinkAgent; nil
// fields take the server defaults (top_k 3, similarity_threshold 0.6,
// retrieval_instruction "").
type LinkKBAgentRequest struct {
	RetrievalInstruction *string  `json:"retrieval_instruction,omitempty"`
	SimilarityThreshold  *float64 `json:"similarity_threshold,omitempty"`
	TopK                 *int     `json:"top_k,omitempty"`
}

// UpdateKBAgentLinkRequest is the body of KnowledgeBases.UpdateAgentLink; nil
// fields keep their current values.
type UpdateKBAgentLinkRequest struct {
	RetrievalInstruction *string  `json:"retrieval_instruction,omitempty"`
	SimilarityThreshold  *float64 `json:"similarity_threshold,omitempty"`
	TopK                 *int     `json:"top_k,omitempty"`
}

// ListKBAgentLinksResponse lists agent-knowledge-base links.
type ListKBAgentLinksResponse struct {
	Items []KBAgentLink `json:"items,omitempty"`
	Total int           `json:"total,omitempty"`
}

// KBCrawlExclusion is a crawl exclusion rule: either a path prefix
// (CrawlPath) or an exact URL (URL); the other field is empty.
type KBCrawlExclusion struct {
	CrawlPath string `json:"crawl_path,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	ID        string `json:"id,omitempty"`
	KBID      string `json:"kb_id,omitempty"`
	TenantID  string `json:"tenant_id,omitempty"`
	URL       string `json:"url,omitempty"`
}

// AddKBExclusionRequest is the body of KnowledgeBases.AddExclusion. Set
// exactly one of CrawlPath and URL (each at most 200 characters).
type AddKBExclusionRequest struct {
	CrawlPath string `json:"crawl_path,omitempty"`
	URL       string `json:"url,omitempty"`
}

// ListKBExclusionsResponse lists the crawl exclusion rules of a knowledge base.
type ListKBExclusionsResponse struct {
	Items []KBCrawlExclusion `json:"items,omitempty"`
	Total int                `json:"total,omitempty"`
}

// KBRetrievalTurn is one conversation turn of a retrieval test; Role is
// KBRetrievalTurnRoleUser, KBRetrievalTurnRoleAssistant or KBRetrievalTurnRoleSystem.
type KBRetrievalTurn struct {
	Content string              `json:"content"`
	Role    KBRetrievalTurnRole `json:"role"`
}

// TestKBRetrievalRequest is the body of KnowledgeBases.TestRetrieval. Set
// exactly one of Question and Turns; TopK (1-10) and Threshold (0.0-1.0)
// optionally override the defaults (3 and 0.0).
type TestKBRetrievalRequest struct {
	Question  string            `json:"question,omitempty"`
	Threshold *float64          `json:"threshold,omitempty"`
	TopK      *int              `json:"top_k,omitempty"`
	Turns     []KBRetrievalTurn `json:"turns,omitempty"`
}

// KBRetrievalHit is a chunk returned by a retrieval test.
type KBRetrievalHit struct {
	KBID       string  `json:"kb_id,omitempty"`
	Similarity float64 `json:"similarity,omitempty"`
	SourceID   string  `json:"source_id,omitempty"`
	Text       string  `json:"text,omitempty"`
}

// TestKBRetrievalResponse is the result of KnowledgeBases.TestRetrieval. An
// empty result (Empty true, no Hits) is a success, not an error.
type TestKBRetrievalResponse struct {
	CondensedQuery string           `json:"condensed_query,omitempty"`
	Empty          bool             `json:"empty,omitempty"`
	Hits           []KBRetrievalHit `json:"hits,omitempty"`
}
