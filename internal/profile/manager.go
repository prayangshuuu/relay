package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/prayangshuuu/relay/internal/tool"
)

var (
	ErrNotFound = errors.New("profile not found")
	ErrExists   = errors.New("profile already exists")
)

// DefaultProvider implements provider.Provider for the launcher.
type DefaultProvider struct {
	cfg config.ProviderConfig
}

func (p *DefaultProvider) Name() string { return p.cfg.Name }
func (p *DefaultProvider) EnvironmentVariables() map[string]string {
	return make(map[string]string)
}
func (p *DefaultProvider) Validate() error               { return provider.Validate(p.cfg) }
func (p *DefaultProvider) Metadata() map[string]string   { return p.cfg.Metadata }
func (p *DefaultProvider) Config() config.ProviderConfig { return p.cfg }

// DefaultProfile implements the Profile interface.
type DefaultProfile struct {
	cfg      config.ProfileConfig
	tool     tool.Tool
	provider provider.Provider
}

func (p *DefaultProfile) Name() string                 { return p.cfg.Name }
func (p *DefaultProfile) Tool() tool.Tool              { return p.tool }
func (p *DefaultProfile) Provider() provider.Provider  { return p.provider }
func (p *DefaultProfile) Model() string                { return p.cfg.Model }
func (p *DefaultProfile) Overrides() map[string]string { return p.cfg.Environment }
func (p *DefaultProfile) Config() config.ProfileConfig { return p.cfg }

// Manager handles CRUD for profiles.
type Manager struct {
	paths       config.PathManager
	storage     storage.Storage
	providerMgr *provider.Manager
	toolMgr     *tool.Manager
}

func NewManager(paths config.PathManager, s storage.Storage, pm *provider.Manager, tm *tool.Manager) *Manager {
	return &Manager{paths: paths, storage: s, providerMgr: pm, toolMgr: tm}
}

func (m *Manager) profileFile(name string) string {
	return filepath.Join(m.paths.ProfilesDir(), fmt.Sprintf("%s.yaml", name))
}

// Get Active or specific profile.
func (m *Manager) Get(name string) (Profile, error) {
	var cfg config.ProfileConfig
	path := m.profileFile(name)

	err := m.storage.ReadYAML(path, &cfg)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s", ErrNotFound, name)
		}
		return nil, err
	}

	// Resolve tool
	var t tool.Tool
	if cfg.Tool != "" {
		t, err = m.toolMgr.Get(cfg.Tool)
		if err != nil {
			return nil, fmt.Errorf("failed to load tool '%s' for profile '%s': %w", cfg.Tool, name, err)
		}
	}

	// Resolve provider
	var p provider.Provider
	if cfg.Provider != "" {
		pcfg, err := m.providerMgr.Get(cfg.Provider)
		if err != nil {
			return nil, fmt.Errorf("failed to load provider '%s' for profile '%s': %w", cfg.Provider, name, err)
		}
		p = &DefaultProvider{cfg: pcfg}
	}

	return &DefaultProfile{
		cfg:      cfg,
		tool:     t,
		provider: p,
	}, nil
}

// List returns all configured profiles.
func (m *Manager) List() ([]Profile, error) {
	files, err := os.ReadDir(m.paths.ProfilesDir())
	if err != nil {
		if os.IsNotExist(err) {
			return []Profile{}, nil
		}
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}

	var profiles []Profile
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		name := strings.TrimSuffix(file.Name(), ".yaml")
		p, err := m.Get(name)
		if err != nil {
			return nil, fmt.Errorf("failed to load profile %s: %w", name, err)
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// Add saves a new profile configuration.
func (m *Manager) Add(cfg config.ProfileConfig) error {
	if cfg.Name == "" {
		return errors.New("profile name cannot be empty")
	}

	var temp config.ProfileConfig
	err := m.storage.ReadYAML(m.profileFile(cfg.Name), &temp)
	if err == nil {
		return fmt.Errorf("%w: %s", ErrExists, cfg.Name)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return m.storage.WriteYAML(m.profileFile(cfg.Name), cfg)
}

// Edit updates an existing profile.
func (m *Manager) Edit(name string, cfg config.ProfileConfig) error {
	if name != cfg.Name {
		return errors.New("cannot change profile name during edit")
	}
	if _, err := m.Get(name); err != nil {
		return err
	}
	return m.storage.WriteYAML(m.profileFile(name), cfg)
}

// Remove deletes a profile.
func (m *Manager) Remove(name string) error {
	if _, err := m.Get(name); err != nil {
		return err
	}
	if err := os.Remove(m.profileFile(name)); err != nil {
		return fmt.Errorf("failed to remove profile file: %w", err)
	}
	return nil
}

// Clone copies an existing profile to a new one.
func (m *Manager) Clone(srcName, targetName string) error {
	srcProf, err := m.Get(srcName)
	if err != nil {
		return err
	}

	if _, err := m.Get(targetName); err == nil {
		return fmt.Errorf("%w: %s", ErrExists, targetName)
	} else if !errors.Is(err, ErrNotFound) {
		return err
	}

	// Get base config from the concrete type we defined
	if sp, ok := srcProf.(*DefaultProfile); ok {
		newCfg := sp.cfg
		newCfg.Name = targetName
		return m.Add(newCfg)
	}

	return errors.New("unsupported profile type for cloning")
}

// Validate checks if a profile's dependencies exist.
func (m *Manager) Validate(name string) error {
	prof, err := m.Get(name)
	if err != nil {
		return err
	}

	if sp, ok := prof.(*DefaultProfile); ok {
		if sp.cfg.Tool != "" {
			if _, err := m.toolMgr.Get(sp.cfg.Tool); err != nil {
				return fmt.Errorf("validation failed: tool '%s' is missing or invalid: %w", sp.cfg.Tool, err)
			}
		}

		if sp.cfg.Provider != "" {
			if _, err := m.providerMgr.Get(sp.cfg.Provider); err != nil {
				return fmt.Errorf("validation failed: provider '%s' is missing or invalid: %w", sp.cfg.Provider, err)
			}
		}
	}
	return nil
}
