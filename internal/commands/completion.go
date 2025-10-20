// Package commands contains all CLI command implementations for gcloud-switcher.
package commands

import (
	"gcloud-switch/internal/config"

	"github.com/spf13/cobra"
)

// GetConfigNames returns a list of all configuration names for autocompletion
func GetConfigNames(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	store, err := config.LoadConfigStore()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var names []string
	for _, cfg := range store.Configurations {
		names = append(names, cfg.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
