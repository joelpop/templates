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
// file. Keys are the placeholder names as they appear inside angle
// brackets in the templates (e.g., "JAVA_VERSION" for <JAVA_VERSION>).
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

// placeholderPattern matches <UPPERCASE_NAME> with at least one character
// and allowing digits and underscores after the first letter.
var placeholderPattern = regexp.MustCompile(`<([A-Z][A-Z0-9_]*)>`)

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

	// Remove orphans first (re-setup only).
	if removeOrphans {
		for k := range vars {
			if _, keep := req[k]; !keep {
				delete(vars, k)
			}
		}
	}

	// Fill in every required placeholder.
	for _, name := range required {
		if _, present := vars[name]; present {
			continue
		}
		// Prefer an auto-discovered value if one was provided.
		if v, ok := discovered[name]; ok && v != "" {
			vars[name] = v
			continue
		}
		// Fall back to an interactive prompt.
		def, ok := knownPlaceholders[name]
		if !ok {
			return fmt.Errorf("no source for placeholder %q — add it to knownPlaceholders", name)
		}
		value, err := def.Resolve()
		if err != nil {
			return fmt.Errorf("resolve %s: %w", name, err)
		}
		vars[name] = value
	}
	return nil
}
