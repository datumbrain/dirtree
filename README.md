# dirtree

[![CI](https://github.com/datumbrain/dirtree/actions/workflows/ci.yml/badge.svg)](https://github.com/datumbrain/dirtree/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/datumbrain/dirtree)](https://goreportcard.com/report/github.com/datumbrain/dirtree)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A lightweight command-line utility that displays directory tree structures while intelligently respecting `.gitignore` patterns. Perfect for documenting project structures, exploring codebases, and generating clean directory listings.

## Features

- **Gitignore-Aware**: Automatically respects `.gitignore` files at all directory levels
- **Hierarchical Patterns**: Inherits and combines gitignore rules from parent directories
- **Smart Sorting**: Organizes output with directories first, then files, then dotfiles (all alphabetically sorted)
- **Clean Output**: Uses ASCII tree characters (├──, └──, │) for clear visual hierarchy
- **Fast & Lightweight**: Single binary with minimal dependencies
- **Cross-Platform**: Works on Linux, macOS, and Windows

## Installation

### Using Go Install (Recommended)

```bash
go install github.com/datumbrain/dirtree@latest
```

### Build from Source

```bash
git clone https://github.com/datumbrain/dirtree.git
cd dirtree
go build
```

## Usage

### Basic Usage

Display the tree structure of the current directory:

```bash
dirtree
```

### Specify a Directory

Display the tree structure of a specific directory:

```bash
dirtree /path/to/directory
```

### Example Output

```
.
├── cmd
│   └── main.go
├── internal
│   ├── config.go
│   └── utils.go
├── pkg
│   └── api
│       └── handler.go
├── go.mod
├── go.sum
├── README.md
└── .gitignore
```

## How It Works

### Gitignore Handling

dirtree reads `.gitignore` files at each directory level and respects all standard gitignore patterns:

- **Wildcard patterns**: `*.log`, `*.tmp`
- **Directory patterns**: `node_modules/`, `build/`
- **Nested patterns**: `src/**/*.test.js`
- **Negation patterns**: `!important.log`

Gitignore rules are inherited from parent directories, just like Git itself handles them.

### Automatic Exclusions

The following are always excluded from the output:
- `.git` directory

### File Sorting

Files and directories are sorted in the following order:
1. **Directories** (alphabetically)
2. **Regular files** (alphabetically)
3. **Dotfiles** (alphabetically)

This ensures consistent, readable output across different systems.

## Use Cases

- **Documentation**: Generate directory structure for README files or documentation
- **Code Review**: Quickly understand project organization
- **Project Planning**: Visualize codebase structure before making changes
- **Onboarding**: Help new team members understand project layout
- **Git-Aware Exploration**: See only tracked/relevant files, ignoring build artifacts

## Comparison with `tree` Command

| Feature | dirtree | tree |
|---------|---------|------|
| Respects .gitignore | ✅ Yes | ❌ No |
| Cross-platform | ✅ Yes | ⚠️ Limited |
| Installation | Go install | OS package manager |
| Size | ~2MB | Varies |
| Hierarchical gitignore | ✅ Yes | N/A |

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Start for Contributors

```bash
# Clone the repository
git clone https://github.com/datumbrain/dirtree.git
cd dirtree

# Run tests
go test -v ./...

# Build
go build
```

## Troubleshooting

### Binary not in PATH

After installation, ensure your Go bin directory is in your PATH:

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:$(go env GOPATH)/bin
```

### Permission Denied Errors

If you encounter permission errors when reading directories, ensure you have read permissions for the target directory.

### Gitignore Not Working

Make sure your `.gitignore` file:
- Is named exactly `.gitignore` (case-sensitive)
- Uses proper gitignore pattern syntax
- Is in the correct directory relative to the files you want to ignore

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Authors

- [Fahad Siddiqui](https://github.com/fahadsiddiqui)

## Acknowledgments

- Uses [go-gitignore](https://github.com/sabhiram/go-gitignore) for gitignore pattern matching
- Inspired by the Unix `tree` command with Git-awareness added
