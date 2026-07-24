// Package tool provides the abstraction for AI CLI tools.
// Future responsibilities include defining how tools like Claude Code,
// Codex CLI, or Aider are executed, what environment variables they require,
// and how they are configured via Relay.
package tool

// Tool represents an AI CLI application (e.g., Claude Code, Aider, Codex CLI).
// It acts as the bridge between Relay and the actual executable being launched.
type Tool interface {
	// ExecutableName returns the base name of the executable (e.g., "claude", "aider").
	ExecutableName() string

	// SupportedEnvironmentVariables returns a list of environment variable names
	// that this tool recognizes or requires.
	SupportedEnvironmentVariables() []string

	// LaunchMethod defines how the tool should be executed.
	// This could indicate standard execution, shell wrapping, or
	// other specialized launch behaviors.
	LaunchMethod() string
}
