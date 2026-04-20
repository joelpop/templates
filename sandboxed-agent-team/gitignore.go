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
	".sandbox/ssh/",
	".sandbox/ssh.source",
	".sandbox/platform-api.env",
	".sandbox/.oauth-token",
	".sandbox/.last-directive",
	".claude/.last-onboarded",
	".claude/.team-active",
	".claude/tasks/",
	".claude/progress.md",
	".claude/worktrees/",
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

// RemoveKitGitignore excises the bracketed kit block from the target
// project's .gitignore. Deletes .gitignore entirely if the block was
// its only content. No-op if the file or markers are absent.
func RemoveKitGitignore(projectRoot string) error {
	path := filepath.Join(projectRoot, ".gitignore")
	existing, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	beginIdx := bytes.Index(existing, []byte(gitignoreBegin))
	endIdx := bytes.Index(existing, []byte(gitignoreEnd))
	if beginIdx < 0 || endIdx < 0 {
		return nil
	}
	endAfter := endIdx + len(gitignoreEnd)
	if beginIdx > 0 && existing[beginIdx-1] == '\n' {
		beginIdx--
	}
	if endAfter < len(existing) && existing[endAfter] == '\n' {
		endAfter++
	}

	out := make([]byte, 0, len(existing))
	out = append(out, existing[:beginIdx]...)
	out = append(out, existing[endAfter:]...)

	if len(bytes.TrimSpace(out)) == 0 {
		return os.Remove(path)
	}
	return os.WriteFile(path, out, 0o644)
}
