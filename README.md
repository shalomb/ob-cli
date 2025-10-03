# ob-cli

Fast CLI tool for Obsidian and Tips note management with frecency-based file selection and async git synchronization.

## Features

- Frecency-based file selection (recent files first)
- Async git synchronization (non-blocking)
- Smart file creation with directory support
- Editor-agnostic (respects $EDITOR)
- Zero configuration (works out of the box)

## Installation

```bash
# Build from source
make build
make install

# Or install directly
go install github.com/unop/ob-cli/cmd/ob-cli@latest
```

## Usage

### Obsidian Mode
```bash
# Set your Obsidian vault location
export OBSIDIAN_VAULT="/path/to/your/vault"

# Use the tool
ob-cli --mode=obsidian
ob-cli notes/daily.md
ob-cli project-search-term
```

### Tips Mode
```bash
# Set your Tips vault location
export TIPS_VAULT="/path/to/your/tips"

# Use the tool
ob-cli --mode=tips
ob-cli notes/daily.md
ob-cli project-search-term
```

### Command Options
```bash
ob-cli                    # Interactive file selection
ob-cli file.md            # Open specific file
ob-cli search-term        # Search for files
ob-cli --list             # List all files
ob-cli --status           # Show git status
ob-cli --sync             # Sync with remote
```

## Configuration

The tool uses environment variables for configuration:

- `OBSIDIAN_VAULT`: Path to your Obsidian vault
- `TIPS_VAULT`: Path to your Tips vault
- `EDITOR`: Your preferred editor (defaults to vim/nano/emacs)

If not set, the tool will look for vaults in common locations:
- `~/obsidian`
- `~/Documents/Obsidian`
- `~/tips`

## Development

```bash
# Run tests
make test

# Run integration tests
make test-integration

# Build and run
make dev

# Lint code
make lint
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

MIT