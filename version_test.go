package audact

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

var semverRe = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

// checkVersionDocs verifies that the changelog and README agree with version.
func checkVersionDocs(version, changelog, readme string) error {
	if !semverRe.MatchString(version) {
		return fmt.Errorf("version %q is not MAJOR.MINOR.PATCH", version)
	}
	heading := "## [" + version + "]"
	n := 0
	for _, l := range strings.Split(changelog, "\n") {
		if strings.HasPrefix(l, heading) {
			n++
		}
	}
	if n != 1 {
		return fmt.Errorf("CHANGELOG.md has %d lines starting %q, want 1", n, heading)
	}
	lines := strings.Split(changelog, "\n")
	found := false
	for i, l := range lines {
		if strings.TrimSpace(l) != "## [Unreleased]" {
			continue
		}
		found = true
		for _, m := range lines[i+1:] {
			if strings.HasPrefix(m, "## ") {
				break
			}
			if strings.TrimSpace(m) != "" {
				return fmt.Errorf("CHANGELOG.md [Unreleased] section is not empty: %q", m)
			}
		}
	}
	if !found {
		return fmt.Errorf("CHANGELOG.md has no ## [Unreleased] heading")
	}
	if !strings.Contains(readme, "@v"+version) {
		return fmt.Errorf("README.md does not contain %q", "@v"+version)
	}
	return nil
}

// TestVersionMatchesChangelog checks Version against CHANGELOG.md and README.md.
func TestVersionMatchesChangelog(t *testing.T) {
	cl, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		t.Fatal(err)
	}
	rd, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatal(err)
	}
	if err := checkVersionDocs(Version, string(cl), string(rd)); err != nil {
		t.Fatal(err)
	}
}

// TestCheckVersionDocs exercises checkVersionDocs with positive controls.
func TestCheckVersionDocs(t *testing.T) {
	const good = "# Changelog\n\n## [Unreleased]\n\n## [1.2.3] - 2026-01-01\n\n### Added\n\n- x\n"
	const readme = "go get example.com/m@v1.2.3\n"
	tests := []struct {
		name, version, changelog, readme string
		wantErr                          bool
	}{
		{"happy", "1.2.3", good, readme, false},
		{"wrong version heading", "1.2.4", good, "@v1.2.4", true},
		{"duplicate heading", "1.2.3", good + "\n## [1.2.3] - 2026-01-02\n", readme, true},
		{"content under unreleased", "1.2.3", strings.Replace(good, "## [Unreleased]\n", "## [Unreleased]\n\n- pending\n", 1), readme, true},
		{"readme pinned elsewhere", "1.2.3", good, "go get example.com/m@v1.2.2\n", true},
		{"non-semver version", "1.2", strings.ReplaceAll(good, "1.2.3", "1.2"), "@v1.2", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkVersionDocs(tc.version, tc.changelog, tc.readme)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
