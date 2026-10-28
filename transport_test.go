package audact

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

var errTest = errors.New("boom")

func TestTransportAuthAndHeaders(t *testing.T) {
	tests := []struct {
		name string
		auth authKind
		want string
	}{
		{"none", authNone, ""},
		{"api key", authAPIKey, "Bearer key"},
		{"org", authOrg, "Bearer org"},
		{"platform", authPlatform, "Bearer plat"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got *http.Request
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = r
				_, _ = w.Write([]byte(`{"a":1}`))
			}))
			defer srv.Close()
			c, err := NewClient("key", WithBaseURL(srv.URL+"/"), WithOrgToken("org"), WithPlatformToken("plat"), WithUserAgent("myapp/1.0"))
			if err != nil {
				t.Fatal(err)
			}
			var out struct{ A int }
			q := url.Values{}
			addString(q, "s", "x")
			addString(q, "empty", "")
			addInt(q, "n", 3)
			addInt(q, "z", 0)
			addBool(q, "b", true)
			addBool(q, "nb", false)
			addTime(q, "t", time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
			addTime(q, "zt", time.Time{})
			if _, err := c.t.do(context.Background(), http.MethodPost, "/things", tc.auth, q, map[string]int{"x": 1}, &out); err != nil {
				t.Fatal(err)
			}
			if got.URL.Path != "/v1/things" || got.Method != http.MethodPost {
				t.Fatalf("bad request %s %s", got.Method, got.URL.Path)
			}
			if got.URL.RawQuery != "b=true&n=3&s=x&t=2026-01-02T03%3A04%3A05Z" {
				t.Fatalf("query %q", got.URL.RawQuery)
			}
			if h := got.Header.Get("Authorization"); h != tc.want {
				t.Fatalf("auth %q want %q", h, tc.want)
			}
			if ua := got.Header.Get("User-Agent"); ua != "audact-go/"+Version+" myapp/1.0" {
				t.Fatalf("ua %q", ua)
			}
			if got.Header.Get("Accept") != "application/json" || got.Header.Get("Content-Type") != "application/json" {
				t.Fatal("missing Accept/Content-Type")
			}
			if out.A != 1 {
				t.Fatalf("decode: %+v", out)
			}
		})
	}
}

func TestTransportNoBodyNoContentType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "" {
			t.Error("unexpected Content-Type")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()
	c, _ := NewClient("k", WithBaseURL(srv.URL))
	var out map[string]any
	resp, err := c.t.do(context.Background(), http.MethodDelete, "/x", authAPIKey, nil, nil, &out)
	if err != nil || resp.StatusCode != 204 {
		t.Fatalf("resp=%v err=%v", resp, err)
	}
}

func TestTransportErrors(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		in      any
	}{
		{"non-2xx", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(404) }, nil},
		{"bad json", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("{")) }, nil},
		{"unencodable", func(w http.ResponseWriter, _ *http.Request) {}, make(chan int)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(tc.handler)
			defer srv.Close()
			c, _ := NewClient("k", WithBaseURL(srv.URL))
			var out json.RawMessage
			if _, err := c.t.do(context.Background(), http.MethodGet, "/x", authAPIKey, nil, tc.in, &out); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestTransportContextCanceled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer srv.Close()
	c, _ := NewClient("k", WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := c.t.do(ctx, http.MethodGet, "/x", authNone, nil, nil, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v", err)
	}
}

func TestOrgTokenConcurrent(t *testing.T) {
	c, _ := NewClient("k")
	done := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			c.t.setOrgToken("a")
		}
		close(done)
	}()
	for i := 0; i < 100; i++ {
		c.t.credential(authOrg)
	}
	<-done
}
