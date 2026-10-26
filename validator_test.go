package audact_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestValidatorService(t *testing.T) {
	tests := []struct {
		name     string
		resp     string
		call     func(c *audact.Client) (any, error)
		wantM    string
		wantPath string
		wantBody string
		wantOut  string
	}{
		{"validate", `{"call_id":"c1","verdict":"ALLOW","turn":1,"checks":[{"check_id":"x","duration_ms":0.5,"layers":["canary"]}]}`,
			func(c *audact.Client) (any, error) {
				return c.Validator.Validate(context.Background(), &audact.ValidateRequest{
					CallID: "c1", Content: "hi", Turn: 1,
					Context: &audact.ValidateContext{CallState: &audact.ValidateCallState{DisclosureDelivered: true}, PreviousTurns: []audact.ValidateTurn{{Content: "a", Turn: 0}}},
				})
			},
			"POST", "/v1/validate", `{"call_id":"c1","content":"hi","turn":1,"context":{"call_state":{"disclosure_delivered":true},"previous_turns":[{"content":"a"}]}}`,
			`{"call_id":"c1","verdict":"ALLOW","turn":1,"checks":[{"check_id":"x","duration_ms":0.5,"layers":["canary"]}]}`},
		{"finalize", `{"status":"finalized","composite_score":0.9,"total_turns":3}`,
			func(c *audact.Client) (any, error) {
				return c.Validator.FinalizeScore(context.Background(), "c/1", &audact.FinalizeScoreRequest{AgentID: "a1"})
			},
			"POST", "/v1/calls/c%2F1/finalize-score", `{"agent_id":"a1"}`, `{"status":"finalized","composite_score":0.9,"total_turns":3}`},
		{"finalize skipped nil req", `{"status":"skipped"}`,
			func(c *audact.Client) (any, error) { return c.Validator.FinalizeScore(context.Background(), "c1", nil) },
			"POST", "/v1/calls/c1/finalize-score", `{}`, `{"status":"skipped"}`},
		{"health", `{"status":"ok","components":{"redis":{"status":"ok","latency":"1ms"}}}`,
			func(c *audact.Client) (any, error) { return c.Validator.Health(context.Background()) },
			"GET", "/v1/health", ``, `{"status":"ok","components":{"redis":{"status":"ok","latency":"1ms"}}}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, got := newAgentsServer(t, 200, tc.resp)
			out, err := tc.call(c)
			if err != nil {
				t.Fatal(err)
			}
			checkSeen(t, got, tc.wantM, tc.wantPath, "", "Bearer key123", tc.wantBody)
			b, _ := json.Marshal(out)
			if !jsonEqual(t, b, tc.wantOut) {
				t.Errorf("out = %s, want %s", b, tc.wantOut)
			}
		})
	}
}

func TestValidatorErrors(t *testing.T) {
	c, _ := newAgentsServer(t, 404, `{"error":{"code":"not_found","message":"nope"}}`)
	if _, err := c.Validator.FinalizeScore(context.Background(), "x", nil); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	if _, err := c.Validator.Validate(context.Background(), nil); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	if _, err := c.Validator.Health(context.Background()); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
}
