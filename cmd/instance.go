package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/spf13/cobra"
)

var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Manage AI provider instances",
}

func init() {
	rootCmd.AddCommand(instanceCmd)
}

func newProviderManager() (*provider.Manager, error) {
	paths, err := config.NewOSPathManager()
	if err != nil {
		return nil, err
	}
	store := storage.NewLocalStorage()
	return provider.NewManager(paths, store), nil
}

// promptInteractive asks the user for a value if it's missing.
func promptInteractive(prompt string) string {
	fmt.Printf("%s: ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}
