package site

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var localReferencePattern = regexp.MustCompile(`(?:href|src)="([^"]+)"`)

func TestPublishedSiteHasCompleteStaticSurface(t *testing.T) {
	t.Parallel()

	root := findSiteRoot(t)
	required := []string{
		"index.html",
		"workflow.html",
		"install.html",
		"releases.html",
		"404.html",
		"assets/styles.css",
		"assets/site.js",
	}
	for _, relative := range required {
		path := filepath.Join(root, relative)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("required published site file %q is unavailable: %v", relative, err)
		}
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("read site directory: %v", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".html") {
			continue
		}
		checkHTMLPage(t, root, filepath.Join(root, entry.Name()))
	}
}

func checkHTMLPage(t *testing.T, root, page string) {
	t.Helper()

	contents, err := os.ReadFile(page)
	if err != nil {
		t.Fatalf("read %s: %v", page, err)
	}
	html := string(contents)
	for _, marker := range []string{`<html lang="en">`, "<title>", `href="assets/styles.css"`, `src="assets/site.js"`} {
		if !strings.Contains(html, marker) {
			t.Errorf("%s is missing %s", filepath.Base(page), marker)
		}
	}

	for _, match := range localReferencePattern.FindAllStringSubmatch(html, -1) {
		reference := match[1]
		if isExternalReference(reference) {
			continue
		}
		pathPart := reference
		if hash := strings.IndexAny(pathPart, "?#"); hash >= 0 {
			pathPart = pathPart[:hash]
		}
		if pathPart == "" {
			continue
		}
		candidate := filepath.Clean(filepath.Join(filepath.Dir(page), filepath.FromSlash(pathPart)))
		if !withinRoot(root, candidate) {
			t.Errorf("%s points outside published site: %q", filepath.Base(page), reference)
			continue
		}
		if _, err := os.Stat(candidate); err != nil {
			t.Errorf("%s points to missing local asset: %q (%v)", filepath.Base(page), reference, err)
		}
	}
}

func findSiteRoot(t *testing.T) string {
	t.Helper()

	directory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		candidate := filepath.Join(directory, "site", "index.html")
		if _, err := os.Stat(candidate); err == nil {
			return filepath.Join(directory, "site")
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			t.Fatal("could not find repository site directory")
		}
		directory = parent
	}
}

func withinRoot(root, path string) bool {
	relative, err := filepath.Rel(root, path)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func isExternalReference(reference string) bool {
	return strings.HasPrefix(reference, "#") ||
		strings.HasPrefix(reference, "http://") ||
		strings.HasPrefix(reference, "https://") ||
		strings.HasPrefix(reference, "mailto:") ||
		strings.HasPrefix(reference, "javascript:") ||
		strings.HasPrefix(reference, "data:")
}
