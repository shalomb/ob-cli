package fzf

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Service interface for fzf operations
type Service interface {
	SelectFile(files []string, query string) (string, error)
}

// RealService handles real fzf integration
type RealService struct{}

// MockService handles mock fzf integration for testing
type MockService struct {
	Selection string
	ShouldExit bool
	Error error
}

// NewService creates a new real fzf service
func NewService() Service {
	return &RealService{}
}

// NewMockService creates a new mock fzf service
func NewMockService(selection string, shouldExit bool, err error) Service {
	return &MockService{
		Selection: selection,
		ShouldExit: shouldExit,
		Error: err,
	}
}

// SelectFile runs fzf to select a file from the given list
func (s *RealService) SelectFile(files []string, query string) (string, error) {
	if len(files) == 0 {
		return "", nil
	}
	
	// Check if fzf is available
	if !s.isFzfAvailable() {
		return "", fmt.Errorf("fzf is not installed. Please install fzf: https://github.com/junegunn/fzf")
	}
	
	// Prepare fzf command with print-query to capture user input
	cmd := exec.Command("fzf", "--height", "40%", "--border", "--print-query")
	
	// Add query if provided
	if query != "" {
		cmd.Args = append(cmd.Args, "--query", query)
	}
	
	// Set up stdin to provide file list
	cmd.Stdin = strings.NewReader(strings.Join(files, "\n"))
	cmd.Stderr = os.Stderr
	
	// Capture stdout to get the selection
	output, err := cmd.Output()
	if err != nil {
		// fzf returns exit code 1 when user cancels (ESC)
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return "", nil // User cancelled, return empty selection
		}
		return "", fmt.Errorf("fzf execution failed: %w", err)
	}
	
	// Parse output: first line is query, second line is selection (if any)
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		return "", nil
	}
	
	// If user typed something but didn't select from list, return what they typed
	if len(lines) == 1 {
		return strings.TrimSpace(lines[0]), nil
	}
	
	// If user selected something from the list, return the selection
	if len(lines) >= 2 && strings.TrimSpace(lines[1]) != "" {
		return strings.TrimSpace(lines[1]), nil
	}
	
	// Fallback to query if no selection
	return strings.TrimSpace(lines[0]), nil
}

// SelectFile returns mock selection for testing
func (s *MockService) SelectFile(files []string, query string) (string, error) {
	if s.Error != nil {
		return "", s.Error
	}
	
	if s.ShouldExit {
		return "", nil // User cancelled
	}
	
	return s.Selection, nil
}

// isFzfAvailable checks if fzf is installed and available
func (s *RealService) isFzfAvailable() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}