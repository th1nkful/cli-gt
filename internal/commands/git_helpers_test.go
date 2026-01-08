package commands

import (
	"testing"
)

func TestSanitizeBranchName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Fix bug in parser", "fix-bug-in-parser"},
		{"Add new feature!", "add-new-feature"},
		{"Update README.md file", "update-readme-md-file"},
		{"feature/test-branch", "feature-test-branch"},
		{"Fix: issue #123", "fix-issue-123"},
		{"UPPERCASE MESSAGE", "uppercase-message"},
		{"Multiple   spaces   between", "multiple-spaces-between"},
		{"---leading-and-trailing---", "leading-and-trailing"},
		{"This is a very long branch name that exceeds the maximum length allowed for branch names", "this-is-a-very-long-branch-name-that-exceeds-the-m"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeBranchName(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeBranchName(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
