package repositories

import (
	"fmt"
	"time"

	"github.com/kumarasakti/passgen/internal/domain/entities"
	"github.com/kumarasakti/passgen/internal/domain/repositories"
	"github.com/kumarasakti/passgen/internal/infrastructure/storage"
)

// EncryptedPasswordStoreRepository implements the PasswordStoreRepository using encrypted storage
type EncryptedPasswordStoreRepository struct {
	storages map[string]*storage.EncryptedStorage
}

// NewEncryptedPasswordStoreRepository creates a new encrypted password store repository
func NewEncryptedPasswordStoreRepository() *EncryptedPasswordStoreRepository {
	return &EncryptedPasswordStoreRepository{
		storages: make(map[string]*storage.EncryptedStorage),
	}
}

// RegisterStorage registers an encrypted storage for a store
func (r *EncryptedPasswordStoreRepository) RegisterStorage(storeName string, encStorage *storage.EncryptedStorage) {
	r.storages[storeName] = encStorage
}

// GetPassword retrieves a password entry from the specified store
func (r *EncryptedPasswordStoreRepository) GetPassword(storeName, service string) (*entities.PasswordEntry, error) {
	storage, exists := r.storages[storeName]
	if !exists {
		return nil, fmt.Errorf("store '%s' not found", storeName)
	}

	return storage.LoadPassword(service)
}

// GetPasswordMetadata retrieves password metadata (no actual password)
func (r *EncryptedPasswordStoreRepository) GetPasswordMetadata(storeName, service string) (*entities.PasswordMetadata, error) {
	storage, exists := r.storages[storeName]
	if !exists {
		return nil, fmt.Errorf("store '%s' not found", storeName)
	}

	// Load the full entry and extract metadata
	entry, err := storage.LoadPassword(service)
	if err != nil {
		return nil, err
	}

	metadata := &entities.PasswordMetadata{
		Service:   entry.Service,
		Username:  entry.Username,
		URL:       entry.URL,
		Notes:     entry.Notes,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
	}

	// Add auto-rotation info if present
	if entry.AutoRotation != nil && entry.AutoRotation.Enabled {
		daysUntilNext := int(entry.AutoRotation.NextRotationAt.Sub(entry.CreatedAt).Hours() / 24)
		metadata.AutoRotation = &entities.AutoRotationInfo{
			Enabled:       true,
			IntervalDays:  entry.AutoRotation.IntervalDays,
			NextRotation:  entry.AutoRotation.NextRotationAt,
			DaysUntilNext: daysUntilNext,
		}
	}

	return metadata, nil
}

// SavePassword saves a password entry to the specified store
func (r *EncryptedPasswordStoreRepository) SavePassword(storeName string, entry *entities.PasswordEntry) error {
	storage, exists := r.storages[storeName]
	if !exists {
		return fmt.Errorf("store '%s' not found", storeName)
	}

	return storage.SavePassword(*entry)
}

// ListPasswords returns all password metadata from the specified store
func (r *EncryptedPasswordStoreRepository) ListPasswords(storeName string) ([]entities.PasswordMetadata, error) {
	storage, exists := r.storages[storeName]
	if !exists {
		return nil, fmt.Errorf("store '%s' not found", storeName)
	}

	return storage.ListPasswords()
}

// DeletePassword removes a password entry from the specified store
func (r *EncryptedPasswordStoreRepository) DeletePassword(storeName, service string) error {
	storage, exists := r.storages[storeName]
	if !exists {
		return fmt.Errorf("store '%s' not found", storeName)
	}

	return storage.DeletePassword(service)
}

// AddPassword adds a password entry to the specified store
func (r *EncryptedPasswordStoreRepository) AddPassword(storeName string, entry entities.PasswordEntry) error {
	return r.SavePassword(storeName, &entry)
}

// UpdatePassword updates an existing password entry
func (r *EncryptedPasswordStoreRepository) UpdatePassword(storeName string, entry entities.PasswordEntry) error {
	return r.SavePassword(storeName, &entry)
}

// CreateStore creates a new password store (placeholder - needs store config management)
func (r *EncryptedPasswordStoreRepository) CreateStore(store entities.PasswordStore) error {
	// This would need to integrate with the configuration system
	// For now, return an error indicating this needs to be implemented
	return fmt.Errorf("CreateStore not implemented - use InitializeStore instead")
}

// GetStore retrieves store information (placeholder)
func (r *EncryptedPasswordStoreRepository) GetStore(name string) (*entities.PasswordStore, error) {
	return nil, fmt.Errorf("GetStore not implemented - use GetStoreInfo instead")
}

// ListStores lists all available stores (placeholder)
func (r *EncryptedPasswordStoreRepository) ListStores() ([]entities.PasswordStore, error) {
	return nil, fmt.Errorf("ListStores not implemented")
}

// DeleteStore removes a store (placeholder)
func (r *EncryptedPasswordStoreRepository) DeleteStore(name string) error {
	return fmt.Errorf("DeleteStore not implemented")
}

// SetDefaultStore sets the default store (placeholder)
func (r *EncryptedPasswordStoreRepository) SetDefaultStore(name string) error {
	return fmt.Errorf("SetDefaultStore not implemented")
}

// CopyPasswordToClipboard copies password to clipboard (placeholder)
func (r *EncryptedPasswordStoreRepository) CopyPasswordToClipboard(storeName, service string, ttl time.Duration) error {
	return fmt.Errorf("CopyPasswordToClipboard not implemented")
}

// ShowPasswordSecure securely shows password (placeholder)
func (r *EncryptedPasswordStoreRepository) ShowPasswordSecure(storeName, service string, confirmation func() bool) error {
	return fmt.Errorf("ShowPasswordSecure not implemented")
}

// SetAutoRotation sets auto-rotation configuration (placeholder)
func (r *EncryptedPasswordStoreRepository) SetAutoRotation(storeName, service string, config entities.AutoRotationConfig) error {
	return fmt.Errorf("SetAutoRotation not implemented")
}

// GetRotationStatus returns rotation status (placeholder)
func (r *EncryptedPasswordStoreRepository) GetRotationStatus(storeName string) ([]entities.RotationStatus, error) {
	return nil, fmt.Errorf("GetRotationStatus not implemented")
}

// RotatePassword rotates a password with reason (placeholder)
func (r *EncryptedPasswordStoreRepository) RotatePassword(storeName, service string, reason string) error {
	return fmt.Errorf("RotatePassword not implemented")
}

// CheckDueRotations checks for due rotations (placeholder)
func (r *EncryptedPasswordStoreRepository) CheckDueRotations(storeName string) ([]entities.RotationStatus, error) {
	return nil, fmt.Errorf("CheckDueRotations not implemented")
}

// SyncStore synchronizes store (same as Sync)
func (r *EncryptedPasswordStoreRepository) SyncStore(storeName string) error {
	return r.Sync(storeName)
}

// PullStore pulls from remote (placeholder)
func (r *EncryptedPasswordStoreRepository) PullStore(storeName string) error {
	storage, exists := r.storages[storeName]
	if !exists {
		return fmt.Errorf("store '%s' not found", storeName)
	}
	return storage.Sync("origin", "main") // For now, same as sync
}

// PushStore pushes to remote (placeholder)
func (r *EncryptedPasswordStoreRepository) PushStore(storeName string) error {
	storage, exists := r.storages[storeName]
	if !exists {
		return fmt.Errorf("store '%s' not found", storeName)
	}
	return storage.Sync("origin", "main") // For now, same as sync
}

// AuditPasswordAccess logs password access (placeholder)
func (r *EncryptedPasswordStoreRepository) AuditPasswordAccess(storeName, service string, action string) error {
	// For now, this is a no-op. In a real implementation, this would log to a secure audit log
	return nil
}

// UpdateAutoRotationConfig updates auto-rotation configuration (placeholder)
func (r *EncryptedPasswordStoreRepository) UpdateAutoRotationConfig(storeName, service string, config entities.AutoRotationConfig) error {
	return fmt.Errorf("UpdateAutoRotationConfig not implemented")
}

// GetPasswordsNeedingRotation returns passwords that need rotation (placeholder)
func (r *EncryptedPasswordStoreRepository) GetPasswordsNeedingRotation(storeName string) ([]entities.PasswordMetadata, error) {
	return nil, fmt.Errorf("GetPasswordsNeedingRotation not implemented")
}

// GetRotationHistory returns rotation history (placeholder)
func (r *EncryptedPasswordStoreRepository) GetRotationHistory(storeName, service string) ([]entities.RotationRecord, error) {
	return nil, fmt.Errorf("GetRotationHistory not implemented")
}

// Sync synchronizes the store with its remote repository
func (r *EncryptedPasswordStoreRepository) Sync(storeName string) error {
	storage, exists := r.storages[storeName]
	if !exists {
		return fmt.Errorf("store '%s' not found", storeName)
	}

	return storage.Sync("origin", "main")
}

// InitializeStore creates a new password store
func (r *EncryptedPasswordStoreRepository) InitializeStore(storeName string, encStorage *storage.EncryptedStorage) error {
	if err := encStorage.InitializeStore(storeName); err != nil {
		return err
	}

	r.RegisterStorage(storeName, encStorage)
	return nil
}

// ConnectRemote connects a store to a remote Git repository
func (r *EncryptedPasswordStoreRepository) ConnectRemote(storeName, remoteName, remoteURL string) error {
	storage, exists := r.storages[storeName]
	if !exists {
		return fmt.Errorf("store '%s' not found", storeName)
	}

	return storage.ConnectRemote(remoteName, remoteURL)
}

// GetStoreInfo returns information about a store
func (r *EncryptedPasswordStoreRepository) GetStoreInfo(storeName string) (map[string]interface{}, error) {
	storage, exists := r.storages[storeName]
	if !exists {
		return nil, fmt.Errorf("store '%s' not found", storeName)
	}

	gitInfo, err := storage.GetStoreInfo()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":        storeName,
		"path":        gitInfo.Path,
		"remote_url":  gitInfo.RemoteURL,
		"branch":      gitInfo.Branch,
		"last_commit": gitInfo.LastCommit,
		"status":      gitInfo.Status,
	}, nil
}

// Ensure EncryptedPasswordStoreRepository implements PasswordStoreRepository
var _ repositories.PasswordStoreRepository = (*EncryptedPasswordStoreRepository)(nil)
