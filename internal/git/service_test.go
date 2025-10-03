package git

import (
	"errors"
	"testing"
)

func TestMockService_FetchAsync(t *testing.T) {
	service := NewMockService(nil, "", 0, 0, nil)
	
	// Should not panic or error
	service.FetchAsync()
	
	// Test passes if no panic occurs
}

func TestMockService_GetStatus(t *testing.T) {
	tests := []struct {
		name           string
		fetchResult    error
		statusResult   string
		expectedStatus string
		expectedError  error
	}{
		{
			name:           "successful status",
			fetchResult:    nil,
			statusResult:   " M file1.md\nA  file2.md",
			expectedStatus: " M file1.md\nA  file2.md",
			expectedError:  nil,
		},
		{
			name:           "status with error",
			fetchResult:    errors.New("git status failed"),
			statusResult:   "",
			expectedStatus: "",
			expectedError:  errors.New("git status failed"),
		},
		{
			name:           "empty status",
			fetchResult:    nil,
			statusResult:   "",
			expectedStatus: "",
			expectedError:  nil,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewMockService(tt.fetchResult, tt.statusResult, 0, 0, nil)
			
			status, err := service.GetStatus()
			
			if status != tt.expectedStatus {
				t.Errorf("Expected status %q, got %q", tt.expectedStatus, status)
			}
			
			if (err != nil) != (tt.expectedError != nil) {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
			
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error message %q, got %q", tt.expectedError.Error(), err.Error())
			}
		})
	}
}

func TestMockService_GetSyncStatus(t *testing.T) {
	tests := []struct {
		name          string
		fetchResult   error
		behindCount   int
		aheadCount    int
		expectedBehind int
		expectedAhead  int
		expectedError error
	}{
		{
			name:          "up to date",
			fetchResult:   nil,
			behindCount:   0,
			aheadCount:    0,
			expectedBehind: 0,
			expectedAhead:  0,
			expectedError: nil,
		},
		{
			name:          "behind remote",
			fetchResult:   nil,
			behindCount:   3,
			aheadCount:    0,
			expectedBehind: 3,
			expectedAhead:  0,
			expectedError: nil,
		},
		{
			name:          "ahead of remote",
			fetchResult:   nil,
			behindCount:   0,
			aheadCount:    2,
			expectedBehind: 0,
			expectedAhead:  2,
			expectedError: nil,
		},
		{
			name:          "diverged",
			fetchResult:   nil,
			behindCount:   1,
			aheadCount:    1,
			expectedBehind: 1,
			expectedAhead:  1,
			expectedError: nil,
		},
		{
			name:          "git error",
			fetchResult:   errors.New("git fetch failed"),
			behindCount:   0,
			aheadCount:    0,
			expectedBehind: 0,
			expectedAhead:  0,
			expectedError: errors.New("git fetch failed"),
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewMockService(tt.fetchResult, "", tt.behindCount, tt.aheadCount, nil)
			
			behind, ahead, err := service.GetSyncStatus()
			
			if behind != tt.expectedBehind {
				t.Errorf("Expected behind %d, got %d", tt.expectedBehind, behind)
			}
			
			if ahead != tt.expectedAhead {
				t.Errorf("Expected ahead %d, got %d", tt.expectedAhead, ahead)
			}
			
			if (err != nil) != (tt.expectedError != nil) {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
			
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error message %q, got %q", tt.expectedError.Error(), err.Error())
			}
		})
	}
}

func TestMockService_SyncWithRemote(t *testing.T) {
	tests := []struct {
		name         string
		syncError    error
		expectedError error
	}{
		{
			name:         "successful sync",
			syncError:    nil,
			expectedError: nil,
		},
		{
			name:         "sync error",
			syncError:    errors.New("git pull failed"),
			expectedError: errors.New("git pull failed"),
		},
		{
			name:         "stash error",
			syncError:    errors.New("git stash failed"),
			expectedError: errors.New("git stash failed"),
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewMockService(nil, "", 0, 0, tt.syncError)
			
			err := service.SyncWithRemote()
			
			if (err != nil) != (tt.expectedError != nil) {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
			
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error message %q, got %q", tt.expectedError.Error(), err.Error())
			}
		})
	}
}