package services

import (
	"testing"

	"github.com/kumarasakti/passgen/internal/domain/entities"
)

func TestPasswordGenerator_GeneratePassword(t *testing.T) {
	generator := NewPasswordGenerator()

	tests := []struct {
		name    string
		config  entities.PasswordConfig
		wantErr bool
	}{
		{
			name: "valid generation",
			config: entities.PasswordConfig{
				Length:         12,
				IncludeLower:   true,
				IncludeUpper:   true,
				IncludeNumbers: true,
				IncludeSymbols: false,
				Count:          1,
			},
			wantErr: false,
		},
		{
			name: "invalid config",
			config: entities.PasswordConfig{
				Length: 0,
				Count:  1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := generator.GeneratePassword(tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if password.Value == "" {
					t.Error("Generated password should not be empty")
				}
				if len(password.Value) != tt.config.Length {
					t.Errorf("Password length = %v, want %v", len(password.Value), tt.config.Length)
				}
			}
		})
	}
}

func TestPasswordGenerator_GenerateMultiplePasswords(t *testing.T) {
	generator := NewPasswordGenerator()

	config := entities.PasswordConfig{
		Length:         8,
		IncludeLower:   true,
		IncludeUpper:   true,
		IncludeNumbers: true,
		Count:          3,
	}

	passwords, err := generator.GenerateMultiplePasswords(config)
	if err != nil {
		t.Errorf("GenerateMultiplePasswords() error = %v", err)
		return
	}

	if len(passwords) != config.Count {
		t.Errorf("Expected %d passwords, got %d", config.Count, len(passwords))
	}

	// Check that all passwords are unique
	seen := make(map[string]bool)
	for _, password := range passwords {
		if seen[password.Value] {
			t.Errorf("Found duplicate password: %s", password.Value)
		}
		seen[password.Value] = true

		if len(password.Value) != config.Length {
			t.Errorf("Password length = %v, want %v", len(password.Value), config.Length)
		}
	}
}

func TestPasswordAnalyzer_AnalyzePassword(t *testing.T) {
	analyzer := NewPasswordAnalyzer()

	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "short password",
			password: "123",
		},
		{
			name:     "common password",
			password: "password",
		},
		{
			name:     "mixed password",
			password: "Password123",
		},
		{
			name:     "complex password",
			password: "MySecur3P@ssw0rd!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password := entities.Password{Value: tt.password}
			config := entities.PasswordConfig{
				Length:         len(tt.password),
				IncludeLower:   true,
				IncludeUpper:   true,
				IncludeNumbers: true,
				IncludeSymbols: true,
				Count:          1,
			}

			analysis := analyzer.AnalyzePassword(password, config)

			if analysis.Password.Value != tt.password {
				t.Errorf("Analysis password = %v, want %v", analysis.Password.Value, tt.password)
			}

			// Basic checks that analysis contains expected data
			if analysis.TimeToCrack == "" {
				t.Error("TimeToCrack should not be empty")
			}

			if analysis.SecurityLevel == "" {
				t.Error("SecurityLevel should not be empty")
			}

			// Verify that strength is a valid value
			validStrengths := []entities.PasswordStrength{
				entities.VeryWeak, entities.Weak, entities.Medium,
				entities.Strong, entities.VeryStrong, entities.ExtremelyStrong,
			}
			validStrength := false
			for _, validStr := range validStrengths {
				if analysis.Strength == validStr {
					validStrength = true
					break
				}
			}
			if !validStrength {
				t.Errorf("Invalid strength value: %v", analysis.Strength)
			}
		})
	}
}
