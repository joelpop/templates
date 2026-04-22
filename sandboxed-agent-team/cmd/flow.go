package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// detectConventionalDevBranch returns a sensible default for
// DEV_BRANCH_NAME, or "" if it can't confidently pick one.
// Read-only; no network; no working-tree changes.
//
// Tier 1: exact-match against names the kit conventionally uses
// as a dev / integration branch (the GitFlow-style case, which
// is typical for this kit's audience).
//
// Tier 2: the remote's cached HEAD — what `git clone` recorded
// as the default branch. Correct for trunk-based projects; uses
// only local git refs, no network.
//
// Tier 3: nothing confident; return "". The placeholder
// interview then falls through to a prompt with no default.
func detectConventionalDevBranch() string {
	conventional := []string{"develop", "dev", "development", "feature/develop"}
	locals, _ := LocalBranches()
	seen := map[string]bool{}
	for _, b := range locals {
		seen[b] = true
	}
	var matches []string
	for _, name := range conventional {
		if seen[name] {
			matches = append(matches, name)
		}
	}
	if len(matches) == 1 {
		return matches[0]
	}

	if ref, err := git("symbolic-ref", "refs/remotes/origin/HEAD"); err == nil {
		ref = strings.TrimSpace(ref)
		if i := strings.LastIndex(ref, "/"); i >= 0 {
			return ref[i+1:]
		}
	}

	return ""
}

// buildDiscoveredMap collects every auto-discoverable value for the
// current run: pom.xml fields, git identity, dev branch (via
// detectConventionalDevBranch), build tool, stack summary, and
// timestamps. Returns a Variables map that ReconcileVariables uses
// to prefer fresh values over persisted ones for SourceAuto
// placeholders.
func buildDiscoveredMap(projectRoot string) (Variables, error) {
	pomPath := filepath.Join(projectRoot, "pom.xml")
	info, err := DiscoverProject(pomPath)
	if err != nil {
		return nil, err
	}
	d := info.ToVariables()

	if _, err := os.Stat(pomPath); err == nil {
		d["BUILD_TOOL"] = "Maven"
	}

	if name, email := GitIdentity(); name != "" || email != "" {
		if name != "" {
			d["GIT_USER_NAME"] = name
		}
		if email != "" {
			d["GIT_USER_EMAIL"] = email
		}
	}

	if b := detectConventionalDevBranch(); b != "" {
		d["DEV_BRANCH_NAME"] = b
	}

	d["UTC_TIMESTAMP"] = time.Now().UTC().Format(time.RFC3339)

	return d, nil
}
