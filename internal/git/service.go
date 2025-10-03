package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Service interface for git operations
type Service interface {
	FetchAsync()
	GetStatus() (string, error)
	GetSyncStatus() (behind, ahead int, err error)
	SyncWithRemote() error
}

// RealService handles real git operations
type RealService struct {
	repoDir string
}

// MockService handles mock git operations for testing
type MockService struct {
	FetchResult error
	StatusResult string
	BehindCount int
	AheadCount int
	SyncError error
}

// NewService creates a new real git service
func NewService(repoDir string) Service {
	return &RealService{repoDir: repoDir}
}

// NewMockService creates a new mock git service
func NewMockService(fetchResult error, statusResult string, behindCount, aheadCount int, syncError error) Service {
	return &MockService{
		FetchResult: fetchResult,
		StatusResult: statusResult,
		BehindCount: behindCount,
		AheadCount: aheadCount,
		SyncError: syncError,
	}
}

// FetchAsync starts git fetch in background
func (s *RealService) FetchAsync() {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		cmd := exec.CommandContext(ctx, "git", "fetch", "--all", "--quiet")
		cmd.Dir = s.repoDir
		cmd.Run() // Ignore errors for async operation
	}()
}

// FetchAsync mock implementation
func (s *MockService) FetchAsync() {
	// Mock implementation - does nothing
}

// GetStatus returns git status output
func (s *RealService) GetStatus() (string, error) {
	cmd := exec.Command("git", "status", "--short", "--porcelain")
	cmd.Dir = s.repoDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// GetStatus mock implementation
func (s *MockService) GetStatus() (string, error) {
	return s.StatusResult, s.FetchResult
}

// GetSyncStatus returns behind/ahead counts
func (s *RealService) GetSyncStatus() (behind, ahead int, err error) {
	// Get behind count
	cmd := exec.Command("git", "rev-list", "--count", "HEAD..@{upstream}")
	cmd.Dir = s.repoDir
	output, err := cmd.Output()
	if err != nil {
		behind = 0 // Ignore errors, assume up to date
	} else {
		fmt.Sscanf(strings.TrimSpace(string(output)), "%d", &behind)
	}

	// Get ahead count
	cmd = exec.Command("git", "rev-list", "--count", "@{upstream}..HEAD")
	cmd.Dir = s.repoDir
	output, err = cmd.Output()
	if err != nil {
		ahead = 0 // Ignore errors, assume up to date
	} else {
		fmt.Sscanf(strings.TrimSpace(string(output)), "%d", &ahead)
	}

	return behind, ahead, nil
}

// GetSyncStatus mock implementation
func (s *MockService) GetSyncStatus() (behind, ahead int, err error) {
	return s.BehindCount, s.AheadCount, s.FetchResult
}

// SyncWithRemote performs stash, pull, pop workflow
func (s *RealService) SyncWithRemote() error {
	// Stash
	if err := s.runGitCommand("stash"); err != nil {
		return fmt.Errorf("git stash failed: %w", err)
	}

	// Pull with rebase
	if err := s.runGitCommand("pull", "--rebase"); err != nil {
		// Try to pop stash even if pull failed
		s.runGitCommand("stash", "pop")
		return fmt.Errorf("git pull --rebase failed: %w", err)
	}

	// Pop stash
	if err := s.runGitCommand("stash", "pop"); err != nil {
		return fmt.Errorf("git stash pop failed: %w", err)
	}

	return nil
}

// SyncWithRemote mock implementation
func (s *MockService) SyncWithRemote() error {
	return s.SyncError
}

func (s *RealService) runGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = s.repoDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}