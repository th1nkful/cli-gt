package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the workspace configuration
type Config struct {
	TrunkBranch    string            `json:"trunk_branch"`
	ManagedBranches map[string]Branch `json:"managed_branches"`
}

// Branch represents a managed branch
type Branch struct {
	Name        string `json:"name"`
	Parent      string `json:"parent"`
	Description string `json:"description"`
}

const (
	configDirName  = "gt"
	configFileName = "config.json"
)

// Load loads the configuration from the git workspace
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return &Config{
			TrunkBranch:     "main",
			ManagedBranches: make(map[string]Branch),
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if cfg.ManagedBranches == nil {
		cfg.ManagedBranches = make(map[string]Branch)
	}

	return &cfg, nil
}

// Save saves the configuration to the git workspace
func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Ensure the config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getConfigPath returns the path to the config file in the git workspace
// Config is stored inside .git/gt/ directory to keep it invisible and device-specific
func getConfigPath() (string, error) {
	gitDir, err := findGitDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(gitDir, configDirName, configFileName), nil
}

// findGitDir finds the .git directory of the git repository
func findGitDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	for {
		gitDir := filepath.Join(dir, ".git")
		if info, err := os.Stat(gitDir); err == nil {
			// Handle both regular .git directory and .git file (for worktrees)
			if info.IsDir() {
				return gitDir, nil
			}
			// If .git is a file (worktree), read it to find the actual git dir
			data, err := os.ReadFile(gitDir)
			if err != nil {
				return "", fmt.Errorf("failed to read .git file: %w", err)
			}
			// Parse "gitdir: /path/to/git/dir" format
			gitdirPrefix := "gitdir: "
			content := string(data)
			if len(content) > len(gitdirPrefix) && content[:len(gitdirPrefix)] == gitdirPrefix {
				actualGitDir := content[len(gitdirPrefix):]
				actualGitDir = filepath.Clean(actualGitDir[:len(actualGitDir)-1]) // Remove trailing newline
				return actualGitDir, nil
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not in a git repository")
		}
		dir = parent
	}
}
