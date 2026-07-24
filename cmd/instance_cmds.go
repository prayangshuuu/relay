package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/spf13/cobra"
)

var instanceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new provider instance",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProviderManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		providerType := promptInteractive("Provider Template (e.g. agentrouter, openrouter)")
		if providerType == "" {
			fmt.Println("Error: Provider template is required.")
			os.Exit(1)
		}

		instanceName := promptInteractive("Instance Name")
		if instanceName == "" {
			fmt.Println("Error: Instance name is required.")
			os.Exit(1)
		}

		baseURL := promptInteractive("Base URL (optional)")
		apiKey := promptInteractive("API Key")

		cfg := config.ProviderConfig{
			ID:      instanceName,
			Name:    instanceName,
			Type:    providerType,
			BaseURL: baseURL,
			Enabled: true,
		}

		if apiKey != "" {
			cfg.EnvironmentVariables = []string{fmt.Sprintf("API_KEY=%s", apiKey)}
		}

		if err := manager.Add(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating instance: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Saved successfully.")
	},
}

var instanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List provider instances",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProviderManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		providers, err := manager.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing instances: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tTYPE\tSTATUS\tCREATED")
		for _, p := range providers {
			status := "Disabled"
			if p.Enabled {
				status = "Active"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Name, p.Type, status, p.CreatedAt)
		}
		w.Flush()
	},
}

var instanceRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a provider instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProviderManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := manager.Remove(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing instance: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Instance '%s' removed successfully.\n", args[0])
	},
}

func init() {
	instanceCmd.AddCommand(instanceCreateCmd)
	instanceCmd.AddCommand(instanceListCmd)
	instanceCmd.AddCommand(instanceRemoveCmd)
}
