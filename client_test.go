package audact

import (
	"testing"
)

func TestNewClientBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		opt     *string
		want    string
		wantErr bool
	}{
		{name: "default", want: DefaultBaseURL},
		{name: "env", env: "https://env.example", want: "https://env.example"},
		{name: "option beats env", env: "https://env.example", opt: ptr("https://opt.example"), want: "https://opt.example"},
		{name: "plain", opt: ptr("https://h"), want: "https://h"},
		{name: "slash", opt: ptr("https://h/"), want: "https://h"},
		{name: "v1", opt: ptr("https://h/v1"), want: "https://h"},
		{name: "v1 slash", opt: ptr("https://h/v1/"), want: "https://h"},
		{name: "ftp", opt: ptr("ftp://h"), wantErr: true},
		{name: "unparseable", opt: ptr("http://[::1"), wantErr: true},
		{name: "no host", opt: ptr("https://"), wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("AUDACT_BASE_URL", tc.env)
			var opts []Option
			if tc.opt != nil {
				opts = append(opts, WithBaseURL(*tc.opt))
			}
			c, err := NewClient("k", opts...)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if c.BaseURL() != tc.want {
				t.Fatalf("got %q want %q", c.BaseURL(), tc.want)
			}
		})
	}
}

func TestNewClientAPIKeyEnv(t *testing.T) {
	t.Setenv("AUDACT_API_KEY", "envkey")
	c, err := NewClient("")
	if err != nil {
		t.Fatal(err)
	}
	if c.t.apiKey != "envkey" {
		t.Fatalf("got %q", c.t.apiKey)
	}
	t.Setenv("AUDACT_API_KEY", "")
	if _, err := NewClient(""); err != nil {
		t.Fatalf("missing key must not fail: %v", err)
	}
}

func TestNewClientOptionErrors(t *testing.T) {
	if _, err := NewClient("k", nil); err == nil {
		t.Fatal("nil option must error")
	}
	if _, err := NewClient("k", func(*Client) error { return errTest }); err == nil {
		t.Fatal("failing option must error")
	}
}

func TestServicesNonNil(t *testing.T) {
	c, err := NewClient("k")
	if err != nil {
		t.Fatal(err)
	}
	svcs := map[string]bool{
		"Agents": c.Agents != nil, "Calls": c.Calls != nil, "Numbers": c.Numbers != nil,
		"Webhooks": c.Webhooks != nil, "Tools": c.Tools != nil, "Integrations": c.Integrations != nil,
		"Evidence": c.Evidence != nil, "PII": c.PII != nil, "Validator": c.Validator != nil,
		"LLMProxy": c.LLMProxy != nil, "Tenants": c.Tenants != nil, "Org": c.Org != nil,
		"Platform": c.Platform != nil, "Aggregation": c.Aggregation != nil, "Policy": c.Policy != nil,
		"Audio": c.Audio != nil, "Auth": c.Auth != nil, "Platform.Orgs": c.Platform.Orgs != nil,
	}
	if len(svcs) != 18 {
		t.Fatalf("expected 18 entries, got %d", len(svcs))
	}
	for n, ok := range svcs {
		if !ok {
			t.Errorf("%s is nil", n)
		}
	}
}

func ptr[T any](v T) *T { return &v }
