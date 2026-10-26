package audact_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

type numbersSeen struct {
	method, path, rawQuery, auth, idem string
	body                               []byte
}

func newNumbersServer(t *testing.T, status int, resp string) (*audact.Client, *numbersSeen) {
	t.Helper()
	var got numbersSeen
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		got = numbersSeen{r.Method, r.URL.EscapedPath(), r.URL.RawQuery, r.Header.Get("Authorization"), r.Header.Get("Idempotency-Key"), b}
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

func TestNumbersService(t *testing.T) {
	ctx := context.Background()
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
		{"list nil", 200, `{"numbers":[{"id":"n1","did":"+31600000000","status":"available","origin":"audact"}],"total_count":1,"next_cursor":"c2"}`,
			func(c *audact.Client) (any, error) { return c.Numbers.List(ctx, nil) },
			"GET", "/v1/numbers", "", "", "", `{"next_cursor":"c2","numbers":[{"did":"+31600000000","id":"n1","origin":"audact","status":"available"}],"total_count":1}`},
		{"list params", 200, `{}`,
			func(c *audact.Client) (any, error) {
				return c.Numbers.List(ctx, &audact.ListNumbersParams{Country: "NL", Status: audact.NumberStatusAssigned, Limit: 10, Cursor: "c1"})
			},
			"GET", "/v1/numbers", "country=NL&cursor=c1&limit=10&status=assigned", "", "", `{}`},
		{"create with idem", 201, `{"id":"n1"}`,
			func(c *audact.Client) (any, error) {
				return c.Numbers.Create(ctx, &audact.CreateNumberRequest{DID: "+31600000000", CountryCode: "NL", Origin: audact.NumberOriginBYO, Connection: &audact.BYOConnectionRequest{Provider: audact.NumberProviderTwilio, SIPHost: "h"}}, "idem-1")
			},
			"POST", "/v1/numbers", "", "idem-1", `{"connection":{"provider":"twilio","sip_host":"h"},"country_code":"NL","did":"+31600000000","origin":"byo"}`, `{"id":"n1"}`},
		{"create nil no idem", 201, `{}`,
			func(c *audact.Client) (any, error) { return c.Numbers.Create(ctx, nil, "") },
			"POST", "/v1/numbers", "", "", `{}`, `{}`},
		{"list available", 200, `{"numbers":[{"phone_number":"+3161","country_code":"NL","phone_number_type":"mobile"}]}`,
			func(c *audact.Client) (any, error) {
				return c.Numbers.ListAvailable(ctx, &audact.ListAvailableNumbersParams{Country: "NL", NumberType: "mobile"})
			},
			"GET", "/v1/numbers/available", "country=NL&number_type=mobile", "", "", `{"numbers":[{"country_code":"NL","phone_number":"+3161","phone_number_type":"mobile"}]}`},
		{"lookup escapes query", 200, `{"id":"n1"}`,
			func(c *audact.Client) (any, error) { return c.Numbers.Lookup(ctx, "+31600000000") },
			"GET", "/v1/numbers/lookup", "did=%2B31600000000", "", "", `{"id":"n1"}`},
		{"list orders", 200, `{"orders":[{"id":"o1","status":"success"}]}`,
			func(c *audact.Client) (any, error) { return c.Numbers.ListOrders(ctx) },
			"GET", "/v1/numbers/orders", "", "", "", `{"orders":[{"id":"o1","status":"success"}]}`},
		{"create order", 201, `{"id":"o1","status":"pending"}`,
			func(c *audact.Client) (any, error) {
				return c.Numbers.CreateOrder(ctx, &audact.CreateNumberOrderRequest{CountryCode: "NL", NumberType: "mobile", SetID: "s1"}, "idem-2")
			},
			"POST", "/v1/numbers/orders", "", "idem-2", `{"country_code":"NL","number_type":"mobile","set_id":"s1"}`, `{"id":"o1","status":"pending"}`},
		{"get order escapes", 200, `{"id":"o/1"}`,
			func(c *audact.Client) (any, error) { return c.Numbers.GetOrder(ctx, "o/1") },
			"GET", "/v1/numbers/orders/o%2F1", "", "", "", `{"id":"o/1"}`},
		{"delete 204", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.Numbers.Delete(ctx, "n/1") },
			"DELETE", "/v1/numbers/n%2F1", "", "", "", ""},
		{"assign", 204, ``,
			func(c *audact.Client) (any, error) {
				return nil, c.Numbers.Assign(ctx, "n1", &audact.AssignNumberRequest{AgentID: "a1"}, "idem-3")
			},
			"POST", "/v1/numbers/n1/assign", "", "idem-3", `{"agent_id":"a1"}`, ""},
		{"assign 200 empty body", 200, ``,
			func(c *audact.Client) (any, error) { return nil, c.Numbers.Assign(ctx, "n1", nil, "") },
			"POST", "/v1/numbers/n1/assign", "", "", `{}`, ""},
		{"release", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.Numbers.Release(ctx, "n1", "idem-4") },
			"DELETE", "/v1/numbers/n1/release", "", "idem-4", "", ""},
		{"resume", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.Numbers.Resume(ctx, "n1", "") },
			"POST", "/v1/numbers/n1/resume", "", "", "", ""},
		{"suspend", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.Numbers.Suspend(ctx, "n1", "idem-5") },
			"POST", "/v1/numbers/n1/suspend", "", "idem-5", "", ""},
		{"set trunk", 200, `{"id":"n1","trunk_id":"t1"}`,
			func(c *audact.Client) (any, error) {
				return c.Numbers.SetTrunk(ctx, "n1", &audact.SetNumberTrunkRequest{TrunkID: "t1"})
			},
			"PUT", "/v1/numbers/n1/trunk", "", "", `{"trunk_id":"t1"}`, `{"id":"n1","trunk_id":"t1"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, got := newNumbersServer(t, tt.status, tt.resp)
			out, err := tt.call(c)
			if err != nil {
				t.Fatal(err)
			}
			if got.method != tt.wantM || got.path != tt.wantPath || got.rawQuery != tt.wantQuery {
				t.Errorf("request = %s %s ?%s, want %s %s ?%s", got.method, got.path, got.rawQuery, tt.wantM, tt.wantPath, tt.wantQuery)
			}
			if got.auth != "Bearer key123" {
				t.Errorf("auth = %q", got.auth)
			}
			if got.idem != tt.wantIdem {
				t.Errorf("idempotency = %q, want %q", got.idem, tt.wantIdem)
			}
			if tt.wantBody == "" {
				if len(got.body) != 0 {
					t.Errorf("body = %q, want none", got.body)
				}
			} else if string(got.body) != tt.wantBody {
				t.Errorf("body = %s, want %s", got.body, tt.wantBody)
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

func TestNumbersService_NotFound(t *testing.T) {
	c, _ := newNumbersServer(t, 404, `{"error":{"code":"not_found","message":"number not found","request_id":"req_123"}}`)
	_, err := c.Numbers.Lookup(context.Background(), "+1")
	if !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	if err := c.Numbers.Delete(context.Background(), "x"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("delete err = %v", err)
	}
}
