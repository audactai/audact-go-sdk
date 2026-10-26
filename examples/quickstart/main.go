// Command quickstart lists agents, places an outbound call and prints its summary.
//
// It reads AUDACT_API_KEY (and optionally AUDACT_BASE_URL) from the
// environment, plus AUDACT_FROM_NUMBER and AUDACT_TO_NUMBER for the call.
// Both numbers must be strict E.164 with no spaces or separators (e.g.
// +442079460958) — they are not normalised — and AUDACT_FROM_NUMBER must be
// one of your own numbers.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	audact "github.com/audactai/audact-go-sdk"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := audact.NewClient("")
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}

	agents, err := client.Agents.List(ctx, &audact.ListAgentsParams{Limit: 10})
	if err != nil {
		return fmt.Errorf("list agents: %w", err)
	}
	for _, a := range agents.Agents {
		fmt.Printf("agent %s (%s)\n", a.AgentID, a.Name)
	}
	if len(agents.Agents) == 0 {
		return fmt.Errorf("no agents found")
	}

	call, err := client.Calls.CreateOutbound(ctx, &audact.OutboundCallRequest{
		AgentID:           agents.Agents[0].AgentID,
		FromNumber:        os.Getenv("AUDACT_FROM_NUMBER"),
		DestinationNumber: os.Getenv("AUDACT_TO_NUMBER"),
	}, "")
	if err != nil {
		return fmt.Errorf("create outbound call: %w", err)
	}
	fmt.Println("call placed:", call.CallID)

	summary, err := client.Calls.Summary(ctx, call.CallID)
	if err != nil {
		return fmt.Errorf("summary: %w", err)
	}
	fmt.Printf("summary: disposition=%q events=%d\n", summary.Disposition, summary.EventCount)
	return nil
}
