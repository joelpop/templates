package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RunSandboxStart invokes .sandbox/start.sh from projectRoot. This
// builds the Docker image (if needed) and starts the sandbox. The
// script is interactive — stdin/stdout/stderr are passed through.
func RunSandboxStart(projectRoot string) error {
	script := filepath.Join(projectRoot, ".sandbox", "start.sh")
	if _, err := os.Stat(script); err != nil {
		return fmt.Errorf("sandbox start script not found at %s: %w", script, err)
	}
	cmd := exec.Command(script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = projectRoot
	return cmd.Run()
}

// RunSandboxStop invokes .sandbox/stop.sh. No-op if the script is
// absent (e.g., called before the kit is installed).
func RunSandboxStop(projectRoot string) error {
	script := filepath.Join(projectRoot, ".sandbox", "stop.sh")
	if _, err := os.Stat(script); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	cmd := exec.Command(script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = projectRoot
	return cmd.Run()
}
