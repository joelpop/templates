package main

import (
	"fmt"
	"sort"
	"strings"
)

// PlaceholderSource describes how the tool obtains a value for a
// placeholder when the variables file doesn't already have one.
type PlaceholderSource int

const (
	// SourcePrompt: ask the human interactively.
	SourcePrompt PlaceholderSource = iota
	// SourceAuto: value is expected to be supplied in the `discovered`
	// map passed to ReconcileVariables (e.g., from pom.xml or git config).
	// If the discovered map lacks the value, we fall back to Prompt.
	SourceAuto
)

// PlaceholderDef describes one known placeholder the kit uses.
type PlaceholderDef struct {
	Name               string
	Source             PlaceholderSource
	Prompt             string   // used when Source is Prompt (or Auto fallback)
	Options            []string // non-empty → choice prompt
	OptionDescriptions []string // optional; parallel to Options, renders " — <desc>" after each
	Default            string   // default on empty input
}

// Resolve returns a value for the placeholder by prompting the user.
// Only called when the value isn't already in vars or discovered.
func (p PlaceholderDef) Resolve() (string, error) {
	if p.Prompt == "" {
		return "", fmt.Errorf("placeholder %s has no Prompt text", p.Name)
	}
	if len(p.Options) > 0 {
		if p.OptionDescriptions != nil {
			return PromptChoiceWithDescriptions(p.Prompt, p.Options, p.OptionDescriptions, p.Default)
		}
		return PromptChoice(p.Prompt, p.Options, p.Default)
	}
	return PromptWithDefault(p.Prompt, p.Default)
}

// knownPlaceholders is the authoritative registry: every placeholder the
// templates may contain must be listed here with its source info. Unknown
// placeholders cause reconciliation to fail — better to fail loudly than
// to ship a template referencing a placeholder the tool doesn't know how
// to resolve.
var knownPlaceholders = map[string]PlaceholderDef{
	"PROJECT_NAME": {
		Name:   "PROJECT_NAME",
		Source: SourceAuto,
		Prompt: "Project name (used in CLAUDE_TEAM.md)",
	},
	"JAVA_VERSION": {
		Name:   "JAVA_VERSION",
		Source: SourcePrompt,
		Prompt: "Java version (e.g., 21)",
	},
	"VAADIN_VERSION": {
		Name:   "VAADIN_VERSION",
		Source: SourcePrompt,
		Prompt: "Vaadin version",
	},
	"SPRING_BOOT_VERSION": {
		Name:    "SPRING_BOOT_VERSION",
		Source:  SourcePrompt,
		Prompt:  "Spring Boot version (leave blank if not used)",
		Default: "",
	},
	"JUNIT_VERSION": {
		Name:    "JUNIT_VERSION",
		Source:  SourcePrompt,
		Prompt:  "JUnit version",
		Default: "5",
	},
	"DATABASE": {
		Name:    "DATABASE",
		Source:  SourcePrompt,
		Prompt:  "Database (e.g., PostgreSQL, H2, leave blank if none)",
		Default: "",
	},
	"GIT_USER_NAME": {
		Name:   "GIT_USER_NAME",
		Source: SourceAuto,
		Prompt: "Git user name",
	},
	"GIT_USER_EMAIL": {
		Name:   "GIT_USER_EMAIL",
		Source: SourceAuto,
		Prompt: "Git user email",
	},
	"DEV_BRANCH_NAME": {
		Name:   "DEV_BRANCH_NAME",
		Source: SourcePrompt,
		Prompt: "Team's development branch name",
	},
	"MERGE_METHOD": {
		Name:    "MERGE_METHOD",
		Source:  SourcePrompt,
		Prompt:  "How should completed work reach the development branch on origin?",
		Options: []string{"PR", "Integrator merge", "Human merge"},
		OptionDescriptions: []string{
			`- Integrator pushes the task branch to origin.
- Integrator creates a PR via the platform's REST API.
- Reviewers approve on the platform.
- Integrator merges the PR via API.
- Platform deletes the remote task branch; Integrator deletes the local one.
Requires a platform API token (GitHub / Bitbucket / GitLab).`,
			`- Integrator squash-merges the task branch into the local dev branch.
- Integrator pushes the dev branch to origin.
- Integrator deletes the local task branch.
The task branch is never pushed; no PR is created.`,
			`- Integrator pauses after the task's local tests pass and tells you the task is ready.
- You perform the final merge however you prefer (local squash + push, a platform PR, etc.).
- You confirm to the team once the merge has landed.
- Integrator deletes the local task branch.`,
		},
		Default: "Integrator merge",
	},
	"COST_IN_COMMIT": {
		Name:    "COST_IN_COMMIT",
		Source:  SourcePrompt,
		Prompt:  "Append a per-model token/cost report to squash-merge commit messages?",
		Options: []string{"yes", "no"},
		Default: "no",
	},
	"COST_IN_LOG": {
		Name:    "COST_IN_LOG",
		Source:  SourcePrompt,
		Prompt:  "Append per-task token/cost reports to a project log file (.claude/.cost-log.md)?",
		Options: []string{"yes", "no"},
		Default: "no",
	},
	"CI_PLATFORM": {
		Name:    "CI_PLATFORM",
		Source:  SourcePrompt,
		Prompt:  "CI platform in use",
		Options: []string{"GitHub Actions", "GitLab CI", "Bitbucket Pipelines", "Jenkins", "none"},
		Default: "GitHub Actions",
	},
	"BUILD_TOOL": {
		Name:    "BUILD_TOOL",
		Source:  SourcePrompt,
		Prompt:  "Build tool",
		Default: "Maven",
	},
	"STACK_SUMMARY": {
		Name:   "STACK_SUMMARY",
		Source: SourceAuto,
		Prompt: "One-line stack summary",
	},
	"UTC_TIMESTAMP": {
		Name:   "UTC_TIMESTAMP",
		Source: SourceAuto,
		Prompt: "Current UTC timestamp (ISO 8601)",
	},
}

// placeholderOrder defines the canonical user-facing order for
// placeholders — used both for the interview (prompt order) and for
// the review summary. Any placeholder not listed here sorts after
// the listed ones, alphabetically. Keep this in sync with the
// printInstallIntro roadmap and the confirmProceedWithInstall
// review layout.
var placeholderOrder = []string{
	"DEV_BRANCH_NAME",
	"JAVA_VERSION",
	"VAADIN_VERSION",
	"SPRING_BOOT_VERSION",
	"JUNIT_VERSION",
	"DATABASE",
	"BUILD_TOOL",
	"MERGE_METHOD",
	"CI_PLATFORM",
	"COST_IN_COMMIT",
	"COST_IN_LOG",
	// Auto-discovered fall-backs — rarely prompt:
	"PROJECT_NAME",
	"STACK_SUMMARY",
	"GIT_USER_NAME",
	"GIT_USER_EMAIL",
	"UTC_TIMESTAMP",
}

// SortByUserOrder sorts the given placeholder names in place by
// placeholderOrder. Names not in placeholderOrder go to the end,
// sorted alphabetically.
func SortByUserOrder(names []string) {
	rank := map[string]int{}
	for i, n := range placeholderOrder {
		rank[n] = i
	}
	sort.SliceStable(names, func(i, j int) bool {
		ri, iIn := rank[names[i]]
		rj, jIn := rank[names[j]]
		switch {
		case iIn && jIn:
			return ri < rj
		case iIn:
			return true
		case jIn:
			return false
		default:
			return names[i] < names[j]
		}
	})
}

// CheckUnknownPlaceholders returns an error if any placeholder name in
// `found` is not declared in knownPlaceholders. Used as a build-time
// sanity check to keep the registry in sync with the templates.
func CheckUnknownPlaceholders(found []string) error {
	var unknown []string
	for _, name := range found {
		if _, ok := knownPlaceholders[name]; !ok {
			unknown = append(unknown, name)
		}
	}
	if len(unknown) == 0 {
		return nil
	}
	return fmt.Errorf("templates reference unknown placeholders: %s "+
		"(add each to knownPlaceholders in placeholders.go)",
		strings.Join(unknown, ", "))
}
