package app

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/unop/ob-cli/internal/editor"
	"github.com/unop/ob-cli/internal/frecency"
	"github.com/unop/ob-cli/internal/fzf"
	"github.com/unop/ob-cli/internal/git"
)

func TestApp_New(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "tips mode",
			config: &Config{
				Mode:  "tips",
				Debug: false,
			},
			wantErr: false,
		},
		{
			name: "obsidian mode",
			config: &Config{
				Mode:  "obsidian",
				Debug: true,
			},
			wantErr: false,
		},
		{
			name: "auto mode",
			config: &Config{
				Mode:  "auto",
				Debug: false,
			},
			wantErr: false,
		},
		{
			name: "invalid mode",
			config: &Config{
				Mode:  "invalid",
				Debug: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && app == nil {
				t.Error("New() returned nil app when no error expected")
			}
		})
	}
}

func TestApp_ListFiles(t *testing.T) {
	tests := []struct {
		name        string
		files       []string
		filesError  error
		wantErr     bool
	}{
		{
			name:       "successful list",
			files:      []string{"file1.md", "file2.md", "file3.md"},
			filesError: nil,
			wantErr:    false,
		},
		{
			name:       "empty list",
			files:      []string{},
			filesError: nil,
			wantErr:    false,
		},
		{
			name:       "frecency error",
			files:      nil,
			filesError: errors.New("frecency service error"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock services
			frecencyService := frecency.NewMockService(tt.files, tt.filesError)
			gitService := git.NewMockService(nil, "", 0, 0, nil)
			editorService := editor.NewMockService([]string{}, 0, nil)
			fzfService := fzf.NewMockService("", false, nil)

			// Create app with mock services
			app := &App{
				config:     &Config{Mode: "tips", Debug: false},
				gitService: gitService,
				editor:     editorService,
				fzf:        fzfService,
				frecency:   frecencyService,
				notesDir:   "/tmp/test",
				mode:       "tips",
			}

			err := app.ListFiles()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_ShowGitStatus(t *testing.T) {
	tests := []struct {
		name        string
		statusResult string
		statusError error
		wantErr     bool
	}{
		{
			name:        "successful status",
			statusResult: " M file1.md\nA  file2.md",
			statusError: nil,
			wantErr:    false,
		},
		{
			name:        "empty status",
			statusResult: "",
			statusError: nil,
			wantErr:    false,
		},
		{
			name:        "git error",
			statusResult: "",
			statusError: errors.New("git status failed"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock services
			frecencyService := frecency.NewMockService([]string{}, nil)
			gitService := git.NewMockService(tt.statusError, tt.statusResult, 0, 0, nil)
			editorService := editor.NewMockService([]string{}, 0, nil)
			fzfService := fzf.NewMockService("", false, nil)

			// Create app with mock services
			app := &App{
				config:     &Config{Mode: "tips", Debug: false},
				gitService: gitService,
				editor:     editorService,
				fzf:        fzfService,
				frecency:   frecencyService,
				notesDir:   "/tmp/test",
				mode:       "tips",
			}

			err := app.ShowGitStatus()
			if (err != nil) != tt.wantErr {
				t.Errorf("ShowGitStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_SyncWithRemote(t *testing.T) {
	tests := []struct {
		name     string
		syncError error
		wantErr  bool
	}{
		{
			name:     "successful sync",
			syncError: nil,
			wantErr:  false,
		},
		{
			name:     "sync error",
			syncError: errors.New("git sync failed"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock services
			frecencyService := frecency.NewMockService([]string{}, nil)
			gitService := git.NewMockService(nil, "", 0, 0, tt.syncError)
			editorService := editor.NewMockService([]string{}, 0, nil)
			fzfService := fzf.NewMockService("", false, nil)

			// Create app with mock services
			app := &App{
				config:     &Config{Mode: "tips", Debug: false},
				gitService: gitService,
				editor:     editorService,
				fzf:        fzfService,
				frecency:   frecencyService,
				notesDir:   "/tmp/test",
				mode:       "tips",
			}

			err := app.SyncWithRemote()
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncWithRemote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_RunInteractive_DirectFile(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	
	// Create a test file
	testFile := filepath.Join(tempDir, "test-file.md")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create mock services
	frecencyService := frecency.NewMockService([]string{}, nil)
	gitService := git.NewMockService(nil, "", 0, 0, nil)
	editorService := editor.NewMockService([]string{}, 0, nil)
	fzfService := fzf.NewMockService("", false, nil)

	// Create app with mock services
	app := &App{
		config:     &Config{Mode: "tips", Debug: false},
		gitService: gitService,
		editor:     editorService,
		fzf:        fzfService,
		frecency:   frecencyService,
		notesDir:   tempDir,
		mode:       "tips",
	}

	// Test direct file access
	err = app.RunInteractive("test-file.md")
	if err != nil {
		t.Errorf("RunInteractive() with direct file error = %v", err)
	}

	// Check that editor was called
	mockEditor := editorService.(*editor.MockService)
	if len(mockEditor.OpenedFiles) != 1 {
		t.Errorf("Expected 1 opened file, got %d", len(mockEditor.OpenedFiles))
	}
	if mockEditor.OpenedFiles[0] != "test-file.md" {
		t.Errorf("Expected opened file 'test-file.md', got %s", mockEditor.OpenedFiles[0])
	}
}

func TestApp_RunInteractive_FileCreation(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create mock services
	frecencyService := frecency.NewMockService([]string{}, nil)
	gitService := git.NewMockService(nil, "", 0, 0, nil)
	editorService := editor.NewMockService([]string{}, 0, nil)
	fzfService := fzf.NewMockService("new-file.md", false, nil)

	// Create app with mock services
	app := &App{
		config:     &Config{Mode: "tips", Debug: false},
		gitService: gitService,
		editor:     editorService,
		fzf:        fzfService,
		frecency:   frecencyService,
		notesDir:   tempDir,
		mode:       "tips",
	}

	// Test file creation
	err := app.RunInteractive("")
	if err != nil {
		t.Errorf("RunInteractive() with file creation error = %v", err)
	}

	// Check that file was created
	expectedFile := filepath.Join(tempDir, "new-file.md")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s to be created", expectedFile)
	}

	// Check that editor was called
	mockEditor := editorService.(*editor.MockService)
	if len(mockEditor.OpenedFiles) != 1 {
		t.Errorf("Expected 1 opened file, got %d", len(mockEditor.OpenedFiles))
	}
	if mockEditor.OpenedFiles[0] != "new-file.md" {
		t.Errorf("Expected opened file 'new-file.md', got %s", mockEditor.OpenedFiles[0])
	}
}

func TestApp_RunInteractive_WithSearchTerm(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create mock services
	files := []string{"project-notes.md", "daily-log.md", "project-ideas.md"}
	frecencyService := frecency.NewMockService(files, nil)
	gitService := git.NewMockService(nil, "", 0, 0, nil)
	editorService := editor.NewMockService([]string{}, 0, nil)
	fzfService := fzf.NewMockService("project-notes.md", false, nil)

	// Create app with mock services
	app := &App{
		config:     &Config{Mode: "tips", Debug: false},
		gitService: gitService,
		editor:     editorService,
		fzf:        fzfService,
		frecency:   frecencyService,
		notesDir:   tempDir,
		mode:       "tips",
	}

	// Test search term
	err := app.RunInteractive("project")
	if err != nil {
		t.Errorf("RunInteractive() with search term error = %v", err)
	}

	// Check that editor was called with selected file
	mockEditor := editorService.(*editor.MockService)
	if len(mockEditor.OpenedFiles) != 1 {
		t.Errorf("Expected 1 opened file, got %d", len(mockEditor.OpenedFiles))
	}
	if mockEditor.OpenedFiles[0] != "project-notes.md" {
		t.Errorf("Expected opened file 'project-notes.md', got %s", mockEditor.OpenedFiles[0])
	}
}

func TestApp_RunInteractive_UserCancels(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create mock services
	files := []string{"file1.md", "file2.md", "file3.md"}
	frecencyService := frecency.NewMockService(files, nil)
	gitService := git.NewMockService(nil, "", 0, 0, nil)
	editorService := editor.NewMockService([]string{}, 0, nil)
	fzfService := fzf.NewMockService("", true, nil) // User cancels

	// Create app with mock services
	app := &App{
		config:     &Config{Mode: "tips", Debug: false},
		gitService: gitService,
		editor:     editorService,
		fzf:        fzfService,
		frecency:   frecencyService,
		notesDir:   tempDir,
		mode:       "tips",
	}

	// Test user cancellation
	err := app.RunInteractive("")
	if err != nil {
		t.Errorf("RunInteractive() with user cancellation error = %v", err)
	}

	// Check that no files were opened
	mockEditor := editorService.(*editor.MockService)
	if len(mockEditor.OpenedFiles) != 0 {
		t.Errorf("Expected 0 opened files when user cancels, got %d", len(mockEditor.OpenedFiles))
	}
}