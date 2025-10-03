package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shalomb/ob-cli/internal/git"
	"github.com/shalomb/ob-cli/internal/editor"
	"github.com/shalomb/ob-cli/internal/fzf"
	"github.com/shalomb/ob-cli/internal/frecency"
	"github.com/shalomb/ob-cli/internal/vault"
)

// Config holds application configuration
type Config struct {
	Mode  string // tips, obsidian, or auto
	Debug bool
}

// App represents the main application
type App struct {
	config     *Config
	gitService git.Service
	editor     editor.Service
	fzf        fzf.Service
	frecency   frecency.Service
	notesDir   string
	mode       string
}

// New creates a new App instance
func New(config *Config) (*App, error) {
	// Determine mode and notes directory
	mode, notesDir, err := determineModeAndDir(config.Mode)
	if err != nil {
		return nil, err
	}

	// Create services
	gitService := git.NewService(notesDir)
	editorService := editor.NewService() // No mode-specific config
	fzfService := fzf.NewService()
	frecencyService := frecency.NewService(notesDir)

	return &App{
		config:     config,
		gitService: gitService,
		editor:     editorService,
		fzf:        fzfService,
		frecency:   frecencyService,
		notesDir:   notesDir,
		mode:       mode,
	}, nil
}

// RunInteractive runs the interactive file selection mode
func (a *App) RunInteractive(target string) error {
	// If target is provided, check if it's a direct file path
	if target != "" {
		if a.isDirectFile(target) {
			return a.handleFileSelection(target)
		}
		// Otherwise, use as search term
	}

	// Start async git fetch
	go a.gitService.FetchAsync()

	// Get frecency-sorted file list
	files, err := a.frecency.GetSortedFiles()
	if err != nil {
		return fmt.Errorf("failed to get file list: %w", err)
	}

	// Run fzf selection
	selection, err := a.fzf.SelectFile(files, target)
	if err != nil {
		return fmt.Errorf("fzf selection failed: %w", err)
	}

	// Check git status if fetch completed
	if err := a.checkGitStatus(); err != nil && a.config.Debug {
		fmt.Fprintf(os.Stderr, "Git status check failed: %v\n", err)
	}

	// Handle file creation or editing
	return a.handleFileSelection(selection)
}

// ListFiles lists all files in a column format
func (a *App) ListFiles() error {
	files, err := a.frecency.GetSortedFiles()
	if err != nil {
		return fmt.Errorf("failed to get file list: %w", err)
	}

	// Print files in column format
	for _, file := range files {
		fmt.Println(file)
	}

	return nil
}

// ShowGitStatus shows git status
func (a *App) ShowGitStatus() error {
	status, err := a.gitService.GetStatus()
	if err != nil {
		return fmt.Errorf("failed to get git status: %w", err)
	}

	fmt.Print(status)
	return nil
}

// SyncWithRemote syncs with remote repository
func (a *App) SyncWithRemote() error {
	return a.gitService.SyncWithRemote()
}

// Private methods

func (a *App) isDirectFile(target string) bool {
	// Check if it looks like a file path (contains .md extension or has directory separators)
	return strings.HasSuffix(target, ".md") || strings.Contains(target, "/")
}

func (a *App) checkGitStatus() error {
	behind, ahead, err := a.gitService.GetSyncStatus()
	if err != nil {
		return err
	}

	if behind > 0 {
		fmt.Printf("⚠️  Repository is %d commits behind origin\n", behind)
		fmt.Println("   Run 'ob-cli --sync' to update")
		fmt.Println()
	}

	if ahead > 0 {
		fmt.Printf("📤 Repository is %d commits ahead of origin\n", ahead)
		fmt.Println()
	}

	return nil
}

func (a *App) handleFileSelection(selection string) error {
	if selection == "" {
		return nil // User cancelled
	}

	// Check if file exists
	fullPath := filepath.Join(a.notesDir, selection)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// Create file and parent directories
		if err := a.createFileWithDirs(fullPath); err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		fmt.Printf("Creating new file: %s\n", selection)
	}

	// Open file in editor
	return a.editor.OpenFile(selection)
}

func (a *App) createFileWithDirs(fullPath string) error {
	// Create parent directories
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create empty file
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	return file.Close()
}

func determineModeAndDir(mode string) (string, string, error) {
	discoverer := vault.NewDiscoverer()
	
	switch mode {
	case "tips":
		vaultPath, err := discoverer.DiscoverTipsVault()
		if err != nil {
			return "", "", fmt.Errorf("failed to discover Tips vault: %w", err)
		}
		return "tips", vaultPath, nil
	case "obsidian":
		vaultPath, err := discoverer.DiscoverObsidianVault()
		if err != nil {
			return "", "", fmt.Errorf("failed to discover Obsidian vault: %w", err)
		}
		return "obsidian", vaultPath, nil
	case "auto":
		// Auto-detect based on command name or environment
		if strings.Contains(os.Args[0], "tips") {
			vaultPath, err := discoverer.DiscoverTipsVault()
			if err != nil {
				return "", "", fmt.Errorf("failed to discover Tips vault: %w", err)
			}
			return "tips", vaultPath, nil
		}
		vaultPath, err := discoverer.DiscoverObsidianVault()
		if err != nil {
			return "", "", fmt.Errorf("failed to discover Obsidian vault: %w", err)
		}
		return "obsidian", vaultPath, nil
	default:
		return "", "", fmt.Errorf("invalid mode: %s", mode)
	}
}

// SetNotesDir sets the notes directory (for testing)
func (a *App) SetNotesDir(dir string) {
	a.notesDir = dir
}

// GetNotesDir returns the notes directory (for testing)
func (a *App) GetNotesDir() string {
	return a.notesDir
}