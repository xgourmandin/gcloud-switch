package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	return buf.String(), err
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close() //nolint:errcheck
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r) //nolint:errcheck
	return buf.String()
}

func TestRootCommand(t *testing.T) {
	output, err := executeCommand(rootCmd, "--help")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !strings.Contains(output, "gcloud-switcher") {
		t.Errorf("Expected output to contain 'gcloud-switcher', got: %s", output)
	}

	if !strings.Contains(output, "GCloud Switcher") {
		t.Errorf("Expected output to contain 'GCloud Switcher', got: %s", output)
	}
}

func TestVersionCommand(t *testing.T) {
	// Set version info for testing
	Version = "1.0.0-test"
	Commit = "abc123"
	Date = "2025-10-20"
	BuiltBy = "test"

	output := captureStdout(func() {
		versionCmd.Run(versionCmd, []string{})
	})

	if !strings.Contains(output, "1.0.0-test") {
		t.Errorf("Expected output to contain version '1.0.0-test', got: %s", output)
	}
	if !strings.Contains(output, "abc123") {
		t.Errorf("Expected output to contain commit 'abc123', got: %s", output)
	}
	if !strings.Contains(output, "2025-10-20") {
		t.Errorf("Expected output to contain date '2025-10-20', got: %s", output)
	}
}

func TestVersionCommandDefaults(t *testing.T) {
	// Reset to defaults
	Version = "dev"
	Commit = "none"
	Date = "unknown"
	BuiltBy = "unknown"

	output := captureStdout(func() {
		versionCmd.Run(versionCmd, []string{})
	})

	if !strings.Contains(output, "dev") {
		t.Errorf("Expected output to contain version 'dev', got: %s", output)
	}
	if !strings.Contains(output, "none") {
		t.Errorf("Expected output to contain commit 'none', got: %s", output)
	}
}

func TestCommandsExist(t *testing.T) {
	commands := rootCmd.Commands()

	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		commandNames[cmd.Name()] = true
	}

	expectedCommands := []string{"list", "switch", "add", "edit", "remove", "current", "version", "completion"}

	for _, expected := range expectedCommands {
		if !commandNames[expected] {
			t.Errorf("Expected command '%s' to be registered, but it wasn't found", expected)
		}
	}
}

func TestGetConfigNamesWithNoConfig(t *testing.T) {
	// This test will fail if there's no config, which is expected
	// We're testing that the function doesn't panic
	names, directive := GetConfigNames(nil, []string{}, "")

	// Should return empty list or some names depending on config state
	if directive != 4 { // ShellCompDirectiveNoFileComp = 4
		t.Errorf("Expected directive to be 4 (NoFileComp), got: %d", directive)
	}

	// names can be nil or empty, both are valid
	if names == nil {
		names = []string{}
	}

	// Just checking it doesn't panic - length should be non-negative
	_ = names // Use names to avoid unused variable warning
}

func TestAddCommandFlags(t *testing.T) {
	// Test that add command has the expected flags
	projectFlag := addCmd.Flags().Lookup("project")
	if projectFlag == nil {
		t.Error("Expected 'project' flag to exist on add command")
		return
	}
	if projectFlag.Shorthand != "p" {
		t.Errorf("Expected 'project' flag shorthand to be 'p', got '%s'", projectFlag.Shorthand)
	}

	saFlag := addCmd.Flags().Lookup("service-account")
	if saFlag == nil {
		t.Error("Expected 'service-account' flag to exist on add command")
		return
	}
	if saFlag.Shorthand != "s" {
		t.Errorf("Expected 'service-account' flag shorthand to be 's', got '%s'", saFlag.Shorthand)
	}
}

func TestEditCommandFlags(t *testing.T) {
	// Test that edit command has the expected flags
	projectFlag := editCmd.Flags().Lookup("project")
	if projectFlag == nil {
		t.Error("Expected 'project' flag to exist on edit command")
	}

	saFlag := editCmd.Flags().Lookup("service-account")
	if saFlag == nil {
		t.Error("Expected 'service-account' flag to exist on edit command")
	}
}

func TestSwitchCommandArgs(t *testing.T) {
	// Test that switch command requires exactly 1 argument
	err := switchCmd.Args(switchCmd, []string{})
	if err == nil {
		t.Error("Expected error when switch command called with no arguments")
	}

	err = switchCmd.Args(switchCmd, []string{"config1"})
	if err != nil {
		t.Errorf("Expected no error when switch command called with 1 argument, got: %v", err)
	}

	err = switchCmd.Args(switchCmd, []string{"config1", "config2"})
	if err == nil {
		t.Error("Expected error when switch command called with 2 arguments")
	}
}

func TestRemoveCommandArgs(t *testing.T) {
	// Test that remove command requires exactly 1 argument
	err := removeCmd.Args(removeCmd, []string{})
	if err == nil {
		t.Error("Expected error when remove command called with no arguments")
	}

	err = removeCmd.Args(removeCmd, []string{"config1"})
	if err != nil {
		t.Errorf("Expected no error when remove command called with 1 argument, got: %v", err)
	}
}

func TestEditCommandArgs(t *testing.T) {
	// Test that edit command requires exactly 1 argument
	err := editCmd.Args(editCmd, []string{})
	if err == nil {
		t.Error("Expected error when edit command called with no arguments")
	}

	err = editCmd.Args(editCmd, []string{"config1"})
	if err != nil {
		t.Errorf("Expected no error when edit command called with 1 argument, got: %v", err)
	}
}
