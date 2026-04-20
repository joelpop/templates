package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// gitPreflight runs the repo/remote/uncommitted checks that precede
// every setup-path flow. Returns an error if the directory is not a
// git repo (hard fail); warns (but does not fail) on missing remote
// or uncommitted changes unless the user declines to proceed.
func gitPreflight(projectRoot string) error {
	if !IsGitRepo() {
		return fmt.Errorf(
			"this directory is not a git repository.\n" +
				"  Run `git init` (and optionally `git remote add origin <url>`),\n" +
				"  then re-run setup.")
	}

	if !HasRemote() {
		fmt.Println("Warning: no git remote configured. Setup will proceed,")
		fmt.Println("but you'll need a remote before you can use the PR merge")
		fmt.Println("method or push branches. Add one with:")
		fmt.Println("    git remote add origin <url>")
		fmt.Println()
	}

	dirty, err := HasUncommittedChanges()
	if err != nil {
		return err
	}
	if dirty {
		fmt.Println("You have uncommitted changes. Setup will write new files and")
		fmt.Println("may switch branches. Commit or stash your changes first.")
		ok, err := PromptYesNo("Continue anyway?", false)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("aborted by user")
		}
	}
	return nil
}

// identifyDevBranch fetches remotes, looks for a conventional dev
// branch name, asks the user if ambiguous, and validates freshness
// against its remote counterpart. Returns the confirmed branch name.
func identifyDevBranch() (string, error) {
	if HasRemote() {
		if err := FetchAllPrune(); err != nil {
			fmt.Printf("Warning: git fetch --all --prune failed: %s\n", err)
			ok, perr := PromptYesNo("Continue without fresh remote refs?", false)
			if perr != nil {
				return "", perr
			}
			if !ok {
				return "", fmt.Errorf("aborted")
			}
		}
	}

	conventional := []string{"develop", "dev", "development", "feature/develop"}
	locals, _ := LocalBranches()
	remotes, _ := RemoteBranches()

	// Collect all branch names (local + remote short names).
	all := map[string]bool{}
	for _, b := range locals {
		all[b] = true
	}
	for _, b := range remotes {
		// Strip leading "origin/" or similar.
		if i := strings.Index(b, "/"); i >= 0 {
			all[b[i+1:]] = true
		}
	}

	var matches []string
	for _, name := range conventional {
		if all[name] {
			matches = append(matches, name)
		}
	}

	var chosen string
	switch len(matches) {
	case 1:
		ok, err := PromptYesNo(fmt.Sprintf("I found %q. Is this your development branch?", matches[0]), true)
		if err != nil {
			return "", err
		}
		if ok {
			chosen = matches[0]
		}
	case 0:
		// Nothing conventional found — ask.
	default:
		// Multiple conventional matches — ask which.
	}
	if chosen == "" {
		raw, err := PromptWithDefault(
			"Which branch is your development branch? (e.g., develop, feature/develop, or main for trunk-based)",
			"")
		if err != nil {
			return "", err
		}
		if raw == "" {
			return "", fmt.Errorf("development branch is required")
		}
		chosen = raw
	}

	// Refuse to default to main/master silently.
	if (chosen == "main" || chosen == "master") && !confirmTrunk(chosen) {
		return "", fmt.Errorf("setup will not treat %s as the development branch without explicit trunk-based confirmation", chosen)
	}

	// If branch exists only on remote, offer to check it out locally.
	localHas, remoteHas, err := BranchExists(chosen)
	if err != nil {
		return "", err
	}
	if !localHas && remoteHas {
		ok, perr := PromptYesNo(fmt.Sprintf("%q exists on the remote but not locally. Check it out?", chosen), true)
		if perr != nil {
			return "", perr
		}
		if ok {
			if err := Checkout(chosen); err != nil {
				return "", err
			}
			localHas = true
		}
	}
	if !localHas && !remoteHas {
		return "", fmt.Errorf("branch %q does not exist locally or on the remote", chosen)
	}

	// Freshness check (if the local branch has a remote counterpart).
	if remoteHas {
		ahead, behind, err := AheadBehind(chosen)
		if err == nil {
			switch {
			case ahead == 0 && behind == 0:
				// current
			case behind > 0 && ahead == 0:
				fmt.Printf("Local %s is %d commit(s) behind origin/%s.\n", chosen, behind, chosen)
				ok, perr := PromptYesNo("Fast-forward now?", true)
				if perr != nil {
					return "", perr
				}
				if ok {
					if err := Checkout(chosen); err != nil {
						return "", err
					}
					if err := FastForward("origin/" + chosen); err != nil {
						return "", err
					}
				}
			case ahead > 0 && behind == 0:
				fmt.Printf("Local %s is %d commit(s) ahead of origin/%s — you'll want to push at some point.\n", chosen, ahead, chosen)
			case ahead > 0 && behind > 0:
				return "", fmt.Errorf("local and remote %s have diverged (%d ahead, %d behind); resolve before re-running setup", chosen, ahead, behind)
			}
		}
	}

	return chosen, nil
}

func confirmTrunk(name string) bool {
	fmt.Printf("You named %q as the development branch. That's normally a trunk-based choice.\n", name)
	ok, _ := PromptYesNo("Confirm trunk-based development?", false)
	return ok
}

// placeOnBranch ensures the setup flow runs on an acceptable branch.
// If the current branch is not the dev branch, prompts the user to
// switch, create a new branch off dev, or stay. Returns the final
// branch the flow will run on.
func placeOnBranch(devBranch, defaultNewName string) (string, error) {
	current, err := CurrentBranch()
	if err != nil {
		return "", err
	}
	if current == devBranch {
		return current, nil
	}

	fmt.Printf("\nSetup writes versioned config that must live on the development branch.\n")
	fmt.Printf("Your development branch is %q; you're currently on %q.\n", devBranch, current)

	options := []string{
		fmt.Sprintf("Switch to %s and run setup there", devBranch),
		fmt.Sprintf("Create a new branch off %s (default name: %s)", devBranch, defaultNewName),
		fmt.Sprintf("Stay on %s (only if you're sure this is where the kit should land)", current),
	}
	choice, err := PromptChoice("How should I proceed?", options, options[1])
	if err != nil {
		return "", err
	}
	switch choice {
	case options[0]:
		if err := Checkout(devBranch); err != nil {
			return "", err
		}
		return devBranch, nil
	case options[1]:
		name, err := PromptWithDefault("New branch name", defaultNewName)
		if err != nil {
			return "", err
		}
		if err := CheckoutNewBranchOff(name, devBranch); err != nil {
			return "", err
		}
		return name, nil
	case options[2]:
		return current, nil
	}
	return "", fmt.Errorf("unexpected choice: %s", choice)
}

// recordOnboarding writes a timestamp marker indicating that the
// current developer has onboarded this machine to this project.
func recordOnboarding(projectRoot string) error {
	path := filepath.Join(projectRoot, ".claude", ".last-onboarded")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(time.Now().UTC().Format(time.RFC3339)+"\n"), 0o644)
}
