package cmd

import (
	"fmt"
	"os"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View recently used environments",
	Run: func(cmd *cobra.Command, args []string) {
		paths, err := config.NewOSPathManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		store := storage.NewLocalStorage()
		cfgMgr := config.NewManager(paths, store)
		globalCfg, _ := cfgMgr.Load()

		fmt.Println("=== Relay Usage History ===")
		printHistory("Profiles", globalCfg.RecentProfiles)
		printHistory("Instances", globalCfg.RecentProviders)
		printHistory("Tools", globalCfg.RecentTools)
		printHistory("Models", globalCfg.RecentModels)
	},
}

func printHistory(name string, items []string) {
	fmt.Printf("\nRecent %s:\n", name)
	if len(items) == 0 {
		fmt.Println("  (none)")
		return
	}
	for i, item := range items {
		if i == 0 {
			fmt.Printf("  -> %s (current)\n", item)
		} else {
			fmt.Printf("     %s\n", item)
		}
	}
}

var historyClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all history",
	Run: func(cmd *cobra.Command, args []string) {
		paths, _ := config.NewOSPathManager()
		store := storage.NewLocalStorage()
		cfgMgr := config.NewManager(paths, store)
		globalCfg, _ := cfgMgr.Load()

		globalCfg.RecentProfiles = nil
		globalCfg.RecentProviders = nil
		globalCfg.RecentTools = nil
		globalCfg.RecentModels = nil

		cfgMgr.Save(globalCfg)
		fmt.Println("History cleared.")
	},
}

func init() {
	historyCmd.AddCommand(historyClearCmd)
	rootCmd.AddCommand(historyCmd)
}
