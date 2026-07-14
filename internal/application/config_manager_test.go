package application

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kumarasakti/passgen/internal/domain/entities"
)

func TestConfigManager_LoadDefaultWhenNoFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nonexistent.yaml")
	cm := NewConfigManagerWithPath(configPath)

	config, err := cm.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	defaults := entities.DefaultPassgenConfig()
	if config.Generation.Length != defaults.Generation.Length {
		t.Errorf("Length = %d, want %d", config.Generation.Length, defaults.Generation.Length)
	}
	if config.Generation.IncludeLower != defaults.Generation.IncludeLower {
		t.Errorf("IncludeLower = %v, want %v", config.Generation.IncludeLower, defaults.Generation.IncludeLower)
	}
}

func TestConfigManager_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	cm := NewConfigManagerWithPath(configPath)

	customConfig := entities.PassgenConfig{
		Generation: entities.GenerationConfig{
			Length:         24,
			IncludeLower:   true,
			IncludeUpper:   true,
			IncludeNumbers: true,
			IncludeSymbols: true,
			ExcludeSimilar: true,
			ExcludeChars:   "!@#",
			NoRepeat:       true,
			Count:          3,
		},
	}

	if err := cm.Save(customConfig); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("Config file not created: %v", err)
	}

	loaded, err := cm.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.Generation.Length != 24 {
		t.Errorf("Length = %d, want 24", loaded.Generation.Length)
	}
	if loaded.Generation.IncludeNumbers != true {
		t.Errorf("IncludeNumbers = %v, want true", loaded.Generation.IncludeNumbers)
	}
	if loaded.Generation.ExcludeChars != "!@#" {
		t.Errorf("ExcludeChars = %q, want %q", loaded.Generation.ExcludeChars, "!@#")
	}
	if loaded.Generation.NoRepeat != true {
		t.Errorf("NoRepeat = %v, want true", loaded.Generation.NoRepeat)
	}
	if loaded.Generation.Count != 3 {
		t.Errorf("Count = %d, want 3", loaded.Generation.Count)
	}
}

func TestConfigManager_InitCreatesDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	cm := NewConfigManagerWithPath(configPath)

	if err := cm.Init(); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	loaded, err := cm.Load()
	if err != nil {
		t.Fatalf("Load() after Init() error: %v", err)
	}

	defaults := entities.DefaultPassgenConfig()
	if loaded.Generation.Length != defaults.Generation.Length {
		t.Errorf("Length = %d, want %d", loaded.Generation.Length, defaults.Generation.Length)
	}
}

func TestConfigManager_InitFailsIfExists(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	cm := NewConfigManagerWithPath(configPath)

	if err := cm.Init(); err != nil {
		t.Fatalf("First Init() error: %v", err)
	}

	err := cm.Init()
	if err == nil {
		t.Error("Second Init() should return error")
	}
}

func TestConfigManager_Set(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		wantValue interface{}
		checkFunc func(entities.PassgenConfig) bool
	}{
		{
			name:      "set length",
			key:       "length",
			value:     "32",
			checkFunc: func(c entities.PassgenConfig) bool { return c.Generation.Length == 32 },
		},
		{
			name:      "set include_numbers true",
			key:       "include_numbers",
			value:     "true",
			checkFunc: func(c entities.PassgenConfig) bool { return c.Generation.IncludeNumbers == true },
		},
		{
			name:      "set include_symbols false",
			key:       "include_symbols",
			value:     "false",
			checkFunc: func(c entities.PassgenConfig) bool { return c.Generation.IncludeSymbols == false },
		},
		{
			name:      "set no_repeat true",
			key:       "no_repeat",
			value:     "true",
			checkFunc: func(c entities.PassgenConfig) bool { return c.Generation.NoRepeat == true },
		},
		{
			name:      "set exclude_chars",
			key:       "exclude_chars",
			value:     "il1Lo0O",
			checkFunc: func(c entities.PassgenConfig) bool { return c.Generation.ExcludeChars == "il1Lo0O" },
		},
		{
			name:      "set count",
			key:       "count",
			value:     "5",
			checkFunc: func(c entities.PassgenConfig) bool { return c.Generation.Count == 5 },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			cm := NewConfigManagerWithPath(configPath)

			if err := cm.Init(); err != nil {
				t.Fatalf("Init() error: %v", err)
			}

			if err := cm.Set(tt.key, tt.value); err != nil {
				t.Fatalf("Set(%s, %s) error: %v", tt.key, tt.value, err)
			}

			loaded, err := cm.Load()
			if err != nil {
				t.Fatalf("Load() error: %v", err)
			}

			if !tt.checkFunc(loaded) {
				t.Errorf("Set(%s, %s) did not produce expected result", tt.key, tt.value)
			}
		})
	}
}

func TestConfigManager_SetInvalidKey(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	cm := NewConfigManagerWithPath(configPath)

	if err := cm.Init(); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	err := cm.Set("nonexistent_key", "value")
	if err == nil {
		t.Error("Set() with invalid key should return error")
	}
}

func TestConfigManager_SavePreservesExistingValues(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	cm := NewConfigManagerWithPath(configPath)

	if err := cm.Init(); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	if err := cm.Set("length", "40"); err != nil {
		t.Fatalf("Set() error: %v", err)
	}
	if err := cm.Set("no_repeat", "true"); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	loaded, err := cm.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.Generation.Length != 40 {
		t.Errorf("Length = %d, want 40", loaded.Generation.Length)
	}
	if loaded.Generation.NoRepeat != true {
		t.Errorf("NoRepeat = %v, want true", loaded.Generation.NoRepeat)
	}

	defaults := entities.DefaultPassgenConfig()
	if loaded.Generation.IncludeLower != defaults.Generation.IncludeLower {
		t.Errorf("IncludeLower = %v, want %v (default should be preserved)",
			loaded.Generation.IncludeLower, defaults.Generation.IncludeLower)
	}
	if loaded.Generation.IncludeSymbols != defaults.Generation.IncludeSymbols {
		t.Errorf("IncludeSymbols = %v, want %v (default should be preserved)",
			loaded.Generation.IncludeSymbols, defaults.Generation.IncludeSymbols)
	}
}

func TestPassgenConfig_ToPasswordConfig(t *testing.T) {
	config := entities.PassgenConfig{
		Generation: entities.GenerationConfig{
			Length:         20,
			IncludeLower:   true,
			IncludeUpper:   false,
			IncludeNumbers: true,
			IncludeSymbols: true,
			ExcludeSimilar: true,
			ExcludeChars:   "xyz",
			NoRepeat:       true,
			Count:          2,
		},
	}

	pc := config.ToPasswordConfig()

	if pc.Length != 20 {
		t.Errorf("Length = %d, want 20", pc.Length)
	}
	if pc.IncludeUpper != false {
		t.Errorf("IncludeUpper = %v, want false", pc.IncludeUpper)
	}
	if pc.ExcludeChars != "xyz" {
		t.Errorf("ExcludeChars = %q, want %q", pc.ExcludeChars, "xyz")
	}
	if pc.NoRepeat != true {
		t.Errorf("NoRepeat = %v, want true", pc.NoRepeat)
	}
	if pc.Count != 2 {
		t.Errorf("Count = %d, want 2", pc.Count)
	}
}
