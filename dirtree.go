package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	gitignore "github.com/sabhiram/go-gitignore"
)

var (
	// Color configuration
	useColor     bool
	dirColor     *color.Color
	fileColor    *color.Color
	dotfileColor *color.Color
	execColor    *color.Color
	imageColor   *color.Color
	docColor     *color.Color
)

// main is the entry point of the dirtree application.
// It accepts an optional directory path as a command-line argument.
// If no argument is provided, it defaults to the current directory.
func main() {
	// Define flags
	maxDepth := flag.Int("L", 0, "Maximum depth of directory tree (0 = unlimited)")
	flag.IntVar(maxDepth, "level", 0, "Maximum depth of directory tree (0 = unlimited)")
	colorMode := flag.String("color", "auto", "Colorize output: auto, always, never")
	flag.Parse()

	// Configure colors based on flag and environment
	configureColors(*colorMode)

	// Get directory from remaining args
	args := flag.Args()
	var rootDir string

	if len(args) == 0 {
		rootDir = "."
	} else if len(args) == 1 {
		rootDir = args[0]
		_, err := os.Stat(rootDir)
		if os.IsNotExist(err) {
			fmt.Println("directory does not exist")
			os.Exit(1)
		}
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Usage: dirtree [options] [directory]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println(rootDir)
	printTree(rootDir, "", nil, 0, *maxDepth)
}

// configureColors sets up color output based on the color mode and NO_COLOR environment variable.
func configureColors(mode string) {
	// Check NO_COLOR environment variable (https://no-color.org/)
	if os.Getenv("NO_COLOR") != "" {
		useColor = false
		return
	}

	switch mode {
	case "always":
		useColor = true
		color.NoColor = false
	case "never":
		useColor = false
		color.NoColor = true
	case "auto":
		// Auto mode: use color if stdout is a terminal
		fi, err := os.Stdout.Stat()
		useColor = err == nil && (fi.Mode()&os.ModeCharDevice) != 0
		color.NoColor = !useColor
	default:
		fmt.Printf("Invalid color mode: %s (use auto, always, or never)\n", mode)
		os.Exit(1)
	}

	if useColor {
		// Define colors for different file types
		dirColor = color.New(color.FgBlue, color.Bold)
		fileColor = color.New(color.FgWhite)
		dotfileColor = color.New(color.FgHiBlack)
		execColor = color.New(color.FgGreen, color.Bold)
		imageColor = color.New(color.FgMagenta)
		docColor = color.New(color.FgCyan)
	}
}

// getFileColor returns the appropriate color for a file based on its type and properties.
func getFileColor(file os.DirEntry, fullPath string) *color.Color {
	if !useColor {
		return nil
	}

	// Directories
	if file.IsDir() {
		return dirColor
	}

	name := file.Name()

	// Dotfiles
	if strings.HasPrefix(name, ".") {
		return dotfileColor
	}

	// Check if executable
	info, err := os.Stat(fullPath)
	if err == nil && info.Mode()&0111 != 0 {
		return execColor
	}

	// Check file extension for specific types
	ext := strings.ToLower(filepath.Ext(name))

	// Image files
	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".bmp": true, ".svg": true, ".ico": true, ".webp": true,
	}
	if imageExts[ext] {
		return imageColor
	}

	// Document files
	docExts := map[string]bool{
		".md": true, ".txt": true, ".pdf": true, ".doc": true,
		".docx": true, ".rst": true, ".org": true,
	}
	if docExts[ext] {
		return docColor
	}

	// Default to regular file color
	return fileColor
}

// formatFileName formats a filename with the appropriate color.
func formatFileName(file os.DirEntry, fullPath string) string {
	name := file.Name()
	if !useColor {
		return name
	}

	fileColor := getFileColor(file, fullPath)
	if fileColor != nil {
		return fileColor.Sprint(name)
	}
	return name
}

// shouldIgnore checks if a given path should be ignored based on the provided gitignore matchers.
// It returns true if any of the matchers match the path, false otherwise.
func shouldIgnore(path string, ignoreMatchers []*gitignore.GitIgnore) bool {
	for _, matcher := range ignoreMatchers {
		if matcher.MatchesPath(path) {
			return true
		}
	}
	return false
}

// loadIgnoreMatchers loads .gitignore patterns from the specified directory and combines them
// with parent matchers. This enables hierarchical gitignore inheritance.
// If no .gitignore file exists in the directory, it returns only the parent matchers.
func loadIgnoreMatchers(path string, parentMatchers []*gitignore.GitIgnore) []*gitignore.GitIgnore {
	ignoreFile := filepath.Join(path, ".gitignore")
	if _, err := os.Stat(ignoreFile); os.IsNotExist(err) {
		return parentMatchers
	}
	ignoreMatcher, err := gitignore.CompileIgnoreFile(ignoreFile)
	if err != nil {
		fmt.Printf("Error reading .gitignore in %s: %v\n", path, err)
		return parentMatchers
	}
	return append(parentMatchers, ignoreMatcher)
}

// sortFiles organizes directory entries into a specific order: directories first,
// then regular files, then dot files. Each category is sorted alphabetically.
// The .git directory is automatically filtered out.
func sortFiles(files []os.DirEntry) []os.DirEntry {
	// Create slices for different categories
	var dirs, regularFiles, dotFiles []os.DirEntry

	for _, file := range files {
		if file.Name() == ".git" {
			continue
		}

		if file.IsDir() {
			dirs = append(dirs, file)
		} else if strings.HasPrefix(file.Name(), ".") {
			dotFiles = append(dotFiles, file)
		} else {
			regularFiles = append(regularFiles, file)
		}
	}

	// Sort each category alphabetically
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})
	sort.Slice(regularFiles, func(i, j int) bool {
		return regularFiles[i].Name() < regularFiles[j].Name()
	})
	sort.Slice(dotFiles, func(i, j int) bool {
		return dotFiles[i].Name() < dotFiles[j].Name()
	})

	// Combine in order: directories, regular files, dot files
	result := make([]os.DirEntry, 0, len(dirs)+len(regularFiles)+len(dotFiles))
	result = append(result, dirs...)
	result = append(result, regularFiles...)
	result = append(result, dotFiles...)

	return result
}

// printTree recursively prints the directory tree structure starting from the root directory.
// It respects .gitignore patterns from the current directory and all parent directories.
// The prefix parameter is used to build the ASCII tree structure (├──, └──, │).
// The currentDepth parameter tracks how deep we are in the tree.
// The maxDepth parameter limits recursion depth (0 = unlimited).
func printTree(root string, prefix string, parentMatchers []*gitignore.GitIgnore, currentDepth int, maxDepth int) {
	// Check if we've reached the maximum depth
	if maxDepth > 0 && currentDepth >= maxDepth {
		return
	}
	ignoreMatchers := loadIgnoreMatchers(root, parentMatchers)
	files, err := os.ReadDir(root)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", root, err)
		return
	}

	// Filter out ignored files first
	var filteredFiles []os.DirEntry
	for _, file := range files {
		if file.Name() == ".git" {
			continue
		}
		relPath, _ := filepath.Rel(".", filepath.Join(root, file.Name()))
		if !shouldIgnore(relPath, ignoreMatchers) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	// Sort the filtered files
	sortedFiles := sortFiles(filteredFiles)

	for i, file := range sortedFiles {
		fullPath := filepath.Join(root, file.Name())
		formattedName := formatFileName(file, fullPath)

		if i == len(sortedFiles)-1 {
			fmt.Printf("%s└── %s\n", prefix, formattedName)
		} else {
			fmt.Printf("%s├── %s\n", prefix, formattedName)
		}

		if file.IsDir() {
			newPrefix := prefix
			if i == len(sortedFiles)-1 {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			printTree(fullPath, newPrefix, ignoreMatchers, currentDepth+1, maxDepth)
		}
	}
}
