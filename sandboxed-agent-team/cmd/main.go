// Command agent-team-install installs the sandboxed Claude Code
// agent team kit in a target project. Once installed, lifecycle
// commands (join, leave, start, stop, uninstall) live in the
// project's team/ directory and are invoked directly from there.
//
// The binary is self-contained — every template file is bundled
// via go:embed (see embed.go). At runtime it needs nothing from
// the kit source directory, so the binary can live anywhere.
// Do NOT rename it to "install" on PATH: that shadows GNU
// install(1) and causes confusing failures in any shell script
// that invokes install(1).
package main

import (
	"fmt"
	"os"
)

const usage = `agent-team-install — install the sandboxed Claude Code agent team kit

Usage:
  agent-team-install [--help]

Running "agent-team-install" installs (or updates) the kit on the
project in the current directory. Before anything is written,
you'll see a roadmap of what's about to happen and a review of
your choices; Ctrl-C aborts at any time.

Install Flow:
On a fresh project, agent-team-install prompts for any inputs it
can't derive automatically and writes the kit files to the
current directory. It then offers (Y/n) to have you join the
team by provisioning your workstation (set up your Docker
sandbox, SSH keys, platform API token) and starting the team.
No git operations — you stage, commit, and push on your own
schedule.

Update Flow:
On a project that already has the kit, agent-team-install
reconciles the variables file against the current templates —
prompting for new placeholders and dropping ones no longer used
— then regenerates the files the tool owns. Sample files the
kit produced (e.g., docs/INDEX.md) are left as-is. No git
operations — you stage, commit, and push on your own schedule.

If your workstation is already provisioned (./team/join.sh was
previously run), agent-team-install automatically re-runs
./team/join.sh to keep your local sandbox in sync with the
updated kit. If not, it prints a pointer to ./team/join.sh
so you can run it when you're ready.

Everything else lives in the project's ./team/ directory
(placed there by install):

  ./team/join.sh       Provision your workstation and launch the team.
  ./team/leave.sh      Discard your workstation's local sandbox state.
  ./team/create.sh     Build the sandbox image and launch a fresh sandbox.
  ./team/attach.sh     Reattach to an already-running sandbox.
  ./team/shell.sh      Drop into a bash shell inside the running sandbox.
  ./team/destroy.sh    Destroy the sandbox VM.
  ./team/uninstall.sh  Remove the kit from the project.

Running these scripts directly (rather than going through the
Go binary) keeps the lifecycle commands in lockstep with the
kit version committed to the project — no chance of an older
agent-team-install binary calling a newer script it doesn't match.

Other developers working on the project run

  ./team/join.sh

directly to provision their workstation and start the team.
`

func main() {
	// Only --help/-h/help are recognized. Any other argument is an
	// error; with no subcommands, a bare invocation runs install.
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--help", "-h", "help":
			fmt.Print(usage)
			return
		default:
			fmt.Fprintf(os.Stderr, "unknown argument: %s\n\n%s", arg, usage)
			os.Exit(2)
		}
	}
	os.Exit(runInstall())
}
