package audact

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func doErr(t *testing.T, status int, hdr map[string]string, body string) error {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range hdr {
			w.Header().Set(k, v)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	c, err := NewClient("k", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.t.do(context.Background(), http.MethodGet, "/x", authAPIKey, nil, nil, nil)
	return err
}

func TestParseError(t *testing.T) {
	tests := []struct {
		name             string
		status           int
		hdr              map[string]string
		body             string
		code, msg, reqID string
		sentinel         error
	}{
		{"nested", 404, nil, `{"error":{"code":"not_found","message":"call not found","hint":"h","details":{"a":1},"request_id":"req_123"}}`, "not_found", "call not found", "req_123", ErrNotFound},
		{"flat", 400, nil, `{"code":"bad","message":"nope","request_id":"r1"}`, "bad", "nope", "r1", nil},
		{"non-json", 502, nil, `<html>`, "", "Bad Gateway", "", nil},
		{"empty", 500, nil, ``, "", "Internal Server Error", "", nil},
		{"header fallback", 401, map[string]string{"X-Request-ID": "hdr"}, `{"error":{"code":"u","message":"m"}}`, "u", "m", "hdr", ErrUnauthorized},
		{"body wins", 403, map[string]string{"X-Request-ID": "hdr"}, `{"error":{"message":"m","request_id":"body"}}`, "", "m", "body", ErrForbidden},
		{"conflict", 409, nil, `{}`, "", "Conflict", "", ErrConflict},
		{"500 no sentinel", 500, nil, `{"error":{"message":"x"}}`, "", "x", "", nil},
	}
	sentinels := []error{ErrUnauthorized, ErrForbidden, ErrNotFound, ErrConflict, ErrRateLimited}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := doErr(t, tc.status, tc.hdr, tc.body)
			var ae *APIError
			if !errors.As(err, &ae) {
				t.Fatalf("not APIError: %v", err)
			}
			if ae.StatusCode != tc.status || ae.Code != tc.code || ae.Message != tc.msg || ae.RequestID != tc.reqID {
				t.Fatalf("got %+v", ae)
			}
			if string(ae.Body) != tc.body {
				t.Fatalf("body %q", ae.Body)
			}
			for _, s := range sentinels {
				if got, want := errors.Is(err, s), (tc.sentinel != nil && s.Error() == tc.sentinel.Error()); got != want {
					t.Errorf("Is(%v)=%v want %v", s, got, want)
				}
			}
		})
	}
}

func TestAPIErrorString(t *testing.T) {
	e := &APIError{StatusCode: 404, Code: "not_found", Message: "call not found", RequestID: "req_123"}
	if got, want := e.Error(), "audact: 404 not_found: call not found (request_id=req_123)"; got != want {
		t.Fatalf("got %q", got)
	}
}

func TestRateLimit(t *testing.T) {
	future := time.Now().Add(30 * time.Second).UTC().Format(http.TimeFormat)
	tests := []struct {
		name, ra string
		min, max time.Duration
	}{
		{"seconds", "7", 7 * time.Second, 7 * time.Second},
		{"http-date", future, 20 * time.Second, 31 * time.Second},
		{"absent", "", 0, 0},
		{"garbage", "soon", 0, 0},
		{"negative", "-3", 0, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := doErr(t, 429, map[string]string{"Retry-After": tc.ra}, `{"error":{"code":"rl","message":"slow"}}`)
			var rl *RateLimitError
			if !errors.As(err, &rl) {
				t.Fatalf("not RateLimitError: %v", err)
			}
			if rl.RetryAfter < tc.min || rl.RetryAfter > tc.max {
				t.Fatalf("RetryAfter=%v", rl.RetryAfter)
			}
			var ae *APIError
			if !errors.As(err, &ae) || ae.Code != "rl" {
				t.Fatalf("As APIError failed: %v", ae)
			}
			if !errors.Is(err, ErrRateLimited) || errors.Is(err, ErrNotFound) {
				t.Fatal("sentinel mismatch")
			}
		})
	}
}

func TestTimeoutIsDeadlineExceeded(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()
	c, err := NewClient("k", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, err = c.t.do(ctx, http.MethodGet, "/x", authNone, nil, nil, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got %v", err)
	}
}
