package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage Relay aliases",
}

var aliasAddCmd = &cobra.Command{
	Use:   "add [alias] [target]",
	Short: "Add a new alias",
	Args:  cobra.ExactArgs(2),
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

		if globalCfg.Aliases == nil {
			globalCfg.Aliases = make(map[string]string)
		}

		globalCfg.Aliases[args[0]] = args[1]
		if err := cfgMgr.Save(globalCfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Alias '%s' -> '%s' added successfully.\n", args[0], args[1])
	},
}

var aliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases",
	Run: func(cmd *cobra.Command, args []string) {
		paths, _ := config.NewOSPathManager()
		store := storage.NewLocalStorage()
		cfgMgr := config.NewManager(paths, store)
		globalCfg, _ := cfgMgr.Load()

		if len(globalCfg.Aliases) == 0 {
			fmt.Println("No aliases configured.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ALIAS\tTARGET")
		for k, v := range globalCfg.Aliases {
			fmt.Fprintf(w, "%s\t%s\n", k, v)
		}
		w.Flush()
	},
}

var aliasRemoveCmd = &cobra.Command{
	Use:   "remove [alias]",
	Short: "Remove an alias",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paths, _ := config.NewOSPathManager()
		store := storage.NewLocalStorage()
		cfgMgr := config.NewManager(paths, store)
		globalCfg, _ := cfgMgr.Load()

		if globalCfg.Aliases != nil {
			delete(globalCfg.Aliases, args[0])
			cfgMgr.Save(globalCfg)
		}
		fmt.Printf("Alias '%s' removed.\n", args[0])
	},
}

func init() {
	aliasCmd.AddCommand(aliasAddCmd)
	aliasCmd.AddCommand(aliasListCmd)
	aliasCmd.AddCommand(aliasRemoveCmd)
	rootCmd.AddCommand(aliasCmd)
}
