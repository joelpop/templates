// Command agent-team installs and manages a sandboxed Claude Code
// agent team kit in a target project.
//
// See the plan at
// /Users/joel/.claude/plans/i-think-this-setup-re-setup-pure-blossom.md
// for the full requirements.
package main

import (
	"fmt"
	"os"
)

const usage = `agent-team — install and manage a sandboxed Claude Code agent team kit

Usage:
  agent-team <command> [--help]

Commands:
  install     Installs or updates the kit on a project.
  uninstall   Removes the kit from a project.

For details on any command, run ` + "`agent-team <command> --help`" + `.

After installing (or updating), subsequent developers run
./team/join.sh in the project to provision (or reprovision) their
local sandbox and start the team.
`

func main() {
	// No args → show help on stdout and exit 0. Same as `--help`.
	// Unknown command → error on stderr and exit 2.
	if len(os.Args) < 2 {
		fmt.Print(usage)
		return
	}

	cmd, args := os.Args[1], os.Args[2:]

	switch cmd {
	case "install":
		os.Exit(runInstall(args))
	case "uninstall":
		os.Exit(runUninstall(args))
	case "--help", "-h", "help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n%s", cmd, usage)
		os.Exit(2)
	}
}
