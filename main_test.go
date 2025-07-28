package main

import (
	"testing"

	"github.com/kumarasakti/passgen/internal/infrastructure/cli"
)

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}

	if Version[0] != 'v' {
		t.Error("Version should start with 'v'")
	}
}

func TestMainIntegration(t *testing.T) {
	// Test that we can create the CLI handler without panicking
	handler := cli.NewHandler()
	if handler == nil {
		t.Error("CLI handler should not be nil")
	}

	// Test that we can create the root command
	rootCmd := handler.CreateRootCommand(Version)
	if rootCmd == nil {
		t.Error("Root command should not be nil")
		return
	}

	// Test that the version is properly set
	if rootCmd.Version != Version {
		t.Errorf("Expected version %s, got %s", Version, rootCmd.Version)
	}

	// Test that required commands exist
	commands := rootCmd.Commands()
	var hasWord, hasCheck bool

	for _, cmd := range commands {
		if cmd.Name() == "word" {
			hasWord = true
		}
		if cmd.Name() == "check" {
			hasCheck = true
		}
	}

	if !hasWord {
		t.Error("Root command should have 'word' subcommand")
	}

	if !hasCheck {
		t.Error("Root command should have 'check' subcommand")
	}
}
