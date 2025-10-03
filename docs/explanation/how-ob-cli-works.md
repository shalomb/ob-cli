# How ob-cli Works

This document explains the internal architecture and behavior of ob-cli.

## Architecture

ob-cli follows a clean architecture pattern with these components:

### Services

- **Vault Discovery**: Finds and validates vault locations
- **Frecency**: Sorts files by access time (recent first)
- **Git**: Handles async synchronization and status
- **FZF**: Provides fuzzy file selection
- **Editor**: Opens files in user's preferred editor

### Data Flow

1. **Vault Discovery**: Resolve vault location from environment or auto-discovery
2. **File Listing**: Get all markdown files sorted by frecency
3. **Git Sync**: Start background git fetch
4. **File Selection**: Present files via fzf or direct access
5. **File Operations**: Create files/directories as needed
6. **Editor Launch**: Open file in user's editor

## Vault Discovery

### Priority Order

1. **Environment Variables**: `OBSIDIAN_VAULT` or `TIPS_VAULT`
2. **Common Locations**: Check standard paths (2 levels deep)
3. **Default Fallback**: Use `~/obsidian` or `~/tips`

### Validation

- **Obsidian Vaults**: Must contain `.obsidian` directory
- **Tips Vaults**: Must contain at least one `.md` file

## Frecency Algorithm

Simple access-time sorting:

- Files sorted by modification time (most recent first)
- No complex weighting or frequency calculations
- Fast and predictable behavior

## Git Integration

### Async Operations

- Git fetch runs in background
- Non-blocking file selection
- Status reported after selection

### Sync Workflow

1. `git stash` - Save local changes
2. `git pull --rebase` - Update from remote
3. `git stash pop` - Restore local changes

## Editor Integration

### Editor Selection

1. Check `$EDITOR` environment variable
2. Fallback to common editors: `edit`, `vim`, `nano`, `emacs`
3. Error if no editor found

### Editor Agnostic

- No editor-specific configuration
- Respects user's editor settings
- Works with any editor

## Performance

### Design Goals

- fzf appears within 100ms
- Git operations don't block UI
- File listing cached for repeated access

### Optimization Strategies

- Fast vault discovery (2-level deep scan)
- Efficient file sorting (access time only)
- Background git operations
- Minimal filesystem access