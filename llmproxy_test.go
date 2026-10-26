package audact_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestLLMProxyComplete(t *testing.T) {
	c, got := newAgentsServer(t, 200, `{"id":"x","choices":[{"index":0,"message":{"role":"assistant","content":"hi"}}],"audact_validation":{"verdict":"ALLOW"},"usage":{"total_tokens":3}}`)
	temp := 0.2
	out, err := c.LLMProxy.Complete(context.Background(), &audact.ProxyParams{ConfigurationID: "cfg"},
		&audact.ChatCompletionRequest{Model: "gpt-4o", Temperature: &temp, Messages: []audact.ChatMessage{{Role: "user", Content: "yo"}}})
	if err != nil {
		t.Fatal(err)
	}
	checkSeen(t, got, "POST", "/v1/proxy/openai", "configuration_id=cfg", "Bearer key123",
		`{"model":"gpt-4o","temperature":0.2,"messages":[{"role":"user","content":"yo"}]}`)
	if out.Choices[0].Message.Content != "hi" || out.AudactValidation.Verdict != "ALLOW" || out.Usage.TotalTokens != 3 {
		t.Errorf("out = %+v", out)
	}
}

func TestLLMProxyHeaders(t *testing.T) {
	var hdr http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdr = r.Header.Clone()
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)
	c, _ := audact.NewClient("k", audact.WithBaseURL(srv.URL))
	p := &audact.ProxyParams{CallID: "c", TenantID: "t", ConfigID: "g", Turn: 4}
	if _, err := c.LLMProxy.Complete(context.Background(), p, nil); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"X-Audact-Call-Id": "c", "X-Audact-Tenant-Id": "t", "X-Audact-Config-Id": "g", "X-Audact-Turn": "4"}
	for k, v := range want {
		if hdr.Get(k) != v {
			t.Errorf("%s = %q", k, hdr.Get(k))
		}
	}
	if _, err := c.LLMProxy.Complete(context.Background(), nil, nil); err != nil {
		t.Fatal(err)
	}
	if hdr.Get("X-Audact-Call-Id") != "" || hdr.Get("X-Audact-Turn") != "" {
		t.Error("headers sent for nil params")
	}
}

func sseServer(t *testing.T, handler func(w http.ResponseWriter, f http.Flusher)) *audact.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		if r.URL.Path != "/v1/proxy/openai/stream" || r.Header.Get("Authorization") != "Bearer key123" {
			http.Error(w, "bad", 500)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		handler(w, w.(http.Flusher))
	}))
	t.Cleanup(srv.Close)
	c, err := audact.NewClient("key123", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func collect(t *testing.T, st *audact.Stream) []string {
	t.Helper()
	var got []string
	for st.Next() {
		got = append(got, st.Current())
	}
	return got
}

// Server-shaped SSE events, as written by the proxy's streaming handler.
const (
	chunkAllow  = `data: {"id":"chatcmpl-1","object":"chat.completion.chunk","created":1700000000,"model":"gpt-4o","choices":[{"index":0,"delta":{"content":"Hello there."},"finish_reason":null}],"audact_validation":{"call_id":"c1","verdict":"ALLOW","fail_open":false,"modified":false}}`
	chunkModify = `data: {"id":"chatcmpl-1","object":"chat.completion.chunk","created":1700000000,"model":"gpt-4o","choices":[{"index":0,"delta":{"content":"Your rate is [redacted]."},"finish_reason":null}],"audact_validation":{"call_id":"c1","verdict":"MODIFY","fail_open":false,"modified":true}}`
	chunkTool   = `data: {"id":"chatcmpl-1","object":"chat.completion.chunk","created":1700000000,"model":"gpt-4o","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"end_call","arguments":"{}"}}]},"finish_reason":null}]}`
	chunkBlock  = `data: {"id":"chatcmpl-1","object":"chat.completion.chunk","created":1700000000,"model":"gpt-4o","choices":[{"index":0,"delta":{"content":"I can't help with that."},"finish_reason":null}],"audact_validation":{"call_id":"c1","verdict":"BLOCK","fail_open":false,"modified":false}}`
)

func TestLLMProxyStream(t *testing.T) {
	c := sseServer(t, func(w http.ResponseWriter, f http.Flusher) {
		for _, e := range []string{chunkAllow, `: comment`, chunkModify, chunkTool, chunkBlock, `data: [DONE]`, chunkAllow} {
			_, _ = io.WriteString(w, e+"\n\n")
			f.Flush()
		}
	})
	st, err := c.LLMProxy.Stream(context.Background(), nil, &audact.ChatCompletionRequest{Model: "m"})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()

	type event struct {
		text, verdict string
		modified      bool
		tool          string
	}
	var got []event
	for st.Next() {
		ch := st.Chunk()
		if ch == nil || ch.Object != "chat.completion.chunk" || ch.Model != "gpt-4o" {
			t.Fatalf("chunk = %+v", ch)
		}
		e := event{text: st.Current()}
		if v := ch.AudactValidation; v != nil {
			e.verdict, e.modified = v.Verdict, v.Modified
		}
		if tc := ch.Choices[0].Delta.ToolCalls; len(tc) > 0 {
			e.tool = tc[0].Function.Name
		}
		if !strings.HasPrefix("data: "+st.Raw(), chunkAllow[:40]) {
			t.Errorf("Raw() = %q", st.Raw())
		}
		got = append(got, e)
	}
	want := []event{
		{text: "Hello there.", verdict: "ALLOW"},
		{text: "Your rate is [redacted].", verdict: "MODIFY", modified: true},
		{tool: "end_call"},
		{text: "I can't help with that.", verdict: "BLOCK"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
	if st.Err() != nil || st.Next() {
		t.Errorf("err = %v", st.Err())
	}
}

func TestLLMProxyStreamUndecodableEvent(t *testing.T) {
	c := sseServer(t, func(w http.ResponseWriter, f http.Flusher) {
		_, _ = io.WriteString(w, chunkAllow+"\n\ndata: plain\n\n"+chunkAllow+"\n\n")
		f.Flush()
	})
	st, err := c.LLMProxy.Stream(context.Background(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()
	got := collect(t, st)
	if !reflect.DeepEqual(got, []string{"Hello there."}) {
		t.Errorf("got %q", got)
	}
	if st.Err() == nil || !strings.Contains(st.Err().Error(), "decode chunk") {
		t.Errorf("err = %v, want a decode error", st.Err())
	}
}

func TestLLMProxyStreamMidStreamError(t *testing.T) {
	c := sseServer(t, func(w http.ResponseWriter, f http.Flusher) {
		w.Header().Set("Content-Length", "1000")
		_, _ = io.WriteString(w, chunkAllow+"\n\n")
		f.Flush()
		// Server closes the connection short of Content-Length.
	})
	st, err := c.LLMProxy.Stream(context.Background(), &audact.ProxyParams{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()
	got := collect(t, st)
	if len(got) != 1 || got[0] != "Hello there." {
		t.Errorf("got %q", got)
	}
	if st.Err() == nil {
		t.Error("want mid-stream error")
	}
}

func TestLLMProxyStreamContextCancel(t *testing.T) {
	release := make(chan struct{})
	c := sseServer(t, func(w http.ResponseWriter, f http.Flusher) {
		_, _ = io.WriteString(w, chunkAllow+"\n\n")
		f.Flush()
		<-release
	})
	defer close(release)
	ctx, cancel := context.WithCancel(context.Background())
	st, err := c.LLMProxy.Stream(ctx, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()
	if !st.Next() || st.Current() != "Hello there." {
		t.Fatal("first event")
	}
	cancel()
	if st.Next() {
		t.Fatal("Next after cancel")
	}
	if !errors.Is(st.Err(), context.Canceled) {
		t.Errorf("err = %v", st.Err())
	}
}

func TestLLMProxyStreamHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(409)
		_, _ = io.WriteString(w, `{"error":{"code":"conflict","message":"c"}}`)
	}))
	t.Cleanup(srv.Close)
	c, _ := audact.NewClient("k", audact.WithBaseURL(srv.URL))
	_, err := c.LLMProxy.Stream(context.Background(), nil, nil)
	if !errors.Is(err, audact.ErrConflict) {
		t.Fatalf("err = %v", err)
	}
	_, err = c.LLMProxy.Complete(context.Background(), nil, nil)
	if !errors.Is(err, audact.ErrConflict) {
		t.Fatalf("err = %v", err)
	}
	srv.Close()
	if _, err := c.LLMProxy.Stream(context.Background(), nil, nil); err == nil {
		t.Fatal("want transport error")
	}
}

func TestChatCompletionRequestJSON(t *testing.T) {
	n := 5
	b, _ := json.Marshal(audact.ChatCompletionRequest{Messages: []audact.ChatMessage{}, MaxTokens: &n, ToolChoice: json.RawMessage(`"auto"`),
		Tools: []audact.ToolDefinition{{Type: "function", Function: audact.ToolFunctionSchema{Name: "f", Parameters: json.RawMessage(`{"type":"object"}`)}}}})
	want := `{"max_tokens":5,"messages":[],"tool_choice":"auto","tools":[{"function":{"name":"f","parameters":{"type":"object"}},"type":"function"}]}`
	if !jsonEqual(t, b, want) {
		t.Errorf("got %s", b)
	}
}
