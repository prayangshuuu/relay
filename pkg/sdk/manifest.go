package sdk

type PluginType string

const (
	TypeProvider PluginType = "provider"
	TypeTool     PluginType = "tool"
)

type Manifest struct {
	ID                  string     `yaml:"id" json:"id"`
	Name                string     `yaml:"name" json:"name"`
	Version             string     `yaml:"version" json:"version"`
	Author              string     `yaml:"author" json:"author"`
	Description         string     `yaml:"description" json:"description"`
	MinimumRelayVersion string     `yaml:"minimum_relay_version" json:"minimum_relay_version"`
	SupportedPlatforms  []string   `yaml:"supported_platforms" json:"supported_platforms"`
	Type                PluginType `yaml:"type" json:"type"`
	Entrypoint          string     `yaml:"entrypoint" json:"entrypoint"`
}
