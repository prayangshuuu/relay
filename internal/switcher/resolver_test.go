package switcher

import (
	"os"
	"path/filepath"
	"testing"
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

func TestResolver_SingleMatch(t *testing.T) {
	tempDir := t.TempDir()
	pm := &MockPathManager{base: tempDir}

	os.MkdirAll(pm.ProfilesDir(), 0755)
	os.WriteFile(filepath.Join(pm.ProfilesDir(), "work.yaml"), []byte("name: work"), 0644)

	resolver := NewResolver(pm)

	match, err := resolver.Resolve("work")
	if err != nil {
		t.Fatalf("Expected match, got error: %v", err)
	}

	if match.Type != TypeProfile {
		t.Errorf("Expected TypeProfile, got %v", match.Type)
	}
}

func TestResolver_NoMatch(t *testing.T) {
	tempDir := t.TempDir()
	pm := &MockPathManager{base: tempDir}

	resolver := NewResolver(pm)

	_, err := resolver.Resolve("missing")
	if err == nil {
		t.Fatal("Expected error for missing entity, got nil")
	}
}
