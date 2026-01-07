package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var restackCmd = &cobra.Command{
	Use:   "restack",
	Short: "Restack branches",
	Long:  `Restack all managed branches to ensure they are up to date with their parents.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("restack command - not yet implemented")
		return nil
	},
}
