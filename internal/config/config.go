// Package config defines the configuration formats and YAML schemas for Relay.
package config

// GlobalConfig represents the schema for config.yaml
type GlobalConfig struct {
	Version        int    `yaml:"version"`
	CurrentProfile string `yaml:"current_profile"`
	DefaultTool    string `yaml:"default_tool"`
	CreatedAt      string `yaml:"created_at,omitempty"`
	UpdatedAt      string `yaml:"updated_at,omitempty"`
}

// ProviderConfig represents the schema for files in providers/
type ProviderConfig struct {
	Name        string            `yaml:"name"`
	Type        string            `yaml:"type"` // e.g., "anthropic", "openrouter"
	BaseURL     string            `yaml:"base_url,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Metadata    map[string]string `yaml:"metadata,omitempty"`
}

// ToolConfig represents the schema for files in tools/
type ToolConfig struct {
	Name        string            `yaml:"name"`
	Executable  string            `yaml:"executable"`
	Arguments   []string          `yaml:"arguments,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
}

// ProfileConfig represents the schema for files in profiles/
type ProfileConfig struct {
	Name        string            `yaml:"name"`
	Tool        string            `yaml:"tool"`
	Provider    string            `yaml:"provider"`
	Model       string            `yaml:"model"`
	Environment map[string]string `yaml:"environment,omitempty"`
}
