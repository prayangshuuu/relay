package tool

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/storage"
)

var (
	ErrNotFound = errors.New("tool not found")
)

// DefaultTool implements the Tool interface based on config.ToolConfig
type DefaultTool struct {
	cfg config.ToolConfig
}

func (t *DefaultTool) ExecutableName() string                  { return t.cfg.Executable }
func (t *DefaultTool) SupportedEnvironmentVariables() []string { return []string{} } // Extendable later
func (t *DefaultTool) LaunchMethod() string                    { return "exec" }

// Manager handles CRUD for tools
type Manager struct {
	paths   config.PathManager
	storage storage.Storage
}

func NewManager(paths config.PathManager, s storage.Storage) *Manager {
	return &Manager{paths: paths, storage: s}
}

func (m *Manager) Get(name string) (Tool, error) {
	var cfg config.ToolConfig
	path := filepath.Join(m.paths.ToolsDir(), fmt.Sprintf("%s.yaml", name))
	err := m.storage.ReadYAML(path, &cfg)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Fallback: If not found, create an on-the-fly tool based on the name.
			// This allows 'relay run claude' to work without a strict tools/claude.yaml
			return &DefaultTool{cfg: config.ToolConfig{
				Name:       name,
				Executable: name,
			}}, nil
		}
		return nil, err
	}
	return &DefaultTool{cfg: cfg}, nil
}
