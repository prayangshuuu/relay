package cmd

import (
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:               "switch [args...]",
	Short:             "Alias for 'relay use'",
	Args:              useCmd.Args,
	Run:               useCmd.Run,
	ValidArgsFunction: useCmd.ValidArgsFunction,
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
