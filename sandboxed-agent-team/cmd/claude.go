package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// The kit adds an @CLAUDE_TEAM.md import to the project's CLAUDE.md,
// bracketed with HTML comment markers so --remove can find and excise
// just the kit's contribution without touching user content.
const (
	claudeImportBegin = "<!-- sandboxed-agent-team: begin -->"
	claudeImportEnd   = "<!-- sandboxed-agent-team: end -->"
	claudeImportLine  = "@CLAUDE_TEAM.md"
)

// AddClaudeImport ensures CLAUDE.md at projectRoot contains the
// kit's bracketed import block. Creates CLAUDE.md if it doesn't
// exist. Idempotent: if the begin marker is already present, does
// nothing.
func AddClaudeImport(projectRoot string) error {
	path := filepath.Join(projectRoot, "CLAUDE.md")
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read %s: %w", path, err)
	}

	if bytes.Contains(existing, []byte(claudeImportBegin)) {
		return nil
	}

	block := []byte(fmt.Sprintf("%s\n%s\n%s\n",
		claudeImportBegin, claudeImportLine, claudeImportEnd))

	var out []byte
	switch {
	case len(existing) == 0:
		out = block
	case bytes.HasSuffix(existing, []byte("\n\n")):
		out = append(existing, block...)
	case bytes.HasSuffix(existing, []byte("\n")):
		out = append(existing, '\n')
		out = append(out, block...)
	default:
		out = append(existing, '\n', '\n')
		out = append(out, block...)
	}
	return os.WriteFile(path, out, 0o644)
}

// Note: CLAUDE.md import REMOVAL lives in team/uninstall.sh (bash,
// committed to target projects). It's canonically bash because a
// developer who no longer has the Go installer still needs to be
// able to uninstall.
