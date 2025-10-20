package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version information (injected at build time)
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
	BuiltBy = "unknown"
)
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Display version, commit, and build information for gcloud-switcher.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gcloud-switcher version %s\n", Version)
		fmt.Printf("  commit: %s\n", Commit)
		fmt.Printf("  built at: %s\n", Date)
		fmt.Printf("  built by: %s\n", BuiltBy)
	},
}
