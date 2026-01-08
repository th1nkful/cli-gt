package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/th1nkful/cli-gt/internal/config"
)

var (
	createAll     bool
	createMessage string
)

var createCmd = &cobra.Command{
	Use:   "create [commit-message]",
	Short: "Create a new branch and commit",
	Long:  `Create a new branch and commit. The commit message can be provided as a positional argument or via -m flag.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if we're on trunk branch and load config
		onTrunk, cfg, err := isOnTrunkBranch()
		if err != nil {
			return err
		}
		if !onTrunk {
			return fmt.Errorf("create command can only be run on trunk branch (%s)", cfg.TrunkBranch)
		}

		// Get commit message from positional arg or -m flag
		var commitMessage string
		if len(args) > 0 {
			// Positional argument takes precedence
			commitMessage = args[0]
		} else if createMessage != "" {
			// Fall back to -m flag
			commitMessage = createMessage
		} else {
			return fmt.Errorf("commit message is required (provide as argument or via -m flag)")
		}

		// Generate branch name from commit message
		branchName := sanitizeBranchName(commitMessage)

		// Check if branch already exists
		localExists, remoteExists, err := branchExists(branchName)
		if err != nil {
			return fmt.Errorf("failed to check if branch exists: %w", err)
		}
		if localExists {
			return fmt.Errorf("branch '%s' already exists locally", branchName)
		}
		if remoteExists {
			return fmt.Errorf("branch '%s' already exists on origin", branchName)
		}

		// Stage files if -a flag is used
		if createAll {
			if err := stageAllFiles(); err != nil {
				return err
			}
		}

		// Create and checkout the new branch first (before committing)
		if err := createBranch(branchName); err != nil {
			return err
		}

		// Create the commit on the new branch
		if err := createCommit(commitMessage); err != nil {
			return err
		}

		// Add this branch as a managed branch in config
		cfg.ManagedBranches[branchName] = config.Branch{
			Name:        branchName,
			Parent:      cfg.TrunkBranch,
			Description: commitMessage,
		}

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Created branch '%s' with commit: %s\n", branchName, commitMessage)
		return nil
	},
}

func init() {
	createCmd.Flags().BoolVarP(&createAll, "all", "a", false, "Stage all changes before committing")
	createCmd.Flags().StringVarP(&createMessage, "message", "m", "", "Commit message (used to generate branch name)")
}
