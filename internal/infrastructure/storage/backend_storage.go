package storage

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kumarasakti/passgen/internal/domain/entities"
	"github.com/kumarasakti/passgen/internal/domain/interfaces"
	"github.com/kumarasakti/passgen/internal/infrastructure/gpg"
	"github.com/kumarasakti/passgen/internal/infrastructure/storage/backends"
)

// Provides encrypted password storage with flexible backend support (local, cloud, etc.)
type BackendStorage struct {
	backend     interfaces.StorageBackend
	gpgService  *gpg.GPGService
	storeName   string
	initialized bool
}

// Initializes storage with specified backend and GPG encryption capabilities
func NewBackendStorage(backend interfaces.StorageBackend, gpgService *gpg.GPGService, storeName string) *BackendStorage {
	return &BackendStorage{
		backend:    backend,
		gpgService: gpgService,
		storeName:  storeName,
	}
}

// Configures storage with local file system backend for encrypted password storage
func NewLocalBackendStorage(basePath string, gpgService *gpg.GPGService, storeName string) *BackendStorage {
	backend := backends.NewLocalBackend(basePath)
	return NewBackendStorage(backend, gpgService, storeName)
}

// Configures storage with Cloudflare R2 cloud backend for encrypted password storage
func NewR2BackendStorage(config backends.R2Config, storePrefix string, gpgService *gpg.GPGService, storeName string) (*BackendStorage, error) {
	backend, err := backends.NewR2Backend(config, storePrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to create R2 backend: %w", err)
	}

	return NewBackendStorage(backend, gpgService, storeName), nil
}

// InitializeStore initializes a new password store
func (bs *BackendStorage) InitializeStore(storeName string) error {
	err := bs.backend.Initialize(storeName)
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}

	bs.initialized = true
	return nil
}

// SetInitialized sets the initialized flag for existing stores
func (bs *BackendStorage) SetInitialized(initialized bool) {
	bs.initialized = initialized
}

// Verifies store initialization status through backend validation
func (bs *BackendStorage) IsInitialized() (bool, error) {
	return bs.backend.IsInitialized(bs.storeName)
}

// SavePassword saves an encrypted password entry
func (bs *BackendStorage) SavePassword(entry entities.PasswordEntry) error {
	if !bs.initialized {
		return fmt.Errorf("store not initialized")
	}

	// Create stored entry using the type from encrypted_storage
	storedEntry := StoredPasswordEntry{
		Service:         entry.Service,
		Username:        entry.Username,
		Password:        entry.Password,
		URL:             entry.URL,
		Notes:           entry.Notes,
		Metadata:        entry.Metadata,
		CreatedAt:       entry.CreatedAt,
		UpdatedAt:       entry.UpdatedAt,
		GeneratedBy:     entry.GeneratedBy,
		AutoRotation:    entry.AutoRotation,
		RotationHistory: entry.RotationHistory,
	}

	// Serialize to JSON
	jsonData, err := json.MarshalIndent(storedEntry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize password entry: %w", err)
	}

	// Encrypt the JSON data
	encryptedData, err := bs.gpgService.Encrypt(jsonData, "")
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	// Save to backend
	fileName := bs.sanitizeFileName(entry.Service) + ".gpg"
	err = bs.backend.SaveFile(fileName, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to save password file: %w", err)
	}

	return nil
}

// LoadPassword loads and decrypts a password entry
func (bs *BackendStorage) LoadPassword(name string) (*entities.PasswordEntry, error) {
	if !bs.initialized {
		return nil, fmt.Errorf("store not initialized")
	}

	fileName := bs.sanitizeFileName(name) + ".gpg"

	// Load from backend
	encryptedData, err := bs.backend.LoadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to load password file: %w", err)
	}

	// Decrypt the data
	decryptedData, err := bs.gpgService.Decrypt(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}

	// Parse JSON
	var storedEntry StoredPasswordEntry
	if err := json.Unmarshal(decryptedData, &storedEntry); err != nil {
		return nil, fmt.Errorf("failed to parse password data: %w", err)
	}

	return &entities.PasswordEntry{
		Service:         storedEntry.Service,
		Username:        storedEntry.Username,
		Password:        storedEntry.Password,
		URL:             storedEntry.URL,
		Notes:           storedEntry.Notes,
		Metadata:        storedEntry.Metadata,
		CreatedAt:       storedEntry.CreatedAt,
		UpdatedAt:       storedEntry.UpdatedAt,
		GeneratedBy:     storedEntry.GeneratedBy,
		AutoRotation:    storedEntry.AutoRotation,
		RotationHistory: storedEntry.RotationHistory,
	}, nil
}

// Provides comprehensive overview of stored passwords without exposing sensitive data
func (bs *BackendStorage) ListPasswords() ([]entities.PasswordMetadata, error) {
	if !bs.initialized {
		return nil, fmt.Errorf("store not initialized")
	}

	var passwords []entities.PasswordMetadata

	// List all .gpg files
	files, err := bs.backend.ListFiles("")
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	for _, file := range files {
		// Skip non-password files
		if !strings.HasSuffix(file, ".gpg") || strings.HasPrefix(file, ".passgen-") {
			continue
		}

		// Extract name from filename
		name := bs.unsanitizeFileName(strings.TrimSuffix(file, ".gpg"))

		// Load the password entry to get metadata
		entry, err := bs.LoadPassword(name)
		if err != nil {
			// Skip entries that can't be loaded
			continue
		}

		// Convert to PasswordMetadata
		metadata := entities.PasswordMetadata{
			Service:   entry.Service,
			Username:  entry.Username,
			URL:       entry.URL,
			Notes:     entry.Notes,
			CreatedAt: entry.CreatedAt,
			UpdatedAt: entry.UpdatedAt,
		}

		// Add auto-rotation info if present
		if entry.AutoRotation != nil {
			metadata.AutoRotation = &entities.AutoRotationInfo{
				Enabled:      true,
				IntervalDays: entry.AutoRotation.IntervalDays,
				// Note: NextRotation would need to be calculated based on interval
			}
		}

		passwords = append(passwords, metadata)
	}

	return passwords, nil
}

// DeletePassword removes a password entry
func (bs *BackendStorage) DeletePassword(name string) error {
	if !bs.initialized {
		return fmt.Errorf("store not initialized")
	}

	fileName := bs.sanitizeFileName(name) + ".gpg"

	// Check if file exists
	exists, err := bs.backend.FileExists(fileName)
	if err != nil {
		return fmt.Errorf("failed to check if password exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("password '%s' not found", name)
	}

	// Delete from backend
	err = bs.backend.DeleteFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to delete password: %w", err)
	}

	return nil
}

// Sync synchronizes with the backend (if supported)
func (bs *BackendStorage) Sync() error {
	return bs.backend.Sync()
}

// Provides comprehensive backend configuration and connection details
func (bs *BackendStorage) GetBackendInfo() map[string]string {
	info := bs.backend.GetConnectionInfo()
	info["store_name"] = bs.storeName
	return info
}

// Identifies passwords exceeding their rotation schedule for security maintenance
func (bs *BackendStorage) GetPasswordsNeedingRotation() ([]entities.PasswordMetadata, error) {
	if !bs.initialized {
		return nil, fmt.Errorf("store not initialized")
	}

	passwords, err := bs.ListPasswords()
	if err != nil {
		return nil, fmt.Errorf("failed to list passwords: %w", err)
	}

	var needingRotation []entities.PasswordMetadata
	now := time.Now()

	for _, password := range passwords {
		if password.AutoRotation != nil && password.AutoRotation.Enabled {
			// Check if password needs rotation
			if password.AutoRotation.NextRotation.Before(now) || password.AutoRotation.NextRotation.Equal(now) {
				needingRotation = append(needingRotation, password)
			}
		}
	}

	return needingRotation, nil
}

// UpdateAutoRotation updates the auto-rotation configuration for a password
func (bs *BackendStorage) UpdateAutoRotation(service string, config entities.AutoRotationConfig) error {
	if !bs.initialized {
		return fmt.Errorf("store not initialized")
	}

	// Load existing password
	entry, err := bs.LoadPassword(service)
	if err != nil {
		return fmt.Errorf("failed to load password: %w", err)
	}

	// Update auto-rotation config
	entry.AutoRotation = &config
	entry.UpdatedAt = time.Now()

	// Save updated entry
	return bs.SavePassword(*entry)
}

// AddRotationRecord adds a rotation record to a password's history
func (bs *BackendStorage) AddRotationRecord(service string, record entities.RotationRecord) error {
	if !bs.initialized {
		return fmt.Errorf("store not initialized")
	}

	// Load existing password
	entry, err := bs.LoadPassword(service)
	if err != nil {
		return fmt.Errorf("failed to load password: %w", err)
	}

	// Add rotation record
	if entry.RotationHistory == nil {
		entry.RotationHistory = []entities.RotationRecord{}
	}
	entry.RotationHistory = append(entry.RotationHistory, record)
	entry.UpdatedAt = time.Now()

	// Save updated entry
	return bs.SavePassword(*entry)
}

// sanitizeFileName converts a password name to a safe filename
func (bs *BackendStorage) sanitizeFileName(name string) string {
	// Replace unsafe characters with underscores
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		" ", "_",
	)
	return replacer.Replace(name)
}

// unsanitizeFileName converts a filename back to original name (basic implementation)
func (bs *BackendStorage) unsanitizeFileName(filename string) string {
	// This is a basic implementation - in real use, you might want to store
	// the original name in metadata to avoid this conversion
	return strings.ReplaceAll(filename, "_", " ")
}
