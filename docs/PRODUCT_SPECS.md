# ob-cli Product Specifications

## Vision
**Simple tool, low-friction, low-config, "do the right thing for me"**

ob-cli is a fast, intelligent CLI tool that replaces the bash `ob`/`tips` script with better performance, testing, and maintainability while keeping the same simple, opinionated approach.

## Core Philosophy

### 1. **Zero Configuration**
- Works out of the box with sensible defaults
- No config files to manage
- Auto-detects Tips vs Obsidian mode
- **Intelligent vault discovery** - finds your vaults automatically

### 2. **"Do the Right Thing"**
- Recent files appear first (simple access-time sorting)
- Git sync happens automatically in background
- Creates directories when needed
- Opens files in the right editor with right settings

### 3. **Low Friction**
- `ob-cli` → just works (interactive mode)
- `ob-cli file.md` → opens directly
- `ob-cli project` → smart search fallback
- Fast response times (< 100ms)

## Feature Specifications

### File Selection
- **Simple sorting**: By access time (mtime/atime)
- **No complex algorithms**: Just "recent files first"
- **Fast and predictable**: Users understand the behavior
- **Cached results**: Repeated access is instant

### Git Integration
- **Background sync**: Non-blocking git fetch
- **Status reporting**: Shows behind/ahead after selection
- **Simple sync**: `--sync` does stash/pull/pop
- **Graceful failures**: Continues working if git fails

### File Operations
- **Smart creation**: Creates parent directories automatically
- **Direct access**: `ob-cli file.md` opens immediately
- **Search fallback**: `ob-cli term` searches if file doesn't exist
- **Editor-agnostic**: Respects `$EDITOR` and user's editor config

### Command Interface
```bash
ob-cli                    # Interactive file selection
ob-cli notes/daily.md     # Open specific file
ob-cli project            # Search for files containing "project"
ob-cli --mode=tips        # Use Tips mode
ob-cli --list             # List all files
ob-cli --status           # Show git status
ob-cli --sync             # Sync with remote
```

## Technical Requirements

### Performance
- **fzf appears within 100ms**
- **Git operations are non-blocking**
- **File listing is cached**

### Simplicity
- **No configuration files**
- **Sensible defaults everywhere**
- **Simple, predictable behavior**
- **Clear error messages**

### Reliability
- **Graceful error handling**
- **Fallback behavior for missing tools**
- **Comprehensive testing**
- **Mockable external dependencies**

## Success Metrics

1. **Speed**: fzf appears within 100ms
2. **Simplicity**: Zero configuration required
3. **Reliability**: Works consistently across environments
4. **Compatibility**: Drop-in replacement for bash version
5. **Maintainability**: Clean, testable code structure

## Non-Goals

- Complex frecency algorithms
- Extensive configuration options
- Advanced git workflows
- Plugin systems
- Multiple repository support
- Custom editor integration

## Migration Strategy

1. **Parallel testing**: Run alongside bash version
2. **Gradual replacement**: Replace `ob` and `tips` commands
3. **Backward compatibility**: Same CLI interface
4. **Rollback plan**: Keep bash version available

## Implementation Approach

### Phase 1: Core Services
- [ ] Simple frecency service (access-time sorting)
- [ ] fzf integration service
- [ ] Complete git service
- [ ] Complete editor service

### Phase 2: Testing
- [ ] Unit tests for all services
- [ ] Integration tests with real files
- [ ] Mock implementations
- [ ] BDD test framework

### Phase 3: Integration
- [ ] Complete app integration
- [ ] Performance optimization
- [ ] Error handling refinement
- [ ] Documentation

### Phase 4: Migration
- [ ] Build and install scripts
- [ ] Wrapper command creation
- [ ] Parallel testing
- [ ] Gradual rollout