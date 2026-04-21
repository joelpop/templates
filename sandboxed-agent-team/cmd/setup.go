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

	offerToJoin(projectRoot)
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

	offerToJoin(projectRoot)
	return 0
}

// runUninstall shells out to the target project's own team/uninstall.sh.
// The script is the single source of truth for the uninstall flow so a
// developer who no longer has the Go installer on their machine can
// still remove the kit. agent-team uninstall is just a convenience for
// people who do have the binary — it invokes the same script they
// would run directly.
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
	script := filepath.Join(projectRoot, "team", "uninstall.sh")
	if _, err := os.Stat(script); err != nil {
		return fail(fmt.Errorf(
			"uninstall script not found at %s.\n" +
				"  The kit may not be installed on this branch,\n" +
				"  or team/uninstall.sh was deleted.",
			script))
	}

	cmd := exec.Command(script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = projectRoot
	if err := cmd.Run(); err != nil {
		return fail(fmt.Errorf("team/uninstall.sh: %w", err))
	}
	return 0
}

// offerToJoin is the last step of install. It asks whether the
// developer wants to run team/join.sh now, and if yes, shells out.
// If no (or on any failure), prints instructions to run it later.
// Install itself carries no joining logic — it just invokes the
// same script the developer would run.
func offerToJoin(projectRoot string) {
	fmt.Println()
	fmt.Println("Kit installed and committed on this branch.")
	fmt.Println()

	ok, err := PromptYesNo("Run ./team/join.sh now to provision your workstation and start the team?", true)
	if err != nil || !ok {
		fmt.Println()
		fmt.Println("When you're ready, run:")
		fmt.Println()
		fmt.Println("    ./team/join.sh")
		fmt.Println()
		return
	}

	join := filepath.Join(projectRoot, "team", "join.sh")
	if _, err := os.Stat(join); err != nil {
		fmt.Fprintf(os.Stderr, "team/join.sh not found at %s; skipping.\n", join)
		return
	}
	cmd := exec.Command(join)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = projectRoot
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "team/join.sh exited with error: %s\n", err)
	}
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

// fail prints the error to stderr and returns a non-zero exit code
// suitable for returning from main-level run* functions.
func fail(err error) int {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	return 1
}
