package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage Relay profiles",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("profile called")
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
}
