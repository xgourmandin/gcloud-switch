# Golang CLI Application: Architecture & Development Best Practices

This document outlines a set of best practices for building high-quality Command Line Interface (CLI) applications using Golang. It is intended to serve as a guide for developers and as a context file for GitHub Copilot to promote idiomatic and robust code.

## 1. Project Structure

A well-organized project structure is crucial for maintainability and scalability. The standard Go project layout is a great starting point.

```
/my-cli
├── /cmd
│   └── /my-cli           # Main application entry point
│       └── main.go
├── /internal             # Private application and library code
│   ├── /app              # Application-specific business logic
│   │   └── run.go
│   ├── /config           # Configuration management
│   │   └── config.go
│   ├── /errors           # Custom error types
│   │   └── errors.go
│   └── /commands         # CLI command definitions (using Cobra)
│       ├── root.go
│       ├── version.go
│       └── mycommand.go
├── /pkg                  # Public library code (if any)
├── go.mod                # Go modules file
├── go.sum
└── README.md
```

- `/cmd`: Contains the main package for your executable. The entry point should be minimal, responsible only for initializing and executing the root command.
- `/internal`: The bulk of your application logic resides here. It's not importable by other projects, which enforces clean separation.
- `/pkg`: Use this directory for code that you intend to be shared with other projects. If you don't have any, you don't need this directory.

## 2. CLI Framework: Cobra

While the standard library's `flag` package is sufficient for simple CLIs, a dedicated framework like Cobra provides a more robust foundation for complex applications with subcommands.

### Why Cobra?

- Easy creation of commands and subcommands.
- Powerful flag parsing (POSIX-compliant).
- Intelligent suggestions (e.g., `app commnad` -> `app command`).
- Automatic generation of help text, shell autocompletion, and man pages.

### Example root.go:

```go
// internal/commands/root.go
package commands

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "my-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
}
```

## 3. Configuration Management: Viper

CLIs often need configuration from flags, environment variables, and config files. Viper is a perfect companion to Cobra for this.

### Why Viper?

- Seamlessly handles multiple configuration formats (JSON, YAML, TOML, etc.).
- Prioritizes configuration sources: flags > env variables > config file > defaults.
- Binds flags directly to configuration keys.
- Live-reloads configuration files.

### Integration with Cobra:

```go
// internal/commands/root.go
import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.my-cli.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.SetDefault("author", "default author")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigName(".my-cli")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
```

## 4. Graceful Error Handling

Good error handling is critical for a good user experience.

- **Exit Codes**: Use non-zero exit codes for errors. `os.Exit(1)` is a common practice.
- **User-Friendly Messages**: Don't just print the raw error. Provide context and potential solutions.
- **Structured Errors**: For complex applications, consider creating custom error types to handle different error scenarios programmatically.
- **Logging**: Use a structured logger like `log/slog` (Go 1.21+) or `zerolog` to distinguish between user-facing errors (printed to stderr) and debug/internal logs.

```go
// internal/errors/errors.go
package errors

import "fmt"

type UserError struct {
	Message string
}

func (e *UserError) Error() string {
	return e.Message
}

func NewUserError(format string, a ...interface{}) error {
	return &UserError{
		Message: fmt.Sprintf(format, a...),
	}
}
```

In your command logic, you can check the error type:

```go
// In a cobra command's RunE function
if err != nil {
	if _, ok := err.(*errors.UserError); ok {
		// This is an error we want to show the user directly
		return err
	}
	// This is an unexpected/internal error
	log.Error("an internal error occurred", "error", err)
	return errors.New("an unexpected internal error occurred; please check the logs for more details")
}
```

Use Cobra's `RunE` field for commands, which allows you to return an error. Cobra will then print the error and exit.

## 5. Testing

- **Unit Tests**: Test business logic in `/internal` packages thoroughly, independent of the CLI framework.
- **Command Tests**: Cobra makes it easy to test commands by simulating command-line arguments and capturing output.

### Example Command Test:

```go
// internal/commands/version_test.go
package commands

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {
	// Redirect stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the command
	rootCmd.SetArgs([]string{"version"})
	Execute()

	// Restore stdout
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Assert the output
	output := buf.String()
	assert.True(t, strings.HasPrefix(output, "my-cli version:"))
}
```

## 6. Cross-Compilation and Distribution

Go's built-in cross-compilation is a major advantage. Use a tool like GoReleaser to automate the build, release, and packaging process.

### GoReleaser Features:

- Cross-compiles for multiple OS/architectures in parallel.
- Creates archives (.tar.gz, .zip).
- Generates checksums.
- Creates Linux packages (.deb, .rpm) and Homebrew taps.
- Pushes releases to GitHub.

A simple `.goreleaser.yml` is all you need to get started.

## 7. Concurrency

- Use goroutines for long-running tasks to keep the UI responsive, but be mindful. A CLI is often a short-lived process.
- Use a `sync.WaitGroup` to ensure all goroutines finish before the main function exits.
- Use channels for communication between goroutines.
- Provide feedback to the user for long operations, using a spinner or progress bar library like `github.com/schollz/progressbar`.

## 8. Code Style and Linting

- **gofmt**: Always format your code with `gofmt`.
- **golangci-lint**: Use `golangci-lint` as a comprehensive linter. It runs many linters in parallel and is highly configurable. Integrate it into your CI pipeline.
