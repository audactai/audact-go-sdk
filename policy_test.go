package audact_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestPolicyService(t *testing.T) {
	yes := true
	tests := []struct {
		name     string
		resp     string
		call     func(c *audact.Client) (any, error)
		wantM    string
		wantPath string
		wantQ    string
		wantBody string
		wantOut  string
	}{
		{"evaluate", `{"verdict":"BLOCK","reasons":["r"],"applicable_rules":["u1"],"require_disclosure":true}`,
			func(c *audact.Client) (any, error) {
				return c.Policy.Evaluate(context.Background(), &audact.EvaluateRequest{Country: "NL", Context: audact.EvaluationContext{Direction: "outbound", DNCRegistered: &yes}})
			},
			"POST", "/v1/evaluate", "", `{"country":"NL","context":{"direction":"outbound","dnc_registered":true}}`,
			`{"verdict":"BLOCK","reasons":["r"],"applicable_rules":["u1"],"require_disclosure":true}`},
		{"get", `{"country":"NL","version":"1","use_cases":[{"id":"u","conditions":[{"field":"f","operator":"eq","value":"v"}]}]}`,
			func(c *audact.Client) (any, error) { return c.Policy.Get(context.Background(), "N/L") },
			"GET", "/v1/policies/N%2FL", "", ``,
			`{"country":"NL","version":"1","use_cases":[{"id":"u","conditions":[{"field":"f","operator":"eq","value":"v"}]}]}`},
		{"calling hours", `{"allowed":true}`,
			func(c *audact.Client) (any, error) {
				return c.Policy.CheckCallingHours(context.Background(), audact.CheckCallingHoursParams{Country: "DE"})
			},
			"GET", "/v1/calling-hours/check", "country=DE", ``, `{"allowed":true}`},
		{"calling hours time", `{"allowed":false,"reason":"quiet_hours"}`,
			func(c *audact.Client) (any, error) {
				return c.Policy.CheckCallingHours(context.Background(), audact.CheckCallingHoursParams{Country: "DE", CallTime: "2026-01-01T00:00:00Z"})
			},
			"GET", "/v1/calling-hours/check", "call_time=2026-01-01T00%3A00%3A00Z&country=DE", ``, `{"reason":"quiet_hours"}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, got := newAgentsServer(t, 200, tc.resp)
			out, err := tc.call(c)
			if err != nil {
				t.Fatal(err)
			}
			// Matches the TS client, which sends the API key on every call.
			checkSeen(t, got, tc.wantM, tc.wantPath, tc.wantQ, "Bearer key123", tc.wantBody)
			b, _ := json.Marshal(out)
			if !jsonEqual(t, b, tc.wantOut) {
				t.Errorf("out = %s, want %s", b, tc.wantOut)
			}
		})
	}
}

func TestPolicyErrors(t *testing.T) {
	c, _ := newAgentsServer(t, 404, `{"error":{"code":"not_found","message":"nope"}}`)
	if _, err := c.Policy.Get(context.Background(), "XX"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	if _, err := c.Policy.Evaluate(context.Background(), nil); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	if _, err := c.Policy.CheckCallingHours(context.Background(), audact.CheckCallingHoursParams{}); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
}
