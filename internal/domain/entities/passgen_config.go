package entities

import (
	"os"
	"path/filepath"
)

// PassgenConfig represents the user's passgen configuration stored at ~/.passgen/config.yaml
type PassgenConfig struct {
	Generation GenerationConfig `yaml:"generation"`
}

// GenerationConfig holds default password generation preferences
type GenerationConfig struct {
	Length         int    `yaml:"length"`
	IncludeLower   bool   `yaml:"include_lower"`
	IncludeUpper   bool   `yaml:"include_upper"`
	IncludeNumbers bool   `yaml:"include_numbers"`
	IncludeSymbols bool   `yaml:"include_symbols"`
	ExcludeSimilar bool   `yaml:"exclude_similar"`
	ExcludeChars   string `yaml:"exclude_chars"`
	NoRepeat       bool   `yaml:"no_repeat"`
	Count          int    `yaml:"count"`
}

// DefaultPassgenConfig returns the default configuration matching the CLI defaults
func DefaultPassgenConfig() PassgenConfig {
	return PassgenConfig{
		Generation: GenerationConfig{
			Length:         DefaultLength,
			IncludeLower:   true,
			IncludeUpper:   true,
			IncludeNumbers: false,
			IncludeSymbols: true,
			ExcludeSimilar: false,
			ExcludeChars:   "",
			NoRepeat:       false,
			Count:          1,
		},
	}
}

// ToPasswordConfig converts PassgenConfig's generation settings to a PasswordConfig
func (pc PassgenConfig) ToPasswordConfig() PasswordConfig {
	g := pc.Generation
	return PasswordConfig{
		Length:         g.Length,
		IncludeLower:   g.IncludeLower,
		IncludeUpper:   g.IncludeUpper,
		IncludeNumbers: g.IncludeNumbers,
		IncludeSymbols: g.IncludeSymbols,
		ExcludeSimilar: g.ExcludeSimilar,
		ExcludeChars:   g.ExcludeChars,
		NoRepeat:       g.NoRepeat,
		Count:          g.Count,
	}
}

// ConfigPath returns the path to the config file at ~/.passgen/config.yaml
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".passgen", "config.yaml"), nil
}

// ConfigDir returns the passgen config directory at ~/.passgen
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".passgen"), nil
}
