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

func TestUnknownCommandIsRejected(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	exitCode := run([]string{"publish-everything"}, &stderr, &stderr)
	if exitCode != 2 {
		t.Fatalf("expected usage error, got %d", exitCode)
	}
}
