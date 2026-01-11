package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gt",
	Short: "A Git workflow CLI tool",
	Long: `gt is a CLI tool that augments git with opinionated workflow commands.
It helps manage branches, track settings per workspace, and streamline common git operations.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add all subcommands
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(popCmd)
	rootCmd.AddCommand(modifyCmd)
	rootCmd.AddCommand(checkoutCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(restackCmd)
	rootCmd.AddCommand(submitCmd)
	rootCmd.AddCommand(continueCmd)
}
