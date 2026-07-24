package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/spf13/cobra"
)

func init() {
	profileCmd.AddCommand(profileCreateCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileShowCmd)
	profileCmd.AddCommand(profileEditCmd)
	profileCmd.AddCommand(profileRemoveCmd)
	profileCmd.AddCommand(profileCloneCmd)
	profileCmd.AddCommand(profileUseCmd)
	profileCmd.AddCommand(profileCurrentCmd)
	profileCmd.AddCommand(profileValidateCmd)
}

var profileCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		manager, err := newProfileManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Configuring profile '%s'...\n", name)
		var cfg config.ProfileConfig
		cfg.Name = name
		cfg.Tool = promptInteractive("Tool (e.g., claude, aider)")
		cfg.Provider = promptInteractive("Provider (e.g., openrouter, anthropic)")
		cfg.Model = promptInteractive("Model (e.g., claude-3-5-sonnet-20240620)")

		if err := manager.Add(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating profile: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Profile '%s' created successfully.\n", name)
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProfileManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		profiles, err := manager.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing profiles: %v\n", err)
			os.Exit(1)
		}

		if len(profiles) == 0 {
			fmt.Println("No profiles configured.")
			return
		}

		fmt.Printf("%-20s %-20s %-20s %-30s\n", "NAME", "TOOL", "PROVIDER", "MODEL")
		fmt.Println(strings.Repeat("-", 95))
		for _, p := range profiles {
			toolName := "none"
			if t := p.Tool(); t != nil {
				toolName = t.ExecutableName()
			}
			provName := "none"
			if pr := p.Provider(); pr != nil {
				provName = pr.Name()
			}
			fmt.Printf("%-20s %-20s %-20s %-30s\n", p.Name(), toolName, provName, p.Model())
		}
	},
}

var profileShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show details of a profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProfileManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		p, err := manager.Get(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		toolName := "none"
		if t := p.Tool(); t != nil {
			toolName = t.ExecutableName()
		}
		provName := "none"
		if pr := p.Provider(); pr != nil {
			provName = pr.Name()
		}

		fmt.Printf("Name: %s\n", p.Name())
		fmt.Printf("Tool: %s\n", toolName)
		fmt.Printf("Provider: %s\n", provName)
		fmt.Printf("Model: %s\n", p.Model())
		fmt.Printf("Environment Overrides: %v\n", p.Overrides())
	},
}

var profileEditCmd = &cobra.Command{
	Use:   "edit [name]",
	Short: "Edit a profile interactively",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		manager, err := newProfileManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		p, err := manager.Get(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Access the underlying config
		var cfg config.ProfileConfig
		if cp, ok := p.(interface{ Config() config.ProfileConfig }); ok {
			cfg = cp.Config()
		} else {
			fmt.Println("Error: unsupported profile type")
			os.Exit(1)
		}

		fmt.Printf("Editing profile '%s' (leave blank to keep current value)\n", name)

		if t := promptInteractive(fmt.Sprintf("Tool [%s]", cfg.Tool)); t != "" {
			cfg.Tool = t
		}
		if pr := promptInteractive(fmt.Sprintf("Provider [%s]", cfg.Provider)); pr != "" {
			cfg.Provider = pr
		}
		if m := promptInteractive(fmt.Sprintf("Model [%s]", cfg.Model)); m != "" {
			cfg.Model = m
		}

		if err := manager.Edit(name, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving edits: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Profile '%s' updated successfully.\n", name)
	},
}

var profileRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProfileManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := manager.Remove(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing profile: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Profile '%s' removed successfully.\n", args[0])
	},
}

var profileCloneCmd = &cobra.Command{
	Use:   "clone [src] [target]",
	Short: "Clone a profile",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProfileManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := manager.Clone(args[0], args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error cloning profile: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Profile '%s' cloned to '%s'.\n", args[0], args[1])
	},
}

var profileUseCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Set the global active profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		paths, err := config.NewOSPathManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		store := storage.NewLocalStorage()

		// Validate profile exists
		profMgr, _ := newProfileManager()
		if _, err := profMgr.Get(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error: profile '%s' does not exist.\n", name)
			os.Exit(1)
		}

		cfgMgr := config.NewManager(paths, store)
		cfg, err := cfgMgr.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading global config: %v\n", err)
			os.Exit(1)
		}

		cfg.CurrentProfile = name
		if err := cfgMgr.Save(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving global config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Active profile set to '%s'.\n", name)
	},
}

var profileCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the currently active profile",
	Run: func(cmd *cobra.Command, args []string) {
		paths, err := config.NewOSPathManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		store := storage.NewLocalStorage()
		cfgMgr := config.NewManager(paths, store)
		cfg, err := cfgMgr.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		if cfg.CurrentProfile == "" {
			fmt.Println("No active profile set. Using default.")
		} else {
			fmt.Println(cfg.CurrentProfile)
		}
	},
}

var profileValidateCmd = &cobra.Command{
	Use:   "validate [name]",
	Short: "Validate a profile's tool and provider dependencies",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newProfileManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := manager.Validate(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Validation failed:\n%v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Profile '%s' is valid.\n", args[0])
	},
}
