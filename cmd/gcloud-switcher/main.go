package main

import "gcloud-switch/internal/commands"

// Version information (set by goreleaser at build time)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	// Set version information for the version command
	commands.Version = version
	commands.Commit = commit
	commands.Date = date
	commands.BuiltBy = builtBy

	commands.Execute()
}
