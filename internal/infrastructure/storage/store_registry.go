package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// StoreConfig represents configuration for a store
type StoreConfig struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	GPGKeyID  string `json:"gpg_key_id"`
	LocalOnly bool   `json:"local_only"`
}

// Provides centralized management of password store configurations and metadata
type StoreRegistry struct {
	configPath string
	stores     map[string]StoreConfig
}

// Initializes store configuration management with persistent storage
func NewStoreRegistry() *StoreRegistry {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".passgen", "stores.json")

	registry := &StoreRegistry{
		configPath: configPath,
		stores:     make(map[string]StoreConfig),
	}

	registry.load()
	return registry
}

// RegisterStore registers a store in the registry
func (r *StoreRegistry) RegisterStore(config StoreConfig) error {
	r.stores[config.Name] = config
	return r.save()
}

// Retrieves specific store configuration by name from registry
func (r *StoreRegistry) GetStore(name string) (StoreConfig, error) {
	config, exists := r.stores[name]
	if !exists {
		return StoreConfig{}, fmt.Errorf("store '%s' not found", name)
	}
	return config, nil
}

// Provides complete inventory of all registered store configurations
func (r *StoreRegistry) ListStores() []StoreConfig {
	var configs []StoreConfig
	for _, config := range r.stores {
		configs = append(configs, config)
	}
	return configs
}

// load loads store configurations from disk
func (r *StoreRegistry) load() error {
	data, err := os.ReadFile(r.configPath)
	if os.IsNotExist(err) {
		return nil // No config file yet
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &r.stores)
}

// save saves store configurations to disk
func (r *StoreRegistry) save() error {
	// Ensure directory exists
	dir := filepath.Dir(r.configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(r.stores, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.configPath, data, 0600)
}
