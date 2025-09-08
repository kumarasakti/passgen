package entities

import (
	"regexp"
	"strings"
)

// PasswordStrength represents the strength level of a password
type PasswordStrength int

const (
	VeryWeak PasswordStrength = iota
	Weak
	Medium
	Strong
	VeryStrong
	ExtremelyStrong
)

// String returns the string representation of password strength
func (ps PasswordStrength) String() string {
	switch ps {
	case VeryWeak:
		return "Very Weak"
	case Weak:
		return "Weak"
	case Medium:
		return "Medium"
	case Strong:
		return "Strong"
	case VeryStrong:
		return "Very Strong"
	case ExtremelyStrong:
		return "Extremely Strong"
	default:
		return "Unknown"
	}
}

// PasswordConfig represents configuration for password generation
type PasswordConfig struct {
	Length         int
	IncludeLower   bool
	IncludeUpper   bool
	IncludeNumbers bool
	IncludeSymbols bool
	ExcludeSimilar bool
	ExcludeChars   string
	Count          int
}

// Validate ensures the password configuration is valid
func (pc PasswordConfig) Validate() error {
	if pc.Length <= 0 {
		return NewPasswordError("password length must be positive")
	}

	if !pc.IncludeLower && !pc.IncludeUpper && !pc.IncludeNumbers && !pc.IncludeSymbols {
		return NewPasswordError("at least one character type must be selected")
	}

	if pc.Count <= 0 {
		return NewPasswordError("password count must be positive")
	}

	return nil
}

// Password represents a generated password with its properties
type Password struct {
	Value  string
	Length int
}

// Constructs password entity with automatic length calculation
func NewPassword(value string) Password {
	return Password{
		Value:  value,
		Length: len(value),
	}
}

// Detects presence of lowercase alphabetic characters for complexity validation
func (p Password) HasLowercase() bool {
	matched, _ := regexp.MatchString(`[a-z]`, p.Value)
	return matched
}

// Detects presence of uppercase alphabetic characters for complexity validation
func (p Password) HasUppercase() bool {
	matched, _ := regexp.MatchString(`[A-Z]`, p.Value)
	return matched
}

// Detects presence of numeric characters for complexity validation
func (p Password) HasNumbers() bool {
	matched, _ := regexp.MatchString(`[0-9]`, p.Value)
	return matched
}

// Detects presence of special symbols for complexity validation
func (p Password) HasSymbols() bool {
	matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]`, p.Value)
	return matched
}

// Categorizes all character types present for complexity analysis
func (p Password) GetCharacterTypes() []string {
	var types []string

	if p.HasLowercase() {
		types = append(types, "Lowercase")
	}
	if p.HasUppercase() {
		types = append(types, "Uppercase")
	}
	if p.HasNumbers() {
		types = append(types, "Numbers")
	}
	if p.HasSymbols() {
		types = append(types, "Symbols")
	}

	return types
}

// Validates password contains actual content beyond whitespace
func (p Password) IsEmpty() bool {
	return strings.TrimSpace(p.Value) == ""
}

// PasswordError represents password-related errors
type PasswordError struct {
	Message string
}

// Constructs password-specific error with descriptive message
func NewPasswordError(message string) *PasswordError {
	return &PasswordError{Message: message}
}

// Error implements the error interface
func (pe *PasswordError) Error() string {
	return pe.Message
}
