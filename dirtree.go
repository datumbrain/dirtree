package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

func main() {
	var rootDir string
	switch len(os.Args) {
	case 1:
		rootDir = "."
		fmt.Println(".")
	case 2:
		rootDir = os.Args[1]
		_, err := os.Stat(rootDir)
		if os.IsNotExist(err) {
			fmt.Println("directory does not exist")
			os.Exit(1)
		}
		if err != nil {
			fmt.Println(fmt.Errorf("error: %v", err))
			os.Exit(1)
		}
	default:
		fmt.Println("Usage: dirtree [directory]")
		os.Exit(1)
	}
	printTree(rootDir, "", nil)
}

func shouldIgnore(path string, ignoreMatchers []*gitignore.GitIgnore) bool {
	for _, matcher := range ignoreMatchers {
		if matcher.MatchesPath(path) {
			return true
		}
	}
	return false
}

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

func printTree(root string, prefix string, parentMatchers []*gitignore.GitIgnore) {
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
		if i == len(sortedFiles)-1 {
			fmt.Printf("%s└── %s\n", prefix, file.Name())
		} else {
			fmt.Printf("%s├── %s\n", prefix, file.Name())
		}

		if file.IsDir() {
			newPrefix := prefix
			if i == len(sortedFiles)-1 {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			printTree(filepath.Join(root, file.Name()), newPrefix, ignoreMatchers)
		}
	}
}
