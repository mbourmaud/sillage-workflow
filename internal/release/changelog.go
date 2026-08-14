// Package release validates and extracts Sillage release notes.
package release

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// ErrReleaseNotFound indicates that a requested version has no changelog section.
var ErrReleaseNotFound = errors.New("release section not found")

type section struct {
	version string
	date    string
	body    string
}

var (
	sectionHeading = regexp.MustCompile(`^## \[([^\]]+)\](?: - (.+))?$`)
	semver         = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)
	linkDefinition = regexp.MustCompile(`^\[[^\]]+\]:\s+\S+\s*$`)
)

// Check validates the changelog structure and, when version is non-empty,
// requires a non-empty release section for that version. A leading v is accepted.
func Check(path string, version string) error {
	sections, err := readSections(path)
	if err != nil {
		return err
	}
	if len(sections) == 0 || sections[0].version != "Unreleased" {
		return fmt.Errorf("missing [Unreleased] section")
	}
	seen := make(map[string]struct{}, len(sections))
	for _, item := range sections {
		if _, exists := seen[item.version]; exists {
			return fmt.Errorf("duplicate release section %q", item.version)
		}
		seen[item.version] = struct{}{}
		if item.version == "Unreleased" {
			if item.date != "" {
				return fmt.Errorf("[Unreleased] must not have a date")
			}
			continue
		}
		if strings.TrimSpace(item.body) == "" {
			return fmt.Errorf("empty release section %q", item.version)
		}
	}
	if version == "" {
		return nil
	}
	normalized := normalizeVersion(version)
	for _, item := range sections {
		if item.version == normalized {
			return nil
		}
	}
	return fmt.Errorf("missing release section %q: %w", normalized, ErrReleaseNotFound)
}

// Extract returns the human-readable notes for one version without its heading.
func Extract(path string, version string) (string, error) {
	sections, err := readSections(path)
	if err != nil {
		return "", err
	}
	if err := Check(path, version); err != nil {
		return "", err
	}
	normalized := normalizeVersion(version)
	for _, item := range sections {
		if item.version == normalized {
			return strings.TrimSpace(item.body) + "\n", nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrReleaseNotFound, normalized)
}

func readSections(path string) ([]section, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
	starts := make([]int, 0)
	sections := make([]section, 0)
	for index, line := range lines {
		if !strings.HasPrefix(line, "## ") {
			continue
		}
		match := sectionHeading.FindStringSubmatch(line)
		if match == nil {
			return nil, fmt.Errorf("invalid release heading on line %d", index+1)
		}
		version := match[1]
		date := match[2]
		if version != "Unreleased" {
			if !semver.MatchString(version) {
				return nil, fmt.Errorf("invalid release version %q on line %d", version, index+1)
			}
			if date == "" {
				return nil, fmt.Errorf("invalid release date for %q on line %d", version, index+1)
			}
			if _, err := time.Parse("2006-01-02", date); err != nil {
				return nil, fmt.Errorf("invalid release date for %q on line %d", version, index+1)
			}
		}
		starts = append(starts, index)
		sections = append(sections, section{version: version, date: date})
	}
	for index, start := range starts {
		end := len(lines)
		if index+1 < len(starts) {
			end = starts[index+1]
		}
		sections[index].body = trimLinkDefinitions(strings.Join(lines[start+1:end], "\n"))
	}
	return sections, nil
}

func trimLinkDefinitions(body string) string {
	lines := strings.Split(body, "\n")
	for len(lines) > 0 {
		last := strings.TrimSpace(lines[len(lines)-1])
		if last == "" || linkDefinition.MatchString(last) {
			lines = lines[:len(lines)-1]
			continue
		}
		break
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func normalizeVersion(version string) string {
	return strings.TrimPrefix(strings.TrimSpace(version), "v")
}
