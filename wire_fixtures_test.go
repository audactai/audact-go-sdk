package audact_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

// wireFixtureClient starts a server that answers every request with the named
// fixture file and returns a client pointed at it.
func wireFixtureClient(t *testing.T, name string) (*audact.Client, []byte) {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "wire", name))
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(raw)
	}))
	t.Cleanup(srv.Close)
	c, err := audact.NewClient("key123", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	return c, raw
}

func TestWireFixturesDecode(t *testing.T) {
	ctx := context.Background()

	t.Run("calls_list", func(t *testing.T) {
		c, _ := wireFixtureClient(t, "calls_list.json")
		r, err := c.Calls.List(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(r.Calls) != 2 {
			t.Fatalf("want 2 records, got %d", len(r.Calls))
		}
		var withCountry, empty int
		for _, rec := range r.Calls {
			if rec.Metadata == nil {
				t.Errorf("record %s: Metadata is nil", rec.CallID)
				continue
			}
			if len(rec.Metadata) == 0 {
				empty++
				continue
			}
			withCountry++
			if rec.Metadata["country"] != "DE" {
				t.Errorf("country = %v", rec.Metadata["country"])
			}
			camp, ok := rec.Metadata["campaign"].(map[string]any)
			if !ok || camp["id"] != float64(7) {
				t.Errorf("campaign = %v", rec.Metadata["campaign"])
			}
		}
		if withCountry != 1 || empty != 1 {
			t.Errorf("withCountry=%d empty=%d", withCountry, empty)
		}
	})

	t.Run("call_detail", func(t *testing.T) {
		c, _ := wireFixtureClient(t, "call_detail.json")
		r, err := c.Calls.Get(ctx, "22222222-2222-4222-8222-222222222222")
		if err != nil {
			t.Fatal(err)
		}
		if r.Metadata["country"] != "DE" {
			t.Errorf("country = %v", r.Metadata["country"])
		}
	})

	t.Run("reconstruct", func(t *testing.T) {
		c, raw := wireFixtureClient(t, "reconstruct.json")
		r, err := c.Evidence.Reconstruct(ctx, "22222222-2222-4222-8222-222222222222")
		if err != nil {
			t.Fatal(err)
		}
		var want struct {
			Events []struct {
				Data map[string]any `json:"data"`
			} `json:"events"`
		}
		if err := json.Unmarshal(raw, &want); err != nil {
			t.Fatal(err)
		}
		if len(r.Events) == 0 || len(r.Events) != len(want.Events) {
			t.Fatalf("events: got %d want %d", len(r.Events), len(want.Events))
		}
		for i, ev := range r.Events {
			var got map[string]any
			if err := json.Unmarshal(ev.Data, &got); err != nil {
				t.Fatalf("event %d data: %v", i, err)
			}
			if !reflect.DeepEqual(got, want.Events[i].Data) {
				t.Errorf("event %d data = %v, want %v", i, got, want.Events[i].Data)
			}
		}
	})

	t.Run("violations", func(t *testing.T) {
		c, _ := wireFixtureClient(t, "violations.json")
		r, err := c.Evidence.ListViolations(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(r.Violations) < 1 {
			t.Fatal("no violations")
		}
		for _, v := range r.Violations {
			var obj map[string]any
			if len(v.Data) == 0 || json.Unmarshal(v.Data, &obj) != nil || obj == nil {
				t.Errorf("violation data not a JSON object: %q", v.Data)
			}
		}
	})

	t.Run("policy_de", func(t *testing.T) {
		c, _ := wireFixtureClient(t, "policy_de.json")
		r, err := c.Policy.Get(ctx, "DE")
		if err != nil {
			t.Fatal(err)
		}
		var str, boolean, num, arr bool
		for _, uc := range r.UseCases {
			for _, cond := range uc.Conditions {
				switch cond.Value.(type) {
				case string:
					str = true
				case bool:
					boolean = true
				case float64:
					num = true
				case []any:
					arr = true
				}
			}
		}
		if !str || !boolean || !num || !arr {
			t.Errorf("value kinds: string=%v bool=%v float64=%v array=%v", str, boolean, num, arr)
		}
	})
}

func TestWireFixturesPostEventSendsObject(t *testing.T) {
	var body []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)
	c, err := audact.NewClient("key123", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Evidence.PostEvent(context.Background(), &audact.PostEventRequest{
		CallID: "c1", EventType: "X", Data: json.RawMessage(`{"k":1}`),
	}, "idem")
	if err != nil {
		t.Fatal(err)
	}
	var got struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	var obj map[string]any
	if err := json.Unmarshal(got.Data, &obj); err != nil || obj["k"] != float64(1) {
		t.Fatalf("data is not the object {\"k\":1}: %s (%v)", got.Data, err)
	}
}
