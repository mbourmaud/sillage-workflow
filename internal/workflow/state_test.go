package workflow

import "testing"

func TestTransitionRequiresHumanApprovalBeforeImplementation(t *testing.T) {
	t.Parallel()

	task := Task{Status: StatusDecide}
	result := ValidateTransition(task, StatusImplement)

	if result.OK {
		t.Fatal("expected transition to be rejected without human approval")
	}
	if result.Code != "human_approval_required" {
		t.Fatalf("expected human_approval_required, got %q", result.Code)
	}
}

func TestTransitionRejectsApprovalForEarlierDecisionContent(t *testing.T) {
	t.Parallel()

	task := validTask()
	task.Status = StatusDecide
	task.Approval = &Approval{By: "human:owner", At: "2026-08-14T00:00:00Z", DecisionDigest: DecisionDigest(task)}
	task.Intent.Scope = append(task.Intent.Scope, "changed-after-approval")

	result := ValidateTransition(task, StatusImplement)
	if result.OK || result.Code != "human_approval_required" {
		t.Fatalf("expected stale approval rejection, got %#v", result)
	}
}

func TestDecisionDigestIgnoresOperationalState(t *testing.T) {
	t.Parallel()

	task := validTask()
	digest := DecisionDigest(task)
	task.Status = StatusVerify
	task.Slices[0].Status = "complete"
	if DecisionDigest(task) != digest {
		t.Fatal("operational status must not invalidate an unchanged decision")
	}
}

func TestTransitionRejectsEmptyApprovalBeforeImplementation(t *testing.T) {
	t.Parallel()

	task := Task{Status: StatusDecide, Approval: &Approval{}}
	result := ValidateTransition(task, StatusImplement)

	if result.OK || result.Code != "human_approval_required" {
		t.Fatalf("expected empty approval to be rejected, got %#v", result)
	}
}

func TestTransitionRequiresEvidenceBeforeReview(t *testing.T) {
	t.Parallel()

	task := Task{
		Status: StatusVerify,
		Acceptance: []AcceptanceCriterion{
			{ID: "AC-1", Statement: "behavior is observable"},
		},
	}
	result := ValidateTransition(task, StatusReview)

	if result.OK {
		t.Fatal("expected transition to be rejected without acceptance evidence")
	}
	if result.Code != "acceptance_evidence_missing" {
		t.Fatalf("expected acceptance_evidence_missing, got %q", result.Code)
	}
}

func TestTransitionRequiresAtLeastOneAcceptanceCriterionBeforeReview(t *testing.T) {
	t.Parallel()

	task := Task{Status: StatusVerify}
	result := ValidateTransition(task, StatusReview)

	if result.OK || result.Code != "acceptance_criteria_missing" {
		t.Fatalf("expected missing acceptance criteria to be rejected, got %#v", result)
	}
}

func TestTransitionRejectsEmptyEvidenceBeforeReview(t *testing.T) {
	t.Parallel()

	task := Task{
		Status: StatusVerify,
		Acceptance: []AcceptanceCriterion{
			{ID: "AC-1", Statement: "behavior is observable", Evidence: []Evidence{{}}},
		},
	}
	result := ValidateTransition(task, StatusReview)

	if result.OK || result.Code != "acceptance_evidence_missing" {
		t.Fatalf("expected empty evidence to be rejected, got %#v", result)
	}
}

func TestTransitionAcceptsEvidenceBackedReview(t *testing.T) {
	t.Parallel()

	task := validTask()
	task.Status = StatusVerify
	task.Acceptance = []AcceptanceCriterion{
		{ID: "AC-1", Statement: "behavior is observable", Risk: "invisible behavior"},
	}
	digest := DecisionDigest(task)
	task.Approval = &Approval{By: "human:owner", At: "2026-08-14T00:00:00Z", DecisionDigest: digest}
	task.Acceptance = []AcceptanceCriterion{
		{
			ID: "AC-1", Statement: "behavior is observable", Risk: "invisible behavior",
			Evidence: []Evidence{{Kind: "test", Ref: "run:42", ObservedAt: "2026-08-14T00:00:00Z", Result: "passed", DecisionDigest: digest}},
		},
	}
	result := ValidateTransition(task, StatusReview)

	if !result.OK {
		t.Fatalf("expected transition to be accepted, got %q", result.Code)
	}
}

func TestTransitionRejectsEvidenceFromEarlierDecisionContent(t *testing.T) {
	t.Parallel()

	task := validTask()
	task.Status = StatusVerify
	task.Approval = &Approval{By: "human:owner", At: "2026-08-14T00:00:00Z", DecisionDigest: DecisionDigest(task)}
	task.Acceptance[0].Evidence = []Evidence{{Kind: "test", Ref: "run:old", ObservedAt: "2026-08-14T00:00:00Z", Result: "passed", DecisionDigest: "old"}}
	result := ValidateTransition(task, StatusReview)
	if result.OK || result.Code != "acceptance_evidence_missing" {
		t.Fatalf("expected stale evidence rejection, got %#v", result)
	}
}

func TestTransitionAcceptsExplicitHumanEvidenceWaiver(t *testing.T) {
	t.Parallel()

	task := validTask()
	task.Status = StatusVerify
	digest := DecisionDigest(task)
	task.Approval = &Approval{By: "human:owner", At: "2026-08-14T00:00:00Z", DecisionDigest: digest}
	task.Acceptance[0].Waiver = &Waiver{By: "human:owner", At: "2026-08-14T00:00:00Z", Reason: "environment unavailable", Impact: "runtime behavior remains unproven", DecisionDigest: digest}
	result := ValidateTransition(task, StatusReview)
	if !result.OK {
		t.Fatalf("expected explicit current waiver to pass, got %#v", result)
	}
}

func TestTransitionRejectsEvidenceObservedBeforeApproval(t *testing.T) {
	t.Parallel()

	task := validTask()
	task.Status = StatusVerify
	digest := DecisionDigest(task)
	task.Approval = &Approval{By: "human:owner", At: "2026-08-14T01:00:00Z", DecisionDigest: digest}
	task.Acceptance[0].Evidence = []Evidence{{Kind: "test", Ref: "run:old", ObservedAt: "2026-08-14T00:00:00Z", Result: "passed", DecisionDigest: digest}}
	result := ValidateTransition(task, StatusReview)
	if result.OK || result.Code != "acceptance_evidence_missing" {
		t.Fatalf("expected pre-approval evidence rejection, got %#v", result)
	}
}

func TestTaskContractRejectsBlankOptionalListMembers(t *testing.T) {
	t.Parallel()

	task := validTask()
	task.Intent.NonGoals = []string{" "}
	if result := ValidateTask(task); result.OK {
		t.Fatal("expected blank non-goal rejection")
	}
	task = validTask()
	task.Slices[0].Dependencies = []string{" "}
	if result := ValidateTask(task); result.OK {
		t.Fatal("expected blank dependency rejection")
	}
}

func TestTaskContractRequiresExplicitEmptyLists(t *testing.T) {
	t.Parallel()

	task := validTask()
	if result := ValidateTask(task); result.OK {
		t.Fatal("expected absent non-goals and dependencies to be rejected")
	}
	task.Intent.NonGoals = []string{}
	task.Slices[0].Dependencies = []string{}
	if result := ValidateTask(task); !result.OK {
		t.Fatalf("expected explicit empty lists to satisfy the public schema, got %#v", result)
	}
}

func TestInvalidLifecycleJumpIsRejected(t *testing.T) {
	t.Parallel()

	task := Task{Status: StatusIntake}
	result := ValidateTransition(task, StatusHandoff)

	if result.OK {
		t.Fatal("expected lifecycle jump to be rejected")
	}
	if result.Code != "invalid_transition" {
		t.Fatalf("expected invalid_transition, got %q", result.Code)
	}
}

func TestTaskContractRequiresIdentityIntentAcceptanceAndSlice(t *testing.T) {
	t.Parallel()

	result := ValidateTask(Task{Status: StatusDecide})
	if result.OK || result.Code != "invalid_task_contract" {
		t.Fatalf("expected invalid task contract, got %#v", result)
	}
}

func TestTaskContractRejectsUnknownStatus(t *testing.T) {
	t.Parallel()

	task := validTask()
	task.Status = Status("MAGIC")
	result := ValidateTask(task)
	if result.OK {
		t.Fatal("expected unknown task status to be rejected")
	}
}

func validTask() Task {
	return Task{
		ID:     "task-1",
		Title:  "Test",
		Status: StatusIntake,
		Intent: Intent{Outcome: "ship", Scope: []string{"core"}},
		Acceptance: []AcceptanceCriterion{
			{ID: "AC-1", Statement: "works", Risk: "incorrect behavior"},
		},
		Slices: []Slice{
			{ID: "core", Title: "Core", Status: "active", Acceptance: []string{"AC-1"}},
		},
	}
}

func TestBlockedTaskResumesOnlyToRecordedState(t *testing.T) {
	t.Parallel()

	task := Task{Status: StatusBlocked, Blocked: &Blocker{From: StatusVerify, ResumeCondition: "service is reachable"}}
	wrong := ValidateTransition(task, StatusImplement)
	if wrong.OK || wrong.Code != "invalid_resume_target" {
		t.Fatalf("expected wrong resume target rejection, got %#v", wrong)
	}
	accepted := ValidateTransition(task, StatusVerify)
	if !accepted.OK {
		t.Fatalf("expected recorded resume target, got %#v", accepted)
	}
}

func TestBlockedTaskRequiresResumeCondition(t *testing.T) {
	t.Parallel()

	task := Task{Status: StatusBlocked, Blocked: &Blocker{From: StatusVerify}}
	result := ValidateTransition(task, StatusVerify)
	if result.OK || result.Code != "resume_condition_missing" {
		t.Fatalf("expected missing resume condition rejection, got %#v", result)
	}
}

func TestEnteringBlockedRequiresMatchingOriginAndResumeCondition(t *testing.T) {
	t.Parallel()

	missing := ValidateTransition(Task{Status: StatusImplement}, StatusBlocked)
	if missing.OK || missing.Code != "blocker_metadata_missing" {
		t.Fatalf("expected blocker metadata rejection, got %#v", missing)
	}
	mismatch := ValidateTransition(Task{Status: StatusImplement, Blocked: &Blocker{From: StatusVerify, ResumeCondition: "service returns"}}, StatusBlocked)
	if mismatch.OK || mismatch.Code != "blocker_origin_mismatch" {
		t.Fatalf("expected blocker origin rejection, got %#v", mismatch)
	}
}

func TestEnteringBlockedAcceptsRecordedOriginAndResumeCondition(t *testing.T) {
	t.Parallel()

	task := Task{Status: StatusImplement, Blocked: &Blocker{From: StatusImplement, ResumeCondition: "service returns"}}
	result := ValidateTransition(task, StatusBlocked)
	if !result.OK {
		t.Fatalf("expected complete blocker record to pass, got %#v", result)
	}
}
