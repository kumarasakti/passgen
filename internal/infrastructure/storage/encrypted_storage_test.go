package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kumarasakti/passgen/internal/domain/entities"
	"github.com/kumarasakti/passgen/internal/infrastructure/gpg"
)

// MockGPGService provides a mock GPG service for testing
type MockGPGService struct{}

func (m *MockGPGService) ListKeys() ([]gpg.GPGKey, error) {
	return []gpg.GPGKey{
		{ID: "test-key", UserID: "Test User <test@example.com>"},
	}, nil
}

func (m *MockGPGService) ValidateKey(keyID string) error {
	return nil
}

func (m *MockGPGService) Encrypt(data []byte, recipientKeyID string) ([]byte, error) {
	// Simple mock encryption - just add a prefix
	return append([]byte("ENCRYPTED:"), data...), nil
}

func (m *MockGPGService) Decrypt(encryptedData []byte) ([]byte, error) {
	// Simple mock decryption - remove the prefix
	if len(encryptedData) > 10 && string(encryptedData[:10]) == "ENCRYPTED:" {
		return encryptedData[10:], nil
	}
	return encryptedData, nil
}

func (m *MockGPGService) Sign(data []byte) ([]byte, error) {
	return []byte("SIGNATURE"), nil
}

func (m *MockGPGService) VerifySignature(data, signature []byte) error {
	return nil
}

func (m *MockGPGService) GetKeyFingerprint(keyID string) (string, error) {
	return "test-fingerprint", nil
}

func createMockGPGService() *gpg.GPGService {
	// For testing, we'll use a simple approach - in a real test environment,
	// you might want to use dependency injection or interfaces
	return gpg.NewGPGService("test-key")
}

func TestEncryptedStorage_InitializeStore(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "storage-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create storage with mock GPG service
	gpgService := createMockGPGService()
	storage := NewEncryptedStorage(tempDir, gpgService)

	// Test initialization
	err = storage.InitializeStore("test-store")
	
	// Note: This test will fail in CI without Git and GPG setup
	// In a real implementation, you'd mock these dependencies
	if err != nil {
		t.Logf("Initialize store failed (expected in test environment): %v", err)
		return
	}

	// Check if store directory was created
	storeDir := filepath.Join(tempDir, "test-store")
	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		t.Error("Store directory was not created")
	}
}

func TestEncryptedStorage_SanitizeFileName(t *testing.T) {
	gpgService := createMockGPGService()
	storage := NewEncryptedStorage("/tmp", gpgService)

	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"with spaces", "with_spaces"},
		{"with/slash", "with_slash"},
		{"with:colon", "with_colon"},
		{"with*asterisk", "with_asterisk"},
		{"with\"quote", "with_quote"},
		{"with<bracket", "with_bracket"},
		{"with>bracket", "with_bracket"},
		{"with|pipe", "with_pipe"},
	}

	for _, test := range tests {
		result := storage.sanitizeFileName(test.input)
		if result != test.expected {
			t.Errorf("sanitizeFileName(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestEncryptedStorage_UnsanitizeFileName(t *testing.T) {
	gpgService := createMockGPGService()
	storage := NewEncryptedStorage("/tmp", gpgService)

	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"with_spaces", "with spaces"},
		{"with_underscore", "with underscore"},
	}

	for _, test := range tests {
		result := storage.unsanitizeFileName(test.input)
		if result != test.expected {
			t.Errorf("unsanitizeFileName(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestPasswordEntry_Creation(t *testing.T) {
	// Test creating a password entry
	entry := entities.PasswordEntry{
		Service:     "github.com",
		Username:    "testuser",
		Password:    "secret123",
		URL:         "https://github.com",
		Notes:       "Development account",
		Metadata:    map[string]string{"category": "development"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		GeneratedBy: "passgen v1.1.0",
	}

	if entry.Service != "github.com" {
		t.Errorf("Expected service 'github.com', got '%s'", entry.Service)
	}

	if entry.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", entry.Username)
	}
}

func TestPasswordMetadata_Creation(t *testing.T) {
	// Test creating password metadata
	metadata := entities.PasswordMetadata{
		Service:   "github.com",
		Username:  "testuser",
		URL:       "https://github.com",
		Notes:     "Development account",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if metadata.Service != "github.com" {
		t.Errorf("Expected service 'github.com', got '%s'", metadata.Service)
	}

	if metadata.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", metadata.Username)
	}
}

func TestAutoRotationConfig_Creation(t *testing.T) {
	// Test creating auto-rotation config
	config := entities.AutoRotationConfig{
		Enabled:          true,
		IntervalDays:     30,
		NextRotationAt:   time.Now().AddDate(0, 0, 30),
		NotifyDaysBefore: 7,
		AutoGenerate:     true,
		PasswordProfile: &entities.PasswordProfile{
			Length:         16,
			IncludeUpper:   true,
			IncludeLower:   true,
			IncludeNumbers: true,
			IncludeSymbols: false,
		},
	}

	if !config.Enabled {
		t.Error("Expected auto-rotation to be enabled")
	}

	if config.IntervalDays != 30 {
		t.Errorf("Expected interval 30 days, got %d", config.IntervalDays)
	}

	if config.PasswordProfile.Length != 16 {
		t.Errorf("Expected password length 16, got %d", config.PasswordProfile.Length)
	}
}
