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
	Long: `Switch to a predefined GCloud configuration. This will activate the gcloud 
configuration and handle authentication automatically, reusing stored credentials when possible.`,
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

		// Step 1: Save ADC of current active configuration (if any)
		if store.ActiveConfig != "" && store.ActiveConfig != configName {
			currentCfg, err := store.FindConfig(store.ActiveConfig)
			if err == nil {
				adcPath, err := config.GetADCFileForConfig(store.ActiveConfig)
				if err == nil {
					logger.Info("Saving ADC for current configuration", "name", store.ActiveConfig)
					if err := gcloud.SaveADC(adcPath); err != nil {
						logger.Warning("Failed to save ADC", "error", err)
					} else {
						currentCfg.ADCPath = adcPath
					}
				}
			}
		}

		// Step 2: Ensure gcloud configuration exists, create if not
		if !gcloud.ConfigurationExists(configName) {
			logger.Info("Creating new gcloud configuration", "name", configName)
			if err := gcloud.CreateConfiguration(configName); err != nil {
				return fmt.Errorf("failed to create gcloud configuration: %w", err)
			}
		}

		// Step 3: Activate the gcloud configuration
		logger.Info("Activating gcloud configuration", "name", configName)
		if err := gcloud.ActivateConfiguration(configName); err != nil {
			return err
		}
		logger.Success("Configuration activated")

		// Step 4: Restore ADC if available for this configuration
		if cfg.ADCPath != "" {
			logger.Info("Restoring saved ADC credentials", "name", configName)
			if err := gcloud.RestoreADC(cfg.ADCPath); err != nil {
				logger.Warning("Failed to restore ADC", "error", err)
			} else {
				logger.Success("ADC credentials restored")
			}
		}

		// Step 5: Check if we need to authenticate (check both account and ADC)
		logger.Info("Checking authentication status...")
		accountValid := gcloud.CheckAccountValid()
		adcValid := gcloud.CheckADCValid()
		needsAuth := !accountValid || !adcValid

		if needsAuth {
			if !accountValid {
				logger.Info("Account credentials are invalid or expired")
			}
			if !adcValid {
				logger.Info("ADC credentials are invalid or expired")
			}
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

			// Save the new ADC credentials
			adcPath, err := config.GetADCFileForConfig(configName)
			if err == nil {
				if err := gcloud.SaveADC(adcPath); err != nil {
					logger.Warning("Failed to save new ADC", "error", err)
				} else {
					cfg.ADCPath = adcPath
				}
			}
		} else {
			logger.Success("Using existing valid credentials")
		}

		// Step 6: Set the project for this configuration (after auth is confirmed)
		if err := gcloud.SetProject(cfg.ProjectID); err != nil {
			return err
		}
		logger.Success("Project set successfully")

		// Update active config
		store.ActiveConfig = cfg.Name
		if err := store.Save(); err != nil {
			return fmt.Errorf("failed to save active configuration: %w", err)
		}

		logger.Success("Successfully switched to configuration", "name", cfg.Name)
		return nil
	},
}
