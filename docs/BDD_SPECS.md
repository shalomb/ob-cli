# ob-cli BDD Specifications

## Overview
ob-cli is a fast, intelligent CLI tool for managing Obsidian/Tips notes with frecency-based file selection and async git synchronization.

## Core Features

### 1. File Selection with Frecency
**As a** note-taking user  
**I want** my most recently accessed files to appear first  
**So that** I can quickly access my current work

**Scenario: Recent files appear first**
```gherkin
Given I have files with different access times
When I run "ob-cli" without arguments
Then the most recently accessed files should appear at the top
And the sorting should be fast and consistent
And no complex configuration should be needed
```

**Scenario: Direct file access**
```gherkin
Given I have a file "notes/daily.md"
When I run "ob-cli notes/daily.md"
Then the file should open directly in the editor
And no fzf selection should be shown
```

**Scenario: Search term fallback**
```gherkin
Given I have files containing "project"
When I run "ob-cli project"
And "project" is not an existing file
Then fzf should open with "project" as the search query
And files matching "project" should be highlighted
```

### 2. Async Git Synchronization
**As a** collaborative note-taker  
**I want** to stay synchronized with remote changes  
**So that** I don't miss important updates

**Scenario: Background git sync**
```gherkin
Given I run "ob-cli" without arguments
When fzf appears
Then git fetch should start in the background
And fzf should show "Syncing..." in the header
And the tool should not wait for git to complete
```

**Scenario: Status reporting after selection**
```gherkin
Given git fetch is running in the background
When I select a file from fzf
And git fetch has completed
Then I should see repository status (behind/ahead)
And I should be prompted to sync if behind
```

### 3. File Creation
**As a** note-taking user  
**I want** to create new files easily  
**So that** I can start new notes quickly

**Scenario: Create new file**
```gherkin
Given I run "ob-cli" without arguments
When I type "new-note.md" in fzf
And press Enter
Then the file "new-note.md" should be created
And parent directories should be created if needed
And the file should open in the editor
```

**Scenario: Create nested file**
```gherkin
Given I run "ob-cli" without arguments
When I type "projects/2024/new-project.md" in fzf
And press Enter
Then the directory "projects/2024/" should be created
And the file "projects/2024/new-project.md" should be created
And the file should open in the editor
```

### 4. Mode Support
**As a** user with different note systems  
**I want** to switch between Tips and Obsidian modes  
**So that** I can work with different note repositories

**Scenario: Tips mode with environment variable**
```gherkin
Given I have set TIPS_VAULT="/my/custom/tips"
When I run "ob-cli --mode=tips"
Then the tool should work in the /my/custom/tips directory
And the editor should respect my $EDITOR setting
```

**Scenario: Obsidian mode with environment variable**
```gherkin
Given I have set OBSIDIAN_VAULT="/my/custom/vault"
When I run "ob-cli --mode=obsidian"
Then the tool should work in the /my/custom/vault directory
And the editor should respect my $EDITOR setting
```

**Scenario: Auto-discovery fallback**
```gherkin
Given I have not set any vault environment variables
When I run "ob-cli --mode=obsidian"
Then the tool should look for vaults in common locations
And should find a valid vault or use default location
```

### 5. Command Line Interface
**As a** power user  
**I want** various command options  
**So that** I can customize the tool's behavior

**Scenario: List files**
```gherkin
Given I run "ob-cli --list"
Then I should see all files in a column format
And the output should be paginated with less
```

**Scenario: Git status**
```gherkin
Given I run "ob-cli --status"
Then I should see git status output
And the tool should exit after showing status
```

**Scenario: Git sync**
```gherkin
Given I run "ob-cli --sync"
Then the tool should run "git stash"
And then "git pull --rebase"
And then "git stash pop"
And exit after completion
```

## Technical Requirements

### Performance
- fzf should appear within 100ms
- git operations should not block the UI
- file listing should be cached for repeated access

### Simplicity
- No complex configuration files
- Environment variable configuration (12-factor app)
- Simple, predictable file sorting (by access time)
- Clear, helpful error messages
- Fast vault discovery (2-level deep scan)

### Error Handling
- Graceful handling of git failures
- Clear error messages for file creation failures
- Simple fallback behavior when external tools are missing

### Testing
- All external dependencies (git, fzf, nvim) must be mockable
- Integration tests should use temporary directories
- Unit tests should cover all business logic

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

## Acceptance Criteria

1. **Speed**: fzf appears within 100ms of command execution
2. **Simplicity**: Files are sorted by access time (simple, predictable)
3. **Async**: Git operations don't block user interaction
4. **Creation**: New files and directories are created correctly
5. **Modes**: Both Tips and Obsidian modes work correctly
6. **CLI**: All command-line options work as specified
7. **Testing**: 100% test coverage for business logic
8. **Mocking**: All external dependencies are properly mocked
9. **Zero Config**: Tool works out of the box with sensible defaults