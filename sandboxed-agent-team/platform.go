package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// PlatformInfo describes the hosting platform of the project's Git
// remote. Used to guide the PR API token prompt and to write
// .sandbox/platform-api.env.
type PlatformInfo struct {
	Type     string // GITHUB, BITBUCKET, GITLAB, or empty if unknown
	APIURL   string // base URL for the platform's REST API
	Owner    string // owner / workspace / group
	RepoName string // repo name, no .git suffix
	RepoSlug string // Bitbucket's URL-path slug; elsewhere = RepoName
}

// DetectPlatform parses the project's `origin` remote URL and returns
// the platform info. Returns nil if no remote is set or the host isn't
// recognized.
func DetectPlatform() *PlatformInfo {
	raw := GitRemoteURL("origin")
	if raw == "" {
		return nil
	}

	host, path := splitRemoteURL(raw)
	if host == "" {
		return nil
	}
	path = strings.TrimSuffix(path, ".git")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return nil
	}
	info := &PlatformInfo{
		Owner:    parts[0],
		RepoName: parts[len(parts)-1],
		RepoSlug: parts[len(parts)-1],
	}
	switch {
	case strings.Contains(host, "github.com"):
		info.Type = "GITHUB"
		info.APIURL = "https://api.github.com"
	case strings.Contains(host, "bitbucket.org"):
		info.Type = "BITBUCKET"
		info.APIURL = "https://api.bitbucket.org"
	case strings.Contains(host, "gitlab"):
		info.Type = "GITLAB"
		info.APIURL = fmt.Sprintf("https://%s/api/v4", host)
	default:
		return nil
	}
	return info
}

// splitRemoteURL returns (host, path) for both SSH-style and
// URL-style git remotes.
func splitRemoteURL(raw string) (host, path string) {
	if strings.HasPrefix(raw, "http://") ||
		strings.HasPrefix(raw, "https://") ||
		strings.HasPrefix(raw, "ssh://") {
		u, err := url.Parse(raw)
		if err != nil {
			return "", ""
		}
		return u.Hostname(), strings.TrimPrefix(u.Path, "/")
	}
	if m := sshGitURL.FindStringSubmatch(raw); m != nil {
		return m[2], raw[len(m[0]):]
	}
	return "", ""
}

// WritePlatformAPIEnv prompts the developer for a platform API token
// and writes .sandbox/platform-api.env. Called during onboard only
// when MERGE_METHOD is PR.
func WritePlatformAPIEnv(projectRoot string, p *PlatformInfo) error {
	if p == nil {
		return fmt.Errorf(
			"cannot provision PR API access: no supported platform detected " +
				"from the `origin` remote (expected github.com, bitbucket.org, " +
				"or a gitlab host)")
	}

	fmt.Printf("\nPR merge method is configured — need a %s API token.\n", p.Type)
	fmt.Printf("Repo: %s/%s\n", p.Owner, p.RepoName)
	fmt.Println()
	switch p.Type {
	case "GITHUB":
		fmt.Println("Create a fine-grained PAT at:")
		fmt.Println("    https://github.com/settings/tokens?type=beta")
		fmt.Println("Scope it to this repository with:")
		fmt.Println("    Contents: Read")
		fmt.Println("    Pull requests: Read and write")
	case "BITBUCKET":
		fmt.Println("Create an app password at:")
		fmt.Println("    https://bitbucket.org/account/settings/app-passwords/")
		fmt.Println("Grant: Repositories:Read, Pull requests:Read + Write.")
	case "GITLAB":
		fmt.Println("Create a personal access token at:")
		fmt.Println("    https://gitlab.com/-/user_settings/personal_access_tokens")
		fmt.Println("Scope: api.")
	}
	fmt.Println()

	token, err := Prompt("Paste the token: ")
	if err != nil {
		return err
	}
	if token == "" {
		return fmt.Errorf("token is required for the PR merge method")
	}

	var apiUser string
	if p.Type == "BITBUCKET" {
		apiUser, err = Prompt("Bitbucket username (for basic auth): ")
		if err != nil {
			return err
		}
	}

	envPath := filepath.Join(projectRoot, ".sandbox", "platform-api.env")
	if err := os.MkdirAll(filepath.Dir(envPath), 0o700); err != nil {
		return err
	}
	content := fmt.Sprintf(`# Developer-local. Gitignored. Do NOT commit.
PLATFORM_TYPE=%s
PLATFORM_API_URL=%s
PLATFORM_API_USER=%s
PLATFORM_API_TOKEN=%s
PLATFORM_REPO_OWNER=%s
PLATFORM_REPO_NAME=%s
PLATFORM_REPO_WORKSPACE=%s
PLATFORM_REPO_SLUG=%s
`, p.Type, p.APIURL, apiUser, token, p.Owner, p.RepoName, p.Owner, p.RepoSlug)
	return os.WriteFile(envPath, []byte(content), 0o600)
}
