package sdk

// Plugin is the base interface that all Relay plugins must implement.
type Plugin interface {
	// Manifest returns the plugin's metadata.
	Manifest() Manifest
}

// ProviderPlugin defines how Relay interacts with AI model providers.
type ProviderPlugin interface {
	Plugin
	// InjectEnvironment returns the environment variables needed to authenticate and configure the provider for a given instance.
	InjectEnvironment(instanceConfig map[string]interface{}) (map[string]string, error)
}

// ToolPlugin defines how Relay interacts with AI development tools (e.g. Claude Code, Aider).
type ToolPlugin interface {
	Plugin
	// Executable returns the base executable name or path for the tool.
	Executable() string
	// FormatLaunchArgs returns the CLI arguments required to launch the tool using the provided profile configuration.
	FormatLaunchArgs(profileConfig map[string]interface{}) ([]string, error)
}
