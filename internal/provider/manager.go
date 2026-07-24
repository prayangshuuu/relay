package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/storage"
	"gopkg.in/yaml.v3"
)

var (
	ErrNotFound = errors.New("provider not found")
	ErrExists   = errors.New("provider already exists")
)

// Manager handles CRUD operations for providers.
type Manager struct {
	paths   config.PathManager
	storage storage.Storage
}

// NewManager creates a new provider manager.
func NewManager(paths config.PathManager, s storage.Storage) *Manager {
	return &Manager{
		paths:   paths,
		storage: s,
	}
}

func (m *Manager) providerFile(id string) string {
	return filepath.Join(m.paths.ProvidersDir(), fmt.Sprintf("%s.yaml", id))
}

// List returns all configured providers.
func (m *Manager) List() ([]config.ProviderConfig, error) {
	files, err := os.ReadDir(m.paths.ProvidersDir())
	if err != nil {
		if os.IsNotExist(err) {
			return []config.ProviderConfig{}, nil
		}
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}

	var providers []config.ProviderConfig
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		id := strings.TrimSuffix(file.Name(), ".yaml")
		cfg, err := m.Get(id)
		if err != nil {
			// skip malformed ones in list, or return error? Let's return error
			return nil, fmt.Errorf("failed to load provider %s: %w", id, err)
		}
		providers = append(providers, cfg)
	}
	return providers, nil
}

// Get retrieves a provider by ID.
func (m *Manager) Get(id string) (config.ProviderConfig, error) {
	var cfg config.ProviderConfig
	err := m.storage.ReadYAML(m.providerFile(id), &cfg)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, fmt.Errorf("%w: %s", ErrNotFound, id)
		}
		return cfg, err
	}
	return cfg, nil
}

// Add saves a new provider configuration.
func (m *Manager) Add(cfg config.ProviderConfig) error {
	if err := Validate(cfg); err != nil {
		return err
	}

	if _, err := m.Get(cfg.ID); err == nil {
		return fmt.Errorf("%w: %s", ErrExists, cfg.ID)
	} else if !errors.Is(err, ErrNotFound) {
		return err // other error occurred
	}

	now := time.Now().Format(time.RFC3339)
	cfg.CreatedAt = now
	cfg.UpdatedAt = now

	return m.storage.WriteYAML(m.providerFile(cfg.ID), cfg)
}

// backup creates a backup of a provider file in the backups directory.
func (m *Manager) backup(id string) error {
	cfg, err := m.Get(id)
	if err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102150405")
	backupPath := filepath.Join(m.paths.BackupsDir(), fmt.Sprintf("provider_%s_%s.yaml", id, timestamp))
	return m.storage.WriteYAML(backupPath, cfg)
}

// Edit updates an existing provider.
func (m *Manager) Edit(id string, cfg config.ProviderConfig) error {
	if id != cfg.ID {
		return errors.New("cannot change provider ID during edit")
	}

	if err := Validate(cfg); err != nil {
		return err
	}

	existing, err := m.Get(id)
	if err != nil {
		return err
	}

	if err := m.backup(id); err != nil {
		return fmt.Errorf("failed to backup before edit: %w", err)
	}

	cfg.CreatedAt = existing.CreatedAt
	cfg.UpdatedAt = time.Now().Format(time.RFC3339)

	return m.storage.WriteYAML(m.providerFile(id), cfg)
}

// Remove deletes a provider.
func (m *Manager) Remove(id string) error {
	if _, err := m.Get(id); err != nil {
		return err
	}

	if err := m.backup(id); err != nil {
		return fmt.Errorf("failed to backup before remove: %w", err)
	}

	if err := os.Remove(m.providerFile(id)); err != nil {
		return fmt.Errorf("failed to remove provider file: %w", err)
	}
	return nil
}

// Enable enables a provider.
func (m *Manager) Enable(id string) error {
	cfg, err := m.Get(id)
	if err != nil {
		return err
	}
	cfg.Enabled = true
	return m.Edit(id, cfg)
}

// Disable disables a provider.
func (m *Manager) Disable(id string) error {
	cfg, err := m.Get(id)
	if err != nil {
		return err
	}
	cfg.Enabled = false
	return m.Edit(id, cfg)
}

// Export saves the provider config to the specified path in yaml or json format.
func (m *Manager) Export(id string, dest string, format string) error {
	cfg, err := m.Get(id)
	if err != nil {
		return err
	}

	var data []byte
	switch strings.ToLower(format) {
	case "json":
		data, err = json.MarshalIndent(cfg, "", "  ")
	case "yaml", "yml":
		data, err = yaml.Marshal(cfg)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal provider: %w", err)
	}

	return os.WriteFile(dest, data, 0644)
}

// Import reads a provider configuration from a file and adds it.
func (m *Manager) Import(src string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	var cfg config.ProviderConfig
	if strings.HasSuffix(strings.ToLower(src), ".json") {
		err = json.Unmarshal(data, &cfg)
	} else {
		err = yaml.Unmarshal(data, &cfg)
	}

	if err != nil {
		return fmt.Errorf("failed to unmarshal import data: %w", err)
	}

	return m.Add(cfg)
}
