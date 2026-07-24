package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// PathManager defines how Relay resolves its configuration paths.
type PathManager interface {
	// BaseDir returns the root directory for Relay's configuration.
	BaseDir() string

	// ConfigFile returns the path to the main config.yaml file.
	ConfigFile() string

	// ProvidersDir returns the path to the providers directory.
	ProvidersDir() string

	// ProfilesDir returns the path to the profiles directory.
	ProfilesDir() string

	// ToolsDir returns the path to the tools directory.
	ToolsDir() string

	// BackupsDir returns the path to the backups directory.
	BackupsDir() string

	// LogsDir returns the path to the logs directory.
	LogsDir() string
}

// OSPathManager implements PathManager using OS-specific standards.
type OSPathManager struct {
	base string
}

// NewOSPathManager creates a path manager based on os.UserConfigDir.
func NewOSPathManager() (*OSPathManager, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	var relayDirName string
	if runtime.GOOS == "linux" {
		relayDirName = "relay"
	} else {
		relayDirName = "Relay"
	}

	return &OSPathManager{
		base: filepath.Join(configDir, relayDirName),
	}, nil
}

func (p *OSPathManager) BaseDir() string      { return p.base }
func (p *OSPathManager) ConfigFile() string   { return filepath.Join(p.base, "config.yaml") }
func (p *OSPathManager) ProvidersDir() string { return filepath.Join(p.base, "providers") }
func (p *OSPathManager) ProfilesDir() string  { return filepath.Join(p.base, "profiles") }
func (p *OSPathManager) ToolsDir() string     { return filepath.Join(p.base, "tools") }
func (p *OSPathManager) BackupsDir() string   { return filepath.Join(p.base, "backups") }
func (p *OSPathManager) LogsDir() string      { return filepath.Join(p.base, "logs") }
