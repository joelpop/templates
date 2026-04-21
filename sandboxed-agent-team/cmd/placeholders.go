package main

import (
	"fmt"
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
	Name    string
	Source  PlaceholderSource
	Prompt  string   // used when Source is Prompt (or Auto fallback)
	Options []string // non-empty → choice prompt
	Default string   // default on empty input
}

// Resolve returns a value for the placeholder by prompting the user.
// Only called when the value isn't already in vars or discovered.
func (p PlaceholderDef) Resolve() (string, error) {
	if p.Prompt == "" {
		return "", fmt.Errorf("placeholder %s has no Prompt text", p.Name)
	}
	if len(p.Options) > 0 {
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
		Source: SourceAuto,
		Prompt: "Java version (e.g., 21)",
	},
	"VAADIN_VERSION": {
		Name:   "VAADIN_VERSION",
		Source: SourceAuto,
		Prompt: "Vaadin version",
	},
	"SPRING_BOOT_VERSION": {
		Name:    "SPRING_BOOT_VERSION",
		Source:  SourceAuto,
		Prompt:  "Spring Boot version (leave blank if not used)",
		Default: "",
	},
	"JUNIT_VERSION": {
		Name:    "JUNIT_VERSION",
		Source:  SourceAuto,
		Prompt:  "JUnit version",
		Default: "5",
	},
	"DATABASE": {
		Name:    "DATABASE",
		Source:  SourceAuto,
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
		Source: SourceAuto,
		Prompt: "Development branch name",
	},
	"MERGE_METHOD": {
		Name:    "MERGE_METHOD",
		Source:  SourcePrompt,
		Prompt:  "How should completed work reach the development branch?",
		Options: []string{"PR", "Integrator merge", "Human merge"},
		Default: "Integrator merge",
	},
	"COST_IN_COMMIT": {
		Name:    "COST_IN_COMMIT",
		Source:  SourcePrompt,
		Prompt:  "Append a per-model token/cost report to squash-merge commit messages?",
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
		Source:  SourceAuto,
		Prompt:  "Build tool",
		Default: "Maven",
	},
	"STACK_SUMMARY": {
		Name:   "STACK_SUMMARY",
		Source: SourceAuto,
		Prompt: "One-line stack summary",
	},
	"DATE": {
		Name:   "DATE",
		Source: SourceAuto,
		Prompt: "Date (YYYY-MM-DD)",
	},
	"UTC_TIMESTAMP": {
		Name:   "UTC_TIMESTAMP",
		Source: SourceAuto,
		Prompt: "Current UTC timestamp (ISO 8601)",
	},
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
