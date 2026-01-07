package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync the current branch with its parent",
	Long:  `Sync the current branch with its parent branch (typically rebasing or merging).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("sync command - not yet implemented")
		return nil
	},
}
