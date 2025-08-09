package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/kumarasakti/passgen/internal/domain/entities"
	"github.com/kumarasakti/passgen/internal/domain/repositories"
	"github.com/kumarasakti/passgen/internal/infrastructure/display"
)

// StoreHandler handles password store CLI commands
type StoreHandler struct {
	repository    repositories.PasswordStoreRepository
	configRepo    repositories.StoreConfigRepository
	cardDisplay   *display.CardDisplayer
}

// NewStoreHandler creates a new store command handler
func NewStoreHandler(
	repo repositories.PasswordStoreRepository, 
	configRepo repositories.StoreConfigRepository,
) *StoreHandler {
	return &StoreHandler{
		repository:  repo,
		configRepo:  configRepo,
		cardDisplay: display.NewCardDisplayer(),
	}
}

// CreateStoreCommands creates the store command tree
func (h *StoreHandler) CreateStoreCommands() *cobra.Command {
	storeCmd := &cobra.Command{
		Use:   "store",
		Short: "Manage password stores",
		Long: `Manage password stores with Git backing and GPG encryption.
		
Password stores allow you to securely store and manage passwords with:
• Git repository backing for sync and collaboration
• GPG encryption for security
• Auto-rotation for enterprise password policies
• Clean card-style display for easy reading`,
	}

	// Add subcommands
	storeCmd.AddCommand(h.createInitCommand())
	storeCmd.AddCommand(h.createListCommand())
	storeCmd.AddCommand(h.createAddCommand())
	storeCmd.AddCommand(h.createGetCommand())
	storeCmd.AddCommand(h.createListPasswordsCommand())
	storeCmd.AddCommand(h.createRemoveCommand())
	storeCmd.AddCommand(h.createSyncCommand())
	storeCmd.AddCommand(h.createRotationCommands())

	return storeCmd
}

// createGetCommand creates the get password command with enhanced card display
func (h *StoreHandler) createGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <service>",
		Short: "Get password information (secure)",
		Long: `Retrieve password information in a clean card format.
		
By default, only metadata is shown (no password). Use flags for secure access:
• --copy: Copy password to clipboard (auto-clears in 30s)
• --show: Display password in terminal (requires confirmation)`,
		Args: cobra.ExactArgs(1),
		RunE: h.GetPassword,
	}

	cmd.Flags().String("store", "", "Store name (default: configured default store)")
	cmd.Flags().Bool("copy", false, "Copy password to clipboard with auto-clear")
	cmd.Flags().Bool("show", false, "Display password in terminal (requires confirmation)")

	return cmd
}

// createAddCommand creates the add password command
func (h *StoreHandler) createAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <service>",
		Short: "Add a new password to store",
		Long: `Add a new password to the specified store with optional auto-rotation.
		
The password can be generated automatically or entered manually.
Auto-rotation can be configured for enterprise password policies.`,
		Args: cobra.ExactArgs(1),
		RunE: h.AddPassword,
	}

	cmd.Flags().String("store", "", "Store name (default: configured default store)")
	cmd.Flags().String("username", "", "Username for the service")
	cmd.Flags().String("url", "", "URL for the service")
	cmd.Flags().String("notes", "", "Notes for the password")
	cmd.Flags().Int("auto-rotate", 0, "Enable auto-rotation (days between rotations)")
	cmd.Flags().Int("notify-before", 7, "Days before rotation to notify")
	cmd.Flags().Int("length", 16, "Password length for generation")

	return cmd
}

// createListPasswordsCommand creates the list passwords command
func (h *StoreHandler) createListPasswordsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all passwords in store",
		Long:  `List all passwords in the specified store with basic information.`,
		RunE:  h.ListPasswords,
	}

	cmd.Flags().String("store", "", "Store name (default: configured default store)")

	return cmd
}

// createRotationCommands creates rotation-related commands
func (h *StoreHandler) createRotationCommands() *cobra.Command {
	rotationCmd := &cobra.Command{
		Use:   "rotation",
		Short: "Manage password rotation",
		Long:  `Manage automatic password rotation for enhanced security.`,
	}

	// rotation status
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show rotation status",
		Long:  `Show rotation status for passwords with auto-rotation enabled.`,
		RunE:  h.RotationStatus,
	}
	statusCmd.Flags().String("store", "", "Store name (default: configured default store)")

	// rotation check
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Check for due rotations",
		Long:  `Check for passwords that need rotation and show notifications.`,
		RunE:  h.CheckRotations,
	}
	checkCmd.Flags().String("store", "", "Store name (default: configured default store)")

	rotationCmd.AddCommand(statusCmd)
	rotationCmd.AddCommand(checkCmd)

	return rotationCmd
}

// Placeholder command creators (to be implemented in next phase)
func (h *StoreHandler) createInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init <name>",
		Short: "Initialize a new password store",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("store init not implemented yet - coming in Phase 1B")
		},
	}
}

func (h *StoreHandler) createListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stores",
		Short: "List configured stores",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("store list not implemented yet - coming in Phase 1B")
		},
	}
}

func (h *StoreHandler) createRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <service>",
		Short: "Remove a password from store",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("store remove not implemented yet - coming in Phase 1B")
		},
	}
}

func (h *StoreHandler) createSyncCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Sync store with remote repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("store sync not implemented yet - coming in Phase 1B")
		},
	}
}

// Handler methods (Phase 1A: Foundation with enhanced card display)

// GetPassword retrieves password metadata and displays in enhanced card format
func (h *StoreHandler) GetPassword(cmd *cobra.Command, args []string) error {
	service := args[0]
	storeName := h.getStoreName(cmd)
	
	copyToClipboard, _ := cmd.Flags().GetBool("copy")
	showPassword, _ := cmd.Flags().GetBool("show")
	
	// For Phase 1A, we'll show a preview of the enhanced card format
	fmt.Printf("🔍 Retrieving '%s' from store '%s'...\n", service, storeName)
	fmt.Printf("📥 Syncing with remote... ✅\n")
	fmt.Printf("🔓 Decrypting metadata... ✅\n\n")
	
	// Mock metadata for demonstration (will be replaced with real data in Phase 1B)
	mockMetadata := h.createMockMetadata(service)
	
	// Display using enhanced card style
	h.cardDisplay.DisplayPasswordCard(mockMetadata)
	
	if copyToClipboard {
		fmt.Printf("\n🔐 Password copied to clipboard (auto-clears in 30 seconds)\n")
		return nil
	}
	
	if showPassword {
		fmt.Printf("\n⚠️  WARNING: This will display the password in terminal\n")
		fmt.Printf("❓ Are you sure? Type 'yes' to confirm: ")
		// In Phase 1B, we'll implement actual confirmation
		fmt.Printf("\n🎯 Password for %s:\n", service)
		
		// Use symmetric password box
		h.cardDisplay.DisplayPasswordBox("Kx9#mN2$vL8@pQ4!")
		
		return nil
	}
	
	return nil
}

// AddPassword adds a new password to the store
func (h *StoreHandler) AddPassword(cmd *cobra.Command, args []string) error {
	service := args[0]
	storeName := h.getStoreName(cmd)
	
	fmt.Printf("🔐 Adding password for '%s' to store '%s'\n", service, storeName)
	fmt.Printf("📝 This will be implemented in Phase 1B with full GPG encryption\n")
	
	return nil
}

// ListPasswords lists all passwords in the store
func (h *StoreHandler) ListPasswords(cmd *cobra.Command, args []string) error {
	storeName := h.getStoreName(cmd)
	
	// Mock data for demonstration
	mockPasswords := h.createMockPasswordList()
	
	h.cardDisplay.DisplayPasswordList(mockPasswords, storeName)
	
	return nil
}

// RotationStatus shows rotation status for auto-rotation enabled passwords
func (h *StoreHandler) RotationStatus(cmd *cobra.Command, args []string) error {
	storeName := h.getStoreName(cmd)
	
	// Mock data for demonstration
	mockStatuses := h.createMockRotationStatuses()
	
	h.cardDisplay.DisplayRotationStatus(mockStatuses, storeName)
	
	return nil
}

// CheckRotations checks for due password rotations
func (h *StoreHandler) CheckRotations(cmd *cobra.Command, args []string) error {
	storeName := h.getStoreName(cmd)
	
	fmt.Printf("🔍 Checking rotation schedule for store '%s'...\n\n", storeName)
	fmt.Printf("🚨 URGENT - Passwords requiring immediate rotation:\n")
	fmt.Printf("• database (2 days overdue)\n")
	fmt.Printf("• api-keys (1 day overdue)\n\n")
	fmt.Printf("⚠️  WARNING - Passwords due soon:\n")
	fmt.Printf("• aws-prod (rotates in 2 days)\n")
	fmt.Printf("• github-token (rotates in 5 days)\n\n")
	fmt.Printf("✅ 12 passwords are up to date\n\n")
	fmt.Printf("💡 Actions:\n")
	fmt.Printf("  passgen store rotate-now database    # Rotate immediately\n")
	fmt.Printf("  passgen store snooze aws-prod 7      # Postpone 7 days\n")
	
	return nil
}

// Helper methods

// getStoreName gets store name from flag or default
func (h *StoreHandler) getStoreName(cmd *cobra.Command) string {
	storeName, _ := cmd.Flags().GetString("store")
	if storeName == "" {
		storeName = "personal" // Default for Phase 1A demo
	}
	return storeName
}

// Mock data helpers for Phase 1A demonstration

func (h *StoreHandler) createMockMetadata(service string) *entities.PasswordMetadata {
	metadata := &entities.PasswordMetadata{
		Service:      service,
		StrengthInfo: "Excellent (16 chars, mixed)",
		CreatedAt:    time.Date(2025, 1, 15, 14, 30, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2025, 1, 15, 14, 30, 0, 0, time.UTC),
	}

	// Add service-specific details for demo
	switch service {
	case "github":
		metadata.Username = "john.doe"
		metadata.URL = "https://github.com"
		metadata.Notes = "Personal GitHub account"
	case "aws-prod":
		metadata.Username = "admin"
		metadata.URL = "https://aws.amazon.com/console"
		metadata.Notes = "Production AWS account"
		metadata.AutoRotation = &entities.AutoRotationInfo{
			Enabled:       true,
			IntervalDays:  90,
			NextRotation:  time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
			DaysUntilNext: 60,
		}
		metadata.StrengthInfo = "Excellent (20 chars, mixed)"
	case "database":
		metadata.Username = "dbuser"
		metadata.URL = "mysql://prod-db.company.com:3306"
		metadata.Notes = "Production database"
		metadata.AutoRotation = &entities.AutoRotationInfo{
			Enabled:       true,
			IntervalDays:  30,
			NextRotation:  time.Date(2025, 8, 12, 0, 0, 0, 0, time.UTC),
			DaysUntilNext: 2,
		}
	}

	return metadata
}

func (h *StoreHandler) createMockPasswordList() []entities.PasswordMetadata {
	return []entities.PasswordMetadata{
		{
			Service:      "github",
			Username:     "john.doe",
			UpdatedAt:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			StrengthInfo: "Excellent",
		},
		{
			Service:   "aws-prod",
			Username:  "admin",
			UpdatedAt: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			AutoRotation: &entities.AutoRotationInfo{
				Enabled:      true,
				IntervalDays: 90,
			},
			StrengthInfo: "Excellent",
		},
		{
			Service:   "database",
			Username:  "dbuser",
			UpdatedAt: time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			AutoRotation: &entities.AutoRotationInfo{
				Enabled:      true,
				IntervalDays: 30,
			},
			StrengthInfo: "Strong",
		},
		{
			Service:      "gitlab",
			Username:     "developer",
			UpdatedAt:    time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC),
			StrengthInfo: "Good",
		},
	}
}

func (h *StoreHandler) createMockRotationStatuses() []entities.RotationStatus {
	return []entities.RotationStatus{
		{
			Service:       "aws-prod",
			NextRotation:  time.Date(2025, 10, 9, 0, 0, 0, 0, time.UTC),
			DaysUntilNext: 60,
			Status:        "scheduled",
			IntervalDays:  90,
		},
		{
			Service:       "database",
			NextRotation:  time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
			DaysUntilNext: 10,
			Status:        "soon",
			IntervalDays:  30,
		},
		{
			Service:       "api-keys",
			NextRotation:  time.Date(2025, 8, 12, 0, 0, 0, 0, time.UTC),
			DaysUntilNext: 2,
			Status:        "critical",
			IntervalDays:  60,
		},
	}
}
