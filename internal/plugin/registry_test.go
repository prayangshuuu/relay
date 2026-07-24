package plugin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/prayangshuuu/relay/pkg/sdk"
)

type MockPathManager struct {
	base string
}

func (m *MockPathManager) BaseDir() string      { return m.base }
func (m *MockPathManager) ConfigFile() string   { return filepath.Join(m.base, "config.yaml") }
func (m *MockPathManager) ProvidersDir() string { return filepath.Join(m.base, "providers") }
func (m *MockPathManager) ProfilesDir() string  { return filepath.Join(m.base, "profiles") }
func (m *MockPathManager) ToolsDir() string     { return filepath.Join(m.base, "tools") }
func (m *MockPathManager) BackupsDir() string   { return filepath.Join(m.base, "backups") }
func (m *MockPathManager) LogsDir() string      { return filepath.Join(m.base, "logs") }

func TestRegistry_BuiltIn(t *testing.T) {
	pm := &MockPathManager{base: t.TempDir()}
	r := NewRegistry(pm)
	RegisterBuiltIns(r)

	plugins := r.List()
	if len(plugins) != 3 {
		t.Fatalf("Expected 3 built-in plugins, got %d", len(plugins))
	}

	p, err := r.Get("anthropic")
	if err != nil {
		t.Fatalf("Expected to find anthropic plugin, got error: %v", err)
	}

	if p.Manifest().Type != sdk.TypeProvider {
		t.Errorf("Expected TypeProvider, got %v", p.Manifest().Type)
	}
}

func TestRegistry_DiscoverExternal(t *testing.T) {
	tempDir := t.TempDir()
	pm := &MockPathManager{base: tempDir}
	r := NewRegistry(pm)

	// Create mock external plugin
	pluginDir := filepath.Join(tempDir, "plugins", "my-plugin")
	os.MkdirAll(pluginDir, 0755)

	manifestYAML := `
id: my-plugin
name: My Plugin
version: 1.0.0
type: provider
entrypoint: run.sh
`
	os.WriteFile(filepath.Join(pluginDir, "plugin.yaml"), []byte(manifestYAML), 0644)

	err := r.Discover()
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	p, err := r.Get("my-plugin")
	if err != nil {
		t.Fatalf("Expected to find my-plugin, got error: %v", err)
	}

	if p.Manifest().ID != "my-plugin" {
		t.Errorf("Expected ID my-plugin, got %s", p.Manifest().ID)
	}
	if p.Manifest().Type != sdk.TypeProvider {
		t.Errorf("Expected TypeProvider, got %s", p.Manifest().Type)
	}
}
