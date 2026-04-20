// Command sandboxed-agent-team bootstraps and manages a disciplined
// Claude Code agent team kit in a target project.
//
// See the plan at
// /Users/joel/.claude/plans/i-think-this-setup-re-setup-pure-blossom.md
// for the full requirements.
package main

import (
	"fmt"
	"os"
)

const usage = `sandboxed-agent-team — set up and manage a Claude Code agent team kit

Usage:
  sandboxed-agent-team <command> [flags]

Commands:
  setup     Bootstrap the kit on a fresh project, or update an
            existing installation. Auto-detects state.
  onboard   Set up developer-local state on a project that
            already has the kit installed.

Flags:
  --remove  Destructive counterpart. With setup, uninstalls the
            kit from the project. With onboard, removes the current
            developer's local state only.
  --help    Show this message (or, after a command, that command's
            usage).

Each command auto-detects project/workspace state and does the
right thing; re-running is safe and idempotent.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	cmd, args := os.Args[1], os.Args[2:]

	switch cmd {
	case "setup":
		os.Exit(runSetup(args))
	case "onboard":
		os.Exit(runOnboard(args))
	case "--help", "-h", "help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n%s", cmd, usage)
		os.Exit(2)
	}
}
