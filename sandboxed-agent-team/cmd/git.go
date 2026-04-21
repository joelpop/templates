package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
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
// working tree.
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

// GitRemoteURL returns the fetch URL of the named remote, or "" if the
// remote does not exist.
func GitRemoteURL(remote string) string {
	url, err := git("remote", "get-url", remote)
	if err != nil {
		return ""
	}
	return url
}

// HasRemote returns true if at least one remote is configured.
func HasRemote() bool {
	out, err := git("remote")
	return err == nil && strings.TrimSpace(out) != ""
}

// IsSSHRemote returns true if the given URL is an SSH-style remote
// (starts with git@ or ssh://).
func IsSSHRemote(url string) bool {
	return strings.HasPrefix(url, "git@") || strings.HasPrefix(url, "ssh://")
}

// HasUncommittedChanges returns true if the working tree or index has
// unstaged or uncommitted changes.
func HasUncommittedChanges() (bool, error) {
	out, err := git("status", "--porcelain")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

// CurrentBranch returns the abbreviated name of the current branch.
func CurrentBranch() (string, error) {
	return git("rev-parse", "--abbrev-ref", "HEAD")
}

// FetchAllPrune runs `git fetch --all --prune`.
func FetchAllPrune() error {
	_, err := git("fetch", "--all", "--prune")
	return err
}

// LocalBranches returns the list of local branch names.
func LocalBranches() ([]string, error) {
	out, err := git("for-each-ref", "--format=%(refname:short)", "refs/heads")
	if err != nil {
		return nil, err
	}
	return splitLines(out), nil
}

// RemoteBranches returns the list of remote-tracking branch names (e.g.,
// "origin/develop"). Requires a prior fetch to be current.
func RemoteBranches() ([]string, error) {
	out, err := git("for-each-ref", "--format=%(refname:short)", "refs/remotes")
	if err != nil {
		return nil, err
	}
	return splitLines(out), nil
}

// BranchExists reports whether `name` exists locally and/or on a remote.
func BranchExists(name string) (local, remote bool, err error) {
	locals, err := LocalBranches()
	if err != nil {
		return false, false, err
	}
	for _, b := range locals {
		if b == name {
			local = true
			break
		}
	}
	remotes, err := RemoteBranches()
	if err != nil {
		return local, false, err
	}
	for _, b := range remotes {
		// b looks like "origin/develop"; match the "/<name>" suffix.
		if strings.HasSuffix(b, "/"+name) {
			remote = true
			break
		}
	}
	return local, remote, nil
}

// AheadBehind returns how far `local` is ahead of / behind
// `origin/<local>`. Both zero means the local branch matches the
// remote counterpart. If no remote counterpart exists, the returned
// error wraps "unknown revision" — callers should handle that case.
func AheadBehind(local string) (ahead, behind int, err error) {
	spec := fmt.Sprintf("%s...origin/%s", local, local)
	out, err := git("rev-list", "--left-right", "--count", spec)
	if err != nil {
		return 0, 0, err
	}
	parts := strings.Fields(out)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected output from rev-list: %q", out)
	}
	ahead, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	behind, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	return ahead, behind, nil
}

// Checkout switches to an existing branch.
func Checkout(branch string) error {
	_, err := git("checkout", branch)
	return err
}

// CheckoutNewBranchOff creates `name` off `base` and switches to it.
// It always runs `git checkout <base>` first, then `git checkout -b
// <name>`. This is deliberate — never trust the current HEAD to be the
// intended base; past sessions have landed setup branches on `main`
// instead of on `develop` precisely because of that implicit-base
// assumption.
func CheckoutNewBranchOff(name, base string) error {
	if _, err := git("checkout", base); err != nil {
		return fmt.Errorf("checkout base %s: %w", base, err)
	}
	if _, err := git("checkout", "-b", name); err != nil {
		return fmt.Errorf("create branch %s off %s: %w", name, base, err)
	}
	return nil
}

// FastForward runs `git merge --ff-only <upstream>` on the current
// branch.
func FastForward(upstream string) error {
	_, err := git("merge", "--ff-only", upstream)
	return err
}

// GitAddForce stages the given paths with --force, bypassing the
// project's .gitignore. The kit owns these paths by design.
func GitAddForce(paths ...string) error {
	if len(paths) == 0 {
		return nil
	}
	args := append([]string{"add", "--force", "--"}, paths...)
	_, err := git(args...)
	return err
}

// GitCommit creates a commit with the given message. Returns nil if
// there's nothing to commit (e.g., a re-run produced no changes).
func GitCommit(message string) error {
	// If the index has no staged changes, git commit will fail with
	// "nothing to commit" — treat that as success (it's a no-op rerun).
	if clean, err := indexIsClean(); err != nil {
		return err
	} else if clean {
		return nil
	}
	_, err := git("commit", "-m", message)
	return err
}

func indexIsClean() (bool, error) {
	out, err := git("diff", "--cached", "--name-only")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) == "", nil
}

func splitLines(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}
