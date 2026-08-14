package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDoctorRequiresCanonicalEntryPoints(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	report := Inspect(dir)

	if report.OK {
		t.Fatal("expected an empty project to fail inspection")
	}
	assertFinding(t, report, "missing_product")
	assertFinding(t, report, "missing_design")
	assertFinding(t, report, "missing_agents")
	assertFinding(t, report, "missing_domain_index")
}

func TestDoctorRequiresClaudeSymlinkToAgents(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeRequiredFiles(t, dir)
	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("duplicate"), 0o600); err != nil {
		t.Fatal(err)
	}

	report := Inspect(dir)
	assertFinding(t, report, "claude_not_agents_symlink")
}

func TestDoctorAcceptsPortableProjectContract(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeRequiredFiles(t, dir)
	if err := os.Symlink("AGENTS.md", filepath.Join(dir, "CLAUDE.md")); err != nil {
		t.Fatal(err)
	}

	report := Inspect(dir)
	if !report.OK {
		t.Fatalf("expected project contract to pass: %#v", report.Findings)
	}
}

func TestDoctorUsesConfiguredEntryPoints(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	paths := []string{"docs/product.md", "docs/design.md", "BOT.md", "docs/vocabulary/index.md"}
	for _, path := range paths {
		absolute := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(absolute), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(absolute, []byte("# Test\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	profilePath := filepath.Join(dir, ".sillage", "project.json")
	if err := os.MkdirAll(filepath.Dir(profilePath), 0o700); err != nil {
		t.Fatal(err)
	}
	profile := `{"version":1,"project":{"name":"Custom","product":"docs/product.md","design":"docs/design.md","agents":"BOT.md","domain":"docs/vocabulary/index.md"},"task_store":{"provider":"local"},"authority":{"human_required":["product_decision"]},"verification":{"fast":[{"name":"test","command":"make test"}],"complete":[{"name":"check","command":"make check"}]}}`
	if err := os.WriteFile(profilePath, []byte(profile), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("BOT.md", filepath.Join(dir, "CLAUDE.md")); err != nil {
		t.Fatal(err)
	}

	report := Inspect(dir)
	if !report.OK {
		t.Fatalf("expected configured project contract to pass: %#v", report.Findings)
	}
}

func TestDoctorAllowsDeclaredAdapterExtensions(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeRequiredFiles(t, dir)
	if err := os.Symlink("AGENTS.md", filepath.Join(dir, "CLAUDE.md")); err != nil {
		t.Fatal(err)
	}
	profilePath := filepath.Join(dir, ".sillage", "project.json")
	if err := os.MkdirAll(filepath.Dir(profilePath), 0o700); err != nil {
		t.Fatal(err)
	}
	profile := `{"version":1,"project":{"name":"Custom","product":"PRODUCT.md","design":"DESIGN.md","agents":"AGENTS.md","domain":"docs/domain/index.md"},"task_store":{"provider":"gitlab","project_id":123},"authority":{"human_required":["product_decision"]},"verification":{"fast":[{"name":"test","command":"make test"}],"complete":[{"name":"check","command":"make check"}],"ci_provider":"gitlab"}}`
	if err := os.WriteFile(profilePath, []byte(profile), 0o600); err != nil {
		t.Fatal(err)
	}
	if report := Inspect(dir); !report.OK {
		t.Fatalf("expected adapter extensions to pass: %#v", report.Findings)
	}
}

func TestDoctorRejectsProfileOutsidePublishedContract(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		profile string
	}{
		{"unknown field", `{"version":1,"unexpected":true}`},
		{"blank authority", `{"version":1,"project":{"name":"Custom","product":"PRODUCT.md","design":"DESIGN.md","agents":"AGENTS.md","domain":"docs/domain/index.md"},"task_store":{"provider":"local"},"authority":{"human_required":[" "]},"verification":{"fast":[{"name":"test","command":"make test"}],"complete":[{"name":"check","command":"make check"}]}}`},
		{"blank command", `{"version":1,"project":{"name":"Custom","product":"PRODUCT.md","design":"DESIGN.md","agents":"AGENTS.md","domain":"docs/domain/index.md"},"task_store":{"provider":"local"},"authority":{"human_required":["product_decision"]},"verification":{"fast":[{"name":"test","command":" "}],"complete":[{"name":"check","command":"make check"}]}}`},
		{"unknown command field", `{"version":1,"project":{"name":"Custom","product":"PRODUCT.md","design":"DESIGN.md","agents":"AGENTS.md","domain":"docs/domain/index.md"},"task_store":{"provider":"local"},"authority":{"human_required":["product_decision"]},"verification":{"fast":[{"name":"test","command":"make test","typo":true}],"complete":[{"name":"check","command":"make check"}]}}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			profilePath := filepath.Join(dir, ".sillage", "project.json")
			if err := os.MkdirAll(filepath.Dir(profilePath), 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(profilePath, []byte(test.profile), 0o600); err != nil {
				t.Fatal(err)
			}
			report := Inspect(dir)
			assertFinding(t, report, "invalid_project_profile")
		})
	}
}

func assertFinding(t *testing.T, report Report, code string) {
	t.Helper()
	for _, finding := range report.Findings {
		if finding.Code == code {
			return
		}
	}
	t.Fatalf("expected finding %q in %#v", code, report.Findings)
}

func writeRequiredFiles(t *testing.T, dir string) {
	t.Helper()
	paths := []string{"PRODUCT.md", "DESIGN.md", "AGENTS.md", "docs/domain/index.md"}
	for _, path := range paths {
		absolute := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(absolute), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(absolute, []byte("# Test\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
