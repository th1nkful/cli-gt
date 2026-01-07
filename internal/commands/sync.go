package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Update trunk and rebase tracked branches",
	Long:  `Updates trunk branch from origin, rebases local tracked branches onto trunk again. If a local tracked branch no longer exists on origin, prompts for confirmation (y/n) to delete the branch.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("sync command - not yet implemented")
		return nil
	},
}
