package commands

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var checkoutCmd = &cobra.Command{
	Use:     "checkout [branch]",
	Aliases: []string{"co"},
	Short:   "Checkout a branch",
	Long:    `Checkout to a branch. If no branch is supplied, list available branches with trunk branch at the bottom and most recently used above that, which you can navigate using up/down arrows to select from the list.`,
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// If branch is provided as argument, switch directly
		if len(args) > 0 {
			branchName := args[0]
			if err := switchBranch(branchName); err != nil {
				return err
			}
			fmt.Printf("Switched to branch '%s'\n", branchName)
			return nil
		}

		// No branch provided - show interactive selection
		branches, err := getBranches()
		if err != nil {
			return err
		}

		// Use survey to prompt for branch selection
		var selectedBranch string
		prompt := &survey.Select{
			Message: "checkout>",
			Options: branches,
		}

		if err := survey.AskOne(prompt, &selectedBranch); err != nil {
			return fmt.Errorf("branch selection cancelled or failed: %w", err)
		}

		// Switch to the selected branch
		if err := switchBranch(selectedBranch); err != nil {
			return err
		}

		fmt.Printf("Switched to branch '%s'\n", selectedBranch)
		return nil
	},
}
