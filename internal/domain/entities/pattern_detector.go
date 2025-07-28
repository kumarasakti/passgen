package entities

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// PasswordPattern represents a detected pattern in a password
type PasswordPattern struct {
	Type        string
	Description string
	Severity    string // "high", "medium", "low"
	Suggestion  string
}

// PasswordPatternDetector detects common patterns in passwords
type PasswordPatternDetector struct{}

// NewPasswordPatternDetector creates a new password pattern detector
func NewPasswordPatternDetector() *PasswordPatternDetector {
	return &PasswordPatternDetector{}
}

// DetectPatterns analyzes a password for common patterns
func (ppd *PasswordPatternDetector) DetectPatterns(password string) []PasswordPattern {
	var patterns []PasswordPattern

	// Check for common patterns
	patterns = append(patterns, ppd.checkSequential(password)...)
	patterns = append(patterns, ppd.checkRepeating(password)...)
	patterns = append(patterns, ppd.checkCommonWords(password)...)
	patterns = append(patterns, ppd.checkKeyboardPatterns(password)...)
	patterns = append(patterns, ppd.checkDatePatterns(password)...)
	patterns = append(patterns, ppd.checkNumberPatterns(password)...)

	return patterns
}

// checkSequential detects sequential characters (abc, 123, etc.)
func (ppd *PasswordPatternDetector) checkSequential(password string) []PasswordPattern {
	var patterns []PasswordPattern

	// Check for sequential letters
	if matched, _ := regexp.MatchString(`(?i)[a-z]{3,}`, password); matched {
		// Check if they are actually sequential
		lower := strings.ToLower(password)
		for i := 0; i < len(lower)-2; i++ {
			if lower[i+1] == lower[i]+1 && lower[i+2] == lower[i]+2 {
				patterns = append(patterns, PasswordPattern{
					Type:        "sequential_letters",
					Description: "Contains sequential letters (e.g., abc, def)",
					Severity:    "medium",
					Suggestion:  "Avoid using sequential letters in passwords",
				})
				break
			}
		}
	}

	// Check for sequential numbers
	if matched, _ := regexp.MatchString(`\d{3,}`, password); matched {
		for i := 0; i < len(password)-2; i++ {
			if unicode.IsDigit(rune(password[i])) && unicode.IsDigit(rune(password[i+1])) && unicode.IsDigit(rune(password[i+2])) {
				if password[i+1] == password[i]+1 && password[i+2] == password[i]+2 {
					patterns = append(patterns, PasswordPattern{
						Type:        "sequential_numbers",
						Description: "Contains sequential numbers (e.g., 123, 456)",
						Severity:    "medium",
						Suggestion:  "Avoid using sequential numbers in passwords",
					})
					break
				}
			}
		}
	}

	return patterns
}

// checkRepeating detects repeating characters
func (ppd *PasswordPatternDetector) checkRepeating(password string) []PasswordPattern {
	var patterns []PasswordPattern

	// Check for 3+ repeating characters
	if matched, _ := regexp.MatchString(`(.)\\1{2,}`, password); matched {
		patterns = append(patterns, PasswordPattern{
			Type:        "repeating_chars",
			Description: "Contains repeating characters (e.g., aaa, 111)",
			Severity:    "high",
			Suggestion:  "Avoid repeating the same character multiple times",
		})
	}

	return patterns
}

// checkCommonWords detects common words and patterns
func (ppd *PasswordPatternDetector) checkCommonWords(password string) []PasswordPattern {
	var patterns []PasswordPattern
	lower := strings.ToLower(password)

	commonWords := []string{
		"password", "admin", "user", "login", "welcome", "qwerty",
		"letmein", "monkey", "dragon", "master", "shadow", "123456",
		"password123", "admin123", "root", "toor", "guest",
	}

	for _, word := range commonWords {
		if strings.Contains(lower, word) {
			patterns = append(patterns, PasswordPattern{
				Type:        "common_word",
				Description: fmt.Sprintf("Contains common word: '%s'", word),
				Severity:    "high",
				Suggestion:  "Avoid using common words in passwords",
			})
		}
	}

	return patterns
}

// checkKeyboardPatterns detects keyboard patterns
func (ppd *PasswordPatternDetector) checkKeyboardPatterns(password string) []PasswordPattern {
	var patterns []PasswordPattern
	lower := strings.ToLower(password)

	keyboardPatterns := []string{
		"qwerty", "asdf", "zxcv", "1234", "qwertyuiop",
		"asdfghjkl", "zxcvbnm", "!@#$", "!@#$%^&*",
	}

	for _, pattern := range keyboardPatterns {
		if strings.Contains(lower, pattern) {
			patterns = append(patterns, PasswordPattern{
				Type:        "keyboard_pattern",
				Description: fmt.Sprintf("Contains keyboard pattern: '%s'", pattern),
				Severity:    "medium",
				Suggestion:  "Avoid using keyboard patterns in passwords",
			})
		}
	}

	return patterns
}

// checkDatePatterns detects date-like patterns
func (ppd *PasswordPatternDetector) checkDatePatterns(password string) []PasswordPattern {
	var patterns []PasswordPattern

	datePatterns := []string{
		`\d{2}/\d{2}/\d{4}`, // MM/DD/YYYY
		`\d{2}-\d{2}-\d{4}`, // MM-DD-YYYY
		`\d{4}/\d{2}/\d{2}`, // YYYY/MM/DD
		`\d{4}-\d{2}-\d{2}`, // YYYY-MM-DD
		`\d{2}\d{2}\d{4}`,   // MMDDYYYY
		`\d{4}\d{2}\d{2}`,   // YYYYMMDD
	}

	for _, pattern := range datePatterns {
		if matched, _ := regexp.MatchString(pattern, password); matched {
			patterns = append(patterns, PasswordPattern{
				Type:        "date_pattern",
				Description: "Contains date-like pattern",
				Severity:    "medium",
				Suggestion:  "Avoid using dates in passwords",
			})
			break
		}
	}

	return patterns
}

// checkNumberPatterns detects simple number patterns
func (ppd *PasswordPatternDetector) checkNumberPatterns(password string) []PasswordPattern {
	var patterns []PasswordPattern

	// Check for all numbers at the end (common pattern)
	if matched, _ := regexp.MatchString(`[a-zA-Z]+\d+$`, password); matched {
		patterns = append(patterns, PasswordPattern{
			Type:        "numbers_at_end",
			Description: "Numbers only at the end of password",
			Severity:    "low",
			Suggestion:  "Consider mixing numbers throughout the password",
		})
	}

	// Check for simple year patterns
	if matched, _ := regexp.MatchString(`(19|20)\d{2}`, password); matched {
		patterns = append(patterns, PasswordPattern{
			Type:        "year_pattern",
			Description: "Contains year-like pattern",
			Severity:    "medium",
			Suggestion:  "Avoid using years in passwords",
		})
	}

	return patterns
}
