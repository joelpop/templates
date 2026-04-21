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

// RemoveClaudeImport excises the bracketed block from CLAUDE.md.
// Deletes CLAUDE.md entirely if the block was its only content.
// No-op if the file or the markers don't exist.
func RemoveClaudeImport(projectRoot string) error {
	path := filepath.Join(projectRoot, "CLAUDE.md")
	existing, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	beginIdx := bytes.Index(existing, []byte(claudeImportBegin))
	endIdx := bytes.Index(existing, []byte(claudeImportEnd))
	if beginIdx < 0 || endIdx < 0 {
		return nil
	}
	endAfter := endIdx + len(claudeImportEnd)

	// Also strip the newline before begin and after end, if present,
	// so we don't leave a widening blank-line gap on repeated add/remove.
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
