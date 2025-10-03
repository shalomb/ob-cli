package editor

import (
	"fmt"
	"os"
	"os/exec"
)

// Service interface for editor operations
type Service interface {
	OpenFile(filePath string) error
}

// RealService handles real editor operations
type RealService struct{}

// MockService handles mock editor operations for testing
type MockService struct {
	OpenedFiles []string
	ConcealLevel int
	Error error
}

// NewService creates a new real editor service
func NewService() Service {
	return &RealService{}
}

// NewMockService creates a new mock editor service
func NewMockService(openedFiles []string, concealLevel int, err error) Service {
	return &MockService{
		OpenedFiles: openedFiles,
		ConcealLevel: concealLevel,
		Error: err,
	}
}

// OpenFile opens a file in the user's preferred editor
func (s *RealService) OpenFile(filePath string) error {
	// Get editor from environment or use fallback
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// Try common editors in order of preference
		editors := []string{"edit", "vim", "nano", "emacs"}
		for _, e := range editors {
			if _, err := exec.LookPath(e); err == nil {
				editor = e
				break
			}
		}
		if editor == "" {
			return fmt.Errorf("no editor found. Please set $EDITOR or install vim/nano/emacs")
		}
	}

	// Run the editor with the file
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// OpenFile mock implementation
func (s *MockService) OpenFile(filePath string) error {
	if s.Error != nil {
		return s.Error
	}
	
	// Add to opened files list
	s.OpenedFiles = append(s.OpenedFiles, filePath)
	return nil
}