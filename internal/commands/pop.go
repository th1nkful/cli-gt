package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var popCmd = &cobra.Command{
	Use:   "pop",
	Short: "Pop the current branch",
	Long:  `Pop the current branch (typically to remove and archive a completed branch).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("pop command - not yet implemented")
		return nil
	},
}
