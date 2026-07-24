// Package config defines the configuration formats and YAML schemas for Relay.
// Future responsibilities include parsing global configuration files located in
// ~/.relay/, and maintaining the schema representations for providers, profiles,
// and tools without implementing the execution logic.
package config

// GlobalConfig represents the schema for ~/.relay/config.yaml
type GlobalConfig struct {
	DefaultProfile string `yaml:"default_profile"`
	Telemetry      bool   `yaml:"telemetry"` // Expected to be false per architecture principles
}

// ProviderConfig represents the schema for files in ~/.relay/providers/
type ProviderConfig struct {
	Name     string            `yaml:"name"`
	Type     string            `yaml:"type"` // e.g., "anthropic", "openrouter"
	EnvVars  map[string]string `yaml:"env"`
	Metadata map[string]string `yaml:"metadata,omitempty"`
}

// ToolConfig represents the schema for files in ~/.relay/tools/
type ToolConfig struct {
	Name         string   `yaml:"name"`
	Executable   string   `yaml:"executable"`
	SupportedEnv []string `yaml:"supported_env"`
	LaunchMethod string   `yaml:"launch_method"`
}

// ProfileConfig represents the schema for files in ~/.relay/profiles/
type ProfileConfig struct {
	Name      string            `yaml:"name"`
	Tool      string            `yaml:"tool"`
	Provider  string            `yaml:"provider"`
	Model     string            `yaml:"model"`
	Overrides map[string]string `yaml:"overrides,omitempty"`
}
