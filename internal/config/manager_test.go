package config

import (
	"path/filepath"
	"testing"

	"github.com/prayangshuuu/relay/internal/storage"
)

// MockPathManager implements PathManager for testing
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

func TestManager_Initialize(t *testing.T) {
	tempDir := t.TempDir()
	pm := &MockPathManager{base: tempDir}
	store := storage.NewLocalStorage()
	manager := NewManager(pm, store)

	if manager.Exists() {
		t.Fatal("Expected manager not to exist before initialization")
	}

	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	if !manager.Exists() {
		t.Fatal("Expected manager to exist after initialization")
	}

	if err := manager.Initialize(); err != ErrAlreadyInitialized {
		t.Fatalf("Expected ErrAlreadyInitialized, got: %v", err)
	}
}

func TestManager_LoadAndSave(t *testing.T) {
	tempDir := t.TempDir()
	pm := &MockPathManager{base: tempDir}
	store := storage.NewLocalStorage()
	manager := NewManager(pm, store)

	// Attempting to load before initialize should fail
	_, err := manager.Load()
	if err != ErrNotInitialized {
		t.Fatalf("Expected ErrNotInitialized, got: %v", err)
	}

	err = manager.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	cfg, err := manager.Load()
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if cfg.Version != 1 {
		t.Errorf("Expected default version 1, got %d", cfg.Version)
	}

	cfg.DefaultTool = "aider"
	if err := manager.Save(cfg); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	loadedCfg, err := manager.Load()
	if err != nil {
		t.Fatalf("Failed to reload: %v", err)
	}
	if loadedCfg.DefaultTool != "aider" {
		t.Errorf("Expected tool to be updated to 'aider', got '%s'", loadedCfg.DefaultTool)
	}
}

func TestManager_Validate(t *testing.T) {
	manager := NewManager(nil, nil)

	cfg := &GlobalConfig{
		Version:     1,
		DefaultTool: "claude",
	}

	if err := manager.Validate(cfg); err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}

	cfg.Version = 0
	if err := manager.Validate(cfg); err == nil {
		t.Error("Expected error for version < 1")
	}

	cfg.Version = 1
	cfg.DefaultTool = ""
	if err := manager.Validate(cfg); err == nil {
		t.Error("Expected error for empty default_tool")
	}
}
