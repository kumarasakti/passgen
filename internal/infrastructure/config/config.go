package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// StorageConfig represents the storage configuration
type StorageConfig struct {
	Backend  string                    `yaml:"backend"`
	Settings map[string]string         `yaml:"settings"`
	Stores   map[string]StoreConfig    `yaml:"stores"`
}

// StoreConfig represents configuration for a single store
type StoreConfig struct {
	Name     string            `yaml:"name"`
	Backend  string            `yaml:"backend"`
	Settings map[string]string `yaml:"settings"`
	GPGKeyID string            `yaml:"gpg_key_id"`
}

// PassgenConfig represents the main passgen configuration
type PassgenConfig struct {
	Storage StorageConfig `yaml:"storage"`
	GPG     GPGConfig     `yaml:"gpg"`
}

// GPGConfig represents GPG configuration
type GPGConfig struct {
	DefaultKeyID string `yaml:"default_key_id"`
	GPGBinary    string `yaml:"gpg_binary"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *PassgenConfig {
	return &PassgenConfig{
		Storage: StorageConfig{
			Backend: "local",
			Settings: map[string]string{
				"base_path": "~/.passgen/stores",
			},
			Stores: make(map[string]StoreConfig),
		},
		GPG: GPGConfig{
			GPGBinary: "gpg",
		},
	}
}

// LoadConfig loads configuration from file
func LoadConfig() (*PassgenConfig, error) {
	configPath := getConfigPath()
	
	// If config doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := DefaultConfig()
		if err := SaveConfig(config); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return config, nil
	}
	
	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config PassgenConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return &config, nil
}

// SaveConfig saves configuration to file
func SaveConfig(config *PassgenConfig) error {
	configPath := getConfigPath()
	
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Marshal config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	// Write config file
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// AddStoreConfig adds a store configuration
func (config *PassgenConfig) AddStoreConfig(name string, storeConfig StoreConfig) {
	if config.Storage.Stores == nil {
		config.Storage.Stores = make(map[string]StoreConfig)
	}
	config.Storage.Stores[name] = storeConfig
}

// GetStoreConfig gets a store configuration
func (config *PassgenConfig) GetStoreConfig(name string) (StoreConfig, bool) {
	store, exists := config.Storage.Stores[name]
	return store, exists
}

// getConfigPath returns the configuration file path
func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".passgen", "config.yaml")
}

// expandPath expands ~ to user home directory
func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, path[1:])
	}
	return path
}
