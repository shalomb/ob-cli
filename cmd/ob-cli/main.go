package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unop/ob-cli/internal/app"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "ob-cli [file|search-term]",
	Short: "Fast Obsidian/Tips CLI with frecency and async git sync",
	Long: `ob-cli is a fast, intelligent CLI tool for managing Obsidian/Tips notes.

Features:
- Frecency-based file selection (recent + frequent files first)
- Async git synchronization (non-blocking)
- Smart file creation with directory support
- Dual mode support (Tips/Obsidian)
- Direct file access or fuzzy search

Examples:
  ob-cli                    # Interactive file selection
  ob-cli notes/daily.md     # Open specific file
  ob-cli project            # Search for files containing "project"
  ob-cli --mode=tips        # Use Tips mode
  ob-cli --list             # List all files
  ob-cli --status           # Show git status
  ob-cli --sync             # Sync with remote`,
	Args: cobra.MaximumNArgs(1),
	RunE: runObCli,
}

var (
	modeFlag     string
	listFlag     bool
	statusFlag   bool
	syncFlag     bool
	versionFlag  bool
	debugFlag    bool
)

func init() {
	rootCmd.Flags().StringVarP(&modeFlag, "mode", "m", "auto", "Mode: tips, obsidian, or auto")
	rootCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List all files")
	rootCmd.Flags().BoolVarP(&statusFlag, "status", "s", false, "Show git status")
	rootCmd.Flags().BoolVarP(&syncFlag, "sync", "", false, "Sync with remote (stash, pull, pop)")
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Show version")
	rootCmd.Flags().BoolVarP(&debugFlag, "debug", "d", false, "Enable debug output")

	// Bind flags to viper
	viper.BindPFlag("mode", rootCmd.Flags().Lookup("mode"))
	viper.BindPFlag("debug", rootCmd.Flags().Lookup("debug"))
}

func runObCli(cmd *cobra.Command, args []string) error {
	// Handle version flag
	if versionFlag {
		fmt.Printf("ob-cli version %s (commit: %s, built: %s)\n", version, commit, date)
		return nil
	}

	// Create app configuration
	config := &app.Config{
		Mode:  modeFlag,
		Debug: debugFlag,
	}

	// Create app instance
	obApp, err := app.New(config)
	if err != nil {
		return fmt.Errorf("failed to create app: %w", err)
	}

	// Handle different command modes
	switch {
	case listFlag:
		return obApp.ListFiles()
	case statusFlag:
		return obApp.ShowGitStatus()
	case syncFlag:
		return obApp.SyncWithRemote()
	default:
		// Interactive mode or direct file access
		var target string
		if len(args) > 0 {
			target = args[0]
		}
		return obApp.RunInteractive(target)
	}
}