package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var checkoutCmd = &cobra.Command{
	Use:     "checkout [branch]",
	Aliases: []string{"co"},
	Short:   "Checkout a branch",
	Long:    `Checkout to a branch. If no branch is supplied, list available branches with trunk branch at the bottom and most recently used above that, which you can navigate using up/down arrows to select from the list.`,
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("checkout command - not yet implemented")
		if len(args) > 0 {
			fmt.Printf("Branch: %s\n", args[0])
		}
		return nil
	},
}
