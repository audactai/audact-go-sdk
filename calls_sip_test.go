package audact_test

import (
	"context"
	"errors"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestCallsSIPService(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name, resp                                              string
		status                                                  int
		call                                                    func(c *audact.Client) (string, error)
		wantM, wantPath, wantQuery, wantBody, wantIdem, wantOut string
	}{
		{"handover", `{"call_id":"c1","handover_id":"h1","state":"initiated"}`, 201,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.Handover(ctx, "c1", &audact.HandoverRequest{Destination: "+49", HandoverReason: audact.HandoverReasonEscalation, WhisperMode: audact.WhisperModeBoth})
				return r.HandoverID, err
			},
			"POST", "/v1/sip/calls/c1/handover", "", `{"destination":"+49","handover_reason":"escalation","whisper_mode":"both"}`, "", "h1"},
		{"get handover", `{"handover_id":"h/1","state":"ended"}`, 200,
			func(c *audact.Client) (string, error) { r, err := c.Calls.GetHandover(ctx, "h/1"); return r.State, err },
			"GET", "/v1/sip/handovers/h%2F1", "", "", "", "ended"},
		{"transfer idem", `{"transfer_id":"t1","status":"started"}`, 201,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.Transfer(ctx, "c1", &audact.TransferRequest{Destination: "+1", TransferOption: audact.TransferOptionWarm}, "k1")
				return r.TransferID, err
			},
			"POST", "/v1/sip/calls/c1/transfer", "", `{"destination":"+1","transfer_option":"warm"}`, "k1", "t1"},
		{"transfer nil req", `{"transfer_id":"t2"}`, 201,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.Transfer(ctx, "c1", nil, "")
				return r.TransferID, err
			},
			"POST", "/v1/sip/calls/c1/transfer", "", `{}`, "", "t2"},
		{"get transfer", `{"transfer_id":"t1","status":"bridged"}`, 200,
			func(c *audact.Client) (string, error) { r, err := c.Calls.GetTransfer(ctx, "t1"); return r.Status, err },
			"GET", "/v1/sip/transfers/t1", "", "", "", "bridged"},
		{"bridge", `{"transfer_id":"t1","status":"bridged"}`, 200,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.BridgeTransfer(ctx, "t1")
				return r.Status, err
			},
			"POST", "/v1/sip/transfers/t1/bridge", "", "", "", "bridged"},
		{"cancel", `{"status":"cancelled"}`, 200,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.CancelTransfer(ctx, "t1")
				return r.Status, err
			},
			"POST", "/v1/sip/transfers/t1/cancel", "", "", "", "cancelled"},
		{"end", `{"status":"ended"}`, 200,
			func(c *audact.Client) (string, error) { r, err := c.Calls.EndTransfer(ctx, "t1"); return r.Status, err },
			"POST", "/v1/sip/transfers/t1/end", "", "", "", "ended"},
		{"list routes", `{"routes":[{"did":"+1","dispatch_rule_id":"r1"}],"has_more":true,"next_cursor":"n"}`, 200,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.ListRoutes(ctx, &audact.ListPageParams{Cursor: "cu", Limit: 5})
				return r.Routes[0].DispatchRuleID + r.NextCursor, err
			},
			"GET", "/v1/sip/routes", "cursor=cu&limit=5", "", "", "r1n"},
		{"list routes nil", `{}`, 200,
			func(c *audact.Client) (string, error) { _, err := c.Calls.ListRoutes(ctx, nil); return "", err },
			"GET", "/v1/sip/routes", "", "", "", ""},
		{"create route", `{"dispatch_rule_id":"r2","status":"created"}`, 201,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.CreateRoute(ctx, &audact.RouteInboundRequest{AgentID: "a", DID: "+1"}, "i")
				return r.DispatchRuleID, err
			},
			"POST", "/v1/sip/routes", "", `{"agent_id":"a","did":"+1"}`, "i", "r2"},
		{"list trunks", `{"inbound_trunks":[{"sip_trunk_id":"in1"}],"outbound_trunks":[{"sip_trunk_id":"out1"}]}`, 200,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.ListTrunks(ctx, &audact.ListPageParams{Limit: 2})
				return r.InboundTrunks[0].SIPTrunkID + r.OutboundTrunks[0].SIPTrunkID, err
			},
			"GET", "/v1/sip/trunks", "limit=2", "", "", "in1out1"},
		{"create inbound trunk", `{"sip_trunk_id":"in2","status":"active"}`, 201,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.CreateInboundTrunk(ctx, &audact.RegisterInboundTrunkRequest{Provider: "p", InboundNumbers: []string{"+1"}}, "")
				return r.SIPTrunkID, err
			},
			"POST", "/v1/sip/trunks/inbound", "", `{"inbound_numbers":["+1"],"provider":"p"}`, "", "in2"},
		{"create outbound trunk", `{"sip_trunk_id":"out2"}`, 200,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.CreateOutboundTrunk(ctx, &audact.RegisterOutboundTrunkRequest{Address: "sip.x", Transport: "tls"}, "i2")
				return r.SIPTrunkID, err
			},
			"POST", "/v1/sip/trunks/outbound", "", `{"address":"sip.x","transport":"tls"}`, "i2", "out2"},
		{"get admission url", `{"admission_url":"https://h"}`, 200,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.GetAdmissionURL(ctx)
				return r.AdmissionURL, err
			},
			"GET", "/v1/sip/admission-url", "", "", "", "https://h"},
		{"set admission url", `{"admission_url":"https://h2"}`, 200,
			func(c *audact.Client) (string, error) {
				r, err := c.Calls.SetAdmissionURL(ctx, &audact.SetAdmissionURLRequest{AdmissionURL: "https://h2"})
				return r.AdmissionURL, err
			},
			"PUT", "/v1/sip/admission-url", "", `{"admission_url":"https://h2"}`, "", "https://h2"},
		{"ws ticket", `{"ticket":"tk"}`, 200,
			func(c *audact.Client) (string, error) { r, err := c.Calls.CreateWSTicket(ctx); return r.Ticket, err },
			"POST", "/v1/ws/ticket", "", "", "", "tk"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, got := newCallsServer(t, tc.status, tc.resp)
			out, err := tc.call(c)
			if err != nil {
				t.Fatal(err)
			}
			if got.method != tc.wantM || got.path != tc.wantPath || got.rawQuery != tc.wantQuery {
				t.Errorf("request = %s %s?%s, want %s %s?%s", got.method, got.path, got.rawQuery, tc.wantM, tc.wantPath, tc.wantQuery)
			}
			if got.auth != "Bearer key123" {
				t.Errorf("auth = %q", got.auth)
			}
			if got.body != tc.wantBody {
				t.Errorf("body = %q, want %q", got.body, tc.wantBody)
			}
			if got.idem != tc.wantIdem {
				t.Errorf("idempotency = %q, want %q", got.idem, tc.wantIdem)
			}
			if out != tc.wantOut {
				t.Errorf("out = %q, want %q", out, tc.wantOut)
			}
		})
	}
}

func TestCallsSIPServiceNotFound(t *testing.T) {
	c, _ := newCallsServer(t, 404, `{"error":{"code":"not_found","message":"transfer not found","request_id":"req_123"}}`)
	_, err := c.Calls.GetTransfer(context.Background(), "missing")
	if !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
