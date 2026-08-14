package contracts

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/dlclark/regexp2"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type ecmaRegexp regexp2.Regexp

func (regexp *ecmaRegexp) MatchString(value string) bool {
	matched, err := (*regexp2.Regexp)(regexp).MatchString(value)
	return err == nil && matched
}

func (regexp *ecmaRegexp) String() string {
	return (*regexp2.Regexp)(regexp).String()
}

func compileECMA(pattern string) (jsonschema.Regexp, error) {
	compiled, err := regexp2.Compile(pattern, regexp2.ECMAScript)
	return (*ecmaRegexp)(compiled), err
}

func TestPublishedJSONContractsAcceptRepositoryExamples(t *testing.T) {
	t.Parallel()

	root := repositoryRoot(t)
	tests := []struct {
		name     string
		schema   string
		instance string
	}{
		{"project profile", "schemas/project.schema.json", ".sillage/project.json"},
		{"local project profile", "schemas/project.schema.json", "examples/local/project.json"},
		{"pilot task", "schemas/task.schema.json", "examples/pilot/task.json"},
		{"agent plugin", "schemas/vendor/agent-plugins/1.0.0/plugin.schema.json", "plugin.json"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			schemaPath := filepath.Join(root, test.schema)
			compiler := jsonschema.NewCompiler()
			compiler.UseRegexpEngine(compileECMA)
			compiler.AssertFormat()
			schema, err := compiler.Compile(schemaPath)
			if err != nil {
				t.Fatalf("compile %s: %v", test.schema, err)
			}
			file, err := os.Open(filepath.Join(root, test.instance))
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()
			instance, err := jsonschema.UnmarshalJSON(file)
			if err != nil {
				t.Fatalf("parse %s: %v", test.instance, err)
			}
			if err := schema.Validate(instance); err != nil {
				t.Fatalf("validate %s against %s: %v", test.instance, test.schema, err)
			}
		})
	}
}

func TestPluginBundleMatchesCanonicalSkill(t *testing.T) {
	t.Parallel()

	root := repositoryRoot(t)
	for _, relative := range []string{"SKILL.md", "agents/openai.yaml"} {
		canonical, err := os.ReadFile(filepath.Join(root, "skills", "researching-with-evidence", relative))
		if err != nil {
			t.Fatal(err)
		}
		bundled, err := os.ReadFile(filepath.Join(root, "plugins", "sillage-workflow", "skills", "researching-with-evidence", relative))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(canonical, bundled) {
			t.Fatalf("bundled %s drifted from the canonical skill", relative)
		}
	}
}

func TestClaudeMarketplacePointsToInstallableBundle(t *testing.T) {
	t.Parallel()

	root := repositoryRoot(t)
	var marketplace struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Plugins []struct {
			Name    string `json:"name"`
			Version string `json:"version"`
			Source  string `json:"source"`
		} `json:"plugins"`
	}
	decodeJSON(t, filepath.Join(root, ".claude-plugin", "marketplace.json"), &marketplace)
	if marketplace.Name != "sillage" || marketplace.Version != "0.1.0-rc.1" || len(marketplace.Plugins) != 1 {
		t.Fatalf("unexpected Claude marketplace identity, version, or plugin count")
	}
	plugin := marketplace.Plugins[0]
	if plugin.Name != "sillage-workflow" || plugin.Version != marketplace.Version || plugin.Source != "./plugins/sillage-workflow" {
		t.Fatalf("Claude marketplace does not point to the versioned Sillage plugin bundle")
	}
	if _, err := os.Stat(filepath.Join(root, plugin.Source, ".claude-plugin", "plugin.json")); err != nil {
		t.Fatalf("marketplace target is not an installable Claude plugin: %v", err)
	}
}

func TestCodexMarketplacePointsToInstallableBundle(t *testing.T) {
	t.Parallel()

	root := repositoryRoot(t)
	var marketplace struct {
		Name    string `json:"name"`
		Plugins []struct {
			Name   string `json:"name"`
			Source struct {
				Kind string `json:"source"`
				Path string `json:"path"`
			} `json:"source"`
		} `json:"plugins"`
	}
	decodeJSON(t, filepath.Join(root, ".agents", "plugins", "marketplace.json"), &marketplace)
	if marketplace.Name != "sillage" || len(marketplace.Plugins) != 1 {
		t.Fatalf("unexpected marketplace identity or plugin count")
	}
	plugin := marketplace.Plugins[0]
	if plugin.Name != "sillage-workflow" || plugin.Source.Kind != "local" || plugin.Source.Path != "./plugins/sillage-workflow" {
		t.Fatalf("marketplace does not point to the Sillage plugin bundle")
	}
	if _, err := os.Stat(filepath.Join(root, plugin.Source.Path, ".codex-plugin", "plugin.json")); err != nil {
		t.Fatalf("marketplace target is not an installable Codex plugin: %v", err)
	}
}

func TestEcosystemPluginVersionsMatch(t *testing.T) {
	t.Parallel()

	root := repositoryRoot(t)
	var portable, codex, claude struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	decodeJSON(t, filepath.Join(root, "plugin.json"), &portable)
	decodeJSON(t, filepath.Join(root, "plugins", "sillage-workflow", ".codex-plugin", "plugin.json"), &codex)
	decodeJSON(t, filepath.Join(root, "plugins", "sillage-workflow", ".claude-plugin", "plugin.json"), &claude)
	if portable.Name != codex.Name || portable.Name != claude.Name || portable.Version != codex.Version || portable.Version != claude.Version {
		t.Fatalf("plugin manifests disagree: portable=%s@%s codex=%s@%s claude=%s@%s", portable.Name, portable.Version, codex.Name, codex.Version, claude.Name, claude.Version)
	}
}

func decodeJSON(t *testing.T, path string, target any) {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(target); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}
