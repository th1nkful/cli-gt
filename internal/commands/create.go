package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [branch-name]",
	Short: "Create a new branch and commit",
	Long:  `Create a new branch and commit. This creates a new branch from the current state and makes an initial commit.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("create command - not yet implemented")
		if len(args) > 0 {
			fmt.Printf("Branch name: %s\n", args[0])
		}
		return nil
	},
}
