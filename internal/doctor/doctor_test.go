package doctor

import (
	"path/filepath"
	"testing"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/profile"
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

func TestDoctor_MissingDirs(t *testing.T) {
	tempDir := t.TempDir()
	pm := &MockPathManager{base: tempDir}
	store := storage.NewLocalStorage()
	cfgMgr := config.NewManager(pm, store)
	provMgr := provider.NewManager(pm, store)
	toolMgr := tool.NewManager(pm, store)
	profMgr := profile.NewManager(pm, store, provMgr, toolMgr)

	suite := NewCheckSuite(pm, store, cfgMgr, profMgr, provMgr, toolMgr)
	results := suite.RunChecks()

	errCount := 0
	for _, r := range results {
		if r.Status == StatusError {
			errCount++
		}
	}

	if errCount == 0 {
		t.Fatalf("Expected errors for missing directories, got 0")
	}

	// Now run repair
	repairs := suite.Repair()
	if len(repairs) == 0 {
		t.Fatalf("Expected repairs, got 0")
	}

	resultsAfter := suite.RunChecks()
	errCountAfter := 0
	for _, r := range resultsAfter {
		if r.Status == StatusError && r.Name == "Config Directory" {
			errCountAfter++
		}
	}
	if errCountAfter > 0 {
		t.Fatalf("Expected 0 errors for dirs after repair, got %d", errCountAfter)
	}
}
