// Package commands contains all CLI command implementations for gcloud-switcher.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the version string injected at build time
	Version = "dev"
	// Commit is the git commit hash injected at build time
	Commit = "none"
	// Date is the build date injected at build time
	Date = "unknown"
	// BuiltBy is the builder name injected at build time
	BuiltBy = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Display version, commit, and build information for gcloud-switcher.",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("gcloud-switcher version %s\n", Version)
		fmt.Printf("  commit: %s\n", Commit)
		fmt.Printf("  built at: %s\n", Date)
		fmt.Printf("  built by: %s\n", BuiltBy)
	},
}
