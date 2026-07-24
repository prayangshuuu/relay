// Package config defines the configuration formats and YAML schemas for Relay.
package config

// GlobalConfig represents the schema for config.yaml
type GlobalConfig struct {
	Version         int      `yaml:"version"`
	CurrentProfile  string   `yaml:"current_profile"`
	DefaultTool     string   `yaml:"default_tool"`
	RecentProfiles  []string `yaml:"recent_profiles,omitempty"`
	RecentProviders []string `yaml:"recent_providers,omitempty"`
	RecentTools     []string `yaml:"recent_tools,omitempty"`
	UndoStack       []string `yaml:"undo_stack,omitempty"`
	CreatedAt       string   `yaml:"created_at,omitempty"`
	UpdatedAt       string   `yaml:"updated_at,omitempty"`
}

// ProviderConfig represents the schema for files in providers/
type ProviderConfig struct {
	ID                   string            `yaml:"id"`
	Name                 string            `yaml:"name"`
	Type                 string            `yaml:"type"` // e.g., "anthropic", "openrouter"
	BaseURL              string            `yaml:"base_url,omitempty"`
	AuthenticationType   string            `yaml:"authentication_type,omitempty"`
	EnvironmentVariables []string          `yaml:"environment_variables,omitempty"`
	DefaultHeaders       map[string]string `yaml:"default_headers,omitempty"`
	Metadata             map[string]string `yaml:"metadata,omitempty"`
	CreatedAt            string            `yaml:"created_at,omitempty"`
	UpdatedAt            string            `yaml:"updated_at,omitempty"`
	Enabled              bool              `yaml:"enabled"`
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
	Name           string            `yaml:"name"`
	Tool           string            `yaml:"tool"`
	Provider       string            `yaml:"provider"`
	Model          string            `yaml:"model"`
	Environment    map[string]string `yaml:"environment,omitempty"`
	WorkspaceRules map[string]string `yaml:"workspace_rules,omitempty"`
	LaunchArgs     []string          `yaml:"launch_args,omitempty"`
}

// WorkspaceConfig represents overrides for a specific project directory.
type WorkspaceConfig struct {
	Path        string            `yaml:"path"`
	Profile     string            `yaml:"profile"`
	Environment map[string]string `yaml:"environment,omitempty"`
}
