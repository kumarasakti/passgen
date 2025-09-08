package services

import (
	"fmt"
	"math"

	"github.com/kumarasakti/passgen/internal/domain/entities"
)

// PasswordAnalysis represents the result of password analysis
type PasswordAnalysis struct {
	Password       entities.Password
	CharsetSize    int
	CharacterTypes []string
	Entropy        float64
	Strength       entities.PasswordStrength
	StrengthEmoji  string
	TimeToCrack    string
	SecurityLevel  string
	Tips           []string
	Celebration    string
	// Word-based password specific fields
	WordBased             bool
	OriginalWord          string
	TransformationQuality string
}

// Provides comprehensive password security evaluation and strength assessment
type PasswordAnalyzer struct {
	charsetManager *entities.CharacterSet
}

// Initializes password security analysis with character set management
func NewPasswordAnalyzer() *PasswordAnalyzer {
	return &PasswordAnalyzer{
		charsetManager: entities.NewCharacterSet(),
	}
}

// Delivers detailed security metrics including entropy, strength rating, and improvement recommendations
func (pa *PasswordAnalyzer) AnalyzePassword(password entities.Password, config entities.PasswordConfig) PasswordAnalysis {
	charsetSize := pa.charsetManager.CalculateCharsetSize(config)
	characterTypes := password.GetCharacterTypes()

	// Calculate entropy: log2(charset^length)
	entropy := float64(password.Length) * math.Log2(float64(charsetSize))

	// Determine strength and related properties
	strength, strengthEmoji, securityLevel, celebration, tips := pa.determineStrength(entropy, password.Length, len(characterTypes))

	// Calculate time to crack
	timeToCrack := pa.calculateTimeToCrack(charsetSize, password.Length)

	return PasswordAnalysis{
		Password:       password,
		CharsetSize:    charsetSize,
		CharacterTypes: characterTypes,
		Entropy:        entropy,
		Strength:       strength,
		StrengthEmoji:  strengthEmoji,
		TimeToCrack:    timeToCrack,
		SecurityLevel:  securityLevel,
		Tips:           tips,
		Celebration:    celebration,
	}
}

// determineStrength determines password strength based on entropy and other factors
func (pa *PasswordAnalyzer) determineStrength(entropy float64, length, charTypeCount int) (entities.PasswordStrength, string, string, string, []string) {
	var strength entities.PasswordStrength
	var strengthEmoji, securityLevel, celebration string
	var tips []string

	switch {
	case entropy >= 100:
		strength = entities.ExtremelyStrong
		strengthEmoji = "🔥"
		securityLevel = "Quantum-resistant for the foreseeable future!"
		celebration = "Brr, that's ice cold security! Even hackers are shivering! 🥶"
	case entropy >= 80:
		strength = entities.VeryStrong
		strengthEmoji = "💪"
		securityLevel = "Exceeds security standards for high-value accounts"
		celebration = "Someone's taking this security thing seriously! 🌟"
	case entropy >= 60:
		strength = entities.Strong
		strengthEmoji = "💯"
		securityLevel = "Great for securing important accounts"
		celebration = "Not bad, you actually read the security guidelines! 🎯"
	case entropy >= 40:
		strength = entities.Medium
		strengthEmoji = "⚡"
		securityLevel = "Adequate for most general purposes"
		celebration = "Well, it's... adequate. I guess that's something! 👍"
		if length < 12 {
			tips = append(tips, "Consider using 12+ characters for better security")
		}
		if charTypeCount < 3 {
			tips = append(tips, "Add more character types (symbols, numbers) for stronger security")
		}
	case entropy >= 25:
		strength = entities.Weak
		strengthEmoji = "😰"
		securityLevel = "Suitable only for low-security uses"
		celebration = "Oh honey, we need to talk about your password choices... 💪"
		tips = append(tips, "Use at least 12 characters")
		tips = append(tips, "Include uppercase, lowercase, numbers, and symbols")
		tips = append(tips, "Try `passgen --secure` for maximum protection!")
	default:
		strength = entities.VeryWeak
		strengthEmoji = "🚨"
		securityLevel = "Not recommended for any security purposes"
		celebration = "Yikes! Even my grandma would crack this in her sleep! 🚀"
		tips = append(tips, "Use at least 12 characters")
		tips = append(tips, "Include multiple character types")
		tips = append(tips, "Try `passgen --secure -l 16` for excellent security!")
	}

	return strength, strengthEmoji, securityLevel, celebration, tips
}

// calculateTimeToCrack calculates time to crack the password
func (pa *PasswordAnalyzer) calculateTimeToCrack(charsetSize, length int) string {
	// Assuming 1 trillion guesses per second
	guessesPerSecond := 1e12
	possibleCombinations := math.Pow(float64(charsetSize), float64(length))
	secondsToCrack := possibleCombinations / (2 * guessesPerSecond) // Average case

	if secondsToCrack < 60 {
		return "Less than a minute"
	} else if secondsToCrack < 3600 {
		return fmt.Sprintf("%.1f minutes", secondsToCrack/60)
	} else if secondsToCrack < 86400 {
		return fmt.Sprintf("%.1f hours", secondsToCrack/3600)
	} else if secondsToCrack < 31536000 {
		return fmt.Sprintf("%.1f days", secondsToCrack/86400)
	} else if secondsToCrack < 31536000000 {
		return fmt.Sprintf("%.1f years", secondsToCrack/31536000)
	} else {
		// For very large numbers, use scientific notation
		years := secondsToCrack / 31536000
		if years > 1e15 {
			return fmt.Sprintf("%.1e years", years)
		} else {
			return fmt.Sprintf("%.0f years", years)
		}
	}
}
