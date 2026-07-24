package plugin

import (
	"github.com/prayangshuuu/relay/pkg/sdk"
)

// RegisterBuiltIns injects the default native plugins into the registry.
func RegisterBuiltIns(r *Registry) {
	r.RegisterBuiltIn(&builtinProviderPlugin{
		manifest: sdk.Manifest{
			ID:          "anthropic",
			Name:        "Anthropic",
			Version:     "1.0.0",
			Author:      "Relay Core Team",
			Description: "Built-in Anthropic provider",
			Type:        sdk.TypeProvider,
		},
	})

	r.RegisterBuiltIn(&builtinProviderPlugin{
		manifest: sdk.Manifest{
			ID:          "openrouter",
			Name:        "OpenRouter",
			Version:     "1.0.0",
			Author:      "Relay Core Team",
			Description: "Built-in OpenRouter provider",
			Type:        sdk.TypeProvider,
		},
	})

	r.RegisterBuiltIn(&builtinToolPlugin{
		manifest: sdk.Manifest{
			ID:          "claude-code",
			Name:        "Claude Code",
			Version:     "1.0.0",
			Author:      "Relay Core Team",
			Description: "Built-in Claude Code tool wrapper",
			Type:        sdk.TypeTool,
		},
	})
}

type builtinProviderPlugin struct {
	manifest sdk.Manifest
}

func (p *builtinProviderPlugin) Manifest() sdk.Manifest {
	return p.manifest
}

func (p *builtinProviderPlugin) InjectEnvironment(instanceConfig map[string]interface{}) (map[string]string, error) {
	// Built-in logic goes here
	return make(map[string]string), nil
}

type builtinToolPlugin struct {
	manifest sdk.Manifest
}

func (p *builtinToolPlugin) Manifest() sdk.Manifest {
	return p.manifest
}

func (p *builtinToolPlugin) Executable() string {
	return "claude"
}

func (p *builtinToolPlugin) FormatLaunchArgs(profileConfig map[string]interface{}) ([]string, error) {
	return []string{}, nil
}
