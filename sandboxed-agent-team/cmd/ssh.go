package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// SSHRemoteInfo holds what we derive about a Git SSH remote URL.
type SSHRemoteInfo struct {
	// Host portion of the URL — may be a host alias from ~/.ssh/config.
	Host string
	// User portion (typically "git").
	User string
	// Absolute path to the identity file resolved from ~/.ssh/config,
	// or a default key fallback. Empty if not found.
	IdentityFile string
	// Real hostname — if Host is an alias, this is what ~/.ssh/config
	// resolves it to. Empty if Host is already a real hostname.
	RealHostName string
}

var sshGitURL = regexp.MustCompile(`^([^@]+)@([^:]+):`)

// ParseSSHURL extracts (user, host) from an SSH-style git URL. Accepts
// both `git@host:path` and `ssh://user@host/path`. Returns ok=false
// for anything else.
func ParseSSHURL(raw string) (user, host string, ok bool) {
	if strings.HasPrefix(raw, "ssh://") {
		u, err := url.Parse(raw)
		if err != nil {
			return "", "", false
		}
		return u.User.Username(), u.Hostname(), true
	}
	if m := sshGitURL.FindStringSubmatch(raw); m != nil {
		return m[1], m[2], true
	}
	return "", "", false
}

// LookupSSHHost parses ~/.ssh/config looking for a Host entry whose
// pattern list includes host. Returns the IdentityFile and
// HostName values from the matched block. Both can be empty if the
// block doesn't declare them.
func LookupSSHHost(host string) (identityFile, realHostName string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	f, err := os.Open(filepath.Join(home, ".ssh", "config"))
	if os.IsNotExist(err) {
		return "", "", nil
	}
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	inBlock := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lower := strings.ToLower(line)
		switch {
		case strings.HasPrefix(lower, "host "):
			inBlock = false
			for _, pat := range strings.Fields(line[len("host "):]) {
				if pat == host {
					inBlock = true
					break
				}
			}
		case inBlock && strings.HasPrefix(lower, "identityfile "):
			identityFile = expandPath(strings.TrimSpace(line[len("identityfile "):]))
		case inBlock && strings.HasPrefix(lower, "hostname "):
			realHostName = strings.TrimSpace(line[len("hostname "):])
		}
	}
	return identityFile, realHostName, scanner.Err()
}

// defaultSSHKey returns the first existing default key path, preferring
// ed25519 over rsa. Returns "" if neither exists.
func defaultSSHKey() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	for _, name := range []string{"id_ed25519", "id_rsa"} {
		path := filepath.Join(home, ".ssh", name)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func expandPath(p string) string {
	if strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, p[2:])
		}
	}
	return p
}

// DiscoverSSHRemote returns SSHRemoteInfo for the project's `origin`
// remote if it uses SSH. Returns nil for HTTPS remotes or if no remote
// is configured.
func DiscoverSSHRemote() *SSHRemoteInfo {
	raw := GitRemoteURL("origin")
	if raw == "" {
		return nil
	}
	user, host, ok := ParseSSHURL(raw)
	if !ok {
		return nil
	}
	info := &SSHRemoteInfo{Host: host, User: user}
	info.IdentityFile, info.RealHostName, _ = LookupSSHHost(host)
	if info.IdentityFile == "" {
		info.IdentityFile = defaultSSHKey()
	}
	return info
}

// ProvisionSSH populates .sandbox/ssh/ for the current developer so
// start.sh can inject the material into the sandbox on each run.
// Writes: the private key (+ .pub if present), a sandbox-scoped SSH
// config, a known_hosts entry from ssh-keyscan, and .sandbox/ssh.source
// pointing at the host key path (so start.sh can re-sync later if the
// key rotates). No-op if info is nil or has no identity file.
func ProvisionSSH(projectRoot string, info *SSHRemoteInfo) error {
	if info == nil || info.IdentityFile == "" {
		return nil
	}
	if _, err := os.Stat(info.IdentityFile); err != nil {
		return fmt.Errorf("ssh key %s is not readable: %w", info.IdentityFile, err)
	}

	sshDir := filepath.Join(projectRoot, ".sandbox", "ssh")
	if err := os.MkdirAll(sshDir, 0o700); err != nil {
		return err
	}

	keyName := filepath.Base(info.IdentityFile)
	if err := copyFile(info.IdentityFile, filepath.Join(sshDir, keyName), 0o600); err != nil {
		return fmt.Errorf("copy ssh key: %w", err)
	}
	if pub := info.IdentityFile + ".pub"; fileExists(pub) {
		_ = copyFile(pub, filepath.Join(sshDir, keyName+".pub"), 0o644)
	}

	realHost := info.RealHostName
	if realHost == "" {
		realHost = info.Host
	}
	sshConfig := fmt.Sprintf(`Host %s
    HostName %s
    User %s
    IdentityFile ~/.ssh/%s
    IdentitiesOnly yes
`, info.Host, realHost, info.User, keyName)
	if err := os.WriteFile(filepath.Join(sshDir, "config"), []byte(sshConfig), 0o644); err != nil {
		return err
	}

	out, err := exec.Command("ssh-keyscan", realHost).Output()
	if err != nil {
		return fmt.Errorf("ssh-keyscan %s: %w", realHost, err)
	}
	if err := os.WriteFile(filepath.Join(sshDir, "known_hosts"), out, 0o644); err != nil {
		return err
	}

	sourcePath := filepath.Join(projectRoot, ".sandbox", "ssh.source")
	return os.WriteFile(sourcePath, []byte(info.IdentityFile+"\n"), 0o644)
}

func copyFile(src, dst string, mode os.FileMode) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, mode)
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
