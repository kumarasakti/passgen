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
	categories, err := cs.BuildCategories(config)
	if err != nil {
		return "", err
	}

	charset := strings.Join(categories, "")

	cs.charset = charset
	return charset, nil
}

// BuildCategories returns individual character categories after applying exclusions.
// Each enabled category (lowercase, uppercase, numbers, symbols) is returned as a
// separate string with similar and explicitly excluded characters removed.
func (cs *CharacterSet) BuildCategories(config PasswordConfig) ([]string, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	applyExclusions := func(s string) string {
		if config.ExcludeSimilar {
			similar := "il1Lo0O"
			for _, char := range similar {
				s = strings.ReplaceAll(s, string(char), "")
			}
		}
		if config.ExcludeChars != "" {
			for _, char := range config.ExcludeChars {
				s = strings.ReplaceAll(s, string(char), "")
			}
		}
		return s
	}

	var categories []string
	if config.IncludeLower {
		categories = append(categories, applyExclusions(Lowercase))
	}
	if config.IncludeUpper {
		categories = append(categories, applyExclusions(Uppercase))
	}
	if config.IncludeNumbers {
		categories = append(categories, applyExclusions(Numbers))
	}
	if config.IncludeSymbols {
		categories = append(categories, applyExclusions(Symbols))
	}

	// Filter out empty categories (all characters excluded)
	var nonEmpty []string
	for _, cat := range categories {
		if cat != "" {
			nonEmpty = append(nonEmpty, cat)
		}
	}

	if len(nonEmpty) == 0 {
		return nil, NewPasswordError("no characters available after exclusions")
	}

	return nonEmpty, nil
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
