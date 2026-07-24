package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Display current active configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("current called")
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
