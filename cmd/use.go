package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/profile"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/prayangshuuu/relay/internal/switcher"
	"github.com/prayangshuuu/relay/internal/tool"
	"github.com/spf13/cobra"
)

// UndoState tracks the previous state for a single 'use' invocation
type UndoState struct {
	Global   config.GlobalConfig             `json:"global"`
	Profiles map[string]config.ProfileConfig `json:"profiles"`
}

var useCmd = &cobra.Command{
	Use:   "use [args...]",
	Short: "Smart switch context (profiles, providers, tools)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paths, err := config.NewOSPathManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		store := storage.NewLocalStorage()
		cfgMgr := config.NewManager(paths, store)

		globalCfg, err := cfgMgr.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		provMgr := provider.NewManager(paths, store)
		toolMgr := tool.NewManager(paths, store)
		profMgr := profile.NewManager(paths, store, provMgr, toolMgr)
		resolver := switcher.NewResolver(paths, globalCfg.Aliases)

		// Record initial state for undo
		undoState := UndoState{
			Global:   *globalCfg,
			Profiles: make(map[string]config.ProfileConfig),
		}

		// Ensure we don't infinitely recurse the undo stack in the snapshot
		undoState.Global.UndoStack = nil

		mutatedProfiles := make(map[string]config.ProfileConfig)

		for _, arg := range args {
			match, err := resolver.Resolve(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			switch match.Type {
			case switcher.TypeProfile:
				fmt.Printf("Switched profile to '%s'\n", match.Name)
				globalCfg.CurrentProfile = match.Name
				globalCfg.RecentProfiles = prependUnique(globalCfg.RecentProfiles, match.Name, 10)

			case switcher.TypeProvider:
				activeProf := globalCfg.CurrentProfile
				if activeProf == "" {
					activeProf = "default"
				}

				// Fetch current profile to mutate
				p, err := profMgr.Get(activeProf)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error fetching active profile '%s': %v\n", activeProf, err)
					os.Exit(1)
				}

				var pCfg config.ProfileConfig
				if cp, ok := p.(interface{ Config() config.ProfileConfig }); ok {
					pCfg = cp.Config()
				} else {
					fmt.Fprintln(os.Stderr, "Error: unsupported profile type")
					os.Exit(1)
				}

				// Snapshot original for undo if not already snapshotted
				if _, exists := undoState.Profiles[activeProf]; !exists {
					undoState.Profiles[activeProf] = pCfg
				}

				fmt.Printf("Switched provider to '%s' (in profile '%s')\n", match.Name, activeProf)
				pCfg.Provider = match.Name
				mutatedProfiles[activeProf] = pCfg
				globalCfg.RecentProviders = prependUnique(globalCfg.RecentProviders, match.Name, 10)

			case switcher.TypeTool:
				activeProf := globalCfg.CurrentProfile
				if activeProf == "" {
					activeProf = "default"
				}

				p, err := profMgr.Get(activeProf)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error fetching active profile '%s': %v\n", activeProf, err)
					os.Exit(1)
				}

				var pCfg config.ProfileConfig
				if cp, ok := p.(interface{ Config() config.ProfileConfig }); ok {
					pCfg = cp.Config()
				} else {
					fmt.Fprintln(os.Stderr, "Error: unsupported profile type")
					os.Exit(1)
				}

				if _, exists := undoState.Profiles[activeProf]; !exists {
					undoState.Profiles[activeProf] = pCfg
				}

				fmt.Printf("Switched tool to '%s' (in profile '%s')\n", match.Name, activeProf)
				pCfg.Tool = match.Name
				mutatedProfiles[activeProf] = pCfg
				globalCfg.RecentTools = prependUnique(globalCfg.RecentTools, match.Name, 10)

			case switcher.TypeModel:
				activeProf := globalCfg.CurrentProfile
				if activeProf == "" {
					activeProf = "default"
				}

				p, err := profMgr.Get(activeProf)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error fetching active profile '%s': %v\n", activeProf, err)
					os.Exit(1)
				}

				var pCfg config.ProfileConfig
				if cp, ok := p.(interface{ Config() config.ProfileConfig }); ok {
					pCfg = cp.Config()
				} else {
					fmt.Fprintln(os.Stderr, "Error: unsupported profile type")
					os.Exit(1)
				}

				if _, exists := undoState.Profiles[activeProf]; !exists {
					undoState.Profiles[activeProf] = pCfg
				}

				fmt.Printf("Switched model to '%s' (in profile '%s')\n", match.Name, activeProf)
				pCfg.Model = match.Name
				mutatedProfiles[activeProf] = pCfg
				globalCfg.RecentModels = prependUnique(globalCfg.RecentModels, match.Name, 10)
			}
		}

		// Save profile mutations
		for name, pCfg := range mutatedProfiles {
			if err := profMgr.Edit(name, pCfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating profile '%s': %v\n", name, err)
				os.Exit(1)
			}
		}

		// Serialize undo state and push to stack
		stateBytes, _ := json.Marshal(undoState)
		globalCfg.UndoStack = append(globalCfg.UndoStack, string(stateBytes))

		// Cap undo stack at 50 to prevent unbounded growth
		if len(globalCfg.UndoStack) > 50 {
			globalCfg.UndoStack = globalCfg.UndoStack[len(globalCfg.UndoStack)-50:]
		}

		if err := cfgMgr.Save(globalCfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
	},
}

func prependUnique(slice []string, item string, max int) []string {
	var result []string
	result = append(result, item)
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	if len(result) > max {
		return result[:max]
	}
	return result
}

func init() {
	useCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		paths, err := config.NewOSPathManager()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var completions []string

		if entries, _ := os.ReadDir(paths.ProfilesDir()); err == nil {
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".yaml") {
					completions = append(completions, strings.TrimSuffix(e.Name(), ".yaml"))
				}
			}
		}
		if entries, _ := os.ReadDir(paths.ProvidersDir()); err == nil {
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".yaml") {
					completions = append(completions, strings.TrimSuffix(e.Name(), ".yaml"))
				}
			}
		}
		if entries, _ := os.ReadDir(paths.ToolsDir()); err == nil {
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".yaml") {
					completions = append(completions, strings.TrimSuffix(e.Name(), ".yaml"))
				}
			}
		}

		return completions, cobra.ShellCompDirectiveNoFileComp
	}
	rootCmd.AddCommand(useCmd)
}
