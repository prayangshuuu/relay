package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/prayangshuuu/relay/internal/storage"
)

var (
	ErrAlreadyInitialized = errors.New("relay is already initialized")
	ErrNotInitialized     = errors.New("relay is not initialized")
	ErrInvalidConfig      = errors.New("invalid configuration")
)

// Manager orchestrates the configuration loading, saving, and initialization.
type Manager struct {
	paths   PathManager
	storage storage.Storage
}

// NewManager creates a new Manager instance.
func NewManager(paths PathManager, s storage.Storage) *Manager {
	return &Manager{
		paths:   paths,
		storage: s,
	}
}

// Exists checks if the configuration directory exists.
func (m *Manager) Exists() bool {
	info, err := os.Stat(m.paths.ConfigFile())
	return err == nil && !info.IsDir()
}

// Initialize creates the initial directory structure and a default config.yaml.
func (m *Manager) Initialize() error {
	if m.Exists() {
		return ErrAlreadyInitialized
	}

	// Create directories
	dirs := []string{
		m.paths.BaseDir(),
		m.paths.ProvidersDir(),
		m.paths.ProfilesDir(),
		m.paths.ToolsDir(),
		m.paths.BackupsDir(),
		m.paths.LogsDir(),
	}

	for _, dir := range dirs {
		if err := m.storage.EnsureDir(dir); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create default config
	defaultConfig := GlobalConfig{
		Version:        1,
		CurrentProfile: "default",
		DefaultTool:    "claude",
		CreatedAt:      time.Now().Format(time.RFC3339),
		UpdatedAt:      time.Now().Format(time.RFC3339),
	}

	if err := m.storage.WriteYAML(m.paths.ConfigFile(), defaultConfig); err != nil {
		return fmt.Errorf("failed to create default config.yaml: %w", err)
	}

	return nil
}

// Load reads the config.yaml file.
func (m *Manager) Load() (*GlobalConfig, error) {
	if !m.Exists() {
		return nil, ErrNotInitialized
	}

	var cfg GlobalConfig
	if err := m.storage.ReadYAML(m.paths.ConfigFile(), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes the global configuration.
func (m *Manager) Save(cfg *GlobalConfig) error {
	if err := m.Validate(cfg); err != nil {
		return err
	}

	cfg.UpdatedAt = time.Now().Format(time.RFC3339)
	return m.storage.WriteYAML(m.paths.ConfigFile(), cfg)
}

// Validate ensures the GlobalConfig is well-formed.
func (m *Manager) Validate(cfg *GlobalConfig) error {
	if cfg.Version < 1 {
		return fmt.Errorf("%w: version must be >= 1", ErrInvalidConfig)
	}
	if cfg.DefaultTool == "" {
		return fmt.Errorf("%w: default_tool cannot be empty", ErrInvalidConfig)
	}
	return nil
}

// CreateBackup creates a timestamped copy of the current config.yaml.
func (m *Manager) CreateBackup() error {
	if !m.Exists() {
		return ErrNotInitialized
	}

	cfg, err := m.Load()
	if err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102150405")
	backupPath := filepath.Join(m.paths.BackupsDir(), fmt.Sprintf("config_%s.yaml", timestamp))

	return m.storage.WriteYAML(backupPath, cfg)
}

// RestoreBackup restores a backup configuration file.
func (m *Manager) RestoreBackup(filename string) error {
	backupPath := filepath.Join(m.paths.BackupsDir(), filename)
	var backupCfg GlobalConfig
	if err := m.storage.ReadYAML(backupPath, &backupCfg); err != nil {
		return fmt.Errorf("failed to load backup %s: %w", filename, err)
	}

	if err := m.Validate(&backupCfg); err != nil {
		return fmt.Errorf("backup validation failed: %w", err)
	}

	return m.storage.WriteYAML(m.paths.ConfigFile(), backupCfg)
}
