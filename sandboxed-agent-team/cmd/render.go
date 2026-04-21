package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// templateRoot is the directory inside the embedded FS that holds the
// template files.
const templateRoot = "templates"

// sampleFiles lists template paths (relative to the templates/ root)
// that are seeded once at initial setup and never overwritten
// afterward. These are the project's to own and evolve. See plan F2.
var sampleFiles = map[string]bool{
	"docs/INDEX.md": true,
}

// RenderTemplate reads one template (path is relative to the embed FS
// root, so it starts with "templates/") and substitutes every
// <PLACEHOLDER> with its value from vars.
func RenderTemplate(embedPath string, vars Variables) ([]byte, error) {
	raw, err := templateFS.ReadFile(embedPath)
	if err != nil {
		return nil, fmt.Errorf("read embedded template %s: %w", embedPath, err)
	}
	return Substitute(raw, vars)
}

// Substitute replaces every <PLACEHOLDER> (uppercase + digits +
// underscores) in content with the corresponding value from vars.
// Returns an error listing any placeholders for which vars has no
// value — silently leaving them in output would be a data bug.
func Substitute(content []byte, vars Variables) ([]byte, error) {
	var missing []string
	out := placeholderPattern.ReplaceAllFunc(content, func(match []byte) []byte {
		name := string(match[1 : len(match)-1])
		v, ok := vars[name]
		if !ok {
			missing = append(missing, name)
			return match
		}
		return []byte(v)
	})
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing variable(s): %s", strings.Join(missing, ", "))
	}
	return out, nil
}

// WriteAllTemplates renders every template in the embedded FS and
// writes the result to its canonical path under projectRoot. Sample
// files (docs/INDEX.md, ...) are only written if their target
// doesn't already exist.
//
// Returns the list of target paths that were written (relative to
// projectRoot), so callers can commit them to git.
func WriteAllTemplates(projectRoot string, vars Variables) ([]string, error) {
	var written []string
	err := fs.WalkDir(templateFS, templateRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel := strings.TrimPrefix(p, templateRoot+"/")
		targetPath := filepath.Join(projectRoot, rel)

		if sampleFiles[rel] {
			if _, err := os.Stat(targetPath); err == nil {
				// Sample file already present — leave it alone.
				return nil
			} else if !os.IsNotExist(err) {
				return err
			}
		}

		rendered, err := RenderTemplate(p, vars)
		if err != nil {
			return fmt.Errorf("%s: %w", rel, err)
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return err
		}

		mode := os.FileMode(0o644)
		if strings.HasSuffix(rel, ".sh") {
			mode = 0o755
		}

		if err := os.WriteFile(targetPath, rendered, mode); err != nil {
			return err
		}
		written = append(written, rel)
		return nil
	})
	return written, err
}
