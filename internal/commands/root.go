package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gcloud-switcher",
	Short: "A CLI tool to simplify switching between GCloud configurations",
	Long: `GCloud Switcher helps you manage multiple GCloud configurations and 
switch between them effortlessly. Instead of manually running multiple 
gcloud commands, you can define configurations with project IDs and 
optional service accounts, then switch between them with a single command.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func init() {
	// Add all subcommands here
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(versionCmd)
}
