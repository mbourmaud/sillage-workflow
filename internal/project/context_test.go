package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mbourmaud/sillage-workflow/internal/workflow"
)

func TestContextReportsReadyProjectAndNextGate(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeRequiredFiles(t, dir)
	if err := os.Symlink("AGENTS.md", filepath.Join(dir, "CLAUDE.md")); err != nil {
		t.Fatal(err)
	}
	task := workflow.Task{
		ID: "task-1", Title: "Expose status", Status: workflow.StatusImplement,
		Intent:     workflow.Intent{Outcome: "status is visible", Scope: []string{"cli"}, NonGoals: []string{}},
		Acceptance: []workflow.AcceptanceCriterion{{ID: "AC-1", Statement: "status is visible", Risk: "hidden state"}},
		Slices:     []workflow.Slice{{ID: "slice-1", Title: "Status view", Status: "active", Acceptance: []string{"AC-1"}, Dependencies: []string{}}},
	}

	report := Context(dir, &task)
	if !report.Project.Ready || report.Project.EntryPoints.Product != "PRODUCT.md" {
		t.Fatalf("expected ready project context, got %#v", report)
	}
	if report.Task == nil || report.Task.Status != workflow.StatusImplement || report.Task.NextStatus != workflow.StatusVerify {
		t.Fatalf("expected IMPLEMENT context with VERIFY next gate, got %#v", report.Task)
	}
	if len(report.Task.ActiveSlices) != 1 || report.Task.ActiveSlices[0] != "slice-1: Status view" {
		t.Fatalf("expected active slice summary, got %#v", report.Task.ActiveSlices)
	}
}

func TestContextReportsBlockedResumeCondition(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeRequiredFiles(t, dir)
	if err := os.Symlink("AGENTS.md", filepath.Join(dir, "CLAUDE.md")); err != nil {
		t.Fatal(err)
	}
	task := workflow.Task{
		ID: "task-1", Title: "Wait for provider", Status: workflow.StatusBlocked,
		Intent:     workflow.Intent{Outcome: "provider is reachable", Scope: []string{"adapter"}, NonGoals: []string{}},
		Acceptance: []workflow.AcceptanceCriterion{{ID: "AC-1", Statement: "health passes", Risk: "unavailable provider"}},
		Slices:     []workflow.Slice{{ID: "slice-1", Title: "Provider check", Status: "blocked", Acceptance: []string{"AC-1"}, Dependencies: []string{}}},
		Blocked:    &workflow.Blocker{From: workflow.StatusVerify, ResumeCondition: "GET /health returns 200"},
	}

	report := Context(dir, &task)
	if report.Task == nil || report.Task.NextStatus != workflow.StatusVerify || report.Task.NextAction != "GET /health returns 200" {
		t.Fatalf("expected blocked resume context, got %#v", report.Task)
	}
}

func TestContextPreservesProjectFindingsWithoutInventingTask(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	report := Context(dir, nil)
	if report.Project.Ready || len(report.Project.Findings) == 0 {
		t.Fatalf("expected project findings, got %#v", report.Project)
	}
	if report.Task != nil {
		t.Fatalf("expected no invented task context, got %#v", report.Task)
	}
}
