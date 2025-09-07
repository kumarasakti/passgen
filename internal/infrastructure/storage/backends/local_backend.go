package backends

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalBackend implements StorageBackend for local file system storage
type LocalBackend struct {
	basePath string
}

// LocalConfig holds the configuration for local backend
type LocalConfig struct {
	BasePath string `yaml:"base_path" json:"base_path"`
}

// NewLocalBackend creates a new local storage backend
func NewLocalBackend(basePath string) *LocalBackend {
	return &LocalBackend{
		basePath: basePath,
	}
}

// SaveFile saves data to local file system with the given key
func (l *LocalBackend) SaveFile(key string, data []byte) error {
	fullPath := filepath.Join(l.basePath, key)
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Write file with secure permissions
	if err := os.WriteFile(fullPath, data, 0600); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	
	return nil
}

// LoadFile loads data from local file system with the given key
func (l *LocalBackend) LoadFile(key string) ([]byte, error) {
	fullPath := filepath.Join(l.basePath, key)
	
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load file: %w", err)
	}
	
	return data, nil
}

// ListFiles lists all files with the given prefix
func (l *LocalBackend) ListFiles(prefix string) ([]string, error) {
	var files []string
	
	err := filepath.Walk(l.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			return nil
		}
		
		// Get relative path from base path
		relPath, err := filepath.Rel(l.basePath, path)
		if err != nil {
			return err
		}
		
		// Check if file matches prefix
		if prefix == "" || strings.HasPrefix(relPath, prefix) {
			files = append(files, relPath)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	
	return files, nil
}

// DeleteFile deletes a file from local file system
func (l *LocalBackend) DeleteFile(key string) error {
	fullPath := filepath.Join(l.basePath, key)
	
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	
	return nil
}

// FileExists checks if a file exists in local file system
func (l *LocalBackend) FileExists(key string) (bool, error) {
	fullPath := filepath.Join(l.basePath, key)
	
	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}
	
	return true, nil
}

// Initialize creates the store metadata locally
func (l *LocalBackend) Initialize(storeName string) error {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(l.basePath, 0700); err != nil {
		return fmt.Errorf("failed to create store directory: %w", err)
	}
	
	// Create a metadata file to mark the store as initialized
	metadataKey := ".passgen-store.json"
	metadata := fmt.Sprintf(`{
  "name": "%s",
  "backend": "local",
  "created_at": "%s",
  "version": "1.0"
}`, storeName, time.Now().Format(time.RFC3339))
	
	return l.SaveFile(metadataKey, []byte(metadata))
}

// IsInitialized checks if the store is initialized locally
func (l *LocalBackend) IsInitialized(storeName string) (bool, error) {
	return l.FileExists(".passgen-store.json")
}

// Sync is a no-op for local storage since it's always in sync
func (l *LocalBackend) Sync() error {
	// Local storage is always synchronized, no action needed
	return nil
}

// GetBackendType returns the backend type
func (l *LocalBackend) GetBackendType() string {
	return "local"
}

// GetConnectionInfo returns connection information for debugging
func (l *LocalBackend) GetConnectionInfo() map[string]string {
	return map[string]string{
		"type":      "local",
		"base_path": l.basePath,
	}
}
