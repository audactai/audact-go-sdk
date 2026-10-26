// Command smoke runs an end-to-end check of the SDK against a live environment.
//
// Required environment: AUDACT_API_KEY, AUDACT_BASE_URL, AUDACT_SMOKE_AGENT_ID,
// AUDACT_SMOKE_FROM_NUMBER, AUDACT_SMOKE_TO_NUMBER. A missing variable exits 2;
// any other failure exits 1.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	audact "github.com/audactai/audact-go-sdk"
)

var required = []string{
	"AUDACT_API_KEY",
	"AUDACT_BASE_URL",
	"AUDACT_SMOKE_AGENT_ID",
	"AUDACT_SMOKE_FROM_NUMBER",
	"AUDACT_SMOKE_TO_NUMBER",
}

func main() {
	for _, name := range required {
		if os.Getenv(name) == "" {
			fmt.Fprintf(os.Stderr, "missing required environment variable %s\n", name)
			os.Exit(2)
		}
	}
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "FAIL:", err)
		var apiErr *audact.APIError
		if errors.As(err, &apiErr) {
			fmt.Fprintln(os.Stderr, "request_id:", apiErr.RequestID)
		}
		os.Exit(1)
	}
	fmt.Println("OK")
}

// isTerminal reports whether a GET /v1/sip/calls/{call_id}/status value ends
// the call. The server writes "terminated" on hang-up and the stale-call sweep writes "ended"; anything else,
// including the default "active", means the call is still live.
func isTerminal(status string) bool {
	return status == "terminated" || status == "ended"
}

// waitTerminal polls the call's live status every interval until it is
// terminal and returns the last status seen. A normally-ended call has its
// live record deleted (RoomLifecycleManager.EndCall), so a 404 after the call
// was seen live also means it ended; a 404 before that is an error.
func waitTerminal(ctx context.Context, client *audact.Client, callID string, interval time.Duration) (string, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	seen := ""
	for {
		st, err := client.Calls.Status(ctx, callID)
		switch {
		case err == nil:
			seen = st.Status
			fmt.Println("status:", seen)
			if isTerminal(seen) {
				return seen, nil
			}
		case errors.Is(err, audact.ErrNotFound) && seen != "":
			fmt.Println("status: live record gone (call ended)")
			return seen, nil
		default:
			return seen, fmt.Errorf("call status: %w", withRID(err))
		}
		select {
		case <-ctx.Done():
			return seen, fmt.Errorf("waiting for terminal status: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	client, err := audact.NewClient("")
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}
	return smoke(ctx, client, os.Getenv("AUDACT_SMOKE_AGENT_ID"),
		os.Getenv("AUDACT_SMOKE_FROM_NUMBER"), os.Getenv("AUDACT_SMOKE_TO_NUMBER"), 5*time.Second)
}

// smoke runs the four steps against client; pollInterval spaces the status polls.
func smoke(ctx context.Context, client *audact.Client, agentID, from, to string, pollInterval time.Duration) error {
	// Step 1: the smoke agent is visible.
	agents, err := client.Agents.List(ctx, nil)
	if err != nil {
		return fmt.Errorf("list agents: %w", withRID(err))
	}
	found := false
	for _, a := range agents.Agents {
		if a.AgentID == agentID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("agent %s not in agents list", agentID)
	}
	fmt.Println("agent found:", agentID)

	// Step 2: place the call.
	call, err := client.Calls.CreateOutbound(ctx, &audact.OutboundCallRequest{
		AgentID:           agentID,
		FromNumber:        from,
		DestinationNumber: to,
	}, "")
	if err != nil {
		return fmt.Errorf("create outbound: %w", withRID(err))
	}
	fmt.Println("call created:", call.CallID)

	// Step 3: poll until terminal.
	if _, err := waitTerminal(ctx, client, call.CallID, pollInterval); err != nil {
		return err
	}

	// Step 4: the call record is readable and the summary is available.
	if _, err := client.Calls.Get(ctx, call.CallID); err != nil {
		return fmt.Errorf("get call: %w", withRID(err))
	}

	sum, err := client.Calls.Summary(ctx, call.CallID)
	if err != nil {
		return fmt.Errorf("summary: %w", withRID(err))
	}
	fmt.Printf("summary: disposition=%q events=%d\n", sum.Disposition, sum.EventCount)
	return nil
}

// withRID prints the request ID of an API error and returns the error.
func withRID(err error) error {
	var apiErr *audact.APIError
	if errors.As(err, &apiErr) {
		fmt.Println("request_id:", apiErr.RequestID)
	}
	return err
}
