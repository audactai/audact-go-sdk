package audact_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

// TestOfflineFlow walks the quickstart's calls. Statuses are CallDetail
// values (active | completed | failed | abandoned); the smoke example's
// live-status polling is tested in examples/smoke.
func TestOfflineFlow(t *testing.T) {
	var mu sync.Mutex
	var seen []string
	var outboundBody map[string]any
	statusCalls := 0
	summaryMissing := false

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		seen = append(seen, r.Method+" "+r.URL.Path)
		if got := r.Header.Get("Authorization"); got != "Bearer sk_offline" {
			t.Errorf("%s %s: Authorization = %q", r.Method, r.URL.Path, got)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Method + " " + r.URL.Path {
		case "GET /v1/agents":
			_, _ = io.WriteString(w, `{"agents":[{"agent_id":"ag_1","name":"Bot"}]}`)
		case "POST /v1/sip/calls/outbound":
			b, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(b, &outboundBody)
			_, _ = io.WriteString(w, `{"call_id":"call_1"}`)
		case "GET /v1/calls/call_1":
			statusCalls++
			st := "active"
			if statusCalls > 1 {
				st = "completed"
			}
			_, _ = io.WriteString(w, `{"call_id":"call_1","status":"`+st+`"}`)
		case "GET /v1/calls/call_1/summary":
			if summaryMissing {
				w.Header().Set("X-Request-ID", "hdr")
				w.WriteHeader(http.StatusNotFound)
				_, _ = io.WriteString(w, `{"error":{"code":"not_found","message":"call not found","request_id":"req_123"}}`)
				return
			}
			_, _ = io.WriteString(w, `{"call_id":"call_1","disposition":"answered","event_count":4}`)
		default:
			w.WriteHeader(http.StatusTeapot)
		}
	}))
	defer srv.Close()

	c, err := audact.NewClient("sk_offline", audact.WithBaseURL(srv.URL+"/v1/"))
	if err != nil {
		t.Fatal(err)
	}
	ctx := t.Context()

	agents, err := c.Agents.List(ctx, nil)
	if err != nil || len(agents.Agents) != 1 || agents.Agents[0].AgentID != "ag_1" {
		t.Fatalf("list agents: %+v %v", agents, err)
	}
	call, err := c.Calls.CreateOutbound(ctx, &audact.OutboundCallRequest{
		AgentID: "ag_1", FromNumber: "+10000000001", DestinationNumber: "+10000000002",
	}, "")
	if err != nil || call.CallID != "call_1" {
		t.Fatalf("outbound: %+v %v", call, err)
	}
	if outboundBody["agent_id"] != "ag_1" || outboundBody["from_number"] != "+10000000001" || outboundBody["destination_number"] != "+10000000002" {
		t.Errorf("outbound body = %v", outboundBody)
	}
	d, err := c.Calls.Get(ctx, "call_1")
	if err != nil || d.Status != "active" {
		t.Fatalf("get 1: %+v %v", d, err)
	}
	d, err = c.Calls.Get(ctx, "call_1")
	if err != nil || d.Status != "completed" {
		t.Fatalf("get 2: %+v %v", d, err)
	}
	sum, err := c.Calls.Summary(ctx, "call_1")
	if err != nil || sum.Disposition != "answered" || sum.EventCount != 4 {
		t.Fatalf("summary: %+v %v", sum, err)
	}

	want := []string{
		"GET /v1/agents", "POST /v1/sip/calls/outbound", "GET /v1/calls/call_1",
		"GET /v1/calls/call_1", "GET /v1/calls/call_1/summary",
	}
	mu.Lock()
	if len(seen) != len(want) {
		t.Fatalf("requests = %v, want %v", seen, want)
	}
	for i := range want {
		if seen[i] != want[i] {
			t.Errorf("request %d = %q, want %q", i, seen[i], want[i])
		}
	}
	summaryMissing = true
	mu.Unlock()

	_, err = c.Calls.Summary(ctx, "call_1")
	if !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	var apiErr *audact.APIError
	if !errors.As(err, &apiErr) || apiErr.RequestID != "req_123" {
		t.Fatalf("apiErr = %+v", apiErr)
	}
}

// ExampleAgentsService_Create_elevenLabs creates an agent that speaks through
// ElevenLabs. Pick the voice from ListVoices()["elevenlabs"]; the voice must
// belong to the chosen provider.
func ExampleAgentsService_Create_elevenLabs() {
	c, err := audact.NewClient("sk_example")
	if err != nil {
		return
	}
	_, _ = c.Agents.Create(context.Background(), &audact.CreateAgentRequest{
		Name:        "Concierge",
		TTSProvider: audact.TTSProviderElevenLabs,
		TTSModel:    audact.TTSModelElevenFlashV25, // or TTSModelDefault for the provider default
		VoiceID:     "EXAVITQu4vr4xnSDxMaL",
	})
}
