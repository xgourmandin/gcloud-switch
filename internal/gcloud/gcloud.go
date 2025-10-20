package gcloud

import (
	"fmt"
	"os/exec"
)

// SetProject sets the active GCloud project
func SetProject(projectID string) error {
	cmd := exec.Command("gcloud", "config", "set", "project", projectID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set project: %w\nOutput: %s", err, output)
	}
	return nil
}

// AuthLogin performs a standard gcloud auth login with ADC update
func AuthLogin() error {
	cmd := exec.Command("gcloud", "auth", "login", "--update-adc")
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}
	return nil
}

// AuthLoginWithServiceAccount performs authentication and sets up ADC for service account impersonation
func AuthLoginWithServiceAccount(serviceAccount string) error {
	// First, ensure user is logged in
	cmd := exec.Command("gcloud", "auth", "login")
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}
	// Then set up ADC with impersonation
	cmd = exec.Command("gcloud", "auth", "application-default", "login", "--impersonate-service-account", serviceAccount)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
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
