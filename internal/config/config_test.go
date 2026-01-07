package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultConfig(t *testing.T) {
	// Create a temporary directory to act as git root
	tempDir, err := os.MkdirTemp("", "gt-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .git directory to simulate git repo
	gitDir := filepath.Join(tempDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Load config (should return default)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.TrunkBranch != "main" {
		t.Errorf("Expected default trunk branch 'main', got '%s'", cfg.TrunkBranch)
	}

	if cfg.ManagedBranches == nil {
		t.Error("Expected ManagedBranches to be initialized")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary directory to act as git root
	tempDir, err := os.MkdirTemp("", "gt-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .git directory to simulate git repo
	gitDir := filepath.Join(tempDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Create and save config
	cfg := &Config{
		TrunkBranch: "develop",
		ManagedBranches: map[string]Branch{
			"feature-1": {
				Name:        "feature-1",
				Parent:      "develop",
				Description: "Feature 1",
			},
		},
	}

	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify config file is stored in .git/gt/ directory
	configPath := filepath.Join(gitDir, "gt", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file not created at expected path: %s", configPath)
	}

	// Load config
	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loadedCfg.TrunkBranch != "develop" {
		t.Errorf("Expected trunk branch 'develop', got '%s'", loadedCfg.TrunkBranch)
	}

	if len(loadedCfg.ManagedBranches) != 1 {
		t.Errorf("Expected 1 managed branch, got %d", len(loadedCfg.ManagedBranches))
	}

	branch, ok := loadedCfg.ManagedBranches["feature-1"]
	if !ok {
		t.Error("Expected to find 'feature-1' branch")
	}

	if branch.Parent != "develop" {
		t.Errorf("Expected parent 'develop', got '%s'", branch.Parent)
	}
}

func TestFindGitDir(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "gt-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .git directory
	gitDir := filepath.Join(tempDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	// Create subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// Change to subdirectory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(subDir); err != nil {
		t.Fatalf("Failed to change to subdir: %v", err)
	}

	// Find git directory
	foundGitDir, err := findGitDir()
	if err != nil {
		t.Fatalf("Failed to find git dir: %v", err)
	}

	if foundGitDir != gitDir {
		t.Errorf("Expected git dir '%s', got '%s'", gitDir, foundGitDir)
	}
}
