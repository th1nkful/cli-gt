package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var checkoutCmd = &cobra.Command{
	Use:     "checkout [branch]",
	Aliases: []string{"co"},
	Short:   "Checkout a branch",
	Long:    `Checkout a branch managed by gt.`,
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("checkout command - not yet implemented")
		if len(args) > 0 {
			fmt.Printf("Branch: %s\n", args[0])
		}
		return nil
	},
}
