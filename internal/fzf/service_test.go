package fzf

import (
	"errors"
	"testing"
)

func TestMockService_SelectFile(t *testing.T) {
	// Test with empty file list
	service := NewMockService("", false, nil)
	files := []string{}
	selection, err := service.SelectFile(files, "")
	if err != nil {
		t.Fatalf("SelectFile failed: %v", err)
	}
	if selection != "" {
		t.Errorf("Expected empty selection for empty file list, got %s", selection)
	}
	
	// Test with mock selection
	service = NewMockService("file1.md", false, nil)
	files = []string{"file1.md", "file2.md", "file3.md"}
	selection, err = service.SelectFile(files, "")
	if err != nil {
		t.Fatalf("SelectFile failed: %v", err)
	}
	if selection != "file1.md" {
		t.Errorf("Expected selection 'file1.md', got %s", selection)
	}
}

func TestMockService_SelectFile_WithQuery(t *testing.T) {
	service := NewMockService("project-notes.md", false, nil)
	files := []string{"project-notes.md", "daily-log.md", "project-ideas.md"}
	selection, err := service.SelectFile(files, "project")
	if err != nil {
		t.Fatalf("SelectFile with query failed: %v", err)
	}
	if selection != "project-notes.md" {
		t.Errorf("Expected selection 'project-notes.md', got %s", selection)
	}
}

func TestMockService_SelectFile_UserCancels(t *testing.T) {
	service := NewMockService("", true, nil) // ShouldExit = true
	files := []string{"file1.md", "file2.md", "file3.md"}
	selection, err := service.SelectFile(files, "")
	if err != nil {
		t.Fatalf("SelectFile failed: %v", err)
	}
	if selection != "" {
		t.Errorf("Expected empty selection when user cancels, got %s", selection)
	}
}

func TestMockService_SelectFile_Error(t *testing.T) {
	expectedError := errors.New("fzf not found")
	service := NewMockService("", false, expectedError)
	files := []string{"file1.md", "file2.md", "file3.md"}
	_, err := service.SelectFile(files, "")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}
}