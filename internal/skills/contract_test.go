package skills

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var portableName = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func TestEveryPublishedSkillHasValidPortableFrontmatter(t *testing.T) {
	t.Parallel()

	root := skillsRoot(t)
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		t.Run(entry.Name(), func(t *testing.T) {
			content := readSkill(t, entry.Name())
			lines := strings.Split(content, "\n")
			if len(lines) < 6 || lines[0] != "---" {
				t.Fatal("skill must start with YAML frontmatter")
			}
			end := -1
			for index := 1; index < len(lines); index++ {
				if lines[index] == "---" {
					end = index
					break
				}
			}
			if end < 3 {
				t.Fatal("skill frontmatter must be closed and contain name and description")
			}
			var name string
			var description string
			for _, line := range lines[1:end] {
				if strings.HasPrefix(line, "name: ") {
					name = strings.TrimPrefix(line, "name: ")
				}
				if strings.HasPrefix(line, "description: ") {
					description = strings.TrimPrefix(line, "description: ")
				}
			}
			if name != entry.Name() || !portableName.MatchString(name) {
				t.Fatalf("skill name %q must match its portable directory name", name)
			}
			if strings.TrimSpace(description) == "" {
				t.Fatal("skill description must not be blank")
			}
			if strings.TrimSpace(strings.Join(lines[end+1:], "\n")) == "" {
				t.Fatal("skill instructions must not be blank")
			}
		})
	}
}

func TestResearchSkillHasPortableContract(t *testing.T) {
	t.Parallel()

	content := readSkill(t, "research")
	required := []string{
		"name: research",
		"description: Use automatically when",
		"namespace: sillage",
		"qualified-name: \"sillage:research\"",
		"task record",
		"primary sources",
		"durable knowledge",
		"inference",
		"freshness",
	}
	for _, phrase := range required {
		if !strings.Contains(content, phrase) {
			t.Errorf("skill contract is missing %q", phrase)
		}
	}
	forbidden := []string{"GitHub issue", "pnpm", "Spin up a background agent"}
	for _, phrase := range forbidden {
		if strings.Contains(content, phrase) {
			t.Errorf("portable skill contains provider assumption %q", phrase)
		}
	}
}

func TestWorkflowSkillHasPortableLifecycleContract(t *testing.T) {
	t.Parallel()

	content := readSkill(t, "sillage")
	normalized := strings.Join(strings.Fields(content), " ")
	required := []string{
		"INTAKE",
		"INVESTIGATE",
		"DECIDE",
		"IMPLEMENT",
		"VERIFY",
		"REVIEW",
		"HANDOFF",
		"BLOCKED",
		"name: sillage",
		"namespace: sillage",
		"qualified-name: sillage",
		"one small",
		"clean worktree",
		"decision digest",
		"research",
		"Never declare success from code inspection alone",
		"capability",
		"effort",
		"fallback",
		"Delegation",
		"subagent",
		"isolated",
		"return shape",
		"human",
		"optional reference tooling",
		"Engineering doctrine",
		"one primary lens per risk",
		"External skills",
		"sillage:architecture",
		"sillage:testing",
		"sillage:interface",
		"sillage:systems",
		"sillage:security",
		"sillage:platform",
	}
	for _, phrase := range required {
		if !strings.Contains(normalized, phrase) {
			t.Errorf("workflow skill contract is missing %q", phrase)
		}
	}
	forbidden := []string{"superpowers", "Matt Pocock", "obra/"}
	for _, phrase := range forbidden {
		if strings.Contains(content, phrase) {
			t.Errorf("workflow skill contains competing workflow reference %q", phrase)
		}
	}
}

func TestEveryPublishedSkillAllowsImplicitInvocation(t *testing.T) {
	t.Parallel()

	root := skillsRoot(t)
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		t.Run(entry.Name(), func(t *testing.T) {
			yaml, err := os.ReadFile(filepath.Join(root, entry.Name(), "agents", "openai.yaml"))
			if err != nil {
				t.Fatalf("implicit invocation metadata is required: %v", err)
			}
			if !strings.Contains(string(yaml), "allow_implicit_invocation: true") {
				t.Fatal("published skills must permit host-driven implicit invocation")
			}
		})
	}
}

func TestSpecialistSkillsCarryPragmaticContracts(t *testing.T) {
	t.Parallel()

	required := map[string][]string{
		"architecture":          {"Pressure", "Boundary", "Option", "Cost", "Proof", "patterns", "Clean Architecture"},
		"testing":               {"Behavior", "Risk", "Seam", "Owner", "Proof", "failing test", "tautology"},
		"solid":                 {"Single Responsibility", "Open/Closed", "Liskov Substitution", "Interface Segregation", "Dependency Inversion", "KISS", "DRY", "YAGNI", "Demeter", "GRASP", "MoSCoW", "Clean Architecture", "pragmatically"},
		"ddd":                   {"ubiquitous language", "bounded contexts", "aggregates", "invariants", "domain core"},
		"interface":             {"Consumer", "Contract", "Boundary", "Failure", "Evolution", "HTTP", "compatibility"},
		"systems":               {"Failure domain", "Guarantee", "Timing", "Recovery", "Limits", "backpressure"},
		"security":              {"Asset", "Threat", "Boundary", "Control", "Residual risk", "least privilege"},
		"platform":              {"Environment", "Constraint", "Boundary", "Failure/recovery", "Proof", "desktop"},
		"frontend-architecture": {"loading", "empty", "success", "error", "Accessibility", "Browser proof"},
		"relational-data":       {"constraints", "transactions", "indexes", "Migration", "rollback"},
		"document-data":         {"access patterns", "Schema", "Consistency", "partition", "reconciliation"},
		"audit":                 {"Observed", "Risk", "Unknown", "Remediate", "Proof"},
		"migrate":               {"Current", "Target", "Compatibility", "Rollback", "Proof per stage"},
		"debug":                 {"Reproduction", "Facts", "Hypotheses", "Cause", "regression test"},
		"test-hygiene":          {"primary owning layer", "duplicates", "flaky", "coverage", "Runtime cost"},
	}
	for name, phrases := range required {
		name, phrases := name, phrases
		t.Run(name, func(t *testing.T) {
			content := strings.ToLower(readSkill(t, name))
			for _, phrase := range phrases {
				if !strings.Contains(content, strings.ToLower(phrase)) {
					t.Errorf("specialist skill is missing %q", phrase)
				}
			}
		})
	}
}

func TestUsingSillageHasSimpleEntryContract(t *testing.T) {
	t.Parallel()

	content := strings.Join(strings.Fields(readSkill(t, "using-sillage")), " ")
	for _, phrase := range []string{
		"name: using-sillage",
		"$using-sillage",
		"three-line start",
		"PRODUCT.md",
		"docs/domain/index.md",
		"state card",
		"sillage:orient",
		"sillage:shape",
		"optional",
		"Do not implement",
	} {
		if !strings.Contains(content, phrase) {
			t.Errorf("using-sillage entry contract is missing %q", phrase)
		}
	}
}

func readSkill(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join(skillsRoot(t), name, "SKILL.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func skillsRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve test path")
	}
	return filepath.Join(filepath.Dir(filename), "..", "..", "skills")
}
