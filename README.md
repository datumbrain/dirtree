# dirtree

[![CI](https://github.com/datumbrain/dirtree/actions/workflows/ci.yml/badge.svg)](https://github.com/datumbrain/dirtree/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/datumbrain/dirtree)](https://goreportcard.com/report/github.com/datumbrain/dirtree)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A lightweight command-line utility that displays directory tree structures while intelligently respecting `.gitignore` patterns. Perfect for documenting project structures, exploring codebases, and generating clean directory listings.

## Features

- **Multiple Output Formats**: Text (traditional), JSON (machine-readable), and Markdown (LLM-optimized)
- **LLM-Optimized Mode**: Special `--llm` flag for AI agent consumption with sensible defaults
- **Gitignore-Aware**: Automatically respects `.gitignore` files at all directory levels
- **Hierarchical Patterns**: Inherits and combines gitignore rules from parent directories
- **Smart Sorting**: Organizes output with directories first, then files, then dotfiles (all alphabetically sorted)
- **Colorized Output**: Beautiful ANSI colors for directories, files, and file types with auto-detection
- **Depth Limiting**: Control how deep to traverse with `-L` flag
- **File Count Limiting**: Limit total files displayed with `--max-files` flag
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

After building, install system-wide:
```bash
go install
```

## Usage

### Basic Usage

Display the tree structure of the current directory:

```bash
dirtree
```

### Command-Line Options

```bash
dirtree [options] [directory]

Options:
  -L int
        Maximum depth of directory tree (0 = unlimited)
  -level int
        Maximum depth of directory tree (0 = unlimited)
  -color string
        Colorize output: auto, always, never (default "auto")
  -format string
        Output format: text, json, markdown (or md) (default "text")
  -o string
        Output format: text, json, markdown (or md) (default "text")
  -max-files int
        Maximum number of files to display (0 = unlimited)
  -llm
        LLM-optimized output (equivalent to: --format markdown --level 3 --max-files 200)
```

### Examples

**Specify a directory:**
```bash
dirtree /path/to/directory
```

**Limit depth to 2 levels:**
```bash
dirtree -L 2
```

**Show only top-level items:**
```bash
dirtree --level 1
```

**Force colored output:**
```bash
dirtree --color always
```

**Disable colors:**
```bash
dirtree --color never
```

**Combine options:**
```bash
dirtree --color always -L 3 /path/to/directory
```

**Respect NO_COLOR environment variable:**
```bash
NO_COLOR=1 dirtree
```

**JSON output for programmatic use:**
```bash
dirtree --format json
# or using the short flag
dirtree -o json
```

**Markdown output for documentation:**
```bash
dirtree --format markdown
# or using the alias
dirtree -o md
```

**LLM-optimized mode (markdown, depth 3, max 200 files):**
```bash
dirtree --llm
```

**Limit number of files displayed:**
```bash
dirtree --max-files 50
```

**Combine format with other options:**
```bash
dirtree --format json -L 2 --max-files 100
```

### Example Output

**Without colors:**
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

**With colors** (uses different colors for directories, files, and file types):
- Directories appear in **blue and bold**
- Regular files appear in white
- Dotfiles appear in gray
- Markdown/documentation files appear in cyan
- Image files appear in magenta
- Executable files appear in green and bold

## Output Formats

dirtree supports three output formats to suit different use cases:

### Text Format (Default)

Traditional tree output with ASCII art. Perfect for terminal display and human reading.

```bash
dirtree --format text
```

Output:
```
.
├── cmd
│   └── main.go
├── pkg
│   └── api
│       └── handler.go
├── README.md
└── go.mod
```

### JSON Format

Machine-readable format ideal for programmatic processing and integration with other tools.

```bash
dirtree --format json
```

Output:
```json
{
  "root": ".",
  "depth": 0,
  "stats": {
    "files": 3,
    "dirs": 3
  },
  "nodes": [
    {"path": "cmd", "type": "dir"},
    {"path": "cmd/main.go", "type": "file"},
    {"path": "pkg", "type": "dir"},
    {"path": "pkg/api", "type": "dir"},
    {"path": "pkg/api/handler.go", "type": "file"},
    {"path": "README.md", "type": "file"},
    {"path": "go.mod", "type": "file"}
  ]
}
```

### Markdown Format

Optimized for documentation and LLM consumption. Groups files by directory with clear sections.

```bash
dirtree --format markdown
```

Output:
```markdown
# Directory Structure: .

## Files
### Root Directory
- `README.md`
- `go.mod`

### cmd/
- `main.go`

### pkg/api/
- `handler.go`

## Directories
- `cmd/`
- `pkg/`
- `pkg/api/`

---
**Total:** 3 files, 3 directories
```

### LLM Mode

Special mode optimized for AI agents with sensible defaults:

```bash
dirtree --llm
```

This is equivalent to: `--format markdown --level 3 --max-files 200`

Benefits for LLM usage:
- Markdown format is easier for LLMs to parse
- Depth limit (3) prevents context overflow
- File limit (200) keeps output manageable
- Grouped structure helps with understanding project organization

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

### Color Output

dirtree uses intelligent color coding to make output easier to scan:

- **Directories**: Blue and bold
- **Regular files**: White
- **Dotfiles**: Gray (dim)
- **Executable files**: Green and bold
- **Image files** (.png, .jpg, .gif, etc.): Magenta
- **Documentation** (.md, .txt, .pdf, etc.): Cyan

Color output respects:
- `--color` flag (`auto`, `always`, `never`)
- `NO_COLOR` environment variable ([no-color.org](https://no-color.org/))
- Terminal capabilities (auto-detection in `auto` mode)

### Depth Limiting

Control traversal depth with the `-L` or `--level` flag:

```bash
dirtree -L 2    # Only show 2 levels deep
```

This is useful for:
- Getting a quick overview of large projects
- Focusing on top-level structure
- Avoiding deeply nested directories

## Use Cases

- **Documentation**: Generate directory structure for README files (Markdown format) or documentation
- **LLM Context**: Feed project structure to AI assistants with `--llm` flag
- **Code Review**: Quickly understand project organization
- **Project Planning**: Visualize codebase structure before making changes
- **Onboarding**: Help new team members understand project layout
- **Git-Aware Exploration**: See only tracked/relevant files, ignoring build artifacts
- **CI/CD Integration**: Use JSON format for automated project analysis
- **Build Tools**: Parse JSON output for custom tooling and scripts

## Comparison with `tree` Command

| Feature | dirtree | tree |
|---------|---------|------|
| Respects .gitignore | ✅ Yes | ❌ No |
| Hierarchical gitignore | ✅ Yes | N/A |
| Multiple output formats | ✅ Yes (text/json/md) | ⚠️ Limited |
| LLM-optimized mode | ✅ Yes | ❌ No |
| JSON output | ✅ Yes | ⚠️ Limited |
| Markdown output | ✅ Yes | ❌ No |
| File count limiting | ✅ Yes | ❌ No |
| Colorized output | ✅ Yes | ✅ Yes |
| Depth limiting | ✅ Yes (-L) | ✅ Yes (-L) |
| Cross-platform | ✅ Yes | ⚠️ Limited |
| Installation | Go install | OS package manager |
| Size | ~2MB | Varies |
| NO_COLOR support | ✅ Yes | ⚠️ Varies |

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
- Uses [fatih/color](https://github.com/fatih/color) for cross-platform colored output
- Inspired by the Unix `tree` command with Git-awareness added
