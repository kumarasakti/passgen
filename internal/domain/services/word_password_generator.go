package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/kumarasakti/passgen/internal/domain/entities"
)

// WordPasswordGenerator handles word-based password generation
type WordPasswordGenerator struct {
	analyzer *PasswordAnalyzer
}

// NewWordPasswordGenerator creates a new WordPasswordGenerator instance
func NewWordPasswordGenerator(analyzer *PasswordAnalyzer) *WordPasswordGenerator {
	return &WordPasswordGenerator{
		analyzer: analyzer,
	}
}

// GenerateWordPassword generates a password based on a word pattern
func (wpg *WordPasswordGenerator) GenerateWordPassword(pattern *entities.WordPattern) (string, error) {
	if err := pattern.Validate(); err != nil {
		return "", err
	}

	basePassword := pattern.GetTransformedWord()
	
	// Add additional randomization if not preserving length
	if !pattern.PreserveLength {
		basePassword = wpg.addRandomElements(basePassword, pattern)
	}

	return basePassword, nil
}

// GenerateMultipleWordPasswords generates multiple word-based passwords
func (wpg *WordPasswordGenerator) GenerateMultipleWordPasswords(pattern *entities.WordPattern, count int) ([]string, error) {
	if count <= 0 {
		return nil, entities.NewPasswordError("count must be greater than 0")
	}

	if count > 100 {
		return nil, entities.NewPasswordError("count cannot exceed 100")
	}

	passwords := make([]string, 0, count)
	
	for i := 0; i < count; i++ {
		// Create slight variations for each password
		variantPattern := wpg.createVariantPattern(pattern, i)
		
		password, err := wpg.GenerateWordPassword(variantPattern)
		if err != nil {
			return nil, fmt.Errorf("failed to generate password %d: %w", i+1, err)
		}
		
		passwords = append(passwords, password)
	}

	return passwords, nil
}

// addRandomElements adds random numbers, symbols, or years to enhance security
func (wpg *WordPasswordGenerator) addRandomElements(base string, pattern *entities.WordPattern) string {
	switch pattern.Complexity {
	case entities.ComplexityLow:
		return wpg.addSimpleElements(base)
	case entities.ComplexityMedium:
		return wpg.addMediumElements(base)
	case entities.ComplexityHigh:
		return wpg.addComplexElements(base)
	default:
		return base
	}
}

// addSimpleElements adds basic random elements
func (wpg *WordPasswordGenerator) addSimpleElements(base string) string {
	// Add a random number 1-9
	num, _ := rand.Int(rand.Reader, big.NewInt(9))
	return base + strconv.FormatInt(num.Int64()+1, 10)
}

// addMediumElements adds medium complexity random elements
func (wpg *WordPasswordGenerator) addMediumElements(base string) string {
	// Add random 2-digit number and a symbol
	num, _ := rand.Int(rand.Reader, big.NewInt(90))
	symbols := []string{"!", "@", "#", "$", "%"}
	symbolIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(symbols))))
	
	return base + strconv.FormatInt(num.Int64()+10, 10) + symbols[symbolIdx.Int64()]
}

// addComplexElements adds high complexity random elements
func (wpg *WordPasswordGenerator) addComplexElements(base string) string {
	// Add current year, random number, and multiple symbols
	currentYear := time.Now().Year()
	num, _ := rand.Int(rand.Reader, big.NewInt(100))
	symbols := []string{"!", "@", "#", "$", "%", "^", "&", "*"}
	symbolIdx1, _ := rand.Int(rand.Reader, big.NewInt(int64(len(symbols))))
	symbolIdx2, _ := rand.Int(rand.Reader, big.NewInt(int64(len(symbols))))
	
	return fmt.Sprintf("%s%d%s%02d%s", 
		base, 
		currentYear%100, // Last 2 digits of year
		symbols[symbolIdx1.Int64()],
		num.Int64(),
		symbols[symbolIdx2.Int64()],
	)
}

// createVariantPattern creates a slight variation of the pattern for multiple password generation
func (wpg *WordPasswordGenerator) createVariantPattern(original *entities.WordPattern, index int) *entities.WordPattern {
	variant := *original // Copy the original pattern
	
	// Create variations based on index
	switch index % 4 {
	case 0:
		// Use original strategy
	case 1:
		// Switch to different strategy if hybrid
		if original.Strategy == entities.StrategyHybrid {
			variant.Strategy = entities.StrategyLeetspeak
		}
	case 2:
		// Increase complexity for this variant
		if original.Complexity == entities.ComplexityLow {
			variant.Complexity = entities.ComplexityMedium
		} else if original.Complexity == entities.ComplexityMedium {
			variant.Complexity = entities.ComplexityHigh
		}
	case 3:
		// Use mixed case strategy
		if original.Strategy == entities.StrategyHybrid {
			variant.Strategy = entities.StrategyMixedCase
		}
	}
	
	return &variant
}

// AnalyzeWordPassword analyzes a word-based password and provides insights
func (wpg *WordPasswordGenerator) AnalyzeWordPassword(password, originalWord string) (*PasswordAnalysis, error) {
	// Create password entity and basic config for analysis
	passwordEntity := entities.NewPassword(password)
	
	// Create a basic config based on password characteristics
	config := entities.PasswordConfig{
		Length:         len(password),
		IncludeLower:   passwordEntity.HasLowercase(),
		IncludeUpper:   passwordEntity.HasUppercase(),
		IncludeNumbers: passwordEntity.HasNumbers(),
		IncludeSymbols: passwordEntity.HasSymbols(),
		Count:          1,
	}
	
	// Use the existing password analyzer
	analysis := wpg.analyzer.AnalyzePassword(passwordEntity, config)
	
	// Add word-specific insights
	analysis.WordBased = true
	analysis.OriginalWord = originalWord
	analysis.TransformationQuality = wpg.assessTransformationQuality(password, originalWord)
	
	return &analysis, nil
}

// assessTransformationQuality evaluates how well the word was transformed
func (wpg *WordPasswordGenerator) assessTransformationQuality(password, originalWord string) string {
	password = strings.ToLower(password)
	originalWord = strings.ToLower(originalWord)
	
	// Check if original word is still clearly visible
	if strings.Contains(password, originalWord) {
		baseWordIndex := strings.Index(password, originalWord)
		beforeWord := password[:baseWordIndex]
		afterWord := password[baseWordIndex+len(originalWord):]
		
		hasPrefix := len(beforeWord) > 0
		hasSuffix := len(afterWord) > 0
		hasNumbers := strings.ContainsAny(password, "0123456789")
		hasSymbols := strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?")
		hasMixedCase := password != strings.ToLower(password)
		
		transformationCount := 0
		if hasPrefix { transformationCount++ }
		if hasSuffix { transformationCount++ }
		if hasNumbers { transformationCount++ }
		if hasSymbols { transformationCount++ }
		if hasMixedCase { transformationCount++ }
		
		switch {
		case transformationCount >= 4:
			return "Excellent transformation - highly secure while memorable"
		case transformationCount >= 3:
			return "Good transformation - secure and readable"
		case transformationCount >= 2:
			return "Moderate transformation - could be more complex"
		default:
			return "Basic transformation - consider adding more complexity"
		}
	}
	
	// If original word is heavily modified (leetspeak, etc.)
	return "Advanced transformation - original word well-disguised"
}

// GetWordStrategySuggestions provides suggestions for improving word-based passwords
func (wpg *WordPasswordGenerator) GetWordStrategySuggestions(word string, currentStrategy entities.TransformationStrategy) []string {
	suggestions := []string{}
	
	switch currentStrategy {
	case entities.StrategyLeetspeak:
		suggestions = append(suggestions, 
			"Try mixed-case strategy for better readability",
			"Add suffix strategy for additional security",
			"Consider hybrid approach for best of both worlds")
	case entities.StrategyMixedCase:
		suggestions = append(suggestions,
			"Add leetspeak for stronger character variety",
			"Include numbers and symbols with suffix strategy",
			"Try hybrid for maximum security")
	case entities.StrategySuffix:
		suggestions = append(suggestions,
			"Combine with leetspeak for character substitution",
			"Add mixed-case for visual variety",
			"Try insert strategy for more subtle changes")
	case entities.StrategyPrefix:
		suggestions = append(suggestions,
			"Combine with suffix for balanced approach",
			"Add character transformations with leetspeak",
			"Try hybrid for comprehensive transformation")
	case entities.StrategyInsert:
		suggestions = append(suggestions,
			"Add more complexity with hybrid strategy",
			"Include leetspeak for character variety",
			"Try suffix for additional length")
	case entities.StrategyHybrid:
		suggestions = append(suggestions,
			"Excellent choice! Try increasing complexity level",
			"Generate multiple variations for different uses",
			"Consider longer words for even better security")
	}
	
	// Add word-specific suggestions
	if len(word) < 6 {
		suggestions = append(suggestions, "Consider using longer words (6+ characters) for better security")
	}
	
	if len(word) > 12 {
		suggestions = append(suggestions, "Long word detected - insert strategy works well with longer words")
	}
	
	return suggestions
}
