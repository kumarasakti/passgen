package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitService handles Git repository operations
type GitService struct {
	repoPath string
}

// NewGitService creates a new Git service instance
func NewGitService(repoPath string) *GitService {
	return &GitService{
		repoPath: repoPath,
	}
}

// RepositoryInfo contains information about a Git repository
type RepositoryInfo struct {
	Path      string
	RemoteURL string
	Branch    string
	LastCommit string
	Status    string
}

// InitializeRepository initializes a new Git repository
func (g *GitService) InitializeRepository() error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(g.repoPath, 0755); err != nil {
		return fmt.Errorf("failed to create repository directory: %w", err)
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = g.repoPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to initialize git repository: %s - %w", string(output), err)
	}

	// Create initial .gitignore
	gitignoreContent := `# Temporary files
*.tmp
*.swp
*.bak

# OS generated files
.DS_Store
Thumbs.db

# Editor files
.vscode/
.idea/
`
	gitignorePath := filepath.Join(g.repoPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	return nil
}

// CloneRepository clones a remote repository
func (g *GitService) CloneRepository(remoteURL string) error {
	// Create parent directory
	parentDir := filepath.Dir(g.repoPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	repoName := filepath.Base(g.repoPath)
	cmd := exec.Command("git", "clone", remoteURL, repoName)
	cmd.Dir = parentDir
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %s - %w", string(output), err)
	}

	return nil
}

// AddRemote adds a remote repository
func (g *GitService) AddRemote(name, url string) error {
	cmd := exec.Command("git", "remote", "add", name, url)
	cmd.Dir = g.repoPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add remote: %s - %w", string(output), err)
	}

	return nil
}

// Pull pulls changes from remote repository
func (g *GitService) Pull(remote, branch string) error {
	if remote == "" {
		remote = "origin"
	}
	if branch == "" {
		branch = "main"
	}

	cmd := exec.Command("git", "pull", remote, branch)
	cmd.Dir = g.repoPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to pull from remote: %s - %w", string(output), err)
	}

	return nil
}

// Push pushes changes to remote repository
func (g *GitService) Push(remote, branch string) error {
	if remote == "" {
		remote = "origin"
	}
	if branch == "" {
		branch = "main"
	}

	cmd := exec.Command("git", "push", remote, branch)
	cmd.Dir = g.repoPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to push to remote: %s - %w", string(output), err)
	}

	return nil
}

// AddFiles adds files to Git staging area
func (g *GitService) AddFiles(files []string) error {
	if len(files) == 0 {
		files = []string{"."}
	}

	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	cmd.Dir = g.repoPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add files: %s - %w", string(output), err)
	}

	return nil
}

// Commit creates a new commit
func (g *GitService) Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = g.repoPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to commit: %s - %w", string(output), err)
	}

	return nil
}

// GetStatus returns the current repository status
func (g *GitService) GetStatus() (*RepositoryInfo, error) {
	info := &RepositoryInfo{
		Path: g.repoPath,
	}

	// Get current branch
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err == nil {
		info.Branch = strings.TrimSpace(string(output))
	}

	// Get remote URL
	cmd = exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = g.repoPath
	output, err = cmd.Output()
	if err == nil {
		info.RemoteURL = strings.TrimSpace(string(output))
	}

	// Get last commit
	cmd = exec.Command("git", "log", "-1", "--format=%H %s")
	cmd.Dir = g.repoPath
	output, err = cmd.Output()
	if err == nil {
		info.LastCommit = strings.TrimSpace(string(output))
	}

	// Get status
	cmd = exec.Command("git", "status", "--porcelain")
	cmd.Dir = g.repoPath
	output, err = cmd.Output()
	if err == nil {
		if len(output) == 0 {
			info.Status = "clean"
		} else {
			info.Status = "modified"
		}
	}

	return info, nil
}

// HasChanges checks if there are uncommitted changes
func (g *GitService) HasChanges() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = g.repoPath
	
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}

	return len(strings.TrimSpace(string(output))) > 0, nil
}

// IsRepository checks if the path is a Git repository
func (g *GitService) IsRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = g.repoPath
	
	err := cmd.Run()
	return err == nil
}

// ConfigureUser sets Git user configuration
func (g *GitService) ConfigureUser(name, email string) error {
	// Set user name
	cmd := exec.Command("git", "config", "user.name", name)
	cmd.Dir = g.repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set git user name: %w", err)
	}

	// Set user email
	cmd = exec.Command("git", "config", "user.email", email)
	cmd.Dir = g.repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set git user email: %w", err)
	}

	return nil
}

// GetConflicts returns files with merge conflicts
func (g *GitService) GetConflicts() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=U")
	cmd.Dir = g.repoPath
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get conflicts: %w", err)
	}

	conflicts := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(conflicts) == 1 && conflicts[0] == "" {
		return []string{}, nil
	}

	return conflicts, nil
}

// ResolveConflict marks a file as resolved
func (g *GitService) ResolveConflict(filePath string) error {
	cmd := exec.Command("git", "add", filePath)
	cmd.Dir = g.repoPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to resolve conflict for %s: %s - %w", filePath, string(output), err)
	}

	return nil
}

// CreateBranch creates and switches to a new branch
func (g *GitService) CreateBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = g.repoPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create branch %s: %s - %w", branchName, string(output), err)
	}

	return nil
}

// SwitchBranch switches to an existing branch
func (g *GitService) SwitchBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", branchName)
	cmd.Dir = g.repoPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to switch to branch %s: %s - %w", branchName, string(output), err)
	}

	return nil
}
