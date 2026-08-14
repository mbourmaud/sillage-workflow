package project

import "github.com/mbourmaud/sillage-workflow/internal/workflow"

// ContextReport is the concise cold-start view of a project and optional task.
type ContextReport struct {
	Project ProjectContext `json:"project"`
	Task    *TaskContext   `json:"task,omitempty"`
}

// ProjectContext describes project readiness and canonical context locations.
type ProjectContext struct {
	Ready       bool        `json:"ready"`
	EntryPoints EntryPoints `json:"entry_points"`
	Findings    []Finding   `json:"findings"`
}

// TaskContext describes the current task gate without changing task state.
type TaskContext struct {
	ID           string                      `json:"id"`
	Title        string                      `json:"title"`
	Status       workflow.Status             `json:"status"`
	Execution    *workflow.ExecutionProfile  `json:"execution,omitempty"`
	Delegation   *workflow.DelegationRequest `json:"delegation,omitempty"`
	Valid        bool                        `json:"valid"`
	ActiveSlices []string                    `json:"active_slices"`
	NextStatus   workflow.Status             `json:"next_status,omitempty"`
	NextAction   string                      `json:"next_action"`
}

// Context builds a fresh-agent view without creating or mutating artifacts.
func Context(root string, task *workflow.Task) ContextReport {
	projectReport := Inspect(root)
	entries, _ := loadEntryPoints(root)
	report := ContextReport{Project: ProjectContext{
		Ready:       projectReport.OK,
		EntryPoints: entries,
		Findings:    projectReport.Findings,
	}}
	if task == nil {
		return report
	}

	taskReport := TaskStatus(*task)
	report.Task = &taskReport
	return report
}

// TaskStatus builds the next-gate view for one task without inspecting a project.
func TaskStatus(task workflow.Task) TaskContext {
	valid := workflow.ValidateTask(task).OK
	taskReport := TaskContext{
		ID:           task.ID,
		Title:        task.Title,
		Status:       task.Status,
		Valid:        valid,
		ActiveSlices: activeSliceSummaries(task),
	}
	if !valid {
		taskReport.NextAction = "repair invalid task contract"
		return taskReport
	}
	if execution, ok := task.ExecutionFor(task.Status); ok {
		taskReport.Execution = &execution
	}
	if delegation, ok := task.DelegationFor(task.Status); ok {
		taskReport.Delegation = &delegation
	}
	taskReport.NextStatus, taskReport.NextAction = nextGate(task)
	return taskReport
}

func activeSliceSummaries(task workflow.Task) []string {
	summaries := make([]string, 0)
	for _, slice := range task.Slices {
		if slice.Status == "active" || slice.Status == "blocked" {
			summaries = append(summaries, slice.ID+": "+slice.Title)
		}
	}
	return summaries
}

func nextGate(task workflow.Task) (workflow.Status, string) {
	if task.Status == workflow.StatusBlocked {
		if task.Blocked == nil {
			return "", "repair missing blocker metadata"
		}
		return task.Blocked.From, task.Blocked.ResumeCondition
	}
	var next workflow.Status
	switch task.Status {
	case workflow.StatusIntake:
		next = workflow.StatusInvestigate
	case workflow.StatusInvestigate:
		next = workflow.StatusDecide
	case workflow.StatusDecide:
		next = workflow.StatusImplement
	case workflow.StatusImplement:
		next = workflow.StatusVerify
	case workflow.StatusVerify:
		next = workflow.StatusReview
	case workflow.StatusReview:
		next = workflow.StatusHandoff
	case workflow.StatusHandoff:
		return "", "handoff complete; human controls merge and deployment"
	default:
		return "", "repair unknown task status"
	}
	if result := workflow.ValidateTransition(task, next); !result.OK {
		return next, result.Code
	}
	return next, "ready for " + string(next)
}
