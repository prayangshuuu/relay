package cmd

import (
	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/profile"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/prayangshuuu/relay/internal/tool"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage Relay execution profiles",
}

func init() {
	rootCmd.AddCommand(profileCmd)
}

func newProfileManager() (*profile.Manager, error) {
	paths, err := config.NewOSPathManager()
	if err != nil {
		return nil, err
	}
	store := storage.NewLocalStorage()
	provMgr := provider.NewManager(paths, store)
	toolMgr := tool.NewManager(paths, store)
	return profile.NewManager(paths, store, provMgr, toolMgr), nil
}
