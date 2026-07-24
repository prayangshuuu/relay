package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/profile"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/prayangshuuu/relay/internal/tool"
	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Restore the configuration to its previous state",
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

		if len(globalCfg.UndoStack) == 0 {
			fmt.Println("No previous state to undo.")
			return
		}

		// Pop last state
		lastStateStr := globalCfg.UndoStack[len(globalCfg.UndoStack)-1]
		newStack := globalCfg.UndoStack[:len(globalCfg.UndoStack)-1]

		var state UndoState
		if err := json.Unmarshal([]byte(lastStateStr), &state); err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding undo state: %v\n", err)
			os.Exit(1)
		}

		// Restore profiles
		provMgr := provider.NewManager(paths, store)
		toolMgr := tool.NewManager(paths, store)
		profMgr := profile.NewManager(paths, store, provMgr, toolMgr)

		for name, pCfg := range state.Profiles {
			if err := profMgr.Edit(name, pCfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error restoring profile '%s': %v\n", name, err)
			}
		}

		// Restore global config (retaining the now-shrunk undo stack)
		restoredCfg := state.Global
		restoredCfg.UndoStack = newStack

		if err := cfgMgr.Save(&restoredCfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving restored config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Undo successful. Restored previous configuration.")
	},
}

func init() {
	rootCmd.AddCommand(undoCmd)
}
