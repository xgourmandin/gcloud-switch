package commands

import (
	"fmt"
	"gcloud-switch/internal/config"
	"gcloud-switch/internal/logger"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available GCloud configurations",
	Long:  `Display a list of all configured GCloud configurations with their details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := config.LoadConfigStore()
		if err != nil {
			return fmt.Errorf("failed to load configurations: %w", err)
		}

		if len(store.Configurations) == 0 {
			logger.Info("No configurations found. Use 'gcloud-switcher add' to create one.")
			return nil
		}

		logger.Info("Available GCloud Configurations:")
		logger.Info("================================")
		for _, cfg := range store.Configurations {
			activeMarker := ""
			if cfg.Name == store.ActiveConfig {
				activeMarker = " (active)"
			}
			logger.Info("", "name", cfg.Name+activeMarker)
			logger.Info("  Project ID", "project_id", cfg.ProjectID)
			if cfg.ServiceAccount != "" {
				logger.Info("  Service Account", "service_account", cfg.ServiceAccount)
			} else {
				logger.Info("  Service Account: (none - using user credentials)")
			}
		}

		return nil
	},
}
