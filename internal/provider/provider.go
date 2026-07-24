// Package provider defines the abstractions for AI model API endpoints.
// Future responsibilities include modeling API communication, defining
// required configurations for endpoints (e.g., Anthropic, OpenRouter, OpenAI),
// and handling provider-specific environment variables.
package provider

// Provider represents an AI API endpoint (e.g., Anthropic, OpenRouter, OpenAI).
// It defines the contract that all concrete provider implementations must satisfy.
type Provider interface {
	// Name returns the provider's identifier (e.g., "anthropic", "openrouter").
	Name() string

	// EnvironmentVariables returns the required environment variables
	// for this provider to function properly.
	EnvironmentVariables() map[string]string

	// Validate checks if the provider's configuration is valid
	// (e.g., API keys are present and correctly formatted).
	Validate() error

	// Metadata returns additional information about the provider.
	Metadata() map[string]string
}
