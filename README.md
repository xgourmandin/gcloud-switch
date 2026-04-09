# GCloud Switcher

[![CI](https://github.com/xgourmandin/gcloud-switch/workflows/CI/badge.svg)](https://github.com/xgourmandin/gcloud-switch/actions)
[![Release](https://github.com/xgourmandin/gcloud-switch/workflows/Release/badge.svg)](https://github.com/xgourmandin/gcloud-switch/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/xgourmandin/gcloud-switch)](https://goreportcard.com/report/github.com/xgourmandin/gcloud-switch)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A CLI tool to simplify switching between multiple GCloud configurations and projects.

## Overview

The main pain with GCloud when working on multiple organizations/projects is switching between them. This CLI tool simplifies that process by allowing you to define configurations once and switch between them with a single command.

## Installation

### Go Install

```bash
go install github.com/xgourmandin/gcloud-switch/cmd/gcloud-switch@latest
```

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/xgourmandin/gcloud-switch/releases).

Available for:
- **Linux**: amd64, arm64, armv6, armv7
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64


## Features

- ✅ **List configurations**: View all your saved GCloud configurations
- ✅ **Quick switching**: Switch between configurations with a single command
- ✅ **Automatic authentication**: Handles auth and ADC automatically, reusing credentials when valid
- ✅ **Service account support**: Optional service account impersonation per configuration
- ✅ **Simple management**: Add, edit, and remove configurations easily
- ✅ **Shell autocompletion**: Dynamic completion for bash, zsh, fish, and PowerShell
- ✅ **Colorful output**: Beautiful, human-readable terminal output
- ✅ **Cross-platform**: Works on Linux, macOS, and Windows

## Quick Start

```bash
# Add a configuration
gcloud-switch add dev-config -p my-dev-project-123

# List all configurations
gcloud-switch list

# Switch to a configuration
gcloud-switch switch dev-config

# View current configuration
gcloud-switch current
```

## Usage

### Add a new configuration

```bash
# Interactive mode
gcloud-switch add myconfig

# With flags
gcloud-switch add myconfig -p my-project-id -s my-sa@project.iam.gserviceaccount.com
```

### List all configurations

```bash
gcloud-switch list
```

### Switch to a configuration

```bash
gcloud-switch switch myconfig
```

The switch command will:
- Set the GCloud project
- Check if credentials are still valid
- Authenticate only if needed (with or without service account impersonation)
- Remember the active configuration

### Edit a configuration

```bash
# Interactive mode
gcloud-switch edit myconfig

# With flags
gcloud-switch edit myconfig -p new-project-id
```

### Remove a configuration

```bash
gcloud-switch remove myconfig
```

### View current active configuration

```bash
gcloud-switch current
```

## Shell Autocompletion

GCloud Switcher supports **dynamic autocompletion** for configuration names!

### Quick Setup

```bash
# Bash
gcloud-switch completion bash | sudo tee /etc/bash_completion.d/gcloud-switch

# Zsh
gcloud-switch completion zsh > "${fpath[1]}/_gcloud-switcher"

# Fish
gcloud-switch completion fish > ~/.config/fish/completions/gcloud-switch.fish

# PowerShell
gcloud-switch completion powershell | Out-String | Invoke-Expression
```

## Configuration Storage

Configurations are stored in `~/.gcloud-switcher/config.json` as a JSON file containing:
- Configuration name
- Project ID
- Optional service account for impersonation
- Currently active configuration

## Authentication

The tool intelligently handles authentication:

- **Without service account**: Uses standard `gcloud auth login --update-adc`
- **With service account**: Performs user login, then sets up ADC with `--impersonate-service-account`
- **Credentials reuse**: Checks if ADC is still valid before prompting for re-authentication

## Building from Source

```bash
# Clone the repository
git clone https://github.com/xgourmandin/gcloud-switch.git
cd gcloud-switch

# Download dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o gcloud-switch ./cmd/gcloud-switch

# Install
sudo mv gcloud-switch /usr/local/bin/
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Development

### Prerequisites

- Go 1.21 or later
- GCloud CLI (for testing)

### Running Tests

```bash
go test -v ./...
```

### Running Linters

```bash
golangci-lint run
```

### Local GoReleaser Test

```bash
goreleaser release --snapshot --clean
```

## Requirements

- Go 1.21 or later (for building from source)
- GCloud CLI installed and in PATH
- Appropriate GCloud permissions for the projects you're managing

## License

MIT - see [LICENSE](LICENSE) for details.

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [GoReleaser](https://goreleaser.com/) - Release automation
