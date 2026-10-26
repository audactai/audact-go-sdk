package audact_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

type callsSeen struct {
	method, path, rawQuery, auth, idem string
	body                               string
	hits                               int
}

func newCallsServer(t *testing.T, status int, resp string) (*audact.Client, *callsSeen) {
	t.Helper()
	var got callsSeen
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		got = callsSeen{r.Method, r.URL.EscapedPath(), r.URL.RawQuery, r.Header.Get("Authorization"), r.Header.Get("Idempotency-Key"), string(b), got.hits + 1}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(resp))
	}))
	t.Cleanup(srv.Close)
	c, err := audact.NewClient("key123", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	return c, &got
}

func TestCallsService(t *testing.T) {
	yes := true
	ctx := context.Background()
	tests := []struct {
		name, resp                                              string
		status                                                  int
		call                                                    func(c *audact.Client) (string, error)
		wantM, wantPath, wantQuery, wantBody, wantIdem, wantOut string
	}{
		{"list params", `{"calls":[{"call_id":"c1"}],"has_more":true,"next_cursor":"n"}`, 200,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.List(ctx, &audact.ListCallsParams{DateFrom: "a", DateTo: "b", Status: "completed", Verdict: "BLOCK", Shadow: &yes, Country: "DE", ConfigID: "cfg", Cursor: "cu", Limit: 7})
				if err != nil {
					return "", err
				}
				return r.Calls[0].CallID + r.NextCursor, nil
			},
			"GET", "/v1/calls", "config_id=cfg&country=DE&cursor=cu&date_from=a&date_to=b&limit=7&shadow=true&status=completed&verdict=BLOCK", "", "", "c1n"},
		{"list nil", `{}`, 200,
			func(c *audact.Client) (string, error) { _, err := c.Calls.List(ctx, nil); return "", err },
			"GET", "/v1/calls", "", "", "", ""},
		{"get escaped", `{"call_id":"a/b","status":"completed"}`, 200,
			func(c *audact.Client) (string, error) { r, err := c.Calls.Get(ctx, "a/b"); return r.CallID, err },
			"GET", "/v1/calls/a%2Fb", "", "", "", "a/b"},
		{"summary", `{"call_id":"c1","event_count":3}`, 200,
			func(c *audact.Client) (string, error) { r, err := c.Calls.Summary(ctx, "c1"); return r.CallID, err },
			"GET", "/v1/calls/c1/summary", "", "", "", "c1"},
		{"get recording", `{"call_id":"c1","url":"https://x"}`, 200,
			func(c *audact.Client) (string, error) { r, err := c.Calls.GetRecording(ctx, "c1"); return r.URL, err },
			"GET", "/v1/calls/c1/recording", "", "", "", "https://x"},
		{"delete recording", ``, 204,
			func(c *audact.Client) (string, error) { return "", c.Calls.DeleteRecording(ctx, "c1") },
			"DELETE", "/v1/calls/c1/recording", "", "", "", ""},
		{"outbound idem", `{"call_id":"c9","room_name":"r"}`, 201,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.CreateOutbound(ctx, &audact.OutboundCallRequest{AgentID: "a", FromNumber: "+1", DestinationNumber: "+2"}, "idem-1")
				return r.CallID, err
			},
			"POST", "/v1/sip/calls/outbound", "", `{"agent_id":"a","destination_number":"+2","from_number":"+1"}`, "idem-1", "c9"},
		{"inbound no idem", `{"call_id":"c8","trunk_id":"t"}`, 200,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.CreateInbound(ctx, &audact.InboundCallRequest{CalledDID: "+1", CallerNumber: "+2"}, "")
				return r.CallID + r.TrunkID, err
			},
			"POST", "/v1/sip/calls/inbound", "", `{"called_did":"+1","caller_number":"+2"}`, "", "c8t"},
		{"browser nil req", `{"call_id":"c7","livekit_token":"tok"}`, 201,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.CreateBrowser(ctx, nil, "k")
				return r.LiveKitToken, err
			},
			"POST", "/v1/sip/calls/browser", "", `{}`, "k", "tok"},
		{"status", `{"call_id":"c1","status":"active"}`, 200,
			func(c *audact.Client) (string, error) { r, err := c.Calls.Status(ctx, "c1"); return r.Status, err },
			"GET", "/v1/sip/calls/c1/status", "", "", "", "active"},
		{"hangup", `{"call_id":"c1","status":"terminated"}`, 200,
			func(c *audact.Client) (string, error) { r, err := c.Calls.Hangup(ctx, "c1"); return r.Status, err },
			"POST", "/v1/sip/calls/c1/hangup", "", "", "", "terminated"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, got := newCallsServer(t, tt.status, tt.resp)
			out, err := tt.call(c)
			if err != nil {
				t.Fatal(err)
			}
			if got.method != tt.wantM || got.path != tt.wantPath || got.rawQuery != tt.wantQuery {
				t.Errorf("request = %s %s ?%s", got.method, got.path, got.rawQuery)
			}
			if got.auth != "Bearer key123" {
				t.Errorf("auth = %q", got.auth)
			}
			if got.body != tt.wantBody {
				t.Errorf("body = %q, want %q", got.body, tt.wantBody)
			}
			if got.idem != tt.wantIdem {
				t.Errorf("idempotency = %q, want %q", got.idem, tt.wantIdem)
			}
			if out != tt.wantOut {
				t.Errorf("out = %q, want %q", out, tt.wantOut)
			}
		})
	}
}

func TestCallsService_NotFound(t *testing.T) {
	c, _ := newCallsServer(t, 404, `{"error":{"code":"not_found","message":"call not found","request_id":"req_123"}}`)
	_, err := c.Calls.Get(context.Background(), "missing")
	if !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	var apiErr *audact.APIError
	if !errors.As(err, &apiErr) || apiErr.RequestID != "req_123" {
		t.Errorf("apiErr = %+v", apiErr)
	}
}
