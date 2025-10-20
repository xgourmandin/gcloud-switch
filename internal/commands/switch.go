package commands

import (
	"fmt"
	"gcloud-switch/internal/config"
	"gcloud-switch/internal/gcloud"
	"gcloud-switch/internal/logger"

	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch <name>",
	Short: "Switch to the specified GCloud configuration",
	Long: `Switch to a predefined GCloud configuration. This will set the project 
and handle authentication automatically, reusing stored credentials when possible.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: GetConfigNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		configName := args[0]

		store, err := config.LoadConfigStore()
		if err != nil {
			return fmt.Errorf("failed to load configurations: %w", err)
		}

		cfg, err := store.FindConfig(configName)
		if err != nil {
			return fmt.Errorf("configuration '%s' not found", configName)
		}

		logger.Info("Switching to configuration", "name", cfg.Name, "project_id", cfg.ProjectID)

		// Set the project
		if err := gcloud.SetProject(cfg.ProjectID); err != nil {
			return err
		}
		logger.Success("Project set successfully")

		// Check if we need to authenticate
		needsAuth := !gcloud.CheckADCValid()

		if needsAuth {
			logger.Info("Authentication required...")

			if cfg.ServiceAccount != "" {
				logger.Info("Authenticating with service account", "service_account", cfg.ServiceAccount)
				if err := gcloud.AuthLoginWithServiceAccount(cfg.ServiceAccount); err != nil {
					return err
				}
			} else {
				logger.Info("Authenticating with user credentials...")
				if err := gcloud.AuthLogin(); err != nil {
					return err
				}
			}
			logger.Success("Authentication successful")
		} else {
			logger.Success("Using existing valid credentials")
		}

		// Update active config
		store.ActiveConfig = cfg.Name
		if err := store.Save(); err != nil {
			return fmt.Errorf("failed to save active configuration: %w", err)
		}

		logger.Success("Successfully switched to configuration", "name", cfg.Name)
		return nil
	},
}
