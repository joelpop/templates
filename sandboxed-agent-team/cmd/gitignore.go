package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

// The kit manages a bracketed block in the target project's .gitignore
// so developer-local artifacts aren't committed. --remove excises the
// block.
const (
	gitignoreBegin = "# BEGIN sandboxed-agent-team"
	gitignoreEnd   = "# END sandboxed-agent-team"
)

// kitGitignoreEntries lists developer-local paths the kit creates that
// must not be committed to the target project's repo.
var kitGitignoreEntries = []string{
	".sandbox/.repo-platform-api.env",
	".sandbox/.oauth-token",
	".sandbox/.last-directive",
	".claude/.last-onboarded",
	".claude/.team-active",
	".claude/.tasks/",
	".claude/.progress.md",
	".claude/.worktrees/",
	".claude/.cost-log.md",
}

// EnsureKitGitignore appends the kit's developer-local entries to the
// target project's .gitignore inside a marker-bracketed block. Creates
// .gitignore if absent. Idempotent: if the begin marker is already
// present, does nothing. Returns true if the file was modified.
func EnsureKitGitignore(projectRoot string) (bool, error) {
	path := filepath.Join(projectRoot, ".gitignore")
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	if bytes.Contains(existing, []byte(gitignoreBegin)) {
		return false, nil
	}

	var block strings.Builder
	block.WriteString(gitignoreBegin)
	block.WriteString("\n# Developer-local artifacts — do not commit. Managed by `sandboxed-agent-team`.\n")
	for _, entry := range kitGitignoreEntries {
		block.WriteString(entry)
		block.WriteByte('\n')
	}
	block.WriteString(gitignoreEnd)
	block.WriteByte('\n')

	var out []byte
	switch {
	case len(existing) == 0:
		out = []byte(block.String())
	case bytes.HasSuffix(existing, []byte("\n\n")):
		out = append(existing, []byte(block.String())...)
	case bytes.HasSuffix(existing, []byte("\n")):
		out = append(existing, '\n')
		out = append(out, block.String()...)
	default:
		out = append(existing, '\n', '\n')
		out = append(out, block.String()...)
	}
	return true, os.WriteFile(path, out, 0o644)
}

// Note: gitignore block REMOVAL lives in team/uninstall.sh (bash,
// committed to target projects). See claude.go for the rationale.
