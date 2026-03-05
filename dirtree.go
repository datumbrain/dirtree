package main

import (
	"encoding/json"
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

// TreeNode represents a single node in the directory tree.
type TreeNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"is_dir"`
	Children []*TreeNode `json:"children,omitempty"`
}

// TreeStats holds statistics about the tree.
type TreeStats struct {
	FileCount int `json:"files"`
	DirCount  int `json:"dirs"`
}

// TreeData holds the complete tree structure and metadata.
type TreeData struct {
	RootPath string      `json:"root"`
	Depth    int         `json:"depth"`
	Stats    TreeStats   `json:"stats"`
	Nodes    []*TreeNode `json:"nodes"`
}

// main is the entry point of the dirtree application.
// It accepts an optional directory path as a command-line argument.
// If no argument is provided, it defaults to the current directory.
func main() {
	// Define flags
	maxDepth := flag.Int("L", 0, "Maximum depth of directory tree (0 = unlimited)")
	flag.IntVar(maxDepth, "level", 0, "Maximum depth of directory tree (0 = unlimited)")
	colorMode := flag.String("color", "auto", "Colorize output: auto, always, never")
	format := flag.String("format", "text", "Output format: text, json, markdown (or md)")
	flag.StringVar(format, "o", "text", "Output format: text, json, markdown (or md)")
	maxFiles := flag.Int("max-files", 0, "Maximum number of files to display (0 = unlimited)")
	llmMode := flag.Bool("llm", false, "LLM-optimized output (equivalent to: --format markdown --level 3 --max-files 200)")
	flag.Parse()

	// Handle --llm flag (overrides other settings)
	if *llmMode {
		*format = "markdown"
		if *maxDepth == 0 {
			*maxDepth = 3
		}
		if *maxFiles == 0 {
			*maxFiles = 200
		}
	}

	// Normalize format aliases
	if *format == "md" {
		*format = "markdown"
	}

	// Validate format
	validFormats := map[string]bool{"text": true, "json": true, "markdown": true}
	if !validFormats[*format] {
		fmt.Printf("Invalid format: %s (use text, json, or markdown)\n", *format)
		os.Exit(1)
	}

	// Configure colors only for text format
	if *format == "text" {
		configureColors(*colorMode)
	}

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

	// Build tree structure
	fileCount := 0
	treeData := buildTree(rootDir, nil, 0, *maxDepth, *maxFiles, &fileCount)

	// Format output based on selected format
	switch *format {
	case "text":
		fmt.Println(rootDir)
		formatText(treeData.Nodes, "", rootDir)
	case "json":
		formatJSON(treeData, rootDir, *maxDepth)
	case "markdown":
		formatMarkdown(treeData, rootDir, *maxDepth)
	}
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

// buildTree recursively builds a tree structure from the directory hierarchy.
// It respects .gitignore patterns and applies depth and file count limits.
// Returns a TreeData structure containing the tree nodes and statistics.
func buildTree(root string, parentMatchers []*gitignore.GitIgnore, currentDepth int, maxDepth int, maxFiles int, fileCount *int) *TreeData {
	// Check if we've reached the maximum depth
	if maxDepth > 0 && currentDepth >= maxDepth {
		return &TreeData{Nodes: []*TreeNode{}}
	}

	// Check if we've reached the maximum file count
	if maxFiles > 0 && *fileCount >= maxFiles {
		return &TreeData{Nodes: []*TreeNode{}}
	}

	ignoreMatchers := loadIgnoreMatchers(root, parentMatchers)
	files, err := os.ReadDir(root)
	if err != nil {
		return &TreeData{Nodes: []*TreeNode{}}
	}

	// Filter out ignored files
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

	// Build tree nodes
	var nodes []*TreeNode
	stats := TreeStats{}

	for _, file := range sortedFiles {
		// Check file count limit
		if maxFiles > 0 && *fileCount >= maxFiles {
			break
		}

		fullPath := filepath.Join(root, file.Name())
		relPath, _ := filepath.Rel(".", fullPath)

		node := &TreeNode{
			Name:  file.Name(),
			Path:  relPath,
			IsDir: file.IsDir(),
		}

		if file.IsDir() {
			stats.DirCount++
			// Recursively build children
			childData := buildTree(fullPath, ignoreMatchers, currentDepth+1, maxDepth, maxFiles, fileCount)
			node.Children = childData.Nodes
			stats.FileCount += childData.Stats.FileCount
			stats.DirCount += childData.Stats.DirCount
		} else {
			stats.FileCount++
			*fileCount++
		}

		nodes = append(nodes, node)
	}

	return &TreeData{
		Nodes: nodes,
		Stats: stats,
	}
}

// formatText outputs the tree in the traditional text format with ASCII art.
// This maintains backward compatibility with the original printTree output.
func formatText(nodes []*TreeNode, prefix string, rootDir string) {
	for i, node := range nodes {
		fullPath := filepath.Join(rootDir, node.Path)

		// Get file info for color formatting
		info, _ := os.Stat(fullPath)
		var mockEntry dirEntryWrapper
		if info != nil {
			mockEntry = dirEntryWrapper{name: node.Name, isDir: node.IsDir, info: info}
		} else {
			mockEntry = dirEntryWrapper{name: node.Name, isDir: node.IsDir}
		}

		formattedName := formatFileName(mockEntry, fullPath)

		if i == len(nodes)-1 {
			fmt.Printf("%s└── %s\n", prefix, formattedName)
		} else {
			fmt.Printf("%s├── %s\n", prefix, formattedName)
		}

		if node.IsDir && len(node.Children) > 0 {
			newPrefix := prefix
			if i == len(nodes)-1 {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			formatText(node.Children, newPrefix, rootDir)
		}
	}
}

// dirEntryWrapper wraps file info to implement os.DirEntry interface for formatFileName.
type dirEntryWrapper struct {
	name  string
	isDir bool
	info  os.FileInfo
}

func (d dirEntryWrapper) Name() string               { return d.name }
func (d dirEntryWrapper) IsDir() bool                { return d.isDir }
func (d dirEntryWrapper) Type() os.FileMode          { return 0 }
func (d dirEntryWrapper) Info() (os.FileInfo, error) { return d.info, nil }

// formatJSON outputs the tree in JSON format.
func formatJSON(treeData *TreeData, rootDir string, maxDepth int) {
	output := map[string]interface{}{
		"root":  rootDir,
		"depth": maxDepth,
		"stats": treeData.Stats,
		"nodes": flattenNodes(treeData.Nodes),
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

// flattenNodes converts the tree structure to a flat list for JSON output.
func flattenNodes(nodes []*TreeNode) []map[string]interface{} {
	var result []map[string]interface{}

	var flatten func([]*TreeNode)
	flatten = func(nodes []*TreeNode) {
		for _, node := range nodes {
			nodeType := "file"
			if node.IsDir {
				nodeType = "dir"
			}
			result = append(result, map[string]interface{}{
				"path": node.Path,
				"type": nodeType,
			})
			if node.IsDir && len(node.Children) > 0 {
				flatten(node.Children)
			}
		}
	}

	flatten(nodes)
	return result
}

// formatMarkdown outputs the tree in markdown format optimized for LLMs.
func formatMarkdown(treeData *TreeData, rootDir string, maxDepth int) {
	fmt.Printf("# Directory Structure: %s\n\n", rootDir)

	// Group files by directory
	filesByDir := make(map[string][]string)
	var dirs []string

	var collectFiles func([]*TreeNode, string)
	collectFiles = func(nodes []*TreeNode, parentPath string) {
		for _, node := range nodes {
			if node.IsDir {
				dirs = append(dirs, node.Path)
				if len(node.Children) > 0 {
					collectFiles(node.Children, node.Path)
				}
			} else {
				dir := filepath.Dir(node.Path)
				if dir == "." {
					dir = "Root Directory"
				}
				filesByDir[dir] = append(filesByDir[dir], node.Name)
			}
		}
	}

	collectFiles(treeData.Nodes, "")

	// Output files section
	if len(filesByDir) > 0 {
		fmt.Println("## Files")

		// Sort directories for consistent output
		var sortedDirs []string
		for dir := range filesByDir {
			sortedDirs = append(sortedDirs, dir)
		}
		sort.Strings(sortedDirs)

		for _, dir := range sortedDirs {
			files := filesByDir[dir]
			if dir == "Root Directory" {
				fmt.Println("### Root Directory")
			} else {
				fmt.Printf("### %s/\n", dir)
			}
			for _, file := range files {
				fmt.Printf("- `%s`\n", file)
			}
			fmt.Println()
		}
	}

	// Output directories section
	if len(dirs) > 0 {
		fmt.Println("## Directories")
		sort.Strings(dirs)
		for _, dir := range dirs {
			fmt.Printf("- `%s/`\n", dir)
		}
		fmt.Println()
	}

	// Output stats footer
	fmt.Println("---")
	fmt.Printf("**Total:** %d files, %d directories\n", treeData.Stats.FileCount, treeData.Stats.DirCount)
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
