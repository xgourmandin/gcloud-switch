package gcloud

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetADCPath returns the standard location of the ADC file
func GetADCPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "gcloud", "application_default_credentials.json"), nil
}

// SaveADC saves the current ADC file to a specified location
func SaveADC(destPath string) error {
	adcPath, err := GetADCPath()
	if err != nil {
		return err
	}

	// Check if ADC file exists
	if _, err := os.Stat(adcPath); os.IsNotExist(err) {
		// No ADC file to save, this is not an error
		return nil
	}

	// Copy the ADC file
	source, err := os.Open(adcPath)
	if err != nil {
		return fmt.Errorf("failed to open ADC file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("failed to copy ADC file: %w", err)
	}

	return nil
}

// RestoreADC restores an ADC file from a saved location
func RestoreADC(sourcePath string) error {
	// Check if source file exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		// No saved ADC file, this is not an error
		return nil
	}

	adcPath, err := GetADCPath()
	if err != nil {
		return err
	}

	// Ensure the directory exists
	adcDir := filepath.Dir(adcPath)
	if err := os.MkdirAll(adcDir, 0755); err != nil {
		return fmt.Errorf("failed to create ADC directory: %w", err)
	}

	// Copy the saved ADC file back
	source, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open saved ADC file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(adcPath)
	if err != nil {
		return fmt.Errorf("failed to create ADC file: %w", err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("failed to restore ADC file: %w", err)
	}

	return nil
}

// DeleteADC removes the current ADC file
func DeleteADC() error {
	adcPath, err := GetADCPath()
	if err != nil {
		return err
	}

	// Check if ADC file exists
	if _, err := os.Stat(adcPath); os.IsNotExist(err) {
		// No ADC file to delete
		return nil
	}

	return os.Remove(adcPath)
}

// ActivateConfiguration activates a gcloud configuration by name
func ActivateConfiguration(configName string) error {
	cmd := exec.Command("gcloud", "config", "configurations", "activate", configName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to activate configuration: %w\nOutput: %s", err, output)
	}
	return nil
}

// CreateConfiguration creates a new gcloud configuration
func CreateConfiguration(configName string) error {
	cmd := exec.Command("gcloud", "config", "configurations", "create", configName, "--no-activate")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create configuration: %w\nOutput: %s", err, output)
	}
	return nil
}

// ConfigurationExists checks if a gcloud configuration exists
func ConfigurationExists(configName string) bool {
	cmd := exec.Command("gcloud", "config", "configurations", "describe", configName)
	err := cmd.Run()
	return err == nil
}

// GetActiveConfiguration returns the name of the currently active gcloud configuration
func GetActiveConfiguration() (string, error) {
	cmd := exec.Command("gcloud", "config", "configurations", "list", "--filter=is_active:true", "--format=value(name)")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get active configuration: %w", err)
	}
	return string(output), nil
}

// GetAccountFromConfiguration gets the account from a specific gcloud configuration
func GetAccountFromConfiguration(configName string) (string, error) {
	cmd := exec.Command("gcloud", "config", "configurations", "describe", configName, "--format=value(properties.core.account)")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get account from configuration: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetProjectFromConfiguration gets the project ID from a specific gcloud configuration
func GetProjectFromConfiguration(configName string) (string, error) {
	cmd := exec.Command("gcloud", "config", "configurations", "describe", configName, "--format=value(properties.core.project)")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get project from configuration: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// SetProject sets the active GCloud project
func SetProject(projectID string) error {
	// Use --no-user-output-enabled to prevent interactive prompts
	cmd := exec.Command("gcloud", "config", "set", "project", projectID, "--quiet")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set project: %w\nOutput: %s", err, output)
	}
	return nil
}

// AuthLogin performs a standard gcloud auth login with ADC update
func AuthLogin() error {
	cmd := exec.Command("gcloud", "auth", "login", "--update-adc")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}
	return nil
}

// AuthLoginWithServiceAccount performs authentication and sets up ADC for service account impersonation
func AuthLoginWithServiceAccount(serviceAccount string) error {
	// First, ensure user is logged in
	cmd := exec.Command("gcloud", "auth", "login")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}
	// Then set up ADC with impersonation
	cmd = exec.Command("gcloud", "auth", "application-default", "login", "--impersonate-service-account", serviceAccount)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set up service account impersonation: %w", err)
	}
	return nil
}

// GetCurrentProject returns the currently active project
func GetCurrentProject() (string, error) {
	cmd := exec.Command("gcloud", "config", "get-value", "project")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current project: %w", err)
	}
	return string(output), nil
}

// CheckADCValid checks if Application Default Credentials are still valid
func CheckADCValid() bool {
	cmd := exec.Command("gcloud", "auth", "application-default", "print-access-token")
	err := cmd.Run()
	return err == nil
}

// CheckAccountValid checks if the account credentials are still valid
func CheckAccountValid() bool {
	cmd := exec.Command("gcloud", "auth", "print-access-token")
	err := cmd.Run()
	return err == nil
}
