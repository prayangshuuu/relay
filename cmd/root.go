package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "relay",
	Short: "Universal AI provider manager",
	Long: `Relay is a lightweight, cross-platform utility for managing AI providers, profiles, environments, and tool integrations.

Configure once.
Switch instantly.

Relay allows developers to switch between providers like OpenRouter, AgentRouter, Anthropic, OpenAI, Ollama and more without manually editing environment variables.`,
	Version: "0.1.0",
	// Handle dynamic aliases like `relay claude`
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			// Proxy to the runCmd logic
			runCmd.Run(runCmd, args)
			return
		}
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
