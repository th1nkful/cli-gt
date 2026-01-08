package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var popCmd = &cobra.Command{
	Use:   "pop",
	Short: "Undo the current branch and commit",
	Long:  `Undo the current branch and commit, returning the files from the commit/branch to an uncommitted state. This effectively undoes the 'create' command. Will not run on trunk branch.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("pop command - not yet implemented")
		return nil
	},
}
