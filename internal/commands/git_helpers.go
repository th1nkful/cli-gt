package commands

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/th1nkful/cli-gt/internal/config"
)

// getCurrentBranch returns the name of the current git branch
func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// isOnTrunkBranch checks if the current branch is the trunk branch
func isOnTrunkBranch() (bool, error) {
	cfg, err := config.Load()
	if err != nil {
		return false, fmt.Errorf("failed to load config: %w", err)
	}

	currentBranch, err := getCurrentBranch()
	if err != nil {
		return false, err
	}

	return currentBranch == cfg.TrunkBranch, nil
}

// branchExists checks if a branch exists locally or remotely
func branchExists(branchName string) (bool, bool, error) {
	// Check local branches
	localCmd := exec.Command("git", "rev-parse", "--verify", branchName)
	localExists := localCmd.Run() == nil

	// Check remote branches
	remoteCmd := exec.Command("git", "ls-remote", "--heads", "origin", branchName)
	remoteOutput, err := remoteCmd.CombinedOutput()
	if err != nil {
		// Check if the error is due to no remote configured
		if strings.Contains(string(remoteOutput), "does not appear to be a git repository") ||
			strings.Contains(err.Error(), "exit status 128") {
			// No remote configured, that's ok - just return local status
			return localExists, false, nil
		}
		// Other errors (network issues, auth failures) should be reported
		return localExists, false, fmt.Errorf("failed to check remote branches: %w", err)
	}
	remoteExists := len(strings.TrimSpace(string(remoteOutput))) > 0

	return localExists, remoteExists, nil
}

// sanitizeBranchName converts a message into a valid git branch name
func sanitizeBranchName(message string) string {
	// Convert to lowercase
	name := strings.ToLower(message)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	name = reg.ReplaceAllString(name, "-")

	// Remove leading/trailing hyphens
	name = strings.Trim(name, "-")

	// Limit length to 50 characters
	if len(name) > 50 {
		name = name[:50]
	}

	// Remove trailing hyphen if trimmed at a hyphen
	name = strings.TrimRight(name, "-")

	return name
}

// stageAllFiles stages all changes (equivalent to git add -A)
func stageAllFiles() error {
	cmd := exec.Command("git", "add", "-A")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stage files: %w", err)
	}
	return nil
}

// createCommit creates a commit with the given message
func createCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create commit: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// createBranch creates a new branch with the given name
func createBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}
	return nil
}
