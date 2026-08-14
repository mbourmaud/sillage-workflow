package contracts

import (
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

func repositoryRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}
