package audact

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

type piiVector struct {
	Type       string `json:"type"`
	Input      string `json:"input"`
	Normalized string `json:"normalized"`
	ClientHash string `json:"client_hash"`
	Error      bool   `json:"error"`
}

func loadPIIVectors(t *testing.T) []piiVector {
	t.Helper()
	path := os.Getenv("AUDACT_PII_VECTORS_PATH")
	if path == "" {
		path = "testdata/pii_normalize_vectors.json"
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read vectors %s: %v", path, err)
	}
	var vs []piiVector
	if err := json.Unmarshal(raw, &vs); err != nil {
		t.Fatalf("decode vectors: %v", err)
	}
	if len(vs) == 0 {
		t.Fatal("vector file is empty")
	}
	return vs
}

func TestPIIVectors(t *testing.T) {
	for i, v := range loadPIIVectors(t) {
		t.Run(fmt.Sprintf("%d_%s_%q", i, v.Type, v.Input), func(t *testing.T) {
			got, err := NormalizeIdentifier(v.Type, v.Input)
			h, herr := ClientHash(v.Type, v.Input)
			if v.Error {
				if err == nil || herr == nil {
					t.Fatalf("expected errors, got %q / %q", got, h)
				}
				return
			}
			if err != nil || herr != nil {
				t.Fatalf("unexpected errors: %v / %v", err, herr)
			}
			if got != v.Normalized {
				t.Errorf("normalized = %q, want %q", got, v.Normalized)
			}
			if h != v.ClientHash {
				t.Errorf("client_hash = %s, want %s", h, v.ClientHash)
			}
		})
	}
}

func TestNormalizeIdentifierEdgeCases(t *testing.T) {
	tests := []struct {
		typ, in, want string
		wantErr       bool
	}{
		{"name", "ΟΔΥΣΣΕΥΣ", "οδυσσευσ", false},
		{"name", "İstanbul", "istanbul", false},
		{"name", "é", "é", false},
		{"name", "ạ́", "ạ́", false}, // reordering + compose
		{"name", "한", "한", false},  // Hangul jamo
		{"name", "한", "한", false},
		{"email", "A@B.C", "a@b.c", false},
		{"phone", "+1 (415) 555-0100", "+14155550100", false},
		{"phone", "", "", true},
		{"bogus", "x", "", true},
		{"email", "  ", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.typ+"_"+tt.in, func(t *testing.T) {
			got, err := NormalizeIdentifier(tt.typ, tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNFCEmpty(t *testing.T) {
	if got := nfc(""); got != "" {
		t.Errorf("nfc(\"\") = %q", got)
	}
	if got := compose(nil); len(got) != 0 {
		t.Errorf("compose(nil) = %v", got)
	}
}
