package cli

import (
	"fmt"
	"os"

	"github.com/kumarasakti/passgen/internal/application"
	"github.com/kumarasakti/passgen/internal/domain/entities"
	"github.com/spf13/cobra"
)

// Handler manages CLI commands and interactions
type Handler struct {
	passwordService *application.PasswordService
	formatter       *Formatter
	config          entities.PasswordConfig
}

// NewHandler creates a new CLI handler
func NewHandler() *Handler {
	return &Handler{
		passwordService: application.NewPasswordService(),
		formatter:       NewFormatter(),
		config: entities.PasswordConfig{
			Length:         entities.DefaultLength,
			IncludeLower:   true,
			IncludeUpper:   true,
			IncludeNumbers: false,
			IncludeSymbols: true,
			ExcludeSimilar: false,
			Count:          1,
		},
	}
}

// CreateRootCommand creates and configures the root command
func (h *Handler) CreateRootCommand(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "passgen",
		Short:   "Generate secure passwords",
		Long:    "passgen is a command-line tool for generating secure passwords.",
		Version: version,
		Run:     h.HandleGeneratePassword,
	}

	// Add flags
	h.addFlags(rootCmd)

	// Add subcommands
	rootCmd.AddCommand(h.createCheckCommand())
	rootCmd.AddCommand(h.createPresetCommand())
	rootCmd.AddCommand(h.createWordCommand())
	
	// Add store commands (Phase 1A: Foundation)
	rootCmd.AddCommand(h.createStoreCommands())

	return rootCmd
}

// HandleGeneratePassword handles the main password generation
func (h *Handler) HandleGeneratePassword(cmd *cobra.Command, args []string) {
	// Handle convenience flags
	h.handleConvenienceFlags(cmd)

	req := application.GeneratePasswordRequest{Config: h.config}
	resp, err := h.passwordService.GeneratePasswords(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating password: %v\n", err)
		os.Exit(1)
	}

	output := h.formatter.FormatPasswordGeneration(resp.Analyses, h.config.ExcludeSimilar)
	fmt.Print(output)
}

// HandleCheckPassword handles password strength checking
func (h *Handler) HandleCheckPassword(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: exactly one password argument required\n")
		os.Exit(1)
	}

	req := application.CheckPasswordRequest{Password: args[0]}
	resp := h.passwordService.CheckPasswordStrength(req)

	output := h.formatter.FormatPasswordStrengthCheck(resp.Result)
	fmt.Print(output)
}

// HandlePresetPassword handles preset password generation
func (h *Handler) HandlePresetPassword(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: exactly one preset type argument required\n")
		os.Exit(1)
	}

	presetType := args[0]
	resp, err := h.passwordService.GeneratePresetPassword(presetType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating preset password: %v\n", err)
		fmt.Fprintf(os.Stderr, "Available presets: secure, simple, pin, alphanumeric\n")
		os.Exit(1)
	}

	output := h.formatter.FormatPasswordGeneration(resp.Analyses, false)
	fmt.Print(output)
}

// HandleWordPassword handles word-based password generation
func (h *Handler) HandleWordPassword(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: exactly one word argument required\n")
		os.Exit(1)
	}

	word := args[0]

	// Get flags
	strategy, _ := cmd.Flags().GetString("strategy")
	complexity, _ := cmd.Flags().GetString("complexity")
	count, _ := cmd.Flags().GetInt("count")

	// Validate strategy
	var transformationStrategy entities.TransformationStrategy
	switch strategy {
	case "leetspeak":
		transformationStrategy = entities.StrategyLeetspeak
	case "mixed-case":
		transformationStrategy = entities.StrategyMixedCase
	case "suffix":
		transformationStrategy = entities.StrategySuffix
	case "prefix":
		transformationStrategy = entities.StrategyPrefix
	case "insert":
		transformationStrategy = entities.StrategyInsert
	case "hybrid":
		transformationStrategy = entities.StrategyHybrid
	default:
		fmt.Fprintf(os.Stderr, "Error: invalid strategy '%s'. Available: leetspeak, mixed-case, suffix, prefix, insert, hybrid\n", strategy)
		os.Exit(1)
	}

	// Validate complexity
	var complexityLevel entities.ComplexityLevel
	switch complexity {
	case "low":
		complexityLevel = entities.ComplexityLow
	case "medium":
		complexityLevel = entities.ComplexityMedium
	case "high":
		complexityLevel = entities.ComplexityHigh
	default:
		fmt.Fprintf(os.Stderr, "Error: invalid complexity '%s'. Available: low, medium, high\n", complexity)
		os.Exit(1)
	}

	// Create request
	req := application.GenerateWordPasswordRequest{
		Word:       word,
		Strategy:   transformationStrategy,
		Complexity: complexityLevel,
		Count:      count,
	}

	// Generate word-based passwords
	resp, err := h.passwordService.GenerateWordPasswords(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating word-based password: %v\n", err)
		os.Exit(1)
	}

	// Format and display output
	output := h.formatter.FormatWordPasswordGeneration(resp)
	fmt.Print(output)
}

// addFlags adds command line flags to the root command
func (h *Handler) addFlags(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&h.config.Length, "length", "l", entities.DefaultLength, "Password length")
	cmd.Flags().BoolVar(&h.config.IncludeLower, "lower", true, "Include lowercase letters")
	cmd.Flags().BoolVar(&h.config.IncludeUpper, "upper", true, "Include uppercase letters")
	cmd.Flags().BoolVarP(&h.config.IncludeNumbers, "numbers", "n", false, "Include numbers")
	cmd.Flags().BoolVarP(&h.config.IncludeSymbols, "symbols", "s", true, "Include symbols")
	cmd.Flags().BoolVar(&h.config.ExcludeSimilar, "exclude-similar", false, "Exclude similar characters (il1Lo0O)")
	cmd.Flags().StringVar(&h.config.ExcludeChars, "exclude", "", "Characters to exclude from password")
	cmd.Flags().IntVarP(&h.config.Count, "count", "c", 1, "Number of passwords to generate")

	// Add convenience flags
	cmd.Flags().BoolP("secure", "S", false, "Generate secure password (includes all character types)")
	cmd.Flags().BoolP("simple", "m", false, "Generate simple password (only letters and numbers)")
	cmd.Flags().BoolP("alphanumeric", "a", false, "Generate alphanumeric password (letters and numbers)")
}

// handleConvenienceFlags processes convenience flags that modify configuration
func (h *Handler) handleConvenienceFlags(cmd *cobra.Command) {
	if secure, _ := cmd.Flags().GetBool("secure"); secure {
		h.config.IncludeLower = true
		h.config.IncludeUpper = true
		h.config.IncludeNumbers = true
		h.config.IncludeSymbols = true
	}

	if simple, _ := cmd.Flags().GetBool("simple"); simple {
		h.config.IncludeLower = true
		h.config.IncludeUpper = true
		h.config.IncludeNumbers = true
		h.config.IncludeSymbols = false
	}

	if alphanumeric, _ := cmd.Flags().GetBool("alphanumeric"); alphanumeric {
		h.config.IncludeLower = true
		h.config.IncludeUpper = true
		h.config.IncludeNumbers = true
		h.config.IncludeSymbols = false
	}
}

// createCheckCommand creates the check subcommand
func (h *Handler) createCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "check [password]",
		Short: "Check password strength",
		Long:  "Analyze password strength and provide feedback with specific suggestions for improvement.",
		Args:  cobra.ExactArgs(1),
		Run:   h.HandleCheckPassword,
	}
}

// createPresetCommand creates the preset subcommand
func (h *Handler) createPresetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "preset [type]",
		Short: "Generate password using predefined presets",
		Long:  "Generate password using predefined presets: secure, simple, pin, alphanumeric",
		Args:  cobra.ExactArgs(1),
		Run:   h.HandlePresetPassword,
	}
}

// createWordCommand creates the word subcommand
func (h *Handler) createWordCommand() *cobra.Command {
	wordCmd := &cobra.Command{
		Use:   "word [word]",
		Short: "Generate password based on a word",
		Long: `Generate password based on a word with various transformation strategies:
  - leetspeak: Replace characters with numbers/symbols (e→3, a→@, etc.)
  - mixed-case: Apply mixed capitalization patterns
  - suffix: Add numbers and symbols at the end
  - prefix: Add symbols at the beginning
  - insert: Insert characters within the word
  - hybrid: Combine multiple strategies (default)
  
Examples:
  passgen word "security"                    # Default hybrid transformation
  passgen word "security" --strategy leetspeak    # S3cur1ty transformation
  passgen word "security" --complexity high       # Maximum complexity
  passgen word "security" --count 3               # Generate 3 variations`,
		Args: cobra.ExactArgs(1),
		Run:  h.HandleWordPassword,
	}

	// Add word-specific flags
	wordCmd.Flags().String("strategy", "hybrid", "Transformation strategy (leetspeak, mixed-case, suffix, prefix, insert, hybrid)")
	wordCmd.Flags().String("complexity", "medium", "Complexity level (low, medium, high)")
	wordCmd.Flags().IntP("count", "c", 1, "Number of password variations to generate")

	return wordCmd
}

// createStoreCommands creates the store command tree (Phase 1A: Foundation)
func (h *Handler) createStoreCommands() *cobra.Command {
	// For Phase 1A, we create a mock store handler to demonstrate the enhanced card display
	// This will be replaced with real implementations in Phase 1B
	storeHandler := NewStoreHandler(nil, nil) // nil repos for Phase 1A demo
	return storeHandler.CreateStoreCommands()
}
