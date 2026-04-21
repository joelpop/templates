package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSubstitute(t *testing.T) {
	input := []byte("Java <JAVA_VERSION> with Vaadin <VAADIN_VERSION>")
	vars := Variables{
		"JAVA_VERSION":   "21",
		"VAADIN_VERSION": "24.5.0",
	}
	got, err := Substitute(input, vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "Java 21 with Vaadin 24.5.0"
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSubstituteMissing(t *testing.T) {
	input := []byte("<ABSENT> here")
	vars := Variables{}
	_, err := Substitute(input, vars)
	if err == nil {
		t.Fatal("expected error for missing variable")
	}
	if !strings.Contains(err.Error(), "ABSENT") {
		t.Errorf("error should name the missing placeholder: %v", err)
	}
}

func TestSubstituteLeavesNonMatches(t *testing.T) {
	// Lowercase and mixed-case angle-bracket content must not be
	// treated as placeholders; bare identifiers without brackets too.
	input := []byte("<lowercase> <Mixed_Case> JAVA_VERSION no brackets")
	got, err := Substitute(input, Variables{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != string(input) {
		t.Errorf("got %q, want unchanged", got)
	}
}

func TestScanAndRegistryAgree(t *testing.T) {
	// Every placeholder in the bundled templates must be declared in
	// knownPlaceholders — otherwise reconciliation would fail on real
	// runs.
	found, err := ScanPlaceholders(templateFS)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(found) == 0 {
		t.Fatal("expected at least one placeholder in the bundled templates")
	}
	if err := CheckUnknownPlaceholders(found); err != nil {
		t.Error(err)
	}
}

func TestAddClaudeImportIdempotent(t *testing.T) {
	dir := t.TempDir()
	if err := AddClaudeImport(dir); err != nil {
		t.Fatalf("first add: %v", err)
	}
	if err := AddClaudeImport(dir); err != nil {
		t.Fatalf("second add: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	count := strings.Count(string(content), claudeImportBegin)
	if count != 1 {
		t.Errorf("marker count = %d, want 1", count)
	}
	if !strings.Contains(string(content), claudeImportLine) {
		t.Errorf("expected %q in %q", claudeImportLine, content)
	}
}

func TestRemoveClaudeImportPreservesUserContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CLAUDE.md")
	if err := os.WriteFile(path, []byte("User's own content\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := AddClaudeImport(dir); err != nil {
		t.Fatalf("add: %v", err)
	}
	if err := RemoveClaudeImport(dir); err != nil {
		t.Fatalf("remove: %v", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if strings.Contains(string(content), claudeImportBegin) {
		t.Errorf("marker should be removed: %q", content)
	}
	if !strings.Contains(string(content), "User's own content") {
		t.Errorf("user content must be preserved: %q", content)
	}
}

func TestRemoveClaudeImportDeletesEmptyFile(t *testing.T) {
	dir := t.TempDir()
	if err := AddClaudeImport(dir); err != nil {
		t.Fatal(err)
	}
	if err := RemoveClaudeImport(dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "CLAUDE.md")); !os.IsNotExist(err) {
		t.Errorf("CLAUDE.md should have been deleted; got err=%v", err)
	}
}

func TestVariablesRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "team-variables.yaml")
	vars := Variables{
		"JAVA_VERSION":    "21",
		"DEV_BRANCH_NAME": "feature/develop",
		"GIT_USER_NAME":   "Joel Robertson",
	}
	if err := SaveVariables(path, vars); err != nil {
		t.Fatal(err)
	}
	got, err := LoadVariables(path)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range vars {
		if got[k] != v {
			t.Errorf("key %q: got %q, want %q", k, got[k], v)
		}
	}
}

func TestLoadVariablesMissingFile(t *testing.T) {
	dir := t.TempDir()
	got, err := LoadVariables(filepath.Join(dir, "does-not-exist.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}
