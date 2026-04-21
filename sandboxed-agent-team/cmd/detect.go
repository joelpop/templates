package main

import (
	"os"
	"path/filepath"
)

// IsKitInstalled reports whether the project at projectRoot already
// has the kit set up. Uses CLAUDE_TEAM.md at the project root as the
// canonical marker.
func IsKitInstalled(projectRoot string) bool {
	_, err := os.Stat(filepath.Join(projectRoot, "CLAUDE_TEAM.md"))
	return err == nil
}

