package audact

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Sentinel errors matched by APIError.Is according to the HTTP status code.
var (
	// ErrUnauthorized matches HTTP 401 responses.
	ErrUnauthorized = errors.New("audact: unauthorized")
	// ErrForbidden matches HTTP 403 responses.
	ErrForbidden = errors.New("audact: forbidden")
	// ErrNotFound matches HTTP 404 responses.
	ErrNotFound = errors.New("audact: not found")
	// ErrConflict matches HTTP 409 responses.
	ErrConflict = errors.New("audact: conflict")
	// ErrRateLimited matches HTTP 429 responses.
	ErrRateLimited = errors.New("audact: rate limited")
)

// ErrPlatformTokenRequired is returned, without a request being sent, when a
// Platform method is called with no platform token set. Mint one per call with
// `audact-admin auth --action <org|admin_keys|export|import>` and pass it via
// WithPlatformToken or Client.SetPlatformToken.
var ErrPlatformTokenRequired = errors.New("audact: platform token required (mint one with audact-admin auth --action <org|admin_keys|export|import>)")

// APIError is returned for every non-2xx response from the Audact API.
type APIError struct {
	// StatusCode is the HTTP status code.
	StatusCode int
	// Code is the machine-readable error code.
	Code string
	// Message is the human-readable error message.
	Message string
	// Hint is an optional remediation hint.
	Hint string
	// Details holds optional structured error details.
	Details json.RawMessage
	// RequestID identifies the request for support purposes.
	RequestID string
	// Body is the raw response body (capped at 1 MiB).
	Body []byte
}

// Error implements the error interface.
func (e *APIError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "audact: %d", e.StatusCode)
	if e.Code != "" {
		b.WriteString(" " + e.Code)
	}
	if e.Message != "" {
		b.WriteString(": " + e.Message)
	}
	if e.RequestID != "" {
		fmt.Fprintf(&b, " (request_id=%s)", e.RequestID)
	}
	return b.String()
}

// Is reports whether target is the sentinel matching e's status code.
func (e *APIError) Is(target error) bool {
	switch e.StatusCode {
	case http.StatusUnauthorized:
		return target == ErrUnauthorized
	case http.StatusForbidden:
		return target == ErrForbidden
	case http.StatusNotFound:
		return target == ErrNotFound
	case http.StatusConflict:
		return target == ErrConflict
	case http.StatusTooManyRequests:
		return target == ErrRateLimited
	}
	return false
}

// RateLimitError is an APIError for HTTP 429 that carries the Retry-After delay.
type RateLimitError struct {
	*APIError
	// RetryAfter is the delay requested by the server, zero if absent.
	RetryAfter time.Duration
}

// Unwrap returns the underlying *APIError.
func (e *RateLimitError) Unwrap() error { return e.APIError }

type errorFields struct {
	Code      string          `json:"code"`
	Message   string          `json:"message"`
	Hint      string          `json:"hint"`
	Details   json.RawMessage `json:"details"`
	RequestID string          `json:"request_id"`
}

func parseError(resp *http.Response, body []byte) error {
	ae := &APIError{StatusCode: resp.StatusCode, Body: body}
	var env struct {
		errorFields
		Error json.RawMessage `json:"error"`
	}
	f := errorFields{}
	if err := json.Unmarshal(body, &env); err == nil {
		f = env.errorFields
		var nested errorFields
		if len(env.Error) > 0 && json.Unmarshal(env.Error, &nested) == nil && env.Error[0] == '{' {
			f = nested
		} else if len(env.Error) > 0 && env.Error[0] == '"' {
			var s string
			if json.Unmarshal(env.Error, &s) == nil && f.Code == "" {
				f.Code = s
			}
		}
	}
	ae.Code, ae.Message, ae.Hint, ae.Details, ae.RequestID = f.Code, f.Message, f.Hint, f.Details, f.RequestID
	if ae.Message == "" {
		ae.Message = http.StatusText(resp.StatusCode)
	}
	if ae.RequestID == "" {
		ae.RequestID = resp.Header.Get("X-Request-ID")
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return &RateLimitError{APIError: ae, RetryAfter: parseRetryAfter(resp.Header.Get("Retry-After"))}
	}
	return ae
}

func parseRetryAfter(v string) time.Duration {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}
	if n, err := strconv.Atoi(v); err == nil {
		if n < 0 {
			return 0
		}
		return time.Duration(n) * time.Second
	}
	if tm, err := http.ParseTime(v); err == nil {
		if d := time.Until(tm); d > 0 {
			return d
		}
	}
	return 0
}
