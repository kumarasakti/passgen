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
