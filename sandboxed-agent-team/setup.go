package main

import (
	"fmt"
	"os"
)

const setupUsage = `setup — bootstrap or update the Claude Code agent team kit

Usage:
  sandboxed-agent-team setup [--remove]

On a project with no kit installed, sets up the kit end-to-end and
onboards the current user.

On a project that already has a kit installed, regenerates all
generated files from the current templates and the project's
persisted variables, preserves and reconciles the variables file
(amending for new placeholders, removing orphans), then re-onboards
the current user. Sample files (e.g., docs/INDEX.md) are left
untouched.

Flags:
  --remove  Uninstall the kit from this project. Tears down the
            sandbox, deletes generated files, excises the
            @CLAUDE_TEAM.md import line from CLAUDE.md, and
            removes the current developer's local state. Does
            NOT touch docs/ — that belongs to the project.
`

func runSetup(args []string) int {
	remove := false
	for _, arg := range args {
		switch arg {
		case "--remove":
			remove = true
		case "--help", "-h":
			fmt.Print(setupUsage)
			return 0
		default:
			fmt.Fprintf(os.Stderr, "setup: unknown flag %q\n\n%s", arg, setupUsage)
			return 2
		}
	}

	if remove {
		return runSetupRemove()
	}
	return runSetupInstallOrUpdate()
}

func runSetupInstallOrUpdate() int {
	// Detects project state (virgin vs. existing kit) and dispatches to
	// the appropriate flow (F1 or F2). See the plan's "Flow — F1" and
	// "Flow — F2" sections.
	//
	// TODO: implement
	fmt.Fprintln(os.Stderr, "setup: not yet implemented")
	return 1
}

func runSetupRemove() int {
	// F9 remove-setup: tear down sandbox, delete generated kit files,
	// excise @CLAUDE_TEAM.md import from CLAUDE.md, delete
	// developer-local state, commit the deletion. Leaves docs/ alone.
	//
	// TODO: implement
	fmt.Fprintln(os.Stderr, "setup --remove: not yet implemented")
	return 1
}
