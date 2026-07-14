package cli

import (
	"fmt"

	"github.com/kumarasakti/passgen/internal/application"
	"github.com/kumarasakti/passgen/internal/domain/entities"
	"github.com/spf13/cobra"
)

// ConfigHandler handles config subcommands
type ConfigHandler struct {
	manager *application.ConfigManager
}

// NewConfigHandler creates a new config handler
func NewConfigHandler() *ConfigHandler {
	manager, err := application.NewConfigManager()
	if err != nil {
		return &ConfigHandler{}
	}
	return &ConfigHandler{manager: manager}
}

// CreateCommands creates the config subcommand tree
func (h *ConfigHandler) CreateCommands() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage passgen configuration",
		Long: `Manage default password generation settings stored in ~/.passgen/config.yaml.

Configuration is loaded automatically on every passgen run. CLI flags
override config file values, so you only need to set your preferences once.

Examples:
  passgen config init              # Create config with default values
  passgen config show              # Display current configuration
  passgen config set length 20     # Set default password length to 20
  passgen config set no_repeat true  # Enable no-repeat by default`,
	}

	configCmd.AddCommand(h.createInitCommand())
	configCmd.AddCommand(h.createShowCommand())
	configCmd.AddCommand(h.createSetCommand())

	return configCmd
}

func (h *ConfigHandler) createInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create config file with default values",
		Long:  "Creates ~/.passgen/config.yaml with default generation settings.",
		RunE:  h.handleInit,
	}
}

func (h *ConfigHandler) createShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Display current configuration",
		RunE:  h.handleShow,
	}
}

func (h *ConfigHandler) createSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long: `Set a single configuration value. Valid keys:
  length, include_lower, include_upper, include_numbers, include_symbols,
  exclude_similar, exclude_chars, no_repeat, count,
  word_strategy, word_complexity, word_count

Examples:
  passgen config set length 20
  passgen config set include_numbers true
  passgen config set no_repeat true
  passgen config set exclude_chars "!#"
  passgen config set word_strategy leetspeak
  passgen config set word_complexity high
  passgen config set word_count 5`,
		Args: cobra.ExactArgs(2),
		RunE:  h.handleSet,
	}
}

func (h *ConfigHandler) handleInit(cmd *cobra.Command, args []string) error {
	if h.manager == nil {
		return fmt.Errorf("config manager not initialized")
	}

	if err := h.manager.Init(); err != nil {
		return err
	}

	fmt.Printf("Config file created at %s\n", h.manager.ConfigPath())
	fmt.Printf("\nDefault settings:\n")
	config, _ := h.manager.Load()
	printConfig(config)
	fmt.Printf("\nEdit the file directly or use 'passgen config set <key> <value>' to change settings.\n")

	return nil
}

func (h *ConfigHandler) handleShow(cmd *cobra.Command, args []string) error {
	if h.manager == nil {
		return fmt.Errorf("config manager not initialized")
	}

	config, err := h.manager.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Printf("Config file: %s\n\n", h.manager.ConfigPath())
	printConfig(config)

	return nil
}

func (h *ConfigHandler) handleSet(cmd *cobra.Command, args []string) error {
	if h.manager == nil {
		return fmt.Errorf("config manager not initialized")
	}

	key := args[0]
	value := args[1]

	if err := h.manager.Set(key, value); err != nil {
		return err
	}

	fmt.Printf("Set %s = %s\n", key, value)
	return nil
}

func printConfig(config entities.PassgenConfig) {
	g := config.Generation
	fmt.Printf("generation:\n")
	fmt.Printf("  length:           %d\n", g.Length)
	fmt.Printf("  include_lower:    %v\n", g.IncludeLower)
	fmt.Printf("  include_upper:    %v\n", g.IncludeUpper)
	fmt.Printf("  include_numbers:  %v\n", g.IncludeNumbers)
	fmt.Printf("  include_symbols:  %v\n", g.IncludeSymbols)
	fmt.Printf("  exclude_similar:  %v\n", g.ExcludeSimilar)
	fmt.Printf("  exclude_chars:    %q\n", g.ExcludeChars)
	fmt.Printf("  no_repeat:        %v\n", g.NoRepeat)
	fmt.Printf("  count:            %d\n", g.Count)

	w := config.Word
	fmt.Printf("\nword:\n")
	fmt.Printf("  strategy:         %s\n", w.Strategy)
	fmt.Printf("  complexity:       %s\n", w.Complexity)
	fmt.Printf("  count:            %d\n", w.Count)
}
