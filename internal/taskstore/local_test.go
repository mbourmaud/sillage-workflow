package taskstore

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbourmaud/sillage-workflow/internal/workflow"
)

func TestWriteTransitionAtomicallyAdvancesValidatedTask(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "task.json")
	original := []byte(`{"id":"task-1","title":"Ship","status":"IMPLEMENT","intent":{"outcome":"ship","scope":["core"],"non_goals":[]},"acceptance":[{"id":"AC-1","statement":"works","risk":"regression"}],"slices":[{"id":"slice-1","title":"Core","status":"active","acceptance":["AC-1"],"dependencies":[]}]}`)
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}
	task := workflow.Task{
		ID: "task-1", Title: "Ship", Status: workflow.StatusImplement,
		Intent:     workflow.Intent{Outcome: "ship", Scope: []string{"core"}, NonGoals: []string{}},
		Acceptance: []workflow.AcceptanceCriterion{{ID: "AC-1", Statement: "works", Risk: "regression"}},
		Slices:     []workflow.Slice{{ID: "slice-1", Title: "Core", Status: "active", Acceptance: []string{"AC-1"}, Dependencies: []string{}}},
	}

	if err := WriteTransition(path, task, workflow.StatusVerify, original); err != nil {
		t.Fatalf("expected atomic transition write: %v", err)
	}
	updated, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(original, updated) {
		t.Fatal("expected task bytes to change after explicit write")
	}
	if !bytes.Contains(updated, []byte(`"status": "VERIFY"`)) {
		t.Fatalf("expected VERIFY status in written task: %s", updated)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("expected private task permissions, got %o", info.Mode().Perm())
	}
}

func TestWriteTransitionRejectsConcurrentModification(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "task.json")
	original := []byte(`{"id":"task-1","title":"Ship","status":"IMPLEMENT","intent":{"outcome":"ship","scope":["core"],"non_goals":[]},"acceptance":[{"id":"AC-1","statement":"works","risk":"regression"}],"slices":[{"id":"slice-1","title":"Core","status":"active","acceptance":["AC-1"],"dependencies":[]}]}`)
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(original, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	task := minimalTask(workflow.StatusImplement)
	if err := WriteTransition(path, task, workflow.StatusVerify, original); !errors.Is(err, ErrConcurrentModification) {
		t.Fatalf("expected concurrent modification error, got %v", err)
	}
}

func TestWriteTransitionRejectsSymlinkTarget(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	target := filepath.Join(dir, "real.json")
	link := filepath.Join(dir, "task.json")
	if err := os.WriteFile(target, []byte(`{}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if err := WriteTransition(link, minimalTask(workflow.StatusImplement), workflow.StatusVerify, []byte(`{}`)); !errors.Is(err, ErrSymlinkTarget) {
		t.Fatalf("expected symlink rejection, got %v", err)
	}
}

func TestWriteTransitionRevalidatesTransitionBeforeWriting(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "task.json")
	original := []byte(`{"id":"task-1"}`)
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := WriteTransition(path, minimalTask(workflow.StatusIntake), workflow.StatusHandoff, original); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("expected invalid transition error, got %v", err)
	}
}

func minimalTask(status workflow.Status) workflow.Task {
	return workflow.Task{
		ID: "task-1", Title: "Ship", Status: status,
		Intent:     workflow.Intent{Outcome: "ship", Scope: []string{"core"}, NonGoals: []string{}},
		Acceptance: []workflow.AcceptanceCriterion{{ID: "AC-1", Statement: "works", Risk: "regression"}},
		Slices:     []workflow.Slice{{ID: "slice-1", Title: "Core", Status: "active", Acceptance: []string{"AC-1"}, Dependencies: []string{}}},
	}
}
