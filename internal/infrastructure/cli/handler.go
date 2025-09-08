package cli

import (
	"fmt"
	"os"

	"github.com/kumarasakti/passgen/internal/application"
	"github.com/kumarasakti/passgen/internal/domain/entities"
	"github.com/kumarasakti/passgen/internal/infrastructure/display"
	"github.com/kumarasakti/passgen/internal/infrastructure/repositories"
	"github.com/spf13/cobra"
)

// Handler manages CLI commands and interactions
type Handler struct {
	passwordService *application.PasswordService
	formatter       *Formatter
	config          entities.PasswordConfig
}

// Initializes the CLI handler with default password configuration
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

// Sets up the main CLI command with all subcommands and flags
func (h *Handler) CreateRootCommand(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "passgen",
		Short:   "Generate secure passwords",
		Long:    h.createBanner(version),
		Version: version,
		Run:     h.HandleGeneratePassword,
	}

	// Set custom version template with banner
	rootCmd.SetVersionTemplate(h.createBanner(version) + "\n")

	// Add flags
	h.addFlags(rootCmd)

	// Add subcommands
	rootCmd.AddCommand(h.createCheckCommand())
	rootCmd.AddCommand(h.createPresetCommand())
	rootCmd.AddCommand(h.createWordCommand())

	// Add store commands (Phase 1A: Foundation)
	rootCmd.AddCommand(h.createStoreCommands())

	// Add storage backend commands (Phase 2A: Backend Management)
	rootCmd.AddCommand(NewStorageCommand())

	return rootCmd
}

// Displays branded ASCII art banner with version details
func (h *Handler) createBanner(version string) string {
	return fmt.Sprintf(`
  ____   _    ____ ____   ____ _____ _   _ 
 |  _ \ / \  / ___/ ___| / ___| ____| \ | |
 | |_) / _ \ \___ \___ \| |  _|  _| |  \| |
 |  __/ ___ \ ___) |__) | |_| | |___| |\  |
 |_| /_/   \_\____/____/ \____|_____|_| \_| %s

  🔒 Secure Password Generation & Management 
  🚀 High performance, simple commands, secure storage

passgen is a command-line tool for generating secure passwords with
safe storage and management features.`, version)
}

// Main entry point for password generation with configurable options
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

// Provides detailed password security analysis and improvement recommendations
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

// Generates passwords using predefined security templates
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

// Creates passwords by transforming user-provided words with various strategies
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

// Enables comprehensive password generation with all available character sets and security options
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

// Applies quick password templates that override individual character type settings
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

// Enables password security evaluation with detailed vulnerability assessment
func (h *Handler) createCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "check [password]",
		Short: "Check password strength",
		Long:  "Analyze password strength and provide feedback with specific suggestions for improvement.",
		Args:  cobra.ExactArgs(1),
		Run:   h.HandleCheckPassword,
	}
}

// Provides instant access to common password templates for different security needs
func (h *Handler) createPresetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "preset [type]",
		Short: "Generate password using predefined presets",
		Long:  "Generate password using predefined presets: secure, simple, pin, alphanumeric",
		Args:  cobra.ExactArgs(1),
		Run:   h.HandlePresetPassword,
	}
}

// Enables memorable password creation through word transformation with multiple security strategies
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

// Provides complete password store management with local encryption and optional Git synchronization
func (h *Handler) createStoreCommands() *cobra.Command {
	// Create card displayer
	displayer := display.NewCardDisplayer()

	// Create Phase 1B encrypted repository (shared between handlers)
	encryptedRepo := repositories.NewEncryptedPasswordStoreRepository()

	// Create store handler with encrypted repository
	storeHandler := NewStoreHandler(encryptedRepo, nil)

	// Create Phase 1B store initialization handler with shared repository
	storeInitHandler := NewStoreInitHandler(displayer)
	storeInitHandler.repo = encryptedRepo // Share the same repository instance

	// Create the main store command structure from Phase 1B handler
	storeCmd := storeInitHandler.CreateCommands()

	// Add Phase 1A commands that now use Phase 1B encrypted repository
	storeCmd.AddCommand(storeHandler.createListCommand())
	storeCmd.AddCommand(storeHandler.createAddCommand())
	storeCmd.AddCommand(storeHandler.createGetCommand())
	storeCmd.AddCommand(storeHandler.createListPasswordsCommand())
	storeCmd.AddCommand(storeHandler.createRemoveCommand())
	storeCmd.AddCommand(storeHandler.createRotationCommands())

	// Note: Phase 1B commands (init, clone, sync, setup-gpg, remote, info) are from storeInitHandler
	// Phase 1A commands (list, add, get, remove, rotation) now use encrypted repository from Phase 1B

	return storeCmd
}
