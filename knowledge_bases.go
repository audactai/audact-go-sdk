package audact

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// List lists knowledge bases, ordered by kb_id. Calls GET /v1/knowledge-bases.
func (s *KnowledgeBasesService) List(ctx context.Context, params *ListKnowledgeBasesParams) (*ListKnowledgeBasesResponse, error) {
	q := url.Values{}
	if params != nil {
		addInt(q, "limit", params.Limit)
		addString(q, "cursor", params.Cursor)
	}
	var out ListKnowledgeBasesResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/knowledge-bases", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a knowledge base. Calls POST /v1/knowledge-bases.
func (s *KnowledgeBasesService) Create(ctx context.Context, req *CreateKnowledgeBaseRequest) (*KnowledgeBase, error) {
	var out KnowledgeBase
	if _, err := s.t.do(ctx, http.MethodPost, "/knowledge-bases", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns a knowledge base. Calls GET /v1/knowledge-bases/{kb_id}.
func (s *KnowledgeBasesService) Get(ctx context.Context, kbID string) (*KnowledgeBase, error) {
	var out KnowledgeBase
	if _, err := s.t.do(ctx, http.MethodGet, kbPath(kbID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update renames a knowledge base and/or toggles auto-refresh; nil fields are
// not sent. Calls PATCH /v1/knowledge-bases/{kb_id}.
func (s *KnowledgeBasesService) Update(ctx context.Context, kbID string, req *UpdateKnowledgeBaseRequest) (*KnowledgeBase, error) {
	var out KnowledgeBase
	if _, err := s.t.do(ctx, http.MethodPatch, kbPath(kbID), authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a knowledge base with all its sources, chunks and agent
// links. This is irreversible. Calls DELETE /v1/knowledge-bases/{kb_id}.
func (s *KnowledgeBasesService) Delete(ctx context.Context, kbID string) error {
	_, err := s.t.do(ctx, http.MethodDelete, kbPath(kbID), authAPIKey, nil, nil, nil)
	return err
}

// ListSources lists the sources of a knowledge base.
// Calls GET /v1/knowledge-bases/{kb_id}/sources.
func (s *KnowledgeBasesService) ListSources(ctx context.Context, kbID string) (*ListKBSourcesResponse, error) {
	var out ListKBSourcesResponse
	if _, err := s.t.do(ctx, http.MethodGet, kbPath(kbID)+"/sources", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddSource adds a URL or text source and enqueues its ingestion; the server
// answers 202 and ingests asynchronously. Use AddFileSource for files.
// Calls POST /v1/knowledge-bases/{kb_id}/sources.
func (s *KnowledgeBasesService) AddSource(ctx context.Context, kbID string, req *AddKBSourceRequest) (*KBSource, error) {
	var out KBSource
	if _, err := s.t.do(ctx, http.MethodPost, kbPath(kbID)+"/sources", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddFileSource uploads a file source (multipart/form-data, at most 50 MiB)
// and enqueues its ingestion; the server answers 202.
// Calls POST /v1/knowledge-bases/{kb_id}/sources.
func (s *KnowledgeBasesService) AddFileSource(ctx context.Context, kbID string, req *AddKBFileSourceRequest) (*KBSource, error) {
	if req == nil || req.File == nil {
		return nil, errors.New("audact: AddFileSource: file is required")
	}
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	filename := req.Filename
	if filename == "" {
		filename = "upload"
	}
	fw, err := mw.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("audact: AddFileSource: %w", err)
	}
	if _, err := io.Copy(fw, req.File); err != nil {
		return nil, fmt.Errorf("audact: AddFileSource: read file: %w", err)
	}
	if req.Filename != "" {
		if err := mw.WriteField("filename", req.Filename); err != nil {
			return nil, fmt.Errorf("audact: AddFileSource: %w", err)
		}
	}
	if req.Config != nil {
		cfg, err := json.Marshal(req.Config)
		if err != nil {
			return nil, fmt.Errorf("audact: AddFileSource: encode config: %w", err)
		}
		if err := mw.WriteField("config", string(cfg)); err != nil {
			return nil, fmt.Errorf("audact: AddFileSource: %w", err)
		}
	}
	if err := mw.Close(); err != nil {
		return nil, fmt.Errorf("audact: AddFileSource: %w", err)
	}
	var out KBSource
	body := &rawBody{contentType: mw.FormDataContentType(), data: buf.Bytes()}
	if _, err := s.t.do(ctx, http.MethodPost, kbPath(kbID)+"/sources", authAPIKey, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSource returns a source.
// Calls GET /v1/knowledge-bases/{kb_id}/sources/{source_id}.
func (s *KnowledgeBasesService) GetSource(ctx context.Context, kbID, sourceID string) (*KBSource, error) {
	var out KBSource
	if _, err := s.t.do(ctx, http.MethodGet, kbSourcePath(kbID, sourceID), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateSource updates a source's title and/or (text sources only) content;
// editing content re-ingests the source under a new ingest_version. nil
// fields are not sent. Calls PATCH /v1/knowledge-bases/{kb_id}/sources/{source_id}.
func (s *KnowledgeBasesService) UpdateSource(ctx context.Context, kbID, sourceID string, req *UpdateKBSourceRequest) (*KBSource, error) {
	var out KBSource
	if _, err := s.t.do(ctx, http.MethodPatch, kbSourcePath(kbID, sourceID), authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ReingestSource re-queues a source for ingestion under a new ingest_version;
// the server answers 202.
// Calls POST /v1/knowledge-bases/{kb_id}/sources/{source_id}/reingest.
func (s *KnowledgeBasesService) ReingestSource(ctx context.Context, kbID, sourceID string) (*KBSource, error) {
	var out KBSource
	if _, err := s.t.do(ctx, http.MethodPost, kbSourcePath(kbID, sourceID)+"/reingest", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteSource deletes a source and its chunks.
// Calls DELETE /v1/knowledge-bases/{kb_id}/sources/{source_id}.
func (s *KnowledgeBasesService) DeleteSource(ctx context.Context, kbID, sourceID string) error {
	_, err := s.t.do(ctx, http.MethodDelete, kbSourcePath(kbID, sourceID), authAPIKey, nil, nil, nil)
	return err
}

// ListAgentLinks lists the agents linked to a knowledge base.
// Calls GET /v1/knowledge-bases/{kb_id}/agents.
func (s *KnowledgeBasesService) ListAgentLinks(ctx context.Context, kbID string) (*ListKBAgentLinksResponse, error) {
	var out ListKBAgentLinksResponse
	if _, err := s.t.do(ctx, http.MethodGet, kbPath(kbID)+"/agents", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListLinksByAgent lists the knowledge bases linked to an agent.
// Calls GET /v1/knowledge-bases/agent-links?agent_id=.
func (s *KnowledgeBasesService) ListLinksByAgent(ctx context.Context, agentID string) (*ListKBAgentLinksResponse, error) {
	q := url.Values{}
	q.Set("agent_id", agentID)
	var out ListKBAgentLinksResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/knowledge-bases/agent-links", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LinkAgent creates or replaces the link between a knowledge base and an
// agent. A nil req sends no body, so the server defaults apply.
// Calls PUT /v1/knowledge-bases/{kb_id}/agents/{agent_id}.
func (s *KnowledgeBasesService) LinkAgent(ctx context.Context, kbID, agentID string, req *LinkKBAgentRequest) (*KBAgentLink, error) {
	var in any
	if req != nil {
		in = req
	}
	var out KBAgentLink
	if _, err := s.t.do(ctx, http.MethodPut, kbAgentPath(kbID, agentID), authAPIKey, nil, in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateAgentLink partially updates the retrieval tuning of an agent link;
// nil fields keep their current values.
// Calls PATCH /v1/knowledge-bases/{kb_id}/agents/{agent_id}.
func (s *KnowledgeBasesService) UpdateAgentLink(ctx context.Context, kbID, agentID string, req *UpdateKBAgentLinkRequest) (*KBAgentLink, error) {
	var out KBAgentLink
	if _, err := s.t.do(ctx, http.MethodPatch, kbAgentPath(kbID, agentID), authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UnlinkAgent removes the link between a knowledge base and an agent.
// Calls DELETE /v1/knowledge-bases/{kb_id}/agents/{agent_id}.
func (s *KnowledgeBasesService) UnlinkAgent(ctx context.Context, kbID, agentID string) error {
	_, err := s.t.do(ctx, http.MethodDelete, kbAgentPath(kbID, agentID), authAPIKey, nil, nil, nil)
	return err
}

// ListExclusions lists the crawl exclusion rules of a knowledge base.
// Calls GET /v1/knowledge-bases/{kb_id}/exclusions.
func (s *KnowledgeBasesService) ListExclusions(ctx context.Context, kbID string) (*ListKBExclusionsResponse, error) {
	var out ListKBExclusionsResponse
	if _, err := s.t.do(ctx, http.MethodGet, kbPath(kbID)+"/exclusions", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddExclusion adds a crawl exclusion rule.
// Calls POST /v1/knowledge-bases/{kb_id}/exclusions.
func (s *KnowledgeBasesService) AddExclusion(ctx context.Context, kbID string, req *AddKBExclusionRequest) (*KBCrawlExclusion, error) {
	var out KBCrawlExclusion
	if _, err := s.t.do(ctx, http.MethodPost, kbPath(kbID)+"/exclusions", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteExclusion removes a crawl exclusion rule, reinstating the path or URL.
// Calls DELETE /v1/knowledge-bases/{kb_id}/exclusions/{id}.
func (s *KnowledgeBasesService) DeleteExclusion(ctx context.Context, kbID, exclusionID string) error {
	_, err := s.t.do(ctx, http.MethodDelete, kbPath(kbID)+"/exclusions/"+url.PathEscape(exclusionID), authAPIKey, nil, nil, nil)
	return err
}

// TestRetrieval runs an ad-hoc retrieval against a knowledge base. No hits
// above the threshold is a success with Empty set, not an error.
// Calls POST /v1/knowledge-bases/{kb_id}/test-retrieval.
func (s *KnowledgeBasesService) TestRetrieval(ctx context.Context, kbID string, req *TestKBRetrievalRequest) (*TestKBRetrievalResponse, error) {
	var out TestKBRetrievalResponse
	if _, err := s.t.do(ctx, http.MethodPost, kbPath(kbID)+"/test-retrieval", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func kbPath(kbID string) string {
	return "/knowledge-bases/" + url.PathEscape(kbID)
}

func kbSourcePath(kbID, sourceID string) string {
	return kbPath(kbID) + "/sources/" + url.PathEscape(sourceID)
}

func kbAgentPath(kbID, agentID string) string {
	return kbPath(kbID) + "/agents/" + url.PathEscape(agentID)
}
