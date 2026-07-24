package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/profile"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/prayangshuuu/relay/internal/tool"
)

type Status string

const (
	StatusOK      Status = "OK"
	StatusWarning Status = "WARNING"
	StatusError   Status = "ERROR"
)

type Result struct {
	Name    string
	Status  Status
	Message string
	Fixable bool
}

type CheckSuite struct {
	paths   config.PathManager
	storage storage.Storage
	cfgMgr  *config.Manager
	profMgr *profile.Manager
	provMgr *provider.Manager
	toolMgr *tool.Manager
}

func NewCheckSuite(paths config.PathManager, s storage.Storage, c *config.Manager, p *profile.Manager, pr *provider.Manager, t *tool.Manager) *CheckSuite {
	return &CheckSuite{paths, s, c, p, pr, t}
}

func (c *CheckSuite) RunChecks() []Result {
	var results []Result

	// 1. Check directories
	dirs := []struct {
		name string
		path string
	}{
		{"Config", c.paths.BaseDir()},
		{"Profiles", c.paths.ProfilesDir()},
		{"Providers", c.paths.ProvidersDir()},
		{"Tools", c.paths.ToolsDir()},
		{"Backups", c.paths.BackupsDir()},
	}

	for _, d := range dirs {
		if _, err := os.Stat(d.path); os.IsNotExist(err) {
			results = append(results, Result{Name: fmt.Sprintf("%s Directory", d.name), Status: StatusError, Message: "Missing", Fixable: true})
		} else {
			results = append(results, Result{Name: fmt.Sprintf("%s Directory", d.name), Status: StatusOK, Message: "Found", Fixable: false})
		}
	}

	// 2. Global Config
	globalCfg, err := c.cfgMgr.Load()
	if err != nil {
		results = append(results, Result{Name: "Global Configuration", Status: StatusError, Message: "Missing or Invalid YAML", Fixable: true})
	} else {
		results = append(results, Result{Name: "Global Configuration", Status: StatusOK, Message: "Valid", Fixable: false})

		// Version Check
		if globalCfg.Version == 0 {
			results = append(results, Result{Name: "Schema Version", Status: StatusWarning, Message: "Missing version", Fixable: true})
		} else {
			results = append(results, Result{Name: "Schema Version", Status: StatusOK, Message: fmt.Sprintf("v%d", globalCfg.Version), Fixable: false})
		}

		// 3. Active Profile
		if globalCfg.CurrentProfile == "" {
			results = append(results, Result{Name: "Active Profile", Status: StatusWarning, Message: "No active profile set", Fixable: true})
		} else {
			_, err := c.profMgr.Get(globalCfg.CurrentProfile)
			if err != nil {
				results = append(results, Result{Name: "Active Profile", Status: StatusError, Message: fmt.Sprintf("Profile '%s' is missing", globalCfg.CurrentProfile), Fixable: true})
			} else {
				results = append(results, Result{Name: "Active Profile", Status: StatusOK, Message: globalCfg.CurrentProfile, Fixable: false})
			}
		}
	}

	// 4. Validate all profiles
	profiles, _ := c.profMgr.List()
	for _, p := range profiles {
		// Tool Check
		if p.Tool() != nil {
			exe := p.Tool().ExecutableName()
			if path, err := exec.LookPath(exe); err == nil {
				// Try to get version
				out, _ := exec.Command(path, "--version").Output()
				vStr := strings.TrimSpace(string(out))
				if vStr == "" {
					vStr = "Found"
				} else {
					vStr = strings.Split(vStr, "\n")[0] // Just first line
				}
				results = append(results, Result{Name: fmt.Sprintf("Profile '%s' Tool (%s)", p.Name(), exe), Status: StatusOK, Message: vStr, Fixable: false})
			} else {
				results = append(results, Result{Name: fmt.Sprintf("Profile '%s' Tool (%s)", p.Name(), exe), Status: StatusError, Message: "Executable not found in PATH", Fixable: false})
			}
		}

		// Provider Check
		if p.Provider() != nil {
			results = append(results, Result{Name: fmt.Sprintf("Profile '%s' Provider (%s)", p.Name(), p.Provider().Name()), Status: StatusOK, Message: "Configured", Fixable: false})
		}
	}

	return results
}

func (c *CheckSuite) Repair() []Result {
	var repairs []Result

	// Create missing dirs
	dirs := []string{
		c.paths.BaseDir(),
		c.paths.ProfilesDir(),
		c.paths.ProvidersDir(),
		c.paths.ToolsDir(),
		c.paths.BackupsDir(),
	}
	for _, d := range dirs {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			os.MkdirAll(d, 0755)
			repairs = append(repairs, Result{Name: "Created Directory", Status: StatusOK, Message: d})
		}
	}

	// Fix global config
	globalCfg, err := c.cfgMgr.Load()
	if err != nil {
		c.cfgMgr.Initialize()
		repairs = append(repairs, Result{Name: "Initialized Global Config", Status: StatusOK, Message: "Fixed"})
		globalCfg, _ = c.cfgMgr.Load()
	}

	// Fix schema version
	if globalCfg.Version == 0 {
		globalCfg.Version = 1
		c.cfgMgr.Save(globalCfg)
		repairs = append(repairs, Result{Name: "Upgraded Schema Version", Status: StatusOK, Message: "v1"})
	}

	// Fix missing default profile
	if _, err := c.profMgr.Get("default"); err != nil {
		c.profMgr.Add(config.ProfileConfig{Name: "default"})
		repairs = append(repairs, Result{Name: "Created Default Profile", Status: StatusOK, Message: "default"})
	}

	// Set active profile if empty or broken
	if globalCfg.CurrentProfile == "" {
		globalCfg.CurrentProfile = "default"
		c.cfgMgr.Save(globalCfg)
		repairs = append(repairs, Result{Name: "Set Active Profile", Status: StatusOK, Message: "default"})
	} else if _, err := c.profMgr.Get(globalCfg.CurrentProfile); err != nil {
		globalCfg.CurrentProfile = "default"
		c.cfgMgr.Save(globalCfg)
		repairs = append(repairs, Result{Name: "Repaired Active Profile Link", Status: StatusOK, Message: "default"})
	}

	return repairs
}
