package main

import (
	"fmt"
	"os"
)

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

func runOnboardInstall() int {
	// F3: verify kit artifacts exist on current branch, tear down
	// existing local sandbox if any, write developer-local files,
	// build local sandbox, record .claude/.last-onboarded.
	//
	// TODO: implement
	fmt.Fprintln(os.Stderr, "onboard: not yet implemented")
	return 1
}

func runOnboardRemove() int {
	// F9 remove-onboarding: tear down developer's sandbox container,
	// delete developer-local state. Leaves versioned kit files alone.
	//
	// TODO: implement
	fmt.Fprintln(os.Stderr, "onboard --remove: not yet implemented")
	return 1
}
