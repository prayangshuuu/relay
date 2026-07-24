package plugin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/pkg/sdk"
	"gopkg.in/yaml.v3"
)

type Registry struct {
	paths   config.PathManager
	plugins map[string]sdk.Plugin
}

func NewRegistry(paths config.PathManager) *Registry {
	return &Registry{
		paths:   paths,
		plugins: make(map[string]sdk.Plugin),
	}
}

// RegisterBuiltIn allows the CLI to inject native Go plugins into the registry.
func (r *Registry) RegisterBuiltIn(p sdk.Plugin) {
	r.plugins[p.Manifest().ID] = p
}

// Discover scans all plugin directories and loads manifests.
func (r *Registry) Discover() error {
	dirs := []string{
		filepath.Join(r.paths.BaseDir(), "plugins"),
		".relay/plugins",
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			continue // skip unreadable dirs
		}

		for _, entry := range entries {
			if entry.IsDir() {
				manifestPath := filepath.Join(dir, entry.Name(), "plugin.yaml")
				if _, err := os.Stat(manifestPath); err == nil {
					data, err := os.ReadFile(manifestPath)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Warning: could not read plugin manifest %s\n", manifestPath)
						continue
					}

					var manifest sdk.Manifest
					if err := yaml.Unmarshal(data, &manifest); err != nil {
						fmt.Fprintf(os.Stderr, "Warning: could not parse plugin manifest %s\n", manifestPath)
						continue
					}

					// For external plugins, we create a proxy wrapper that handles RPC execution.
					// For now, we'll just register a stub proxy plugin so `plugin list` works.
					r.plugins[manifest.ID] = &RPCProxyPlugin{manifest: manifest, binaryPath: filepath.Join(dir, entry.Name(), manifest.Entrypoint)}
				}
			}
		}
	}

	return nil
}

func (r *Registry) List() []sdk.Plugin {
	var list []sdk.Plugin
	for _, p := range r.plugins {
		list = append(list, p)
	}
	return list
}

func (r *Registry) Get(id string) (sdk.Plugin, error) {
	if p, ok := r.plugins[id]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("plugin %s not found", id)
}

// RPCProxyPlugin represents an external executable plugin.
type RPCProxyPlugin struct {
	manifest   sdk.Manifest
	binaryPath string
}

func (p *RPCProxyPlugin) Manifest() sdk.Manifest {
	return p.manifest
}

// Implement ProviderPlugin stub
func (p *RPCProxyPlugin) InjectEnvironment(instanceConfig map[string]interface{}) (map[string]string, error) {
	// In the future: exec.Command(p.binaryPath) with JSON-RPC
	return nil, fmt.Errorf("RPC execution not implemented yet")
}

// Implement ToolPlugin stub
func (p *RPCProxyPlugin) Executable() string {
	return ""
}

func (p *RPCProxyPlugin) FormatLaunchArgs(profileConfig map[string]interface{}) ([]string, error) {
	return nil, fmt.Errorf("RPC execution not implemented yet")
}
