package services

import (
	"strings"
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

// TestPasswordGenerator_DefaultAllowsDuplicates confirms the default (NoRepeat=false)
// path samples with replacement — duplicates are permitted and entropy is maximized.
func TestPasswordGenerator_DefaultAllowsDuplicates(t *testing.T) {
	generator := NewPasswordGenerator()

	config := entities.PasswordConfig{
		Length:         20,
		IncludeLower:   true,
		IncludeUpper:   true,
		IncludeNumbers: true,
		IncludeSymbols: true,
		Count:          1,
		NoRepeat:       false,
	}

	foundDuplicate := false
	for iter := 0; iter < 1000 && !foundDuplicate; iter++ {
		password, err := generator.GeneratePassword(config)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		seen := make(map[byte]bool)
		for i := 0; i < len(password.Value); i++ {
			c := password.Value[i]
			if seen[c] {
				foundDuplicate = true
				break
			}
			seen[c] = true
		}
	}

	if !foundDuplicate {
		t.Error("Default path should permit duplicate characters (with-replacement sampling)")
	}
}

func TestPasswordGenerator_NoDuplicateCharacters(t *testing.T) {
	generator := NewPasswordGenerator()

	config := entities.PasswordConfig{
		Length:         20,
		IncludeLower:   true,
		IncludeUpper:   true,
		IncludeNumbers: true,
		IncludeSymbols: true,
		Count:          1,
		NoRepeat:       true,
	}

	// Run many iterations to catch any randomness edge cases
	for iter := 0; iter < 100; iter++ {
		password, err := generator.GeneratePassword(config)
		if err != nil {
			t.Fatalf("Iteration %d: unexpected error: %v", iter, err)
		}

		seen := make(map[byte]bool)
		for i := 0; i < len(password.Value); i++ {
			c := password.Value[i]
			if seen[c] {
				t.Fatalf("Iteration %d: duplicate character %q found in password %q",
					iter, string(c), password.Value)
			}
			seen[c] = true
		}
	}
}

func TestPasswordGenerator_GuaranteedCharacterTypeCoverage(t *testing.T) {
	generator := NewPasswordGenerator()

	config := entities.PasswordConfig{
		Length:         16,
		IncludeLower:   true,
		IncludeUpper:   true,
		IncludeNumbers: true,
		IncludeSymbols: true,
		Count:          1,
		NoRepeat:       true,
	}

	// Run many iterations — every password must contain all 4 character types
	for iter := 0; iter < 200; iter++ {
		password, err := generator.GeneratePassword(config)
		if err != nil {
			t.Fatalf("Iteration %d: unexpected error: %v", iter, err)
		}

		if !password.HasLowercase() {
			t.Fatalf("Iteration %d: password %q missing lowercase letters", iter, password.Value)
		}
		if !password.HasUppercase() {
			t.Fatalf("Iteration %d: password %q missing uppercase letters", iter, password.Value)
		}
		if !password.HasNumbers() {
			t.Fatalf("Iteration %d: password %q missing numbers", iter, password.Value)
		}
		if !password.HasSymbols() {
			t.Fatalf("Iteration %d: password %q missing symbols", iter, password.Value)
		}
	}
}

func TestPasswordGenerator_LengthExceedsCharsetSize(t *testing.T) {
	generator := NewPasswordGenerator()

	// Lowercase only = 26 chars, requesting 30 with NoRepeat should fail
	config := entities.PasswordConfig{
		Length:         30,
		IncludeLower:   true,
		IncludeUpper:   false,
		IncludeNumbers: false,
		IncludeSymbols: false,
		Count:          1,
		NoRepeat:       true,
	}

	_, err := generator.GeneratePassword(config)
	if err == nil {
		t.Error("Expected error when length exceeds available unique characters with --no-repeat")
	}
	if err != nil && !strings.Contains(err.Error(), "exceeds available unique characters") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// TestPasswordGenerator_LengthExceedsCharsetDefaultOK confirms the default path
// allows length > charset size (with-replacement sampling has no such constraint).
func TestPasswordGenerator_LengthExceedsCharsetDefaultOK(t *testing.T) {
	generator := NewPasswordGenerator()

	config := entities.PasswordConfig{
		Length:         30,
		IncludeLower:   true,
		IncludeUpper:   false,
		IncludeNumbers: false,
		IncludeSymbols: false,
		Count:          1,
		NoRepeat:       false,
	}

	password, err := generator.GeneratePassword(config)
	if err != nil {
		t.Fatalf("Default path should allow length > charset size: %v", err)
	}
	if len(password.Value) != 30 {
		t.Errorf("Password length = %d, want 30", len(password.Value))
	}
}

func TestPasswordGenerator_ShortPasswordRelaxesGuarantee(t *testing.T) {
	generator := NewPasswordGenerator()

	// Length 2 with 4 categories enabled: can't guarantee all 4, but should
	// still produce a valid password with no duplicates
	config := entities.PasswordConfig{
		Length:         2,
		IncludeLower:   true,
		IncludeUpper:   true,
		IncludeNumbers: true,
		IncludeSymbols: true,
		Count:          1,
		NoRepeat:       true,
	}

	password, err := generator.GeneratePassword(config)
	if err != nil {
		t.Fatalf("Unexpected error for short password: %v", err)
	}

	if len(password.Value) != 2 {
		t.Errorf("Password length = %d, want 2", len(password.Value))
	}

	// Verify no duplicates even in short passwords
	if password.Value[0] == password.Value[1] {
		t.Errorf("Short password has duplicate character: %q", password.Value)
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

func TestPasswordGenerator_GenerateMultiplePasswordsAllUnique(t *testing.T) {
	generator := NewPasswordGenerator()

	// Larger batch to verify uniqueness guarantee holds
	config := entities.PasswordConfig{
		Length:         16,
		IncludeLower:   true,
		IncludeUpper:   true,
		IncludeNumbers: true,
		IncludeSymbols: true,
		Count:          20,
	}

	passwords, err := generator.GenerateMultiplePasswords(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	seen := make(map[string]bool)
	for _, password := range passwords {
		if seen[password.Value] {
			t.Errorf("Duplicate password in batch: %s", password.Value)
		}
		seen[password.Value] = true
	}
}

func TestPasswordGenerator_ExcludeSimilarRespected(t *testing.T) {
	generator := NewPasswordGenerator()

	config := entities.PasswordConfig{
		Length:         12,
		IncludeLower:   true,
		IncludeUpper:   true,
		IncludeNumbers: true,
		ExcludeSimilar: true,
		Count:          1,
	}

	for iter := 0; iter < 50; iter++ {
		password, err := generator.GeneratePassword(config)
		if err != nil {
			t.Fatalf("Iteration %d: unexpected error: %v", iter, err)
		}

		similar := "il1Lo0O"
		for _, char := range similar {
			if strings.ContainsRune(password.Value, char) {
				t.Fatalf("Iteration %d: password %q contains excluded similar char %q",
					iter, password.Value, string(char))
			}
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
