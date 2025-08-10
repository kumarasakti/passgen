package gpg

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// GPGService handles GPG encryption and decryption operations
type GPGService struct {
	keyID string
}

// NewGPGService creates a new GPG service instance
func NewGPGService(keyID string) *GPGService {
	return &GPGService{
		keyID: keyID,
	}
}

// GPGKey represents a GPG key information
type GPGKey struct {
	ID          string
	UserID      string
	Fingerprint string
	KeyType     string
	KeyLength   int
}

// ListKeys returns available GPG keys
func (g *GPGService) ListKeys() ([]GPGKey, error) {
	cmd := exec.Command("gpg", "--list-secret-keys", "--with-colons")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list GPG keys: %w", err)
	}

	return parseGPGKeys(string(output)), nil
}

// ValidateKey checks if the specified key exists and is usable
func (g *GPGService) ValidateKey(keyID string) error {
	cmd := exec.Command("gpg", "--list-secret-keys", keyID)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("GPG key %s not found or not accessible: %w", keyID, err)
	}
	return nil
}

// Encrypt encrypts data using the configured GPG key
func (g *GPGService) Encrypt(data []byte, recipientKeyID string) ([]byte, error) {
	if recipientKeyID == "" {
		recipientKeyID = g.keyID
	}
	
	cmd := exec.Command("gpg", "--armor", "--encrypt", "--recipient", recipientKeyID, "--trust-model", "always")
	cmd.Stdin = bytes.NewReader(data)
	
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("GPG encryption failed: %s - %w", stderr.String(), err)
	}
	
	return out.Bytes(), nil
}

// Decrypt decrypts GPG-encrypted data
func (g *GPGService) Decrypt(encryptedData []byte) ([]byte, error) {
	cmd := exec.Command("gpg", "--quiet", "--batch", "--decrypt")
	cmd.Stdin = bytes.NewReader(encryptedData)
	
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("GPG decryption failed: %s - %w", stderr.String(), err)
	}
	
	return out.Bytes(), nil
}

// Sign creates a detached signature for the data
func (g *GPGService) Sign(data []byte) ([]byte, error) {
	cmd := exec.Command("gpg", "--armor", "--detach-sign", "--local-user", g.keyID)
	cmd.Stdin = bytes.NewReader(data)
	
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("GPG signing failed: %s - %w", stderr.String(), err)
	}
	
	return out.Bytes(), nil
}

// VerifySignature verifies a detached signature
func (g *GPGService) VerifySignature(data, signature []byte) error {
	// Write signature to temporary buffer for verification
	cmd := exec.Command("gpg", "--verify", "-", "-")
	
	// Create combined input: signature then data
	var input bytes.Buffer
	input.Write(signature)
	input.Write(data)
	cmd.Stdin = &input
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("signature verification failed: %s - %w", stderr.String(), err)
	}
	
	return nil
}

// GetKeyFingerprint returns the fingerprint for a key ID
func (g *GPGService) GetKeyFingerprint(keyID string) (string, error) {
	cmd := exec.Command("gpg", "--list-keys", "--with-colons", keyID)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get key fingerprint: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "fpr:") {
			parts := strings.Split(line, ":")
			if len(parts) >= 10 {
				return parts[9], nil
			}
		}
	}
	
	return "", fmt.Errorf("fingerprint not found for key %s", keyID)
}

// parseGPGKeys parses GPG key listing output
func parseGPGKeys(output string) []GPGKey {
	var keys []GPGKey
	lines := strings.Split(output, "\n")
	
	var currentKey *GPGKey
	
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}
		
		switch parts[0] {
		case "sec":
			// Secret key line: sec:u:4096:1:KEYID:CREATED:::u:::scESCA:::+::0:
			if len(parts) >= 5 {
				currentKey = &GPGKey{
					KeyType:   parts[3],
					KeyLength: parseKeyLength(parts[2]),
				}
			}
		case "uid":
			// User ID line: uid:u::::CREATED::KEYID::USERID::::::::::0:
			if currentKey != nil && len(parts) >= 10 && parts[9] != "" {
				currentKey.UserID = parts[9]
				currentKey.ID = extractKeyID(currentKey.UserID)
			}
		case "fpr":
			// Fingerprint line: fpr:::::::::FINGERPRINT:
			if currentKey != nil && len(parts) >= 10 && parts[9] != "" {
				currentKey.Fingerprint = parts[9]
				// Only add key if we have both UserID and Fingerprint
				if currentKey.UserID != "" {
					keys = append(keys, *currentKey)
				}
				currentKey = nil
			}
		}
	}
	
	return keys
}

// parseKeyLength extracts key length from GPG output
func parseKeyLength(keyInfo string) int {
	// Extract numeric part from strings like "4096R" or "2048D"
	var length int
	fmt.Sscanf(keyInfo, "%d", &length)
	return length
}

// extractKeyID extracts a short ID from user ID string
func extractKeyID(userID string) string {
	// Extract email or first part as ID
	if strings.Contains(userID, "<") && strings.Contains(userID, ">") {
		start := strings.Index(userID, "<") + 1
		end := strings.Index(userID, ">")
		if start < end {
			return userID[start:end]
		}
	}
	
	// Fallback to first word
	parts := strings.Fields(userID)
	if len(parts) > 0 {
		return parts[0]
	}
	
	return userID
}
