package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var popCmd = &cobra.Command{
	Use:   "pop",
	Short: "Undo the current branch and commit",
	Long:  `Undo the current branch and commit, returning the files from the commit/branch to an uncommitted state. This effectively undoes the 'create' command. Will not run on trunk branch.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if we're on trunk branch and load config
		onTrunk, cfg, err := isOnTrunkBranch()
		if err != nil {
			return err
		}
		if onTrunk {
			return fmt.Errorf("pop command cannot be run on trunk branch (%s)", cfg.TrunkBranch)
		}

		// Get current branch name
		currentBranch, err := getCurrentBranch()
		if err != nil {
			return err
		}

		// Get parent branch from config (if managed), otherwise default to trunk
		parentBranch := cfg.TrunkBranch
		if branchInfo, exists := cfg.ManagedBranches[currentBranch]; exists {
			parentBranch = branchInfo.Parent
		}

		// Reset the last commit (keeping changes in working directory)
		if err := resetLastCommit(); err != nil {
			return err
		}

		// Checkout to parent branch
		if err := switchBranch(parentBranch); err != nil {
			return err
		}

		// Delete the branch (force delete since we know it's safe - we just reset the commit)
		if err := deleteBranch(currentBranch, false); err != nil {
			return err
		}

		// Remove branch from managed branches if it exists
		if _, exists := cfg.ManagedBranches[currentBranch]; exists {
			delete(cfg.ManagedBranches, currentBranch)
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
		}

		fmt.Printf("Popped branch '%s' and returned to '%s' with uncommitted changes\n", currentBranch, parentBranch)
		return nil
	},
}
