// Package runner provides process execution and signal handling for the GJG launcher.
package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Run executes the given argv with env and working directory. Returns the exit code.
func Run(argv, env []string, workDir string) (int, error) {
	if len(argv) == 0 {
		return 1, errors.New("empty argv")
	}
	// #nosec G204 - This is intentional subprocess execution with user-provided args for launcher
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Env = env
	cmd.Dir = workDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Handle Ctrl+C: forward as kill to child
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	if err := cmd.Start(); err != nil {
		return 1, fmt.Errorf("failed to start process: %w", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	for {
		select {
		case <-sigCh:
			if cmd.Process != nil {
				_ = cmd.Process.Signal(syscall.SIGKILL)
			}
		case err := <-done:
			if err == nil {
				return 0, nil
			}
			var ee *exec.ExitError
			if errors.As(err, &ee) {
				// On Windows, ExitCode available
				return ee.ExitCode(), nil
			}
			return 1, err
		}
	}
}
