package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/kumarasakti/passgen/internal/infrastructure/gpg"
	"github.com/kumarasakti/passgen/internal/infrastructure/repositories"
	"github.com/kumarasakti/passgen/internal/infrastructure/storage"
	"github.com/kumarasakti/passgen/internal/infrastructure/display"
)

// StoreInitHandler handles store initialization commands
type StoreInitHandler struct {
	displayer *display.CardDisplayer
	repo      *repositories.EncryptedPasswordStoreRepository
}

// NewStoreInitHandler creates a new store initialization handler
func NewStoreInitHandler(displayer *display.CardDisplayer) *StoreInitHandler {
	repo := repositories.NewEncryptedPasswordStoreRepository()
	return &StoreInitHandler{
		displayer: displayer,
		repo:      repo,
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
	return &cobra.Command{
		Use:   "init [store-name]",
		Short: "Initialize a new password store",
		Long: `Initialize a new encrypted password store with Git backing.

This command will:
1. Create a new directory for the store
2. Initialize a Git repository
3. Set up GPG encryption
4. Create initial configuration files

Example:
  passgen store init personal`,
		Args: cobra.ExactArgs(1),
		RunE: h.handleInit,
	}
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
	
	fmt.Printf("🔐 Initializing password store: %s\n\n", storeName)

	// Get GPG key
	gpgKeyID, err := h.selectGPGKey()
	if err != nil {
		return fmt.Errorf("failed to setup GPG: %w", err)
	}

	// Create store directory
	homeDir, _ := os.UserHomeDir()
	storePath := filepath.Join(homeDir, ".passgen", "stores", storeName)
	
	// Initialize storage
	gpgService := gpg.NewGPGService(gpgKeyID)
	encryptedStorage := storage.NewEncryptedStorage(storePath, gpgService)
	
	if err := h.repo.InitializeStore(storeName, encryptedStorage); err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}

	fmt.Printf("✅ Successfully initialized store '%s'\n", storeName)
	fmt.Printf("📁 Store location: %s\n", storePath)
	fmt.Printf("🔑 GPG Key: %s\n", gpgKeyID)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  1. Add a remote: passgen store remote add %s origin <git-url>\n", storeName)
	fmt.Printf("  2. Add passwords: passgen add %s <service>\n", storeName)
	
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
		fmt.Printf("Please create a GPG key first:\n")
		fmt.Printf("  gpg --full-generate-key\n\n")
		fmt.Printf("Then run this command again.\n")
		return nil
	}
	
	fmt.Printf("Available GPG Keys:\n\n")
	for i, key := range keys {
		fmt.Printf("%d. %s\n", i+1, key.UserID)
		fmt.Printf("   ID: %s\n", key.ID)
		fmt.Printf("   Type: %s\n", key.KeyType)
		fmt.Printf("   Length: %d bits\n\n", key.KeyLength)
	}
	
	fmt.Printf("Select a key by number for password store encryption.\n")
	fmt.Printf("The selected key will be used to encrypt all passwords in your stores.\n")
	
	return nil
}

// selectGPGKey interactively selects a GPG key
func (h *StoreInitHandler) selectGPGKey() (string, error) {
	gpgService := gpg.NewGPGService("")
	keys, err := gpgService.ListKeys()
	if err != nil {
		return "", fmt.Errorf("failed to list GPG keys: %w", err)
	}
	
	if len(keys) == 0 {
		return "", fmt.Errorf("no GPG keys found - please create one with 'gpg --full-generate-key'")
	}
	
	if len(keys) == 1 {
		fmt.Printf("🔑 Using GPG key: %s\n", keys[0].UserID)
		return keys[0].ID, nil
	}
	
	// For now, use the first key - in a real implementation, you'd prompt the user
	fmt.Printf("🔑 Using GPG key: %s\n", keys[0].UserID)
	return keys[0].ID, nil
}
