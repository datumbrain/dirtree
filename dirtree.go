package main

import (
	"fmt"
	"os"
	"path/filepath"

	gitignore "github.com/sabhiram/go-gitignore"
)

func main() {
	rootDir := "."

	ignoreFile := filepath.Join(rootDir, ".gitignore")
	ignoreMatcher, err := gitignore.CompileIgnoreFile(ignoreFile)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error reading .gitignore: %v\n", err)
		return
	}

	fmt.Println(".")
	printTree(rootDir, "", ignoreMatcher)
}

func printTree(root string, prefix string, ignoreMatcher *gitignore.GitIgnore) {
	files, err := os.ReadDir(root)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", root, err)
		return
	}

	for i, file := range files {
		relPath, _ := filepath.Rel(".", filepath.Join(root, file.Name()))
		if ignoreMatcher != nil && ignoreMatcher.MatchesPath(relPath) {
			continue
		}

		if i == len(files)-1 {
			fmt.Printf("%s└── %s\n", prefix, file.Name())
		} else {
			fmt.Printf("%s├── %s\n", prefix, file.Name())
		}

		if file.IsDir() {
			newPrefix := prefix
			if i == len(files)-1 {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			printTree(filepath.Join(root, file.Name()), newPrefix, ignoreMatcher)
		}
	}
}
