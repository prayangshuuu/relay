package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/spf13/cobra"
)

func init() {
	providerCmd.AddCommand(providerListCmd)
	providerCmd.AddCommand(providerAddCmd)
	providerCmd.AddCommand(providerRemoveCmd)
	providerCmd.AddCommand(providerEditCmd)
	providerCmd.AddCommand(providerShowCmd)
	providerCmd.AddCommand(providerEnableCmd)
	providerCmd.AddCommand(providerDisableCmd)
	providerCmd.AddCommand(providerValidateCmd)
	providerCmd.AddCommand(providerExportCmd)
	providerCmd.AddCommand(providerImportCmd)
}

var providerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured providers",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProviderManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		providers, err := manager.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing providers: %v\n", err)
			os.Exit(1)
		}
		if len(providers) == 0 {
			fmt.Println("No providers configured.")
			return
		}
		fmt.Printf("%-20s %-20s %-20s %s\n", "ID", "NAME", "TYPE", "STATUS")
		fmt.Println(strings.Repeat("-", 70))
		for _, p := range providers {
			status := "Enabled"
			if !p.Enabled {
				status = "Disabled"
			}
			fmt.Printf("%-20s %-20s %-20s %s\n", p.ID, p.Name, p.Type, status)
		}
	},
}

var providerAddCmd = &cobra.Command{
	Use:   "add [id]",
	Short: "Add a new AI provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		manager, err := newProviderManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		templates := provider.Templates()
		var cfg config.ProviderConfig

		if template, exists := templates[id]; exists {
			fmt.Printf("Using built-in template for %s...\n", template.Name)
			cfg = template
		} else {
			cfg.ID = id
			fmt.Println("Configuring custom provider...")
			cfg.Name = promptInteractive("Provider name")
			cfg.Type = promptInteractive("Type (e.g., openai-compatible)")
			cfg.BaseURL = promptInteractive("Base URL")
			cfg.AuthenticationType = promptInteractive("Authentication Type (e.g., bearer, header, none)")

			envVars := promptInteractive("Environment Variables (comma-separated, e.g., API_KEY)")
			if envVars != "" {
				vars := strings.Split(envVars, ",")
				for i, v := range vars {
					vars[i] = strings.TrimSpace(v)
				}
				cfg.EnvironmentVariables = vars
			}
			cfg.Enabled = true
		}

		if err := manager.Add(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding provider: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Provider '%s' added successfully.\n", id)
	},
}

var providerRemoveCmd = &cobra.Command{
	Use:   "remove [id]",
	Short: "Remove a provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProviderManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := manager.Remove(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing provider: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Provider '%s' removed successfully.\n", args[0])
	},
}

var providerEditCmd = &cobra.Command{
	Use:   "edit [id]",
	Short: "Edit a provider interactively",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProviderManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		id := args[0]
		cfg, err := manager.Get(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading provider: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Editing provider '%s' (leave blank to keep current value)\n", id)

		if name := promptInteractive(fmt.Sprintf("Provider name [%s]", cfg.Name)); name != "" {
			cfg.Name = name
		}
		if t := promptInteractive(fmt.Sprintf("Type [%s]", cfg.Type)); t != "" {
			cfg.Type = t
		}
		if url := promptInteractive(fmt.Sprintf("Base URL [%s]", cfg.BaseURL)); url != "" {
			cfg.BaseURL = url
		}
		if auth := promptInteractive(fmt.Sprintf("Authentication Type [%s]", cfg.AuthenticationType)); auth != "" {
			cfg.AuthenticationType = auth
		}
		currentEnv := strings.Join(cfg.EnvironmentVariables, ", ")
		if envVars := promptInteractive(fmt.Sprintf("Environment Variables [%s]", currentEnv)); envVars != "" {
			vars := strings.Split(envVars, ",")
			for i, v := range vars {
				vars[i] = strings.TrimSpace(v)
			}
			cfg.EnvironmentVariables = vars
		}

		if err := manager.Edit(id, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving edits: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Provider '%s' updated successfully.\n", id)
	},
}

var providerShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show details of a provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProviderManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		cfg, err := manager.Get(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading provider: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("ID: %s\n", cfg.ID)
		fmt.Printf("Name: %s\n", cfg.Name)
		fmt.Printf("Type: %s\n", cfg.Type)
		fmt.Printf("Base URL: %s\n", cfg.BaseURL)
		fmt.Printf("Auth Type: %s\n", cfg.AuthenticationType)
		fmt.Printf("Env Vars: %v\n", cfg.EnvironmentVariables)
		fmt.Printf("Enabled: %v\n", cfg.Enabled)
		fmt.Printf("Created At: %s\n", cfg.CreatedAt)
		fmt.Printf("Updated At: %s\n", cfg.UpdatedAt)
	},
}

var providerEnableCmd = &cobra.Command{
	Use:   "enable [id]",
	Short: "Enable a provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProviderManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := manager.Enable(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling provider: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Provider '%s' enabled.\n", args[0])
	},
}

var providerDisableCmd = &cobra.Command{
	Use:   "disable [id]",
	Short: "Disable a provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProviderManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := manager.Disable(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error disabling provider: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Provider '%s' disabled.\n", args[0])
	},
}

var providerValidateCmd = &cobra.Command{
	Use:   "validate [id]",
	Short: "Validate a provider configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProviderManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		cfg, err := manager.Get(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading provider: %v\n", err)
			os.Exit(1)
		}
		if err := provider.Validate(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Validation failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Provider configuration is valid.")
	},
}

var providerExportCmd = &cobra.Command{
	Use:   "export [id] [dest_path]",
	Short: "Export a provider to a file (json or yaml)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProviderManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		format := "yaml"
		if strings.HasSuffix(strings.ToLower(args[1]), ".json") {
			format = "json"
		}

		if err := manager.Export(args[0], args[1], format); err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting provider: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Provider '%s' exported successfully to %s.\n", args[0], args[1])
	},
}

var providerImportCmd = &cobra.Command{
	Use:   "import [src_path]",
	Short: "Import a provider from a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProviderManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := manager.Import(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error importing provider: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Provider imported successfully.")
	},
}
