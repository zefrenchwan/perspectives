package dsl

import (
	"io/fs"
	"os"
	"path/filepath"
)

// FindFilesFrom looks for files matching a predicate since an entry point.
// It returns all the paths matching that predicate (or an error if any).
// If startingPath is a directory, it goes through that directory.
// Otherwise, for a file, it tests the file and just returns it (if matching).
// Nil predicate is considered to match everything.
func FindFilesFrom(startingPath string, predicate func(string) bool) ([]string, error) {
	// startingPath may be a directory (usual case) or not (then test it)
	if info, err := os.Stat(startingPath); err != nil {
		return nil, err
	} else if !info.IsDir() {
		if predicate == nil || predicate(startingPath) {
			return []string{startingPath}, nil
		}
	}

	var elements []string
	addAll := func(currentPath string, currentDir fs.DirEntry, err error) error {
		if stat, errStat := os.Stat(currentPath); err != nil {
			return errStat
		} else if !stat.IsDir() {
			if predicate == nil || predicate(currentPath) {
				elements = append(elements, currentPath)
			}
		}
		return err
	}

	if err := filepath.WalkDir(startingPath, addAll); err != nil {
		return nil, err
	} else {
		return elements, nil
	}
}

// LoadFile returns the content of the file as a string, or error if any
func LoadFile(path string) (string, error) {
	if content, errRead := os.ReadFile(path); errRead != nil {
		return "", errRead
	} else {
		return string(content), nil
	}
}
