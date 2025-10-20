# GCloud configuration switcher

The main pain with GCloud when working on multiple organizations/projects is switching between them. This CLI tool simplifies that process.

Instead of going through the classcial `gcloud config set project <project-id>` dance (then `gcloud auth login`, etc...), this tool allows you to quickly switch between predefined configurations.

## Features
- List available GCloud configurations.
- Switch to a specified configuration.
- Add new configurations / Edit existing ones.
- Remove existing configurations.
- View the current active configuration.

## Installation
To install the CLI tool, run:
```bash
go install github.com/yourusername/gcloud-switcher@latest
```

## Usage
After installation, you can use the CLI tool with the following commands:
```bash
gcloud-switcher list               # List all available configurations
gcloud-switcher switch <name>      # Switch to the specified configuration
gcloud-switcher add <name>         # Add a new configuration
gcloud-switcher edit <name>        # Edit an existing configuration
gcloud-switcher remove <name>      # Remove an existing configuration
gcloud-switcher current            # Show the current active configuration
```

## Features details
### Create and Manage Configurations
The CLI tool allows you to create a GCLoud configuration by setting the traditional project ID, but also setting an optional service account to impersonate.

If you don't specify a service account, the default user credentials will be used.

### Switch Configurations
Switching configurations is as simple as running the `switch` command with the desired configuration name. The tool will handle updating the GCloud settings accordingly.

In order to not ask for a login every time, the ADC (Application Default Credentials) are stored when switching, so when you come back to an existing configuration, a login is asked only if the credentials have expired.

If a login is required : 
- If no service account is set for the configuration, a standard `gcloud auth login --update-adc` is performed.
- If a service account is set, a serie of `gcloud auth login` followed by a `gcloud auth application-default login --impersonate-service-account <sa name>` is performed.

### Edit Configurations
You can edit existing configurations to update the project ID or service account associated with them.

### Other features
The other features behave the same as their `gcloud` counterparts, but are simplified to work with the predefined configurations.