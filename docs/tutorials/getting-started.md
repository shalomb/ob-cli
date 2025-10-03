# Getting Started with ob-cli

This tutorial will get you up and running with ob-cli in under 5 minutes.

## Prerequisites

- Go 1.21 or later
- Git
- fzf (for file selection)
- An editor (vim, nano, emacs, or set $EDITOR)

## Installation

```bash
# Clone and build
git clone https://github.com/unop/ob-cli.git
cd ob-cli
make build
make install
```

## Quick Start

### 1. Set up your vault location

```bash
# For Obsidian
export OBSIDIAN_VAULT="/path/to/your/obsidian/vault"

# For Tips
export TIPS_VAULT="/path/to/your/tips/vault"
```

### 2. Test the tool

```bash
# List files in your vault
ob-cli --list

# Open a specific file
ob-cli notes/daily.md

# Interactive file selection
ob-cli
```

### 3. Add to your shell profile

```bash
# Add to ~/.bashrc or ~/.zshrc
echo 'export OBSIDIAN_VAULT="/path/to/your/vault"' >> ~/.bashrc
echo 'alias ob="ob-cli --mode=obsidian"' >> ~/.bashrc
```

## Next Steps

- Read the [How-to Guides](how-to-guides/) for specific tasks
- Check the [Reference](reference/) for complete command documentation
- See [Explanation](explanation/) for how the tool works