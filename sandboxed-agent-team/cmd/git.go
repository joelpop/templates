package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// git runs a git subcommand in the current working directory and
// returns its combined stdout/stderr. Non-zero exit is returned as an
// error whose Error() includes the command and the output.
func git(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return strings.TrimRight(buf.String(), "\n"),
			fmt.Errorf("git %s: %w (output: %s)",
				strings.Join(args, " "), err,
				strings.TrimSpace(buf.String()))
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}

// IsGitRepo returns true iff the current directory is inside a git
// working tree. The installer does not gate any behavior on this —
// it's used only to tailor post-install / post-uninstall messaging
// (suggest `git status` inside a repo; stay silent about git
// outside one).
func IsGitRepo() bool {
	_, err := git("rev-parse", "--git-dir")
	return err == nil
}

// GitIdentity returns (user.name, user.email). Either may be empty if
// not configured.
func GitIdentity() (string, string) {
	name, _ := git("config", "user.name")
	email, _ := git("config", "user.email")
	return name, email
}

// CurrentBranch returns the abbreviated name of the current branch.
// Used by the review summary to tell the user which branch the kit
// files will land in.
func CurrentBranch() (string, error) {
	return git("rev-parse", "--abbrev-ref", "HEAD")
}

// LocalBranches returns the list of local branch names.
func LocalBranches() ([]string, error) {
	out, err := git("for-each-ref", "--format=%(refname:short)", "refs/heads")
	if err != nil {
		return nil, err
	}
	return splitLines(out), nil
}

func splitLines(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}
