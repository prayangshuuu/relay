package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Manage AI providers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("provider called")
	},
}

func init() {
	rootCmd.AddCommand(providerCmd)
}
