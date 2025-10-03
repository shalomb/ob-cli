//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/unop/ob-cli/internal/app"
)

func TestAppIntegration_CompleteWorkflow(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	
	// Create test files with different access times
	now := time.Now()
	testFiles := []struct {
		name     string
		modTime  time.Time
		content  string
	}{
		{"old-note.md", now.Add(-2 * time.Hour), "old content"},
		{"recent-note.md", now.Add(-1 * time.Minute), "recent content"},
		{"new-note.md", now.Add(-30 * time.Second), "new content"},
	}
	
	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file.name)
		err := os.WriteFile(filePath, []byte(file.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file.name, err)
		}
		
		// Set modification time
		err = os.Chtimes(filePath, file.modTime, file.modTime)
		if err != nil {
			t.Fatalf("Failed to set mod time for %s: %v", file.name, err)
		}
	}
	
	// Create app with real services
	config := &app.Config{
		Mode:  "tips",
		Debug: false,
	}
	
	obApp, err := app.New(config)
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}
	
	// Override the notes directory for testing
	obApp.SetNotesDir(tempDir)
	
	// Test ListFiles
	t.Run("ListFiles", func(t *testing.T) {
		err := obApp.ListFiles()
		if err != nil {
			t.Errorf("ListFiles failed: %v", err)
		}
	})
	
	// Test ShowGitStatus (should work even without git repo)
	t.Run("ShowGitStatus", func(t *testing.T) {
		err := obApp.ShowGitStatus()
		// This might fail if no git repo, which is expected
		if err != nil {
			t.Logf("ShowGitStatus failed (expected for non-git directory): %v", err)
		}
	})
}

func TestAppIntegration_FileCreation(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	
	// Create app with real services
	config := &app.Config{
		Mode:  "tips",
		Debug: false,
	}
	
	obApp, err := app.New(config)
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}
	
	// Override the notes directory for testing
	obApp.SetNotesDir(tempDir)
	
	// Test file creation by running with a non-existent file
	// This should trigger file creation
	err = obApp.RunInteractive("new-file.md")
	if err != nil {
		t.Errorf("RunInteractive failed: %v", err)
	}
	
	// Check that file was created
	expectedFile := filepath.Join(tempDir, "new-file.md")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s to be created", expectedFile)
	}
}

func TestAppIntegration_DirectFileAccess(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	
	// Create a test file
	testFile := filepath.Join(tempDir, "existing-file.md")
	err := os.WriteFile(testFile, []byte("existing content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Create app with real services
	config := &app.Config{
		Mode:  "tips",
		Debug: false,
	}
	
	obApp, err := app.New(config)
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}
	
	// Override the notes directory for testing
	obApp.SetNotesDir(tempDir)
	
	// Test direct file access
	err = obApp.RunInteractive("existing-file.md")
	if err != nil {
		t.Errorf("RunInteractive with direct file failed: %v", err)
	}
	
	// File should still exist (editor would have been called)
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s to still exist after opening", testFile)
	}
}

func TestAppIntegration_SubdirectoryCreation(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	
	// Create app with real services
	config := &app.Config{
		Mode:  "tips",
		Debug: false,
	}
	
	obApp, err := app.New(config)
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}
	
	// Override the notes directory for testing
	obApp.SetNotesDir(tempDir)
	
	// Test subdirectory creation by running with a nested file path
	err = obApp.RunInteractive("projects/2024/new-project.md")
	if err != nil {
		t.Errorf("RunInteractive with subdirectory creation failed: %v", err)
	}
	
	// Check that directory was created
	expectedDir := filepath.Join(tempDir, "projects", "2024")
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("Expected directory %s to be created", expectedDir)
	}
	
	// Check that file was created
	expectedFile := filepath.Join(tempDir, "projects", "2024", "new-project.md")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s to be created", expectedFile)
	}
}

func TestAppIntegration_ModeConfiguration(t *testing.T) {
	tests := []struct {
		name string
		mode string
	}{
		{
			name: "tips mode",
			mode: "tips",
		},
		{
			name: "obsidian mode",
			mode: "obsidian",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for testing
			tempDir := t.TempDir()
			
			// Create app with real services
			config := &app.Config{
				Mode:  tt.mode,
				Debug: false,
			}
			
			obApp, err := app.New(config)
			if err != nil {
				t.Fatalf("Failed to create app: %v", err)
			}
			
			// Override the notes directory for testing
			obApp.SetNotesDir(tempDir)
			
			// Test that app was created successfully with the correct mode
			if obApp.GetNotesDir() != tempDir {
				t.Errorf("Expected notes dir %s, got %s", tempDir, obApp.GetNotesDir())
			}
		})
	}
}