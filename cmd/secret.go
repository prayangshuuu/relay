package cmd

import (
	"fmt"
	"os"

	"github.com/prayangshuuu/relay/internal/keyring"
	"github.com/spf13/cobra"
)

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage instance secrets natively in OS Keyring",
}

var secretSetCmd = &cobra.Command{
	Use:   "set [instance_id]",
	Short: "Set an API key in the OS Keyring",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := promptInteractive("Enter API Key (will not be echoed securely for now)")
		if apiKey == "" {
			fmt.Println("Error: API Key cannot be empty.")
			os.Exit(1)
		}

		km := keyring.NewManager()
		if err := km.Set(args[0], apiKey); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving to OS keyring: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Secret successfully saved to OS Keyring for instance '%s'.\n", args[0])
	},
}

var secretGetCmd = &cobra.Command{
	Use:   "get [instance_id]",
	Short: "Retrieve an API key from the OS Keyring",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		km := keyring.NewManager()
		val, err := km.Get(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching from OS keyring (or not found): %v\n", err)
			os.Exit(1)
		}

		fmt.Println(val)
	},
}

var secretDeleteCmd = &cobra.Command{
	Use:   "delete [instance_id]",
	Short: "Delete an API key from the OS Keyring",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		km := keyring.NewManager()
		if err := km.Delete(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting from OS keyring: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Secret successfully deleted from OS Keyring for instance '%s'.\n", args[0])
	},
}

func init() {
	secretCmd.AddCommand(secretSetCmd)
	secretCmd.AddCommand(secretGetCmd)
	secretCmd.AddCommand(secretDeleteCmd)
	rootCmd.AddCommand(secretCmd)
}
