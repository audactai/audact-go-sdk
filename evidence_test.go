package audact_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestEvidenceService(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		resp      string
		call      func(c *audact.Client) (any, error)
		wantM     string
		wantPath  string
		wantQuery string
		wantBody  string
		wantOut   string
	}{
		{"chain", `{"call_id":"c1","events":[{"sequence":1}],"has_more":true,"next_cursor":1}`,
			func(c *audact.Client) (any, error) {
				return c.Evidence.Chain(ctx, "c/1", &audact.ChainParams{AfterSequence: 5, Limit: 10})
			},
			"GET", "/v1/chains/c%2F1", "after_sequence=5&limit=10", ``, `{"call_id":"c1","events":[{"sequence":1}],"has_more":true,"next_cursor":1}`},
		{"chain nil params", `{"call_id":"c1"}`,
			func(c *audact.Client) (any, error) { return c.Evidence.Chain(ctx, "c1", nil) },
			"GET", "/v1/chains/c1", "", ``, `{"call_id":"c1"}`},
		{"verify", `{"call_id":"c1","status":"valid","event_count":3}`,
			func(c *audact.Client) (any, error) { return c.Evidence.VerifyChain(ctx, "c1") },
			"GET", "/v1/verify/c1", "", ``, `{"call_id":"c1","status":"valid","event_count":3}`},
		{"reconstruct", `{"call_id":"c1","chain_status":"valid","config_versions":{"k":{"config_id":"x","pii_encryption":false}},"event_count":1,"events":[{"event_id":"e","compliance_receipt":{"verdict":"ALLOW"}}],"policy_version":"p1"}`,
			func(c *audact.Client) (any, error) { return c.Evidence.Reconstruct(ctx, "c1") },
			"GET", "/v1/reconstruct/c1", "", ``, `{"call_id":"c1","chain_status":"valid","config_versions":{"k":{"config_id":"x","pii_encryption":false}},"event_count":1,"events":[{"event_id":"e","compliance_receipt":{"verdict":"ALLOW"}}],"policy_version":"p1"}`},
		{"violations", `{"violations":[{"event_id":"e","sequence":2}],"has_more":false}`,
			func(c *audact.Client) (any, error) {
				return c.Evidence.ListViolations(ctx, &audact.ListViolationsParams{CallID: "c1", DateFrom: "2026-01-01", Verdict: "BLOCK", Cursor: "cu", Limit: 5})
			},
			"GET", "/v1/violations", "call_id=c1&cursor=cu&date_from=2026-01-01&limit=5&verdict=BLOCK", ``, `{"violations":[{"event_id":"e","sequence":2}]}`},
		{"receipt", `{"hmac_valid":true,"receipt":{"call_id":"c1","verdict":"MODIFY","input_hash":"a","output_hash":"b"}}`,
			func(c *audact.Client) (any, error) { return c.Evidence.Receipt(ctx, "c1") },
			"GET", "/v1/receipts/c1", "", ``, `{"hmac_valid":true,"receipt":{"call_id":"c1","verdict":"MODIFY","input_hash":"a","output_hash":"b"}}`},
		{"delivery status", `[{"turn_number":1,"status":"confirmed"}]`,
			func(c *audact.Client) (any, error) { return c.Evidence.DeliveryStatus(ctx, "c1") },
			"GET", "/v1/calls/c1/delivery", "", ``, `[{"turn_number":1,"status":"confirmed"}]`},
		{"article50 report", `{"call_id":"c1","direction":"outbound","signature":{"key_id":"k1"}}`,
			func(c *audact.Client) (any, error) {
				exp, err := c.Evidence.Article50Report(ctx, "c1", &audact.Article50Params{Format: "json"})
				if err != nil {
					return nil, err
				}
				return exp.Report, nil
			},
			"GET", "/v1/calls/c1/article50-report", "format=json", ``, `{"call_id":"c1","direction":"outbound","signature":{"key_id":"k1"}}`},
		{"c2pa", c2paManifestJSON,
			func(c *audact.Client) (any, error) { return c.Evidence.GenerateC2PAManifest(ctx, "c1", "") },
			"POST", "/v1/calls/c1/c2pa", "", ``, c2paManifestJSON},
		{"c2pa list", `{"manifests":[{"manifest_id":"m1","call_id":"c1","signature_algorithm":"ES256","asset":{"kind":"recording","url":"/v1/calls/c1/c2pa/m1/asset"},"manifest_store":{"media_type":"application/c2pa","sha256":"ab"},"signing_credential":{"fingerprint_sha256":"fp","trust":"self_signed"},"generated_at":"2026-10-15T10:00:00Z"}],"next_cursor":"n2"}`,
			func(c *audact.Client) (any, error) {
				return c.Evidence.ListC2PAManifests(ctx, "c/1", &audact.ListC2PAManifestsParams{Limit: 20, Cursor: "n1"})
			},
			"GET", "/v1/calls/c%2F1/c2pa", "cursor=n1&limit=20", ``,
			`{"manifests":[{"manifest_id":"m1","call_id":"c1","signature_algorithm":"ES256","asset":{"kind":"recording","url":"/v1/calls/c1/c2pa/m1/asset"},"manifest_store":{"media_type":"application/c2pa","sha256":"ab"},"signing_credential":{"fingerprint_sha256":"fp","trust":"self_signed"},"generated_at":"2026-10-15T10:00:00Z"}],"next_cursor":"n2"}`},
		{"c2pa list nil params", `{"manifests":[]}`,
			func(c *audact.Client) (any, error) { return c.Evidence.ListC2PAManifests(ctx, "c1", nil) },
			"GET", "/v1/calls/c1/c2pa", "", ``, `{}`},
		{"compliance events full", `{"events":[{"call_id":"c1","event_subtype":"fail_open","verdict":"ALLOW","severity":"low","compliance_score":0.5,"turn_sequence":2}],"has_more":true,"next_cursor":"n","total":9}`,
			func(c *audact.Client) (any, error) {
				return c.Evidence.ListComplianceEvents(ctx, &audact.ListComplianceEventsParams{
					EventSubtype: audact.ComplianceEventFailOpen, Verdict: audact.ComplianceVerdictAllow,
					Severity: audact.ComplianceSeverityLow, CallID: "c1", Cursor: "cu", Limit: 50})
			},
			"GET", "/v1/compliance/events", "call_id=c1&cursor=cu&event_subtype=fail_open&limit=50&severity=low&verdict=ALLOW", ``,
			`{"events":[{"call_id":"c1","event_subtype":"fail_open","verdict":"ALLOW","severity":"low","compliance_score":0.5,"turn_sequence":2}],"has_more":true,"next_cursor":"n","total":9}`},
		{"compliance events nil", `{"events":[]}`,
			func(c *audact.Client) (any, error) { return c.Evidence.ListComplianceEvents(ctx, nil) },
			"GET", "/v1/compliance/events", "", ``, `{}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, got := newAgentsServer(t, 200, tc.resp)
			out, err := tc.call(c)
			if err != nil {
				t.Fatal(err)
			}
			checkSeen(t, got, tc.wantM, tc.wantPath, tc.wantQuery, "Bearer key123", tc.wantBody)
			b, _ := json.Marshal(out)
			if !jsonEqual(t, b, tc.wantOut) {
				t.Errorf("out = %s, want %s", b, tc.wantOut)
			}
		})
	}
}

func TestEvidenceWrites(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name     string
		resp     string
		call     func(c *audact.Client) (any, error)
		wantPath string
		wantBody string
		wantKey  string
	}{
		{"post event", `{"event_id":"e1","sequence":1,"hash":"h"}`,
			func(c *audact.Client) (any, error) {
				return c.Evidence.PostEvent(ctx, &audact.PostEventRequest{CallID: "c1", EventType: "CALL_START", TenantID: "t1"}, "idem-1")
			},
			"/v1/events", `{"call_id":"c1","event_type":"CALL_START","tenant_id":"t1"}`, "idem-1"},
		{"post event no key", `{"event_id":"e1"}`,
			func(c *audact.Client) (any, error) {
				return c.Evidence.PostEvent(ctx, &audact.PostEventRequest{CallID: "c1"}, "")
			},
			"/v1/events", `{"call_id":"c1"}`, ""},
		{"confirm delivery", `{"confirmation_grade":"independent","content_match":true}`,
			func(c *audact.Client) (any, error) {
				return c.Evidence.ConfirmDelivery(ctx, &audact.ConfirmDeliveryRequest{CallID: "c1", TurnNumber: 2, ValidatedSegments: []string{"a"}}, "idem-2")
			},
			"/v1/delivery-confirmation", `{"call_id":"c1","turn_number":2,"validated_segments":["a"]}`, "idem-2"},
		{"reconcile", `{"call_id":"c1","total_segments":1,"segment_results":[{"status":"matched"}]}`,
			func(c *audact.Client) (any, error) {
				return c.Evidence.ReconcileTranscript(ctx, map[string]any{"call_id": "c1", "source": "asr"}, "idem-3")
			},
			"/v1/transcript-reconciliation", `{"call_id":"c1","source":"asr"}`, "idem-3"},
		{"reconcile nil", `{"call_id":"c1"}`,
			func(c *audact.Client) (any, error) { return c.Evidence.ReconcileTranscript(ctx, nil, "") },
			"/v1/transcript-reconciliation", `{}`, ""},
		{"c2pa with key", c2paManifestJSON,
			func(c *audact.Client) (any, error) { return c.Evidence.GenerateC2PAManifest(ctx, "c1", "idem-4") },
			"/v1/calls/c1/c2pa", ``, "idem-4"},
		{"c2pa no key", c2paManifestJSON,
			func(c *audact.Client) (any, error) { return c.Evidence.GenerateC2PAManifest(ctx, "c1", "") },
			"/v1/calls/c1/c2pa", ``, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotKey string
			var got seen
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, _ := io.ReadAll(r.Body)
				got = seen{r.Method, r.URL.EscapedPath(), r.URL.RawQuery, r.Header.Get("Authorization"), r.Header.Get("Content-Type"), b}
				gotKey = r.Header.Get("Idempotency-Key")
				_, _ = w.Write([]byte(tc.resp))
			}))
			defer srv.Close()
			c, err := audact.NewClient("key123", audact.WithBaseURL(srv.URL))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := tc.call(c); err != nil {
				t.Fatal(err)
			}
			checkSeen(t, &got, "POST", tc.wantPath, "", "Bearer key123", tc.wantBody)
			if gotKey != tc.wantKey {
				t.Errorf("Idempotency-Key = %q, want %q", gotKey, tc.wantKey)
			}
		})
	}
}

func TestEvidenceUnauthenticated(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		resp      string
		call      func(c *audact.Client) (any, error)
		wantPath  string
		wantQuery string
		wantOut   string
	}{
		{"export json", `{"call_id":"c1"}`,
			func(c *audact.Client) (any, error) {
				exp, err := c.Evidence.Export(ctx, "c 1", "ev1.p.s", &audact.Article50Params{Format: "json"})
				if err != nil {
					return nil, err
				}
				return exp.Report, nil
			},
			"/v1/evidence-exports/c%201", "format=json&token=ev1.p.s", `{"call_id":"c1"}`},
		{"public key", "-----BEGIN PUBLIC KEY-----\nabc\n",
			func(c *audact.Client) (any, error) { return c.Evidence.Article50PublicKey(ctx, "k/1") },
			"/v1/keys/k%2F1", "", `"-----BEGIN PUBLIC KEY-----\nabc\n"`},
		{"c2pa trust anchor", "-----BEGIN CERTIFICATE-----\nroot\n",
			func(c *audact.Client) (any, error) { return c.Evidence.C2PATrustAnchor(ctx) },
			"/v1/c2pa/trust-anchor", "", `"-----BEGIN CERTIFICATE-----\nroot\n"`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, got := newAgentsServer(t, 200, tc.resp)
			out, err := tc.call(c)
			if err != nil {
				t.Fatal(err)
			}
			if got.path != tc.wantPath || got.rawQuery != tc.wantQuery {
				t.Errorf("request = %s ?%s, want %s ?%s", got.path, got.rawQuery, tc.wantPath, tc.wantQuery)
			}
			if got.auth != "" {
				t.Errorf("Authorization = %q, want none", got.auth)
			}
			b, _ := json.Marshal(out)
			if !jsonEqual(t, b, tc.wantOut) {
				t.Errorf("out = %s, want %s", b, tc.wantOut)
			}
		})
	}
}

func TestEvidenceErrors(t *testing.T) {
	ctx := context.Background()
	c, _ := newAgentsServer(t, 404, `{"error":{"code":"not_found","message":"call not found","request_id":"req_123"}}`)
	calls := map[string]func() error{
		"chain":      func() error { _, err := c.Evidence.Chain(ctx, "x", nil); return err },
		"receipt":    func() error { _, err := c.Evidence.Receipt(ctx, "x"); return err },
		"key":        func() error { _, err := c.Evidence.Article50PublicKey(ctx, "x"); return err },
		"compliance": func() error { _, err := c.Evidence.ListComplianceEvents(ctx, nil); return err },
		"post":       func() error { _, err := c.Evidence.PostEvent(ctx, nil, ""); return err },
		"c2pa list":  func() error { _, err := c.Evidence.ListC2PAManifests(ctx, "x", nil); return err },
		"c2pa store": func() error { _, err := c.Evidence.DownloadC2PAManifestStore(ctx, "x", "m"); return err },
		"c2pa asset": func() error { _, err := c.Evidence.DownloadC2PAAsset(ctx, "x", "m"); return err },
	}
	for name, f := range calls {
		t.Run(name, func(t *testing.T) {
			err := f()
			if !errors.Is(err, audact.ErrNotFound) {
				t.Fatalf("err = %v", err)
			}
			var apiErr *audact.APIError
			if !errors.As(err, &apiErr) || apiErr.RequestID != "req_123" {
				t.Errorf("apiErr = %+v", apiErr)
			}
		})
	}
}

// TestEvidenceArticle50Formats serves what the article50 handlers serve
// application/pdf for format=pdf, which is
// also the default of a platform-minted evidence link, else application/json.
func TestEvidenceArticle50Formats(t *testing.T) {
	ctx := context.Background()
	const pdf = "%PDF-1.7\n%\xe2\xe3\xcf\xd3\n1 0 obj\n<<>>\nendobj\n%%EOF\n"
	type req struct{ path, query, auth, accept string }
	var got req
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = req{r.URL.EscapedPath(), r.URL.RawQuery, r.Header.Get("Authorization"), r.Header.Get("Accept")}
		format := r.URL.Query().Get("format")
		if format == "" && r.URL.Query().Get("token") != "" {
			format = "pdf" // the token's embedded format
		}
		if format == "pdf" {
			w.Header().Set("Content-Type", "application/pdf")
			_, _ = io.WriteString(w, pdf)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"call_id":"c1","direction":"outbound"}`)
	}))
	defer srv.Close()
	c, err := audact.NewClient("key123", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		call       func() (*audact.Article50Export, error)
		want       req
		wantPDF    bool
		wantCallID string
	}{
		{"export defaults to the token's pdf",
			func() (*audact.Article50Export, error) { return c.Evidence.Export(ctx, "c1", "ev1.p.s", nil) },
			req{"/v1/evidence-exports/c1", "token=ev1.p.s", "", "application/json, application/pdf"}, true, ""},
		{"export format=pdf",
			func() (*audact.Article50Export, error) {
				return c.Evidence.Export(ctx, "c1", "ev1.p.s", &audact.Article50Params{Format: "pdf"})
			},
			req{"/v1/evidence-exports/c1", "format=pdf&token=ev1.p.s", "", "application/json, application/pdf"}, true, ""},
		{"export format=json",
			func() (*audact.Article50Export, error) {
				return c.Evidence.Export(ctx, "c1", "ev1.p.s", &audact.Article50Params{Format: "json"})
			},
			req{"/v1/evidence-exports/c1", "format=json&token=ev1.p.s", "", "application/json, application/pdf"}, false, "c1"},
		{"report defaults to json",
			func() (*audact.Article50Export, error) { return c.Evidence.Article50Report(ctx, "c1", nil) },
			req{"/v1/calls/c1/article50-report", "", "Bearer key123", "application/json, application/pdf"}, false, "c1"},
		{"report format=pdf",
			func() (*audact.Article50Export, error) {
				return c.Evidence.Article50Report(ctx, "c1", &audact.Article50Params{Format: "pdf"})
			},
			req{"/v1/calls/c1/article50-report", "format=pdf", "Bearer key123", "application/json, application/pdf"}, true, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			exp, err := tc.call()
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("request = %+v, want %+v", got, tc.want)
			}
			if tc.wantPDF {
				if exp.Report != nil || string(exp.PDF) != pdf {
					t.Errorf("export = %+v, want the PDF bytes only", exp)
				}
				return
			}
			if exp.PDF != nil || exp.Report == nil || exp.Report.CallID != tc.wantCallID {
				t.Errorf("export = %+v, want a report for %s only", exp, tc.wantCallID)
			}
		})
	}
}

// TestEvidenceExportTokenNotInTransportError: the capability token is the
// authorisation, so a failed connection must not echo it via the *url.Error.
func TestEvidenceExportTokenNotInTransportError(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	base := srv.URL
	srv.Close() // nothing listens on base any more

	c, err := audact.NewClient("key123", audact.WithBaseURL(base))
	if err != nil {
		t.Fatal(err)
	}
	const token = "ev1.SECRET-CAP-TOKEN.sig"
	_, err = c.Evidence.Export(context.Background(), "c1", token, nil)
	if err == nil {
		t.Fatal("want a transport error")
	}
	if strings.Contains(err.Error(), "SECRET-CAP-TOKEN") {
		t.Errorf("error leaks the token: %v", err)
	}
	var ue *url.Error
	if !errors.As(err, &ue) {
		t.Errorf("err = %v, want a wrapped *url.Error", err)
	}
}

const c2paManifestJSON = `{"manifest_id":"m1","call_id":"c1","claim":{"claim_generator":"audact/2.0.0"},"assertions":[{"label":"l"}],"signature_algorithm":"ES256","asset":{"kind":"evidence_export","media_type":"application/json","sha256":"aa","url":"/v1/calls/c1/c2pa/m1/asset"},"manifest_store":{"media_type":"application/c2pa","sha256":"bb","url":"/v1/calls/c1/c2pa/m1/manifest.c2pa"},"signing_credential":{"subject":"CN=Audact","issuer":"CN=Audact Root","serial":"01","not_before":"2026-01-01T00:00:00Z","not_after":"2027-01-01T00:00:00Z","fingerprint_sha256":"fp","trust":"self_signed"},"generated_at":"2026-10-15T10:00:00Z"}`

// TestEvidenceC2PADownloads covers the two byte endpoints. The asset endpoint
// answers 302 to a presigned URL for recordings;
// the SDK must follow it WITHOUT the Authorization header — the storage host
// stands in for S3, which rejects a second auth mechanism.
func TestEvidenceC2PADownloads(t *testing.T) {
	ctx := context.Background()
	const store = "\x00\x00\x00\x1ejumbc2pa-store-bytes"
	const recording = "RIFF\x24\x00\x00\x00WAVEfmt "

	var storageAuth []string
	storage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		storageAuth = append(storageAuth, r.Header.Get("Authorization"))
		if r.Header.Get("Authorization") != "" || r.URL.Query().Get("X-Amz-Signature") != "sig" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, "<Error><Code>InvalidArgument</Code></Error>")
			return
		}
		w.Header().Set("Content-Type", "audio/wav")
		_, _ = io.WriteString(w, recording)
	}))
	defer storage.Close()

	type req struct{ path, auth, accept string }
	var api []req
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api = append(api, req{r.URL.EscapedPath(), r.Header.Get("Authorization"), r.Header.Get("Accept")})
		switch r.URL.EscapedPath() {
		case "/v1/calls/c%201/c2pa/m%2F1/manifest.c2pa":
			w.Header().Set("Content-Type", "application/c2pa")
			_, _ = io.WriteString(w, store)
		case "/v1/calls/c1/c2pa/export/asset":
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"call_id":"c1"}`)
		case "/v1/calls/c1/c2pa/rec/asset":
			http.Redirect(w, r, storage.URL+"/bucket/rec.wav?X-Amz-Signature=sig", http.StatusFound)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	c, err := audact.NewClient("key123", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("manifest store", func(t *testing.T) {
		api = nil
		b, err := c.Evidence.DownloadC2PAManifestStore(ctx, "c 1", "m/1")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != store {
			t.Errorf("body = %q, want %q", b, store)
		}
		want := req{"/v1/calls/c%201/c2pa/m%2F1/manifest.c2pa", "Bearer key123", "application/c2pa"}
		if len(api) != 1 || api[0] != want {
			t.Errorf("requests = %+v, want [%+v]", api, want)
		}
	})
	t.Run("asset served directly", func(t *testing.T) {
		api, storageAuth = nil, nil
		b, err := c.Evidence.DownloadC2PAAsset(ctx, "c1", "export")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != `{"call_id":"c1"}` {
			t.Errorf("body = %q", b)
		}
		if len(api) != 1 || api[0].auth != "Bearer key123" || len(storageAuth) != 0 {
			t.Errorf("api = %+v storage = %q", api, storageAuth)
		}
	})
	t.Run("asset redirect drops Authorization", func(t *testing.T) {
		api, storageAuth = nil, nil
		b, err := c.Evidence.DownloadC2PAAsset(ctx, "c1", "rec")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != recording {
			t.Errorf("body = %q, want the recording bytes", b)
		}
		if len(api) != 1 || api[0].auth != "Bearer key123" {
			t.Errorf("api requests = %+v, want one authenticated request", api)
		}
		if len(storageAuth) != 1 || storageAuth[0] != "" {
			t.Errorf("storage saw Authorization = %q, want exactly one request with none", storageAuth)
		}
	})
	t.Run("storage error surfaces and hides the presigned query", func(t *testing.T) {
		bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, storage.URL+"/bucket/rec.wav?X-Amz-Signature=WRONG-SECRET", http.StatusFound)
		}))
		defer bad.Close()
		bc, err := audact.NewClient("key123", audact.WithBaseURL(bad.URL))
		if err != nil {
			t.Fatal(err)
		}
		_, err = bc.Evidence.DownloadC2PAAsset(ctx, "c1", "rec")
		var apiErr *audact.APIError
		if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusBadRequest {
			t.Fatalf("err = %v, want a 400 APIError from storage", err)
		}
		if strings.Contains(err.Error(), "WRONG-SECRET") {
			t.Errorf("error leaks the presigned query: %v", err)
		}
	})
}
