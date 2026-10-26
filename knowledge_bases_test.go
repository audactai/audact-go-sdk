package audact_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestKnowledgeBasesService(t *testing.T) {
	ctx := context.Background()
	zero := 0.0
	topK := 5
	empty := ""
	title := "T"
	f := false
	src := `{"source_id":"s1","kb_id":"k1","tenant_id":"t1","type":"url","status":"pending","uri":"https://x","ingest_version":1,"config":{"depth":2},"title":"FAQ","content_preview":"p","error":"boom","created_at":"c","updated_at":"u"}`
	link := `{"kb_id":"k1","agent_id":"a1","tenant_id":"t1","top_k":3,"similarity_threshold":0.6,"retrieval_instruction":"r","created_at":"c"}`
	excl := `{"id":"e1","kb_id":"k1","tenant_id":"t1","crawl_path":"/admin","created_at":"c"}`
	kb := `{"kb_id":"k1","tenant_id":"t1","name":"FAQ","auto_refresh":true,"robots_policy":"honor","created_at":"c","updated_at":"u"}`
	tests := []struct {
		name      string
		status    int
		resp      string
		call      func(c *audact.Client) (any, error)
		wantM     string
		wantPath  string
		wantQuery string
		wantBody  string // JSON, "" = no body
		wantOut   string // JSON of result, "" = none
	}{
		{"list", 200, `{"items":[` + kb + `],"has_more":true,"next_cursor":"n"}`,
			func(c *audact.Client) (any, error) { return c.KnowledgeBases.List(ctx, nil) },
			"GET", "/v1/knowledge-bases", "", ``, `{"items":[` + kb + `],"has_more":true,"next_cursor":"n"}`},
		{"list params", 200, `{}`,
			func(c *audact.Client) (any, error) {
				return c.KnowledgeBases.List(ctx, &audact.ListKnowledgeBasesParams{Limit: 5, Cursor: "c/1"})
			},
			"GET", "/v1/knowledge-bases", "cursor=c%2F1&limit=5", ``, `{}`},
		{"create", 201, kb,
			func(c *audact.Client) (any, error) {
				return c.KnowledgeBases.Create(ctx, &audact.CreateKnowledgeBaseRequest{Name: "FAQ", AutoRefresh: true, RobotsPolicy: audact.KBRobotsPolicyOverride})
			},
			"POST", "/v1/knowledge-bases", "", `{"name":"FAQ","auto_refresh":true,"robots_policy":"override"}`, kb},
		{"create nil sends required name", 201, `{}`,
			func(c *audact.Client) (any, error) { return c.KnowledgeBases.Create(ctx, nil) },
			"POST", "/v1/knowledge-bases", "", `{"name":""}`, `{}`},
		{"get", 200, kb,
			func(c *audact.Client) (any, error) { return c.KnowledgeBases.Get(ctx, "k/1") },
			"GET", "/v1/knowledge-bases/k%2F1", "", ``, kb},
		{"update omits unset", 200, kb,
			func(c *audact.Client) (any, error) {
				return c.KnowledgeBases.Update(ctx, "k1", &audact.UpdateKnowledgeBaseRequest{AutoRefresh: &f})
			},
			"PATCH", "/v1/knowledge-bases/k1", "", `{"auto_refresh":false}`, kb},
		{"delete", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.KnowledgeBases.Delete(ctx, "k/1") },
			"DELETE", "/v1/knowledge-bases/k%2F1", "", ``, ``},
		{"list sources", 200, `{"items":[` + src + `],"total":1}`,
			func(c *audact.Client) (any, error) { return c.KnowledgeBases.ListSources(ctx, "k/1") },
			"GET", "/v1/knowledge-bases/k%2F1/sources", "", ``, `{"items":[` + src + `],"total":1}`},
		{"add source url", 202, src,
			func(c *audact.Client) (any, error) {
				return c.KnowledgeBases.AddSource(ctx, "k1", &audact.AddKBSourceRequest{Type: audact.KBSourceTypeURL, URL: "https://x", Title: "FAQ", Config: map[string]any{"depth": 2}})
			},
			"POST", "/v1/knowledge-bases/k1/sources", "", `{"type":"url","url":"https://x","title":"FAQ","config":{"depth":2}}`, src},
		{"add source text", 202, `{"source_id":"s2","type":"text"}`,
			func(c *audact.Client) (any, error) {
				return c.KnowledgeBases.AddSource(ctx, "k/1", &audact.AddKBSourceRequest{Type: audact.KBSourceTypeText, Content: "hello"})
			},
			"POST", "/v1/knowledge-bases/k%2F1/sources", "", `{"type":"text","content":"hello"}`, `{"source_id":"s2","type":"text"}`},
		{"get source", 200, src,
			func(c *audact.Client) (any, error) { return c.KnowledgeBases.GetSource(ctx, "k/1", "s/1") },
			"GET", "/v1/knowledge-bases/k%2F1/sources/s%2F1", "", ``, src},
		{"update source", 200, src,
			func(c *audact.Client) (any, error) {
				return c.KnowledgeBases.UpdateSource(ctx, "k1", "s/1", &audact.UpdateKBSourceRequest{Title: &title})
			},
			"PATCH", "/v1/knowledge-bases/k1/sources/s%2F1", "", `{"title":"T"}`, src},
		{"reingest source", 202, src,
			func(c *audact.Client) (any, error) { return c.KnowledgeBases.ReingestSource(ctx, "k/1", "s/1") },
			"POST", "/v1/knowledge-bases/k%2F1/sources/s%2F1/reingest", "", ``, src},
		{"delete source", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.KnowledgeBases.DeleteSource(ctx, "k/1", "s/1") },
			"DELETE", "/v1/knowledge-bases/k%2F1/sources/s%2F1", "", ``, ``},
		{"list agent links", 200, `{"items":[` + link + `],"total":1}`,
			func(c *audact.Client) (any, error) { return c.KnowledgeBases.ListAgentLinks(ctx, "k/1") },
			"GET", "/v1/knowledge-bases/k%2F1/agents", "", ``, `{"items":[` + link + `],"total":1}`},
		{"list links by agent", 200, `{"items":[` + link + `],"total":1}`,
			func(c *audact.Client) (any, error) { return c.KnowledgeBases.ListLinksByAgent(ctx, "a/1") },
			"GET", "/v1/knowledge-bases/agent-links", "agent_id=a%2F1", ``, `{"items":[` + link + `],"total":1}`},
		{"link agent defaults sends no body", 200, link,
			func(c *audact.Client) (any, error) { return c.KnowledgeBases.LinkAgent(ctx, "k/1", "a/1", nil) },
			"PUT", "/v1/knowledge-bases/k%2F1/agents/a%2F1", "", ``, link},
		{"link agent explicit zero threshold", 200, link,
			func(c *audact.Client) (any, error) {
				return c.KnowledgeBases.LinkAgent(ctx, "k1", "a1", &audact.LinkKBAgentRequest{TopK: &topK, SimilarityThreshold: &zero, RetrievalInstruction: &empty})
			},
			"PUT", "/v1/knowledge-bases/k1/agents/a1", "", `{"top_k":5,"similarity_threshold":0,"retrieval_instruction":""}`, link},
		{"update agent link", 200, link,
			func(c *audact.Client) (any, error) {
				return c.KnowledgeBases.UpdateAgentLink(ctx, "k/1", "a/1", &audact.UpdateKBAgentLinkRequest{SimilarityThreshold: &zero})
			},
			"PATCH", "/v1/knowledge-bases/k%2F1/agents/a%2F1", "", `{"similarity_threshold":0}`, link},
		{"unlink agent", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.KnowledgeBases.UnlinkAgent(ctx, "k/1", "a/1") },
			"DELETE", "/v1/knowledge-bases/k%2F1/agents/a%2F1", "", ``, ``},
		{"list exclusions", 200, `{"items":[` + excl + `],"total":1}`,
			func(c *audact.Client) (any, error) { return c.KnowledgeBases.ListExclusions(ctx, "k/1") },
			"GET", "/v1/knowledge-bases/k%2F1/exclusions", "", ``, `{"items":[` + excl + `],"total":1}`},
		{"add exclusion", 201, excl,
			func(c *audact.Client) (any, error) {
				return c.KnowledgeBases.AddExclusion(ctx, "k/1", &audact.AddKBExclusionRequest{CrawlPath: "/admin"})
			},
			"POST", "/v1/knowledge-bases/k%2F1/exclusions", "", `{"crawl_path":"/admin"}`, excl},
		{"delete exclusion", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.KnowledgeBases.DeleteExclusion(ctx, "k/1", "e/1") },
			"DELETE", "/v1/knowledge-bases/k%2F1/exclusions/e%2F1", "", ``, ``},
		{"test retrieval question", 200, `{"condensed_query":"refund","empty":false,"hits":[{"kb_id":"k1","source_id":"s1","text":"30 days","similarity":0.87}]}`,
			func(c *audact.Client) (any, error) {
				return c.KnowledgeBases.TestRetrieval(ctx, "k/1", &audact.TestKBRetrievalRequest{Question: "refund?", TopK: &topK, Threshold: &zero})
			},
			"POST", "/v1/knowledge-bases/k%2F1/test-retrieval", "", `{"question":"refund?","top_k":5,"threshold":0}`,
			`{"condensed_query":"refund","hits":[{"kb_id":"k1","source_id":"s1","text":"30 days","similarity":0.87}]}`},
		{"test retrieval turns empty result", 200, `{"condensed_query":"","empty":true,"hits":[]}`,
			func(c *audact.Client) (any, error) {
				return c.KnowledgeBases.TestRetrieval(ctx, "k1", &audact.TestKBRetrievalRequest{Turns: []audact.KBRetrievalTurn{{Role: "user", Content: "hi"}}})
			},
			"POST", "/v1/knowledge-bases/k1/test-retrieval", "", `{"turns":[{"role":"user","content":"hi"}]}`, `{"empty":true}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, got := newAgentsServer(t, tc.status, tc.resp)
			out, err := tc.call(c)
			if err != nil {
				t.Fatal(err)
			}
			if got.method != tc.wantM || got.path != tc.wantPath || got.rawQuery != tc.wantQuery {
				t.Errorf("request = %s %s ?%s, want %s %s ?%s", got.method, got.path, got.rawQuery, tc.wantM, tc.wantPath, tc.wantQuery)
			}
			if got.auth != "Bearer key123" {
				t.Errorf("auth = %q", got.auth)
			}
			if tc.wantBody == "" {
				if len(got.body) != 0 {
					t.Errorf("unexpected body %s", got.body)
				}
			} else {
				if got.ctype != "application/json" {
					t.Errorf("content-type = %q", got.ctype)
				}
				if !jsonEqual(t, got.body, tc.wantBody) {
					t.Errorf("body = %s, want %s", got.body, tc.wantBody)
				}
			}
			if tc.wantOut != "" {
				b, _ := json.Marshal(out)
				if !jsonEqual(t, b, tc.wantOut) {
					t.Errorf("out = %s, want %s", b, tc.wantOut)
				}
			}
		})
	}
}

func TestKnowledgeBasesTypedDecode(t *testing.T) {
	c, _ := newAgentsServer(t, 200, `{"source_id":"s1","type":"text","status":"ready","ingest_version":2}`)
	out, err := c.KnowledgeBases.GetSource(context.Background(), "k1", "s1")
	if err != nil {
		t.Fatal(err)
	}
	if out.Type != audact.KBSourceTypeText || out.Status != audact.KBSourceStatusReady || out.IngestVersion != 2 {
		t.Fatalf("out = %+v", out)
	}
}

func TestKnowledgeBasesAddFileSource(t *testing.T) {
	c, got := newAgentsServer(t, 202, `{"source_id":"s9","type":"file","uri":"faq.pdf","status":"pending"}`)
	out, err := c.KnowledgeBases.AddFileSource(context.Background(), "k/1", &audact.AddKBFileSourceRequest{
		File: strings.NewReader("%PDF-1.7"), Filename: "faq.pdf", Config: map[string]any{"auto_refresh": true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.SourceID != "s9" || out.Type != audact.KBSourceTypeFile {
		t.Fatalf("out = %+v", out)
	}
	if got.method != "POST" || got.path != "/v1/knowledge-bases/k%2F1/sources" || got.auth != "Bearer key123" {
		t.Fatalf("request = %+v", got)
	}
	if !strings.HasPrefix(got.ctype, "multipart/form-data; boundary=") {
		t.Fatalf("content-type = %q", got.ctype)
	}
	r, _ := http.NewRequest("POST", "/", strings.NewReader(string(got.body)))
	r.Header.Set("Content-Type", got.ctype)
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		t.Fatal(err)
	}
	if v := r.FormValue("filename"); v != "faq.pdf" {
		t.Errorf("filename = %q", v)
	}
	if v := r.FormValue("config"); !jsonEqual(t, []byte(v), `{"auto_refresh":true}`) {
		t.Errorf("config = %q", v)
	}
	fh := r.MultipartForm.File["file"]
	if len(fh) != 1 || fh[0].Filename != "faq.pdf" {
		t.Fatalf("file part = %+v", fh)
	}
	fr, _ := fh[0].Open()
	b, _ := io.ReadAll(fr)
	if string(b) != "%PDF-1.7" {
		t.Errorf("file = %q", b)
	}
}

func TestKnowledgeBasesAddFileSourceOmitsUnsetFields(t *testing.T) {
	c, got := newAgentsServer(t, 202, `{}`)
	if _, err := c.KnowledgeBases.AddFileSource(context.Background(), "k1", &audact.AddKBFileSourceRequest{File: strings.NewReader("x")}); err != nil {
		t.Fatal(err)
	}
	r, _ := http.NewRequest("POST", "/", strings.NewReader(string(got.body)))
	r.Header.Set("Content-Type", got.ctype)
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"filename", "config"} {
		if _, ok := r.MultipartForm.Value[k]; ok {
			t.Errorf("%s should be omitted when unset", k)
		}
	}
	if fh := r.MultipartForm.File["file"]; len(fh) != 1 || fh[0].Filename != "upload" {
		t.Fatalf("file part = %+v", fh)
	}
}

func TestKnowledgeBasesAddFileSourceValidation(t *testing.T) {
	c, err := audact.NewClient("k", audact.WithBaseURL("http://127.0.0.1:1"))
	if err != nil {
		t.Fatal(err)
	}
	for _, req := range []*audact.AddKBFileSourceRequest{nil, {Filename: "x"}} {
		if _, err := c.KnowledgeBases.AddFileSource(context.Background(), "k1", req); err == nil || !strings.Contains(err.Error(), "file is required") {
			t.Errorf("req %+v: err = %v", req, err)
		}
	}
}

func TestKnowledgeBasesErrors(t *testing.T) {
	ctx := context.Background()
	c, _ := newAgentsServer(t, 404, `{"error":{"code":"NOT_FOUND","message":"knowledge base not found"}}`)
	if _, err := c.KnowledgeBases.Get(ctx, "x"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("Get err = %v", err)
	}
	if err := c.KnowledgeBases.DeleteSource(ctx, "k", "x"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("DeleteSource err = %v", err)
	}
	if err := c.KnowledgeBases.UnlinkAgent(ctx, "k", "a"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("UnlinkAgent err = %v", err)
	}

	c, _ = newAgentsServer(t, 400, `{"error":{"code":"INVALID_REQUEST","message":"provide exactly one of crawl_path or url"}}`)
	_, err := c.KnowledgeBases.AddExclusion(ctx, "k", &audact.AddKBExclusionRequest{})
	var apiErr *audact.APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 400 || apiErr.Code != "INVALID_REQUEST" {
		t.Fatalf("AddExclusion err = %v", err)
	}
}

func TestKBEnumTypes(t *testing.T) {
	fields := []struct {
		typ   reflect.Type
		field string
		want  reflect.Type
	}{
		{reflect.TypeOf(audact.KnowledgeBase{}), "RobotsPolicy", reflect.TypeOf(audact.KBRobotsPolicy(""))},
		{reflect.TypeOf(audact.CreateKnowledgeBaseRequest{}), "RobotsPolicy", reflect.TypeOf(audact.KBRobotsPolicy(""))},
		{reflect.TypeOf(audact.KBRetrievalTurn{}), "Role", reflect.TypeOf(audact.KBRetrievalTurnRole(""))},
	}
	for _, tc := range fields {
		t.Run(tc.typ.Name()+"."+tc.field, func(t *testing.T) {
			f, ok := tc.typ.FieldByName(tc.field)
			if !ok {
				t.Fatalf("%s.%s not found", tc.typ.Name(), tc.field)
			}
			if f.Type != tc.want {
				t.Errorf("%s.%s type = %v, want %v", tc.typ.Name(), tc.field, f.Type, tc.want)
			}
		})
	}

	decode := []struct {
		body string
		want audact.KBRobotsPolicy
	}{
		{`{"robots_policy":"honor"}`, audact.KBRobotsPolicyHonor},
		{`{"robots_policy":"override"}`, audact.KBRobotsPolicyOverride},
	}
	for _, tc := range decode {
		t.Run(tc.body, func(t *testing.T) {
			var kb audact.KnowledgeBase
			if err := json.Unmarshal([]byte(tc.body), &kb); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if kb.RobotsPolicy != tc.want {
				t.Errorf("RobotsPolicy = %q, want %q", kb.RobotsPolicy, tc.want)
			}
		})
	}
}

// TestKnowledgeBasesSourceIngestionOutcome covers the ingestion outcome fields: the
// latest successful ingestion's chunk count, crawl coverage and finish time.
// A zero chunk count is reported, not dropped; a source that has never
// finished an ingestion carries none of the three.
func TestKnowledgeBasesSourceIngestionOutcome(t *testing.T) {
	body := `{"source_id":"s1","type":"url","status":"ready","chunk_count":0,"crawl_coverage":{"discovered":12,"stored":10,"excluded":1,"failed":1},"last_ingested_at":"2026-06-22T12:05:00Z"}`
	c, _ := newAgentsServer(t, 200, body)
	out, err := c.KnowledgeBases.GetSource(context.Background(), "k1", "s1")
	if err != nil {
		t.Fatal(err)
	}
	if out.ChunkCount == nil || *out.ChunkCount != 0 {
		t.Fatalf("ChunkCount = %v, want pointer to 0", out.ChunkCount)
	}
	want := audact.KBCrawlCoverage{Discovered: 12, Stored: 10, Excluded: 1, Failed: 1}
	if out.CrawlCoverage == nil || *out.CrawlCoverage != want {
		t.Fatalf("CrawlCoverage = %+v, want %+v", out.CrawlCoverage, want)
	}
	if out.LastIngestedAt != "2026-06-22T12:05:00Z" {
		t.Fatalf("LastIngestedAt = %q", out.LastIngestedAt)
	}
	b, err := json.Marshal(out.CrawlCoverage)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"discovered":12,"excluded":1,"failed":1,"stored":10}` {
		t.Fatalf("coverage json = %s", b)
	}

	c2, _ := newAgentsServer(t, 200, `{"source_id":"s2","type":"url","status":"pending"}`)
	fresh, err := c2.KnowledgeBases.GetSource(context.Background(), "k1", "s2")
	if err != nil {
		t.Fatal(err)
	}
	if fresh.ChunkCount != nil || fresh.CrawlCoverage != nil || fresh.LastIngestedAt != "" {
		t.Fatalf("never-ingested source = %+v, want no ingestion outcome", fresh)
	}
}
