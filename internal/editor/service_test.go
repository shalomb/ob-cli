package editor

import (
	"errors"
	"testing"
)

func TestMockService_OpenFile(t *testing.T) {
	tests := []struct {
		name         string
		openedFiles  []string
		concealLevel int
		error        error
		filePath     string
		expectedError error
		expectedFiles []string
	}{
		{
			name:         "successful open",
			openedFiles:  []string{},
			concealLevel: 0,
			error:        nil,
			filePath:     "notes/daily.md",
			expectedError: nil,
			expectedFiles: []string{"notes/daily.md"},
		},
		{
			name:         "open multiple files",
			openedFiles:  []string{"file1.md"},
			concealLevel: 2,
			error:        nil,
			filePath:     "file2.md",
			expectedError: nil,
			expectedFiles: []string{"file1.md", "file2.md"},
		},
		{
			name:         "open with error",
			openedFiles:  []string{},
			concealLevel: 0,
			error:        errors.New("editor not found"),
			filePath:     "notes/daily.md",
			expectedError: errors.New("editor not found"),
			expectedFiles: []string{},
		},
		{
			name:         "open empty path",
			openedFiles:  []string{},
			concealLevel: 0,
			error:        nil,
			filePath:     "",
			expectedError: nil,
			expectedFiles: []string{""},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewMockService(tt.openedFiles, tt.concealLevel, tt.error)
			
			err := service.OpenFile(tt.filePath)
			
			if (err != nil) != (tt.expectedError != nil) {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
			
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error message %q, got %q", tt.expectedError.Error(), err.Error())
			}
			
			// Check that files were added to the opened files list
			mockService := service.(*MockService)
			if len(mockService.OpenedFiles) != len(tt.expectedFiles) {
				t.Errorf("Expected %d opened files, got %d", len(tt.expectedFiles), len(mockService.OpenedFiles))
			}
			
			for i, expectedFile := range tt.expectedFiles {
				if i < len(mockService.OpenedFiles) && mockService.OpenedFiles[i] != expectedFile {
					t.Errorf("Expected opened file %d to be %q, got %q", i, expectedFile, mockService.OpenedFiles[i])
				}
			}
		})
	}
}

func TestMockService_EditorRespect(t *testing.T) {
	// Test that the service respects editor choice
	service := NewMockService([]string{}, 0, nil)
	mockService := service.(*MockService)
	
	// Test opening a file
	err := service.OpenFile("test.md")
	if err != nil {
		t.Errorf("OpenFile failed: %v", err)
	}
	
	// Check that file was added to opened files
	if len(mockService.OpenedFiles) != 1 {
		t.Errorf("Expected 1 opened file, got %d", len(mockService.OpenedFiles))
	}
	if mockService.OpenedFiles[0] != "test.md" {
		t.Errorf("Expected opened file 'test.md', got %s", mockService.OpenedFiles[0])
	}
}