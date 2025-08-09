package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"github.com/kumarasakti/passgen/internal/domain/entities"
)

// ConfigManager handles store configuration file operations
type ConfigManager struct {
	configPath string
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".config", "passgen", "stores.yaml")
	
	return &ConfigManager{
		configPath: configPath,
	}
}

// LoadConfig loads the store configuration from file
func (c *ConfigManager) LoadConfig() (*entities.StoreConfig, error) {
	// Create default config if file doesn't exist
	if _, err := os.Stat(c.configPath); os.IsNotExist(err) {
		return c.createDefaultConfig(), nil
	}

	data, err := os.ReadFile(c.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config entities.StoreConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	config.ConfigPath = c.configPath
	return &config, nil
}

// SaveConfig saves the store configuration to file
func (c *ConfigManager) SaveConfig(config *entities.StoreConfig) error {
	// Ensure config directory exists
	if err := os.MkdirAll(filepath.Dir(c.configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(c.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetDefaultStore returns the default store name
func (c *ConfigManager) GetDefaultStore() (string, error) {
	config, err := c.LoadConfig()
	if err != nil {
		return "", err
	}
	
	if config.DefaultStore == "" {
		return "", fmt.Errorf("no default store configured")
	}
	
	return config.DefaultStore, nil
}

// SetDefaultStore sets the default store
func (c *ConfigManager) SetDefaultStore(storeName string) error {
	config, err := c.LoadConfig()
	if err != nil {
		return err
	}
	
	// Verify store exists
	if _, exists := config.Stores[storeName]; !exists {
		return fmt.Errorf("store '%s' does not exist", storeName)
	}
	
	// Update default and mark store as default
	config.DefaultStore = storeName
	for name, store := range config.Stores {
		store.IsDefault = (name == storeName)
		config.Stores[name] = store
	}
	
	return c.SaveConfig(config)
}

// createDefaultConfig creates a default configuration
func (c *ConfigManager) createDefaultConfig() *entities.StoreConfig {
	return &entities.StoreConfig{
		DefaultStore: "",
		Stores:       make(map[string]entities.PasswordStore),
		ConfigPath:   c.configPath,
		DefaultRotation: &entities.DefaultRotationConfig{
			IntervalDays:     90,
			NotifyDaysBefore: 7,
			AutoGenerate:     true,
			PasswordProfile: &entities.PasswordProfile{
				Length:         16,
				IncludeUpper:   true,
				IncludeLower:   true,
				IncludeNumbers: true,
				IncludeSymbols: true,
				CustomRules:    "no-ambiguous",
			},
		},
		Notifications: &entities.NotificationConfig{
			Enabled: false,
		},
	}
}

// GetConfigPath returns the configuration file path
func (c *ConfigManager) GetConfigPath() string {
	return c.configPath
}
