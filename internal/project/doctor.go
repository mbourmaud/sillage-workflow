// Package project validates the portable repository entry contract.
package project

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Finding describes one project-contract violation.
type Finding struct {
	Code string `json:"code"`
	Path string `json:"path"`
}

// Report is the deterministic result of a project inspection.
type Report struct {
	OK       bool      `json:"ok"`
	Findings []Finding `json:"findings"`
}

// EntryPoints identifies the canonical project context files.
type EntryPoints struct {
	Name    string `json:"name"`
	Product string `json:"product"`
	Design  string `json:"design"`
	Agents  string `json:"agents"`
	Domain  string `json:"domain"`
}

type command struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

type taskStore struct {
	Provider   string
	Extensions map[string]json.RawMessage
}

func (store *taskStore) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	provider, ok := fields["provider"]
	if !ok || json.Unmarshal(provider, &store.Provider) != nil {
		return errors.New("task_store.provider must be a string")
	}
	delete(fields, "provider")
	store.Extensions = fields
	return nil
}

type verification struct {
	Fast       []command
	Complete   []command
	Extensions map[string]json.RawMessage
}

func (verification *verification) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if raw, ok := fields["fast"]; !ok || decodeStrict(raw, &verification.Fast) != nil {
		return errors.New("verification.fast must be a command array")
	}
	if raw, ok := fields["complete"]; !ok || decodeStrict(raw, &verification.Complete) != nil {
		return errors.New("verification.complete must be a command array")
	}
	delete(fields, "fast")
	delete(fields, "complete")
	verification.Extensions = fields
	return nil
}

func decodeStrict(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("expected exactly one JSON value")
	}
	return nil
}

type profile struct {
	Version   int         `json:"version"`
	Project   EntryPoints `json:"project"`
	TaskStore taskStore   `json:"task_store"`
	Authority struct {
		HumanRequired []string `json:"human_required"`
	} `json:"authority"`
	Verification verification `json:"verification"`
}

// Inspect validates the canonical project entry points under root.
func Inspect(root string) Report {
	findings := make([]Finding, 0)
	entries, profileFinding := loadEntryPoints(root)
	if profileFinding != nil {
		findings = append(findings, *profileFinding)
	}
	requiredFiles := []struct {
		path string
		code string
	}{
		{entries.Product, "missing_product"},
		{entries.Design, "missing_design"},
		{entries.Agents, "missing_agents"},
		{entries.Domain, "missing_domain_index"},
	}
	for _, required := range requiredFiles {
		if _, err := os.Stat(filepath.Join(root, required.path)); err != nil {
			findings = append(findings, Finding{Code: required.code, Path: required.path})
		}
	}

	claudePath := filepath.Join(root, "CLAUDE.md")
	target, err := os.Readlink(claudePath)
	if err != nil || filepath.Clean(target) != filepath.Clean(entries.Agents) {
		findings = append(findings, Finding{Code: "claude_not_agents_symlink", Path: "CLAUDE.md"})
	}

	return Report{OK: len(findings) == 0, Findings: findings}
}

func loadEntryPoints(root string) (EntryPoints, *Finding) {
	defaults := EntryPoints{
		Name:    "Project",
		Product: "PRODUCT.md",
		Design:  "DESIGN.md",
		Agents:  "AGENTS.md",
		Domain:  "docs/domain/index.md",
	}
	profilePath := filepath.Join(root, ".sillage", "project.json")
	content, err := os.ReadFile(profilePath)
	if errors.Is(err, os.ErrNotExist) {
		return defaults, nil
	}
	if err != nil {
		return defaults, &Finding{Code: "invalid_project_profile", Path: ".sillage/project.json"}
	}
	var configured profile
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&configured); err != nil || decoder.Decode(&struct{}{}) != io.EOF {
		return defaults, &Finding{Code: "invalid_project_profile", Path: ".sillage/project.json"}
	}
	if !validProfile(configured) {
		return defaults, &Finding{Code: "invalid_project_profile", Path: ".sillage/project.json"}
	}
	return configured.Project, nil
}

func validProfile(configured profile) bool {
	project := configured.Project
	return configured.Version == 1 && nonBlank(project.Name) && nonBlank(project.Product) && nonBlank(project.Design) && nonBlank(project.Agents) && nonBlank(project.Domain) && nonBlank(configured.TaskStore.Provider) && uniqueNonBlank(configured.Authority.HumanRequired) && validCommands(configured.Verification.Fast) && validCommands(configured.Verification.Complete)
}

func nonBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func uniqueNonBlank(values []string) bool {
	if len(values) == 0 {
		return false
	}
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !nonBlank(value) {
			return false
		}
		if _, exists := seen[value]; exists {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}

func validCommands(commands []command) bool {
	if len(commands) == 0 {
		return false
	}
	for _, item := range commands {
		if !nonBlank(item.Name) || !nonBlank(item.Command) {
			return false
		}
	}
	return true
}
