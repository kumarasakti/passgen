package infrastructure

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kumarasakti/passgen/internal/domain/entities"
)

func TestConfigManager_LoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "stores.yaml")
	
	configManager := &ConfigManager{
		configPath: configPath,
	}
	
	t.Run("load default config when file doesn't exist", func(t *testing.T) {
		config, err := configManager.LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig() error = %v, want nil", err)
		}
		
		if config == nil {
			t.Fatal("LoadConfig() returned nil config")
		}
		
		if config.DefaultStore != "" {
			t.Errorf("Default config should have empty DefaultStore, got %q", config.DefaultStore)
		}
		
		if config.Stores == nil {
			t.Error("Default config should have initialized Stores map")
		}
		
		if config.DefaultRotation == nil {
			t.Error("Default config should have DefaultRotation")
		}
		
		if config.DefaultRotation.IntervalDays != 90 {
			t.Errorf("Default rotation interval should be 90 days, got %d", config.DefaultRotation.IntervalDays)
		}
	})
	
	t.Run("save and load config", func(t *testing.T) {
		// Create a test config
		testConfig := &entities.StoreConfig{
			DefaultStore: "personal",
			Stores: map[string]entities.PasswordStore{
				"personal": {
					Name:      "personal",
					GitURL:    "git@github.com:user/personal-passwords.git",
					LocalPath: "/home/user/.password-stores/personal",
					GPGKeyID:  "passgen-personal@localhost",
					IsDefault: true,
					CreatedAt: time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC),
				},
			},
			ConfigPath: configPath,
		}
		
		// Save config
		err := configManager.SaveConfig(testConfig)
		if err != nil {
			t.Fatalf("SaveConfig() error = %v, want nil", err)
		}
		
		// Load config
		loadedConfig, err := configManager.LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig() error = %v, want nil", err)
		}
		
		// Verify loaded config
		if loadedConfig.DefaultStore != testConfig.DefaultStore {
			t.Errorf("DefaultStore = %q, want %q", loadedConfig.DefaultStore, testConfig.DefaultStore)
		}
		
		if len(loadedConfig.Stores) != 1 {
			t.Errorf("Expected 1 store, got %d", len(loadedConfig.Stores))
		}
		
		personalStore, exists := loadedConfig.Stores["personal"]
		if !exists {
			t.Error("Personal store not found in loaded config")
		}
		
		if personalStore.Name != "personal" {
			t.Errorf("Store name = %q, want %q", personalStore.Name, "personal")
		}
		
		if personalStore.GitURL != "git@github.com:user/personal-passwords.git" {
			t.Errorf("Store GitURL = %q, want %q", personalStore.GitURL, "git@github.com:user/personal-passwords.git")
		}
	})
}

func TestConfigManager_GetSetDefaultStore(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "stores.yaml")
	
	configManager := &ConfigManager{
		configPath: configPath,
	}
	
	// Create initial config with stores
	testConfig := &entities.StoreConfig{
		DefaultStore: "",
		Stores: map[string]entities.PasswordStore{
			"personal": {
				Name:      "personal",
				GitURL:    "git@github.com:user/personal.git",
				LocalPath: "/home/user/.stores/personal",
				IsDefault: false,
			},
			"work": {
				Name:      "work",
				GitURL:    "git@company.com:work.git",
				LocalPath: "/home/user/.stores/work",
				IsDefault: false,
			},
		},
	}
	
	err := configManager.SaveConfig(testConfig)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}
	
	t.Run("get default store when none set", func(t *testing.T) {
		_, err := configManager.GetDefaultStore()
		if err == nil {
			t.Error("GetDefaultStore() should return error when no default set")
		}
	})
	
	t.Run("set and get default store", func(t *testing.T) {
		// Set default store
		err := configManager.SetDefaultStore("personal")
		if err != nil {
			t.Fatalf("SetDefaultStore() error = %v", err)
		}
		
		// Get default store
		defaultStore, err := configManager.GetDefaultStore()
		if err != nil {
			t.Fatalf("GetDefaultStore() error = %v", err)
		}
		
		if defaultStore != "personal" {
			t.Errorf("GetDefaultStore() = %q, want %q", defaultStore, "personal")
		}
		
		// Verify the store is marked as default
		config, err := configManager.LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}
		
		if !config.Stores["personal"].IsDefault {
			t.Error("Personal store should be marked as default")
		}
		
		if config.Stores["work"].IsDefault {
			t.Error("Work store should not be marked as default")
		}
	})
	
	t.Run("set default to non-existent store", func(t *testing.T) {
		err := configManager.SetDefaultStore("nonexistent")
		if err == nil {
			t.Error("SetDefaultStore() should return error for non-existent store")
		}
	})
}

func TestConfigManager_CreateDefaultConfig(t *testing.T) {
	configManager := &ConfigManager{
		configPath: "/tmp/test-config.yaml",
	}
	
	config := configManager.createDefaultConfig()
	
	if config == nil {
		t.Fatal("createDefaultConfig() returned nil")
	}
	
	if config.DefaultStore != "" {
		t.Errorf("Default config should have empty DefaultStore, got %q", config.DefaultStore)
	}
	
	if config.Stores == nil {
		t.Error("Default config should have initialized Stores map")
	}
	
	if len(config.Stores) != 0 {
		t.Errorf("Default config should have empty Stores map, got %d entries", len(config.Stores))
	}
	
	if config.DefaultRotation == nil {
		t.Fatal("Default config should have DefaultRotation")
	}
	
	// Verify default rotation settings
	if config.DefaultRotation.IntervalDays != 90 {
		t.Errorf("Default interval should be 90 days, got %d", config.DefaultRotation.IntervalDays)
	}
	
	if config.DefaultRotation.NotifyDaysBefore != 7 {
		t.Errorf("Default notify days should be 7, got %d", config.DefaultRotation.NotifyDaysBefore)
	}
	
	if !config.DefaultRotation.AutoGenerate {
		t.Error("Default config should have AutoGenerate enabled")
	}
	
	if config.DefaultRotation.PasswordProfile == nil {
		t.Fatal("Default config should have PasswordProfile")
	}
	
	// Verify password profile defaults
	profile := config.DefaultRotation.PasswordProfile
	if profile.Length != 16 {
		t.Errorf("Default password length should be 16, got %d", profile.Length)
	}
	
	if !profile.IncludeUpper || !profile.IncludeLower || !profile.IncludeNumbers || !profile.IncludeSymbols {
		t.Error("Default profile should include all character types")
	}
	
	if profile.CustomRules != "no-ambiguous" {
		t.Errorf("Default custom rules should be 'no-ambiguous', got %q", profile.CustomRules)
	}
	
	// Verify notifications config
	if config.Notifications == nil {
		t.Fatal("Default config should have Notifications")
	}
	
	if config.Notifications.Enabled {
		t.Error("Default notifications should be disabled")
	}
}

func TestConfigManager_FilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "stores.yaml")
	
	configManager := &ConfigManager{
		configPath: configPath,
	}
	
	// Create and save a config
	testConfig := &entities.StoreConfig{
		DefaultStore: "test",
		Stores:       make(map[string]entities.PasswordStore),
	}
	
	err := configManager.SaveConfig(testConfig)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}
	
	// Check file permissions
	fileInfo, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}
	
	// Config file should be readable/writable by owner only (600)
	expectedPerm := os.FileMode(0600)
	if fileInfo.Mode().Perm() != expectedPerm {
		t.Errorf("Config file permissions = %o, want %o", fileInfo.Mode().Perm(), expectedPerm)
	}
}
