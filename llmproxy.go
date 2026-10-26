package audact

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (p *ProxyParams) query() url.Values {
	q := url.Values{}
	if p != nil {
		addString(q, "configuration_id", p.ConfigurationID)
	}
	return q
}

func (p *ProxyParams) headers() map[string]string {
	if p == nil {
		return nil
	}
	h := map[string]string{
		"X-Audact-Call-ID":   p.CallID,
		"X-Audact-Tenant-ID": p.TenantID,
		"X-Audact-Config-ID": p.ConfigID,
	}
	if p.Turn != 0 {
		h["X-Audact-Turn"] = strconv.Itoa(p.Turn)
	}
	return h
}

// Complete sends a chat completion request and returns the buffered response.
// Calls POST /v1/proxy/openai.
func (s *LLMProxyService) Complete(ctx context.Context, params *ProxyParams, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	var out ChatCompletionResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/proxy/openai", authAPIKey, params.query(), params.headers(), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Stream is a server-sent-events stream of chat completion chunks. The proxy
// emits one chunk per validated sentence, each carrying that sentence's
// compliance verdict in AudactValidation; tool-call deltas are forwarded
// as raw upstream chunks without one.
// Call Next until it returns false, then check Err, and always Close.
type Stream struct {
	body    io.Closer
	scanner *bufio.Scanner
	ctx     context.Context
	raw     string
	chunk   *ChatCompletionChunk
	err     error
	done    bool
}

// Stream sends a streaming chat completion request. Errors returned before
// streaming begins are *APIError values. Calls POST /v1/proxy/openai/stream.
func (s *LLMProxyService) Stream(ctx context.Context, params *ProxyParams, req *ChatCompletionRequest) (*Stream, error) {
	const path = "/proxy/openai/stream"
	t := s.t
	b, err := json.Marshal(orEmpty(req))
	if err != nil {
		return nil, fmt.Errorf("audact: POST %s: encode request: %w", path, err)
	}
	full := t.baseURL + "/v1" + path
	if q := params.query(); len(q) > 0 {
		full += "?" + q.Encode()
	}
	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, full, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("audact: POST %s: %w", path, err)
	}
	hreq.Header.Set("User-Agent", t.userAgent)
	hreq.Header.Set("Accept", "text/event-stream")
	hreq.Header.Set("Content-Type", "application/json")
	for k, v := range params.headers() {
		if v != "" {
			hreq.Header.Set(k, v)
		}
	}
	if tok, ok := t.credential(authAPIKey); ok {
		hreq.Header.Set("Authorization", "Bearer "+tok)
	}
	resp, err := t.hc.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("audact: POST %s: %w", path, redactURLError(err))
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer func() { _ = resp.Body.Close() }()
		eb, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
		return nil, parseError(resp, eb)
	}
	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	return &Stream{body: resp.Body, scanner: sc, ctx: ctx}, nil
}

// Next advances to the next data event. It returns false at the [DONE]
// sentinel, at end of stream, on error or when the context is cancelled.
func (st *Stream) Next() bool {
	if st.done {
		return false
	}
	for st.scanner.Scan() {
		if err := st.ctx.Err(); err != nil {
			st.finish(fmt.Errorf("audact: read stream: %w", err))
			return false
		}
		line := strings.TrimSpace(st.scanner.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			st.finish(nil)
			return false
		}
		if data == "" {
			continue
		}
		var chunk ChatCompletionChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			st.finish(fmt.Errorf("audact: read stream: decode chunk: %w", err))
			return false
		}
		st.raw, st.chunk = data, &chunk
		return true
	}
	err := st.scanner.Err()
	if err == nil {
		err = st.ctx.Err()
	}
	if err != nil {
		err = fmt.Errorf("audact: read stream: %w", err)
	}
	st.finish(err)
	return false
}

func (st *Stream) finish(err error) {
	st.done = true
	st.err = err
}

// Chunk returns the chunk read by the last successful Next. Its
// AudactValidation reports the verdict of the sentence it carries (MODIFY
// means the text was rewritten; a BLOCK chunk carries the fallback line and
// ends the stream) and is nil on forwarded tool-call chunks.
func (st *Stream) Chunk() *ChatCompletionChunk { return st.chunk }

// Current returns the text of the last chunk: choices[0].delta.content, or ""
// when the chunk carries none (for example a tool-call delta).
func (st *Stream) Current() string {
	if st.chunk == nil || len(st.chunk.Choices) == 0 {
		return ""
	}
	return st.chunk.Choices[0].Delta.Content
}

// Raw returns the undecoded data of the last event, including any fields
// this SDK version does not model.
func (st *Stream) Raw() string { return st.raw }

// Err returns the error that ended the stream, if any.
func (st *Stream) Err() error { return st.err }

// Close releases the underlying connection.
func (st *Stream) Close() error {
	if err := st.body.Close(); err != nil {
		return fmt.Errorf("audact: close stream: %w", err)
	}
	return nil
}
