package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbourmaud/sillage-workflow/internal/workflow"
)

func TestDoctorPrintsMachineReadableReport(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	var stdout bytes.Buffer
	exitCode := run([]string{"doctor", "--root", dir, "--json"}, &stdout, &stdout)

	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	var report struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("expected JSON report: %v\n%s", err, stdout.String())
	}
	if report.OK {
		t.Fatal("expected invalid project report")
	}
}

func TestTransitionValidatesTaskFileWithoutMutatingIt(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	taskPath := filepath.Join(dir, "task.json")
	task := workflow.Task{
		ID: "task-1", Title: "Test", Status: workflow.StatusDecide,
		Intent:     workflow.Intent{Outcome: "ship", Scope: []string{"core"}, NonGoals: []string{}},
		Acceptance: []workflow.AcceptanceCriterion{{ID: "AC-1", Statement: "works", Risk: "incorrect transition"}},
		Slices:     []workflow.Slice{{ID: "core", Title: "Core", Status: "active", Acceptance: []string{"AC-1"}, Dependencies: []string{}}},
	}
	task.Approval = &workflow.Approval{By: "human:owner", At: "2026-08-14T00:00:00Z", DecisionDigest: workflow.DecisionDigest(task)}
	encoded, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(taskPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	exitCode := run([]string{"transition", "--task", taskPath, "--to", "IMPLEMENT", "--json"}, &stdout, &stdout)
	if exitCode != 0 {
		t.Fatalf("expected accepted transition: %s", stdout.String())
	}
	content, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(content, encoded) {
		t.Fatal("validation command must not mutate the task record")
	}
}

func TestTransitionWriteRequiresExplicitFlagAndPersistsNewStatus(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	taskPath := filepath.Join(dir, "task.json")
	task := workflow.Task{
		ID: "task-1", Title: "Test", Status: workflow.StatusImplement,
		Intent:     workflow.Intent{Outcome: "ship", Scope: []string{"core"}, NonGoals: []string{}},
		Acceptance: []workflow.AcceptanceCriterion{{ID: "AC-1", Statement: "works", Risk: "incorrect transition"}},
		Slices:     []workflow.Slice{{ID: "core", Title: "Core", Status: "active", Acceptance: []string{"AC-1"}, Dependencies: []string{}}},
	}
	encoded, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(taskPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	if exitCode := run([]string{"transition", "--task", taskPath, "--to", "VERIFY", "--write", "--json"}, &stdout, &stdout); exitCode != 0 {
		t.Fatalf("expected explicit write to succeed: %s", stdout.String())
	}
	updated, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatal(err)
	}
	var saved workflow.Task
	if err := json.Unmarshal(updated, &saved); err != nil {
		t.Fatal(err)
	}
	if saved.Status != workflow.StatusVerify {
		t.Fatalf("expected persisted VERIFY status, got %s", saved.Status)
	}
}

func TestTransitionRejectsTaskThatViolatesPublicContract(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	taskPath := filepath.Join(dir, "task.json")
	if err := os.WriteFile(taskPath, []byte(`{"status":"DECIDE","approval":{"by":"human","at":"now"}}`), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	exitCode := run([]string{"transition", "--task", taskPath, "--to", "IMPLEMENT", "--json"}, &stdout, &stdout)
	if exitCode != 1 {
		t.Fatalf("expected malformed task rejection, got %d: %s", exitCode, stdout.String())
	}
}

func TestDigestPrintsStableDecisionFingerprintWithoutMutatingTask(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	taskPath := filepath.Join(dir, "task.json")
	task := `{"id":"task-1","title":"Test","status":"DECIDE","intent":{"outcome":"ship","scope":["core"],"non_goals":[]},"acceptance":[{"id":"AC-1","statement":"works","risk":"incorrect transition"}],"slices":[{"id":"core","title":"Core","status":"active","acceptance":["AC-1"],"dependencies":[]}]}`
	if err := os.WriteFile(taskPath, []byte(task), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	exitCode := run([]string{"digest", "--task", taskPath, "--json"}, &stdout, &stdout)
	if exitCode != 0 {
		t.Fatalf("expected digest command to succeed: %s", stdout.String())
	}
	var output struct {
		DecisionDigest string `json:"decision_digest"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("expected JSON output: %v", err)
	}
	const expected = "8232538d334742b5156a4f95a77cb98b703e23626734c074ee61c9181655fc1c"
	if output.DecisionDigest != expected {
		t.Fatalf("expected stable digest %q, got %q", expected, output.DecisionDigest)
	}
	content, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != task {
		t.Fatal("digest command must not mutate the task record")
	}
}

func TestPilotTaskCanAdvanceFromVerifyToReview(t *testing.T) {
	t.Parallel()

	taskPath := filepath.Join("..", "..", "examples", "pilot", "task.json")
	var stdout bytes.Buffer
	exitCode := run([]string{"transition", "--task", taskPath, "--to", "REVIEW", "--json"}, &stdout, &stdout)
	if exitCode != 0 {
		t.Fatalf("expected the public pilot task to pass the verification gate: %s", stdout.String())
	}
	var result workflow.TransitionResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected JSON result: %v", err)
	}
	if !result.OK || result.Code != "accepted" {
		t.Fatalf("expected accepted pilot transition, got %#v", result)
	}
}

func TestFullWorkflowPilotAdvancesOnlyInExplicitWriteCopy(t *testing.T) {
	t.Parallel()

	sourcePath := filepath.Join("..", "..", "examples", "full-workflow", "task.json")
	original, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	var readOnly bytes.Buffer
	if exitCode := run([]string{"transition", "--task", sourcePath, "--to", "HANDOFF", "--json"}, &readOnly, &readOnly); exitCode != 0 {
		t.Fatalf("expected full pilot review gate to pass: %s", readOnly.String())
	}
	afterReadOnly, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(original, afterReadOnly) {
		t.Fatal("read-only pilot transition must not mutate the example")
	}

	copyPath := filepath.Join(t.TempDir(), "task.json")
	if err := os.WriteFile(copyPath, original, 0o600); err != nil {
		t.Fatal(err)
	}
	var written bytes.Buffer
	if exitCode := run([]string{"transition", "--task", copyPath, "--to", "HANDOFF", "--write", "--json"}, &written, &written); exitCode != 0 {
		t.Fatalf("expected explicit pilot write to pass: %s", written.String())
	}
	updated, err := os.ReadFile(copyPath)
	if err != nil {
		t.Fatal(err)
	}
	var task workflow.Task
	if err := json.Unmarshal(updated, &task); err != nil {
		t.Fatal(err)
	}
	if task.Status != workflow.StatusHandoff {
		t.Fatalf("expected HANDOFF in explicit copy, got %s", task.Status)
	}
}

func TestUnknownCommandIsRejected(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	exitCode := run([]string{"publish-everything"}, &stderr, &stderr)
	if exitCode != 2 {
		t.Fatalf("expected usage error, got %d", exitCode)
	}
}

func TestChangelogCheckAndExtract(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "CHANGELOG.md")
	content := "# Changelog\n\n## [Unreleased]\n\n### Added\n\n- Next.\n\n## [0.2.0] - 2026-08-14\n\n### Added\n\n- Full workflow.\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	var checked bytes.Buffer
	if exitCode := run([]string{"changelog", "check", "--file", path, "--version", "v0.2.0"}, &checked, &checked); exitCode != 0 {
		t.Fatalf("expected changelog check to pass: %s", checked.String())
	}
	if !bytes.Contains(checked.Bytes(), []byte("0.2.0")) {
		t.Fatalf("expected version in check output, got %q", checked.String())
	}

	var notes bytes.Buffer
	if exitCode := run([]string{"changelog", "extract", "--file", path, "--version", "0.2.0"}, &notes, &notes); exitCode != 0 {
		t.Fatalf("expected changelog extraction to pass: %s", notes.String())
	}
	if notes.String() != "### Added\n\n- Full workflow.\n" {
		t.Fatalf("unexpected release notes: %q", notes.String())
	}
}

func TestContextPrintsColdStartJSON(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	for _, path := range []string{"PRODUCT.md", "DESIGN.md", "AGENTS.md", "docs/domain/index.md"} {
		absolute := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(absolute), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(absolute, []byte("# Context\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink("AGENTS.md", filepath.Join(dir, "CLAUDE.md")); err != nil {
		t.Fatal(err)
	}
	taskPath := filepath.Join(dir, "task.json")
	task := workflow.Task{
		ID: "task-1", Title: "Expose status", Status: workflow.StatusImplement,
		Intent:     workflow.Intent{Outcome: "status is visible", Scope: []string{"cli"}, NonGoals: []string{}},
		Acceptance: []workflow.AcceptanceCriterion{{ID: "AC-1", Statement: "status is visible", Risk: "hidden state"}},
		Slices:     []workflow.Slice{{ID: "slice-1", Title: "Status view", Status: "active", Acceptance: []string{"AC-1"}, Dependencies: []string{}}},
	}
	encoded, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(taskPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	exitCode := run([]string{"context", "--root", dir, "--task", taskPath, "--json"}, &stdout, &stdout)
	if exitCode != 0 {
		t.Fatalf("expected context success: %s", stdout.String())
	}
	var report struct {
		Project struct {
			Ready bool `json:"ready"`
		} `json:"project"`
		Task struct {
			Status     workflow.Status `json:"status"`
			NextStatus workflow.Status `json:"next_status"`
		} `json:"task"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("expected context JSON: %v", err)
	}
	if !report.Project.Ready || report.Task.Status != workflow.StatusImplement || report.Task.NextStatus != workflow.StatusVerify {
		t.Fatalf("unexpected context report: %#v", report)
	}
}

func TestStatusPrintsGateWithoutMutatingTask(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	taskPath := filepath.Join(dir, "task.json")
	task := workflow.Task{
		ID: "task-1", Title: "Wait", Status: workflow.StatusBlocked,
		Intent:     workflow.Intent{Outcome: "resume", Scope: []string{"adapter"}, NonGoals: []string{}},
		Acceptance: []workflow.AcceptanceCriterion{{ID: "AC-1", Statement: "ready", Risk: "blocked"}},
		Slices:     []workflow.Slice{{ID: "slice-1", Title: "Adapter", Status: "blocked", Acceptance: []string{"AC-1"}, Dependencies: []string{}}},
		Blocked:    &workflow.Blocker{From: workflow.StatusVerify, ResumeCondition: "health returns 200"},
	}
	encoded, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(taskPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	exitCode := run([]string{"status", "--task", taskPath, "--json"}, &stdout, &stdout)
	if exitCode != 0 {
		t.Fatalf("expected status success: %s", stdout.String())
	}
	var report struct {
		NextStatus workflow.Status `json:"next_status"`
		NextAction string          `json:"next_action"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("expected status JSON: %v", err)
	}
	if report.NextStatus != workflow.StatusVerify || report.NextAction != "health returns 200" {
		t.Fatalf("unexpected status report: %#v", report)
	}
	after, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, after) {
		t.Fatal("status must not mutate the task")
	}
}

func TestStatusRejectsInvalidDelegationRequest(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	taskPath := filepath.Join(dir, "task.json")
	task := `{"id":"task-1","title":"Delegation","status":"IMPLEMENT","intent":{"outcome":"ship","scope":["core"],"non_goals":[]},"delegation":{"default":{"mode":"subagent","role":"builder","isolation":"same_context","return":"implementation_patch"}},"acceptance":[{"id":"AC-1","statement":"works","risk":"regression"}],"slices":[{"id":"core","title":"Core","status":"active","acceptance":["AC-1"],"dependencies":[]}]}`
	if err := os.WriteFile(taskPath, []byte(task), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	if exitCode := run([]string{"status", "--task", taskPath, "--json"}, &stdout, &stdout); exitCode != 1 {
		t.Fatalf("expected invalid delegation request rejection, got %d: %s", exitCode, stdout.String())
	}
	var report struct {
		Valid      bool   `json:"valid"`
		NextAction string `json:"next_action"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("expected status JSON: %v", err)
	}
	if report.Valid || report.NextAction != "repair invalid task contract" {
		t.Fatalf("expected invalid task report, got %#v", report)
	}
}
