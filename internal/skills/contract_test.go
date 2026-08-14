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

	content := readSkill(t, "researching-with-evidence")
	required := []string{
		"name: researching-with-evidence",
		"description: Use when",
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

	content := readSkill(t, "working-with-sillage")
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
		"one small task",
		"one dedicated worktree",
		"decision digest",
		"researching-with-evidence",
		"Never declare success from code inspection alone",
		"human",
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
