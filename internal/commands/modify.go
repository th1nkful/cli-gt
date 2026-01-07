package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var modifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "Amend the current commit",
	Long:  `Amend the current commit. This allows you to modify the most recent commit on the current branch. Will not run on trunk branch.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("modify command - not yet implemented")
		return nil
	},
}
