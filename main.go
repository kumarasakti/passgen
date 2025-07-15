package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	DefaultLength = 14
	Lowercase     = "abcdefghijklmnopqrstuvwxyz"
	Uppercase     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers       = "0123456789"
	Symbols       = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

// Version can be overridden at build time using -ldflags "-X main.Version=x.y.z"
var Version = "v1.0.5"

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

type PasswordAnalysis struct {
	Password       string
	Length         int
	CharsetSize    int
	CharacterTypes []string
	Entropy        float64
	Strength       string
	StrengthEmoji  string
	TimeToCrack    string
	SecurityLevel  string
	Tips           []string
	Celebration    string
}

var config PasswordConfig

func calculateCharsetSize(cfg PasswordConfig) int {
	size := 0
	if cfg.IncludeLower {
		size += 26
	}
	if cfg.IncludeUpper {
		size += 26
	}
	if cfg.IncludeNumbers {
		size += 10
	}
	if cfg.IncludeSymbols {
		size += len(Symbols)
	}

	// Adjust for excluded characters
	if cfg.ExcludeSimilar {
		similar := "il1Lo0O"
		for _, char := range similar {
			if cfg.IncludeLower && strings.Contains(Lowercase, string(char)) {
				size--
			}
			if cfg.IncludeUpper && strings.Contains(Uppercase, string(char)) {
				size--
			}
			if cfg.IncludeNumbers && strings.Contains(Numbers, string(char)) {
				size--
			}
		}
	}

	if cfg.ExcludeChars != "" {
		for _, char := range cfg.ExcludeChars {
			if cfg.IncludeLower && strings.Contains(Lowercase, string(char)) {
				size--
			}
			if cfg.IncludeUpper && strings.Contains(Uppercase, string(char)) {
				size--
			}
			if cfg.IncludeNumbers && strings.Contains(Numbers, string(char)) {
				size--
			}
			if cfg.IncludeSymbols && strings.Contains(Symbols, string(char)) {
				size--
			}
		}
	}

	return size
}

func analyzePassword(password string, cfg PasswordConfig) PasswordAnalysis {
	length := len(password)
	charsetSize := calculateCharsetSize(cfg)

	// Calculate entropy: log2(charset^length)
	entropy := float64(length) * math.Log2(float64(charsetSize))

	// Determine character types present
	var charTypes []string
	if cfg.IncludeLower {
		charTypes = append(charTypes, "Lowercase")
	}
	if cfg.IncludeUpper {
		charTypes = append(charTypes, "Uppercase")
	}
	if cfg.IncludeNumbers {
		charTypes = append(charTypes, "Numbers")
	}
	if cfg.IncludeSymbols {
		charTypes = append(charTypes, "Symbols")
	}

	// Determine strength and security level
	var strength, strengthEmoji, securityLevel, celebration string
	var tips []string

	switch {
	case entropy >= 100:
		strength = "Extremely Strong"
		strengthEmoji = "🔥"
		securityLevel = "Quantum-resistant for the foreseeable future!"
		celebration = "Brr, that's ice cold security! Even hackers are shivering! 🥶"
	case entropy >= 80:
		strength = "Very Strong"
		strengthEmoji = "💪"
		securityLevel = "Exceeds security standards for high-value accounts"
		celebration = "Someone's taking this security thing seriously! 🌟"
	case entropy >= 60:
		strength = "Strong"
		strengthEmoji = "💯"
		securityLevel = "Great for securing important accounts"
		celebration = "Not bad, you actually read the security guidelines! 🎯"
	case entropy >= 40:
		strength = "Medium"
		strengthEmoji = "⚡"
		securityLevel = "Adequate for most general purposes"
		celebration = "Well, it's... adequate. I guess that's something! 👍"
		if length < 12 {
			tips = append(tips, "Consider using 12+ characters for better security")
		}
		if len(charTypes) < 3 {
			tips = append(tips, "Add more character types (symbols, numbers) for stronger security")
		}
	case entropy >= 25:
		strength = "Weak"
		strengthEmoji = "😰"
		securityLevel = "Suitable only for low-security uses"
		celebration = "Oh honey, we need to talk about your password choices... 💪"
		tips = append(tips, "Use at least 12 characters")
		tips = append(tips, "Include uppercase, lowercase, numbers, and symbols")
		tips = append(tips, "Try `passgen --secure` for maximum protection!")
	default:
		strength = "Very Weak"
		strengthEmoji = "🚨"
		securityLevel = "Not recommended for any security purposes"
		celebration = "Yikes! Even my grandma would crack this in her sleep! 🚀"
		tips = append(tips, "Use at least 12 characters")
		tips = append(tips, "Include multiple character types")
		tips = append(tips, "Try `passgen --secure -l 16` for excellent security!")
	}

	// Calculate time to crack (assuming 1 trillion guesses per second)
	guessesPerSecond := 1e12
	possibleCombinations := math.Pow(float64(charsetSize), float64(length))
	secondsToCrack := possibleCombinations / (2 * guessesPerSecond) // Average case

	var timeToCrack string
	if secondsToCrack < 60 {
		timeToCrack = "Less than a minute"
	} else if secondsToCrack < 3600 {
		timeToCrack = fmt.Sprintf("%.1f minutes", secondsToCrack/60)
	} else if secondsToCrack < 86400 {
		timeToCrack = fmt.Sprintf("%.1f hours", secondsToCrack/3600)
	} else if secondsToCrack < 31536000 {
		timeToCrack = fmt.Sprintf("%.1f days", secondsToCrack/86400)
	} else if secondsToCrack < 31536000000 {
		timeToCrack = fmt.Sprintf("%.1f years", secondsToCrack/31536000)
	} else {
		// For very large numbers, use scientific notation
		years := secondsToCrack / 31536000
		if years > 1e15 {
			timeToCrack = fmt.Sprintf("%.1e years", years)
		} else {
			timeToCrack = fmt.Sprintf("%.0f years", years)
		}
	}

	return PasswordAnalysis{
		Password:       password,
		Length:         length,
		CharsetSize:    charsetSize,
		CharacterTypes: charTypes,
		Entropy:        entropy,
		Strength:       strength,
		StrengthEmoji:  strengthEmoji,
		TimeToCrack:    timeToCrack,
		SecurityLevel:  securityLevel,
		Tips:           tips,
		Celebration:    celebration,
	}
}

func printPasswordAnalysis(analysis PasswordAnalysis) {
	// Header with appropriate emoji
	if analysis.Entropy >= 60 {
		fmt.Printf("🎉 Password Generated Successfully! 🎉\n\n")
	} else if analysis.Entropy >= 40 {
		fmt.Printf("✨ Password Generated! ✨\n\n")
	} else {
		fmt.Printf("⚠️  Basic Password Generated ⚠️\n\n")
	}

	// Display password in bold
	fmt.Printf("Password: \033[1m%s\033[0m\n\n", analysis.Password)

	// Analysis section
	fmt.Printf("📊 Password Analysis:\n")

	// Length assessment
	lengthStatus := "✅"
	lengthComment := "(Good!)"
	if analysis.Length < 8 {
		lengthStatus = "❌"
		lengthComment = "(Too Short)"
	} else if analysis.Length < 12 {
		lengthStatus = "⚠️ "
		lengthComment = "(Could be longer)"
	} else if analysis.Length >= 16 {
		lengthComment = "(Excellent!)"
	}
	fmt.Printf("%s Length: %d characters %s\n", lengthStatus, analysis.Length, lengthComment)

	// Character sets
	charStatus := "✅"
	if len(analysis.CharacterTypes) < 2 {
		charStatus = "❌"
	} else if len(analysis.CharacterTypes) < 3 {
		charStatus = "⚠️ "
	}
	fmt.Printf("%s Character Sets: %s\n", charStatus, strings.Join(analysis.CharacterTypes, ", "))

	// Entropy
	entropyStatus := "✅"
	if analysis.Entropy < 25 {
		entropyStatus = "❌"
	} else if analysis.Entropy < 40 {
		entropyStatus = "⚠️ "
	}
	fmt.Printf("%s Entropy: %.1f bits (%s!)\n", entropyStatus, analysis.Entropy, analysis.Strength)

	// Overall strength
	fmt.Printf("✅ Strength: %s %s\n\n", analysis.Strength, analysis.StrengthEmoji)

	// Security assessment
	fmt.Printf("🔒 Security Assessment:\n")
	fmt.Printf("• This password would take approximately %s to crack with modern hardware\n", analysis.TimeToCrack)

	if len(analysis.CharacterTypes) > 1 {
		fmt.Printf("• Contains %d different character types for %s complexity\n",
			len(analysis.CharacterTypes),
			map[int]string{2: "good", 3: "excellent", 4: "maximum"}[len(analysis.CharacterTypes)])
	}

	fmt.Printf("• %s\n", analysis.SecurityLevel)

	if config.ExcludeSimilar {
		fmt.Printf("• No similar characters (like i, l, 1, O, 0) to avoid confusion\n")
	}

	// Tips if any
	if len(analysis.Tips) > 0 {
		fmt.Printf("\n💡 Tips for improvement:\n")
		for _, tip := range analysis.Tips {
			fmt.Printf("• %s\n", tip)
		}
	}

	// Celebration message
	fmt.Printf("\n%s\n", analysis.Celebration)
}

func generatePassword(cfg PasswordConfig) (string, error) {
	charset := ""

	if cfg.IncludeLower {
		charset += Lowercase
	}
	if cfg.IncludeUpper {
		charset += Uppercase
	}
	if cfg.IncludeNumbers {
		charset += Numbers
	}
	if cfg.IncludeSymbols {
		charset += Symbols
	}

	if charset == "" {
		return "", fmt.Errorf("no character sets selected")
	}

	// Remove similar characters if requested
	if cfg.ExcludeSimilar {
		similar := "il1Lo0O"
		for _, char := range similar {
			charset = strings.ReplaceAll(charset, string(char), "")
		}
	}

	// Remove excluded characters
	if cfg.ExcludeChars != "" {
		for _, char := range cfg.ExcludeChars {
			charset = strings.ReplaceAll(charset, string(char), "")
		}
	}

	if len(charset) == 0 {
		return "", fmt.Errorf("no characters available after exclusions")
	}

	password := make([]byte, cfg.Length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}

func checkPasswordStrength(password string) (string, int) {
	score := 0
	feedback := []string{}
	sarcasticComments := []string{}

	// Length check
	if len(password) >= 12 {
		score += 2
		if len(password) >= 16 {
			sarcasticComments = append(sarcasticComments, "Wow, someone actually read the security guidelines! 👏")
		}
	} else if len(password) >= 8 {
		score += 1
		sarcasticComments = append(sarcasticComments, "8 characters? How... minimalistic of you 🤔")
	} else {
		feedback = append(feedback, "Password should be at least 8 characters long")
		sarcasticComments = append(sarcasticComments, "Really? That's shorter than most people's names! 😅")
	}

	// Character variety checks
	if matched, _ := regexp.MatchString(`[a-z]`, password); matched {
		score += 1
	} else {
		feedback = append(feedback, "Add lowercase letters")
		sarcasticComments = append(sarcasticComments, "No lowercase letters? Are we SHOUTING all the time? 📢")
	}

	if matched, _ := regexp.MatchString(`[A-Z]`, password); matched {
		score += 1
	} else {
		feedback = append(feedback, "Add uppercase letters")
		sarcasticComments = append(sarcasticComments, "No capitals? I guess we're going for the e.e. cummings aesthetic 🎭")
	}

	if matched, _ := regexp.MatchString(`[0-9]`, password); matched {
		score += 1
	} else {
		feedback = append(feedback, "Add numbers")
		sarcasticComments = append(sarcasticComments, "Numbers are optional now? Math teachers everywhere are crying 😢")
	}

	if matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]`, password); matched {
		score += 2
	} else {
		feedback = append(feedback, "Add special characters")
		sarcasticComments = append(sarcasticComments, "No symbols? Your password is as plain as unseasoned chicken 🐔")
	}

	// Bonus for length
	if len(password) >= 16 {
		score += 1
	}

	var strength, strengthEmoji, celebration string
	switch {
	case score >= 7:
		strength = "Very Strong"
		strengthEmoji = "🔥"
		celebration = "Impressive! Your password could probably withstand a zombie apocalypse! 🧟‍♂️"
	case score >= 5:
		strength = "Strong"
		strengthEmoji = "💪"
		celebration = "Not bad! Your password has some real backbone! 🦴"
	case score >= 3:
		strength = "Medium"
		strengthEmoji = "😐"
		celebration = "It's... adequate. Like a participation trophy for password security 🏆"
	case score >= 1:
		strength = "Weak"
		strengthEmoji = "😰"
		celebration = "Yikes! This password couldn't protect a diary from a nosy sibling! 📖"
	default:
		strength = "Very Weak"
		strengthEmoji = "🚨"
		celebration = "Oh dear... this password is weaker than my WiFi signal in the basement! 📶"
	}

	result := "🔍 Password Analysis Results:\n"
	result += fmt.Sprintf("Strength: %s %s (Score: %d/8)\n", strength, strengthEmoji, score)
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

	return result, score
}

var rootCmd = &cobra.Command{
	Use:   "passgen",
	Short: "A secure password generator CLI tool",
	Long: `passgen is a command-line tool for generating secure passwords.
It supports various character sets, customizable length, and advanced options
like excluding similar characters and generating multiple passwords.

Examples:
  passgen                           # Generate default password
  passgen -l 16 -s                  # Generate 16-char password with symbols
  passgen --secure -l 20            # Generate secure 20-char password
  passgen -c 5 -l 12                # Generate 5 passwords of 12 characters
  passgen --exclude-similar -s      # Exclude similar characters
  passgen --exclude "aeiou"         # Exclude vowels
  passgen check "mypassword"        # Check password strength`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < config.Count; i++ {
			password, err := generatePassword(config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating password: %v\n", err)
				os.Exit(1)
			}

			// Analyze and display password with enhanced output
			analysis := analyzePassword(password, config)
			printPasswordAnalysis(analysis)

			// Add separator between multiple passwords
			if config.Count > 1 && i < config.Count-1 {
				fmt.Println("\n" + strings.Repeat("=", 50) + "\n")
				time.Sleep(100 * time.Millisecond) // Brief pause for dramatic effect
			}
		}
	},
}

var checkCmd = &cobra.Command{
	Use:   "check [password]",
	Short: "Check password strength",
	Long:  `Check the strength of a password and get suggestions for improvement.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		password := args[0]
		result, _ := checkPasswordStrength(password)
		fmt.Println(result)
	},
}

var presetCmd = &cobra.Command{
	Use:   "preset [type]",
	Short: "Generate password using predefined presets",
	Long: `Generate password using predefined presets:
  - secure: All character types, 16 characters
  - simple: Letters and numbers only, 12 characters
  - pin: Numbers only, 6 characters
  - alphanumeric: Letters and numbers, 12 characters`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		preset := args[0]
		var cfg PasswordConfig

		switch preset {
		case "secure":
			cfg = PasswordConfig{
				Length: 16, IncludeLower: true, IncludeUpper: true,
				IncludeNumbers: true, IncludeSymbols: true, Count: 1,
			}
		case "simple":
			cfg = PasswordConfig{
				Length: 12, IncludeLower: true, IncludeUpper: true,
				IncludeNumbers: true, IncludeSymbols: false, Count: 1,
			}
		case "pin":
			cfg = PasswordConfig{
				Length: 6, IncludeLower: false, IncludeUpper: false,
				IncludeNumbers: true, IncludeSymbols: false, Count: 1,
			}
		case "alphanumeric":
			cfg = PasswordConfig{
				Length: 12, IncludeLower: true, IncludeUpper: true,
				IncludeNumbers: true, IncludeSymbols: false, Count: 1,
			}
		default:
			fmt.Fprintf(os.Stderr, "Unknown preset: %s\n", preset)
			fmt.Fprintf(os.Stderr, "Available presets: secure, simple, pin, alphanumeric\n")
			os.Exit(1)
		}

		password, err := generatePassword(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating password: %v\n", err)
			os.Exit(1)
		}

		// Analyze and display password with enhanced output
		analysis := analyzePassword(password, cfg)
		printPasswordAnalysis(analysis)
	},
}

func init() {
	// Set default values
	config = PasswordConfig{
		Length:         DefaultLength,
		IncludeLower:   true,
		IncludeUpper:   true,
		IncludeNumbers: false,
		IncludeSymbols: true,
		ExcludeSimilar: false,
		Count:          1,
	}

	// Add flags to root command
	rootCmd.Flags().IntVarP(&config.Length, "length", "l", DefaultLength, "Password length")
	rootCmd.Flags().BoolVar(&config.IncludeLower, "lower", true, "Include lowercase letters")
	rootCmd.Flags().BoolVar(&config.IncludeUpper, "upper", true, "Include uppercase letters")
	rootCmd.Flags().BoolVarP(&config.IncludeNumbers, "numbers", "n", false, "Include numbers")
	rootCmd.Flags().BoolVarP(&config.IncludeSymbols, "symbols", "s", true, "Include symbols")
	rootCmd.Flags().BoolVar(&config.ExcludeSimilar, "exclude-similar", false, "Exclude similar characters (il1Lo0O)")
	rootCmd.Flags().StringVar(&config.ExcludeChars, "exclude", "", "Characters to exclude from password")
	rootCmd.Flags().IntVarP(&config.Count, "count", "c", 1, "Number of passwords to generate")

	// Add convenience flags
	rootCmd.Flags().BoolP("secure", "S", false, "Generate secure password (includes all character types)")
	rootCmd.Flags().BoolP("simple", "m", false, "Generate simple password (only letters and numbers)")
	rootCmd.Flags().BoolP("alphanumeric", "a", false, "Generate alphanumeric password (letters and numbers)")

	// Handle convenience flags
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if secure, _ := cmd.Flags().GetBool("secure"); secure {
			config.IncludeLower = true
			config.IncludeUpper = true
			config.IncludeNumbers = true
			config.IncludeSymbols = true
		}

		if simple, _ := cmd.Flags().GetBool("simple"); simple {
			config.IncludeLower = true
			config.IncludeUpper = true
			config.IncludeNumbers = true
			config.IncludeSymbols = false
		}

		if alphanumeric, _ := cmd.Flags().GetBool("alphanumeric"); alphanumeric {
			config.IncludeLower = true
			config.IncludeUpper = true
			config.IncludeNumbers = true
			config.IncludeSymbols = false
		}
	}

	// Add subcommands
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(presetCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
