package workflow

import "sort"

const (
	// ClassificationProbe describes exploratory work whose primary lens is none.
	ClassificationProbe = "probe"
	// ClassificationBounded describes one bounded vertical change.
	ClassificationBounded = "bounded"
	// ClassificationCrossCutting describes work that changes several boundaries.
	ClassificationCrossCutting = "cross-cutting"

	// LensNone makes an explicit choice that no specialist lens owns the risk.
	LensNone = "none"
	// LensArchitecture owns boundaries, patterns, and dependency direction.
	LensArchitecture = "architecture"
	// LensTesting owns behavioral seams and proof strategy.
	LensTesting = "testing"
	// LensDDD owns domain language, invariants, and bounded contexts.
	LensDDD = "ddd"
	// LensSolid owns maintainable design and responsibility boundaries.
	LensSolid = "solid"
	// LensInterface owns HTTP, web, RPC, messaging, and compatibility risks.
	LensInterface = "interface"
	// LensSystems owns networks, concurrency, timing, and resource risks.
	LensSystems = "systems"
	// LensSecurity owns assets, trust boundaries, abuse cases, and controls.
	LensSecurity = "security"
	// LensPlatform owns desktop, mobile, CLI, service, and packaging concerns.
	LensPlatform = "platform"
	// LensFrontendArchitecture owns UI, route, accessibility, and browser seams.
	LensFrontendArchitecture = "frontend-architecture"
	// LensRelationalData owns relational constraints, transactions, and migrations.
	LensRelationalData = "relational-data"
	// LensDocumentData owns document access patterns and schema evolution.
	LensDocumentData = "document-data"
	// LensAudit owns evidence-led assessment of code, architecture, and debt.
	LensAudit = "audit"
	// LensMigrate owns staged and observable legacy migrations.
	LensMigrate = "migrate"
	// LensDebug owns reproducible diagnosis and bounded fixes.
	LensDebug = "debug"
	// LensTestHygiene owns proof portfolio clarity and duplication control.
	LensTestHygiene = "test-hygiene"
)

// ConformanceFinding describes one deterministic contract violation.
type ConformanceFinding struct {
	Code   string `json:"code"`
	Path   string `json:"path,omitempty"`
	Detail string `json:"detail"`
}

// ConformanceReport is the machine-readable result of checking a task context.
type ConformanceReport struct {
	OK       bool                 `json:"ok"`
	Findings []ConformanceFinding `json:"findings"`
}

var validClassifications = map[string]struct{}{
	ClassificationProbe: {}, ClassificationBounded: {}, ClassificationCrossCutting: {},
}

var validLenses = map[string]struct{}{
	LensNone: {}, LensArchitecture: {}, LensTesting: {}, LensDDD: {}, LensSolid: {},
	LensInterface: {}, LensSystems: {}, LensSecurity: {}, LensPlatform: {},
	LensFrontendArchitecture: {}, LensRelationalData: {}, LensDocumentData: {},
	LensAudit: {}, LensMigrate: {}, LensDebug: {}, LensTestHygiene: {},
}

// Conformance checks the opt-in engineering context required for a fully
// explicit task. It is stricter than ValidateTask, which remains compatible
// with existing 0.2 task records.
func Conformance(task Task) ConformanceReport {
	report := ConformanceReport{Findings: []ConformanceFinding{}}
	add := func(code, path, detail string) {
		report.Findings = append(report.Findings, ConformanceFinding{Code: code, Path: path, Detail: detail})
	}

	if result := ValidateTask(task); !result.OK {
		add("task_invalid", "", "task does not satisfy the portable task contract")
	}
	if task.Classification == "" {
		add("classification_missing", "classification", "choose probe, bounded, or cross-cutting")
	} else if _, ok := validClassifications[task.Classification]; !ok {
		add("classification_invalid", "classification", "unknown task classification")
	}
	if task.PrimaryLens == "" {
		add("primary_lens_missing", "primary_lens", "choose one owning engineering lens")
	} else if _, ok := validLenses[task.PrimaryLens]; !ok {
		add("primary_lens_invalid", "primary_lens", "unknown engineering lens")
	}
	seen := map[string]struct{}{}
	for _, lens := range task.SecondaryLenses {
		path := "secondary_lenses"
		if lens == "" {
			add("secondary_lens_invalid", path, "secondary lenses cannot be blank")
			continue
		}
		if _, ok := validLenses[lens]; !ok || lens == LensNone {
			add("secondary_lens_invalid", path, "unknown or non-owning secondary lens")
		}
		if lens == task.PrimaryLens {
			add("primary_lens_repeated", path, "primary lens must not be repeated as a secondary lens")
		}
		if _, ok := seen[lens]; ok {
			add("secondary_lens_duplicate", path, "secondary lenses must be unique")
		}
		seen[lens] = struct{}{}
	}
	for _, criterion := range task.Acceptance {
		owner, ok := task.RiskOwners[criterion.ID]
		if !ok || blank(owner) {
			add("risk_owner_missing", "risk_owners["+criterion.ID+"]", "every acceptance criterion needs one owning test or review layer")
		}
	}
	keys := make([]string, 0, len(task.RiskOwners))
	for key := range task.RiskOwners {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	acceptanceIDs := make(map[string]struct{}, len(task.Acceptance))
	for _, criterion := range task.Acceptance {
		acceptanceIDs[criterion.ID] = struct{}{}
	}
	for _, key := range keys {
		if _, ok := acceptanceIDs[key]; !ok {
			add("risk_owner_unknown_acceptance", "risk_owners["+key+"]", "risk owner key must reference an acceptance criterion")
		}
	}
	report.OK = len(report.Findings) == 0
	return report
}
