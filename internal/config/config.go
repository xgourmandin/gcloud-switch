package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// GCloudConfig represents a single GCloud configuration
type GCloudConfig struct {
	Name           string `json:"name"`
	ProjectID      string `json:"project_id"`
	ServiceAccount string `json:"service_account,omitempty"`
}

// ConfigStore manages all configurations
type ConfigStore struct {
	Configurations []GCloudConfig `json:"configurations"`
	ActiveConfig   string         `json:"active_config,omitempty"`
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".gcloud-switcher")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

// LoadConfigStore loads the configuration store from disk
func LoadConfigStore() (*ConfigStore, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If file doesn't exist, return empty store
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &ConfigStore{
			Configurations: []GCloudConfig{},
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var store ConfigStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}

	return &store, nil
}

// Save saves the configuration store to disk
func (cs *ConfigStore) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// FindConfig finds a configuration by name
func (cs *ConfigStore) FindConfig(name string) (*GCloudConfig, error) {
	for i := range cs.Configurations {
		if cs.Configurations[i].Name == name {
			return &cs.Configurations[i], nil
		}
	}
	return nil, errors.New("configuration not found")
}

// AddConfig adds a new configuration
func (cs *ConfigStore) AddConfig(config GCloudConfig) error {
	// Check if already exists
	for _, c := range cs.Configurations {
		if c.Name == config.Name {
			return errors.New("configuration with this name already exists")
		}
	}
	cs.Configurations = append(cs.Configurations, config)
	return nil
}

// RemoveConfig removes a configuration by name
func (cs *ConfigStore) RemoveConfig(name string) error {
	for i, c := range cs.Configurations {
		if c.Name == name {
			cs.Configurations = append(cs.Configurations[:i], cs.Configurations[i+1:]...)
			if cs.ActiveConfig == name {
				cs.ActiveConfig = ""
			}
			return nil
		}
	}
	return errors.New("configuration not found")
}

// UpdateConfig updates an existing configuration
func (cs *ConfigStore) UpdateConfig(name string, projectID, serviceAccount string) error {
	config, err := cs.FindConfig(name)
	if err != nil {
		return err
	}
	config.ProjectID = projectID
	config.ServiceAccount = serviceAccount
	return nil
}
