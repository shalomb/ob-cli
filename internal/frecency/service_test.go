package frecency

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestService_GetSortedFiles(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	
	// Create test files with different access times
	now := time.Now()
	files := []struct {
		name     string
		modTime  time.Time
		content  string
	}{
		{"old-file.md", now.Add(-2 * time.Hour), "old content"},
		{"recent-file.md", now.Add(-1 * time.Minute), "recent content"},
		{"new-file.md", now.Add(-30 * time.Second), "new content"},
	}
	
	for _, file := range files {
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
	
	// Create service
	service := NewService(tempDir)
	
	// Test GetSortedFiles
	sortedFiles, err := service.GetSortedFiles()
	if err != nil {
		t.Fatalf("GetSortedFiles failed: %v", err)
	}
	
	// Verify files are sorted by modification time (most recent first)
	expectedOrder := []string{"new-file.md", "recent-file.md", "old-file.md"}
	if len(sortedFiles) != len(expectedOrder) {
		t.Fatalf("Expected %d files, got %d", len(expectedOrder), len(sortedFiles))
	}
	
	for i, expectedFile := range expectedOrder {
		if sortedFiles[i] != expectedFile {
			t.Errorf("Expected file %d to be %s, got %s", i, expectedFile, sortedFiles[i])
		}
	}
}

func TestService_GetSortedFiles_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()
	
	service := NewService(tempDir)
	sortedFiles, err := service.GetSortedFiles()
	if err != nil {
		t.Fatalf("GetSortedFiles failed: %v", err)
	}
	
	if len(sortedFiles) != 0 {
		t.Errorf("Expected empty directory to return 0 files, got %d", len(sortedFiles))
	}
}

func TestService_GetSortedFiles_NonExistentDirectory(t *testing.T) {
	service := NewService("/non/existent/directory")
	_, err := service.GetSortedFiles()
	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}
}

func TestService_GetSortedFiles_WithSubdirectories(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	
	// Create subdirectories and files
	subDir1 := tempDir + "/notes"
	subDir2 := tempDir + "/projects"
	os.MkdirAll(subDir1, 0755)
	os.MkdirAll(subDir2, 0755)
	
	// Create test files with different access times
	now := time.Now()
	files := []struct {
		name     string
		modTime  time.Time
		content  string
	}{
		{"notes/old-note.md", now.Add(-2 * time.Hour), "old note"},
		{"projects/recent-project.md", now.Add(-1 * time.Minute), "recent project"},
		{"notes/new-note.md", now.Add(-30 * time.Second), "new note"},
		{"projects/old-project.md", now.Add(-1 * time.Hour), "old project"},
	}
	
	for _, file := range files {
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
	
	// Create service
	service := NewService(tempDir)
	
	// Test GetSortedFiles
	sortedFiles, err := service.GetSortedFiles()
	if err != nil {
		t.Fatalf("GetSortedFiles failed: %v", err)
	}
	
	// Verify files are sorted by modification time (most recent first)
	expectedOrder := []string{"notes/new-note.md", "projects/recent-project.md", "projects/old-project.md", "notes/old-note.md"}
	if len(sortedFiles) != len(expectedOrder) {
		t.Fatalf("Expected %d files, got %d", len(expectedOrder), len(sortedFiles))
	}
	
	for i, expectedFile := range expectedOrder {
		if sortedFiles[i] != expectedFile {
			t.Errorf("Expected file %d to be %s, got %s", i, expectedFile, sortedFiles[i])
		}
	}
}

func TestService_GetSortedFiles_IgnoresNonMarkdownFiles(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	
	// Create test files with different extensions
	now := time.Now()
	files := []struct {
		name     string
		modTime  time.Time
		content  string
	}{
		{"file.md", now.Add(-1 * time.Minute), "markdown file"},
		{"file.txt", now.Add(-30 * time.Second), "text file"},
		{"file.json", now.Add(-2 * time.Minute), "json file"},
		{"another.md", now.Add(-45 * time.Second), "another markdown"},
	}
	
	for _, file := range files {
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
	
	// Create service
	service := NewService(tempDir)
	
	// Test GetSortedFiles
	sortedFiles, err := service.GetSortedFiles()
	if err != nil {
		t.Fatalf("GetSortedFiles failed: %v", err)
	}
	
	// Should only return .md files, sorted by modification time
	expectedOrder := []string{"another.md", "file.md"}
	if len(sortedFiles) != len(expectedOrder) {
		t.Fatalf("Expected %d markdown files, got %d", len(expectedOrder), len(sortedFiles))
	}
	
	for i, expectedFile := range expectedOrder {
		if sortedFiles[i] != expectedFile {
			t.Errorf("Expected file %d to be %s, got %s", i, expectedFile, sortedFiles[i])
		}
	}
}

func TestService_GetSortedFiles_IgnoresHiddenFiles(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	
	// Create test files including hidden ones
	now := time.Now()
	files := []struct {
		name     string
		modTime  time.Time
		content  string
	}{
		{"file.md", now.Add(-1 * time.Minute), "visible markdown"},
		{".hidden.md", now.Add(-30 * time.Second), "hidden markdown"},
		{"another.md", now.Add(-2 * time.Minute), "another visible"},
	}
	
	for _, file := range files {
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
	
	// Create service
	service := NewService(tempDir)
	
	// Test GetSortedFiles
	sortedFiles, err := service.GetSortedFiles()
	if err != nil {
		t.Fatalf("GetSortedFiles failed: %v", err)
	}
	
	// Should only return visible .md files, sorted by modification time
	expectedOrder := []string{"file.md", "another.md"}
	if len(sortedFiles) != len(expectedOrder) {
		t.Fatalf("Expected %d visible markdown files, got %d", len(expectedOrder), len(sortedFiles))
	}
	
	for i, expectedFile := range expectedOrder {
		if sortedFiles[i] != expectedFile {
			t.Errorf("Expected file %d to be %s, got %s", i, expectedFile, sortedFiles[i])
		}
	}
}