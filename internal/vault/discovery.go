package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Discoverer handles vault discovery
type Discoverer struct{}

// NewDiscoverer creates a new vault discoverer
func NewDiscoverer() *Discoverer {
	return &Discoverer{}
}

// DiscoverObsidianVault finds the Obsidian vault location (12-factor approach)
func (d *Discoverer) DiscoverObsidianVault() (string, error) {
	// 1. Environment variable (12-factor app principle)
	if vaultPath := os.Getenv("OBSIDIAN_VAULT"); vaultPath != "" {
		if d.IsValidVault(vaultPath) {
			return vaultPath, nil
		}
		return "", fmt.Errorf("invalid vault path in OBSIDIAN_VAULT: %s", vaultPath)
	}

	// 2. Fast fallback discovery (2-level deep only)
	return d.FastObsidianDiscovery()
}

// FastObsidianDiscovery performs fast vault discovery
func (d *Discoverer) FastObsidianDiscovery() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check common locations (2 levels deep max)
	commonPaths := []string{
		filepath.Join(homeDir, "obsidian"),
		filepath.Join(homeDir, "Documents", "Obsidian"),
		filepath.Join(homeDir, "Documents", "Obsidian Vaults"),
	}

	for _, path := range commonPaths {
		if d.IsValidVault(path) {
			return path, nil
		}
	}

	// Default fallback
	return filepath.Join(homeDir, "obsidian"), nil
}

// DiscoverTipsVault finds the Tips vault location (12-factor approach)
func (d *Discoverer) DiscoverTipsVault() (string, error) {
	// 1. Environment variable (12-factor app principle)
	if vaultPath := os.Getenv("TIPS_VAULT"); vaultPath != "" {
		if d.IsValidVault(vaultPath) {
			return vaultPath, nil
		}
		return "", fmt.Errorf("invalid vault path in TIPS_VAULT: %s", vaultPath)
	}

	// 2. Fast fallback discovery (2-level deep only)
	return d.FastTipsDiscovery()
}

// FastTipsDiscovery performs fast Tips vault discovery
func (d *Discoverer) FastTipsDiscovery() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check common locations (2 levels deep max)
	commonPaths := []string{
		filepath.Join(homeDir, "tips"),
		filepath.Join(homeDir, "Documents", "tips"),
		filepath.Join(homeDir, "Documents", "Tips"),
		filepath.Join(homeDir, "Notes"),
		filepath.Join(homeDir, "notes"),
	}

	for _, path := range commonPaths {
		if d.IsValidVault(path) {
			return path, nil
		}
	}

	// Default fallback
	return filepath.Join(homeDir, "tips"), nil
}

// FindVaultInDirectory searches for vaults in a directory
func (d *Discoverer) FindVaultInDirectory(basePath string) (string, error) {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return "", fmt.Errorf("directory does not exist")
	}

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return "", err
	}

	// Look for directories that might be vaults
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		vaultPath := filepath.Join(basePath, entry.Name())
		if d.IsValidVault(vaultPath) {
			return vaultPath, nil
		}
	}

	return "", fmt.Errorf("no vault found in directory")
}

// searchForObsidianConfig searches for .obsidian directories
func (d *Discoverer) searchForObsidianConfig(startPath string) (string, error) {
	var foundVault string
	
	err := filepath.Walk(startPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors, continue searching
		}

		// Look for .obsidian directory
		if info.IsDir() && info.Name() == ".obsidian" {
			// Found .obsidian directory, parent should be the vault
			vaultPath := filepath.Dir(path)
			if d.IsValidVault(vaultPath) {
				foundVault = vaultPath
				return filepath.SkipDir // Stop searching
			}
		}

		// Skip hidden directories and common non-vault directories
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || 
			info.Name() == "node_modules" || 
			info.Name() == ".git") {
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if foundVault != "" {
		return foundVault, nil
	}

	return "", fmt.Errorf("no Obsidian vault found")
}

// isValidVault checks if a path is a valid vault
func (d *Discoverer) isValidVault(path string) bool {
	// Check if directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	// Check if it's a directory
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}

	// For Obsidian vaults, check for .obsidian directory
	obsidianConfig := filepath.Join(path, ".obsidian")
	if _, err := os.Stat(obsidianConfig); err == nil {
		return true
	}

	// For Tips vaults, check for markdown files
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	// Look for at least one .md file
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			return true
		}
	}

	return false
}

// IsValidVault is a public method for external use
func (d *Discoverer) IsValidVault(path string) bool {
	return d.isValidVault(path)
}