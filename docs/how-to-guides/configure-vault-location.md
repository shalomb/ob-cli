# How to Configure Vault Location

This guide shows you how to configure ob-cli to work with your specific vault locations.

## Method 1: Environment Variables (Recommended)

### For Obsidian Vaults

```bash
# Set the environment variable
export OBSIDIAN_VAULT="/path/to/your/obsidian/vault"

# Verify it's set
echo $OBSIDIAN_VAULT

# Use the tool
ob-cli --mode=obsidian
```

### For Tips Vaults

```bash
# Set the environment variable
export TIPS_VAULT="/path/to/your/tips/vault"

# Verify it's set
echo $TIPS_VAULT

# Use the tool
ob-cli --mode=tips
```

### Make it Permanent

Add to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
# Add these lines to your shell profile
export OBSIDIAN_VAULT="/path/to/your/obsidian/vault"
export TIPS_VAULT="/path/to/your/tips/vault"
```

## Method 2: Auto-Discovery

If you don't set environment variables, ob-cli will automatically look for vaults in these locations:

### Obsidian Vaults
- `~/obsidian`
- `~/Documents/Obsidian`
- `~/Documents/Obsidian Vaults`

### Tips Vaults
- `~/tips`
- `~/Documents/tips`
- `~/Documents/Tips`

## Troubleshooting

### Vault Not Found

```bash
# Check if the path exists
ls -la "$OBSIDIAN_VAULT"

# Verify it's a valid vault (has .obsidian directory for Obsidian)
ls -la "$OBSIDIAN_VAULT/.obsidian"
```

### Multiple Vaults

If you have multiple vaults, you can switch between them:

```bash
# Use different vaults for different projects
OBSIDIAN_VAULT="/project1/vault" ob-cli --mode=obsidian
OBSIDIAN_VAULT="/project2/vault" ob-cli --mode=obsidian
```