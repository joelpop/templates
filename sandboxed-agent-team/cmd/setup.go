package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const installUsage = `install — install the agent team kit on a project

Usage:
  agent-team install [--help]

Fresh project → full install: identifies the development branch,
prompts for inputs, renders templates, commits the kit artifacts.

Existing installation → state-aware re-install: reconciles the
variables file against the current templates (prompting for any
new placeholders, cleaning up orphans), regenerates every
generated file, and commits the update. Sample files (e.g.,
docs/INDEX.md) are left untouched.

Install does NOT provision developer-local state (Docker sandbox,
SSH keys, platform API token). Once install finishes, run
./team/join.sh to provision your workstation and start the team.
`

const uninstallUsage = `uninstall — remove the kit from a project

Usage:
  agent-team uninstall [--help]

Stops any running sandbox, deletes the kit's generated files,
excises the @CLAUDE_TEAM.md import from CLAUDE.md and the kit's
block from .gitignore, and commits the removal. Does NOT touch
docs/ — that belongs to the project.
`

func runInstall(args []string) int {
	for _, arg := range args {
		switch arg {
		case "--help", "-h":
			fmt.Print(installUsage)
			return 0
		default:
			fmt.Fprintf(os.Stderr, "install: unknown flag %q\n\n%s", arg, installUsage)
			return 2
		}
	}

	projectRoot := "."
	if IsKitInstalled(projectRoot) {
		return runInstallUpdate(projectRoot)
	}
	return runInstallFresh(projectRoot)
}

// runInstallFresh bootstraps the kit on a project that doesn't yet
// have it. See the plan's "Flow — F1".
func runInstallFresh(projectRoot string) int {
	fmt.Println("Installing the sandboxed-agent-team kit on this project.")

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

	toStage := append([]string{DefaultVariablesPath, "CLAUDE.md", ".gitignore"}, written...)
	if err := GitAddForce(toStage...); err != nil {
		return fail(fmt.Errorf("stage kit files: %w", err))
	}
	if err := GitCommit("Install sandboxed-agent-team kit"); err != nil {
		return fail(fmt.Errorf("commit kit files: %w", err))
	}

	printJoinInstructions()
	return 0
}

// runInstallUpdate re-runs install on a project that already has the
// kit. See the plan's "Flow — F2".
func runInstallUpdate(projectRoot string) int {
	fmt.Println("Updating the sandboxed-agent-team kit on this project.")

	if err := gitPreflight(projectRoot); err != nil {
		return fail(err)
	}

	// Stop any running sandbox before overwriting scripts.
	stopSandboxIfInstalled(projectRoot)

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

	devBranch := vars["DEV_BRANCH_NAME"]
	discovered, err := buildDiscoveredMap(projectRoot, devBranch)
	if err != nil {
		return fail(err)
	}

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

	printJoinInstructions()
	return 0
}

func runUninstall(args []string) int {
	for _, arg := range args {
		switch arg {
		case "--help", "-h":
			fmt.Print(uninstallUsage)
			return 0
		default:
			fmt.Fprintf(os.Stderr, "uninstall: unknown flag %q\n\n%s", arg, uninstallUsage)
			return 2
		}
	}

	projectRoot := "."

	if !IsKitInstalled(projectRoot) {
		fmt.Println("No kit installation detected here. Nothing to remove.")
		return 0
	}

	fmt.Println("Removing the sandboxed-agent-team kit from this project.")
	ok, err := PromptYesNo("This tears down the sandbox and deletes kit-generated files. Continue?", false)
	if err != nil {
		return fail(err)
	}
	if !ok {
		return fail(fmt.Errorf("aborted by user"))
	}

	// Stop any running sandbox BEFORE we delete the team/ directory
	// (which contains stop.sh).
	stopSandboxIfInstalled(projectRoot)

	removeDeveloperLocalState(projectRoot)

	// Kit files to remove. docs/ is intentionally left untouched.
	toRemove := []string{
		"CLAUDE_TEAM.md",
		"ONBOARDING.md",
		"TEAM_GUIDE.md",
		".mcp.json",
		".claude/team-variables.yaml",
		".claude/settings.json",
		".claude/commands/team-start.md",
		".sandbox",
		"team",
	}
	if err := GitRm(toRemove...); err != nil {
		return fail(fmt.Errorf("remove kit files: %w", err))
	}

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

	if err := GitCommit("Uninstall sandboxed-agent-team kit"); err != nil {
		return fail(fmt.Errorf("commit removal: %w", err))
	}

	fmt.Println("Kit removed.")
	return 0
}

// printJoinInstructions tells the developer what to do next after
// install writes and commits the kit. Joining is a separate command
// (team/join.sh) so the install path doesn't carry dev-local logic.
func printJoinInstructions() {
	fmt.Println()
	fmt.Println("Kit installed and committed on this branch.")
	fmt.Println()
	fmt.Println("To provision your local sandbox and start the team, run:")
	fmt.Println()
	fmt.Println("    ./team/join.sh")
	fmt.Println()
}

// stopSandboxIfInstalled invokes team/stop.sh if it exists. Used
// before install-update and uninstall to clean up a running sandbox.
// Non-fatal on error.
func stopSandboxIfInstalled(projectRoot string) {
	stop := filepath.Join(projectRoot, "team", "stop.sh")
	if _, err := os.Stat(stop); err != nil {
		return
	}
	cmd := exec.Command(stop)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = projectRoot
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: team/stop.sh reported: %s\n", err)
	}
}

// removeDeveloperLocalState deletes per-developer files that aren't
// in git. Best-effort: failures are swallowed. The canonical list
// lives in team/leave.sh; this duplicates it for uninstall because
// we're about to delete leave.sh itself.
func removeDeveloperLocalState(projectRoot string) {
	paths := []string{
		".claude/.last-onboarded",
		".claude/.team-active",
		".claude/.tasks",
		".claude/.progress.md",
		".claude/.worktrees",
		".sandbox/.ssh",
		".sandbox/.ssh.source",
		".sandbox/.platform-api.env",
		".sandbox/.oauth-token",
		".sandbox/.last-directive",
	}
	for _, p := range paths {
		_ = os.RemoveAll(filepath.Join(projectRoot, p))
	}
}

// stagePath stages a path for commit, choosing `git add -f` if the
// file exists or `git rm` if it was deleted.
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

// fail prints the error to stderr and returns a non-zero exit code
// suitable for returning from main-level run* functions.
func fail(err error) int {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	return 1
}
