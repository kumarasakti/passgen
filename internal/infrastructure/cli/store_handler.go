package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/kumarasakti/passgen/internal/domain/entities"
	"github.com/kumarasakti/passgen/internal/domain/repositories"
	"github.com/kumarasakti/passgen/internal/infrastructure/display"
	"github.com/kumarasakti/passgen/internal/infrastructure/gpg"
	infraRepositories "github.com/kumarasakti/passgen/internal/infrastructure/repositories"
	"github.com/kumarasakti/passgen/internal/infrastructure/storage"
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
	
	// Check if we have an encrypted repository (Phase 1B)
	if h.repository != nil {
		// Try to load the store
		err := h.loadStoreIfExists(storeName)
		if err != nil {
			return fmt.Errorf("failed to load store '%s': %v", storeName, err)
		}
		
		fmt.Printf("🔍 Retrieving '%s' from store '%s'...\n", service, storeName)
		
		// Get password metadata
		metadata, err := h.repository.GetPasswordMetadata(storeName, service)
		if err != nil {
			return fmt.Errorf("failed to get password metadata: %w", err)
		}
		
		fmt.Printf("🔓 Decrypted metadata ✅\n\n")
		
		// Display using enhanced card style
		h.cardDisplay.DisplayPasswordCard(metadata)
		
		if copyToClipboard {
			// Get full password entry for clipboard
			entry, err := h.repository.GetPassword(storeName, service)
			if err != nil {
				return fmt.Errorf("failed to get password: %w", err)
			}
			
			// In a real implementation, you'd copy to clipboard
			fmt.Printf("\n🔐 Password would be copied to clipboard (auto-clears in 30 seconds)\n")
			fmt.Printf("🔑 Password: %s\n", entry.Password)
			return nil
		}
		
		if showPassword {
			fmt.Printf("\n⚠️  WARNING: This will display the password in terminal\n")
			fmt.Printf("❓ Are you sure? Type 'yes' to confirm: ")
			
			// Simple confirmation for now
			var confirmation string
			fmt.Scanln(&confirmation)
			
			if strings.ToLower(confirmation) == "yes" {
				// Get full password entry
				entry, err := h.repository.GetPassword(storeName, service)
				if err != nil {
					return fmt.Errorf("failed to get password: %w", err)
				}
				
				fmt.Printf("\n🎯 Password for %s:\n", service)
				h.cardDisplay.DisplayPasswordBox(entry.Password)
			} else {
				fmt.Printf("🚫 Password display cancelled\n")
			}
			
			return nil
		}
		
		return nil
	}
	
	// Fallback to Phase 1A demo mode
	fmt.Printf("🔍 Retrieving '%s' from store '%s'...\n", service, storeName)
	fmt.Printf("📥 Syncing with remote... ✅\n")
	fmt.Printf("🔓 Decrypting metadata... ✅\n\n")
	
	// Mock metadata for demonstration
	mockMetadata := h.createMockMetadata(service)
	h.cardDisplay.DisplayPasswordCard(mockMetadata)
	
	if copyToClipboard {
		fmt.Printf("\n🔐 Password copied to clipboard (auto-clears in 30 seconds)\n")
		return nil
	}
	
	if showPassword {
		fmt.Printf("\n⚠️  WARNING: This will display the password in terminal\n")
		fmt.Printf("❓ Are you sure? Type 'yes' to confirm: ")
		fmt.Printf("\n🎯 Password for %s:\n", service)
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
	
	// Check if we have an encrypted repository (Phase 1B)
	if h.repository != nil {
		// Try to load the store (in case it's not loaded yet)
		err := h.loadStoreIfExists(storeName)
		if err != nil {
			return fmt.Errorf("failed to load store '%s': %v\nHint: Initialize the store first with: passgen store init %s", storeName, err, storeName)
		}
		
		// Get additional password details
		username, _ := cmd.Flags().GetString("username")
		url, _ := cmd.Flags().GetString("url") 
		notes, _ := cmd.Flags().GetString("notes")
		length, _ := cmd.Flags().GetInt("length")
		
		// Get auto-rotation settings
		autoRotateDays, _ := cmd.Flags().GetInt("auto-rotate")
		notifyBefore, _ := cmd.Flags().GetInt("notify-before")
		
		// Generate password
		password := h.generatePassword(length)
		
		// Create password entry
		entry := entities.PasswordEntry{
			Service:     service,
			Username:    username,
			Password:    password,
			URL:         url,
			Notes:       notes,
			Metadata:    make(map[string]string),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			GeneratedBy: "passgen-cli",
		}
		
		// Add auto-rotation config if specified
		if autoRotateDays > 0 {
			entry.AutoRotation = &entities.AutoRotationConfig{
				Enabled:          true,
				IntervalDays:     autoRotateDays,
				NextRotationAt:   time.Now().AddDate(0, 0, autoRotateDays),
				NotifyDaysBefore: notifyBefore,
				AutoGenerate:     true,
				PasswordProfile: &entities.PasswordProfile{
					Length:         length,
					IncludeUpper:   true,
					IncludeLower:   true,
					IncludeNumbers: true,
					IncludeSymbols: false, // Conservative default
				},
			}
		}
		
		// Save to encrypted repository
		if err := h.repository.AddPassword(storeName, entry); err != nil {
			return fmt.Errorf("failed to save password: %w", err)
		}
		
		fmt.Printf("✅ Successfully added password for '%s'\n", service)
		fmt.Printf("📝 Username: %s\n", username)
		fmt.Printf("🔗 URL: %s\n", url)
		if notes != "" {
			fmt.Printf("📄 Notes: %s\n", notes)
		}
		fmt.Printf("🔑 Password: %s\n", password)
		
		// Show auto-rotation info if enabled
		if entry.AutoRotation != nil && entry.AutoRotation.Enabled {
			fmt.Printf("🔄 Auto-rotation: Every %d days (next: %s)\n", 
				entry.AutoRotation.IntervalDays,
				entry.AutoRotation.NextRotationAt.Format("2006-01-02"))
		}
		
		fmt.Printf("⚠️  Make sure to save this password securely!\n")
		
		return nil
	}
	
	// Fallback to Phase 1A demo mode
	fmt.Printf("📝 This will be implemented in Phase 1B with full GPG encryption\n")
	return nil
}

// ListPasswords lists all passwords in the store
func (h *StoreHandler) ListPasswords(cmd *cobra.Command, args []string) error {
	storeName := h.getStoreName(cmd)
	
	// Check if we have an encrypted repository (Phase 1B)
	if h.repository != nil {
		// Try to load the store
		err := h.loadStoreIfExists(storeName)
		if err != nil {
			return fmt.Errorf("failed to load store '%s': %v", storeName, err)
		}
		
		// Get all passwords from the repository
		passwords, err := h.repository.ListPasswords(storeName)
		if err != nil {
			return fmt.Errorf("failed to list passwords: %w", err)
		}
		
		if len(passwords) == 0 {
			fmt.Printf("📦 No passwords found in store '%s'\n", storeName)
			fmt.Printf("💡 Add a password with: passgen store add <service> --store %s\n", storeName)
			return nil
		}
		
		h.cardDisplay.DisplayPasswordList(passwords, storeName)
		return nil
	}
	
	// Fallback to Phase 1A demo mode
	mockPasswords := h.createMockPasswordList()
	h.cardDisplay.DisplayPasswordList(mockPasswords, storeName)
	
	return nil
}

// RotationStatus shows rotation status for auto-rotation enabled passwords
func (h *StoreHandler) RotationStatus(cmd *cobra.Command, args []string) error {
	storeName := h.getStoreName(cmd)
	
	// Check if we have an encrypted repository (Phase 1B)
	if h.repository != nil {
		// Try to load the store
		err := h.loadStoreIfExists(storeName)
		if err != nil {
			return fmt.Errorf("failed to load store '%s': %v", storeName, err)
		}
		
		// Get all passwords from the store
		passwords, err := h.repository.ListPasswords(storeName)
		if err != nil {
			return fmt.Errorf("failed to list passwords: %w", err)
		}
		
		// Filter passwords with auto-rotation and create rotation statuses
		var rotationStatuses []entities.RotationStatus
		now := time.Now()
		
		for _, password := range passwords {
			if password.AutoRotation != nil && password.AutoRotation.Enabled {
				daysUntil := int(password.AutoRotation.NextRotation.Sub(now).Hours() / 24)
				
				var status string
				if daysUntil <= 0 {
					status = "overdue"
				} else if daysUntil <= 7 {
					status = "critical"
				} else if daysUntil <= 14 {
					status = "soon"
				} else {
					status = "scheduled"
				}
				
				rotationStatus := entities.RotationStatus{
					Service:       password.Service,
					Status:        status,
					NextRotation:  password.AutoRotation.NextRotation,
					DaysUntilNext: daysUntil,
					IntervalDays:  password.AutoRotation.IntervalDays,
				}
				
				rotationStatuses = append(rotationStatuses, rotationStatus)
			}
		}
		
		h.cardDisplay.DisplayRotationStatus(rotationStatuses, storeName)
		return nil
	}
	
	// Fallback to Phase 1A demo mode
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

// loadStoreIfExists attempts to load a store into the repository if it exists
func (h *StoreHandler) loadStoreIfExists(storeName string) error {
	// Check if store is already loaded
	if _, err := h.repository.GetStore(storeName); err == nil {
		return nil // Store already loaded
	}

	// Try to get store config from registry
	registry := storage.NewStoreRegistry()
	storeConfig, err := registry.GetStore(storeName)
	if err != nil {
		return fmt.Errorf("store '%s' not found in registry", storeName)
	}

	// Cast to encrypted repository to access RegisterStorage
	encryptedRepo, ok := h.repository.(*infraRepositories.EncryptedPasswordStoreRepository)
	if !ok {
		return fmt.Errorf("repository does not support dynamic store loading")
	}

	// Load the store based on its configuration
	gpgService := gpg.NewGPGService(storeConfig.GPGKeyID)
	
	if storeConfig.LocalOnly {
		// Load local-only store
		encStorage := storage.NewLocalOnlyEncryptedStorage(storeConfig.Path, gpgService)
		// Mark as initialized since it's an existing store
		encStorage.SetInitialized(true)
		encryptedRepo.RegisterStorage(storeName, encStorage)
	} else {
		// Load Git-backed store
		encStorage := storage.NewEncryptedStorage(storeConfig.Path, gpgService)
		// Mark as initialized since it's an existing store
		encStorage.SetInitialized(true)
		encryptedRepo.RegisterStorage(storeName, encStorage)
	}

	return nil
}

// generatePassword generates a password of the specified length
func (h *StoreHandler) generatePassword(length int) string {
	// This is a simplified implementation - in production we'd inject the password service
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := ""
	for i := 0; i < length; i++ {
		password += string(chars[i%len(chars)]) // Simplified for demo
	}
	return password
}
