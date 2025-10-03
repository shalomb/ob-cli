package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiscoverer_DiscoverObsidianVault(t *testing.T) {
	discoverer := NewDiscoverer()
	
	t.Run("with environment variable", func(t *testing.T) {
		// Create temporary directory
		tempDir := t.TempDir()
		
		// Create .obsidian directory to make it a valid vault
		obsidianDir := filepath.Join(tempDir, ".obsidian")
		err := os.MkdirAll(obsidianDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create .obsidian directory: %v", err)
		}
		
		// Set environment variable
		os.Setenv("OBSIDIAN_VAULT", tempDir)
		defer os.Unsetenv("OBSIDIAN_VAULT")
		
		vaultPath, err := discoverer.DiscoverObsidianVault()
		if err != nil {
			t.Errorf("DiscoverObsidianVault failed: %v", err)
		}
		
		if vaultPath != tempDir {
			t.Errorf("Expected vault path %s, got %s", tempDir, vaultPath)
		}
	})
	
	t.Run("with invalid environment variable", func(t *testing.T) {
		// Set invalid environment variable
		os.Setenv("OBSIDIAN_VAULT", "/non/existent/path")
		defer os.Unsetenv("OBSIDIAN_VAULT")
		
		_, err := discoverer.DiscoverObsidianVault()
		if err == nil {
			t.Error("Expected error for invalid vault path, got nil")
		}
		
		expectedError := "invalid vault path in OBSIDIAN_VAULT"
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("Expected error to contain %q, got %q", expectedError, err.Error())
		}
	})
}

func TestDiscoverer_DiscoverTipsVault(t *testing.T) {
	discoverer := NewDiscoverer()
	
	t.Run("with environment variable", func(t *testing.T) {
		// Create temporary directory
		tempDir := t.TempDir()
		
		// Create a markdown file to make it a valid vault
		mdFile := filepath.Join(tempDir, "test.md")
		err := os.WriteFile(mdFile, []byte("# Test"), 0644)
		if err != nil {
			t.Fatalf("Failed to create markdown file: %v", err)
		}
		
		// Set environment variable
		os.Setenv("TIPS_VAULT", tempDir)
		defer os.Unsetenv("TIPS_VAULT")
		
		vaultPath, err := discoverer.DiscoverTipsVault()
		if err != nil {
			t.Errorf("DiscoverTipsVault failed: %v", err)
		}
		
		if vaultPath != tempDir {
			t.Errorf("Expected vault path %s, got %s", tempDir, vaultPath)
		}
	})
	
	t.Run("with invalid environment variable", func(t *testing.T) {
		// Set invalid environment variable
		os.Setenv("TIPS_VAULT", "/non/existent/path")
		defer os.Unsetenv("TIPS_VAULT")
		
		_, err := discoverer.DiscoverTipsVault()
		if err == nil {
			t.Error("Expected error for invalid vault path, got nil")
		}
		
		expectedError := "invalid vault path in TIPS_VAULT"
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("Expected error to contain %q, got %q", expectedError, err.Error())
		}
	})
}

func TestDiscoverer_IsValidVault(t *testing.T) {
	discoverer := NewDiscoverer()
	
	tests := []struct {
		name     string
		setup    func(string) error
		expected bool
	}{
		{
			name: "valid Obsidian vault",
			setup: func(path string) error {
				obsidianDir := filepath.Join(path, ".obsidian")
				return os.MkdirAll(obsidianDir, 0755)
			},
			expected: true,
		},
		{
			name: "valid Tips vault",
			setup: func(path string) error {
				mdFile := filepath.Join(path, "test.md")
				return os.WriteFile(mdFile, []byte("# Test"), 0644)
			},
			expected: true,
		},
		{
			name: "empty directory",
			setup: func(path string) error {
				return nil // Just create the directory
			},
			expected: false,
		},
		{
			name: "non-existent directory",
			setup: func(path string) error {
				return os.RemoveAll(path) // Remove the directory
			},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			
			err := tt.setup(tempDir)
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			
			result := discoverer.IsValidVault(tempDir)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDiscoverer_FindVaultInDirectory(t *testing.T) {
	discoverer := NewDiscoverer()
	
	t.Run("finds valid vault", func(t *testing.T) {
		// Create temporary directory structure
		baseDir := t.TempDir()
		vaultDir := filepath.Join(baseDir, "MyVault")
		
		// Create .obsidian directory to make it a valid vault
		obsidianDir := filepath.Join(vaultDir, ".obsidian")
		err := os.MkdirAll(obsidianDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create .obsidian directory: %v", err)
		}
		
		vaultPath, err := discoverer.FindVaultInDirectory(baseDir)
		if err != nil {
			t.Errorf("FindVaultInDirectory failed: %v", err)
		}
		
		if vaultPath != vaultDir {
			t.Errorf("Expected vault path %s, got %s", vaultDir, vaultPath)
		}
	})
	
	t.Run("no vault found", func(t *testing.T) {
		// Create empty directory
		baseDir := t.TempDir()
		
		_, err := discoverer.FindVaultInDirectory(baseDir)
		if err == nil {
			t.Error("Expected error when no vault found, got nil")
		}
	})
}