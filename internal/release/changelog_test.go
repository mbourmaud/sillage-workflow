package release

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckAcceptsKeepAChangelogDocument(t *testing.T) {
	t.Parallel()

	path := writeChangelog(t, `# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- A future change.

## [0.1.0] - 2026-08-14

### Added

- The first public protocol.

[unreleased]: https://example.test/unreleased
[0.1.0]: https://example.test/0.1.0
`)

	if err := Check(path, ""); err != nil {
		t.Fatalf("expected structural changelog check to pass: %v", err)
	}
	if err := Check(path, "v0.1.0"); err != nil {
		t.Fatalf("expected versioned changelog check to pass: %v", err)
	}
	notes, err := Extract(path, "0.1.0")
	if err != nil {
		t.Fatalf("expected release notes extraction to pass: %v", err)
	}
	if !strings.Contains(notes, "The first public protocol") || strings.Contains(notes, "[0.1.0]") || strings.Contains(notes, "example.test") {
		t.Fatalf("unexpected extracted notes: %q", notes)
	}
}

func TestCheckRejectsReleaseWithoutVersionedSection(t *testing.T) {
	t.Parallel()

	path := writeChangelog(t, `# Changelog

## [Unreleased]

### Added

- Pending.
`)
	err := Check(path, "0.2.0")
	if err == nil || !strings.Contains(err.Error(), "missing release section") {
		t.Fatalf("expected missing release section error, got %v", err)
	}
}

func TestCheckRejectsMalformedOrEmptyReleaseSections(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name: "missing unreleased",
			content: `# Changelog

## [0.1.0] - 2026-08-14

### Added

- First.
`,
			want: "missing [Unreleased]",
		},
		{
			name: "empty release",
			content: `# Changelog

## [Unreleased]

## [0.1.0] - 2026-08-14
`,
			want: "empty release section",
		},
		{
			name: "invalid date",
			content: `# Changelog

## [Unreleased]

### Added

- Pending.

## [0.1.0] - soon

### Added

- First.
`,
			want: "invalid release date",
		},
		{
			name: "duplicate version",
			content: `# Changelog

## [Unreleased]

### Added

- Pending.

## [0.1.0] - 2026-08-14

### Added

- First.

## [0.1.0] - 2026-08-15

### Fixed

- Duplicate.
`,
			want: "duplicate release section",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := writeChangelog(t, test.content)
			err := Check(path, "")
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected %q, got %v", test.want, err)
			}
		})
	}
}

func TestExtractReportsUnknownVersion(t *testing.T) {
	t.Parallel()

	path := writeChangelog(t, `# Changelog

## [Unreleased]

### Added

- Pending.
`)
	_, err := Extract(path, "0.1.0")
	if err == nil || !errors.Is(err, ErrReleaseNotFound) {
		t.Fatalf("expected unknown version error, got %v", err)
	}
}

func writeChangelog(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "CHANGELOG.md")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
