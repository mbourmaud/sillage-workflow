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

// Decision records a human decision that bounds the task.
type Decision struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Status   string `json:"status"`
	By       string `json:"by"`
	At       string `json:"at"`
}

// ReviewFinding records one issue found during independent review.
type ReviewFinding struct {
	Severity string `json:"severity"`
	Ref      string `json:"ref"`
	Detail   string `json:"detail"`
}

// Review records an independent assessment of verified work.
type Review struct {
	Status         string          `json:"status"`
	By             string          `json:"by"`
	At             string          `json:"at"`
	Summary        string          `json:"summary"`
	Findings       []ReviewFinding `json:"findings,omitempty"`
	DecisionDigest string          `json:"decision_digest"`
}

// Handoff records the outcome and next safe action for a fresh agent.
type Handoff struct {
	Outcome        string `json:"outcome"`
	NextAction     string `json:"next_action"`
	At             string `json:"at"`
	DecisionDigest string `json:"decision_digest"`
}

// Task contains the portable state needed to validate lifecycle transitions.
type Task struct {
	ID         string                `json:"id"`
	Title      string                `json:"title"`
	Status     Status                `json:"status"`
	Intent     Intent                `json:"intent"`
	Execution  *ExecutionPlan        `json:"execution,omitempty"`
	Delegation *DelegationPlan       `json:"delegation,omitempty"`
	Decisions  []Decision            `json:"decisions,omitempty"`
	Acceptance []AcceptanceCriterion `json:"acceptance"`
	Slices     []Slice               `json:"slices"`
	Approval   *Approval             `json:"approval,omitempty"`
	Blocked    *Blocker              `json:"blocked,omitempty"`
	Review     *Review               `json:"review,omitempty"`
	Handoff    *Handoff              `json:"handoff,omitempty"`
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

// ExecutionPlan records provider-neutral reasoning requirements for a task.
// It recommends capability and effort without naming a model or provider.
type ExecutionPlan struct {
	Default   ExecutionProfile            `json:"default"`
	Overrides map[Status]ExecutionProfile `json:"overrides,omitempty"`
}

// ExecutionProfile describes the minimum reasoning capability and effort a
// task stage should receive from its configured agent adapter.
type ExecutionProfile struct {
	Capability string `json:"capability"`
	Effort     string `json:"effort"`
}

// ExecutionFor returns the requested profile for a lifecycle stage.
func (task Task) ExecutionFor(stage Status) (ExecutionProfile, bool) {
	if task.Execution == nil {
		return ExecutionProfile{}, false
	}
	if profile, ok := task.Execution.Overrides[stage]; ok {
		return profile, true
	}
	return task.Execution.Default, true
}

// DelegationPlan records provider-neutral child-agent requests by lifecycle stage.
// It describes the requested handoff, not the model or the host mechanism.
type DelegationPlan struct {
	Default   DelegationRequest            `json:"default"`
	Overrides map[Status]DelegationRequest `json:"overrides,omitempty"`
}

// DelegationRequest describes how the parent agent should handle one stage.
type DelegationRequest struct {
	Mode      string `json:"mode"`
	Role      string `json:"role"`
	Isolation string `json:"isolation"`
	Return    string `json:"return"`
	Required  bool   `json:"required,omitempty"`
}

// DelegationFor returns the requested delegation policy for a lifecycle stage.
func (task Task) DelegationFor(stage Status) (DelegationRequest, bool) {
	if task.Delegation == nil {
		return DelegationRequest{}, false
	}
	if request, ok := task.Delegation.Overrides[stage]; ok {
		return request, true
	}
	return task.Delegation.Default, true
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
	if task.Status == StatusReview && target == StatusHandoff {
		if !hasAcceptedReview(task) {
			if hasBlockingReviewFinding(task.Review) {
				return TransitionResult{Code: "review_blocking_findings"}
			}
			return TransitionResult{Code: "review_required"}
		}
		if !hasCurrentHandoff(task) {
			return TransitionResult{Code: "handoff_required"}
		}
	}
	return TransitionResult{OK: true, Code: "accepted"}
}

// ValidateTask checks the required portable task contract used by transition gates.
func ValidateTask(task Task) TransitionResult {
	if blank(task.ID) || blank(task.Title) || !validStatus(task.Status) || blank(task.Intent.Outcome) || !nonBlankStrings(task.Intent.Scope) {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	if task.Intent.NonGoals == nil || !blankFreeStrings(task.Intent.NonGoals) {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	if len(task.Acceptance) == 0 || len(task.Slices) == 0 {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	if !validExecutionPlan(task.Execution) {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	if !validDelegationPlan(task.Delegation) {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	for _, decision := range task.Decisions {
		if !validDecisionShape(decision) {
			return TransitionResult{Code: "invalid_task_contract"}
		}
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
		if blank(slice.ID) || blank(slice.Title) || !validSliceStatus(slice.Status) || !nonBlankStrings(slice.Acceptance) || slice.Dependencies == nil || !blankFreeStrings(slice.Dependencies) {
			return TransitionResult{Code: "invalid_task_contract"}
		}
	}
	if task.Approval != nil && (!validHumanActor(task.Approval.By) || !validPastInstant(task.Approval.At) || !validDigest(task.Approval.DecisionDigest)) {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	if task.Blocked != nil && (!validResumeStatus(task.Blocked.From) || blank(task.Blocked.ResumeCondition)) {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	if task.Review != nil && !validReviewShape(task.Review) {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	if task.Handoff != nil && !validHandoffShape(task.Handoff) {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	if task.Status == StatusBlocked && task.Blocked == nil {
		return TransitionResult{Code: "invalid_task_contract"}
	}
	return TransitionResult{OK: true, Code: "accepted"}
}

// DecisionDigest binds approvals and verification records to the exact approved intent and plan.
func DecisionDigest(task Task) string {
	type criterion struct {
		ID        string `json:"id"`
		Statement string `json:"statement"`
		Risk      string `json:"risk"`
	}
	type slice struct {
		ID           string   `json:"id"`
		Title        string   `json:"title"`
		Acceptance   []string `json:"acceptance"`
		Dependencies []string `json:"dependencies"`
	}
	type decision struct {
		ID       string `json:"id"`
		Question string `json:"question"`
		Answer   string `json:"answer"`
		Status   string `json:"status"`
		By       string `json:"by"`
		At       string `json:"at"`
	}
	payload := struct {
		Intent     Intent      `json:"intent"`
		Decisions  []decision  `json:"decisions,omitempty"`
		Acceptance []criterion `json:"acceptance"`
		Slices     []slice     `json:"slices"`
	}{Intent: task.Intent}
	for _, item := range task.Decisions {
		payload.Decisions = append(payload.Decisions, decision{item.ID, item.Question, item.Answer, item.Status, item.By, item.At})
	}
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

func validExecutionPlan(plan *ExecutionPlan) bool {
	if plan == nil {
		return true
	}
	if !validExecutionProfile(plan.Default) {
		return false
	}
	for stage, profile := range plan.Overrides {
		if !validStatus(stage) || !validExecutionProfile(profile) {
			return false
		}
	}
	return true
}

func validExecutionProfile(profile ExecutionProfile) bool {
	return validExecutionCapability(profile.Capability) && validExecutionEffort(profile.Effort)
}

func validExecutionCapability(capability string) bool {
	switch capability {
	case "light", "standard", "advanced", "frontier":
		return true
	default:
		return false
	}
}

func validExecutionEffort(effort string) bool {
	switch effort {
	case "low", "medium", "high", "max":
		return true
	default:
		return false
	}
}

func validDelegationPlan(plan *DelegationPlan) bool {
	if plan == nil {
		return true
	}
	if !validDelegationRequest(plan.Default) {
		return false
	}
	for stage, request := range plan.Overrides {
		if !validStatus(stage) || !validDelegationRequest(request) {
			return false
		}
	}
	return true
}

func validDelegationRequest(request DelegationRequest) bool {
	if !validDelegationMode(request.Mode) || !validDelegationRole(request.Role) ||
		!validDelegationIsolation(request.Isolation) || !validDelegationReturn(request.Return) {
		return false
	}
	if request.Mode == "parent" {
		return request.Isolation == "same_context" && !request.Required
	}
	return request.Isolation != "same_context"
}

func validDelegationMode(mode string) bool {
	return mode == "parent" || mode == "subagent"
}

func validDelegationRole(role string) bool {
	if role == "" {
		return false
	}
	for index, character := range role {
		if (character >= 'a' && character <= 'z') ||
			(character >= '0' && character <= '9' && index > 0) ||
			(character == '-' || character == '_') {
			continue
		}
		return false
	}
	return role[0] >= 'a' && role[0] <= 'z'
}

func validDelegationIsolation(isolation string) bool {
	switch isolation {
	case "same_context", "read_only", "isolated_worktree":
		return true
	default:
		return false
	}
}

func validDelegationReturn(result string) bool {
	switch result {
	case "summary", "decision_packet", "implementation_patch", "verification_evidence", "review_findings", "handoff_packet":
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

func validDecisionShape(decision Decision) bool {
	return !blank(decision.ID) && !blank(decision.Question) && !blank(decision.Answer) &&
		(decision.Status == "accepted" || decision.Status == "rejected") &&
		validHumanActor(decision.By) && validPastInstant(decision.At)
}

func validReviewShape(review *Review) bool {
	if review == nil || (review.Status != "accepted" && review.Status != "changes_requested") ||
		!validReviewer(review.By) || !validPastInstant(review.At) || blank(review.Summary) ||
		!validDigest(review.DecisionDigest) {
		return false
	}
	for _, finding := range review.Findings {
		if (finding.Severity != "blocking" && finding.Severity != "non_blocking") ||
			blank(finding.Ref) || blank(finding.Detail) {
			return false
		}
	}
	return true
}

func validHandoffShape(handoff *Handoff) bool {
	return handoff != nil && !blank(handoff.Outcome) && !blank(handoff.NextAction) &&
		validPastInstant(handoff.At) && validDigest(handoff.DecisionDigest)
}

func validReviewer(reviewer string) bool {
	return (strings.HasPrefix(reviewer, "reviewer:") || strings.HasPrefix(reviewer, "human:")) &&
		!blank(strings.TrimPrefix(strings.TrimPrefix(reviewer, "reviewer:"), "human:"))
}

func hasBlockingReviewFinding(review *Review) bool {
	if review == nil {
		return false
	}
	for _, finding := range review.Findings {
		if finding.Severity == "blocking" {
			return true
		}
	}
	return false
}

func hasAcceptedReview(task Task) bool {
	return task.Review != nil && task.Review.Status == "accepted" &&
		validReviewShape(task.Review) && task.Review.DecisionDigest == DecisionDigest(task) &&
		!hasBlockingReviewFinding(task.Review)
}

func hasCurrentHandoff(task Task) bool {
	return validHandoffShape(task.Handoff) && task.Handoff.DecisionDigest == DecisionDigest(task)
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
