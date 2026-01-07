package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit the current branch for review",
	Long:  `Submit the current branch for review (e.g., create/update a pull request). Will not run on trunk branch.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("submit command - not yet implemented")
		return nil
	},
}
