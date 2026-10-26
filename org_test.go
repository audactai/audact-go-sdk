package audact_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestOrgService(t *testing.T) {
	ctx := context.Background()
	const tenant = `{"tenant":{"tenant_id":"t1","status":"active"}}`
	tests := []struct {
		name      string
		opts      []audact.Option
		status    int
		resp      string
		call      func(c *audact.Client) (any, error)
		wantM     string
		wantPath  string
		wantQuery string
		wantBody  string
		wantAuth  string
		wantIdem  string
	}{
		{"signup", nil, 201, `{"token":"tok","expires_in":900}`,
			func(c *audact.Client) (any, error) {
				return c.Org.Signup(ctx, &audact.OrgSignupRequest{Email: "a@b.c", EntityType: audact.OrgEntityCompany, KYC: &audact.OrgKYC{City: "X"}})
			},
			"POST", "/v1/org/signup", "", `{"email":"a@b.c","entity_type":"company","kyc":{"city":"X"}}`, "", ""},
		{"login", nil, 200, `{"token":"tok","ver":true}`,
			func(c *audact.Client) (any, error) {
				return c.Org.Login(ctx, &audact.OrgLoginRequest{Email: "a@b.c", Password: "pw"})
			},
			"POST", "/v1/org/login", "", `{"email":"a@b.c","password":"pw"}`, "", ""},
		{"refresh", nil, 200, `{"token":"tok2"}`,
			func(c *audact.Client) (any, error) { return c.Org.Refresh(ctx) },
			"POST", "/v1/org/refresh", "", ``, "", ""},
		{"logout", []audact.Option{audact.WithOrgToken("o")}, 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.Org.Logout(ctx) },
			"POST", "/v1/org/logout", "", ``, "Bearer o", ""},
		{"profile", []audact.Option{audact.WithOrgToken("o")}, 200, `{"org":{"org_id":"o1"},"admin":{"email":"e"}}`,
			func(c *audact.Client) (any, error) { return c.Org.Profile(ctx) },
			"GET", "/v1/org/profile", "", ``, "Bearer o", ""},
		{"update profile", []audact.Option{audact.WithOrgToken("o")}, 200, `{}`,
			func(c *audact.Client) (any, error) {
				return c.Org.UpdateProfile(ctx, &audact.OrgUpdateProfileRequest{OrgName: "N"})
			},
			"PATCH", "/v1/org/profile", "", `{"org_name":"N"}`, "Bearer o", ""},
		{"kyc", []audact.Option{audact.WithOrgToken("o")}, 200, `{"kyc":{"city":"X"}}`,
			func(c *audact.Client) (any, error) { return c.Org.KYC(ctx) },
			"GET", "/v1/org/kyc", "", ``, "Bearer o", ""},
		{"update kyc", []audact.Option{audact.WithOrgToken("o")}, 200, `{}`,
			func(c *audact.Client) (any, error) {
				return c.Org.UpdateKYC(ctx, &audact.OrgKYCPatchRequest{City: "Y", Additional: map[string]any{"k": "v"}})
			},
			"PATCH", "/v1/org/kyc", "", `{"additional":{"k":"v"},"city":"Y"}`, "Bearer o", ""},
		{"list tenants", []audact.Option{audact.WithOrgToken("o")}, 200, `{"tenants":[{"tenant_id":"t1"}],"next_cursor":"n"}`,
			func(c *audact.Client) (any, error) {
				return c.Org.ListTenants(ctx, &audact.ListOrgTenantsParams{IncludeDeleted: true, Limit: 5, Cursor: "c"})
			},
			"GET", "/v1/org/tenants", "cursor=c&include_deleted=true&limit=5", ``, "Bearer o", ""},
		{"list tenants nil", []audact.Option{audact.WithOrgToken("o")}, 200, `{}`,
			func(c *audact.Client) (any, error) { return c.Org.ListTenants(ctx, nil) },
			"GET", "/v1/org/tenants", "", ``, "Bearer o", ""},
		{"create tenant", []audact.Option{audact.WithOrgToken("o")}, 201, tenant,
			func(c *audact.Client) (any, error) {
				return c.Org.CreateTenant(ctx, &audact.OrgCreateTenantRequest{Name: "Acme"})
			},
			"POST", "/v1/org/tenants", "", `{"name":"Acme"}`, "Bearer o", ""},
		{"get tenant", []audact.Option{audact.WithOrgToken("o")}, 200, tenant,
			func(c *audact.Client) (any, error) { return c.Org.GetTenant(ctx, "t/1") },
			"GET", "/v1/org/tenants/t%2F1", "", ``, "Bearer o", ""},
		{"delete tenant", []audact.Option{audact.WithOrgToken("o")}, 200, tenant,
			func(c *audact.Client) (any, error) { return c.Org.DeleteTenant(ctx, "t1") },
			"DELETE", "/v1/org/tenants/t1", "", ``, "Bearer o", ""},
		{"disable", []audact.Option{audact.WithOrgToken("o")}, 200, tenant,
			func(c *audact.Client) (any, error) { return c.Org.DisableTenant(ctx, "t1") },
			"POST", "/v1/org/tenants/t1/disable", "", ``, "Bearer o", ""},
		{"enable", []audact.Option{audact.WithOrgToken("o")}, 200, tenant,
			func(c *audact.Client) (any, error) { return c.Org.EnableTenant(ctx, "t1") },
			"POST", "/v1/org/tenants/t1/enable", "", ``, "Bearer o", ""},
		{"restore", []audact.Option{audact.WithOrgToken("o")}, 200, tenant,
			func(c *audact.Client) (any, error) { return c.Org.RestoreTenant(ctx, "t1") },
			"POST", "/v1/org/tenants/t1/restore", "", ``, "Bearer o", ""},
		{"impersonate", []audact.Option{audact.WithOrgToken("o")}, 200, `{"tenant_token":"tt","expires_in":1200}`,
			func(c *audact.Client) (any, error) { return c.Org.ImpersonateTenant(ctx, "t1") },
			"POST", "/v1/org/tenants/t1/impersonate", "", ``, "Bearer o", ""},
		{"initialize", []audact.Option{audact.WithOrgToken("o")}, 201, `{"api_key":{"key":"k"},"tenant":{"tenant_id":"t1"}}`,
			func(c *audact.Client) (any, error) { return c.Org.InitializeTenant(ctx, "t1", "idem-1") },
			"POST", "/v1/org/tenants/t1/initialize", "", ``, "Bearer o", "idem-1"},
		{"issue key", []audact.Option{audact.WithOrgToken("o")}, 201, `{"api_key":{"key":"k"}}`,
			func(c *audact.Client) (any, error) {
				return c.Org.IssueTenantKey(ctx, "t1", &audact.OrgIssueTenantKeyRequest{Name: "ci"}, "")
			},
			"POST", "/v1/org/tenants/t1/keys", "", `{"name":"ci"}`, "Bearer o", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, _ := io.ReadAll(r.Body)
				if r.Method != tc.wantM || r.URL.EscapedPath() != tc.wantPath || r.URL.RawQuery != tc.wantQuery {
					t.Errorf("got %s %s?%s", r.Method, r.URL.EscapedPath(), r.URL.RawQuery)
				}
				if string(b) != tc.wantBody {
					t.Errorf("body = %q, want %q", b, tc.wantBody)
				}
				if got := r.Header.Get("Authorization"); got != tc.wantAuth {
					t.Errorf("auth = %q, want %q", got, tc.wantAuth)
				}
				if got := r.Header.Get("Idempotency-Key"); got != tc.wantIdem {
					t.Errorf("idem = %q, want %q", got, tc.wantIdem)
				}
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.resp))
			}))
			defer srv.Close()
			c, err := audact.NewClient("apikey", append([]audact.Option{audact.WithBaseURL(srv.URL)}, tc.opts...)...)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := tc.call(c); err != nil {
				t.Fatalf("err = %v", err)
			}
		})
	}
}

func TestOrgTokenLifecycle(t *testing.T) {
	var mu sync.Mutex
	var auths []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		auths = append(auths, r.Method+" "+r.URL.Path+" "+r.Header.Get("Authorization"))
		mu.Unlock()
		if r.URL.Path == "/v1/org/login" {
			_, _ = w.Write([]byte(`{"token":"issued"}`))
			return
		}
		if r.URL.Path == "/v1/org/logout" {
			w.WriteHeader(204)
			return
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	c, err := audact.NewClient("k", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if _, err := c.Org.Login(ctx, &audact.OrgLoginRequest{Email: "e", Password: "p"}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Org.Profile(ctx); err != nil {
		t.Fatal(err)
	}
	if err := c.Org.Logout(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Org.Profile(ctx); err != nil {
		t.Fatal(err)
	}
	want := []string{
		"POST /v1/org/login ",
		"GET /v1/org/profile Bearer issued",
		"POST /v1/org/logout Bearer issued",
		"GET /v1/org/profile ",
	}
	if len(auths) != len(want) {
		t.Fatalf("got %v", auths)
	}
	for i := range want {
		if auths[i] != want[i] {
			t.Errorf("req %d = %q, want %q", i, auths[i], want[i])
		}
	}
}

func TestOrgConcurrent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"token":"x"}`))
	}))
	defer srv.Close()
	c, err := audact.NewClient("", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				_, _ = c.Org.Login(context.Background(), nil)
				_, _ = c.Org.Profile(context.Background())
			}
		}()
	}
	wg.Wait()
}

func TestOrgErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(`{"error":{"code":"NOT_FOUND","message":"nope"}}`))
	}))
	defer srv.Close()
	c, _ := audact.NewClient("", audact.WithBaseURL(srv.URL))
	ctx := context.Background()
	if _, err := c.Org.GetTenant(ctx, "x"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	if _, err := c.Org.Login(ctx, nil); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	if err := c.Org.Logout(ctx); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
}
