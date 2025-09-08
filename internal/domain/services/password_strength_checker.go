package services

import (
	"fmt"

	"github.com/kumarasakti/passgen/internal/domain/entities"
)

// StrengthCheckResult represents the result of password strength checking
type StrengthCheckResult struct {
	Password          entities.Password
	Score             int
	MaxScore          int
	Strength          entities.PasswordStrength
	StrengthEmoji     string
	Celebration       string
	SarcasticComments []string
	Feedback          []string
	FormattedResult   string
}

// PasswordStrengthChecker provides sarcastic password strength checking
type PasswordStrengthChecker struct{}

// Evaluates password security with entertaining feedback and improvement suggestions
func NewPasswordStrengthChecker() *PasswordStrengthChecker {
	return &PasswordStrengthChecker{}
}

// CheckPasswordStrength analyzes password strength with sarcastic feedback
func (psc *PasswordStrengthChecker) CheckPasswordStrength(password entities.Password) StrengthCheckResult {
	score := 0
	maxScore := 8
	feedback := []string{}
	sarcasticComments := []string{}

	// Length check
	if password.Length >= 12 {
		score += 2
		if password.Length >= 16 {
			sarcasticComments = append(sarcasticComments, "Wow, someone actually read the security guidelines! 👏")
		}
	} else if password.Length >= 8 {
		score += 1
		sarcasticComments = append(sarcasticComments, "8 characters? How... minimalistic of you 🤔")
	} else {
		feedback = append(feedback, "Password should be at least 8 characters long")
		sarcasticComments = append(sarcasticComments, "Really? That's shorter than most people's names! 😅")
	}

	// Character variety checks
	if password.HasLowercase() {
		score += 1
	} else {
		feedback = append(feedback, "Add lowercase letters")
		sarcasticComments = append(sarcasticComments, "No lowercase letters? Are we SHOUTING all the time? 📢")
	}

	if password.HasUppercase() {
		score += 1
	} else {
		feedback = append(feedback, "Add uppercase letters")
		sarcasticComments = append(sarcasticComments, "No capitals? I guess we're going for the e.e. cummings aesthetic 🎭")
	}

	if password.HasNumbers() {
		score += 1
	} else {
		feedback = append(feedback, "Add numbers")
		sarcasticComments = append(sarcasticComments, "Numbers are optional now? Math teachers everywhere are crying 😢")
	}

	if password.HasSymbols() {
		score += 2
	} else {
		feedback = append(feedback, "Add special characters")
		sarcasticComments = append(sarcasticComments, "No symbols? Your password is as plain as unseasoned chicken 🐔")
	}

	// Bonus for length
	if password.Length >= 16 {
		score += 1
	}

	// Determine strength and celebration
	strength, strengthEmoji, celebration := psc.determineStrengthFromScore(score)

	// Format the result
	formattedResult := psc.formatResult(password, score, maxScore, strength, strengthEmoji, celebration, sarcasticComments, feedback)

	return StrengthCheckResult{
		Password:          password,
		Score:             score,
		MaxScore:          maxScore,
		Strength:          strength,
		StrengthEmoji:     strengthEmoji,
		Celebration:       celebration,
		SarcasticComments: sarcasticComments,
		Feedback:          feedback,
		FormattedResult:   formattedResult,
	}
}

// determineStrengthFromScore determines strength based on score
func (psc *PasswordStrengthChecker) determineStrengthFromScore(score int) (entities.PasswordStrength, string, string) {
	var strength entities.PasswordStrength
	var strengthEmoji, celebration string

	switch {
	case score >= 7:
		strength = entities.VeryStrong
		strengthEmoji = "🔥"
		celebration = "Impressive! Your password could probably withstand a zombie apocalypse! 🧟‍♂️"
	case score >= 5:
		strength = entities.Strong
		strengthEmoji = "💪"
		celebration = "Not bad! Your password has some real backbone! 🦴"
	case score >= 3:
		strength = entities.Medium
		strengthEmoji = "😐"
		celebration = "It's... adequate. Like a participation trophy for password security 🏆"
	case score >= 1:
		strength = entities.Weak
		strengthEmoji = "😰"
		celebration = "Yikes! This password couldn't protect a diary from a nosy sibling! 📖"
	default:
		strength = entities.VeryWeak
		strengthEmoji = "🚨"
		celebration = "Oh dear... this password is weaker than my WiFi signal in the basement! 📶"
	}

	return strength, strengthEmoji, celebration
}

// formatResult formats the strength check result into a string
func (psc *PasswordStrengthChecker) formatResult(password entities.Password, score, maxScore int, strength entities.PasswordStrength, strengthEmoji, celebration string, sarcasticComments, feedback []string) string {
	result := "🔍 Password Analysis Results:\n"
	result += fmt.Sprintf("Strength: %s %s (Score: %d/%d)\n", strength.String(), strengthEmoji, score, maxScore)
	result += fmt.Sprintf("\n%s\n", celebration)

	if len(sarcasticComments) > 0 {
		result += "\n💭 Honest Feedback:\n"
		for _, comment := range sarcasticComments {
			result += fmt.Sprintf("• %s\n", comment)
		}
	}

	if len(feedback) > 0 {
		result += "\n💡 Actionable Suggestions:\n"
		for _, suggestion := range feedback {
			result += fmt.Sprintf("• %s\n", suggestion)
		}
		result += "\nPro tip: Try 'passgen --secure -l 16' for a password that actually means business! 🚀"
	}

	return result
}
