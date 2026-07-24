package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/plugin"
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage Relay extensions and plugins",
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed plugins",
	Run: func(cmd *cobra.Command, args []string) {
		paths, _ := config.NewOSPathManager()
		registry := plugin.NewRegistry(paths)
		plugin.RegisterBuiltIns(registry)
		registry.Discover()
		plugins := registry.List()

		if len(plugins) == 0 {
			fmt.Println("No plugins installed.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tVERSION\tTYPE")
		for _, p := range plugins {
			m := p.Manifest()
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", m.ID, m.Name, m.Version, m.Type)
		}
		w.Flush()
	},
}

var pluginInfoCmd = &cobra.Command{
	Use:   "info [id]",
	Short: "Show details for a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paths, _ := config.NewOSPathManager()
		registry := plugin.NewRegistry(paths)
		plugin.RegisterBuiltIns(registry)
		registry.Discover()

		p, err := registry.Get(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		m := p.Manifest()
		fmt.Printf("ID: %s\n", m.ID)
		fmt.Printf("Name: %s\n", m.Name)
		fmt.Printf("Version: %s\n", m.Version)
		fmt.Printf("Type: %s\n", m.Type)
		fmt.Printf("Author: %s\n", m.Author)
		fmt.Printf("Description: %s\n", m.Description)
	},
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install [path/url]",
	Short: "Install a new plugin (coming soon)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Plugin installation via CLI is coming soon.")
		fmt.Println("For now, copy the plugin directory into ~/.relay/plugins/")
	},
}

var pluginValidateCmd = &cobra.Command{
	Use:   "validate [id]",
	Short: "Validate a plugin's RPC connectivity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paths, _ := config.NewOSPathManager()
		registry := plugin.NewRegistry(paths)
		plugin.RegisterBuiltIns(registry)
		registry.Discover()

		_, err := registry.Get(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Plugin '%s' found. Validation checks passed (mock).\n", args[0])
	},
}

func init() {
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginInfoCmd)
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginCmd.AddCommand(pluginValidateCmd)
	rootCmd.AddCommand(pluginCmd)
}
