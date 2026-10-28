package audact_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestPIIService(t *testing.T) {
	ctx := context.Background()
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
		{"shred", 202, `{"shred_id":"s1","identities_shredded":1,"calls_affected":2,"shredded_at":"t"}`,
			func(c *audact.Client) (any, error) {
				return c.PII.Shred(ctx, &audact.PIIShredRequest{IdentifierType: audact.PIIIdentifierTypeEmail, IdentifierHash: "h"}, "")
			},
			"POST", "/v1/pii/shred", "", `{"identifier_type":"email","identifier_hash":"h"}`,
			`{"shred_id":"s1","identities_shredded":1,"calls_affected":2,"shredded_at":"t"}`},
		{"list shreds", 200, `{"shreds":[{"shred_id":"s1","scope":"call","call_id":"c1"}],"next_cursor":"n"}`,
			func(c *audact.Client) (any, error) { return c.PII.ListShreds(ctx, nil) },
			"GET", "/v1/pii/shreds", "", ``, `{"shreds":[{"shred_id":"s1","scope":"call","call_id":"c1"}],"next_cursor":"n"}`},
		{"list shreds params", 200, `{}`,
			func(c *audact.Client) (any, error) {
				return c.PII.ListShreds(ctx, &audact.ListPIIShredsParams{Cursor: "c", Limit: 5})
			},
			"GET", "/v1/pii/shreds", "cursor=c&limit=5", ``, `{}`},
		{"list call identities", 200, `{"identities":[{"identity_id":"i1","type":"phone","role":"subject","source":"pstn"}]}`,
			func(c *audact.Client) (any, error) { return c.PII.ListCallIdentities(ctx, "c/1") },
			"GET", "/v1/calls/c%2F1/pii/identities", "", ``,
			`{"identities":[{"identity_id":"i1","type":"phone","role":"subject","source":"pstn"}]}`},
		{"link call identity", 201, `{"identity_id":"i1"}`,
			func(c *audact.Client) (any, error) {
				return c.PII.LinkCallIdentity(ctx, "c1", &audact.PIILinkIdentityRequest{IdentifierType: audact.PIIIdentifierTypeName, IdentifierHash: "h"})
			},
			"POST", "/v1/calls/c1/pii/identities", "", `{"identifier_type":"name","identifier_hash":"h"}`, `{"identity_id":"i1"}`},
		{"shred call all", 202, `{"shred_id":"s1"}`,
			func(c *audact.Client) (any, error) { return c.PII.ShredCall(ctx, "c1", nil, "") },
			"POST", "/v1/calls/c1/pii/shred", "", `{}`, `{"shred_id":"s1"}`},
		{"shred call some", 202, `{"shred_id":"s1"}`,
			func(c *audact.Client) (any, error) {
				return c.PII.ShredCall(ctx, "c1", &audact.PIICallShredRequest{IdentityIDs: []string{"a", "b"}}, "")
			},
			"POST", "/v1/calls/c1/pii/shred", "", `{"identity_ids":["a","b"]}`, `{"shred_id":"s1"}`},
		{"decrypt", 200, `{"text":"x","decrypted":1,"erased":0}`,
			func(c *audact.Client) (any, error) {
				return c.PII.Decrypt(ctx, "c1", &audact.PIIDecryptRequest{Text: "tok"})
			},
			"POST", "/v1/calls/c1/pii/decrypt", "", `{"text":"tok"}`, `{"text":"x","decrypted":1}`},
		{"erased keysets one", 200, `{"erased_keyset_ids":["k1"]}`,
			func(c *audact.Client) (any, error) { return c.PII.ListErasedKeysets(ctx, "c1", []string{"k1"}) },
			"GET", "/v1/calls/c1/pii/erased-keysets", "keyset_id=k1", ``, `{"erased_keyset_ids":["k1"]}`},
		{"erased keysets two", 200, `{"erased_keyset_ids":[]}`,
			func(c *audact.Client) (any, error) { return c.PII.ListErasedKeysets(ctx, "c1", []string{"k1", "k2"}) },
			"GET", "/v1/calls/c1/pii/erased-keysets", "keyset_id=k1&keyset_id=k2", ``, `{}`},
		{"list exemptions", 200, `{"exemptions":[{"exempt_type":"email","lookup_prefix":"ab"}],"auto_exempt_dids":3}`,
			func(c *audact.Client) (any, error) { return c.PII.ListExemptions(ctx) },
			"GET", "/v1/pii/exemptions", "", ``, `{"exemptions":[{"exempt_type":"email","lookup_prefix":"ab"}],"auto_exempt_dids":3}`},
		{"add exemption", 201, `{"exempt_type":"phone","lookup_prefix":"ab"}`,
			func(c *audact.Client) (any, error) {
				return c.PII.AddExemption(ctx, &audact.PIIExemptionRequest{ExemptType: audact.PIIExemptTypePhone, IdentifierHash: "h"})
			},
			"POST", "/v1/pii/exemptions", "", `{"exempt_type":"phone","identifier_hash":"h"}`, `{"exempt_type":"phone","lookup_prefix":"ab"}`},
		{"delete exemption body", 204, ``,
			func(c *audact.Client) (any, error) {
				return nil, c.PII.DeleteExemption(ctx, &audact.PIIExemptionRequest{ExemptType: audact.PIIExemptTypeEmail, IdentifierHash: "h"})
			},
			"DELETE", "/v1/pii/exemptions", "", `{"exempt_type":"email","identifier_hash":"h"}`, ``},
		{"delete exemption by hash", 204, ``,
			func(c *audact.Client) (any, error) {
				return nil, c.PII.DeleteExemptionByHash(ctx, audact.PIIExemptTypeEmail, "a/b")
			},
			"DELETE", "/v1/pii/exemptions/email/a%2Fb", "", ``, ``},
		{"settings", 200, `{"encryption_enabled":true,"changed_by":"k","history":[{"encryption_enabled":true,"configurations_resealed":2}]}`,
			func(c *audact.Client) (any, error) { return c.PII.Settings(ctx) },
			"GET", "/v1/pii/settings", "", ``, `{"encryption_enabled":true,"changed_by":"k","history":[{"encryption_enabled":true,"configurations_resealed":2}]}`},
		{"update settings false is sent", 200, `{"configurations_resealed":4}`,
			func(c *audact.Client) (any, error) {
				return c.PII.UpdateSettings(ctx, &audact.PIISettingsPutRequest{EncryptionEnabled: false}, "")
			},
			"PUT", "/v1/pii/settings", "", `{"encryption_enabled":false}`, `{"configurations_resealed":4,"encryption_enabled":false}`},
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

func TestPIIListErasedKeysetsEmptyMakesNoRequest(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(200)
	}))
	defer srv.Close()
	c, err := audact.NewClient("k", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	for _, ids := range [][]string{nil, {}} {
		if _, err := c.PII.ListErasedKeysets(context.Background(), "c1", ids); err == nil {
			t.Error("want error for empty keyset ids")
		}
	}
	if n := atomic.LoadInt32(&hits); n != 0 {
		t.Errorf("server hits = %d, want 0", n)
	}
}

func TestPIIIdempotencyKeys(t *testing.T) {
	ctx := context.Background()
	calls := map[string]func(c *audact.Client) error{
		"shred": func(c *audact.Client) error { _, err := c.PII.Shred(ctx, nil, "idem"); return err },
		"shred call": func(c *audact.Client) error {
			_, err := c.PII.ShredCall(ctx, "c1", nil, "idem")
			return err
		},
		"update settings": func(c *audact.Client) error {
			_, err := c.PII.UpdateSettings(ctx, nil, "idem")
			return err
		},
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			var hdr string
			c, _ := newHeaderServer(t, &hdr, "Idempotency-Key")
			if err := call(c); err != nil {
				t.Fatal(err)
			}
			if hdr != "idem" {
				t.Errorf("Idempotency-Key = %q", hdr)
			}
		})
	}
}

func TestPIINotFound(t *testing.T) {
	ctx := context.Background()
	c, _ := newAgentsServer(t, 404, `{"error":{"code":"not_found","message":"nope"}}`)
	if _, err := c.PII.ListCallIdentities(ctx, "x"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	if err := c.PII.DeleteExemptionByHash(ctx, audact.PIIExemptTypePhone, "h"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
}
