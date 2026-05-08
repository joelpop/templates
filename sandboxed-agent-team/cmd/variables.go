package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Variables is the in-memory representation of the persisted variables
// file. Keys are the placeholder names as they appear inside double
// braces in the templates (e.g., "JAVA_VERSION" for {{JAVA_VERSION}}).
type Variables map[string]string

// DefaultVariablesPath is the canonical location of the variables file
// inside a target project.
const DefaultVariablesPath = ".claude/team-variables.yaml"

// LoadVariables reads the variables file at path. Returns an empty
// Variables (not an error) if the file does not exist — the virgin-project
// case, where setup hasn't persisted anything yet.
func LoadVariables(path string) (Variables, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Variables{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	v := Variables{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return v, nil
}

// SaveVariables writes the variables file to path, creating parent
// directories as needed. Keys are serialized in sorted order so the
// output is deterministic across runs.
func SaveVariables(path string, v Variables) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}

	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString("# Persisted variables for the sandboxed-agent-team kit.\n")
	sb.WriteString("# The setup tool reads this on re-run, adds newly-introduced\n")
	sb.WriteString("# placeholders, and removes orphan keys. Hand-editing is expected.\n\n")
	for _, k := range keys {
		// Quote every value for consistency; yaml accepts unquoted scalars
		// but our values may contain characters that need quoting (spaces,
		// special chars).
		sb.WriteString(fmt.Sprintf("%s: %q\n", k, v[k]))
	}

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

// placeholderPattern matches {{UPPERCASE_NAME}} with at least one
// character and allowing digits and underscores after the first letter.
// Double braces (Mustache / Handlebars / Go-template style) are used
// rather than angle brackets so the pattern doesn't collide with Java
// generic syntax (`<T>`, `<KEY>`, etc.) that appears in template content
// under `templates/docs/patterns/`.
var placeholderPattern = regexp.MustCompile(`\{\{([A-Z][A-Z0-9_]*)\}\}`)

// ScanPlaceholders walks the given filesystem and returns every
// placeholder name referenced by any file, sorted and deduplicated.
func ScanPlaceholders(tfs fs.FS) ([]string, error) {
	seen := map[string]struct{}{}
	err := fs.WalkDir(tfs, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(tfs, p)
		if err != nil {
			return fmt.Errorf("read %s: %w", p, err)
		}
		for _, m := range placeholderPattern.FindAllSubmatch(data, -1) {
			seen[string(m[1])] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(seen))
	for name := range seen {
		out = append(out, name)
	}
	sort.Strings(out)
	return out, nil
}

// ReconcileVariables amends vars in-place so that every required
// placeholder has a value, prompting (via src) for any that can't be
// auto-discovered. Orphan keys (not in required) are removed if
// removeOrphans is true.
func ReconcileVariables(vars Variables, required []string, discovered Variables, removeOrphans bool) error {
	req := map[string]struct{}{}
	for _, r := range required {
		req[r] = struct{}{}
	}

	// Remove orphans first (re-setup only). Announce them to the user
	// before the review step so a hand-edited value isn't silently
	// dropped; they can abort with Ctrl-C if needed.
	if removeOrphans {
		var orphans []string
		for k := range vars {
			if _, keep := req[k]; !keep {
				orphans = append(orphans, k)
			}
		}
		if len(orphans) > 0 {
			sort.Strings(orphans)
			fmt.Println()
			fmt.Println("The current templates no longer reference these variables:")
			for _, k := range orphans {
				fmt.Printf("  %s = %q\n", k, vars[k])
			}
			fmt.Println("They will be removed from .claude/team-variables.yaml.")
			fmt.Println("Abort with Ctrl-C now if any value needs preserving.")
			fmt.Println()
			for _, k := range orphans {
				delete(vars, k)
			}
		}
	}

	// Iterate in the user-facing canonical order so the prompts
	// appear in the order the roadmap advertises. `required` comes
	// in alphabetical (from ScanPlaceholders); we re-sort here to
	// match the review-summary layout.
	ordered := append([]string(nil), required...)
	SortByUserOrder(ordered)

	// Fill in every required placeholder. For SourceAuto placeholders,
	// an up-to-date `discovered` value wins over any persisted value
	// (so pom.xml bumps, timestamps, etc. refresh on every run). For
	// SourcePrompt placeholders, a persisted value wins (user edits
	// to the variables file survive re-runs).
	for _, name := range ordered {
		def, ok := knownPlaceholders[name]
		if !ok {
			return fmt.Errorf("no source for placeholder %q — add it to knownPlaceholders", name)
		}

		if def.Source == SourceAuto {
			if v, has := discovered[name]; has && v != "" {
				vars[name] = v
				continue
			}
		}

		if _, present := vars[name]; present {
			continue
		}

		// Not yet set and no auto value — prompt the user.
		value, err := def.Resolve()
		if err != nil {
			return fmt.Errorf("resolve %s: %w", name, err)
		}
		vars[name] = value
	}
	return nil
}
