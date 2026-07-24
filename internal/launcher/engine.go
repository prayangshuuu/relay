package launcher

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/env"
	"github.com/prayangshuuu/relay/internal/keyring"
	"github.com/prayangshuuu/relay/internal/profile"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/shell"
	"github.com/prayangshuuu/relay/internal/tool"
)

// Engine implements Launcher interface
type Engine struct {
	env   env.Environment
	shell shell.Shell
}

func NewEngine(e env.Environment, s shell.Shell) *Engine {
	return &Engine{env: e, shell: s}
}

// Config wraps runtime launch variables
type Config struct {
	Profile  profile.Profile
	Tool     tool.Tool         // If overriden from CLI (e.g. relay run aider)
	Provider provider.Provider // If explicitly requested
	Verbose  bool
	DryRun   bool
	Args     []string
}

// Launch prepares and runs the tool.
func (e *Engine) Launch(cfg Config) error {
	activeTool := cfg.Tool
	if activeTool == nil && cfg.Profile != nil {
		activeTool = cfg.Profile.Tool()
	}

	if activeTool == nil {
		return fmt.Errorf("no tool specified to launch")
	}

	executable := activeTool.ExecutableName()
	resolvedExe, err := exec.LookPath(executable)
	if err != nil {
		return fmt.Errorf("executable '%s' not found in PATH: %w", executable, err)
	}

	// Build injections map
	injections := make(map[string]string)

	activeProvider := cfg.Provider
	if activeProvider == nil && cfg.Profile != nil {
		activeProvider = cfg.Profile.Provider()
	}

	var provName string
	if activeProvider != nil {
		provName = activeProvider.Name()
		if defaultProv, ok := activeProvider.(interface{ Config() config.ProviderConfig }); ok {
			pCfg := defaultProv.Config()

			for _, envLine := range pCfg.EnvironmentVariables {
				parts := strings.SplitN(envLine, "=", 2)
				if len(parts) == 2 {
					injections[parts[0]] = parts[1]
				}
			}

			if pCfg.UsesKeyring {
				km := keyring.NewManager()
				secret, err := km.Get(pCfg.ID)
				if err == nil && secret != "" {
					injections["API_KEY"] = secret
					injections[strings.ToUpper(pCfg.Type)+"_API_KEY"] = secret
				}
			}
		}
	}

	// Inject Profile overrides
	if cfg.Profile != nil {
		for k, v := range cfg.Profile.Overrides() {
			injections[k] = v
		}
		if cfg.Profile.Model() != "" {
			// Some tools respect specific env vars for models
			injections["RELAY_MODEL"] = cfg.Profile.Model()
		}
	}

	finalEnv := e.env.Build(injections)

	if cfg.Verbose || cfg.DryRun {
		profName := "none"
		if cfg.Profile != nil {
			profName = cfg.Profile.Name()
		}

		fmt.Printf("Loaded profile: %s\n", profName)
		fmt.Printf("Loaded provider: %s\n", provName)
		fmt.Printf("Injected variables:\n")
		for k := range injections {
			fmt.Printf("  %s\n", k)
		}

		if cfg.DryRun {
			fmt.Println("--- DRY RUN ---")
			fmt.Printf("Executable: %s\n", resolvedExe)
			fmt.Printf("Arguments: %v\n", cfg.Args)
			return nil
		}

		fmt.Printf("Launching %s...\n", executable)
	}

	// Delegate to shell
	return e.shell.Exec(resolvedExe, cfg.Args, finalEnv)
}
