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
	projectID      string
	serviceAccount string
)
var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new GCloud configuration",
	Long: `Create a new GCloud configuration with a project ID and optional 
service account for impersonation. If a native gcloud configuration with the 
same name already exists, it will be imported.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configName := args[0]

		store, err := config.LoadConfigStore()
		if err != nil {
			return fmt.Errorf("failed to load configurations: %w", err)
		}

		// Check if a native gcloud configuration already exists
		configExists := gcloud.ConfigurationExists(configName)

		var finalProjectID string

		if configExists {
			// Import existing gcloud configuration
			logger.Info("Found existing gcloud configuration", "name", configName)

			// Get the project ID from the existing configuration
			existingProject, err := gcloud.GetProjectFromConfiguration(configName)
			if err != nil || existingProject == "" {
				// If we can't get the project from the config, still allow user to set it
				logger.Warning("Could not retrieve project from existing configuration")

				reader := bufio.NewReader(os.Stdin)
				if projectID == "" {
					fmt.Print("Enter Project ID: ")
					projectID, _ = reader.ReadString('\n')
					projectID = strings.TrimSpace(projectID)
				}

				if projectID == "" {
					return fmt.Errorf("project ID is required")
				}
				finalProjectID = projectID
			} else {
				finalProjectID = existingProject
				logger.Info("Importing configuration", "project_id", finalProjectID)

				// If user provided a project flag that differs, warn them
				if projectID != "" && projectID != finalProjectID {
					logger.Warning("Ignoring provided project ID, using existing configuration's project", "existing", finalProjectID, "provided", projectID)
				}
			}
		} else {
			// Creating new configuration - prompt for project ID if not provided
			reader := bufio.NewReader(os.Stdin)

			if projectID == "" {
				fmt.Print("Enter Project ID: ")
				projectID, _ = reader.ReadString('\n')
				projectID = strings.TrimSpace(projectID)
			}

			if projectID == "" {
				return fmt.Errorf("project ID is required")
			}
			finalProjectID = projectID
		}

		// Service account is optional and can be set later via edit
		if serviceAccount == "" && !cmd.Flags().Changed("service-account") && !configExists {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter Service Account (optional, press Enter to skip): ")
			serviceAccount, _ = reader.ReadString('\n')
			serviceAccount = strings.TrimSpace(serviceAccount)
		}

		newConfig := config.GCloudConfig{
			Name:           configName,
			ProjectID:      finalProjectID,
			ServiceAccount: serviceAccount,
		}

		if err := store.AddConfig(newConfig); err != nil {
			return fmt.Errorf("failed to add configuration: %w", err)
		}

		// Create the native gcloud configuration only if it doesn't exist
		if !configExists {
			logger.Info("Creating gcloud configuration", "name", configName)
			if err := gcloud.CreateConfiguration(configName); err != nil {
				return fmt.Errorf("failed to create gcloud configuration: %w", err)
			}
		}

		if err := store.Save(); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		if configExists {
			logger.Success("Successfully imported existing configuration", "name", configName, "project_id", finalProjectID)
		} else {
			logger.Success("Successfully added configuration", "name", configName, "project_id", finalProjectID)
		}

		if serviceAccount != "" {
			logger.Info("  Service Account", "service_account", serviceAccount)
		} else if configExists {
			logger.Info("  No service account set. Use 'gcloud-switcher edit " + configName + "' to add one if needed.")
		}

		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&projectID, "project", "p", "", "GCloud Project ID")
	addCmd.Flags().StringVarP(&serviceAccount, "service-account", "s", "", "Service Account to impersonate (optional)")
}
