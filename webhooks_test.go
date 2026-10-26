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

func TestWebhooksService(t *testing.T) {
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
		{"list", 200, `{"items":[{"id":"w1","url":"https://x"}],"has_more":true,"next_cursor":"n"}`,
			func(c *audact.Client) (any, error) { return c.Webhooks.List(context.Background(), nil) },
			"GET", "/v1/webhooks", "", ``, `{"items":[{"id":"w1","url":"https://x"}],"has_more":true,"next_cursor":"n"}`},
		{"list params", 200, `{}`,
			func(c *audact.Client) (any, error) {
				return c.Webhooks.List(context.Background(), &audact.ListWebhooksParams{Limit: 2, Cursor: "c"})
			},
			"GET", "/v1/webhooks", "cursor=c&limit=2", ``, `{}`},
		{"create", 201, `{"id":"w1","signing_secret":"s"}`,
			func(c *audact.Client) (any, error) {
				return c.Webhooks.Create(context.Background(), &audact.CreateWebhookRequest{URL: "https://x", Events: []audact.WebhookEventType{audact.WebhookEventCallEnded}}, "")
			},
			"POST", "/v1/webhooks", "", `{"url":"https://x","events":["call_ended"]}`, `{"id":"w1","signing_secret":"s"}`},
		{"get escapes", 200, `{"id":"a/b"}`,
			func(c *audact.Client) (any, error) { return c.Webhooks.Get(context.Background(), "a/b") },
			"GET", "/v1/webhooks/a%2Fb", "", ``, `{"id":"a/b"}`},
		{"update omits unset", 200, `{"id":"w1"}`,
			func(c *audact.Client) (any, error) {
				return c.Webhooks.Update(context.Background(), "w1", &audact.UpdateWebhookRequest{IsActive: new(bool)})
			},
			"PUT", "/v1/webhooks/w1", "", `{"is_active":false}`, `{"id":"w1"}`},
		{"delete 204", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.Webhooks.Delete(context.Background(), "w1") },
			"DELETE", "/v1/webhooks/w1", "", ``, ``},
		{"deliveries", 200, `{"deliveries":[{"id":"d1","status":"dead_letter","latency_ms":1.5}],"total":1}`,
			func(c *audact.Client) (any, error) {
				return c.Webhooks.ListDeliveries(context.Background(), "w1", &audact.ListWebhookDeliveriesParams{Limit: 5, Offset: 10})
			},
			"GET", "/v1/webhooks/w1/deliveries", "limit=5&offset=10", ``, `{"deliveries":[{"id":"d1","status":"dead_letter","latency_ms":1.5}],"total":1}`},
		{"deliveries nil params", 200, `{"total":2}`,
			func(c *audact.Client) (any, error) { return c.Webhooks.ListDeliveries(context.Background(), "w1", nil) },
			"GET", "/v1/webhooks/w1/deliveries", "", ``, `{"total":2}`},
		{"test", 200, `{"success":true,"status_code":200}`,
			func(c *audact.Client) (any, error) { return c.Webhooks.Test(context.Background(), "w1") },
			"POST", "/v1/webhooks/w1/test", "", ``, `{"success":true,"status_code":200}`},
		{"resend", 200, `{"delivery_id":"d1","requeued":true}`,
			func(c *audact.Client) (any, error) { return c.Webhooks.ResendDelivery(context.Background(), "d1") },
			"POST", "/v1/webhooks/deliveries/d1/resend", "", ``, `{"delivery_id":"d1","requeued":true}`},
		{"event types", 200, `{"event_types":[{"name":"call_ended"}]}`,
			func(c *audact.Client) (any, error) { return c.Webhooks.ListEventTypes(context.Background()) },
			"GET", "/v1/webhooks/event-types", "", ``, `{"event_types":[{"name":"call_ended"}]}`},
		{"event schemas", 200, `{"event_type":"call_ended","example_payload":{"a":1}}`,
			func(c *audact.Client) (any, error) {
				return c.Webhooks.EventSchemas(context.Background(), "call_ended")
			},
			"GET", "/v1/webhooks/event-schemas", "type=call_ended", ``, `{"event_type":"call_ended","example_payload":{"a":1}}`},
		{"rotate", 200, `{"rotated":true,"signing_secret":"s2"}`,
			func(c *audact.Client) (any, error) { return c.Webhooks.RotateSecret(context.Background()) },
			"POST", "/v1/webhooks/secrets/rotate", "", ``, `{"rotated":true,"signing_secret":"s2"}`},
		{"secret status", 200, `{"in_rotation":true}`,
			func(c *audact.Client) (any, error) { return c.Webhooks.SecretStatus(context.Background()) },
			"GET", "/v1/webhooks/secrets/status", "", ``, `{"in_rotation":true}`},
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

func TestWebhooksCreateIdempotencyKey(t *testing.T) {
	var hdr string
	c, _ := newHeaderServer(t, &hdr, "Idempotency-Key")
	if _, err := c.Webhooks.Create(context.Background(), nil, "idem-1"); err != nil {
		t.Fatal(err)
	}
	if hdr != "idem-1" {
		t.Errorf("Idempotency-Key = %q", hdr)
	}
}

func TestWebhooksNotFound(t *testing.T) {
	c, _ := newAgentsServer(t, 404, `{"error":{"code":"not_found","message":"nope","request_id":"req_123"}}`)
	_, err := c.Webhooks.Get(context.Background(), "x")
	if !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	var apiErr *audact.APIError
	if !errors.As(err, &apiErr) || apiErr.RequestID != "req_123" {
		t.Fatalf("apiErr = %+v", apiErr)
	}
}

func newHeaderServer(t *testing.T, dst *string, header string) (*audact.Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*dst = r.Header.Get(header)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)
	c, err := audact.NewClient("key123", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	return c, srv
}
