package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	audact "github.com/audactai/audact-go-sdk"
)

const smokeCallID = "11111111-1111-1111-1111-111111111111"

// smokeServer serves the smoke's four steps. statuses is replayed in order by
// GET /v1/sip/calls/{id}/status; the value "404" answers with a not-found
// envelope, and the last entry repeats once the list is exhausted.
func smokeServer(t *testing.T, statuses []string) (*audact.Client, *[]string) {
	t.Helper()
	var mu sync.Mutex
	var seen []string
	polls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		seen = append(seen, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch r.Method + " " + r.URL.Path {
		case "GET /v1/agents":
			_, _ = io.WriteString(w, `{"agents":[{"agent_id":"ag_1","name":"Bot"}]}`)
		case "POST /v1/sip/calls/outbound":
			_, _ = io.WriteString(w, `{"call_id":"`+smokeCallID+`"}`)
		case "GET /v1/sip/calls/" + smokeCallID + "/status":
			st := statuses[min(polls, len(statuses)-1)]
			polls++
			if st == "404" {
				w.WriteHeader(http.StatusNotFound)
				_, _ = io.WriteString(w, `{"error":{"code":"not_found","message":"call not found","request_id":"req_404"}}`)
				return
			}
			_, _ = io.WriteString(w, `{"call_id":"`+smokeCallID+`","status":"`+st+`"}`)
		case "GET /v1/calls/" + smokeCallID:
			_, _ = io.WriteString(w, `{"call_id":"`+smokeCallID+`","status":"completed"}`)
		case "GET /v1/calls/" + smokeCallID + "/summary":
			_, _ = io.WriteString(w, `{"call_id":"`+smokeCallID+`","disposition":"answered","event_count":4}`)
		default:
			w.WriteHeader(http.StatusTeapot)
		}
	}))
	t.Cleanup(srv.Close)
	c, err := audact.NewClient("sk_smoke", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	return c, &seen
}

func TestIsTerminal(t *testing.T) {
	// Values from the SIP call status response.
	for _, st := range []string{"terminated", "ended"} {
		if !isTerminal(st) {
			t.Errorf("isTerminal(%q) = false, want true", st)
		}
	}
	for _, st := range []string{"active", ""} {
		if isTerminal(st) {
			t.Errorf("isTerminal(%q) = true, want false", st)
		}
	}
}

func TestSmokeEndsOnTerminated(t *testing.T) {
	c, seen := smokeServer(t, []string{"active", "active", "terminated"})
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	if err := smoke(ctx, c, "ag_1", "+10000000001", "+10000000002", time.Millisecond); err != nil {
		t.Fatalf("smoke: %v", err)
	}
	want := []string{
		"GET /v1/agents",
		"POST /v1/sip/calls/outbound",
		"GET /v1/sip/calls/" + smokeCallID + "/status",
		"GET /v1/sip/calls/" + smokeCallID + "/status",
		"GET /v1/sip/calls/" + smokeCallID + "/status",
		"GET /v1/calls/" + smokeCallID,
		"GET /v1/calls/" + smokeCallID + "/summary",
	}
	if strings.Join(*seen, "\n") != strings.Join(want, "\n") {
		t.Fatalf("requests =\n%s\nwant\n%s", strings.Join(*seen, "\n"), strings.Join(want, "\n"))
	}
}

func TestWaitTerminalLiveRecordDeletedAfterActive(t *testing.T) {
	c, _ := smokeServer(t, []string{"active", "404"})
	got, err := waitTerminal(t.Context(), c, smokeCallID, time.Millisecond)
	if err != nil {
		t.Fatalf("waitTerminal: %v", err)
	}
	if got != "active" {
		t.Errorf("last status = %q, want active", got)
	}
}

func TestWaitTerminalNotFoundBeforeLiveIsError(t *testing.T) {
	c, _ := smokeServer(t, []string{"404"})
	_, err := waitTerminal(t.Context(), c, smokeCallID, time.Millisecond)
	if !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestWaitTerminalNeverTerminalTimesOut(t *testing.T) {
	c, _ := smokeServer(t, []string{"active"})
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	_, err := waitTerminal(ctx, c, smokeCallID, time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context.DeadlineExceeded", err)
	}
}
