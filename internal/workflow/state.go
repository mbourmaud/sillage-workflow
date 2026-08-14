// Package workflow defines Sillage's portable task lifecycle.
package workflow

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

// Status identifies the current lifecycle state of a task.
type Status string

const (
	StatusIntake      Status = "INTAKE"
	StatusInvestigate Status = "INVESTIGATE"
	StatusDecide      Status = "DECIDE"
	StatusImplement   Status = "IMPLEMENT"
	StatusVerify      Status = "VERIFY"
	StatusReview      Status = "REVIEW"
	StatusHandoff     Status = "HANDOFF"
	StatusBlocked     Status = "BLOCKED"
)

// Evidence identifies one durable proof for an acceptance criterion.
type Evidence struct {
	Kind           string `json:"kind"`
	Ref            string `json:"ref"`
	ObservedAt     string `json:"observed_at"`
	Result         string `json:"result"`
	DecisionDigest string `json:"decision_digest"`
}

// AcceptanceCriterion describes observable behavior and its evidence.
type AcceptanceCriterion struct {
	ID        string     `json:"id"`
	Statement string     `json:"statement"`
	Risk      string     `json:"risk"`
	Evidence  []Evidence `json:"evidence,omitempty"`
	Waiver    *Waiver    `json:"waiver,omitempty"`
}

// Waiver records a human decision to proceed without deterministic evidence.
type Waiver struct {
	By             string `json:"by"`
	At             string `json:"at"`
	Reason         string `json:"reason"`
	Impact         string `json:"impact"`
	DecisionDigest string `json:"decision_digest"`
}

// Approval records an explicit human decision.
type Approval struct {
	By             string `json:"by"`
	At             string `json:"at"`
	DecisionDigest string `json:"decision_digest"`
}

// Task contains the portable state needed to validate lifecycle transitions.
type Task struct {
	ID         string                `json:"id"`
	Title      string                `json:"title"`
	Status     Status                `json:"status"`
	Intent     Intent                `json:"intent"`
	Acceptance []AcceptanceCriterion `json:"acceptance"`
	Slices     []Slice               `json:"slices"`
	Approval   *Approval             `json:"approval,omitempty"`
	Blocked    *Blocker              `json:"blocked,omitempty"`
}

// Intent bounds the approved outcome and exclusions of a task.
type Intent struct {
	Outcome  string   `json:"outcome"`
	Scope    []string `json:"scope"`
	NonGoals []string `json:"non_goals"`
}

// Slice describes one independently understandable part of a task.
type Slice struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Status       string   `json:"status"`
	Acceptance   []string `json:"acceptance"`
	Dependencies []string `json:"dependencies"`
}

// Blocker records where a blocked task resumes and the observable condition required.
type Blocker struct {
	From            Status `json:"from"`
	ResumeCondition string `json:"resume_condition"`
}

// TransitionResult explains whether a lifecycle transition is permitted.
type TransitionResult struct {
	OK   bool   `json:"ok"`
	Code string `json:"code"`
}

var forwardTransitions = map[Status]Status{
	StatusIntake:      StatusInvestigate,
	StatusInvestigate: StatusDecide,
	StatusDecide:      StatusImplement,
	StatusImplement:   StatusVerify,
	StatusVerify:      StatusReview,
	StatusReview:      StatusHandoff,
}

// ValidateTransition checks lifecycle order, human authority, and acceptance evidence.
func ValidateTransition(task Task, target Status) TransitionResult {
	if target == StatusBlocked {
		if task.Blocked == nil || blank(task.Blocked.ResumeCondition) {
			return TransitionResult{Code: "blocker_metadata_missing"}
		}
		if task.Blocked.From != task.Status {
			return TransitionResult{Code: "blocker_origin_mismatch"}
		}
		return TransitionResult{OK: true, Code: "accepted"}
	}
	if task.Status == StatusBlocked {
		if task.Blocked == nil || blank(task.Blocked.ResumeCondition) {
			return TransitionResult{Code: "resume_condition_missing"}
		}
		if task.Blocked.From != target {
			return TransitionResult{Code: "invalid_resume_target"}
		}
		return TransitionResult{OK: true, Code: "accepted"}
	}
	if forwardTransitions[task.Status] != target {
		return TransitionResult{Code: "invalid_transition"}
	}
	if task.Status == StatusDecide && !hasApproval(task) {
		return TransitionResult{Code: "human_approval_required"}
	}
	if task.Status == StatusVerify {
		if len(task.Acceptance) == 0 {
			return TransitionResult{Code: "acceptance_criteria_missing"}
		}
		for _, criterion := range task.Acceptance {
			if !hasEvidence(task, criterion.Evidence) && !hasWaiver(task, criterion.Waiver) {
				return TransitionResult{Code: "acceptance_evidence_missing"}
			}
		}
	}
	return TransitionResult{OK: true, Code: "accepted"}
}

// ValidateTask checks the required portable task contract used by transition gates.
func ValidateTask(task Task) TransitionResult {
	if blank(task.ID) || blank(task.Title) || !validStatus(task.Status) || blank(task.Intent.Outcome) || !nonBlankStrings(task.Intent.Scope) {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	if !blankFreeStrings(task.Intent.NonGoals) {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	if len(task.Acceptance) == 0 || len(task.Slices) == 0 {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	for _, criterion := range task.Acceptance {
		if blank(criterion.ID) || blank(criterion.Statement) || blank(criterion.Risk) {
			return TransitionResult{Code: "invalid_task_contract"}
		}
		for _, evidence := range criterion.Evidence {
			if blank(evidence.Kind) || blank(evidence.Ref) || !validPastInstant(evidence.ObservedAt) || !validEvidenceResult(evidence.Result) || !validDigest(evidence.DecisionDigest) {
				return TransitionResult{Code: "invalid_task_contract"}
			}
		}
		if criterion.Waiver != nil && !validWaiverShape(criterion.Waiver) {
			return TransitionResult{Code: "invalid_task_contract"}
		}
	}
	for _, slice := range task.Slices {
		if blank(slice.ID) || blank(slice.Title) || !validSliceStatus(slice.Status) || !nonBlankStrings(slice.Acceptance) || !blankFreeStrings(slice.Dependencies) {
			return TransitionResult{Code: "invalid_task_contract"}
		}
	}
	if task.Approval != nil && (!validHumanActor(task.Approval.By) || !validPastInstant(task.Approval.At) || !validDigest(task.Approval.DecisionDigest)) {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	if task.Blocked != nil && (!validResumeStatus(task.Blocked.From) || blank(task.Blocked.ResumeCondition)) {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	if task.Status == StatusBlocked && task.Blocked == nil {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	return TransitionResult{OK: true, Code: "accepted"}
}

// DecisionDigest binds approvals and verification records to the exact approved intent and plan.
func DecisionDigest(task Task) string {
	type criterion struct{ ID, Statement, Risk string }
	type slice struct {
		ID, Title                string
		Acceptance, Dependencies []string
	}
	payload := struct {
		Intent     Intent      `json:"intent"`
		Acceptance []criterion `json:"acceptance"`
		Slices     []slice     `json:"slices"`
	}{Intent: task.Intent}
	for _, item := range task.Acceptance {
		payload.Acceptance = append(payload.Acceptance, criterion{item.ID, item.Statement, item.Risk})
	}
	for _, item := range task.Slices {
		payload.Slices = append(payload.Slices, slice{item.ID, item.Title, item.Acceptance, item.Dependencies})
	}
	encoded, _ := json.Marshal(payload)
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:])
}

func validStatus(status Status) bool {
	switch status {
	case StatusIntake, StatusInvestigate, StatusDecide, StatusImplement, StatusVerify, StatusReview, StatusHandoff, StatusBlocked:
		return true
	default:
		return false
	}
}

func hasApproval(task Task) bool {
	approval := task.Approval
	if approval == nil || !validHumanActor(approval.By) {
		return false
	}
	return validPastInstant(approval.At) && approval.DecisionDigest == DecisionDigest(task)
}

func hasEvidence(task Task, evidence []Evidence) bool {
	approvedAt, ok := approvalInstant(task)
	if !ok {
		return false
	}
	for _, item := range evidence {
		observedAt, valid := parsePastInstant(item.ObservedAt)
		if !blank(item.Kind) && !blank(item.Ref) && valid && !observedAt.Before(approvedAt) && validEvidenceResult(item.Result) && item.DecisionDigest == DecisionDigest(task) {
			return true
		}
	}
	return false
}

func hasWaiver(task Task, waiver *Waiver) bool {
	approvedAt, ok := approvalInstant(task)
	waivedAt, valid := parsePastInstant(waiverInstant(waiver))
	return ok && valid && !waivedAt.Before(approvedAt) && validWaiverShape(waiver) && waiver.DecisionDigest == DecisionDigest(task)
}

func approvalInstant(task Task) (time.Time, bool) {
	if !hasApproval(task) {
		return time.Time{}, false
	}
	return parsePastInstant(task.Approval.At)
}

func waiverInstant(waiver *Waiver) string {
	if waiver == nil {
		return ""
	}
	return waiver.At
}

func validWaiverShape(waiver *Waiver) bool {
	return waiver != nil && validHumanActor(waiver.By) && validPastInstant(waiver.At) && !blank(waiver.Reason) && !blank(waiver.Impact) && validDigest(waiver.DecisionDigest)
}

func validHumanActor(actor string) bool {
	return strings.HasPrefix(actor, "human:") && !blank(strings.TrimPrefix(actor, "human:"))
}

func blank(value string) bool {
	return strings.TrimSpace(value) == ""
}

func nonBlankStrings(values []string) bool {
	if len(values) == 0 {
		return false
	}
	for _, value := range values {
		if blank(value) {
			return false
		}
	}
	return true
}

func blankFreeStrings(values []string) bool {
	for _, value := range values {
		if blank(value) {
			return false
		}
	}
	return true
}

func validPastInstant(value string) bool {
	_, ok := parsePastInstant(value)
	return ok
}

func parsePastInstant(value string) (time.Time, bool) {
	instant, err := time.Parse(time.RFC3339, value)
	return instant, err == nil && !instant.After(time.Now().Add(5*time.Minute))
}

func validDigest(value string) bool {
	decoded, err := hex.DecodeString(value)
	return err == nil && len(decoded) == sha256.Size && value == strings.ToLower(value)
}

func validEvidenceResult(result string) bool {
	return result == "passed" || result == "observed"
}

func validSliceStatus(status string) bool {
	switch status {
	case "planned", "active", "complete", "blocked":
		return true
	default:
		return false
	}
}

func validResumeStatus(status Status) bool {
	return validStatus(status) && status != StatusBlocked
}
