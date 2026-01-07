package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var restackCmd = &cobra.Command{
	Use:   "restack",
	Short: "Restack all managed branches",
	Long:  `Restack all managed branches to ensure they are up to date with their parent branch (assuming always trunk for now).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("restack command - not yet implemented")
		return nil
	},
}
