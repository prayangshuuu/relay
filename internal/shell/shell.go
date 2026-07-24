// Package shell provides abstractions for system command execution.
// Future responsibilities include wrapping cross-platform shell behaviors
// (Windows, macOS, Linux) and providing safe primitives for spawning
// child processes without retaining background execution.
package shell

// Shell abstracts the operating system's process execution mechanism.
type Shell interface {
	// Exec executes the given binary with the specified arguments and environment.
	// It aims to replace the current process (e.g., via syscall.Exec on Unix)
	// or spawn a detached process and exit immediately on platforms where
	// replacement is not possible.
	Exec(executable string, args []string, env []string) error
}
