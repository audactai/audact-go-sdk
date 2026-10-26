package audact_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestToolsService(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		resp      string
		call      func(c *audact.Client) (any, error)
		wantM     string
		wantPath  string
		wantQuery string
		wantBody  string
		wantOut   string
	}{
		{"list", 200, `{"tools":[{"tool_id":"t1"}],"has_more":true,"cursor":"c"}`,
			func(c *audact.Client) (any, error) { return c.Tools.List(context.Background(), nil) },
			"GET", "/v1/tools", "", ``, `{"tools":[{"tool_id":"t1"}],"has_more":true,"cursor":"c"}`},
		{"list params", 200, `{}`,
			func(c *audact.Client) (any, error) {
				return c.Tools.List(context.Background(), &audact.ListToolsParams{Name: "n", Cursor: "c", Limit: 3})
			},
			"GET", "/v1/tools", "cursor=c&limit=3&name=n", ``, `{}`},
		{"create", 201, `{"tool_id":"t1"}`,
			func(c *audact.Client) (any, error) {
				return c.Tools.Create(context.Background(), &audact.CreateToolRequest{Name: "n", Method: "GET", URL: "https://x", QueryParams: []audact.ToolQueryParam{{Key: "k", Mode: "const", Value: "v"}}, Execution: &audact.ToolExecution{TimeoutMs: 5}}, "")
			},
			"POST", "/v1/tools", "", `{"name":"n","method":"GET","url":"https://x","query_params":[{"key":"k","mode":"const","value":"v"}],"execution":{"timeout_ms":5}}`, `{"tool_id":"t1"}`},
		{"get", 200, `{"tool_id":"t/1"}`,
			func(c *audact.Client) (any, error) { return c.Tools.Get(context.Background(), "t/1") },
			"GET", "/v1/tools/t%2F1", "", ``, `{"tool_id":"t/1"}`},
		{"replace", 200, `{"tool_id":"t1"}`,
			func(c *audact.Client) (any, error) {
				return c.Tools.Replace(context.Background(), "t1", &audact.CreateToolRequest{Name: "n"})
			},
			"PUT", "/v1/tools/t1", "", `{"name":"n"}`, `{"tool_id":"t1"}`},
		{"update omits unset", 200, `{"tool_id":"t1"}`,
			func(c *audact.Client) (any, error) {
				return c.Tools.Update(context.Background(), "t1", &audact.UpdateToolRequest{ArgsOnly: new(bool)})
			},
			"PATCH", "/v1/tools/t1", "", `{"args_only":false}`, `{"tool_id":"t1"}`},
		{"delete", 200, `{"deleted":true}`,
			func(c *audact.Client) (any, error) { return c.Tools.Delete(context.Background(), "t1", false) },
			"DELETE", "/v1/tools/t1", "", ``, `{"deleted":true}`},
		{"delete force", 200, `{"deleted":true}`,
			func(c *audact.Client) (any, error) { return c.Tools.Delete(context.Background(), "t1", true) },
			"DELETE", "/v1/tools/t1", "force=true", ``, `{"deleted":true}`},
		{"test", 200, `{"healthy":true,"checks_passed":["dns"]}`,
			func(c *audact.Client) (any, error) { return c.Tools.Test(context.Background(), "t1") },
			"POST", "/v1/tools/t1/test", "", ``, `{"healthy":true,"checks_passed":["dns"]}`},
		{"list for agent", 200, `{"items":[{"tool_id":"t1"}]}`,
			func(c *audact.Client) (any, error) {
				return c.Tools.ListForAgent(context.Background(), "a1", &audact.ListAgentToolsParams{Limit: 1, Cursor: "c"})
			},
			"GET", "/v1/agents/a1/tools", "cursor=c&limit=1", ``, `{"items":[{"tool_id":"t1"}]}`},
		{"attach", 201, `{"attached":true,"agent_id":"a1","tool_id":"t1"}`,
			func(c *audact.Client) (any, error) {
				return c.Tools.AttachToAgent(context.Background(), "a1", &audact.AttachToolRequest{ToolID: "t1"})
			},
			"POST", "/v1/agents/a1/tools", "", `{"tool_id":"t1"}`, `{"attached":true,"agent_id":"a1","tool_id":"t1"}`},
		{"detach", 200, `{"detached":true}`,
			func(c *audact.Client) (any, error) { return c.Tools.DetachFromAgent(context.Background(), "a/1", "t1") },
			"DELETE", "/v1/agents/a%2F1/tools/t1", "", ``, `{"detached":true}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, got := newAgentsServer(t, tc.status, tc.resp)
			out, err := tc.call(c)
			if err != nil {
				t.Fatal(err)
			}
			if got.method != tc.wantM || got.path != tc.wantPath || got.rawQuery != tc.wantQuery {
				t.Errorf("request = %s %s ?%s", got.method, got.path, got.rawQuery)
			}
			if got.auth != "Bearer key123" {
				t.Errorf("auth = %q", got.auth)
			}
			if tc.wantBody == "" {
				if len(got.body) != 0 {
					t.Errorf("unexpected body %s", got.body)
				}
			} else if !jsonEqual(t, got.body, tc.wantBody) {
				t.Errorf("body = %s, want %s", got.body, tc.wantBody)
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

func TestToolsCreateIdempotencyKey(t *testing.T) {
	var hdr string
	c, _ := newHeaderServer(t, &hdr, "Idempotency-Key")
	if _, err := c.Tools.Create(context.Background(), nil, "idem-2"); err != nil {
		t.Fatal(err)
	}
	if hdr != "idem-2" {
		t.Errorf("Idempotency-Key = %q", hdr)
	}
}

func TestToolsNotFound(t *testing.T) {
	c, _ := newAgentsServer(t, 404, `{"error":{"code":"not_found","message":"nope"}}`)
	if _, err := c.Tools.Get(context.Background(), "x"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	if _, err := c.Tools.DetachFromAgent(context.Background(), "a", "x"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
}
