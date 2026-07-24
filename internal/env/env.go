// Package env is responsible for securely handling environment variables.
// Future responsibilities include reading the system environment,
// creating temporary environment payloads for child processes,
// and ensuring that the host OS environment is never permanently modified.
package env

// Environment manages the creation and isolation of environment variables
// required by tools and providers.
type Environment interface {
	// Build takes a set of injected variables and merges them with the
	// current system environment, returning a slice formatted as "KEY=VALUE".
	// It guarantees no side effects to the host system.
	Build(injections map[string]string) []string

	// Clean ensures any sensitive temporary environment data is handled securely,
	// if necessary.
	Clean() error
}
