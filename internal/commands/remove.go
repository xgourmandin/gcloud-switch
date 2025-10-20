package commands

import (
	"fmt"
	"gcloud-switch/internal/config"
	"gcloud-switch/internal/logger"

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

		if err := store.RemoveConfig(configName); err != nil {
			return fmt.Errorf("failed to remove configuration: %w", err)
		}

		if err := store.Save(); err != nil {
			return fmt.Errorf("failed to save changes: %w", err)
		}

		logger.Success("Successfully removed configuration", "name", configName)
		return nil
	},
}
