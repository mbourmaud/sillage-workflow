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

func TestUnknownCommandIsRejected(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	exitCode := run([]string{"publish-everything"}, &stderr, &stderr)
	if exitCode != 2 {
		t.Fatalf("expected usage error, got %d", exitCode)
	}
}
