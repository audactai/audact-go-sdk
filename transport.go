package audact

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type authKind int

const (
	authNone authKind = iota
	authAPIKey
	authOrg
	// authPlatform sends the platform token. A platform route called without
	// one fails locally with ErrPlatformTokenRequired, so the tenant API key is
	// never sent to a super-admin route.
	authPlatform
)

const maxErrorBody = 1 << 20

type transport struct {
	baseURL   string
	hc        *http.Client
	userAgent string
	apiKey    string

	mu            sync.RWMutex
	orgToken      string
	platformToken string
}

// rawBody is a pre-encoded request body (for example multipart/form-data)
// that do sends verbatim with its own Content-Type.
type rawBody struct {
	contentType string
	data        []byte
	// header holds extra (possibly repeated) request headers.
	header http.Header
}

// rawResponse, passed as out, receives the undecoded response body and its
// Content-Type.
type rawResponse struct {
	data        []byte
	contentType string
}

func (t *transport) setOrgToken(tok string) {
	t.mu.Lock()
	t.orgToken = tok
	t.mu.Unlock()
}

func (t *transport) setPlatformToken(tok string) {
	t.mu.Lock()
	t.platformToken = tok
	t.mu.Unlock()
}

func (t *transport) credential(a authKind) (string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	switch a {
	case authAPIKey:
		return t.apiKey, t.apiKey != ""
	case authOrg:
		return t.orgToken, t.orgToken != ""
	case authPlatform:
		return t.platformToken, t.platformToken != ""
	default:
		return "", false
	}
}

// do performs a request to baseURL + "/v1" + path. 2xx bodies are decoded into
// out unless the status is 204 or out is nil.
func (t *transport) do(ctx context.Context, method, path string, auth authKind, q url.Values, in, out any) (*http.Response, error) {
	return t.doHeaders(ctx, method, path, auth, q, nil, in, out)
}

// doHeaders is do with extra request headers (for example Idempotency-Key).
func (t *transport) doHeaders(ctx context.Context, method, path string, auth authKind, q url.Values, extra map[string]string, in, out any) (*http.Response, error) {
	tok, hasTok := t.credential(auth)
	if auth == authPlatform && !hasTok {
		return nil, fmt.Errorf("audact: %s %s: %w", method, path, ErrPlatformTokenRequired)
	}
	full := t.baseURL + "/v1" + path
	if len(q) > 0 {
		full += "?" + q.Encode()
	}
	var body io.Reader
	contentType := "application/json"
	var rawHeader http.Header
	if raw, ok := in.(*rawBody); ok {
		body = bytes.NewReader(raw.data)
		contentType = raw.contentType
		rawHeader = raw.header
	} else if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return nil, fmt.Errorf("audact: %s %s: encode request: %w", method, path, err)
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, full, body)
	if err != nil {
		return nil, fmt.Errorf("audact: %s %s: %w", method, path, redactURLError(err))
	}
	req.Header.Set("User-Agent", t.userAgent)
	req.Header.Set("Accept", "application/json")
	if in != nil {
		req.Header.Set("Content-Type", contentType)
	}
	for k, vs := range rawHeader {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	for k, v := range extra {
		if v != "" {
			req.Header.Set(k, v)
		}
	}
	if hasTok {
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	resp, err := t.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("audact: %s %s: %w", method, path, redactURLError(err))
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
		return resp, parseError(resp, b)
	}
	if resp.StatusCode == http.StatusNoContent || out == nil {
		return resp, nil
	}
	if raw, ok := out.(*rawResponse); ok {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp, fmt.Errorf("audact: %s %s: read response: %w", method, path, err)
		}
		raw.data = b
		raw.contentType = resp.Header.Get("Content-Type")
		return resp, nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return resp, fmt.Errorf("audact: %s %s: decode response: %w", method, path, err)
	}
	return resp, nil
}

// redactURLError drops the query string from the URL a *url.Error carries.
// Query values can be credentials (the evidence-export capability token), and
// callers routinely log transport errors.
func redactURLError(err error) error {
	var ue *url.Error
	if errors.As(err, &ue) {
		if i := strings.IndexByte(ue.URL, '?'); i >= 0 {
			ue.URL = ue.URL[:i] + "?[redacted]"
		}
	}
	return err
}

func addString(q url.Values, k, v string) {
	if v != "" {
		q.Set(k, v)
	}
}

func addInt(q url.Values, k string, v int) {
	if v != 0 {
		q.Set(k, strconv.Itoa(v))
	}
}

func addBool(q url.Values, k string, v bool) {
	if v {
		q.Set(k, "true")
	}
}

func addTime(q url.Values, k string, v time.Time) {
	if !v.IsZero() {
		q.Set(k, v.UTC().Format(time.RFC3339))
	}
}
