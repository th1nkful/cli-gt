package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/th1nkful/cli-gt/internal/config"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Update trunk and rebase tracked branches",
	Long:  `Updates trunk branch from origin, rebases local tracked branches onto trunk again. If a local tracked branch no longer exists on origin, prompts for confirmation (y/n) to delete the branch.`,
	RunE:  runSync,
}

func runSync(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Save the current branch to return to later
	currentBranch, err := getCurrentBranch()
	if err != nil {
		return err
	}

	// Step 1: Fetch from origin to update remote tracking branches
	fmt.Println("Fetching from origin...")
	if err := fetchFromOrigin(); err != nil {
		return err
	}

	// Step 2: Update the trunk branch from origin
	fmt.Printf("Updating %s from origin...\n", cfg.TrunkBranch)
	if err := updateTrunkBranch(cfg.TrunkBranch); err != nil {
		return err
	}

	// Step 3: Process managed branches
	branchesToDelete := []string{}
	branchesToRebase := []string{}

	for branchName := range cfg.ManagedBranches {
		// Check if branch exists locally
		localExists, remoteExists, err := branchExists(branchName)
		if err != nil {
			// If we can't check remote (e.g., network issue), skip delete check
			fmt.Printf("Warning: Could not check remote for branch '%s': %v\n", branchName, err)
			if localExists {
				branchesToRebase = append(branchesToRebase, branchName)
			}
			continue
		}

		if !localExists {
			// Branch doesn't exist locally, skip it
			continue
		}

		if !remoteExists {
			// Branch exists locally but not on origin
			branchesToDelete = append(branchesToDelete, branchName)
		} else {
			branchesToRebase = append(branchesToRebase, branchName)
		}
	}

	// Step 4: Prompt for deletion of branches that don't exist on origin
	reader := bufio.NewReader(os.Stdin)
	for _, branchName := range branchesToDelete {
		fmt.Printf("Branch '%s' no longer exists on origin. Delete local branch? (y/n): ", branchName)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response == "y" || response == "yes" {
			if err := syncDeleteBranch(branchName, cfg.TrunkBranch); err != nil {
				fmt.Printf("Warning: Failed to delete branch '%s': %v\n", branchName, err)
			} else {
				// Remove from managed branches and save config immediately
				delete(cfg.ManagedBranches, branchName)
				if err := cfg.Save(); err != nil {
					fmt.Printf("Warning: Failed to save config after deleting '%s': %v\n", branchName, err)
				}
				fmt.Printf("Deleted branch '%s'\n", branchName)
			}
		} else {
			// Still rebase the branch since user wants to keep it
			branchesToRebase = append(branchesToRebase, branchName)
		}
	}

	// Step 5: Rebase managed branches onto trunk
	for _, branchName := range branchesToRebase {
		fmt.Printf("Rebasing '%s' onto %s...\n", branchName, cfg.TrunkBranch)
		if err := rebaseBranchOntoTrunk(branchName, cfg.TrunkBranch); err != nil {
			fmt.Printf("Warning: Failed to rebase branch '%s': %v\n", branchName, err)
		}
	}

	// Step 6: Return to the original branch (if it still exists)
	localExists, _, _ := branchExists(currentBranch)
	if localExists {
		if err := syncCheckoutBranch(currentBranch); err != nil {
			fmt.Printf("Warning: Failed to return to branch '%s': %v\n", currentBranch, err)
		}
	} else {
		// If original branch was deleted, checkout trunk
		if err := syncCheckoutBranch(cfg.TrunkBranch); err != nil {
			fmt.Printf("Warning: Failed to checkout trunk branch '%s': %v\n", cfg.TrunkBranch, err)
		}
	}

	// Save updated config (in case branches were deleted)
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("Sync complete!")
	return nil
}

// fetchFromOrigin fetches updates from the origin remote
func fetchFromOrigin() error {
	cmd := exec.Command("git", "fetch", "origin")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to fetch from origin: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// updateTrunkBranch updates the trunk branch from origin
func updateTrunkBranch(trunkBranch string) error {
	// First checkout the trunk branch
	if err := syncCheckoutBranch(trunkBranch); err != nil {
		return err
	}

	// Pull latest changes (fast-forward only to avoid merge commits)
	cmd := exec.Command("git", "pull", "--ff-only", "origin", trunkBranch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update trunk branch: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// syncCheckoutBranch checks out the specified branch
func syncCheckoutBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to checkout branch '%s': %w\nOutput: %s", branchName, err, string(output))
	}
	return nil
}

// syncDeleteBranch deletes a local branch
func syncDeleteBranch(branchName, trunkBranch string) error {
	// First checkout trunk to avoid deleting the branch we're on
	if err := syncCheckoutBranch(trunkBranch); err != nil {
		return err
	}

	// Try safe delete first (-d), which fails if branch has unmerged commits
	cmd := exec.Command("git", "branch", "-d", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If safe delete fails, try force delete (-D)
		fmt.Printf("Note: Branch '%s' has unmerged commits, force deleting...\n", branchName)
		cmd = exec.Command("git", "branch", "-D", branchName)
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to delete branch '%s': %w\nOutput: %s", branchName, err, string(output))
		}
	}
	return nil
}

// rebaseBranchOntoTrunk rebases a branch onto the trunk branch
func rebaseBranchOntoTrunk(branchName, trunkBranch string) error {
	// Checkout the branch to rebase
	if err := syncCheckoutBranch(branchName); err != nil {
		return err
	}

	// Rebase onto trunk
	cmd := exec.Command("git", "rebase", trunkBranch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Abort the rebase if it fails
		abortCmd := exec.Command("git", "rebase", "--abort")
		if abortErr := abortCmd.Run(); abortErr != nil {
			fmt.Printf("Warning: Failed to abort rebase for '%s': %v\n", branchName, abortErr)
		}
		return fmt.Errorf("failed to rebase branch '%s' onto '%s': %w\nOutput: %s", branchName, trunkBranch, err, string(output))
	}
	return nil
}
