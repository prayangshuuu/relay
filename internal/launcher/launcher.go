// Package launcher handles the execution of AI tools with configured environments.
// Future responsibilities include resolving profiles, generating the necessary
// environment payloads, wrapping the executable, spawning the process, and ensuring
// Relay exits transparently after launch.
package launcher

import (
	"github.com/prayangshuuu/relay/internal/profile"
)

// Launcher is responsible for securely and minimally starting a Tool
// connected to a Provider via a Profile. Relay must NEVER remain running
// after the child process is launched or handed off.
type Launcher interface {
	// Launch takes a profile, resolves the necessary environment variables,
	// and executes the corresponding tool. It returns an error if the process
	// fails to start.
	Launch(p profile.Profile) error
}
