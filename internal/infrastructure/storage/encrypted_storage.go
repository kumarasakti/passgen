package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kumarasakti/passgen/internal/domain/entities"
	"github.com/kumarasakti/passgen/internal/infrastructure/git"
	"github.com/kumarasakti/passgen/internal/infrastructure/gpg"
)

// EncryptedStorage handles encrypted password storage with Git backing
type EncryptedStorage struct {
	storePath   string
	gpgService  *gpg.GPGService
	gitService  *git.GitService
	initialized bool
}

// NewEncryptedStorage creates a new encrypted storage instance
func NewEncryptedStorage(storePath string, gpgService *gpg.GPGService) *EncryptedStorage {
	gitService := git.NewGitService(storePath)
	
	return &EncryptedStorage{
		storePath:  storePath,
		gpgService: gpgService,
		gitService: gitService,
	}
}

// StoredPasswordEntry represents the stored format of a password entry
type StoredPasswordEntry struct {
	Service         string                     `json:"service"`
	Username        string                     `json:"username,omitempty"`
	Password        string                     `json:"password"`
	URL             string                     `json:"url,omitempty"`
	Notes           string                     `json:"notes,omitempty"`
	Metadata        map[string]string          `json:"metadata"`
	CreatedAt       time.Time                  `json:"created_at"`
	UpdatedAt       time.Time                  `json:"updated_at"`
	GeneratedBy     string                     `json:"generated_by"`
	AutoRotation    *entities.AutoRotationConfig `json:"auto_rotation,omitempty"`
	RotationHistory []entities.RotationRecord    `json:"rotation_history,omitempty"`
}

// InitializeStore initializes a new password store
func (es *EncryptedStorage) InitializeStore(storeName string) error {
	storeDir := filepath.Join(es.storePath, storeName)
	
	// Create store directory
	if err := os.MkdirAll(storeDir, 0700); err != nil {
		return fmt.Errorf("failed to create store directory: %w", err)
	}

	// Update git service path to store directory
	es.gitService = git.NewGitService(storeDir)

	// Initialize Git repository
	if !es.gitService.IsRepository() {
		if err := es.gitService.InitializeRepository(); err != nil {
			return fmt.Errorf("failed to initialize git repository: %w", err)
		}
	}

	// Create store metadata
	storeMetadata := entities.PasswordStore{
		Name:      storeName,
		LocalPath: storeDir,
		CreatedAt: time.Now(),
	}

	metadataPath := filepath.Join(storeDir, ".passgen-store")
	if err := es.saveStoreMetadata(metadataPath, storeMetadata); err != nil {
		return fmt.Errorf("failed to save store metadata: %w", err)
	}

	// Add and commit initial files
	if err := es.gitService.AddFiles([]string{"."}); err != nil {
		return fmt.Errorf("failed to add files to git: %w", err)
	}

	if err := es.gitService.Commit("Initialize password store: " + storeName); err != nil {
		return fmt.Errorf("failed to commit initial files: %w", err)
	}

	es.initialized = true
	return nil
}

// ConnectRemote connects the store to a remote Git repository
func (es *EncryptedStorage) ConnectRemote(remoteName, remoteURL string) error {
	if !es.initialized {
		return fmt.Errorf("store not initialized")
	}

	if err := es.gitService.AddRemote(remoteName, remoteURL); err != nil {
		return fmt.Errorf("failed to add remote: %w", err)
	}

	return nil
}

// SavePassword saves an encrypted password entry
func (es *EncryptedStorage) SavePassword(entry entities.PasswordEntry) error {
	if !es.initialized {
		return fmt.Errorf("store not initialized")
	}

	// Create stored entry
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
	encryptedData, err := es.gpgService.Encrypt(jsonData, "")
	if err != nil {
		return fmt.Errorf("failed to encrypt password entry: %w", err)
	}

	// Save to file
	fileName := es.sanitizeFileName(entry.Service) + ".gpg"
	filePath := filepath.Join(es.storePath, fileName)
	
	if err := os.WriteFile(filePath, encryptedData, 0600); err != nil {
		return fmt.Errorf("failed to write encrypted file: %w", err)
	}

	// Add to git and commit
	if err := es.gitService.AddFiles([]string{fileName}); err != nil {
		return fmt.Errorf("failed to add file to git: %w", err)
	}

	commitMsg := fmt.Sprintf("Add password entry: %s", entry.Service)
	if err := es.gitService.Commit(commitMsg); err != nil {
		return fmt.Errorf("failed to commit password entry: %w", err)
	}

	return nil
}

// LoadPassword loads and decrypts a password entry
func (es *EncryptedStorage) LoadPassword(name string) (*entities.PasswordEntry, error) {
	if !es.initialized {
		return nil, fmt.Errorf("store not initialized")
	}

	fileName := es.sanitizeFileName(name) + ".gpg"
	filePath := filepath.Join(es.storePath, fileName)

	// Read encrypted file
	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("password entry '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to read encrypted file: %w", err)
	}

	// Decrypt the data
	decryptedData, err := es.gpgService.Decrypt(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password entry: %w", err)
	}

	// Parse JSON
	var storedEntry StoredPasswordEntry
	if err := json.Unmarshal(decryptedData, &storedEntry); err != nil {
		return nil, fmt.Errorf("failed to parse password entry: %w", err)
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

// ListPasswords returns metadata for all stored passwords
func (es *EncryptedStorage) ListPasswords() ([]entities.PasswordMetadata, error) {
	if !es.initialized {
		return nil, fmt.Errorf("store not initialized")
	}

	var passwords []entities.PasswordMetadata

	entries, err := os.ReadDir(es.storePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read store directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".gpg") {
			continue
		}

		// Extract name from filename
		name := strings.TrimSuffix(entry.Name(), ".gpg")
		name = es.unsanitizeFileName(name)

		// Load just metadata by loading the full entry
		passwordEntry, err := es.LoadPassword(name)
		if err != nil {
			// Skip entries that can't be decrypted
			continue
		}

		// Convert to PasswordMetadata
		metadata := entities.PasswordMetadata{
			Service:   passwordEntry.Service,
			Username:  passwordEntry.Username,
			URL:       passwordEntry.URL,
			Notes:     passwordEntry.Notes,
			CreatedAt: passwordEntry.CreatedAt,
			UpdatedAt: passwordEntry.UpdatedAt,
		}

		// Add auto-rotation info if present
		if passwordEntry.AutoRotation != nil && passwordEntry.AutoRotation.Enabled {
			daysUntilNext := int(time.Until(passwordEntry.AutoRotation.NextRotationAt).Hours() / 24)
			metadata.AutoRotation = &entities.AutoRotationInfo{
				Enabled:       true,
				IntervalDays:  passwordEntry.AutoRotation.IntervalDays,
				NextRotation:  passwordEntry.AutoRotation.NextRotationAt,
				DaysUntilNext: daysUntilNext,
			}
		}

		passwords = append(passwords, metadata)
	}

	return passwords, nil
}

// DeletePassword removes a password entry
func (es *EncryptedStorage) DeletePassword(name string) error {
	if !es.initialized {
		return fmt.Errorf("store not initialized")
	}

	fileName := es.sanitizeFileName(name) + ".gpg"
	filePath := filepath.Join(es.storePath, fileName)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("password entry '%s' not found", name)
	}

	// Remove file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to remove password file: %w", err)
	}

	// Add removal to git and commit
	if err := es.gitService.AddFiles([]string{fileName}); err != nil {
		return fmt.Errorf("failed to add file removal to git: %w", err)
	}

	commitMsg := fmt.Sprintf("Remove password entry: %s", name)
	if err := es.gitService.Commit(commitMsg); err != nil {
		return fmt.Errorf("failed to commit password removal: %w", err)
	}

	return nil
}

// Sync synchronizes with remote repository
func (es *EncryptedStorage) Sync(remote, branch string) error {
	if !es.initialized {
		return fmt.Errorf("store not initialized")
	}

	// Pull changes from remote
	if err := es.gitService.Pull(remote, branch); err != nil {
		return fmt.Errorf("failed to pull from remote: %w", err)
	}

	// Check for conflicts
	conflicts, err := es.gitService.GetConflicts()
	if err != nil {
		return fmt.Errorf("failed to check for conflicts: %w", err)
	}

	if len(conflicts) > 0 {
		return fmt.Errorf("merge conflicts detected in files: %v", conflicts)
	}

	// Push local changes
	hasChanges, err := es.gitService.HasChanges()
	if err != nil {
		return fmt.Errorf("failed to check for changes: %w", err)
	}

	if hasChanges {
		if err := es.gitService.Push(remote, branch); err != nil {
			return fmt.Errorf("failed to push to remote: %w", err)
		}
	}

	return nil
}

// GetStoreInfo returns information about the store
func (es *EncryptedStorage) GetStoreInfo() (*git.RepositoryInfo, error) {
	if !es.initialized {
		return nil, fmt.Errorf("store not initialized")
	}

	return es.gitService.GetStatus()
}

// saveStoreMetadata saves store metadata as encrypted JSON
func (es *EncryptedStorage) saveStoreMetadata(filePath string, metadata entities.PasswordStore) error {
	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize store metadata: %w", err)
	}

	encryptedData, err := es.gpgService.Encrypt(jsonData, "")
	if err != nil {
		return fmt.Errorf("failed to encrypt store metadata: %w", err)
	}

	return os.WriteFile(filePath+".gpg", encryptedData, 0600)
}

// sanitizeFileName converts a password name to a safe filename
func (es *EncryptedStorage) sanitizeFileName(name string) string {
	// Replace unsafe characters with underscores
	unsafe := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " "}
	result := name
	
	for _, char := range unsafe {
		result = strings.ReplaceAll(result, char, "_")
	}
	
	return result
}

// unsanitizeFileName converts a filename back to original name (basic implementation)
func (es *EncryptedStorage) unsanitizeFileName(filename string) string {
	// This is a basic implementation - in real use, you might want to store 
	// the original name in metadata to avoid this conversion
	return strings.ReplaceAll(filename, "_", " ")
}
