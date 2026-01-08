package commands

import (
	"testing"
)

func TestRootCommandHasSubcommands(t *testing.T) {
	// Verify all expected commands are registered
	expectedCommands := []string{"create", "pop", "modify", "checkout", "sync", "restack", "submit"}
	
	for _, cmdName := range expectedCommands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == cmdName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command '%s' to be registered", cmdName)
		}
	}
}

func TestCreateCommandExists(t *testing.T) {
	if createCmd.Use != "create [commit-message]" {
		t.Errorf("create command Use string is incorrect: %s", createCmd.Use)
	}
}

func TestPopCommandExists(t *testing.T) {
	if popCmd.Use != "pop" {
		t.Errorf("pop command Use string is incorrect: %s", popCmd.Use)
	}
}

func TestModifyCommandExists(t *testing.T) {
	if modifyCmd.Use != "modify" {
		t.Errorf("modify command Use string is incorrect: %s", modifyCmd.Use)
	}
}

func TestModifyCommandHasAllFlag(t *testing.T) {
	// Verify that the -a flag exists for modify command
	flag := modifyCmd.Flags().Lookup("all")
	if flag == nil {
		t.Error("Expected 'all' flag to exist for modify command")
	}
	
	// Verify the shorthand
	if flag != nil && flag.Shorthand != "a" {
		t.Errorf("Expected 'all' flag shorthand to be 'a', got '%s'", flag.Shorthand)
	}
}

func TestCheckoutCommandExists(t *testing.T) {
	if checkoutCmd.Use != "checkout [branch]" {
		t.Errorf("checkout command Use string is incorrect: %s", checkoutCmd.Use)
	}
}

func TestCheckoutAlias(t *testing.T) {
	// Verify that 'co' is an alias for checkout
	found := false
	for _, alias := range checkoutCmd.Aliases {
		if alias == "co" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected 'co' to be an alias for checkout")
	}
}

func TestSyncCommandExists(t *testing.T) {
	if syncCmd.Use != "sync" {
		t.Errorf("sync command Use string is incorrect: %s", syncCmd.Use)
	}
}

func TestRestackCommandExists(t *testing.T) {
	if restackCmd.Use != "restack" {
		t.Errorf("restack command Use string is incorrect: %s", restackCmd.Use)
	}
}

func TestSubmitCommandExists(t *testing.T) {
	if submitCmd.Use != "submit" {
		t.Errorf("submit command Use string is incorrect: %s", submitCmd.Use)
	}
}
