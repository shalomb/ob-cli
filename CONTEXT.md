# ob-cli Context

## Project Status
This is a Go rewrite of the `ob`/`tips` bash script to provide better performance, testing, and maintainability.

## Current State
- ✅ Project structure created (`cmd/`, `internal/`, `pkg/`, `test/`, `docs/`)
- ✅ `go.mod` with dependencies (cobra, viper, etc.)
- ✅ `Makefile` with build, test, install targets
- ✅ BDD specifications written (`docs/BDD_SPECS.md`)
- ✅ Main command structure (`cmd/ob-cli/main.go`)
- ✅ App architecture (`internal/app/app.go`)
- ✅ Service interfaces started (`internal/git/`, `internal/editor/`, `internal/fzf/`)

## What's Working in Bash Version
The current bash `ob` script has these features:
- ✅ Async git sync (non-blocking fzf)
- ✅ Search term support (`ob project` → fzf with "project" query)
- ✅ Direct file access (`ob file.md` → opens directly)
- ✅ File creation with parent directories
- ✅ Tips/Obsidian mode support
- ✅ Command flags (`--list`, `--status`, `--sync`)

## Next Steps for Go Implementation

### 1. Complete Service Implementations
- [ ] `internal/frecency/service.go` - File sorting by access time
- [ ] `internal/fzf/service.go` - Complete fzf integration
- [ ] `internal/editor/service.go` - Complete editor integration
- [ ] `internal/git/service.go` - Complete git operations

### 2. Testing Infrastructure
- [ ] Mock interfaces for all external dependencies
- [ ] Unit tests for business logic
- [ ] Integration tests with temporary directories
- [ ] BDD tests using Ginkgo

### 3. Migration Plan
- [ ] Update `.config/go-tools/*` to include ob-cli
- [ ] Create wrapper scripts for `ob` and `tips` commands
- [ ] Test in parallel with bash version
- [ ] Gradual migration

## Key Design Decisions

### Architecture
- **Clean Architecture**: Services are interfaces, implementations are swappable
- **Dependency Injection**: All external dependencies are mockable
- **Async Operations**: Git operations don't block UI

### Testing Strategy
- **TDD**: Write tests first, then implementation
- **BDD**: Behavior-driven development with Ginkgo
- **Mocks**: All external tools (git, fzf, nvim) are mocked
- **Integration**: Real file operations in isolated test directories

### Performance Goals
- fzf appears within 100ms
- Git operations are non-blocking
- File listing is cached for repeated access

## Current Bash Script Features to Replicate

```bash
# Direct file access
ob existing-file.md

# Search term fallback
ob project-search-term

# Interactive selection
ob

# Command flags
ob --list
ob --status  
ob --sync

# Mode support
ob --mode=tips
ob --mode=obsidian
```

## Mock Requirements

### Git Mock
```go
type GitMock struct {
    FetchResult error
    StatusResult string
    BehindCount int
    AheadCount int
}
```

### Fzf Mock
```go
type FzfMock struct {
    Selection string
    Query string
    ShouldExit bool
}
```

### Editor Mock
```go
type EditorMock struct {
    OpenedFiles []string
    ConcealLevel int
}
```

## File Structure
```
ob-cli/
├── cmd/ob-cli/main.go          # CLI entry point
├── internal/
│   ├── app/app.go              # Main application logic
│   ├── git/service.go          # Git operations
│   ├── editor/service.go       # Editor integration
│   ├── fzf/service.go          # Fzf integration
│   └── frecency/service.go     # File sorting
├── test/
│   ├── bdd/                    # BDD tests
│   └── integration/             # Integration tests
├── docs/BDD_SPECS.md           # Specifications
├── Makefile                     # Build system
└── go.mod                       # Dependencies
```

## Development Commands
```bash
make build          # Build binary
make test           # Run tests
make test-integration # Integration tests
make install        # Install to ~/.local/bin
make dev            # Build and run
make mocks          # Generate mocks
make bdd            # Run BDD tests
```

## Integration Points
- Replace `.local/bin/tips` with Go binary
- Update `.config/go-tools/*` for build/install
- Create symlinks: `ob` → `ob-cli`, `tips` → `ob-cli --mode=tips`
- Maintain same CLI interface for compatibility

## Success Criteria
1. **Performance**: fzf appears within 100ms
2. **Compatibility**: Same CLI interface as bash version
3. **Testing**: 100% test coverage for business logic
4. **Reliability**: Proper error handling and recovery
5. **Maintainability**: Clean, testable code structure