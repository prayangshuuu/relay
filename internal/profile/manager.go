package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/prayangshuuu/relay/internal/tool"
)

var (
	ErrNotFound = errors.New("profile not found")
)

// DefaultProvider implements provider.Provider for the launcher.
type DefaultProvider struct {
	cfg config.ProviderConfig
}

func (p *DefaultProvider) Name() string { return p.cfg.Name }
func (p *DefaultProvider) EnvironmentVariables() map[string]string {
	// For launcher, the values actually come from the environment injection process,
	// but this could store default values if needed.
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

// Get Active or specific profile.
func (m *Manager) Get(name string) (Profile, error) {
	var cfg config.ProfileConfig
	path := filepath.Join(m.paths.ProfilesDir(), fmt.Sprintf("%s.yaml", name))

	err := m.storage.ReadYAML(path, &cfg)
	if err != nil {
		if os.IsNotExist(err) {
			// Fallback: If not found, use a sensible empty profile to allow raw provider running
			return &DefaultProfile{
				cfg: config.ProfileConfig{
					Name: "default",
				},
			}, nil
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
