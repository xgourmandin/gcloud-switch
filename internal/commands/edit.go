package commands

import (
	"bufio"
	"fmt"
	"gcloud-switch/internal/config"
	"gcloud-switch/internal/gcloud"
	"gcloud-switch/internal/logger"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	editProjectID      string
	editServiceAccount string
)
var editCmd = &cobra.Command{
	Use:               "edit <name>",
	Short:             "Edit an existing GCloud configuration",
	Long:              `Update the project ID or service account for an existing configuration.`,
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

		logger.Info("Editing configuration", "name", configName)
		logger.Info("Current Project ID", "project_id", cfg.ProjectID)
		logger.Info("Current Service Account", "service_account", cfg.ServiceAccount)
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)

		// If flags not provided, prompt for them
		if editProjectID == "" && !cmd.Flags().Changed("project") {
			fmt.Printf("Enter new Project ID (or press Enter to keep current): ")
			editProjectID, _ = reader.ReadString('\n')
			editProjectID = strings.TrimSpace(editProjectID)
		}

		if editServiceAccount == "" && !cmd.Flags().Changed("service-account") {
			fmt.Printf("Enter new Service Account (or press Enter to keep current): ")
			editServiceAccount, _ = reader.ReadString('\n')
			editServiceAccount = strings.TrimSpace(editServiceAccount)
		}

		projectChanged := false

		// Update only if new values provided
		if editProjectID != "" {
			cfg.ProjectID = editProjectID
			projectChanged = true
		}

		if cmd.Flags().Changed("service-account") || editServiceAccount != "" {
			cfg.ServiceAccount = editServiceAccount
		}

		if err := store.Save(); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		// If project ID changed, update it in the native gcloud configuration
		if projectChanged && gcloud.ConfigurationExists(configName) {
			logger.Info("Updating native gcloud configuration with new project ID")

			// Get current active configuration to restore it later
			currentActive, _ := gcloud.GetActiveConfiguration()
			currentActive = strings.TrimSpace(currentActive)

			// Only activate the configuration if it's not already active
			needsRestore := false
			if currentActive != configName {
				if err := gcloud.ActivateConfiguration(configName); err != nil {
					logger.Warning("Failed to activate configuration for update", "error", err)
					// Continue anyway, might still work
				} else {
					needsRestore = true
				}
			}

			// Set the project ID
			if err := gcloud.SetProject(cfg.ProjectID); err != nil {
				logger.Warning("Failed to update project in native gcloud configuration", "error", err)
			} else {
				logger.Success("Native gcloud configuration updated")
			}

			// Restore the previously active configuration if we changed it
			if needsRestore && currentActive != "" {
				if err := gcloud.ActivateConfiguration(currentActive); err != nil {
					logger.Warning("Failed to restore previous active configuration", "error", err)
				}
			}
		}

		logger.Success("Successfully updated configuration", "name", configName, "project_id", cfg.ProjectID)
		if cfg.ServiceAccount != "" {
			logger.Info("  Service Account", "service_account", cfg.ServiceAccount)
		} else {
			logger.Info("  Service Account: (none)")
		}

		return nil
	},
}

func init() {
	editCmd.Flags().StringVarP(&editProjectID, "project", "p", "", "New GCloud Project ID")
	editCmd.Flags().StringVarP(&editServiceAccount, "service-account", "s", "", "New Service Account to impersonate")
}
