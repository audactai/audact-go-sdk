package audact_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestNumbersBYO(t *testing.T) {
	ctx := context.Background()
	empty := []string{}
	tests := []struct {
		name      string
		status    int
		resp      string
		call      func(c *audact.Client) (any, error)
		wantM     string
		wantPath  string
		wantQuery string
		wantIdem  string
		wantBody  string
		wantOut   string
	}{
		{"config", 200, `{"signaling_ips":["203.0.113.10/32"],"sip_ingress_host":"sip.audact.ai"}`,
			func(c *audact.Client) (any, error) { return c.Numbers.BYOConfig(ctx) },
			"GET", "/v1/byo/config", "", "", "", `{"signaling_ips":["203.0.113.10/32"],"sip_ingress_host":"sip.audact.ai"}`},
		{"list nil", 200, `{"connections":[{"id":"c1","connection_health":"healthy"}],"next_cursor":"n2"}`,
			func(c *audact.Client) (any, error) { return c.Numbers.ListBYOConnections(ctx, nil) },
			"GET", "/v1/byo/connections", "", "", "", `{"connections":[{"connection_health":"healthy","id":"c1"}],"next_cursor":"n2"}`},
		{"list params", 200, `{}`,
			func(c *audact.Client) (any, error) {
				return c.Numbers.ListBYOConnections(ctx, &audact.ListBYOConnectionsParams{Limit: 5, Cursor: "c"})
			},
			"GET", "/v1/byo/connections", "cursor=c&limit=5", "", "", `{}`},
		{"create", 201, `{"id":"c1","has_credentials":true,"last_error_code":503}`,
			func(c *audact.Client) (any, error) {
				return c.Numbers.CreateBYOConnection(ctx, &audact.CreateBYOConnectionRequest{
					Name: "HQ", Provider: audact.NumberProviderTwilio, SIPHost: "h.example.com", SIPUsername: "u", SIPPassword: "p",
				}, "idem-1")
			},
			"POST", "/v1/byo/connections", "", "idem-1",
			`{"name":"HQ","provider":"twilio","sip_host":"h.example.com","sip_password":"p","sip_username":"u"}`,
			`{"has_credentials":true,"id":"c1","last_error_code":503}`},
		{"get escaped", 200, `{"id":"a b"}`,
			func(c *audact.Client) (any, error) { return c.Numbers.GetBYOConnection(ctx, "a b") },
			"GET", "/v1/byo/connections/a%20b", "", "", "", `{"id":"a b"}`},
		{"update partial", 200, `{"id":"c1","name":"new"}`,
			func(c *audact.Client) (any, error) {
				return c.Numbers.UpdateBYOConnection(ctx, "c1", &audact.UpdateBYOConnectionRequest{Name: "new"}, "")
			},
			"PUT", "/v1/byo/connections/c1", "", "", `{"name":"new"}`, `{"id":"c1","name":"new"}`},
		{"update clears allowlist", 200, `{"id":"c1"}`,
			func(c *audact.Client) (any, error) {
				return c.Numbers.UpdateBYOConnection(ctx, "c1", &audact.UpdateBYOConnectionRequest{InboundAllowedIPs: &empty, SIPTransport: audact.SIPTransportTLS}, "i2")
			},
			"PUT", "/v1/byo/connections/c1", "", "i2", `{"inbound_allowed_ips":[],"sip_transport":"tls"}`, `{"id":"c1"}`},
		{"update nil", 200, `{"id":"c1"}`,
			func(c *audact.Client) (any, error) { return c.Numbers.UpdateBYOConnection(ctx, "c1", nil, "") },
			"PUT", "/v1/byo/connections/c1", "", "", `{}`, `{"id":"c1"}`},
		{"delete", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.Numbers.DeleteBYOConnection(ctx, "c1") },
			"DELETE", "/v1/byo/connections/c1", "", "", "", ``},
		{"verify", 200, `{"id":"c1","connection_health":"unreachable"}`,
			func(c *audact.Client) (any, error) { return c.Numbers.VerifyBYOConnection(ctx, "c1") },
			"POST", "/v1/byo/connections/c1/verify", "", "", "", `{"connection_health":"unreachable","id":"c1"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, got := newNumbersServer(t, tt.status, tt.resp)
			out, err := tt.call(c)
			if err != nil {
				t.Fatal(err)
			}
			if got.method != tt.wantM || got.path != tt.wantPath || got.rawQuery != tt.wantQuery {
				t.Errorf("got %s %s ?%s", got.method, got.path, got.rawQuery)
			}
			if got.auth != "Bearer key123" {
				t.Errorf("auth = %q", got.auth)
			}
			if got.idem != tt.wantIdem {
				t.Errorf("idem = %q", got.idem)
			}
			if tt.wantBody != "" && string(got.body) != tt.wantBody {
				t.Errorf("body = %s", got.body)
			}
			if tt.wantBody == "" && len(got.body) != 0 {
				t.Errorf("unexpected body %s", got.body)
			}
			if tt.wantOut != "" {
				b, _ := json.Marshal(out)
				if string(b) != tt.wantOut {
					t.Errorf("out = %s, want %s", b, tt.wantOut)
				}
			}
		})
	}
}

func TestNumbersBYO_NotFound(t *testing.T) {
	c, _ := newNumbersServer(t, 404, `{"error":{"code":"not_found","message":"connection not found","request_id":"req_123"}}`)
	ctx := context.Background()
	if _, err := c.Numbers.VerifyBYOConnection(ctx, "x"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	if err := c.Numbers.DeleteBYOConnection(ctx, "x"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
}
