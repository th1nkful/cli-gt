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
	Use:   "create [branch-name]",
	Short: "Create a new branch and commit",
	Long:  `Create a new branch and commit. This creates a new branch from the current state and makes an initial commit.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if we're on trunk branch
		onTrunk, err := isOnTrunkBranch()
		if err != nil {
			return err
		}
		if !onTrunk {
			cfg, _ := config.Load()
			trunkBranch := "main"
			if cfg != nil {
				trunkBranch = cfg.TrunkBranch
			}
			return fmt.Errorf("create command can only be run on trunk branch (%s)", trunkBranch)
		}

		// Determine branch name
		var branchName string
		if len(args) > 0 {
			// Use provided branch name
			branchName = args[0]
		} else if createMessage != "" {
			// Convert message to branch name
			branchName = sanitizeBranchName(createMessage)
		} else {
			return fmt.Errorf("must provide either a branch name or a message (-m)")
		}

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

		// Determine commit message
		commitMessage := createMessage
		if commitMessage == "" && len(args) > 0 {
			commitMessage = args[0]
		}
		if commitMessage == "" {
			return fmt.Errorf("must provide a commit message (-m)")
		}

		// Create the commit first (on trunk)
		if err := createCommit(commitMessage); err != nil {
			return err
		}

		// Create and checkout the new branch
		if err := createBranch(branchName); err != nil {
			return err
		}

		// Load config and add this branch as a managed branch
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

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
	createCmd.Flags().StringVarP(&createMessage, "message", "m", "", "Commit message (also used to generate branch name if not provided)")
}
