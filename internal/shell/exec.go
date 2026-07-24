package shell

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// OSExecShell implements the Shell interface using os/exec to ensure cross-platform
// compatibility while propagating standard streams, signals, and exit codes.
type OSExecShell struct{}

// NewOSExecShell creates a new shell execution wrapper.
func NewOSExecShell() *OSExecShell {
	return &OSExecShell{}
}

// Exec spawns the executable with the provided arguments and environment.
// It maps the current process's standard I/O to the child, forwards interrupt
// signals, waits for the child to exit, and then mirrors the exit code.
func (s *OSExecShell) Exec(executable string, args []string, env []string) error {
	cmd := exec.Command(executable, args...)
	cmd.Env = env

	// Directly wire standard descriptors to preserve TTY and colors
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Setup signal forwarding
	sigChan := make(chan os.Signal, 1)
	// Catch all signals typically sent to CLI apps
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	// Start the child process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start executable '%s': %w", executable, err)
	}

	// Forward signals to the child in a background goroutine
	go func() {
		for sig := range sigChan {
			if cmd.Process != nil {
				_ = cmd.Process.Signal(sig)
			}
		}
	}()

	// Wait for the child process to complete
	err := cmd.Wait()

	// Extract the exit code
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			// A non-exit error means it failed to wait or similar
			return fmt.Errorf("execution error: %w", err)
		}
	}

	// Terminate the relay process mirroring the child's exit code
	os.Exit(exitCode)
	return nil // Should never be reached
}
