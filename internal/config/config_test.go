package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGCloudConfig(t *testing.T) {
	config := GCloudConfig{
		Name:           "test-config",
		ProjectID:      "test-project-123",
		ServiceAccount: "test@project.iam.gserviceaccount.com",
	}
	if config.Name != "test-config" {
		t.Errorf("Expected Name to be 'test-config', got '%s'", config.Name)
	}
	if config.ProjectID != "test-project-123" {
		t.Errorf("Expected ProjectID to be 'test-project-123', got '%s'", config.ProjectID)
	}
	if config.ServiceAccount != "test@project.iam.gserviceaccount.com" {
		t.Errorf("Expected ServiceAccount to be 'test@project.iam.gserviceaccount.com', got '%s'", config.ServiceAccount)
	}
}
func TestConfigStoreSaveAndLoad(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gcloud-switcher-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	configPath := filepath.Join(tempDir, "config.json")
	store := &ConfigStore{
		Configurations: []GCloudConfig{
			{
				Name:           "dev",
				ProjectID:      "dev-project-123",
				ServiceAccount: "dev@project.iam.gserviceaccount.com",
			},
			{
				Name:      "prod",
				ProjectID: "prod-project-456",
			},
		},
		ActiveConfig: "dev",
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
	loadedData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	var loadedStore ConfigStore
	err = json.Unmarshal(loadedData, &loadedStore)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}
	if len(loadedStore.Configurations) != 2 {
		t.Errorf("Expected 2 configurations, got %d", len(loadedStore.Configurations))
	}
	if loadedStore.ActiveConfig != "dev" {
		t.Errorf("Expected ActiveConfig to be 'dev', got '%s'", loadedStore.ActiveConfig)
	}
}
func TestConfigStoreFindConfig(t *testing.T) {
	store := &ConfigStore{
		Configurations: []GCloudConfig{
			{Name: "dev", ProjectID: "dev-project"},
			{Name: "prod", ProjectID: "prod-project"},
		},
	}
	config, err := store.FindConfig("dev")
	if err != nil {
		t.Errorf("Expected to find 'dev' config, got error: %v", err)
	}
	if config.ProjectID != "dev-project" {
		t.Errorf("Expected ProjectID 'dev-project', got '%s'", config.ProjectID)
	}
	_, err = store.FindConfig("staging")
	if err == nil {
		t.Error("Expected error when finding non-existing config, got nil")
	}
}
func TestConfigStoreAddConfig(t *testing.T) {
	store := &ConfigStore{
		Configurations: []GCloudConfig{
			{Name: "dev", ProjectID: "dev-project"},
		},
	}
	newConfig := GCloudConfig{
		Name:      "prod",
		ProjectID: "prod-project",
	}
	err := store.AddConfig(newConfig)
	if err != nil {
		t.Errorf("Expected to add config successfully, got error: %v", err)
	}
	if len(store.Configurations) != 2 {
		t.Errorf("Expected 2 configurations after adding, got %d", len(store.Configurations))
	}
	err = store.AddConfig(newConfig)
	if err == nil {
		t.Error("Expected error when adding duplicate config, got nil")
	}
}
func TestConfigStoreRemoveConfig(t *testing.T) {
	store := &ConfigStore{
		Configurations: []GCloudConfig{
			{Name: "dev", ProjectID: "dev-project"},
			{Name: "prod", ProjectID: "prod-project"},
		},
		ActiveConfig: "dev",
	}
	err := store.RemoveConfig("prod")
	if err != nil {
		t.Errorf("Expected to remove config successfully, got error: %v", err)
	}
	if len(store.Configurations) != 1 {
		t.Errorf("Expected 1 configuration after removing, got %d", len(store.Configurations))
	}
	err = store.RemoveConfig("dev")
	if err != nil {
		t.Errorf("Expected to remove config successfully, got error: %v", err)
	}
	if store.ActiveConfig != "" {
		t.Errorf("Expected ActiveConfig to be empty after removing active config, got '%s'", store.ActiveConfig)
	}
	err = store.RemoveConfig("staging")
	if err == nil {
		t.Error("Expected error when removing non-existing config, got nil")
	}
}
func TestConfigStoreUpdateConfig(t *testing.T) {
	store := &ConfigStore{
		Configurations: []GCloudConfig{
			{Name: "dev", ProjectID: "dev-project", ServiceAccount: "old@project.iam.gserviceaccount.com"},
		},
	}
	err := store.UpdateConfig("dev", "new-dev-project", "new@project.iam.gserviceaccount.com")
	if err != nil {
		t.Errorf("Expected to update config successfully, got error: %v", err)
	}
	config, _ := store.FindConfig("dev")
	if config.ProjectID != "new-dev-project" {
		t.Errorf("Expected ProjectID to be 'new-dev-project', got '%s'", config.ProjectID)
	}
	if config.ServiceAccount != "new@project.iam.gserviceaccount.com" {
		t.Errorf("Expected ServiceAccount to be 'new@project.iam.gserviceaccount.com', got '%s'", config.ServiceAccount)
	}
	err = store.UpdateConfig("staging", "staging-project", "")
	if err == nil {
		t.Error("Expected error when updating non-existing config, got nil")
	}
}
