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

func TestTenantsService(t *testing.T) {
	ctx := context.Background()
	tr := true
	days := 30
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
		{"get", 200, `{"tenant_id":"t1","name":"Acme","plan":"pro"}`,
			func(c *audact.Client) (any, error) { return c.Tenants.Get(ctx, "t/1") },
			"GET", "/v1/tenants/t%2F1", "", ``, `{"tenant_id":"t1","name":"Acme","plan":"pro"}`},
		{"update", 200, `{"tenant_id":"t1","name":"New"}`,
			func(c *audact.Client) (any, error) {
				return c.Tenants.Update(ctx, "t1", &audact.UpdateTenantRequest{Name: "New"})
			},
			"PUT", "/v1/tenants/t1", "", `{"name":"New"}`, `{"tenant_id":"t1","name":"New"}`},
		{"usage", 200, `{"tenant_id":"t1","calls_total":3,"validator_checks":9}`,
			func(c *audact.Client) (any, error) { return c.Tenants.Usage(ctx, "t1", nil) },
			"GET", "/v1/tenants/t1/usage", "", ``, `{"tenant_id":"t1","calls_total":3,"validator_checks":9}`},
		{"usage params", 200, `{}`,
			func(c *audact.Client) (any, error) {
				return c.Tenants.Usage(ctx, "t1", &audact.GetTenantUsageParams{PeriodStart: "2026-10-01T00:00:00Z", PeriodEnd: "2026-10-31T00:00:00Z"})
			},
			"GET", "/v1/tenants/t1/usage", "period_end=2026-10-31T00%3A00%3A00Z&period_start=2026-10-01T00%3A00%3A00Z", ``, `{}`},
		{"list keys", 200, `{"keys":[{"api_key_id":"k1","scopes":["a"]}],"has_more":true,"next_cursor":"n"}`,
			func(c *audact.Client) (any, error) {
				return c.Tenants.ListKeys(ctx, "t1", &audact.ListTenantKeysParams{Cursor: "c", Limit: 2})
			},
			"GET", "/v1/tenants/t1/keys", "cursor=c&limit=2", ``,
			`{"keys":[{"api_key_id":"k1","scopes":["a"]}],"has_more":true,"next_cursor":"n"}`},
		{"list keys nil", 200, `{}`,
			func(c *audact.Client) (any, error) { return c.Tenants.ListKeys(ctx, "t1", nil) },
			"GET", "/v1/tenants/t1/keys", "", ``, `{}`},
		{"generate key", 201, `{"api_key_id":"k1","key":"secret"}`,
			func(c *audact.Client) (any, error) {
				return c.Tenants.GenerateKey(ctx, "t1", &audact.GenerateKeyRequest{Name: "ci", Role: "developer"})
			},
			"POST", "/v1/tenants/t1/keys", "", `{"name":"ci","role":"developer"}`, `{"api_key_id":"k1","key":"secret"}`},
		{"revoke key", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.Tenants.RevokeKey(ctx, "t1", "k/1") },
			"DELETE", "/v1/tenants/t1/keys/k%2F1", "", ``, ``},
		{"list configurations", 200, `{"configurations":[{"config_id":"c1","version":2}],"has_more":false}`,
			func(c *audact.Client) (any, error) {
				return c.Tenants.ListConfigurations(ctx, "t1", &audact.ListConfigurationsParams{Limit: 1})
			},
			"GET", "/v1/tenants/t1/configurations", "limit=1", ``, `{"configurations":[{"config_id":"c1","pii_encryption_enabled":false,"version":2}]}`},
		{"create configuration", 201, `{"config_id":"c1","name":"n","validator_modes":{"x":1}}`,
			func(c *audact.Client) (any, error) {
				return c.Tenants.CreateConfiguration(ctx, &audact.CreateConfigurationRequest{
					TenantID: "t1", Name: "n", RecordCalls: &tr, RecordingRetentionDays: &days,
					BlocklistPatterns: []audact.BlocklistPatternRequest{{ID: "b", Terms: []string{"x"}}},
					ValidatorModes:    json.RawMessage(`{"x":1}`),
				})
			},
			"POST", "/v1/configurations", "",
			`{"tenant_id":"t1","name":"n","record_calls":true,"recording_retention_days":30,"blocklist_patterns":[{"id":"b","terms":["x"]}],"validator_modes":{"x":1}}`,
			`{"config_id":"c1","name":"n","pii_encryption_enabled":false,"validator_modes":{"x":1}}`},
		{"get configuration", 200, `{"config_id":"c/1","end_call":{"closing_line":{"en":"bye"}}}`,
			func(c *audact.Client) (any, error) { return c.Tenants.GetConfiguration(ctx, "c/1") },
			"GET", "/v1/configurations/c%2F1", "", ``, `{"config_id":"c/1","pii_encryption_enabled":false,"end_call":{"closing_line":{"en":"bye"}}}`},
		{"get configuration version", 200, `{"config_id":"c/1","version":3,"validator_modes":{"x":1}}`,
			func(c *audact.Client) (any, error) { return c.Tenants.GetConfigurationVersion(ctx, "c/1", 3) },
			"GET", "/v1/configurations/c%2F1/versions/3", "", ``, `{"config_id":"c/1","pii_encryption_enabled":false,"version":3,"validator_modes":{"x":1}}`},
		{"update configuration omits unset", 200, `{"config_id":"c1"}`,
			func(c *audact.Client) (any, error) {
				return c.Tenants.UpdateConfiguration(ctx, "c1", &audact.UpdateConfigurationRequest{Name: "z"})
			},
			"PUT", "/v1/configurations/c1", "", `{"name":"z"}`, `{"config_id":"c1","pii_encryption_enabled":false}`},
		{"delete configuration", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.Tenants.DeleteConfiguration(ctx, "c1") },
			"DELETE", "/v1/configurations/c1", "", ``, ``},
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

func TestTenantsCreateSendsNoAuth(t *testing.T) {
	tests := []struct {
		name, idem string
	}{
		{"no idempotency key", ""},
		{"idempotency key", "idem-1"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotAuth, gotIdem, gotPath, gotMethod string
			var gotBody map[string]any
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotAuth, gotIdem = r.Header.Get("Authorization"), r.Header.Get("Idempotency-Key")
				gotPath, gotMethod = r.URL.Path, r.Method
				_ = json.NewDecoder(r.Body).Decode(&gotBody)
				w.WriteHeader(201)
				_, _ = w.Write([]byte(`{"tenant_id":"t1","api_key":{"key":"secret","key_id":"k1"}}`))
			}))
			defer srv.Close()
			c, err := audact.NewClient("key123", audact.WithBaseURL(srv.URL))
			if err != nil {
				t.Fatal(err)
			}
			out, err := c.Tenants.Create(context.Background(), &audact.CreateTenantRequest{Name: "Acme", Plan: "pro"}, tc.idem)
			if err != nil {
				t.Fatal(err)
			}
			if gotAuth != "" || gotIdem != tc.idem || gotPath != "/v1/tenants" || gotMethod != "POST" {
				t.Errorf("auth=%q idem=%q %s %s", gotAuth, gotIdem, gotMethod, gotPath)
			}
			if gotBody["name"] != "Acme" || gotBody["plan"] != "pro" || len(gotBody) != 2 {
				t.Errorf("body = %v", gotBody)
			}
			if out.TenantID != "t1" || out.APIKey == nil || out.APIKey.Key != "secret" {
				t.Errorf("out = %+v", out)
			}
		})
	}
}

func TestTenantsNotFound(t *testing.T) {
	c, _ := newAgentsServer(t, 404, `{"error":{"code":"not_found","message":"tenant not found"}}`)
	if _, err := c.Tenants.Get(context.Background(), "x"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	if err := c.Tenants.RevokeKey(context.Background(), "x", "y"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
}
