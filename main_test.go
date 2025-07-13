package main

import (
	"regexp"
	"strings"
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	tests := []struct {
		name   string
		config PasswordConfig
		want   struct {
			length     int
			hasLower   bool
			hasUpper   bool
			hasNumbers bool
			hasSymbols bool
		}
	}{
		{
			name: "default config",
			config: PasswordConfig{
				Length:         12,
				IncludeLower:   true,
				IncludeUpper:   true,
				IncludeNumbers: true,
				IncludeSymbols: false,
			},
			want: struct {
				length     int
				hasLower   bool
				hasUpper   bool
				hasNumbers bool
				hasSymbols bool
			}{
				length:     12,
				hasLower:   true,
				hasUpper:   true,
				hasNumbers: true,
				hasSymbols: false,
			},
		},
		{
			name: "secure config",
			config: PasswordConfig{
				Length:         16,
				IncludeLower:   true,
				IncludeUpper:   true,
				IncludeNumbers: true,
				IncludeSymbols: true,
			},
			want: struct {
				length     int
				hasLower   bool
				hasUpper   bool
				hasNumbers bool
				hasSymbols bool
			}{
				length:     16,
				hasLower:   true,
				hasUpper:   true,
				hasNumbers: true,
				hasSymbols: true,
			},
		},
		{
			name: "numbers only",
			config: PasswordConfig{
				Length:         6,
				IncludeLower:   false,
				IncludeUpper:   false,
				IncludeNumbers: true,
				IncludeSymbols: false,
			},
			want: struct {
				length     int
				hasLower   bool
				hasUpper   bool
				hasNumbers bool
				hasSymbols bool
			}{
				length:     6,
				hasLower:   false,
				hasUpper:   false,
				hasNumbers: true,
				hasSymbols: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate multiple passwords to test character inclusion probabilistically
			var hasLower, hasUpper, hasNumbers, hasSymbols bool

			for i := 0; i < 10; i++ {
				password, err := generatePassword(tt.config)
				if err != nil {
					t.Errorf("generatePassword() error = %v", err)
					return
				}

				// Check length
				if len(password) != tt.want.length {
					t.Errorf("generatePassword() length = %v, want %v", len(password), tt.want.length)
				}

				// Check character types
				if regexp.MustCompile(`[a-z]`).MatchString(password) {
					hasLower = true
				}
				if regexp.MustCompile(`[A-Z]`).MatchString(password) {
					hasUpper = true
				}
				if regexp.MustCompile(`[0-9]`).MatchString(password) {
					hasNumbers = true
				}
				if regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]`).MatchString(password) {
					hasSymbols = true
				}
			}

			// Check that required character types appeared at least once
			if tt.config.IncludeLower && !hasLower {
				t.Errorf("generatePassword() should include lowercase letters but didn't in 10 attempts")
			}
			if tt.config.IncludeUpper && !hasUpper {
				t.Errorf("generatePassword() should include uppercase letters but didn't in 10 attempts")
			}
			if tt.config.IncludeNumbers && !hasNumbers {
				t.Errorf("generatePassword() should include numbers but didn't in 10 attempts")
			}
			if tt.config.IncludeSymbols && !hasSymbols {
				t.Errorf("generatePassword() should include symbols but didn't in 10 attempts")
			}
		})
	}
}

func TestGeneratePasswordExcludeSimilar(t *testing.T) {
	config := PasswordConfig{
		Length:         100, // Large length to increase chances of hitting excluded chars
		IncludeLower:   true,
		IncludeUpper:   true,
		IncludeNumbers: true,
		IncludeSymbols: false,
		ExcludeSimilar: true,
	}

	password, err := generatePassword(config)
	if err != nil {
		t.Errorf("generatePassword() error = %v", err)
		return
	}

	// Check that similar characters are excluded
	similarChars := "il1Lo0O"
	for _, char := range similarChars {
		if strings.Contains(password, string(char)) {
			t.Errorf("generatePassword() should exclude similar character %c but found it in password", char)
		}
	}
}

func TestGeneratePasswordExcludeCustom(t *testing.T) {
	config := PasswordConfig{
		Length:         50,
		IncludeLower:   true,
		IncludeUpper:   true,
		IncludeNumbers: true,
		IncludeSymbols: false,
		ExcludeChars:   "aeiou",
	}

	password, err := generatePassword(config)
	if err != nil {
		t.Errorf("generatePassword() error = %v", err)
		return
	}

	// Check that excluded characters are not present
	excludedChars := "aeiou"
	for _, char := range excludedChars {
		if strings.Contains(password, string(char)) {
			t.Errorf("generatePassword() should exclude character %c but found it in password", char)
		}
	}
}

func TestGeneratePasswordError(t *testing.T) {
	// Test with no character sets selected
	config := PasswordConfig{
		Length:         12,
		IncludeLower:   false,
		IncludeUpper:   false,
		IncludeNumbers: false,
		IncludeSymbols: false,
	}

	_, err := generatePassword(config)
	if err == nil {
		t.Error("generatePassword() should return error when no character sets are selected")
	}
}

func TestCheckPasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantMin  int // Minimum expected score
		wantMax  int // Maximum expected score
	}{
		{
			name:     "very weak password",
			password: "123",
			wantMin:  0,
			wantMax:  2,
		},
		{
			name:     "weak password",
			password: "password",
			wantMin:  1,
			wantMax:  3,
		},
		{
			name:     "medium password",
			password: "Password123",
			wantMin:  3,
			wantMax:  5,
		},
		{
			name:     "strong password",
			password: "Password123!",
			wantMin:  5,
			wantMax:  7,
		},
		{
			name:     "very strong password",
			password: "MyVerySecurePassword123!@#",
			wantMin:  7,
			wantMax:  8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, score := checkPasswordStrength(tt.password)

			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("checkPasswordStrength() score = %v, want between %v and %v", score, tt.wantMin, tt.wantMax)
			}

			if result == "" {
				t.Error("checkPasswordStrength() should return non-empty result")
			}
		})
	}
}

func TestPasswordUniqueness(t *testing.T) {
	config := PasswordConfig{
		Length:         12,
		IncludeLower:   true,
		IncludeUpper:   true,
		IncludeNumbers: true,
		IncludeSymbols: true,
	}

	passwords := make(map[string]bool)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		password, err := generatePassword(config)
		if err != nil {
			t.Errorf("generatePassword() error = %v", err)
			return
		}

		if passwords[password] {
			t.Errorf("generatePassword() generated duplicate password: %s", password)
		}
		passwords[password] = true
	}
}
