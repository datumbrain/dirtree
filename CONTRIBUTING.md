# Contributing to dirtree

Thank you for considering contributing to dirtree! We welcome contributions from the community.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior vs actual behavior
- Your environment (OS, Go version, dirtree version)
- Any relevant logs or output

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- A clear and descriptive title
- A detailed description of the proposed functionality
- Example use cases
- Why this enhancement would be useful

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following the code style guidelines below
3. **Add tests** for any new functionality
4. **Ensure all tests pass** by running `go test ./...`
5. **Update documentation** if needed
6. **Submit a pull request** with a clear description of your changes

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git

### Setting Up Your Development Environment

1. Fork and clone the repository:
   ```bash
   git clone https://github.com/YOUR_USERNAME/dirtree.git
   cd dirtree
   ```

2. Download dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   go build
   ```

4. Run tests:
   ```bash
   go test -v ./...
   ```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run tests with verbose output
go test -v ./...
```

### Code Style Guidelines

- Follow standard Go conventions and idioms
- Run `go fmt` on your code before committing
- Run `go vet` to catch common mistakes
- Add comments to exported functions following godoc conventions
- Keep functions focused and single-purpose
- Write clear, descriptive variable names

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests after the first line

Example:
```
Add depth limiting functionality

Implements a -depth flag to control maximum recursion depth.
This addresses feature request #123.
```

## Project Structure

```
dirtree/
├── dirtree.go       # Main implementation
├── dirtree_test.go  # Test suite
├── go.mod           # Go module definition
├── go.sum           # Dependency checksums
└── README.md        # User documentation
```

## Testing Philosophy

- Write tests for all new functionality
- Maintain or improve code coverage
- Include both unit tests and integration tests
- Test edge cases and error conditions

## Release Process

Releases are managed by the maintainers. Version numbers follow [Semantic Versioning](https://semver.org/).

## Questions?

Feel free to open an issue for questions or discussion before starting work on significant changes.

## License

By contributing to dirtree, you agree that your contributions will be licensed under the MIT License.
