package workflow

import "testing"

func TestConformanceAcceptsACompleteEngineeringContext(t *testing.T) {
	t.Parallel()

	task := conformanceTask()
	task.Classification = "bounded"
	task.PrimaryLens = "interface"
	task.SecondaryLenses = []string{"systems", "security"}
	task.RiskOwners = map[string]string{"AC-1": "integration-test"}

	report := Conformance(task)
	if !report.OK {
		t.Fatalf("expected complete context to conform, got %#v", report.Findings)
	}
}

func TestConformanceReportsMissingContextAndUnknownRiskOwner(t *testing.T) {
	t.Parallel()

	task := conformanceTask()
	task.Classification = "bogus"
	task.PrimaryLens = "not-a-lens"
	task.RiskOwners = map[string]string{"AC-9": "test"}

	report := Conformance(task)
	if report.OK {
		t.Fatal("expected incomplete context to fail conformance")
	}
	for _, code := range []string{"classification_invalid", "primary_lens_invalid", "risk_owner_unknown_acceptance", "risk_owner_missing"} {
		if !hasConformanceCode(report, code) {
			t.Fatalf("expected finding %q, got %#v", code, report.Findings)
		}
	}
}

func TestConformanceRejectsOverlappingLenses(t *testing.T) {
	t.Parallel()

	task := conformanceTask()
	task.Classification = "cross-cutting"
	task.PrimaryLens = "interface"
	task.SecondaryLenses = []string{"interface", "systems", "systems"}
	task.RiskOwners = map[string]string{"AC-1": "integration-test"}

	report := Conformance(task)
	if report.OK {
		t.Fatal("expected overlapping lenses to fail conformance")
	}
	for _, code := range []string{"primary_lens_repeated", "secondary_lens_duplicate"} {
		if !hasConformanceCode(report, code) {
			t.Fatalf("expected finding %q, got %#v", code, report.Findings)
		}
	}
}

func TestDecisionDigestChangesWhenEngineeringContextChanges(t *testing.T) {
	t.Parallel()

	task := conformanceTask()
	task.Classification = "bounded"
	task.PrimaryLens = "interface"
	task.RiskOwners = map[string]string{"AC-1": "integration-test"}
	digest := DecisionDigest(task)
	task.PrimaryLens = "systems"
	if DecisionDigest(task) == digest {
		t.Fatal("engineering context changes must invalidate the decision digest")
	}
}

func conformanceTask() Task {
	task := validTask()
	task.Intent.NonGoals = []string{}
	task.Slices[0].Dependencies = []string{}
	return task
}

func hasConformanceCode(report ConformanceReport, code string) bool {
	for _, finding := range report.Findings {
		if finding.Code == code {
			return true
		}
	}
	return false
}
