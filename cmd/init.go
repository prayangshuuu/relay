package cmd

import (
	"fmt"
	"os"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/profile"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/prayangshuuu/relay/internal/tool"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Relay configuration",
	Run: func(cmd *cobra.Command, args []string) {
		paths, err := config.NewOSPathManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving config paths: %v\n", err)
			os.Exit(1)
		}

		store := storage.NewLocalStorage()
		manager := config.NewManager(paths, store)

		fmt.Printf("Initializing Relay in %s...\n", paths.BaseDir())
		if err := manager.Initialize(); err != nil {
			if err == config.ErrAlreadyInitialized {
				fmt.Println("Relay is already initialized. Skipping to avoid overwriting existing files.")
				return
			}
			fmt.Fprintf(os.Stderr, "Error initializing Relay: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Relay initialized successfully.")

		provMgr := provider.NewManager(paths, store)
		toolMgr := tool.NewManager(paths, store)
		profMgr := profile.NewManager(paths, store, provMgr, toolMgr)

		// Ensure default profile exists
		if _, err := profMgr.Get("default"); err != nil {
			err = profMgr.Add(config.ProfileConfig{
				Name: "default",
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to create default profile: %v\n", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
