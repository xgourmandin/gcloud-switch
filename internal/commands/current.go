package commands

import (
	"fmt"
	"gcloud-switch/internal/config"
	"gcloud-switch/internal/gcloud"
	"gcloud-switch/internal/logger"
	"strings"

	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the current active GCloud configuration",
	Long:  `Display information about the currently active configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := config.LoadConfigStore()
		if err != nil {
			return fmt.Errorf("failed to load configurations: %w", err)
		}

		// Get the actual active gcloud configuration
		activeGcloudConfig, err := gcloud.GetActiveConfiguration()
		if err == nil {
			activeGcloudConfig = strings.TrimSpace(activeGcloudConfig)
			logger.Info("Active GCloud Configuration", "name", activeGcloudConfig)
		}

		if store.ActiveConfig == "" {
			logger.Info("No active configuration tracked by gcloud-switcher.")
			return nil
		}

		cfg, err := store.FindConfig(store.ActiveConfig)
		if err != nil {
			return fmt.Errorf("active configuration not found: %w", err)
		}

		logger.Info("Current Active Configuration (gcloud-switcher):")
		logger.Info("================================================")
		logger.Info("Name", "name", cfg.Name)
		logger.Info("Project ID", "project_id", cfg.ProjectID)
		if cfg.ServiceAccount != "" {
			logger.Info("Service Account", "service_account", cfg.ServiceAccount)
		} else {
			logger.Info("Service Account: (none - using user credentials)")
		}

		// Also show current gcloud project
		currentProject, err := gcloud.GetCurrentProject()
		if err == nil {
			currentProject = strings.TrimSpace(currentProject)
			logger.Info("Current GCloud Project", "project", currentProject)
		}

		// Show if ADC is valid
		if gcloud.CheckADCValid() {
			logger.Success("ADC credentials are valid")
		} else {
			logger.Warning("ADC credentials are invalid or expired")
		}

		return nil
	},
}
