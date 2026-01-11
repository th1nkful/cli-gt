# cli-gt

A Go CLI application that augments Git with opinionated workflow commands.

## Overview

`gt` (Git Tool) is a command-line tool that extends Git with workspace-aware commands to manage branches and streamline common Git operations. It tracks settings per Git workspace, including trunk branch configuration and metadata about managed branches.

## Features

- **Workspace Configuration**: Tracks per-workspace settings like trunk branch and managed branch metadata
- **Command Scaffolding**: Extensible architecture with support for multiple commands
- **Git Integration**: Works within Git repositories to augment standard Git workflows

## Installation

### From Source

```bash
go build -o gt ./cmd/gt
# Move to a directory in your PATH, e.g.:
sudo mv gt /usr/local/bin/
```

## Usage

```bash
gt [command] [flags]
```

### Available Commands

- **`create [branch-name]`** - Create a new branch and commit
- **`pop`** - Undo the current branch and commit, returning the files from the commit/branch to an uncommitted state (effectively undoes "create"). Will not run on trunk branch.
- **`modify`** - Amend the current commit. Will not run on trunk branch.
- **`checkout [branch]`** (alias: `co`) - Checkout to a branch. If no branch is supplied, lists available branches with trunk branch at the bottom and most recently used above that, which you can navigate using up/down arrows to select from the list.
- **`sync`** - Updates trunk branch from origin, rebases local tracked branches onto trunk again. If a local tracked branch no longer exists on origin, prompts for confirmation (y/n) to delete the branch.
- **`restack`** - Restack all managed branches to ensure they are up to date with their parent branch (assuming always trunk for now).
- **`submit`** - Submit the current branch for review (e.g., create/update a pull request). Will not run on trunk branch.
- **`continue`** - Continue a rebase in progress. This is equivalent to `git rebase --continue` and only works if a rebase is currently in progress.

### Configuration

The tool stores workspace settings internally within the `.git` directory (specifically at `.git/gt/config.json`). This ensures:

- **Invisible to users**: Configuration is stored in the same hidden location as other git metadata
- **Device-specific**: Settings are unique to each local workspace and never committed to the repository
- **Seamless**: No need to add entries to `.gitignore` or manage configuration files manually

The configuration includes:

- `trunk_branch`: The main/trunk branch for the repository (default: "main")
- `managed_branches`: A map of branches managed by gt with their metadata

Configuration is automatically created and managed by the tool when you use commands like `create` or `modify`.

## Development

### Building

```bash
go build -o gt ./cmd/gt
```

### Testing

```bash
go test ./...
```

### Project Structure

```
.
├── cmd/
│   └── gt/           # Main CLI entry point
├── internal/
│   ├── commands/     # Command implementations
│   └── config/       # Configuration management
└── go.mod
```

## License

See [LICENSE](LICENSE) file for details.