package cli

import (
	"fmt"
	"strings"

	"github.com/kumarasakti/passgen/internal/application"
	"github.com/kumarasakti/passgen/internal/domain/services"
)

// Formatter handles output formatting for the CLI
type Formatter struct{}

// NewFormatter creates a new Formatter instance
func NewFormatter() *Formatter {
	return &Formatter{}
}

// FormatPasswordGeneration formats password generation results for display
func (f *Formatter) FormatPasswordGeneration(analyses []services.PasswordAnalysis, excludeSimilar bool) string {
	var output strings.Builder

	for i, analysis := range analyses {
		// Make password super prominent and easy to read
		if len(analyses) > 1 {
			output.WriteString(fmt.Sprintf("🎯 Password %d:\n", i+1))
		} else {
			output.WriteString("🎯 Your Password:\n")
		}

		// Create a box around the password for maximum visibility
		password := analysis.Password.Value
		output.WriteString("┌" + strings.Repeat("─", len(password)+2) + "┐\n")
		output.WriteString(fmt.Sprintf("│ %s │\n", password))
		output.WriteString("└" + strings.Repeat("─", len(password)+2) + "┘\n\n")

		// Brief one-line summary
		output.WriteString(fmt.Sprintf("📊 Length: %d | Character types: %s | Strength: %s %s\n",
			analysis.Password.Length,
			strings.Join(analysis.CharacterTypes, ", "),
			analysis.Strength.String(),
			analysis.StrengthEmoji))

		// Optional: Show detailed analysis only if single password
		if len(analyses) == 1 {
			output.WriteString(fmt.Sprintf("\n🔒 Security info: %.1f bits entropy, cracks in %s\n",
				analysis.Entropy, analysis.TimeToCrack))

			// Tips if password is weak
			if len(analysis.Tips) > 0 {
				output.WriteString("\n💡 Suggestions:\n")
				for _, tip := range analysis.Tips {
					output.WriteString(fmt.Sprintf("   • %s\n", tip))
				}
			}

			// Add the sarcastic comment for fun
			if analysis.Celebration != "" {
				output.WriteString(fmt.Sprintf("\n💬 %s\n", analysis.Celebration))
			}
		}

		// Add separator for multiple passwords
		if i < len(analyses)-1 {
			output.WriteString("\n" + strings.Repeat("─", 60) + "\n\n")
		}
	}

	return output.String()
}

// FormatPasswordStrengthCheck formats password strength check results
func (f *Formatter) FormatPasswordStrengthCheck(result services.StrengthCheckResult) string {
	return result.FormattedResult
}

// getLengthStatus returns status and comment for password length
func (f *Formatter) getLengthStatus(length int) (string, string) {
	status := "✅"
	comment := "(Good!)"

	if length < 8 {
		status = "❌"
		comment = "(Too Short)"
	} else if length < 12 {
		status = "⚠️ "
		comment = "(Could be longer)"
	} else if length >= 16 {
		comment = "(Excellent!)"
	}

	return status, comment
}

// getCharacterSetStatus returns status for character set diversity
func (f *Formatter) getCharacterSetStatus(charTypeCount int) string {
	if charTypeCount < 2 {
		return "❌"
	} else if charTypeCount < 3 {
		return "⚠️ "
	}
	return "✅"
}

// getEntropyStatus returns status for entropy level
func (f *Formatter) getEntropyStatus(entropy float64) string {
	if entropy < 25 {
		return "❌"
	} else if entropy < 40 {
		return "⚠️ "
	}
	return "✅"
}

// FormatWordPasswordGeneration formats word-based password generation results
func (f *Formatter) FormatWordPasswordGeneration(resp application.GenerateWordPasswordResponse) string {
	var output strings.Builder

	// Display each password prominently
	for i, password := range resp.Passwords {
		if len(resp.Passwords) > 1 {
			output.WriteString(fmt.Sprintf("🎯 Password %d:\n", i+1))
		} else {
			output.WriteString("🎯 Your Password:\n")
		}

		// Make password VERY prominent and easy to read
		output.WriteString("┌" + strings.Repeat("─", len(password)+2) + "┐\n")
		output.WriteString(fmt.Sprintf("│ %s │\n", password))
		output.WriteString("└" + strings.Repeat("─", len(password)+2) + "┘\n\n")

		// Brief info on one line
		analysis := resp.Analyses[i]
		output.WriteString(fmt.Sprintf("📝 Based on: \"%s\" | Strategy: %s | Length: %d | Strength: %s %s\n",
			resp.Pattern.Word,
			string(resp.Pattern.Strategy),
			analysis.Password.Length,
			analysis.Strength,
			analysis.StrengthEmoji))

		// Add separator for multiple passwords
		if i < len(resp.Passwords)-1 {
			output.WriteString("\n" + strings.Repeat("─", 60) + "\n\n")
		}
	}

	// Optional: Show detailed analysis only if single password
	if len(resp.Passwords) == 1 {
		analysis := resp.Analyses[0]

		// Minimal analysis section (collapsible feel)
		output.WriteString("\n� Details (security geek info):\n")
		output.WriteString(fmt.Sprintf("   Entropy: %.1f bits | Character types: %s\n",
			analysis.Entropy, strings.Join(analysis.CharacterTypes, ", ")))
		output.WriteString(fmt.Sprintf("   Time to crack: %s\n", analysis.TimeToCrack))

		// Add the sarcastic comment at the end for fun
		if analysis.Celebration != "" {
			output.WriteString(fmt.Sprintf("\n💬 %s\n", analysis.Celebration))
		}
	}

	return output.String()
}

// getLengthDescription returns a description for password length
func (f *Formatter) getLengthDescription(length int) string {
	switch {
	case length < 8:
		return "Too Short"
	case length < 12:
		return "Acceptable"
	case length < 16:
		return "Good"
	case length < 20:
		return "Excellent"
	default:
		return "Outstanding"
	}
}
