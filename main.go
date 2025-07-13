package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

const (
	DefaultLength = 12
	Lowercase     = "abcdefghijklmnopqrstuvwxyz"
	Uppercase     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers       = "0123456789"
	Symbols       = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	Version       = "1.0.0"
)

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

var config PasswordConfig

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

	if len(password) >= 12 {
		score += 2
	} else if len(password) >= 8 {
		score += 1
	} else {
		feedback = append(feedback, "Password should be at least 8 characters long")
	}

	if matched, _ := regexp.MatchString(`[a-z]`, password); matched {
		score += 1
	} else {
		feedback = append(feedback, "Add lowercase letters")
	}

	if matched, _ := regexp.MatchString(`[A-Z]`, password); matched {
		score += 1
	} else {
		feedback = append(feedback, "Add uppercase letters")
	}

	if matched, _ := regexp.MatchString(`[0-9]`, password); matched {
		score += 1
	} else {
		feedback = append(feedback, "Add numbers")
	}

	if matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]`, password); matched {
		score += 2
	} else {
		feedback = append(feedback, "Add special characters")
	}

	if len(password) >= 16 {
		score += 1
	}

	var strength string
	switch {
	case score >= 7:
		strength = "Very Strong"
	case score >= 5:
		strength = "Strong"
	case score >= 3:
		strength = "Medium"
	case score >= 1:
		strength = "Weak"
	default:
		strength = "Very Weak"
	}

	result := fmt.Sprintf("Strength: %s (Score: %d/8)", strength, score)
	if len(feedback) > 0 {
		result += "\nSuggestions: " + strings.Join(feedback, ", ")
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
			fmt.Println(password)
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
		fmt.Println(password)
	},
}

func init() {
	// Set default values
	config = PasswordConfig{
		Length:         DefaultLength,
		IncludeLower:   true,
		IncludeUpper:   true,
		IncludeNumbers: true,
		IncludeSymbols: false,
		ExcludeSimilar: false,
		Count:          1,
	}

	// Add flags to root command
	rootCmd.Flags().IntVarP(&config.Length, "length", "l", DefaultLength, "Password length")
	rootCmd.Flags().BoolVar(&config.IncludeLower, "lower", true, "Include lowercase letters")
	rootCmd.Flags().BoolVar(&config.IncludeUpper, "upper", true, "Include uppercase letters")
	rootCmd.Flags().BoolVarP(&config.IncludeNumbers, "numbers", "n", true, "Include numbers")
	rootCmd.Flags().BoolVarP(&config.IncludeSymbols, "symbols", "s", false, "Include symbols")
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
