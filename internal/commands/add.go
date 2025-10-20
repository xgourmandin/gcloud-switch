package commands

import (
	"bufio"
	"fmt"
	"gcloud-switch/internal/config"
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
service account for impersonation.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configName := args[0]

		store, err := config.LoadConfigStore()
		if err != nil {
			return fmt.Errorf("failed to load configurations: %w", err)
		}

		// If flags not provided, prompt for them
		reader := bufio.NewReader(os.Stdin)

		if projectID == "" {
			fmt.Print("Enter Project ID: ")
			projectID, _ = reader.ReadString('\n')
			projectID = strings.TrimSpace(projectID)
		}

		if projectID == "" {
			return fmt.Errorf("project ID is required")
		}

		// Service account is optional
		if serviceAccount == "" && !cmd.Flags().Changed("service-account") {
			fmt.Print("Enter Service Account (optional, press Enter to skip): ")
			serviceAccount, _ = reader.ReadString('\n')
			serviceAccount = strings.TrimSpace(serviceAccount)
		}

		newConfig := config.GCloudConfig{
			Name:           configName,
			ProjectID:      projectID,
			ServiceAccount: serviceAccount,
		}

		if err := store.AddConfig(newConfig); err != nil {
			return fmt.Errorf("failed to add configuration: %w", err)
		}

		if err := store.Save(); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		logger.Success("Successfully added configuration", "name", configName, "project_id", projectID)
		if serviceAccount != "" {
			logger.Info("  Service Account", "service_account", serviceAccount)
		}

		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&projectID, "project", "p", "", "GCloud Project ID")
	addCmd.Flags().StringVarP(&serviceAccount, "service-account", "s", "", "Service Account to impersonate (optional)")
}
