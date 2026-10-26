package audact_test

import (
	"errors"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestAggregationService(t *testing.T) {
	zero := 0
	tests := []svcCase{
		{name: "trigger", status: 200, resp: `{"aggregation_id":"ag1","leaf_count":3}`,
			call: func(c *audact.Client) (any, error) {
				return c.Aggregation.Trigger(bg, &audact.AdminAggregateRequest{WindowStart: "s", WindowEnd: "e"})
			},
			wantM: "POST", wantPath: "/v1/admin/aggregate", wantBody: `{"window_end":"e","window_start":"s"}`, wantAuth: "Bearer apikey"},
		{name: "trigger nil", status: 200, resp: `{}`,
			call:  func(c *audact.Client) (any, error) { return c.Aggregation.Trigger(bg, nil) },
			wantM: "POST", wantPath: "/v1/admin/aggregate", wantBody: `{}`, wantAuth: "Bearer apikey"},
		{name: "list proofs", status: 200, resp: `{"proofs":[{"window_id":"w1"}],"has_more":true}`,
			call: func(c *audact.Client) (any, error) {
				return c.Aggregation.ListProofs(bg, &audact.ListProofsParams{Limit: 10, Cursor: "c"})
			},
			wantM: "GET", wantPath: "/v1/proofs", wantQuery: "cursor=c&limit=10", wantAuth: "Bearer apikey"},
		{name: "proof", status: 200, resp: `{"window_id":"w/1","branch_proof":{"steps":[{"hash":"h","position":"left"}]}}`,
			call: func(c *audact.Client) (any, error) {
				p, err := c.Aggregation.Proof(bg, "w/1")
				if err == nil && p.BranchProof.Steps[0].Position != "left" {
					t.Errorf("proof = %+v", p)
				}
				return p, err
			},
			wantM: "GET", wantPath: "/v1/proofs/w%2F1", wantAuth: "Bearer apikey"},
		{name: "metachain", status: 200, resp: `{"entries":[{"sequence":1}],"total":1}`,
			call: func(c *audact.Client) (any, error) {
				return c.Aggregation.Metachain(bg, &audact.MetachainParams{Limit: 2, Cursor: "c", Order: "desc"})
			},
			wantM: "GET", wantPath: "/v1/metachain", wantQuery: "cursor=c&limit=2&order=desc", wantAuth: "Bearer apikey"},
		{name: "metachain nil", status: 200, resp: `{}`,
			call:  func(c *audact.Client) (any, error) { return c.Aggregation.Metachain(bg, nil) },
			wantM: "GET", wantPath: "/v1/metachain", wantAuth: "Bearer apikey"},
		{name: "verify inclusion zero index kept", status: 200, resp: `{"valid":true}`,
			call: func(c *audact.Client) (any, error) {
				return c.Aggregation.VerifyInclusion(bg, audact.VerifyInclusionRequest{
					LeafHash: "l", Root: "r", LeafIndex: &zero, PadCount: &zero, Steps: []audact.ProofStep{{Hash: "h", Position: "right"}}})
			},
			wantM: "POST", wantPath: "/v1/verify-inclusion", wantBody: `{"leaf_hash":"l","leaf_index":0,"pad_count":0,"root":"r","steps":[{"hash":"h","position":"right"}]}`, wantAuth: "Bearer apikey"},
	}
	runSvcCases(t, tests)
}

func TestAggregationServiceNotFound(t *testing.T) {
	c := notFoundServer(t)
	calls := map[string]func() error{
		"trigger": func() error { _, err := c.Aggregation.Trigger(bg, nil); return err },
		"list":    func() error { _, err := c.Aggregation.ListProofs(bg, nil); return err },
		"proof":   func() error { _, err := c.Aggregation.Proof(bg, "w"); return err },
		"meta":    func() error { _, err := c.Aggregation.Metachain(bg, nil); return err },
		"verify":  func() error { _, err := c.Aggregation.VerifyInclusion(bg, audact.VerifyInclusionRequest{}); return err },
	}
	for name, f := range calls {
		t.Run(name, func(t *testing.T) {
			if err := f(); !errors.Is(err, audact.ErrNotFound) {
				t.Fatalf("err = %v", err)
			}
		})
	}
}
