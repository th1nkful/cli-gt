package commands

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var continueCmd = &cobra.Command{
	Use:   "continue",
	Short: "Continue a rebase in progress",
	Long:  `Continue a rebase in progress. This is equivalent to 'git rebase --continue' and only works if a rebase is currently in progress.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if a rebase is in progress
		inProgress, err := isRebaseInProgress()
		if err != nil {
			return err
		}

		if !inProgress {
			return fmt.Errorf("Error: no rebase in progress")
		}

		// Execute git rebase --continue
		continueCmd := exec.Command("git", "rebase", "--continue")
		output, err := continueCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to continue rebase: %w\nOutput: %s", err, string(output))
		}

		fmt.Println("✔ gt continue: rebase continued successfully")
		return nil
	},
}
