package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGitService_InitializeRepository(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize Git service
	service := NewGitService(tempDir)
	
	// Test initialization
	err = service.InitializeRepository()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Check if .git directory exists
	gitDir := filepath.Join(tempDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Error("Git directory was not created")
	}

	// Check if .gitignore exists
	gitignore := filepath.Join(tempDir, ".gitignore")
	if _, err := os.Stat(gitignore); os.IsNotExist(err) {
		t.Error("Gitignore file was not created")
	}
}

func TestGitService_IsRepository(t *testing.T) {
	// Test with non-repository directory
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	service := NewGitService(tempDir)
	
	// Should not be a repository initially
	if service.IsRepository() {
		t.Error("Expected false for non-repository directory")
	}

	// Initialize repository
	err = service.InitializeRepository()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Should be a repository now
	if !service.IsRepository() {
		t.Error("Expected true for repository directory")
	}
}

func TestGitService_ConfigureUser(t *testing.T) {
	// Create temporary directory and initialize git
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	service := NewGitService(tempDir)
	err = service.InitializeRepository()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Configure user
	err = service.ConfigureUser("Test User", "test@example.com")
	if err != nil {
		t.Errorf("Failed to configure user: %v", err)
	}
}

func TestGitService_AddFilesAndCommit(t *testing.T) {
	// Create temporary directory and initialize git
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	service := NewGitService(tempDir)
	err = service.InitializeRepository()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Configure user for commits
	err = service.ConfigureUser("Test User", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to configure user: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add files
	err = service.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Errorf("Failed to add files: %v", err)
	}

	// Commit
	err = service.Commit("Test commit")
	if err != nil {
		t.Errorf("Failed to commit: %v", err)
	}
}

func TestGitService_GetStatus(t *testing.T) {
	// Create temporary directory and initialize git
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	service := NewGitService(tempDir)
	err = service.InitializeRepository()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Get status
	info, err := service.GetStatus()
	if err != nil {
		t.Errorf("Failed to get status: %v", err)
	}

	if info.Path != tempDir {
		t.Errorf("Expected path %s, got %s", tempDir, info.Path)
	}
}

func TestGitService_HasChanges(t *testing.T) {
	// Create temporary directory and initialize git
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	service := NewGitService(tempDir)
	err = service.InitializeRepository()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Should have changes initially (gitignore was created)
	hasChanges, err := service.HasChanges()
	if err != nil {
		t.Errorf("Failed to check changes: %v", err)
	}

	t.Logf("Has changes: %v", hasChanges)
}
