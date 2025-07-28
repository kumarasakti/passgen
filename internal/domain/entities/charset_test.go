package entities

import (
	"strings"
	"testing"
)

func TestCharacterSet_BuildCharset(t *testing.T) {
	cs := NewCharacterSet()

	tests := []struct {
		name           string
		config         PasswordConfig
		wantContains   []string
		wantNotContains []string
		wantErr        bool
	}{
		{
			name: "all character sets",
			config: PasswordConfig{
				Length:         12,
				IncludeLower:   true,
				IncludeUpper:   true,
				IncludeNumbers: true,
				IncludeSymbols: true,
				Count:          1,
			},
			wantContains:   []string{"a", "Z", "5", "!"},
			wantNotContains: []string{},
			wantErr:        false,
		},
		{
			name: "exclude similar characters",
			config: PasswordConfig{
				Length:         12,
				IncludeLower:   true,
				IncludeUpper:   true,
				IncludeNumbers: true,
				ExcludeSimilar: true,
				Count:          1,
			},
			wantContains:   []string{"a", "Z", "5"},
			wantNotContains: []string{"i", "l", "1", "L", "o", "0", "O"},
			wantErr:        false,
		},
		{
			name: "custom exclusions",
			config: PasswordConfig{
				Length:       12,
				IncludeLower: true,
				IncludeUpper: true,
				ExcludeChars: "aeiou",
				Count:        1,
			},
			wantContains:   []string{"b", "Z"},
			wantNotContains: []string{"a", "e", "i", "o", "u"},
			wantErr:        false,
		},
		{
			name: "lowercase only",
			config: PasswordConfig{
				Length:       8,
				IncludeLower: true,
				Count:        1,
			},
			wantContains:   []string{"a", "z"},
			wantNotContains: []string{"A", "1", "!"},
			wantErr:        false,
		},
		{
			name: "invalid config",
			config: PasswordConfig{
				Length: 0,
				Count:  1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charset, err := cs.BuildCharset(tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("CharacterSet.BuildCharset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			for _, char := range tt.wantContains {
				if !strings.Contains(charset, char) {
					t.Errorf("Expected charset to contain %s", char)
				}
			}

			for _, char := range tt.wantNotContains {
				if strings.Contains(charset, char) {
					t.Errorf("Expected charset to not contain %s", char)
				}
			}
		})
	}
}

func TestCharacterSet_Constants(t *testing.T) {
	// Test that character constants are not empty
	if Lowercase == "" {
		t.Error("Lowercase constant should not be empty")
	}
	if Uppercase == "" {
		t.Error("Uppercase constant should not be empty")
	}
	if Numbers == "" {
		t.Error("Numbers constant should not be empty")
	}
	if Symbols == "" {
		t.Error("Symbols constant should not be empty")
	}

	// Test default length
	if DefaultLength <= 0 {
		t.Error("DefaultLength should be positive")
	}
}
