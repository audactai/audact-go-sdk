package audact_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestIntegrationsService(t *testing.T) {
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
		{"list", 200, `{"connections":[{"connection_id":"c1","health":"connected"}]}`,
			func(c *audact.Client) (any, error) { return c.Integrations.List(context.Background(), nil) },
			"GET", "/v1/integrations", "", ``, `{"connections":[{"connection_id":"c1","health":"connected"}]}`},
		{"list provider", 200, `{}`,
			func(c *audact.Client) (any, error) {
				return c.Integrations.List(context.Background(), &audact.ListConnectionsParams{Provider: "cal_com"})
			},
			"GET", "/v1/integrations", "provider=cal_com", ``, `{}`},
		{"create", 201, `{"connection_id":"c1","status":"active"}`,
			func(c *audact.Client) (any, error) {
				return c.Integrations.Create(context.Background(), &audact.CreateConnectionRequest{Provider: "cal_com", AuthMethod: "api_key", APIKey: "k", Name: "n"})
			},
			"POST", "/v1/integrations", "", `{"api_key":"k","auth_method":"api_key","name":"n","provider":"cal_com"}`, `{"connection_id":"c1","status":"active"}`},
		{"get", 200, `{"connection_id":"c/1"}`,
			func(c *audact.Client) (any, error) { return c.Integrations.Get(context.Background(), "c/1") },
			"GET", "/v1/integrations/c%2F1", "", ``, `{"connection_id":"c/1"}`},
		{"delete", 200, `{"deleted":true}`,
			func(c *audact.Client) (any, error) { return c.Integrations.Delete(context.Background(), "c1", false) },
			"DELETE", "/v1/integrations/c1", "", ``, `{"deleted":true}`},
		{"delete force", 200, `{"deleted":true}`,
			func(c *audact.Client) (any, error) { return c.Integrations.Delete(context.Background(), "c1", true) },
			"DELETE", "/v1/integrations/c1", "force=true", ``, `{"deleted":true}`},
		{"event types", 200, `{"event_types":[{"id":1,"slug":"s","title":"t","lengthInMinutes":30}]}`,
			func(c *audact.Client) (any, error) { return c.Integrations.ListEventTypes(context.Background(), "c1") },
			"GET", "/v1/integrations/c1/event-types", "", ``, `{"event_types":[{"id":1,"slug":"s","title":"t","lengthInMinutes":30}]}`},
		{"test", 200, `{"connection_id":"c1","healthy":true,"event_type_count":2}`,
			func(c *audact.Client) (any, error) { return c.Integrations.Test(context.Background(), "c1") },
			"POST", "/v1/integrations/c1/test", "", ``, `{"connection_id":"c1","healthy":true,"event_type_count":2}`},
		{"oauth start", 200, `{"authorize_url":"https://cal.com/auth"}`,
			func(c *audact.Client) (any, error) { return c.Integrations.CalcomOAuthStart(context.Background(), nil) },
			"GET", "/v1/integrations/cal_com/oauth/start", "", ``, `{"authorize_url":"https://cal.com/auth"}`},
		{"oauth start params", 200, `{"authorize_url":"u"}`,
			func(c *audact.Client) (any, error) {
				return c.Integrations.CalcomOAuthStart(context.Background(), &audact.CalcomOAuthStartParams{Name: "n", CalUsername: "bob"})
			},
			"GET", "/v1/integrations/cal_com/oauth/start", "cal_username=bob&name=n", ``, `{"authorize_url":"u"}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, got := newAgentsServer(t, tc.status, tc.resp)
			out, err := tc.call(c)
			if err != nil {
				t.Fatal(err)
			}
			checkSeen(t, got, tc.wantM, tc.wantPath, tc.wantQuery, "Bearer key123", tc.wantBody)
			b, _ := json.Marshal(out)
			if !jsonEqual(t, b, tc.wantOut) {
				t.Errorf("out = %s, want %s", b, tc.wantOut)
			}
		})
	}
}

func checkSeen(t *testing.T, got *seen, method, path, query, auth, body string) {
	t.Helper()
	if got.method != method || got.path != path || got.rawQuery != query {
		t.Errorf("request = %s %s ?%s, want %s %s ?%s", got.method, got.path, got.rawQuery, method, path, query)
	}
	if got.auth != auth {
		t.Errorf("auth = %q, want %q", got.auth, auth)
	}
	if body == "" {
		if len(got.body) != 0 {
			t.Errorf("unexpected body %s", got.body)
		}
	} else if !jsonEqual(t, got.body, body) {
		t.Errorf("body = %s, want %s", got.body, body)
	}
}

func TestIntegrationsCalcomOAuthCallback(t *testing.T) {
	var gotAuth, gotQuery string
	redirects := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/dashboard" {
			redirects++
		}
		gotAuth = r.Header.Get("Authorization")
		gotQuery = r.URL.RawQuery
		if r.URL.Query().Get("state") == "bad" {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(`{"error":{"code":"bad_state","message":"bad state"}}`))
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}))
	t.Cleanup(srv.Close)
	c, err := audact.NewClient("key123", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	out, err := c.Integrations.CalcomOAuthCallback(context.Background(), audact.CalcomOAuthCallbackParams{State: "s1", Code: "c1"})
	if err != nil {
		t.Fatal(err)
	}
	if out.StatusCode != 302 || out.Location != "/dashboard" || redirects != 0 {
		t.Errorf("out = %+v redirects=%d", out, redirects)
	}
	if gotAuth != "" || gotQuery != "code=c1&state=s1" {
		t.Errorf("auth=%q query=%q", gotAuth, gotQuery)
	}
	_, err = c.Integrations.CalcomOAuthCallback(context.Background(), audact.CalcomOAuthCallbackParams{State: "bad"})
	var apiErr *audact.APIError
	if !errors.As(err, &apiErr) || apiErr.Code != "bad_state" {
		t.Fatalf("err = %v", err)
	}
}

func TestIntegrationsNotFound(t *testing.T) {
	c, _ := newAgentsServer(t, 404, `{"error":{"code":"not_found","message":"nope"}}`)
	if _, err := c.Integrations.Get(context.Background(), "x"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
}
