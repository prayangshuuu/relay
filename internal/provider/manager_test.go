package provider

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/storage"
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

func setupTestManager(t *testing.T) (*Manager, string) {
	tempDir := t.TempDir()
	pm := &MockPathManager{base: tempDir}
	store := storage.NewLocalStorage()

	// Create required directories
	store.EnsureDir(pm.ProvidersDir())
	store.EnsureDir(pm.BackupsDir())

	return NewManager(pm, store), tempDir
}

func TestManager_AddGetListRemove(t *testing.T) {
	m, _ := setupTestManager(t)

	cfg := config.ProviderConfig{
		ID:   "test-prov",
		Name: "Test Prov",
		Type: "test",
	}

	// Add
	if err := m.Add(cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Add duplicate
	if err := m.Add(cfg); !errors.Is(err, ErrExists) {
		t.Fatalf("Expected ErrExists, got: %v", err)
	}

	// Get
	got, err := m.Get("test-prov")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Name != "Test Prov" {
		t.Errorf("Expected name 'Test Prov', got %s", got.Name)
	}

	// List
	list, err := m.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("Expected 1 provider, got %d", len(list))
	}

	// Remove
	if err := m.Remove("test-prov"); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	// Verify removed
	if _, err := m.Get("test-prov"); err == nil {
		t.Fatal("Expected error after remove, got nil")
	}
}

func TestManager_EditEnableDisable(t *testing.T) {
	m, _ := setupTestManager(t)

	cfg := config.ProviderConfig{
		ID:   "prov",
		Name: "Prov",
		Type: "test",
	}

	if err := m.Add(cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Enable
	if err := m.Enable("prov"); err != nil {
		t.Fatalf("Enable failed: %v", err)
	}

	got, _ := m.Get("prov")
	if !got.Enabled {
		t.Fatal("Expected provider to be enabled")
	}

	// Disable
	if err := m.Disable("prov"); err != nil {
		t.Fatalf("Disable failed: %v", err)
	}

	got, _ = m.Get("prov")
	if got.Enabled {
		t.Fatal("Expected provider to be disabled")
	}

	// Edit
	got.Name = "Updated"
	if err := m.Edit("prov", got); err != nil {
		t.Fatalf("Edit failed: %v", err)
	}

	got2, _ := m.Get("prov")
	if got2.Name != "Updated" {
		t.Fatal("Expected name to be updated")
	}
}

func TestManager_ImportExport(t *testing.T) {
	m, tempDir := setupTestManager(t)

	cfg := config.ProviderConfig{
		ID:   "export-test",
		Name: "Export Test",
		Type: "test",
	}
	m.Add(cfg)

	exportPath := filepath.Join(tempDir, "exported.yaml")
	if err := m.Export("export-test", exportPath, "yaml"); err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if _, err := os.Stat(exportPath); os.IsNotExist(err) {
		t.Fatal("Exported file not found")
	}

	// Remove original
	m.Remove("export-test")

	// Import
	if err := m.Import(exportPath); err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	// Verify import
	got, err := m.Get("export-test")
	if err != nil {
		t.Fatalf("Failed to get imported provider: %v", err)
	}
	if got.Name != "Export Test" {
		t.Errorf("Expected 'Export Test', got %s", got.Name)
	}
}
