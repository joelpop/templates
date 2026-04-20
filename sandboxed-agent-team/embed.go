package main

import "embed"

// templateFS holds every template file under ./templates/ bundled
// into the compiled binary. The `all:` prefix is required to include
// files under dot-directories (.sandbox/, .claude/) — go:embed excludes
// those by default.
//
//go:embed all:templates
var templateFS embed.FS
