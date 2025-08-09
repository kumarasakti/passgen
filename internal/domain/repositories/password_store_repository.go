package repositories

import (
	"time"

	"github.com/kumarasakti/passgen/internal/domain/entities"
)

// PasswordStoreRepository defines the interface for password store operations
type PasswordStoreRepository interface {
	// Store management
	CreateStore(store entities.PasswordStore) error
	GetStore(name string) (*entities.PasswordStore, error)
	ListStores() ([]entities.PasswordStore, error)
	DeleteStore(name string) error
	SetDefaultStore(name string) error

	// Password operations - secure access
	AddPassword(storeName string, entry entities.PasswordEntry) error
	GetPasswordMetadata(storeName, service string) (*entities.PasswordMetadata, error)
	GetPassword(storeName, service string) (*entities.PasswordEntry, error)
	ListPasswords(storeName string) ([]entities.PasswordMetadata, error)
	UpdatePassword(storeName string, entry entities.PasswordEntry) error
	DeletePassword(storeName, service string) error

	// Secure password access
	CopyPasswordToClipboard(storeName, service string, ttl time.Duration) error
	ShowPasswordSecure(storeName, service string, confirmation func() bool) error

	// Auto-rotation management
	SetAutoRotation(storeName, service string, config entities.AutoRotationConfig) error
	GetRotationStatus(storeName string) ([]entities.RotationStatus, error)
	RotatePassword(storeName, service string, reason string) error
	CheckDueRotations(storeName string) ([]entities.RotationStatus, error)
	GetRotationHistory(storeName, service string) ([]entities.RotationRecord, error)

	// Sync operations
	SyncStore(storeName string) error
	PullStore(storeName string) error
	PushStore(storeName string) error

	// Audit and logging
	AuditPasswordAccess(storeName, service string, action string) error
}

// StoreConfigRepository defines the interface for store configuration management
type StoreConfigRepository interface {
	LoadConfig() (*entities.StoreConfig, error)
	SaveConfig(config *entities.StoreConfig) error
	GetDefaultStore() (string, error)
	SetDefaultStore(storeName string) error
}
