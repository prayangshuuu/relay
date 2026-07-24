package cmd

import (
	"fmt"
	"os"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/env"
	"github.com/prayangshuuu/relay/internal/launcher"
	"github.com/prayangshuuu/relay/internal/profile"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/shell"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/prayangshuuu/relay/internal/tool"
	"github.com/spf13/cobra"
)

var (
	runVerbose bool
	runDryRun  bool
	runProfile string
)

var runCmd = &cobra.Command{
	Use:   "run [tool] [args...]",
	Short: "Run an AI tool with the configured profile",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: must specify a tool to run")
			os.Exit(1)
		}

		toolName := args[0]
		toolArgs := args[1:]

		paths, err := config.NewOSPathManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		store := storage.NewLocalStorage()
		provMgr := provider.NewManager(paths, store)
		toolMgr := tool.NewManager(paths, store)
		profMgr := profile.NewManager(paths, store, provMgr, toolMgr)

		// Determine active profile
		var activeProfileName string
		if runProfile != "" {
			activeProfileName = runProfile
		} else {
			cfgMgr := config.NewManager(paths, store)
			cfg, err := cfgMgr.Load()
			if err == nil && cfg.CurrentProfile != "" {
				activeProfileName = cfg.CurrentProfile
			} else {
				activeProfileName = "default"
			}
		}

		prof, err := profMgr.Get(activeProfileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading profile: %v\n", err)
			os.Exit(1)
		}

		// If the user specified a tool different from the profile's tool, we use that tool.
		activeTool, err := toolMgr.Get(toolName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading tool: %v\n", err)
			os.Exit(1)
		}

		engine := launcher.NewEngine(env.NewOSEnvironment(), shell.NewOSExecShell())
		lCfg := launcher.Config{
			Profile: prof,
			Tool:    activeTool,
			Verbose: runVerbose,
			DryRun:  runDryRun,
			Args:    toolArgs,
		}

		if err := engine.Launch(lCfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	runCmd.Flags().BoolVarP(&runVerbose, "verbose", "v", false, "Enable verbose output")
	runCmd.Flags().BoolVar(&runDryRun, "dry-run", false, "Simulate execution without launching")
	runCmd.Flags().StringVarP(&runProfile, "profile", "p", "", "Override the active profile")
	rootCmd.AddCommand(runCmd)
}
