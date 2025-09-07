package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kumarasakti/passgen/internal/infrastructure/display"
	"github.com/kumarasakti/passgen/internal/infrastructure/gpg"
	"github.com/kumarasakti/passgen/internal/infrastructure/repositories"
	"github.com/kumarasakti/passgen/internal/infrastructure/storage"
	"github.com/spf13/cobra"
)

// StoreInitHandler handles store initialization commands
type StoreInitHandler struct {
	displayer *display.CardDisplayer
	repo      *repositories.EncryptedPasswordStoreRepository
	registry  *storage.StoreRegistry
}

// NewStoreInitHandler creates a new store initialization handler
func NewStoreInitHandler(displayer *display.CardDisplayer) *StoreInitHandler {
	repo := repositories.NewEncryptedPasswordStoreRepository()
	registry := storage.NewStoreRegistry()
	return &StoreInitHandler{
		displayer: displayer,
		repo:      repo,
		registry:  registry,
	}
}

// CreateCommands creates the store initialization commands
func (h *StoreInitHandler) CreateCommands() *cobra.Command {
	storeCmd := &cobra.Command{
		Use:   "store",
		Short: "Password store management",
		Long:  "Initialize, configure, and manage encrypted password stores with Git backing",
	}

	// Add subcommands
	storeCmd.AddCommand(h.createInitCommand())
	storeCmd.AddCommand(h.createCloneCommand())
	storeCmd.AddCommand(h.createSyncCommand())
	storeCmd.AddCommand(h.createRemoteCommand())
	storeCmd.AddCommand(h.createInfoCommand())
	storeCmd.AddCommand(h.createSetupGPGCommand())

	return storeCmd
}

// createInitCommand creates the store init command
func (h *StoreInitHandler) createInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [store-name]",
		Short: "Initialize a new password store",
		Long: `Initialize a new encrypted password store for local use.

This command will:
1. Create a new directory for the store
2. Set up GPG encryption
3. Create initial configuration files

Example:
  passgen store init personal

Use --git flag to enable Git backing for synchronization:
  passgen store init personal --git`,
		Args: cobra.ExactArgs(1),
		RunE: h.handleInit,
	}

	// Add flag for Git backing
	cmd.Flags().Bool("git", false, "Enable Git backing for synchronization")
	
	return cmd
}

// createCloneCommand creates the store clone command
func (h *StoreInitHandler) createCloneCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clone [git-url] [store-name]",
		Short: "Clone an existing password store",
		Long: `Clone an existing password store from a Git repository.

You will need:
1. Access to the Git repository
2. The GPG private key used to encrypt the store

Example:
  passgen store clone https://github.com/user/passwords.git personal`,
		Args: cobra.ExactArgs(2),
		RunE: h.handleClone,
	}
}

// createSyncCommand creates the store sync command
func (h *StoreInitHandler) createSyncCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sync [store-name]",
		Short: "Synchronize store with remote repository",
		Long: `Synchronize the password store with its remote Git repository.

This will pull changes from the remote and push any local changes.

Example:
  passgen store sync personal`,
		Args: cobra.ExactArgs(1),
		RunE: h.handleSync,
	}
}

// createRemoteCommand creates the store remote command
func (h *StoreInitHandler) createRemoteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote",
		Short: "Manage store remotes",
		Long:  "Add, remove, or list remote repositories for a store",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "add [store-name] [remote-name] [git-url]",
		Short: "Add a remote repository",
		Args:  cobra.ExactArgs(3),
		RunE:  h.handleRemoteAdd,
	})

	return cmd
}

// createInfoCommand creates the store info command
func (h *StoreInitHandler) createInfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "info [store-name]",
		Short: "Show store information",
		Long:  "Display information about a password store including Git status",
		Args:  cobra.ExactArgs(1),
		RunE:  h.handleInfo,
	}
}

// createSetupGPGCommand creates the GPG setup command
func (h *StoreInitHandler) createSetupGPGCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "setup-gpg",
		Short: "Setup GPG for password stores",
		Long: `Interactive setup for GPG encryption.

This command will:
1. List available GPG keys
2. Help you select or create a key for password encryption
3. Validate the key setup

Example:
  passgen store setup-gpg`,
		RunE: h.handleSetupGPG,
	}
}

// handleInit handles store initialization
func (h *StoreInitHandler) handleInit(cmd *cobra.Command, args []string) error {
	storeName := args[0]
	enableGit, _ := cmd.Flags().GetBool("git")

	fmt.Printf("🔐 Initializing password store: %s\n", storeName)
	if enableGit {
		fmt.Printf("📁 Mode: Local storage with Git backing\n\n")
	} else {
		fmt.Printf("📁 Mode: Local storage only\n\n")
	}

	// Get GPG key
	gpgKeyID, err := h.selectGPGKey()
	if err != nil {
		return fmt.Errorf("failed to setup GPG: %w", err)
	}

	// Create store directory
	homeDir, _ := os.UserHomeDir()
	storePath := filepath.Join(homeDir, ".passgen", "stores")

	if enableGit {
		// Initialize with Git backing
		gpgService := gpg.NewGPGService(gpgKeyID)
		encryptedStorage := storage.NewEncryptedStorage(storePath, gpgService)

		if err := h.repo.InitializeStore(storeName, encryptedStorage); err != nil {
			return fmt.Errorf("failed to initialize store: %w", err)
		}

		// Register store in registry
		storeConfig := storage.StoreConfig{
			Name:      storeName,
			Path:      filepath.Join(storePath, storeName),
			GPGKeyID:  gpgKeyID,
			LocalOnly: false,
		}
		if err := h.registry.RegisterStore(storeConfig); err != nil {
			return fmt.Errorf("failed to register store: %w", err)
		}

		fmt.Printf("✅ Successfully initialized store '%s' with Git backing\n", storeName)
		fmt.Printf("📁 Store location: %s/%s\n", storePath, storeName)
		fmt.Printf("🔑 GPG Key: %s\n", gpgKeyID)
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("  1. Add a remote: passgen store remote add %s origin <git-url>\n", storeName)
		fmt.Printf("  2. Add passwords: passgen store add gmail --store %s\n", storeName)
	} else {
		// Initialize local-only
		if err := h.repo.InitializeLocalStore(storeName, storePath, gpgKeyID); err != nil {
			return fmt.Errorf("failed to initialize local store: %w", err)
		}

		// Register store in registry
		storeConfig := storage.StoreConfig{
			Name:      storeName,
			Path:      filepath.Join(storePath, storeName),
			GPGKeyID:  gpgKeyID,
			LocalOnly: true,
		}
		if err := h.registry.RegisterStore(storeConfig); err != nil {
			return fmt.Errorf("failed to register store: %w", err)
		}

		fmt.Printf("✅ Successfully initialized local store '%s'\n", storeName)
		fmt.Printf("📁 Store location: %s/%s\n", storePath, storeName)
		fmt.Printf("🔑 GPG Key: %s\n", gpgKeyID)
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("  1. Add passwords: passgen store add gmail --store %s\n", storeName)
		fmt.Printf("  2. Enable Git later: passgen store init %s --git (to add Git backing)\n", storeName)
	}

	return nil
}

// handleClone handles store cloning
func (h *StoreInitHandler) handleClone(cmd *cobra.Command, args []string) error {
	gitURL := args[0]
	storeName := args[1]

	fmt.Printf("📥 Cloning password store from: %s\n", gitURL)

	// Get GPG key
	gpgKeyID, err := h.selectGPGKey()
	if err != nil {
		return fmt.Errorf("failed to setup GPG: %w", err)
	}

	// Create store directory
	homeDir, _ := os.UserHomeDir()
	storePath := filepath.Join(homeDir, ".passgen", "stores", storeName)

	// Clone and setup
	// gpgService := gpg.NewGPGService(gpgKeyID)
	// encryptedStorage := storage.NewEncryptedStorage(storePath, gpgService)

	// TODO: Implement actual cloning logic with encryptedStorage
	fmt.Printf("⚠️  Clone functionality not yet implemented\n")
	fmt.Printf("For now, use: git clone %s %s\n", gitURL, storePath)
	fmt.Printf("🔑 GPG Key ready: %s\n", gpgKeyID)

	return nil
}

// handleSync handles store synchronization
func (h *StoreInitHandler) handleSync(cmd *cobra.Command, args []string) error {
	storeName := args[0]

	fmt.Printf("🔄 Synchronizing store: %s\n", storeName)

	if err := h.repo.SyncStore(storeName); err != nil {
		return fmt.Errorf("failed to sync store: %w", err)
	}

	fmt.Printf("✅ Store synchronized successfully\n")
	return nil
}

// handleRemoteAdd handles adding a remote
func (h *StoreInitHandler) handleRemoteAdd(cmd *cobra.Command, args []string) error {
	storeName := args[0]
	remoteName := args[1]
	gitURL := args[2]

	fmt.Printf("🌐 Adding remote '%s' to store '%s'\n", remoteName, storeName)

	if err := h.repo.ConnectRemote(storeName, remoteName, gitURL); err != nil {
		return fmt.Errorf("failed to add remote: %w", err)
	}

	fmt.Printf("✅ Remote added successfully\n")
	return nil
}

// handleInfo handles store information display
func (h *StoreInitHandler) handleInfo(cmd *cobra.Command, args []string) error {
	storeName := args[0]

	info, err := h.repo.GetStoreInfo(storeName)
	if err != nil {
		return fmt.Errorf("failed to get store info: %w", err)
	}

	fmt.Printf("📊 Store Information: %s\n\n", storeName)
	fmt.Printf("📁 Path: %v\n", info["path"])
	fmt.Printf("🌐 Remote: %v\n", info["remote_url"])
	fmt.Printf("🌿 Branch: %v\n", info["branch"])
	fmt.Printf("📝 Status: %v\n", info["status"])
	fmt.Printf("🕐 Last Commit: %v\n", info["last_commit"])

	return nil
}

// handleSetupGPG handles GPG setup
func (h *StoreInitHandler) handleSetupGPG(cmd *cobra.Command, args []string) error {
	fmt.Printf("🔑 GPG Setup for Password Stores\n\n")
	
	gpgService := gpg.NewGPGService("")
	keys, err := gpgService.ListKeys()
	if err != nil {
		return fmt.Errorf("failed to list GPG keys: %w", err)
	}
	
	if len(keys) == 0 {
		fmt.Printf("❌ No GPG keys found.\n\n")
		fmt.Printf("🔐 GPG keys are essential for secure password encryption.\n")
		fmt.Printf("Each password will be encrypted with your GPG key before storage.\n\n")
		fmt.Printf("Would you like to create one now?\n\n")
		fmt.Printf("Options:\n")
		fmt.Printf("  1. Create a new GPG key interactively (Recommended)\n")
		fmt.Printf("  2. Show manual creation instructions\n")
		fmt.Printf("  3. Exit\n\n")
		fmt.Printf("Choose option (1-3): ")
		
		var choice int
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			return fmt.Errorf("invalid input: %w", err)
		}
		
		switch choice {
		case 1:
			_, err := h.createGPGKeyInteractive()
			if err != nil {
				return fmt.Errorf("failed to create GPG key: %w", err)
			}
			fmt.Printf("✅ GPG key created successfully!\n")
			fmt.Printf("🚀 You can now create a store with: passgen store init <store-name>\n")
			return nil
		case 2:
			return h.createGPGKey()
		case 3:
			fmt.Printf("👋 You can run this command again anytime to set up GPG.\n")
			return nil
		default:
			return fmt.Errorf("invalid choice: please select 1, 2, or 3")
		}
	}
	
	fmt.Printf("✅ Found %d GPG key(s):\n\n", len(keys))
	for i, key := range keys {
		status := "✅ Ready to use"
		if key.KeyLength < 2048 && key.KeyType != "22" { // 22 is Ed25519
			status = "⚠️  Key length below recommended 2048 bits"
		}
		
		fmt.Printf("%d. %s\n", i+1, key.UserID)
		fmt.Printf("   ID: %s\n", key.ID)
		fmt.Printf("   Type: %s | Length: %d bits\n", key.KeyType, key.KeyLength)
		fmt.Printf("   Status: %s\n\n", status)
	}
	
	fmt.Printf("💡 These keys are ready for password store encryption.\n")
	fmt.Printf("🚀 You can now create a store with: passgen store init <store-name>\n\n")
	fmt.Printf("Would you like to create an additional GPG key? (y/N): ")
	
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
		fmt.Printf("\n")
		_, err := h.createGPGKeyInteractive()
		if err != nil {
			return fmt.Errorf("failed to create additional GPG key: %w", err)
		}
		fmt.Printf("✅ Additional GPG key created successfully!\n")
	}
	
	return nil
}// selectGPGKey interactively selects a GPG key
func (h *StoreInitHandler) selectGPGKey() (string, error) {
	gpgService := gpg.NewGPGService("")
	keys, err := gpgService.ListKeys()
	if err != nil {
		return "", fmt.Errorf("failed to list GPG keys: %w", err)
	}

	if len(keys) == 0 {
		fmt.Printf("❌ No GPG keys found.\n\n")
		fmt.Printf("🔑 GPG keys are required for secure password encryption.\n")
		fmt.Printf("Would you like to create one now?\n\n")
		fmt.Printf("Options:\n")
		fmt.Printf("  1. Create a new GPG key automatically (Recommended)\n")
		fmt.Printf("  2. Create manually with 'gpg --full-generate-key'\n")
		fmt.Printf("  3. Exit and create key later\n\n")
		fmt.Printf("Choose option (1-3): ")
		
		var choice int
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			return "", fmt.Errorf("invalid input: %w", err)
		}
		
		switch choice {
		case 1:
			// Create GPG key automatically
			keyID, err := h.createGPGKeyInteractive()
			if err != nil {
				return "", fmt.Errorf("failed to create GPG key: %w", err)
			}
			fmt.Printf("✅ GPG key created successfully!\n")
			return keyID, nil
		case 2:
			err := h.createGPGKey()
			if err != nil {
				return "", err
			}
			// After manual creation, we need to detect the new key
			keysAfter, err := gpgService.ListKeys()
			if err != nil {
				return "", fmt.Errorf("failed to list keys after creation: %w", err)
			}
			
			if len(keysAfter) > 0 {
				// If there's exactly one key now, use it
				if len(keysAfter) == 1 {
					return keysAfter[0].ID, nil
				}
				// If multiple keys, ask user to select
				fmt.Printf("\n🔑 Please select your newly created key:\n\n")
				for i, key := range keysAfter {
					fmt.Printf("  %d. %s\n", i+1, key.UserID)
					fmt.Printf("     ID: %s | Type: %s | Length: %d bits\n\n", key.ID, key.KeyType, key.KeyLength)
				}
				fmt.Printf("Which is your new key? (1-%d): ", len(keysAfter))
				
				var choice int
				_, err = fmt.Scanf("%d", &choice)
				if err != nil || choice < 1 || choice > len(keysAfter) {
					return "", fmt.Errorf("invalid selection")
				}
				
				return keysAfter[choice-1].ID, nil
			}
			return "", fmt.Errorf("no keys found after creation")
		case 3:
			return "", fmt.Errorf("GPG key required - please create one and try again")
		default:
			return "", fmt.Errorf("invalid choice: please select 1, 2, or 3")
		}
	}

	if len(keys) == 1 {
		fmt.Printf("🔑 Found GPG key: %s\n", keys[0].UserID)
		fmt.Printf("   Using this key for encryption.\n")
		return keys[0].ID, nil
	}

	// Multiple keys available - show selection
	fmt.Printf("🔑 Multiple GPG keys found:\n\n")
	for i, key := range keys {
		status := "✅ Ready"
		if key.KeyLength < 2048 && key.KeyType != "22" { // 22 is Ed25519
			status = "⚠️  Short key"
		}
		fmt.Printf("  %d. %s\n", i+1, key.UserID)
		fmt.Printf("     ID: %s | Type: %s | Length: %d bits | %s\n\n", key.ID, key.KeyType, key.KeyLength, status)
	}
	
	fmt.Printf("  %d. 🆕 Create a new GPG key\n\n", len(keys)+1)
	
	fmt.Printf("Which option would you like? (1-%d): ", len(keys)+1)
	
	// Read user input for key selection
	var choice int
	_, err = fmt.Scanf("%d", &choice)
	if err != nil || choice < 1 || choice > len(keys)+1 {
		return "", fmt.Errorf("invalid selection: please choose a number between 1 and %d", len(keys)+1)
	}
	
	// Check if user wants to create a new key
	if choice == len(keys)+1 {
		fmt.Printf("\n")
		keyID, err := h.createGPGKeyInteractive()
		if err != nil {
			return "", fmt.Errorf("failed to create GPG key: %w", err)
		}
		fmt.Printf("✅ New GPG key created and selected!\n")
		return keyID, nil
	}
	
	selectedKey := keys[choice-1]
	fmt.Printf("🔑 Selected: %s\n", selectedKey.UserID)
	
	return selectedKey.ID, nil
}

// createGPGKeyInteractive creates a GPG key using GPG's interactive prompts
func (h *StoreInitHandler) createGPGKeyInteractive() (string, error) {
	fmt.Printf("🔑 Creating a new GPG key for passgen...\n\n")
	fmt.Printf("I'll launch GPG's interactive key generation process.\n")
	fmt.Printf("GPG will guide you through creating a secure key.\n\n")
	fmt.Printf("💡 Recommendations for the prompts you'll see:\n")
	fmt.Printf("   • Key type: (1) RSA and RSA (default) or (9) ECC and ECC\n")
	fmt.Printf("   • Key size: 4096 bits (for RSA) or accept default (for ECC)\n")
	fmt.Printf("   • Expiration: 2y (2 years) - recommended\n")
	fmt.Printf("   • Use a strong passphrase you'll remember\n\n")
	fmt.Printf("🚀 Starting GPG key generation...\n")
	fmt.Printf("=====================================\n\n")
	
	// Get the number of keys before creation to identify the new one
	gpgService := gpg.NewGPGService("")
	keysBefore, err := gpgService.ListKeys()
	if err != nil {
		// If we can't list keys, continue anyway
		keysBefore = []gpg.GPGKey{}
	}
	
	// Launch GPG's interactive key generation
	cmd := exec.Command("gpg", "--full-generate-key")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("GPG key generation failed: %w", err)
	}
	
	fmt.Printf("\n=====================================\n")
	fmt.Printf("✅ GPG key generation completed!\n\n")
	
	// Find the newly created key
	keysAfter, err := gpgService.ListKeys()
	if err != nil {
		return "", fmt.Errorf("failed to list keys after creation: %w", err)
	}
	
	// Find the new key by comparing before and after
	var newKey *gpg.GPGKey
	for _, keyAfter := range keysAfter {
		found := false
		for _, keyBefore := range keysBefore {
			if keyAfter.ID == keyBefore.ID {
				found = true
				break
			}
		}
		if !found {
			newKey = &keyAfter
			break
		}
	}
	
	if newKey == nil {
		// Fallback: if we can't detect the new key, ask user to select
		if len(keysAfter) > len(keysBefore) {
			fmt.Printf("� Key created successfully! Please select your new key:\n\n")
			for i, key := range keysAfter {
				fmt.Printf("  %d. %s\n", i+1, key.UserID)
				fmt.Printf("     ID: %s | Type: %s | Length: %d bits\n\n", key.ID, key.KeyType, key.KeyLength)
			}
			fmt.Printf("Which is your new key? (1-%d): ", len(keysAfter))
			
			var choice int
			_, err = fmt.Scanf("%d", &choice)
			if err != nil || choice < 1 || choice > len(keysAfter) {
				return "", fmt.Errorf("invalid selection")
			}
			
			newKey = &keysAfter[choice-1]
		} else {
			return "", fmt.Errorf("could not detect newly created key")
		}
	}
	
	fmt.Printf("🔑 New GPG key ready:\n")
	fmt.Printf("   Name: %s\n", newKey.UserID)
	fmt.Printf("   Key ID: %s\n", newKey.ID)
	fmt.Printf("   Type: %s | Length: %d bits\n\n", newKey.KeyType, newKey.KeyLength)
	
	return newKey.ID, nil
}

// createGPGKey guides the user through GPG key creation using GPG's interactive prompts
func (h *StoreInitHandler) createGPGKey() error {
	fmt.Printf("🔑 Creating a new GPG key for passgen...\n\n")
	fmt.Printf("I'll launch GPG's interactive key generation process.\n")
	fmt.Printf("GPG will guide you through creating a secure key.\n\n")
	fmt.Printf("💡 Recommendations for the prompts you'll see:\n")
	fmt.Printf("   • Key type: (1) RSA and RSA (default) or (9) ECC and ECC\n")
	fmt.Printf("   • Key size: 4096 bits (for RSA) or accept default (for ECC)\n")
	fmt.Printf("   • Expiration: 2y (2 years) - recommended\n")
	fmt.Printf("   • Use a strong passphrase you'll remember\n\n")
	fmt.Printf("🚀 Starting GPG key generation...\n")
	fmt.Printf("=====================================\n\n")
	
	// Launch GPG's interactive key generation
	cmd := exec.Command("gpg", "--full-generate-key")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("GPG key generation failed: %w", err)
	}
	
	fmt.Printf("\n=====================================\n")
	fmt.Printf("✅ GPG key generation completed!\n")
	fmt.Printf("🚀 You can now run 'passgen store init <store-name>' to create your password store.\n")
	
	return nil
}
