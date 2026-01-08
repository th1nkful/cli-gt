package commands

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	modifyAll bool
)

var modifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "Amend the current commit",
	Long:  `Amend the current commit. This allows you to modify the most recent commit on the current branch. Will not run on trunk branch.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get current branch
		currentBranch, err := getCurrentBranch()
		if err != nil {
			return err
		}

		// Check for detached HEAD
		if currentBranch == "HEAD" {
			return fmt.Errorf("Error: detached HEAD")
		}

		// Check if we're on trunk branch
		onTrunk, cfg, err := isOnTrunkBranch()
		if err != nil {
			return err
		}
		if onTrunk {
			return fmt.Errorf("Error: gt modify cannot be run on %s", cfg.TrunkBranch)
		}

		// Stage all files if -a flag is used
		if modifyAll {
			if err := stageAllFiles(); err != nil {
				return err
			}
		}

		// Check if branch exists on origin
		_, remoteExists, err := branchExists(currentBranch)
		if err != nil {
			// If we can't check remote, just continue (might not have origin configured)
			// Don't fail the command because of this
		} else if remoteExists {
			fmt.Printf("⚠️  Branch '%s' exists on origin.\n", currentBranch)
			fmt.Println("    Amending rewrites history; you'll likely need:")
			fmt.Println("    git push --force-with-lease")
			fmt.Println()
		}

		// Amend the commit
		amendCmd := exec.Command("git", "commit", "--amend", "--no-edit")
		output, err := amendCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to amend commit: %w\nOutput: %s", err, string(output))
		}

		fmt.Println("✔ gt modify: amended latest commit")
		return nil
	},
}

func init() {
	modifyCmd.Flags().BoolVarP(&modifyAll, "all", "a", false, "Stage all changes before amending")
}
