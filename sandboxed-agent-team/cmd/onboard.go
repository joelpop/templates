package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// maybeProvisionPlatformAPI writes .sandbox/platform-api.env when the
// project's merge method is PR and the file is not already present.
// Skips the prompt if the developer already has an env file (common
// when re-running onboard).
func maybeProvisionPlatformAPI(projectRoot string) error {
	envPath := filepath.Join(projectRoot, ".sandbox", "platform-api.env")
	if _, err := os.Stat(envPath); err == nil {
		return nil // already provisioned
	}

	vars, err := LoadVariables(filepath.Join(projectRoot, DefaultVariablesPath))
	if err != nil {
		return err
	}
	if vars["MERGE_METHOD"] != "PR" {
		return nil
	}

	return WritePlatformAPIEnv(projectRoot, DetectPlatform())
}

const onboardUsage = `onboard — set up developer-local state on a kit-installed project

Usage:
  sandboxed-agent-team onboard [--remove]

On a project where the kit is installed but the current developer's
local state (sandbox container, .claude/.last-onboarded, etc.) is
missing, this builds that local state. Idempotent — running again
tears down and rebuilds local state.

Flags:
  --remove  Tear down the current developer's sandbox container and
            delete per-developer local state. Leaves the project's
            versioned kit files intact.
`

func runOnboard(args []string) int {
	remove := false
	for _, arg := range args {
		switch arg {
		case "--remove":
			remove = true
		case "--help", "-h":
			fmt.Print(onboardUsage)
			return 0
		default:
			fmt.Fprintf(os.Stderr, "onboard: unknown flag %q\n\n%s", arg, onboardUsage)
			return 2
		}
	}

	if remove {
		return runOnboardRemove()
	}
	return runOnboardInstall()
}

// runOnboardInstall implements the plan's "Flow — F3":
// verify kit → tear down existing local sandbox → build + start
// sandbox → record .last-onboarded.
func runOnboardInstall() int {
	projectRoot := "."

	if !IsKitInstalled(projectRoot) {
		return fail(fmt.Errorf(
			"kit is not installed on this branch.\n" +
				"  Either you're on the wrong branch, or the kit hasn't been set up yet.\n" +
				"  Run `sandboxed-agent-team setup` to install the kit, or switch to\n" +
				"  the branch where it lives."))
	}

	// Tear down an existing sandbox container if present. Safe to call
	// even when nothing is running.
	if err := RunSandboxStop(projectRoot); err != nil {
		// Non-fatal: the sandbox may simply not be running yet.
		fmt.Printf("Warning: sandbox teardown reported: %s\n", err)
	}

	// If the project uses an SSH remote, provision the developer's
	// SSH material into .sandbox/ssh/ so start.sh can inject it.
	// No-op for HTTPS remotes.
	if sshInfo := DiscoverSSHRemote(); sshInfo != nil {
		if err := ProvisionSSH(projectRoot, sshInfo); err != nil {
			fmt.Printf("Warning: SSH provisioning failed: %s\n", err)
			fmt.Println("         Sandbox will start, but git operations over SSH may fail.")
		}
	}

	// If the project uses the PR merge method, prompt for a platform
	// API token so the Integrator can create/merge PRs via the
	// platform's REST API. Skip if a platform-api.env is already in
	// place (don't re-prompt on re-onboard).
	if err := maybeProvisionPlatformAPI(projectRoot); err != nil {
		fmt.Printf("Warning: platform API provisioning failed: %s\n", err)
		fmt.Println("         Sandbox will start, but PRs won't go through the API.")
	}

	if err := RunSandboxStart(projectRoot); err != nil {
		return fail(fmt.Errorf("start sandbox: %w", err))
	}

	if err := recordOnboarding(projectRoot); err != nil {
		return fail(fmt.Errorf("record onboarding: %w", err))
	}

	fmt.Println("Onboarding complete.")
	return 0
}

// runOnboardRemove implements `onboard --remove`: tear down the
// developer's sandbox container and delete per-developer state.
// Leaves the project's versioned kit files intact.
func runOnboardRemove() int {
	projectRoot := "."

	fmt.Println("Removing developer-local state for this project.")

	if err := RunSandboxStop(projectRoot); err != nil {
		fmt.Printf("Warning: sandbox teardown reported: %s\n", err)
	}

	if err := removeDeveloperLocalState(projectRoot); err != nil {
		return fail(err)
	}

	fmt.Println("Onboarding removed.")
	return 0
}
