package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var modifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "Modify branch settings",
	Long:  `Modify settings for the current or specified branch.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("modify command - not yet implemented")
		return nil
	},
}
