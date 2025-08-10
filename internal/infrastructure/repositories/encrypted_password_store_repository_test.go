package repositories

import (
	"testing"
	"time"

	"github.com/kumarasakti/passgen/internal/domain/entities"
	"github.com/kumarasakti/passgen/internal/infrastructure/gpg"
	"github.com/kumarasakti/passgen/internal/infrastructure/storage"
)

func TestEncryptedPasswordStoreRepository_Creation(t *testing.T) {
	repo := NewEncryptedPasswordStoreRepository()
	
	if repo == nil {
		t.Error("Expected repository to be created")
	}
	
	if repo.storages == nil {
		t.Error("Expected storages map to be initialized")
	}
}

func TestEncryptedPasswordStoreRepository_RegisterStorage(t *testing.T) {
	repo := NewEncryptedPasswordStoreRepository()
	
	// Create a mock storage
	gpgService := gpg.NewGPGService("test-key")
	mockStorage := storage.NewEncryptedStorage("/tmp/test", gpgService)
	
	// Register storage
	repo.RegisterStorage("test-store", mockStorage)
	
	// Check if storage is registered
	if len(repo.storages) != 1 {
		t.Errorf("Expected 1 storage, got %d", len(repo.storages))
	}
	
	if _, exists := repo.storages["test-store"]; !exists {
		t.Error("Expected test-store to be registered")
	}
}

func TestEncryptedPasswordStoreRepository_NotFoundErrors(t *testing.T) {
	repo := NewEncryptedPasswordStoreRepository()
	
	// Test GetPassword with non-existent store
	_, err := repo.GetPassword("non-existent", "service")
	if err == nil {
		t.Error("Expected error for non-existent store")
	}
	
	// Test GetPasswordMetadata with non-existent store
	_, err = repo.GetPasswordMetadata("non-existent", "service")
	if err == nil {
		t.Error("Expected error for non-existent store")
	}
	
	// Test SavePassword with non-existent store
	entry := entities.PasswordEntry{
		Service:     "test",
		Password:    "secret",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = repo.SavePassword("non-existent", &entry)
	if err == nil {
		t.Error("Expected error for non-existent store")
	}
	
	// Test ListPasswords with non-existent store
	_, err = repo.ListPasswords("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent store")
	}
	
	// Test DeletePassword with non-existent store
	err = repo.DeletePassword("non-existent", "service")
	if err == nil {
		t.Error("Expected error for non-existent store")
	}
	
	// Test Sync with non-existent store
	err = repo.Sync("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent store")
	}
}

func TestEncryptedPasswordStoreRepository_PlaceholderMethods(t *testing.T) {
	repo := NewEncryptedPasswordStoreRepository()
	
	// Test placeholder methods that should return errors
	store := entities.PasswordStore{Name: "test"}
	err := repo.CreateStore(store)
	if err == nil {
		t.Error("Expected CreateStore to return error (not implemented)")
	}
	
	_, err = repo.GetStore("test")
	if err == nil {
		t.Error("Expected GetStore to return error (not implemented)")
	}
	
	_, err = repo.ListStores()
	if err == nil {
		t.Error("Expected ListStores to return error (not implemented)")
	}
	
	err = repo.DeleteStore("test")
	if err == nil {
		t.Error("Expected DeleteStore to return error (not implemented)")
	}
	
	err = repo.SetDefaultStore("test")
	if err == nil {
		t.Error("Expected SetDefaultStore to return error (not implemented)")
	}
	
	err = repo.CopyPasswordToClipboard("test", "service", time.Minute)
	if err == nil {
		t.Error("Expected CopyPasswordToClipboard to return error (not implemented)")
	}
	
	err = repo.ShowPasswordSecure("test", "service", nil)
	if err == nil {
		t.Error("Expected ShowPasswordSecure to return error (not implemented)")
	}
}
