package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	projectRoot := "."
	if IsKitInstalled(projectRoot) {
		return runF2Existing(projectRoot)
	}
	return runF1Fresh(projectRoot)
}

// runF1Fresh implements the plan's "Flow — F1 (fresh setup)":
// preflight → dev-branch ID → branch placement → discovery → inputs
// → render → commit → onboard (which builds and starts the sandbox).
func runF1Fresh(projectRoot string) int {
	fmt.Println("Setting up the sandboxed-agent-team kit on this project.")

	if err := gitPreflight(projectRoot); err != nil {
		return fail(err)
	}

	devBranch, err := identifyDevBranch()
	if err != nil {
		return fail(err)
	}

	if _, err := placeOnBranch(devBranch, "chore/initial-setup"); err != nil {
		return fail(err)
	}

	info, err := DiscoverProject(filepath.Join(projectRoot, "pom.xml"))
	if err != nil {
		return fail(err)
	}
	discovered := info.ToVariables()
	if name, email := GitIdentity(); name != "" || email != "" {
		if name != "" {
			discovered["GIT_USER_NAME"] = name
		}
		if email != "" {
			discovered["GIT_USER_EMAIL"] = email
		}
	}
	discovered["DEV_BRANCH_NAME"] = devBranch

	vars, err := LoadVariables(filepath.Join(projectRoot, DefaultVariablesPath))
	if err != nil {
		return fail(err)
	}

	required, err := ScanPlaceholders(templateFS)
	if err != nil {
		return fail(err)
	}
	if err := CheckUnknownPlaceholders(required); err != nil {
		return fail(err)
	}

	if err := ReconcileVariables(vars, required, discovered, false); err != nil {
		return fail(err)
	}

	if err := SaveVariables(filepath.Join(projectRoot, DefaultVariablesPath), vars); err != nil {
		return fail(err)
	}

	written, err := WriteAllTemplates(projectRoot, vars)
	if err != nil {
		return fail(err)
	}

	if err := AddClaudeImport(projectRoot); err != nil {
		return fail(err)
	}

	// Stage the kit files plus the variables file plus CLAUDE.md.
	toStage := append([]string{DefaultVariablesPath, "CLAUDE.md"}, written...)
	if err := GitAddForce(toStage...); err != nil {
		return fail(fmt.Errorf("stage kit files: %w", err))
	}
	if err := GitCommit("Initial sandboxed-agent-team kit setup"); err != nil {
		return fail(fmt.Errorf("commit kit files: %w", err))
	}

	fmt.Println("\nKit files written and committed. Onboarding the current user…")
	return runOnboardInstall()
}

// runF2Existing implements re-running setup on a project that already
// has the kit installed. See plan "Flow — F2".
// TODO: implement (Task 6).
func runF2Existing(projectRoot string) int {
	fmt.Fprintln(os.Stderr, "setup (existing installation): not yet implemented")
	return 1
}

// runSetupRemove implements `setup --remove`. See plan F9.
// TODO: implement (Task 7).
func runSetupRemove() int {
	fmt.Fprintln(os.Stderr, "setup --remove: not yet implemented")
	return 1
}

// fail prints the error to stderr and returns a non-zero exit code
// suitable for returning from main-level run* functions.
func fail(err error) int {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	return 1
}
