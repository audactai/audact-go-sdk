package audact_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestPlatformService(t *testing.T) {
	pt := []audact.Option{audact.WithPlatformToken("plat")}
	tests := []svcCase{
		{name: "export agent yaml", opts: pt, status: 200, resp: "kind: AgentBundle\n", respCT: "application/yaml",
			call: func(c *audact.Client) (any, error) {
				b, err := c.Platform.ExportAgent(bg, "a/1", &audact.PlatformExportAgentOptions{TenantID: "t1"})
				if err == nil && (string(b.Data) != "kind: AgentBundle\n" || b.ContentType != "application/yaml") {
					t.Errorf("bundle = %+v", b)
				}
				return b, err
			},
			wantM: "GET", wantPath: "/v1/platform/export/agents/a%2F1", wantQuery: "tenant_id=t1", wantAuth: "Bearer plat",
			wantHdr: map[string][]string{"Accept": {"application/yaml"}}},
		{name: "export tenant json", opts: pt, status: 200, resp: `{"kind":"TenantBundle"}`, respCT: "application/json",
			call: func(c *audact.Client) (any, error) {
				return c.Platform.ExportTenant(bg, "t1", &audact.PlatformExportOptions{Format: audact.PlatformBundleJSON})
			},
			wantM: "GET", wantPath: "/v1/platform/export/tenants/t1", wantAuth: "Bearer plat",
			wantHdr: map[string][]string{"Accept": {"application/json"}}},
		{name: "import tenant", opts: pt, status: 200, resp: `{"dry_run":true,"counts":{"agents":1},"entities":[{"kind":"Agent","needs_secret":true}]}`,
			call: func(c *audact.Client) (any, error) {
				r, err := c.Platform.ImportTenant(bg, "t2", []byte("kind: TenantBundle\n"), &audact.PlatformImportOptions{
					DryRun: true, Mode: "merge", SetSecret: []string{"A=1", "B=2"}, IdempotencyKey: "k1"})
				if err == nil && (!r.DryRun || r.Counts["agents"] != 1 || !r.Entities[0].NeedsSecret) {
					t.Errorf("resp = %+v", r)
				}
				return r, err
			},
			wantM: "POST", wantPath: "/v1/platform/import/tenants/t2", wantQuery: "dryRun=true&mode=merge", wantBody: "kind: TenantBundle\n", wantAuth: "Bearer plat",
			wantHdr: map[string][]string{"Content-Type": {"application/yaml"}, "X-Audact-Set-Secret": {"A=1", "B=2"}, "Idempotency-Key": {"k1"}}},
		{name: "list admin keys", opts: pt, status: 200, resp: `{"admin_keys":[{"kid":"k","active":true}]}`,
			call:  func(c *audact.Client) (any, error) { return c.Platform.ListAdminKeys(bg) },
			wantM: "GET", wantPath: "/v1/platform/admin-keys", wantAuth: "Bearer plat"},
		{name: "add admin key", opts: pt, status: 201, resp: `{"admin_key":{"kid":"k1","operator":"op"}}`,
			call: func(c *audact.Client) (any, error) {
				return c.Platform.AddAdminKey(bg, audact.AddAdminKeyRequest{KID: "k1", PublicKeyPEM: "PEM", Operator: "op"})
			},
			wantM: "POST", wantPath: "/v1/platform/admin-keys", wantBody: `{"kid":"k1","operator":"op","public_key_pem":"PEM"}`, wantAuth: "Bearer plat"},
		{name: "disable admin key", opts: pt, status: 200, resp: `{"admin_key":{"kid":"k 1","active":false}}`,
			call:  func(c *audact.Client) (any, error) { return c.Platform.DisableAdminKey(bg, "k 1") },
			wantM: "POST", wantPath: "/v1/platform/admin-keys/k%201/disable", wantAuth: "Bearer plat"},
		{name: "orgs list", opts: pt, status: 200, resp: `{"organizations":[{"org_id":"o1"}],"next_cursor":"n"}`,
			call: func(c *audact.Client) (any, error) {
				return c.Platform.Orgs.List(bg, &audact.PlatformListOrganizationsParams{Status: "pending", Email: "a@b.c", Limit: 5, Cursor: "c"})
			},
			wantM: "GET", wantPath: "/v1/platform/organizations", wantQuery: "cursor=c&email=a%40b.c&limit=5&status=pending", wantAuth: "Bearer plat"},
		{name: "orgs list nil params", opts: pt, status: 200, resp: `{}`,
			call:  func(c *audact.Client) (any, error) { return c.Platform.Orgs.List(bg, nil) },
			wantM: "GET", wantPath: "/v1/platform/organizations", wantAuth: "Bearer plat"},
		{name: "orgs verify", opts: pt, status: 200, resp: `{"organization":{"org_id":"o1"}}`,
			call:  func(c *audact.Client) (any, error) { return c.Platform.Orgs.Verify(bg, "o1") },
			wantM: "POST", wantPath: "/v1/platform/organizations/o1/verify", wantAuth: "Bearer plat"},
		{name: "orgs reject", opts: pt, status: 200, resp: `{"organization":{"org_id":"o1"}}`,
			call: func(c *audact.Client) (any, error) {
				return c.Platform.Orgs.Reject(bg, "o1", audact.PlatformRejectOrganizationRequest{Reason: "bad"})
			},
			wantM: "POST", wantPath: "/v1/platform/organizations/o1/reject", wantBody: `{"reason":"bad"}`, wantAuth: "Bearer plat"},
		{name: "orgs find admins", opts: pt, status: 200, resp: `{"org_admins":[{"admin_id":"a1"}]}`,
			call: func(c *audact.Client) (any, error) {
				return c.Platform.Orgs.FindAdmins(bg, &audact.PlatformFindOrgAdminsParams{Email: "a@b.c", OrgID: "o1"})
			},
			wantM: "GET", wantPath: "/v1/platform/org-admins", wantQuery: "email=a%40b.c&org_id=o1", wantAuth: "Bearer plat"},
		{name: "orgs change admin email", opts: pt, status: 200, resp: `{"org_admin":{"admin_id":"a1"}}`,
			call: func(c *audact.Client) (any, error) {
				return c.Platform.Orgs.ChangeAdminEmail(bg, "a1", audact.PlatformChangeOrgAdminEmailRequest{Email: "n@b.c"})
			},
			wantM: "PATCH", wantPath: "/v1/platform/org-admins/a1", wantBody: `{"email":"n@b.c"}`, wantAuth: "Bearer plat"},
		{name: "orgs reset password empty body", opts: pt, status: 200, resp: `{"reset":"ok","expires_in":60}`,
			call:  func(c *audact.Client) (any, error) { return c.Platform.Orgs.ResetAdminPassword(bg, "a1", nil) },
			wantM: "POST", wantPath: "/v1/platform/org-admins/a1/reset-password", wantBody: `{}`, wantAuth: "Bearer plat"},
		{name: "orgs reset password with body", opts: pt, status: 200, resp: `{"reset":"ok"}`,
			call: func(c *audact.Client) (any, error) {
				return c.Platform.Orgs.ResetAdminPassword(bg, "a1", &audact.PlatformResetOrgAdminPasswordRequest{NewPassword: "pw"})
			},
			wantM: "POST", wantPath: "/v1/platform/org-admins/a1/reset-password", wantBody: `{"new_password":"pw"}`, wantAuth: "Bearer plat"},
	}
	runSvcCases(t, tests)
}

// platformCalls returns one invocation of every Platform method on c.
func platformCalls(c *audact.Client) map[string]func() error {
	return map[string]func() error{
		"export agent":  func() error { _, err := c.Platform.ExportAgent(bg, "a", nil); return err },
		"export tenant": func() error { _, err := c.Platform.ExportTenant(bg, "t", nil); return err },
		"import":        func() error { _, err := c.Platform.ImportTenant(bg, "t", nil, nil); return err },
		"list keys":     func() error { _, err := c.Platform.ListAdminKeys(bg); return err },
		"add key":       func() error { _, err := c.Platform.AddAdminKey(bg, audact.AddAdminKeyRequest{}); return err },
		"disable key":   func() error { _, err := c.Platform.DisableAdminKey(bg, "k"); return err },
		"orgs list":     func() error { _, err := c.Platform.Orgs.List(bg, nil); return err },
		"orgs verify":   func() error { _, err := c.Platform.Orgs.Verify(bg, "o"); return err },
		"orgs reject": func() error {
			_, err := c.Platform.Orgs.Reject(bg, "o", audact.PlatformRejectOrganizationRequest{})
			return err
		},
		"orgs find": func() error { _, err := c.Platform.Orgs.FindAdmins(bg, nil); return err },
		"orgs email": func() error {
			_, err := c.Platform.Orgs.ChangeAdminEmail(bg, "a", audact.PlatformChangeOrgAdminEmailRequest{})
			return err
		},
		"orgs reset": func() error { _, err := c.Platform.Orgs.ResetAdminPassword(bg, "a", nil); return err },
	}
}

func TestPlatformServiceNotFound(t *testing.T) {
	c := notFoundServer(t)
	c.SetPlatformToken("plat")
	for name, f := range platformCalls(c) {
		t.Run(name, func(t *testing.T) {
			err := f()
			if !errors.Is(err, audact.ErrNotFound) {
				t.Fatalf("err = %v, want ErrNotFound", err)
			}
			var apiErr *audact.APIError
			if !errors.As(err, &apiErr) || apiErr.RequestID != "req_123" {
				t.Fatalf("apiErr = %+v", apiErr)
			}
		})
	}
}

// TestPlatformMissingTokenFailsLocally: with no platform token every Platform
// method fails with ErrPlatformTokenRequired and sends nothing, so the tenant
// API key never reaches a platform route.
func TestPlatformMissingTokenFailsLocally(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	c, err := audact.NewClient("apikey", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	for name, f := range platformCalls(c) {
		t.Run(name, func(t *testing.T) {
			if err := f(); !errors.Is(err, audact.ErrPlatformTokenRequired) {
				t.Fatalf("err = %v, want ErrPlatformTokenRequired", err)
			}
		})
	}
	if n := hits.Load(); n != 0 {
		t.Fatalf("server received %d requests, want 0", n)
	}

	// Clearing a token restores the local failure.
	c.SetPlatformToken("plat")
	c.SetPlatformToken("")
	if _, err := c.Platform.ListAdminKeys(bg); !errors.Is(err, audact.ErrPlatformTokenRequired) {
		t.Fatalf("after clear: err = %v, want ErrPlatformTokenRequired", err)
	}
}

// TestPlatformSetTokenRotates: platform tokens are single-use, so successive
// calls on one client must carry the token most recently set.
func TestPlatformSetTokenRotates(t *testing.T) {
	var got []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = append(got, r.Header.Get("Authorization"))
		_, _ = w.Write([]byte(`{"admin_keys":[]}`))
	}))
	defer srv.Close()
	c, err := audact.NewClient("apikey", audact.WithBaseURL(srv.URL), audact.WithPlatformToken("tok-1"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Platform.ListAdminKeys(bg); err != nil {
		t.Fatal(err)
	}
	c.SetPlatformToken("tok-2")
	if _, err := c.Platform.Orgs.List(bg, nil); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != "Bearer tok-1" || got[1] != "Bearer tok-2" {
		t.Fatalf("bearers = %q, want [Bearer tok-1 Bearer tok-2]", got)
	}
}
