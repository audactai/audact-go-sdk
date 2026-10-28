package audact_test

import (
	"encoding/json"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

// TestResponseRoundTripPreservesFalse checks that response types keep a
// boolean false when decoded and re-encoded.
func TestResponseRoundTripPreservesFalse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		field string
		newV  func() any
	}{
		{"WebhookEndpoint", "is_active", func() any { return &audact.WebhookEndpoint{} }},
		{"CreateWebhookResponse", "is_active", func() any { return &audact.CreateWebhookResponse{} }},
		{"AdminKey", "active", func() any { return &audact.AdminKey{} }},
		{"ToolSpeakDuring", "enabled", func() any { return &audact.ToolSpeakDuring{} }},
		{"PIISettingsHistoryEntry", "encryption_enabled", func() any { return &audact.PIISettingsHistoryEntry{} }},
		{"PIISettingsResponse", "encryption_enabled", func() any { return &audact.PIISettingsResponse{} }},
		{"PIISettingsPutResponse", "encryption_enabled", func() any { return &audact.PIISettingsPutResponse{} }},
		{"Configuration", "pii_encryption_enabled", func() any { return &audact.Configuration{} }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := tt.newV()
			if err := json.Unmarshal([]byte(`{"`+tt.field+`":false}`), v); err != nil {
				t.Fatalf("decode: %v", err)
			}
			b, err := json.Marshal(v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var m map[string]any
			if err := json.Unmarshal(b, &m); err != nil {
				t.Fatalf("decode map: %v", err)
			}
			got, ok := m[tt.field]
			if !ok || got != false {
				t.Fatalf("%s: key %q = %v (present=%v) in %s, want false", tt.name, tt.field, got, ok, b)
			}
		})
	}
}
