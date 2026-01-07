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

const configFileName = ".gt-config.json"

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
func getConfigPath() (string, error) {
	gitRoot, err := findGitRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(gitRoot, configFileName), nil
}

// findGitRoot finds the root of the git repository
func findGitRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	for {
		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not in a git repository")
		}
		dir = parent
	}
}
