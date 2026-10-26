package audact_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestNumbersKYC(t *testing.T) {
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
		{"requirements", 200, `{"requirements":[{"id":"document.passport","required":true}]}`,
			func(c *audact.Client) (any, error) {
				return c.Numbers.KYCRequirements(ctx, &audact.KYCRequirementsParams{Country: "NL", NumberType: "mobile"})
			},
			"GET", "/v1/kyc/requirements", "country=NL&number_type=mobile", "", "", `{"requirements":[{"id":"document.passport","required":true}]}`},
		{"requirements nil", 200, `{}`,
			func(c *audact.Client) (any, error) { return c.Numbers.KYCRequirements(ctx, nil) },
			"GET", "/v1/kyc/requirements", "", "", "", `{}`},
		{"list documents", 200, `{"documents":[{"id":"d1","size_bytes":5}]}`,
			func(c *audact.Client) (any, error) { return c.Numbers.ListKYCDocuments(ctx) },
			"GET", "/v1/kyc/documents", "", "", "", `{"documents":[{"id":"d1","size_bytes":5}]}`},
		{"get document escaped", 200, `{"id":"a/b"}`,
			func(c *audact.Client) (any, error) { return c.Numbers.GetKYCDocument(ctx, "a/b") },
			"GET", "/v1/kyc/documents/a%2Fb", "", "", "", `{"id":"a/b"}`},
		{"delete document", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.Numbers.DeleteKYCDocument(ctx, "d1") },
			"DELETE", "/v1/kyc/documents/d1", "", "", "", ``},
		{"list sets", 200, `{"document_sets":[{"id":"s1","status":"submitted"}]}`,
			func(c *audact.Client) (any, error) { return c.Numbers.ListKYCDocumentSets(ctx) },
			"GET", "/v1/kyc/document-sets", "", "", "", `{"document_sets":[{"id":"s1","status":"submitted"}]}`},
		{"create set", 201, `{"id":"s1","status":"draft"}`,
			func(c *audact.Client) (any, error) {
				return c.Numbers.CreateKYCDocumentSet(ctx, &audact.CreateKYCDocumentSetRequest{
					Name: "NL", CountryCode: "NL", NumberType: "mobile",
					Items: []audact.CreateKYCDocumentSetItemRequest{{RequirementID: "document.passport", DocumentID: "d1"}},
				}, "idem-9")
			},
			"POST", "/v1/kyc/document-sets", "", "idem-9",
			`{"country_code":"NL","items":[{"document_id":"d1","requirement_id":"document.passport"}],"name":"NL","number_type":"mobile"}`,
			`{"id":"s1","status":"draft"}`},
		{"get set", 200, `{"id":"s1"}`,
			func(c *audact.Client) (any, error) { return c.Numbers.GetKYCDocumentSet(ctx, "s1") },
			"GET", "/v1/kyc/document-sets/s1", "", "", "", `{"id":"s1"}`},
		{"delete set", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.Numbers.DeleteKYCDocumentSet(ctx, "s 1") },
			"DELETE", "/v1/kyc/document-sets/s%201", "", "", "", ``},
		{"rotate dek", 200, `{"new_version":2}`,
			func(c *audact.Client) (any, error) { return c.Numbers.RotateKYCDEK(ctx) },
			"POST", "/v1/kyc/dek/rotate", "", "", "", `{"new_version":2}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, got := newNumbersServer(t, tt.status, tt.resp)
			out, err := tt.call(c)
			if err != nil {
				t.Fatal(err)
			}
			if got.method != tt.wantM || got.path != tt.wantPath || got.rawQuery != tt.wantQuery {
				t.Errorf("got %s %s ?%s", got.method, got.path, got.rawQuery)
			}
			if got.auth != "Bearer key123" {
				t.Errorf("auth = %q", got.auth)
			}
			if got.idem != tt.wantIdem {
				t.Errorf("idem = %q", got.idem)
			}
			if tt.wantBody != "" && string(got.body) != tt.wantBody {
				t.Errorf("body = %s", got.body)
			}
			if tt.wantBody == "" && len(got.body) != 0 {
				t.Errorf("unexpected body %s", got.body)
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

func TestNumbersKYC_Multipart(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name       string
		method     string
		path       string
		wantFields map[string]string
		wantFile   string
		wantIdem   string
		call       func(c *audact.Client) (*audact.KYCDocument, error)
	}{
		{"upload", "POST", "/v1/kyc/documents", map[string]string{"doc_type": "passport"}, "scan.pdf", "k1",
			func(c *audact.Client) (*audact.KYCDocument, error) {
				return c.Numbers.UploadKYCDocument(ctx, strings.NewReader("PDFDATA"), "scan.pdf", "passport", &audact.KYCUploadOptions{IdempotencyKey: "k1"})
			}},
		{"upload default filename nil opts", "POST", "/v1/kyc/documents", map[string]string{"doc_type": "poa"}, "upload", "",
			func(c *audact.Client) (*audact.KYCDocument, error) {
				return c.Numbers.UploadKYCDocument(ctx, strings.NewReader("PDFDATA"), "", "poa", nil)
			}},
		{"replace", "PUT", "/v1/kyc/documents/d1", map[string]string{}, "new.pdf", "",
			func(c *audact.Client) (*audact.KYCDocument, error) {
				return c.Numbers.ReplaceKYCDocument(ctx, "d1", strings.NewReader("PDFDATA"), "new.pdf", nil)
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := map[string]string{}
			var file, filename, idem, method, path, auth string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				method, path, idem, auth = r.Method, r.URL.EscapedPath(), r.Header.Get("Idempotency-Key"), r.Header.Get("Authorization")
				mt, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
				if err != nil || mt != "multipart/form-data" {
					t.Errorf("content-type = %q", r.Header.Get("Content-Type"))
					return
				}
				mr := multipart.NewReader(r.Body, params["boundary"])
				for {
					p, err := mr.NextPart()
					if err != nil {
						break
					}
					b, _ := io.ReadAll(p)
					if p.FormName() == "file" {
						file, filename = string(b), p.FileName()
					} else {
						fields[p.FormName()] = string(b)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"id":"d1","doc_type":"passport"}`))
			}))
			t.Cleanup(srv.Close)
			c, err := audact.NewClient("key123", audact.WithBaseURL(srv.URL))
			if err != nil {
				t.Fatal(err)
			}
			doc, err := tt.call(c)
			if err != nil {
				t.Fatal(err)
			}
			if doc.ID != "d1" || method != tt.method || path != tt.path || auth != "Bearer key123" || idem != tt.wantIdem {
				t.Errorf("doc=%+v %s %s auth=%q idem=%q", doc, method, path, auth, idem)
			}
			if file != "PDFDATA" || filename != tt.wantFile {
				t.Errorf("file=%q filename=%q", file, filename)
			}
			if len(fields) != len(tt.wantFields) {
				t.Errorf("fields = %v", fields)
			}
			for k, v := range tt.wantFields {
				if fields[k] != v {
					t.Errorf("field %s = %q", k, fields[k])
				}
			}
		})
	}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errors.New("boom") }

func TestNumbersKYC_ReadError(t *testing.T) {
	c, got := newNumbersServer(t, 200, `{}`)
	if _, err := c.Numbers.UploadKYCDocument(context.Background(), errReader{}, "f", "t", nil); err == nil {
		t.Fatal("want error")
	}
	if got.method != "" {
		t.Error("request should not be sent")
	}
}

func TestNumbersKYC_NotFound(t *testing.T) {
	c, _ := newNumbersServer(t, 404, `{"error":{"code":"not_found","message":"document not found","request_id":"req_123"}}`)
	if _, err := c.Numbers.GetKYCDocument(context.Background(), "x"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	if _, err := c.Numbers.ReplaceKYCDocument(context.Background(), "x", strings.NewReader("a"), "f", nil); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
}
