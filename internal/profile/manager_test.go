package profile

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/prayangshuuu/relay/internal/tool"
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

func setupTestManager(t *testing.T) (*Manager, *provider.Manager, *tool.Manager) {
	tempDir := t.TempDir()
	pm := &MockPathManager{base: tempDir}
	store := storage.NewLocalStorage()

	store.EnsureDir(pm.ProfilesDir())
	store.EnsureDir(pm.ProvidersDir())
	store.EnsureDir(pm.ToolsDir())

	provMgr := provider.NewManager(pm, store)
	toolMgr := tool.NewManager(pm, store)

	return NewManager(pm, store, provMgr, toolMgr), provMgr, toolMgr
}

func TestManager_AddGetListRemove(t *testing.T) {
	m, _, _ := setupTestManager(t)

	cfg := config.ProfileConfig{
		Name: "test-prof",
		Tool: "dummy",
	}

	if err := m.Add(cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if err := m.Add(cfg); !errors.Is(err, ErrExists) {
		t.Fatalf("Expected ErrExists, got: %v", err)
	}

	got, err := m.Get("test-prof")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Name() != "test-prof" {
		t.Errorf("Expected name 'test-prof', got %s", got.Name())
	}

	// Tool will be mocked by tool manager fallback in this simple test scenario
	if got.Tool() == nil || got.Tool().ExecutableName() != "dummy" {
		t.Errorf("Expected tool 'dummy', got %v", got.Tool())
	}

	list, err := m.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("Expected 1 profile, got %d", len(list))
	}

	if err := m.Remove("test-prof"); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	if _, err := m.Get("test-prof"); err == nil {
		t.Fatal("Expected error after remove, got nil")
	}
}

func TestManager_Clone(t *testing.T) {
	m, _, _ := setupTestManager(t)

	cfg := config.ProfileConfig{
		Name:  "src-prof",
		Model: "test-model",
	}

	m.Add(cfg)

	if err := m.Clone("src-prof", "cloned-prof"); err != nil {
		t.Fatalf("Clone failed: %v", err)
	}

	got, err := m.Get("cloned-prof")
	if err != nil {
		t.Fatalf("Get cloned failed: %v", err)
	}

	if got.Model() != "test-model" {
		t.Errorf("Expected cloned model 'test-model', got %s", got.Model())
	}
}

func TestManager_Validate(t *testing.T) {
	m, provMgr, _ := setupTestManager(t)

	cfg := config.ProfileConfig{
		Name:     "prof-val",
		Provider: "real-prov",
	}
	m.Add(cfg)

	// Since 'real-prov' is not created in provMgr, validation should fail
	err := m.Validate("prof-val")
	if err == nil {
		t.Fatal("Expected validation to fail due to missing provider")
	}

	// Create provider
	provMgr.Add(config.ProviderConfig{
		ID:   "real-prov",
		Name: "Real Prov",
		Type: "test",
	})

	err = m.Validate("prof-val")
	if err != nil {
		t.Fatalf("Expected validation to pass, got: %v", err)
	}
}
