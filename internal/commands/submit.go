package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit the current branch for review",
	Long:  `Submit the current branch for review (typically by creating or updating a pull request).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("submit command - not yet implemented")
		return nil
	},
}
