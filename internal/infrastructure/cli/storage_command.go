package cli

import (
	"fmt"
	"os"

	"github.com/kumarasakti/passgen/internal/infrastructure/config"
	"github.com/kumarasakti/passgen/internal/infrastructure/storage/backends"
	"github.com/spf13/cobra"
)

// Enables comprehensive storage backend management and configuration
func NewStorageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage",
		Short: "Manage storage backends",
		Long:  "Configure and test storage backends like local, R2, etc.",
	}

	cmd.AddCommand(
		newStorageConfigCommand(),
		newStorageTestCommand(),
		newStorageListCommand(),
	)

	return cmd
}

// Provides interactive storage backend configuration with validation
func newStorageConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config [backend]",
		Short: "Configure storage backend",
		Long:  "Configure storage backend settings (local, r2)",
		Args:  cobra.ExactArgs(1),
		RunE:  runStorageConfig,
	}

	// R2 specific flags
	cmd.Flags().String("account-id", "", "Cloudflare R2 Account ID")
	cmd.Flags().String("access-key", "", "R2 Access Key ID")
	cmd.Flags().String("secret-key", "", "R2 Secret Access Key")
	cmd.Flags().String("bucket", "", "R2 Bucket Name")
	cmd.Flags().String("region", "auto", "R2 Region")

	// Local specific flags
	cmd.Flags().String("base-path", "~/.passgen/stores", "Base path for local storage")

	return cmd
}

// Validates storage backend connectivity and configuration accuracy
func newStorageTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test [backend]",
		Short: "Test storage backend connection",
		Long:  "Test if the storage backend is properly configured and accessible",
		Args:  cobra.ExactArgs(1),
		RunE:  runStorageTest,
	}

	return cmd
}

// Displays overview of all configured storage backends with their settings
func newStorageListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List configured storage backends",
		Long:  "List all configured storage backends and their settings",
		RunE:  runStorageList,
	}

	return cmd
}

// Processes and saves storage backend configuration based on user input
func runStorageConfig(cmd *cobra.Command, args []string) error {
	backend := args[0]

	// Load existing config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	switch backend {
	case "r2":
		return configureR2Backend(cmd, cfg)
	case "local":
		return configureLocalBackend(cmd, cfg)
	default:
		return fmt.Errorf("unsupported backend: %s (supported: local, r2)", backend)
	}
}

// Sets up Cloudflare R2 cloud storage with credentials and bucket configuration
func configureR2Backend(cmd *cobra.Command, cfg *config.PassgenConfig) error {
	accountID, _ := cmd.Flags().GetString("account-id")
	accessKey, _ := cmd.Flags().GetString("access-key")
	secretKey, _ := cmd.Flags().GetString("secret-key")
	bucket, _ := cmd.Flags().GetString("bucket")
	region, _ := cmd.Flags().GetString("region")

	// Validate required flags
	if accountID == "" || accessKey == "" || secretKey == "" || bucket == "" {
		return fmt.Errorf("R2 configuration requires: --account-id, --access-key, --secret-key, --bucket")
	}

	// Update config
	cfg.Storage.Backend = "r2"
	cfg.Storage.Settings = map[string]string{
		"account_id":        accountID,
		"access_key_id":     accessKey,
		"secret_access_key": secretKey,
		"bucket_name":       bucket,
		"region":            region,
	}

	// Save config
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("✅ R2 backend configured successfully!")
	fmt.Printf("   Account ID: %s\n", accountID)
	fmt.Printf("   Bucket: %s\n", bucket)
	fmt.Printf("   Region: %s\n", region)

	return nil
}

// Establishes local file system storage with base directory configuration
func configureLocalBackend(cmd *cobra.Command, cfg *config.PassgenConfig) error {
	basePath, _ := cmd.Flags().GetString("base-path")

	// Update config
	cfg.Storage.Backend = "local"
	cfg.Storage.Settings = map[string]string{
		"base_path": basePath,
	}

	// Save config
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("✅ Local backend configured successfully!")
	fmt.Printf("   Base Path: %s\n", basePath)

	return nil
}

// Validates storage backend functionality and reports connection status
func runStorageTest(cmd *cobra.Command, args []string) error {
	backend := args[0]

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	switch backend {
	case "r2":
		return testR2Backend(cfg)
	case "local":
		return testLocalBackend(cfg)
	default:
		return fmt.Errorf("unsupported backend: %s (supported: local, r2)", backend)
	}
}

// testR2Backend tests R2 backend connection
func testR2Backend(cfg *config.PassgenConfig) error {
	if cfg.Storage.Backend != "r2" {
		return fmt.Errorf("R2 backend not configured. Run: passgen storage config r2 --help")
	}

	// Create R2 config from settings
	r2Config := backends.R2Config{
		AccountID:       cfg.Storage.Settings["account_id"],
		AccessKeyID:     cfg.Storage.Settings["access_key_id"],
		SecretAccessKey: cfg.Storage.Settings["secret_access_key"],
		BucketName:      cfg.Storage.Settings["bucket_name"],
		Region:          cfg.Storage.Settings["region"],
	}

	// Test connection
	fmt.Print("🔍 Testing R2 connection... ")

	backend, err := backends.NewR2Backend(r2Config, "test")
	if err != nil {
		fmt.Println("❌ Failed")
		return fmt.Errorf("failed to create R2 backend: %w", err)
	}

	if err := backend.TestConnection(); err != nil {
		fmt.Println("❌ Failed")
		return fmt.Errorf("R2 connection test failed: %w", err)
	}

	fmt.Println("✅ Success!")
	fmt.Printf("   Connected to bucket: %s\n", r2Config.BucketName)
	fmt.Printf("   Account ID: %s\n", r2Config.AccountID)

	return nil
}

// testLocalBackend tests local backend
func testLocalBackend(cfg *config.PassgenConfig) error {
	if cfg.Storage.Backend != "local" {
		return fmt.Errorf("local backend not configured. Run: passgen storage config local --help")
	}

	basePath := cfg.Storage.Settings["base_path"]
	if basePath == "" {
		return fmt.Errorf("base_path not configured for local backend")
	}

	// Expand path
	expandedPath := expandPath(basePath)

	fmt.Print("🔍 Testing local storage... ")

	// Test if we can create directory and write files
	testDir := expandedPath + "/test"
	if err := os.MkdirAll(testDir, 0700); err != nil {
		fmt.Println("❌ Failed")
		return fmt.Errorf("cannot create test directory: %w", err)
	}

	testFile := testDir + "/test.txt"
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		fmt.Println("❌ Failed")
		return fmt.Errorf("cannot write test file: %w", err)
	}

	// Clean up
	os.RemoveAll(testDir)

	fmt.Println("✅ Success!")
	fmt.Printf("   Base Path: %s\n", expandedPath)

	return nil
}

// runStorageList lists configured storage backends
func runStorageList(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("🗄️  Storage Configuration:")
	fmt.Printf("   Default Backend: %s\n", cfg.Storage.Backend)

	fmt.Println("\n📋 Backend Settings:")
	for key, value := range cfg.Storage.Settings {
		// Mask sensitive values
		if key == "secret_access_key" {
			value = "***masked***"
		}
		fmt.Printf("   %s: %s\n", key, value)
	}

	if len(cfg.Storage.Stores) > 0 {
		fmt.Println("\n🏪 Configured Stores:")
		for name, store := range cfg.Storage.Stores {
			fmt.Printf("   %s: %s backend\n", name, store.Backend)
		}
	}

	return nil
}

// expandPath expands ~ to user home directory
func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		homeDir, _ := os.UserHomeDir()
		return homeDir + path[1:]
	}
	return path
}
