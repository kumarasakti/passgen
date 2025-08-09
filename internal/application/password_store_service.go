package application

import (
	"fmt"

	"github.com/kumarasakti/passgen/internal/domain/entities"
	"github.com/kumarasakti/passgen/internal/domain/repositories"
)

// PasswordStoreService handles password store business logic
type PasswordStoreService struct {
	storeRepo  repositories.PasswordStoreRepository
	configRepo repositories.StoreConfigRepository
}

// NewPasswordStoreService creates a new password store service
func NewPasswordStoreService(
	storeRepo repositories.PasswordStoreRepository,
	configRepo repositories.StoreConfigRepository,
) *PasswordStoreService {
	return &PasswordStoreService{
		storeRepo:  storeRepo,
		configRepo: configRepo,
	}
}

// InitializeStore initializes a new password store
func (s *PasswordStoreService) InitializeStore(name, gitURL string) error {
	// Will be implemented in Phase 1B
	return fmt.Errorf("store initialization not implemented yet - coming in Phase 1B")
}

// AddPassword adds a new password to the store
func (s *PasswordStoreService) AddPassword(storeName, service string, req AddPasswordRequest) error {
	// Will be implemented in Phase 1B
	return fmt.Errorf("add password not implemented yet - coming in Phase 1B")
}

// GetPasswordMetadata retrieves password metadata (no actual password)
func (s *PasswordStoreService) GetPasswordMetadata(storeName, service string) (*entities.PasswordMetadata, error) {
	// Will be implemented in Phase 1B with real repository calls
	return nil, fmt.Errorf("get password metadata not implemented yet - coming in Phase 1B")
}

// ListPasswords lists all passwords in a store
func (s *PasswordStoreService) ListPasswords(storeName string) ([]entities.PasswordMetadata, error) {
	// Will be implemented in Phase 1B
	return nil, fmt.Errorf("list passwords not implemented yet - coming in Phase 1B")
}

// SetupAutoRotation configures auto-rotation for a password
func (s *PasswordStoreService) SetupAutoRotation(storeName, service string, config entities.AutoRotationConfig) error {
	// Will be implemented in Phase 1C
	return fmt.Errorf("auto-rotation setup not implemented yet - coming in Phase 1C")
}

// CheckDueRotations checks for passwords that need rotation
func (s *PasswordStoreService) CheckDueRotations(storeName string) ([]entities.RotationStatus, error) {
	// Will be implemented in Phase 1C
	return nil, fmt.Errorf("rotation checking not implemented yet - coming in Phase 1C")
}

// Request/Response types

// AddPasswordRequest represents a request to add a password
type AddPasswordRequest struct {
	Username         string
	URL              string
	Notes            string
	Password         string // If empty, will be generated
	AutoRotate       bool
	RotationInterval int // Days
	NotifyBefore     int // Days before rotation
	PasswordLength   int
}

// StoreInitRequest represents a request to initialize a store
type StoreInitRequest struct {
	Name      string
	GitURL    string
	IsDefault bool
}

// RotationCheckResult represents the result of checking due rotations
type RotationCheckResult struct {
	Urgent   []entities.RotationStatus // Overdue passwords
	Soon     []entities.RotationStatus // Due soon
	Upcoming []entities.RotationStatus // Future rotations
}
