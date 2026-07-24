package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/profile"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/setup"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/prayangshuuu/relay/internal/tool"
	"github.com/spf13/cobra"
)

var setupReset bool

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup wizard for Relay",
	Run: func(cmd *cobra.Command, args []string) {
		paths, err := config.NewOSPathManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if setupReset {
			fmt.Println("⚠ WARNING: This will completely delete your Relay configuration.")
			fmt.Print("Are you absolutely sure? (type 'yes'): ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			if strings.TrimSpace(input) == "yes" {
				os.RemoveAll(paths.BaseDir())
				fmt.Println("Configuration purged.")
			} else {
				fmt.Println("Reset cancelled.")
				os.Exit(0)
			}
		}

		// Ensure config dir exists
		if _, err := os.Stat(paths.BaseDir()); os.IsNotExist(err) {
			os.MkdirAll(paths.BaseDir(), 0755)
		}

		store := storage.NewLocalStorage()
		cfgMgr := config.NewManager(paths, store)

		// If no config exists, initialize to prevent silent failures
		if _, err := cfgMgr.Load(); err != nil {
			cfgMgr.Initialize()
		}

		provMgr := provider.NewManager(paths, store)
		toolMgr := tool.NewManager(paths, store)
		profMgr := profile.NewManager(paths, store, provMgr, toolMgr)

		wizard := setup.NewWizard(cfgMgr, profMgr, provMgr, toolMgr)
		wizard.Run()
	},
}

func init() {
	setupCmd.Flags().BoolVar(&setupReset, "reset", false, "Delete existing configuration and start fresh")
	rootCmd.AddCommand(setupCmd)
}
