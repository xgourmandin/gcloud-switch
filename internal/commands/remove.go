package commands

import (
	"fmt"
	"gcloud-switch/internal/config"
	"gcloud-switch/internal/logger"
	"os"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:               "remove <name>",
	Short:             "Remove an existing GCloud configuration",
	Long:              `Delete a configuration from the configuration store.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: GetConfigNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		configName := args[0]

		store, err := config.LoadConfigStore()
		if err != nil {
			return fmt.Errorf("failed to load configurations: %w", err)
		}

		// Get the config to check for saved ADC
		cfg, err := store.FindConfig(configName)
		if err == nil && cfg.ADCPath != "" {
			// Clean up saved ADC file
			if err := os.Remove(cfg.ADCPath); err != nil && !os.IsNotExist(err) {
				logger.Warning("Failed to remove saved ADC file", "error", err)
			}
		}

		if err := store.RemoveConfig(configName); err != nil {
			return fmt.Errorf("failed to remove configuration: %w", err)
		}

		if err := store.Save(); err != nil {
			return fmt.Errorf("failed to save changes: %w", err)
		}

		// Note: We don't delete the native gcloud configuration as the user might want to keep it
		// They can manually delete it with: gcloud config configurations delete <name>
		logger.Success("Successfully removed configuration", "name", configName)
		logger.Info("Note: Native gcloud configuration still exists. Delete manually if needed with: gcloud config configurations delete " + configName)
		return nil
	},
}
