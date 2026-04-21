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
  agent-team <command> [flags]

Commands:
  install     Install the kit on a project. On a fresh project,
              bootstraps from scratch. On a project that already
              has the kit, reconciles variables, regenerates files,
              and re-commits. Auto-detects state.
  uninstall   Remove the kit from the project — deletes generated
              files, excises the CLAUDE.md import block and the
              kit's .gitignore block, commits the removal. Does
              NOT touch docs/.

Flags:
  --help      Show this message (or, after a command, that command's
              usage).

After installing, developers run ./team/join.sh in the project to
provision their local sandbox and start the team.
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
