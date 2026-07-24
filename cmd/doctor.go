package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose Relay installation",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("doctor called")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
