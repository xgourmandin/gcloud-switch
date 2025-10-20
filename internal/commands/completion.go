package commands

import (
	"gcloud-switch/internal/config"

	"github.com/spf13/cobra"
)

// GetConfigNames returns a list of all configuration names for autocompletion
func GetConfigNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
