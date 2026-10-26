package audact_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

type svcCase struct {
	name      string
	opts      []audact.Option
	status    int
	resp      string
	respCT    string
	call      func(c *audact.Client) (any, error)
	wantM     string
	wantPath  string
	wantQuery string
	wantBody  string
	wantAuth  string
	wantHdr   map[string][]string
}

func runSvcCases(t *testing.T, tests []svcCase) {
	t.Helper()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tc.wantM {
					t.Errorf("method = %s, want %s", r.Method, tc.wantM)
				}
				if r.URL.EscapedPath() != tc.wantPath {
					t.Errorf("path = %s, want %s", r.URL.EscapedPath(), tc.wantPath)
				}
				if r.URL.RawQuery != tc.wantQuery {
					t.Errorf("query = %q, want %q", r.URL.RawQuery, tc.wantQuery)
				}
				b, _ := io.ReadAll(r.Body)
				if string(b) != tc.wantBody {
					t.Errorf("body = %q, want %q", b, tc.wantBody)
				}
				if got := r.Header.Get("Authorization"); got != tc.wantAuth {
					t.Errorf("auth = %q, want %q", got, tc.wantAuth)
				}
				for k, want := range tc.wantHdr {
					got := r.Header.Values(k)
					if len(got) != len(want) {
						t.Errorf("header %s = %v, want %v", k, got, want)
						continue
					}
					for i := range want {
						if got[i] != want[i] {
							t.Errorf("header %s = %v, want %v", k, got, want)
						}
					}
				}
				if tc.respCT != "" {
					w.Header().Set("Content-Type", tc.respCT)
				}
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.resp))
			}))
			defer srv.Close()
			c, err := audact.NewClient("apikey", append([]audact.Option{audact.WithBaseURL(srv.URL)}, tc.opts...)...)
			if err != nil {
				t.Fatal(err)
			}
			got, err := tc.call(c)
			if err != nil {
				t.Fatalf("call: %v", err)
			}
			if got == nil && tc.resp != "" && tc.respCT == "" {
				t.Fatal("nil result")
			}
		})
	}
}

func notFoundServer(t *testing.T) *audact.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(`{"error":{"code":"not_found","message":"nope","request_id":"req_123"}}`))
	}))
	t.Cleanup(srv.Close)
	c, err := audact.NewClient("apikey", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

var bg = context.Background()
