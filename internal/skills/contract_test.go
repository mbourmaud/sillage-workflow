package skills

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

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

func readSkill(t *testing.T, name string) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve test path")
	}
	path := filepath.Join(filepath.Dir(filename), "..", "..", "skills", name, "SKILL.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}
