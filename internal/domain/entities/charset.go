package entities

import "strings"

// Character set constants
const (
	DefaultLength = 14
	Lowercase     = "abcdefghijklmnopqrstuvwxyz"
	Uppercase     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers       = "0123456789"
	Symbols       = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

// CharacterSet manages character sets for password generation
type CharacterSet struct {
	charset string
}

// NewCharacterSet creates a new CharacterSet instance
func NewCharacterSet() *CharacterSet {
	return &CharacterSet{}
}

// BuildCharset builds a character set based on the provided configuration
func (cs *CharacterSet) BuildCharset(config PasswordConfig) (string, error) {
	if err := config.Validate(); err != nil {
		return "", err
	}

	charset := ""

	if config.IncludeLower {
		charset += Lowercase
	}
	if config.IncludeUpper {
		charset += Uppercase
	}
	if config.IncludeNumbers {
		charset += Numbers
	}
	if config.IncludeSymbols {
		charset += Symbols
	}

	if charset == "" {
		return "", NewPasswordError("no character sets selected")
	}

	// Remove similar characters if requested
	if config.ExcludeSimilar {
		similar := "il1Lo0O"
		for _, char := range similar {
			charset = strings.ReplaceAll(charset, string(char), "")
		}
	}

	// Remove excluded characters
	if config.ExcludeChars != "" {
		for _, char := range config.ExcludeChars {
			charset = strings.ReplaceAll(charset, string(char), "")
		}
	}

	if len(charset) == 0 {
		return "", NewPasswordError("no characters available after exclusions")
	}

	cs.charset = charset
	return charset, nil
}

// CalculateCharsetSize calculates the size of the character set for entropy calculation
func (cs *CharacterSet) CalculateCharsetSize(config PasswordConfig) int {
	size := 0

	if config.IncludeLower {
		size += 26
	}
	if config.IncludeUpper {
		size += 26
	}
	if config.IncludeNumbers {
		size += 10
	}
	if config.IncludeSymbols {
		size += len(Symbols)
	}

	// Adjust for excluded characters
	if config.ExcludeSimilar {
		similar := "il1Lo0O"
		for _, char := range similar {
			if config.IncludeLower && strings.Contains(Lowercase, string(char)) {
				size--
			}
			if config.IncludeUpper && strings.Contains(Uppercase, string(char)) {
				size--
			}
			if config.IncludeNumbers && strings.Contains(Numbers, string(char)) {
				size--
			}
		}
	}

	if config.ExcludeChars != "" {
		for _, char := range config.ExcludeChars {
			if config.IncludeLower && strings.Contains(Lowercase, string(char)) {
				size--
			}
			if config.IncludeUpper && strings.Contains(Uppercase, string(char)) {
				size--
			}
			if config.IncludeNumbers && strings.Contains(Numbers, string(char)) {
				size--
			}
			if config.IncludeSymbols && strings.Contains(Symbols, string(char)) {
				size--
			}
		}
	}

	return size
}

// GetCharset returns the current charset
func (cs *CharacterSet) GetCharset() string {
	return cs.charset
}
