# Command Line Interface Reference

Complete reference for ob-cli command-line options and usage.

## Synopsis

```bash
ob-cli [file|search-term] [options]
```

## Arguments

- `file`: Direct path to a file to open
- `search-term`: Search term for fuzzy file selection

## Options

### Mode Selection

- `--mode, -m`: Set the mode (tips, obsidian, auto)
  - `tips`: Use Tips vault
  - `obsidian`: Use Obsidian vault  
  - `auto`: Auto-detect based on command name

### Operations

- `--list, -l`: List all files in column format
- `--status, -s`: Show git status
- `--sync`: Sync with remote (stash, pull, pop)

### Information

- `--version, -v`: Show version information
- `--help, -h`: Show help message

### Debugging

- `--debug, -d`: Enable debug output

## Examples

### Interactive Mode

```bash
# Interactive file selection
ob-cli

# Interactive with search term
ob-cli project
```

### Direct File Access

```bash
# Open specific file
ob-cli notes/daily.md

# Open file in subdirectory
ob-cli projects/2024/planning.md
```

### Command Operations

```bash
# List all files
ob-cli --list

# Show git status
ob-cli --status

# Sync with remote
ob-cli --sync
```

### Mode Selection

```bash
# Use Tips mode
ob-cli --mode=tips

# Use Obsidian mode
ob-cli --mode=obsidian

# Auto-detect mode
ob-cli --mode=auto
```

## Environment Variables

- `OBSIDIAN_VAULT`: Path to Obsidian vault
- `TIPS_VAULT`: Path to Tips vault
- `EDITOR`: Preferred editor (defaults to vim/nano/emacs)

## Exit Codes

- `0`: Success
- `1`: General error
- `2`: Invalid arguments
- `3`: Vault not found
- `4`: Editor not found