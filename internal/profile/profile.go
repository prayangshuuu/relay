// Package profile connects Tools and Providers together.
// Future responsibilities include loading and managing user-defined profiles,
// applying specific models (e.g., Claude Sonnet 3.5), and handling overrides
// when switching tools.
package profile

import (
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/tool"
)

// Profile represents a specific configuration linking a Tool to a Provider.
// Example: "Work" profile connecting "Claude Code" (Tool) to "OpenRouter" (Provider)
// with the model "Claude Sonnet 4".
type Profile interface {
	// Name is the unique identifier for the profile.
	Name() string

	// Tool returns the CLI tool configured for this profile.
	Tool() tool.Tool

	// Provider returns the AI provider configured for this profile.
	Provider() provider.Provider

	// Model returns the specific AI model to be used.
	Model() string

	// Overrides returns any profile-specific environment variable overrides.
	Overrides() map[string]string
}
