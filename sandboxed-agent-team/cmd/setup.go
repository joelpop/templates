package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// runInstall is the binary's single entry point. agent-team-install
// has no subcommands: the bare invocation runs install (fresh or
// update, auto-detected). See main.go for the top-level usage text.
func runInstall() int {
	projectRoot := "."
	if IsKitInstalled(projectRoot) {
		return runInstallUpdate(projectRoot)
	}
	return runInstallFresh(projectRoot)
}

// runInstallFresh bootstraps the kit on a project that doesn't yet
// have it. The installer writes files only — no git operations.
// The user commits, merges, and pushes on their own schedule.
func runInstallFresh(projectRoot string) int {
	printInstallIntro(true)

	discovered, err := buildDiscoveredMap(projectRoot)
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

	seedPromptDefaultsFromDiscovered(discovered, interviewPromptDefaults...)

	if err := ReconcileVariables(vars, required, discovered, false); err != nil {
		return fail(err)
	}

	runJoin, err := PromptYesNo(
		"Run ./team/join.sh automatically after install to set up your local sandbox?", true)
	if err != nil {
		return fail(err)
	}

	ok, err := confirmProceedWithInstall(vars, true, runJoin)
	if err != nil {
		return fail(err)
	}
	if !ok {
		fmt.Println("Aborted. Nothing written to disk.")
		return 0
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

	return finishFreshInstall(projectRoot, runJoin, written)
}

// finishFreshInstall closes out the fresh-install flow. If the user
// opted in during the review, run team/join.sh automatically.
// Otherwise print a pointer so they can run it when ready. Returns
// 0 on success, 1 if team/join.sh fails — the kit files are on
// disk either way, but "install + join" as a whole did not succeed.
func finishFreshInstall(projectRoot string, runJoin bool, written []string) int {
	fmt.Println()
	printPostInstallMessage(projectRoot, written)
	fmt.Println()
	if runJoin {
		fmt.Println("Running ./team/join.sh to set up your local sandbox...")
		fmt.Println()
		if err := runJoinScript(projectRoot); err != nil {
			fmt.Fprintf(os.Stderr, "\nError: team/join.sh exited with %s\n", err)
			fmt.Fprintln(os.Stderr, "The kit files are written, but your local sandbox is not set up.")
			fmt.Fprintln(os.Stderr, "Fix the underlying issue and re-run ./team/join.sh.")
			return 1
		}
		return 0
	}
	fmt.Println("When you're ready to set up your workstation, run:")
	fmt.Println()
	fmt.Println("    ./team/join.sh")
	fmt.Println()
	return 0
}

// runInstallUpdate re-runs install on a project that already has the
// kit. Like runInstallFresh, the installer writes files only — no
// git operations. The user commits, merges, and pushes on their own
// schedule.
func runInstallUpdate(projectRoot string) int {
	printInstallIntro(false)

	// Stop any running sandbox before overwriting scripts.
	destroySandboxIfInstalled(projectRoot)

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

	discovered, err := buildDiscoveredMap(projectRoot)
	if err != nil {
		return fail(err)
	}

	seedPromptDefaultsFromDiscovered(discovered, interviewPromptDefaults...)

	if err := ReconcileVariables(vars, required, discovered, true); err != nil {
		return fail(err)
	}

	ok, err := confirmProceedWithInstall(vars, false, false)
	if err != nil {
		return fail(err)
	}
	if !ok {
		fmt.Println("Aborted. Nothing written to disk.")
		return 0
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

	return finishUpdate(projectRoot, written)
}

// finishUpdate closes out the install-update flow. If the developer
// has already joined this workstation (marker present), re-run
// team/join.sh automatically so the local sandbox picks up any kit
// changes. Otherwise, leave a pointer so they can join on their own
// schedule. Returns 0 on success, 1 if the auto-rerun of
// team/join.sh fails — the kit files are written either way, but
// the local sandbox is out of sync.
func finishUpdate(projectRoot string, written []string) int {
	fmt.Println()
	printPostInstallMessage(projectRoot, written)
	fmt.Println()
	if isWorkstationJoined(projectRoot) {
		fmt.Println("Your workstation is already provisioned — re-running")
		fmt.Println("./team/join.sh to sync the local sandbox with the")
		fmt.Println("updated kit.")
		fmt.Println()
		if err := runJoinScript(projectRoot); err != nil {
			fmt.Fprintf(os.Stderr, "\nError: team/join.sh exited with %s\n", err)
			fmt.Fprintln(os.Stderr, "The kit files are written, but your local sandbox was not resynced.")
			fmt.Fprintln(os.Stderr, "Fix the underlying issue and re-run ./team/join.sh.")
			return 1
		}
		return 0
	}
	fmt.Println("When you're ready to provision your workstation, run:")
	fmt.Println()
	fmt.Println("    ./team/join.sh")
	fmt.Println()
	return 0
}

// printPostInstallMessage tells the user what files were touched
// and — if they're in a git repo — reminds them to review + commit.
// Outside a git repo the message is shorter and git-free.
func printPostInstallMessage(projectRoot string, written []string) {
	if IsGitRepo() {
		fmt.Println("Kit files written. No git operations were performed.")
		fmt.Println()
		fmt.Println("Review with:")
		fmt.Println("    git status")
		fmt.Println("    git diff")
		fmt.Println()
		fmt.Println("Stage, commit, and push on your own schedule — the install")
		fmt.Println("touched these paths:")
	} else {
		fmt.Printf("Kit files written to %s. These paths were touched:\n", projectRoot)
	}
	for _, p := range postInstallPathList(written) {
		fmt.Printf("  %s\n", p)
	}
}

// postInstallPathList returns a de-duplicated, sorted list of paths
// the installer touched, suitable for user-facing output.
// Includes the two auxiliary files that WriteAllTemplates doesn't
// return (CLAUDE.md import block, .gitignore kit block).
func postInstallPathList(written []string) []string {
	paths := append([]string{
		DefaultVariablesPath,
		"CLAUDE.md         (import block added)",
		".gitignore        (kit block added)",
	}, written...)
	return paths
}

// printInstallIntro prints a short roadmap so the developer knows
// what to expect before any prompts start.
func printInstallIntro(fresh bool) {
	if fresh {
		fmt.Println("Installing the sandboxed-agent-team kit on this project.")
	} else {
		fmt.Println("Updating the sandboxed-agent-team kit on this project.")
	}
	fmt.Println()
	fmt.Println("Steps:")
	fmt.Println("  1. Detect and collect project information:")
	fmt.Println("     * Team's development branch name")
	fmt.Println("     * Stack details:")
	fmt.Println("       - Java version")
	fmt.Println("       - Vaadin version")
	fmt.Println("       - Spring Boot version")
	fmt.Println("       - JUnit version")
	fmt.Println("       - Database")
	fmt.Println("       - Build tool")
	fmt.Println("     * Merge method")
	fmt.Println("     * CI platform")
	fmt.Println("     * Preference to include a cost report in commit messages")
	fmt.Println("     * Preference to set up your local sandbox after install")
	fmt.Println("  2. Show all your choices for confirmation.")
	fmt.Println("  3. Write the kit files to the target directory.")
	fmt.Println("  4. Set up your local sandbox (if you opted in).")
	fmt.Println()
	fmt.Println("Nothing is written until you approve at step 2.")
	fmt.Println("Ctrl-C to abort at any time.")
	fmt.Println()
	fmt.Println("After step 3 the new files appear in the target directory. No")
	fmt.Println("git operations happen — if the directory is a git repo, use")
	fmt.Println("`git status` to see what changed and stage/commit/push when")
	fmt.Println("you're ready.")
	fmt.Println()
}

// confirmProceedWithInstall shows a summary of the resolved variables
// and asks the user to approve before anything is written to disk.
// runJoinAfterInstall is surfaced in the summary only on fresh
// installs (the only flow that asks the question); update-flow
// callers should pass false — it is ignored when fresh is false.
func confirmProceedWithInstall(vars Variables, fresh bool, runJoinAfterInstall bool) (bool, error) {
	currentBranch, _ := CurrentBranch()

	rule := strings.Repeat("─", 64)
	fmt.Println()
	fmt.Println(rule)
	fmt.Println("Ready to install. Review your choices:")
	fmt.Println()

	printKV := func(label, value string) {
		if value == "" {
			return
		}
		fmt.Printf("  %-22s %s\n", label+":", value)
	}

	printKV("Project", vars["PROJECT_NAME"])
	if currentBranch != "" {
		printKV("Current branch", currentBranch)
	}
	if d := vars["DEV_BRANCH_NAME"]; d != "" && d != currentBranch {
		printKV("Team dev branch", d)
	}
	printKV("Stack", vars["STACK_SUMMARY"])

	identity := ""
	if n := vars["GIT_USER_NAME"]; n != "" {
		identity = n
		if e := vars["GIT_USER_EMAIL"]; e != "" {
			identity += " <" + e + ">"
		}
	}
	printKV("Git identity", identity)

	printKV("Merge method", vars["MERGE_METHOD"])
	printKV("CI platform", vars["CI_PLATFORM"])
	printKV("Show cost in commit", vars["COST_IN_COMMIT"])
	printKV("Append cost to project log", vars["COST_IN_LOG"])
	if fresh {
		if runJoinAfterInstall {
			printKV("Run join after install", "yes")
		} else {
			printKV("Run join after install", "no")
		}
	}

	fmt.Println()
	if fresh {
		fmt.Println("Writes the kit files, adds the @CLAUDE_TEAM.md import to")
		fmt.Println("CLAUDE.md, and adds the kit block to .gitignore.")
		fmt.Println("No git operations — you stage and commit when ready.")
	} else {
		fmt.Println("Regenerates the kit files and refreshes the CLAUDE.md")
		fmt.Println("import and the .gitignore block.")
		fmt.Println("No git operations — you stage and commit when ready.")
	}
	fmt.Println(rule)
	fmt.Println()

	return PromptYesNo("Proceed?", true)
}

// runJoinScript execs team/join.sh with stdio inherited. A missing
// script is treated as a non-fatal condition (warning + nil).
// A non-zero exit from the script is returned to the caller.
func runJoinScript(projectRoot string) error {
	join := filepath.Join(projectRoot, "team", "join.sh")
	if _, err := os.Stat(join); err != nil {
		fmt.Fprintf(os.Stderr, "team/join.sh not found at %s; skipping.\n", join)
		return nil
	}
	cmd := exec.Command(join)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = projectRoot
	return cmd.Run()
}

// isWorkstationJoined reports whether this checkout has the local
// .claude/.last-onboarded marker (written by team/join.sh, removed
// by team/leave.sh). Used to decide whether an install-update should
// auto-re-run join to keep the sandbox in sync.
func isWorkstationJoined(projectRoot string) bool {
	_, err := os.Stat(filepath.Join(projectRoot, ".claude", ".last-onboarded"))
	return err == nil
}

// destroySandboxIfInstalled invokes team/destroy.sh --yes if it
// exists. Used before install-update so the regenerated scripts
// aren't trying to manage a sandbox created by an older version of
// the kit. Non-fatal on error; --yes bypasses the interactive
// confirmation (the user already approved the install in the
// review summary).
func destroySandboxIfInstalled(projectRoot string) {
	destroy := filepath.Join(projectRoot, "team", "destroy.sh")
	if _, err := os.Stat(destroy); err != nil {
		return
	}
	cmd := exec.Command(destroy, "--yes")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = projectRoot
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: team/destroy.sh reported: %s\n", err)
	}
}

// fail prints the error to stderr and returns a non-zero exit code
// suitable for returning from main-level run* functions.
func fail(err error) int {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	return 1
}

// seedPromptDefaultsFromDiscovered mutates the named placeholders'
// static Defaults to match the auto-discovered values, so the
// resulting user prompts offer the discovered values as their
// defaults. Used for placeholders we want to confirm interactively
// (not silently auto-fill) but still pre-populate with a best-guess.
// No-op for any name the discovered map doesn't have.
func seedPromptDefaultsFromDiscovered(discovered Variables, names ...string) {
	for _, name := range names {
		v, has := discovered[name]
		if !has || v == "" {
			continue
		}
		def, ok := knownPlaceholders[name]
		if !ok {
			continue
		}
		def.Default = v
		knownPlaceholders[name] = def
	}
}

// interviewPromptDefaults lists every placeholder whose interview
// prompt should take its default from the discovered map (rather
// than the static Default in knownPlaceholders). These are the
// user-facing values the review summary confirms: the user sees
// each one as a prompt with the auto-guessed value pre-filled,
// and can Enter through to accept or type an override.
var interviewPromptDefaults = []string{
	"DEV_BRANCH_NAME",
	"JAVA_VERSION",
	"VAADIN_VERSION",
	"SPRING_BOOT_VERSION",
	"JUNIT_VERSION",
	"DATABASE",
	"BUILD_TOOL",
}
