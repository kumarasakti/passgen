package entities

import (
	"regexp"
	"strings"
	"unicode"
)

// TransformationStrategy defines different word transformation approaches
type TransformationStrategy string

const (
	StrategyLeetspeak  TransformationStrategy = "leetspeak"
	StrategyMixedCase  TransformationStrategy = "mixed-case"
	StrategySuffix     TransformationStrategy = "suffix"
	StrategyPrefix     TransformationStrategy = "prefix"
	StrategyInsert     TransformationStrategy = "insert"
	StrategyHybrid     TransformationStrategy = "hybrid"
)

// ComplexityLevel defines the complexity of transformations
type ComplexityLevel string

const (
	ComplexityLow    ComplexityLevel = "low"
	ComplexityMedium ComplexityLevel = "medium"
	ComplexityHigh   ComplexityLevel = "high"
)

// WordPattern manages word-based password transformation rules
type WordPattern struct {
	Word               string
	Strategy           TransformationStrategy
	Complexity         ComplexityLevel
	PreserveLength     bool
	LeetSpeakMappings  map[rune]rune
	NumberSuffixLength int
	SymbolCount        int
}

// NewWordPattern creates a new WordPattern with default settings
func NewWordPattern(word string) *WordPattern {
	return &WordPattern{
		Word:               strings.ToLower(strings.TrimSpace(word)),
		Strategy:           StrategyHybrid,
		Complexity:         ComplexityMedium,
		PreserveLength:     false,
		LeetSpeakMappings:  getDefaultLeetSpeakMappings(),
		NumberSuffixLength: 2,
		SymbolCount:        1,
	}
}

// getDefaultLeetSpeakMappings returns common leetspeak character substitutions
func getDefaultLeetSpeakMappings() map[rune]rune {
	return map[rune]rune{
		'a': '@',
		'e': '3',
		'i': '1',
		'o': '0',
		's': '$',
		't': '7',
		'l': '1',
		'g': '9',
		'b': '6',
		'z': '2',
	}
}

// Validate checks if the word is valid for password generation
func (wp *WordPattern) Validate() error {
	if wp.Word == "" {
		return NewPasswordError("word cannot be empty")
	}
	
	if len(wp.Word) < 3 {
		return NewPasswordError("word must be at least 3 characters long")
	}
	
	if len(wp.Word) > 50 {
		return NewPasswordError("word must be less than 50 characters long")
	}
	
	// Check if word contains only letters (basic validation)
	if !regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(wp.Word) {
		return NewPasswordError("word must contain only letters")
	}
	
	return nil
}

// SetStrategy sets the transformation strategy
func (wp *WordPattern) SetStrategy(strategy TransformationStrategy) *WordPattern {
	wp.Strategy = strategy
	return wp
}

// SetComplexity sets the complexity level
func (wp *WordPattern) SetComplexity(complexity ComplexityLevel) *WordPattern {
	wp.Complexity = complexity
	
	// Adjust parameters based on complexity
	switch complexity {
	case ComplexityLow:
		wp.NumberSuffixLength = 1
		wp.SymbolCount = 1
	case ComplexityMedium:
		wp.NumberSuffixLength = 2
		wp.SymbolCount = 1
	case ComplexityHigh:
		wp.NumberSuffixLength = 3
		wp.SymbolCount = 2
	}
	
	return wp
}

// SetPreserveLength sets whether to preserve the original word length
func (wp *WordPattern) SetPreserveLength(preserve bool) *WordPattern {
	wp.PreserveLength = preserve
	return wp
}

// GetTransformedWord returns the transformed word based on the pattern settings
func (wp *WordPattern) GetTransformedWord() string {
	switch wp.Strategy {
	case StrategyLeetspeak:
		return wp.applyLeetspeak()
	case StrategyMixedCase:
		return wp.applyMixedCase()
	case StrategySuffix:
		return wp.applySuffix()
	case StrategyPrefix:
		return wp.applyPrefix()
	case StrategyInsert:
		return wp.applyInsert()
	case StrategyHybrid:
		return wp.applyHybrid()
	default:
		return wp.Word
	}
}

// applyLeetspeak applies leetspeak transformations
func (wp *WordPattern) applyLeetspeak() string {
	result := make([]rune, 0, len(wp.Word))
	
	for _, char := range wp.Word {
		if replacement, exists := wp.LeetSpeakMappings[unicode.ToLower(char)]; exists {
			// Apply leetspeak based on complexity
			switch wp.Complexity {
			case ComplexityLow:
				// Only replace some characters
				if char == 'e' || char == 'a' {
					result = append(result, replacement)
				} else {
					result = append(result, char)
				}
			case ComplexityMedium:
				// Replace most common leetspeak characters
				if char == 'e' || char == 'a' || char == 'i' || char == 'o' || char == 's' {
					result = append(result, replacement)
				} else {
					result = append(result, char)
				}
			case ComplexityHigh:
				// Replace all available leetspeak characters
				result = append(result, replacement)
			}
		} else {
			result = append(result, char)
		}
	}
	
	return string(result)
}

// applyMixedCase applies mixed case transformations
func (wp *WordPattern) applyMixedCase() string {
	result := make([]rune, 0, len(wp.Word))
	
	for i, char := range wp.Word {
		switch wp.Complexity {
		case ComplexityLow:
			// Just capitalize first letter
			if i == 0 {
				result = append(result, unicode.ToUpper(char))
			} else {
				result = append(result, unicode.ToLower(char))
			}
		case ComplexityMedium:
			// Capitalize first and middle
			if i == 0 || i == len(wp.Word)/2 {
				result = append(result, unicode.ToUpper(char))
			} else {
				result = append(result, unicode.ToLower(char))
			}
		case ComplexityHigh:
			// Alternating case
			if i%2 == 0 {
				result = append(result, unicode.ToUpper(char))
			} else {
				result = append(result, unicode.ToLower(char))
			}
		}
	}
	
	return string(result)
}

// applySuffix adds suffix based on complexity
func (wp *WordPattern) applySuffix() string {
	base := strings.Title(wp.Word)
	
	switch wp.Complexity {
	case ComplexityLow:
		return base + "1"
	case ComplexityMedium:
		return base + "123"
	case ComplexityHigh:
		return base + "123!"
	}
	
	return base
}

// applyPrefix adds prefix based on complexity
func (wp *WordPattern) applyPrefix() string {
	base := strings.Title(wp.Word)
	
	switch wp.Complexity {
	case ComplexityLow:
		return "!" + base
	case ComplexityMedium:
		return "@" + base + "1"
	case ComplexityHigh:
		return "#!" + base + "42"
	}
	
	return base
}

// applyInsert inserts characters within the word
func (wp *WordPattern) applyInsert() string {
	if len(wp.Word) < 4 {
		return wp.applySuffix() // Fallback for short words
	}
	
	result := make([]rune, 0, len(wp.Word)+3)
	wordRunes := []rune(wp.Word)
	
	// Insert character in the middle
	middle := len(wordRunes) / 2
	
	for i, char := range wordRunes {
		if i == 0 {
			result = append(result, unicode.ToUpper(char))
		} else if i == middle {
			switch wp.Complexity {
			case ComplexityLow:
				result = append(result, char, '1')
			case ComplexityMedium:
				result = append(result, char, '1', '!')
			case ComplexityHigh:
				result = append(result, char, '@', '3')
			}
		} else {
			result = append(result, char)
		}
	}
	
	return string(result)
}

// applyHybrid combines multiple strategies
func (wp *WordPattern) applyHybrid() string {
	// Start with mixed case
	result := wp.applyMixedCase()
	
	// Apply some leetspeak
	tempPattern := &WordPattern{
		Word:              result,
		Complexity:        wp.Complexity,
		LeetSpeakMappings: wp.LeetSpeakMappings,
	}
	result = tempPattern.applyLeetspeak()
	
	// Add suffix based on complexity
	switch wp.Complexity {
	case ComplexityLow:
		result += "!"
	case ComplexityMedium:
		result += "!23"
	case ComplexityHigh:
		result += "@42!"
	}
	
	return result
}
