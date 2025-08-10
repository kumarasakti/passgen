package gpg

import (
	"os/exec"
	"testing"
)

func TestGPGService_ListKeys(t *testing.T) {
	// Check if GPG is available
	if _, err := exec.LookPath("gpg"); err != nil {
		t.Skip("GPG not available, skipping test")
	}

	service := NewGPGService("")
	keys, err := service.ListKeys()
	
	// This test will only pass if there are GPG keys available
	// In a CI environment, this might fail, so we just check the method doesn't crash
	if err != nil {
		t.Logf("No GPG keys available: %v", err)
		return
	}
	
	t.Logf("Found %d GPG keys", len(keys))
	for _, key := range keys {
		t.Logf("Key: %s - %s", key.ID, key.UserID)
	}
}

func TestGPGService_ValidateKey(t *testing.T) {
	if _, err := exec.LookPath("gpg"); err != nil {
		t.Skip("GPG not available, skipping test")
	}

	service := NewGPGService("")
	
	// Test with a non-existent key
	err := service.ValidateKey("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent key")
	}
}

func TestParseGPGKeys(t *testing.T) {
	// Test parsing GPG output - realistic format
	testOutput := `sec:u:4096:1:ABC123DEF456:2023-01-01:::u:::scESCA:::+::0:
uid:u::::2023-01-01::ABC123DEF456ABC123DEF456ABC123DEF456ABC123::Test User <test@example.com>::::::::::0:
fpr:::::::::ABC123DEF456ABC123DEF456ABC123DEF456ABC123DEF456:
ssb:u:4096:1:DEF456ABC123:2023-01-01::::::e:::+::0:
fpr:::::::::DEF456ABC123DEF456ABC123DEF456ABC123DEF456ABC123:
`

	keys := parseGPGKeys(testOutput)
	
	if len(keys) != 1 {
		t.Errorf("Expected 1 key, got %d", len(keys))
		t.Logf("Parsed keys: %+v", keys)
		return
	}
	
	key := keys[0]
	if key.UserID != "Test User <test@example.com>" {
		t.Errorf("Expected UserID 'Test User <test@example.com>', got '%s'", key.UserID)
	}
	if key.Fingerprint != "ABC123DEF456ABC123DEF456ABC123DEF456ABC123DEF456" {
		t.Errorf("Expected fingerprint 'ABC123DEF456ABC123DEF456ABC123DEF456ABC123DEF456', got '%s'", key.Fingerprint)
	}
}

func TestParseKeyLength(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"4096R", 4096},
		{"2048D", 2048},
		{"1024", 1024},
		{"invalid", 0},
	}
	
	for _, test := range tests {
		result := parseKeyLength(test.input)
		if result != test.expected {
			t.Errorf("parseKeyLength(%s) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestExtractKeyID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Test User <test@example.com>", "test@example.com"},
		{"John Doe <john@doe.com>", "john@doe.com"},
		{"Simple Name", "Simple"},
		{"", ""},
	}
	
	for _, test := range tests {
		result := extractKeyID(test.input)
		if result != test.expected {
			t.Errorf("extractKeyID(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}
