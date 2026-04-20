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

	discovered, err := buildDiscoveredMap(projectRoot, devBranch)
	if err != nil {
		return fail(err)
	}

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

	if _, err := EnsureKitGitignore(projectRoot); err != nil {
		return fail(fmt.Errorf("update .gitignore: %w", err))
	}

	// Stage the kit files plus the variables file, CLAUDE.md, and
	// .gitignore.
	toStage := append([]string{DefaultVariablesPath, "CLAUDE.md", ".gitignore"}, written...)
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
// has the kit installed. See plan "Flow — F2": shut down sandbox →
// reconcile variables → regenerate → commit → re-onboard.
func runF2Existing(projectRoot string) int {
	fmt.Println("Updating existing sandboxed-agent-team kit installation.")

	if err := gitPreflight(projectRoot); err != nil {
		return fail(err)
	}

	// Tear down the sandbox before overwriting its scripts.
	if err := RunSandboxStop(projectRoot); err != nil {
		fmt.Printf("Warning: sandbox teardown reported: %s\n", err)
	}

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

	// Re-discover auto-derivable values (git identity may have
	// changed on a new machine; pom.xml may reflect a version bump;
	// timestamps refresh each run). DEV_BRANCH_NAME is pulled from
	// vars, not re-discovered — it was recorded during initial setup.
	devBranch := vars["DEV_BRANCH_NAME"]
	discovered, err := buildDiscoveredMap(projectRoot, devBranch)
	if err != nil {
		return fail(err)
	}

	// Remove orphans, prompt for newly-introduced placeholders.
	if err := ReconcileVariables(vars, required, discovered, true); err != nil {
		return fail(err)
	}
	if err := SaveVariables(filepath.Join(projectRoot, DefaultVariablesPath), vars); err != nil {
		return fail(err)
	}

	written, err := WriteAllTemplates(projectRoot, vars)
	if err != nil {
		return fail(err)
	}

	// Re-adding the import and refreshing the gitignore block are
	// both idempotent; they ensure those files are correct even if
	// they were edited since initial setup.
	if err := AddClaudeImport(projectRoot); err != nil {
		return fail(err)
	}
	if _, err := EnsureKitGitignore(projectRoot); err != nil {
		return fail(fmt.Errorf("update .gitignore: %w", err))
	}

	toStage := append([]string{DefaultVariablesPath, "CLAUDE.md", ".gitignore"}, written...)
	if err := GitAddForce(toStage...); err != nil {
		return fail(fmt.Errorf("stage kit files: %w", err))
	}
	if err := GitCommit("Update sandboxed-agent-team kit"); err != nil {
		return fail(fmt.Errorf("commit kit files: %w", err))
	}

	fmt.Println("\nKit regenerated. Re-onboarding the current user…")
	return runOnboardInstall()
}

// runSetupRemove implements `setup --remove`: tear down sandbox,
// delete generated kit files (but leave docs/ alone), excise the
// CLAUDE.md import, cascade to onboard-remove, commit.
func runSetupRemove() int {
	projectRoot := "."

	if !IsKitInstalled(projectRoot) {
		fmt.Println("No kit installation detected here. Nothing to remove.")
		return 0
	}

	fmt.Println("Removing sandboxed-agent-team kit from this project.")
	ok, err := PromptYesNo("This tears down the sandbox and deletes kit-generated files. Continue?", false)
	if err != nil {
		return fail(err)
	}
	if !ok {
		return fail(fmt.Errorf("aborted by user"))
	}

	if err := RunSandboxStop(projectRoot); err != nil {
		fmt.Printf("Warning: sandbox teardown reported: %s\n", err)
	}

	// Cascade to onboarding-remove (developer-local state).
	if err := removeDeveloperLocalState(projectRoot); err != nil {
		return fail(err)
	}

	// Delete the kit's generated files. Intentionally NOT touching
	// docs/ — per plan, that belongs to the project.
	toRemove := []string{
		"CLAUDE_TEAM.md",
		"ONBOARDING.md",
		"TEAM_GUIDE.md",
		".mcp.json",
		".claude/team-variables.yaml",
		".claude/settings.json",
		".claude/commands/team-start.md",
		".sandbox",
	}
	if err := GitRm(toRemove...); err != nil {
		return fail(fmt.Errorf("remove kit files: %w", err))
	}

	// Excise the @CLAUDE_TEAM.md import block from CLAUDE.md (or
	// delete CLAUDE.md entirely if the block was its only content)
	// and the kit's block from .gitignore.
	if err := RemoveClaudeImport(projectRoot); err != nil {
		return fail(fmt.Errorf("remove CLAUDE.md import: %w", err))
	}
	if err := stagePath(projectRoot, "CLAUDE.md"); err != nil {
		fmt.Printf("Warning: couldn't stage CLAUDE.md: %s\n", err)
	}
	if err := RemoveKitGitignore(projectRoot); err != nil {
		fmt.Printf("Warning: couldn't excise kit block from .gitignore: %s\n", err)
	}
	if err := stagePath(projectRoot, ".gitignore"); err != nil {
		fmt.Printf("Warning: couldn't stage .gitignore: %s\n", err)
	}

	if err := GitCommit("Remove sandboxed-agent-team kit"); err != nil {
		return fail(fmt.Errorf("commit removal: %w", err))
	}

	fmt.Println("Kit removed.")
	return 0
}

// stagePath stages a path for commit, choosing `git add -f` if the
// file exists or `git rm` if it was deleted. Handles the case where a
// remove operation may have left the path either modified or deleted.
func stagePath(projectRoot, rel string) error {
	full := filepath.Join(projectRoot, rel)
	if _, err := os.Stat(full); err == nil {
		return GitAddForce(rel)
	} else if os.IsNotExist(err) {
		return GitRm(rel)
	} else {
		return err
	}
}

// removeDeveloperLocalState deletes the per-developer files the kit
// writes during onboarding.
func removeDeveloperLocalState(projectRoot string) error {
	// Only explicit per-developer artifact: the onboarding timestamp.
	// SSH material, OAuth tokens, and the sandbox container itself
	// are handled by RunSandboxStop.
	path := filepath.Join(projectRoot, ".claude", ".last-onboarded")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove %s: %w", path, err)
	}
	return nil
}

// fail prints the error to stderr and returns a non-zero exit code
// suitable for returning from main-level run* functions.
func fail(err error) int {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	return 1
}
