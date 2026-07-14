package application

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kumarasakti/passgen/internal/domain/entities"
	"gopkg.in/yaml.v3"
)

// ConfigManager handles loading and saving the passgen config file
type ConfigManager struct {
	configPath string
}

// NewConfigManager creates a ConfigManager using the default config path
func NewConfigManager() (*ConfigManager, error) {
	configPath, err := entities.ConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to determine config path: %w", err)
	}
	return &ConfigManager{configPath: configPath}, nil
}

// NewConfigManagerWithPath creates a ConfigManager with a custom config path (for testing)
func NewConfigManagerWithPath(configPath string) *ConfigManager {
	return &ConfigManager{configPath: configPath}
}

// Load reads the config file and returns the configuration.
// If the file doesn't exist, it returns the default config without error.
func (cm *ConfigManager) Load() (entities.PassgenConfig, error) {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return entities.DefaultPassgenConfig(), nil
		}
		return entities.PassgenConfig{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config entities.PassgenConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return entities.PassgenConfig{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// Save writes the config to disk, creating parent directories if needed
func (cm *ConfigManager) Save(config entities.PassgenConfig) error {
	dir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cm.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Init creates the config file with default values if it doesn't exist
func (cm *ConfigManager) Init() error {
	if _, err := os.Stat(cm.configPath); err == nil {
		return fmt.Errorf("config file already exists at %s", cm.configPath)
	}

	defaultConfig := entities.DefaultPassgenConfig()
	if err := cm.Save(defaultConfig); err != nil {
		return err
	}

	return nil
}

// Set updates a single config key and saves
func (cm *ConfigManager) Set(key string, value string) error {
	config, err := cm.Load()
	if err != nil {
		return err
	}

	g := &config.Generation
	w := &config.Word
	switch key {
	case "length":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil {
			return fmt.Errorf("invalid value for length: %s", value)
		}
		g.Length = n
	case "include_lower":
		g.IncludeLower = value == "true"
	case "include_upper":
		g.IncludeUpper = value == "true"
	case "include_numbers":
		g.IncludeNumbers = value == "true"
	case "include_symbols":
		g.IncludeSymbols = value == "true"
	case "exclude_similar":
		g.ExcludeSimilar = value == "true"
	case "exclude_chars":
		g.ExcludeChars = value
	case "no_repeat":
		g.NoRepeat = value == "true"
	case "count":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil {
			return fmt.Errorf("invalid value for count: %s", value)
		}
		g.Count = n
	case "word_strategy":
		w.Strategy = value
	case "word_complexity":
		w.Complexity = value
	case "word_count":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil {
			return fmt.Errorf("invalid value for word_count: %s", value)
		}
		w.Count = n
	default:
		return fmt.Errorf("unknown config key: %s (valid keys: length, include_lower, include_upper, include_numbers, include_symbols, exclude_similar, exclude_chars, no_repeat, count, word_strategy, word_complexity, word_count)", key)
	}

	return cm.Save(config)
}

// ConfigPath returns the path to the config file
func (cm *ConfigManager) ConfigPath() string {
	return cm.configPath
}
