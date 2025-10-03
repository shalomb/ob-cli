package frecency

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Service interface for file sorting operations
type Service interface {
	GetSortedFiles() ([]string, error)
}

// RealService handles real file sorting by access time
type RealService struct {
	notesDir string
}

// MockService handles mock file sorting for testing
type MockService struct {
	Files []string
	Error error
}

// NewService creates a new real frecency service
func NewService(notesDir string) Service {
	return &RealService{notesDir: notesDir}
}

// NewMockService creates a new mock frecency service
func NewMockService(files []string, err error) Service {
	return &MockService{
		Files: files,
		Error: err,
	}
}

// FileInfo represents a file with its modification time
type FileInfo struct {
	Name    string
	ModTime time.Time
}

// GetSortedFiles returns files sorted by modification time (most recent first)
func (s *RealService) GetSortedFiles() ([]string, error) {
	var files []FileInfo
	
	err := s.walkDirectory(s.notesDir, &files)
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", s.notesDir, err)
	}
	
	// Sort by modification time (most recent first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.After(files[j].ModTime)
	})
	
	// Extract just the file names
	result := make([]string, len(files))
	for i, file := range files {
		result[i] = file.Name
	}
	
	return result, nil
}

// walkDirectory recursively walks a directory and collects markdown files
func (s *RealService) walkDirectory(dirPath string, files *[]FileInfo) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		// Skip hidden files and directories
		if entry.Name()[0] == '.' {
			continue
		}
		
		fullPath := filepath.Join(dirPath, entry.Name())
		
		if entry.IsDir() {
			// Recursively walk subdirectories
			if err := s.walkDirectory(fullPath, files); err != nil {
				return err
			}
		} else {
			// Check if it's a markdown file
			if filepath.Ext(entry.Name()) == ".md" {
				// Get file info
				info, err := entry.Info()
				if err != nil {
					continue // Skip files we can't get info for
				}
				
				// Get relative path from notes directory
				relPath, err := filepath.Rel(s.notesDir, fullPath)
				if err != nil {
					continue // Skip files we can't get relative path for
				}
				
				*files = append(*files, FileInfo{
					Name:    relPath,
					ModTime: info.ModTime(),
				})
			}
		}
	}
	
	return nil
}

// GetSortedFiles mock implementation
func (s *MockService) GetSortedFiles() ([]string, error) {
	if s.Error != nil {
		return nil, s.Error
	}
	return s.Files, nil
}