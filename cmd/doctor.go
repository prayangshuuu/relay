package cmd

import (
	"fmt"
	"os"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/doctor"
	"github.com/prayangshuuu/relay/internal/profile"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/storage"
	"github.com/prayangshuuu/relay/internal/tool"
	"github.com/spf13/cobra"
)

var doctorFix bool
var doctorVerbose bool

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Validate and diagnose Relay configuration",
	Run: func(cmd *cobra.Command, args []string) {
		paths, err := config.NewOSPathManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		store := storage.NewLocalStorage()
		cfgMgr := config.NewManager(paths, store)
		provMgr := provider.NewManager(paths, store)
		toolMgr := tool.NewManager(paths, store)
		profMgr := profile.NewManager(paths, store, provMgr, toolMgr)

		suite := doctor.NewCheckSuite(paths, store, cfgMgr, profMgr, provMgr, toolMgr)

		if doctorFix {
			fmt.Println("Running Auto-Repair...")
			repairs := suite.Repair()
			for _, r := range repairs {
				fmt.Printf("  ✓ %s: %s\n", r.Name, r.Message)
			}
			fmt.Println("Repair complete. Running checks...")
			fmt.Println()
		}

		results := suite.RunChecks()

		var warnings, errs int

		fmt.Println("Relay Doctor:")
		for _, r := range results {
			switch r.Status {
			case doctor.StatusOK:
				if doctorVerbose {
					fmt.Printf("  ✓ %-30s : %s\n", r.Name, r.Message)
				} else {
					fmt.Printf("  ✓ %-30s : %s\n", r.Name, r.Message)
				}
			case doctor.StatusWarning:
				fmt.Printf("  ⚠ %-30s : %s\n", r.Name, r.Message)
				warnings++
			case doctor.StatusError:
				fmt.Printf("  ✗ %-30s : %s\n", r.Name, r.Message)
				errs++
			}
		}

		if errs > 0 || warnings > 0 {
			fmt.Printf("\nFound %d errors and %d warnings.\n", errs, warnings)
			if !doctorFix && errs > 0 {
				fmt.Println("Run 'relay doctor --fix' to attempt automatic repair.")
			}
			if errs > 0 {
				os.Exit(1)
			}
		} else {
			fmt.Println("\nConfiguration is healthy!")
		}
	},
}

func init() {
	doctorCmd.Flags().BoolVar(&doctorFix, "fix", false, "Attempt to automatically fix issues")
	doctorCmd.Flags().BoolVar(&doctorFix, "repair", false, "Alias for --fix")
	doctorCmd.Flags().BoolVarP(&doctorVerbose, "verbose", "v", false, "Show detailed output")
	rootCmd.AddCommand(doctorCmd)
}
