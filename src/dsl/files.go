package dsl

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// SourceFile is a file that contains source code
type SourceFile struct {
	AbsolutePath string
	Content      []ParsingElement
}

// LoadFile loads a content from a file
func LoadFile(contentPath string) (SourceFile, error) {
	var result SourceFile
	var file *os.File
	if f, err := os.Open(contentPath); err != nil {
		return result, err
	} else if f != nil {
		defer f.Close()
		file = f
	}

	if absPath, err := filepath.Abs(contentPath); err != nil {
		return result, err
	} else if parsed, err := Load(file); err != nil {
		return result, err
	} else {
		result = SourceFile{AbsolutePath: absPath, Content: parsed}
		return result, nil
	}
}

// LoadAllFilesFromBase loads either a file or a directory and regroups all content into modules.
// acceptFile returns true if file should be read. It applies only on regular files.
// Result is then a map of module and linked source files
func LoadAllFilesFromBase(sourceBase string, acceptFile func(path string) bool) (map[string][]SourceFile, error) {
	result := make(map[string][]SourceFile)
	if res, err := os.Stat(sourceBase); err != nil {
		return nil, err
	} else if !res.IsDir() {
		if !acceptFile(sourceBase) {
			return nil, nil
		} else if content, err := LoadFile(sourceBase); err != nil {
			return nil, err
		} else if m, err := content.Module(); err != nil {
			return nil, err
		} else {
			result[m] = []SourceFile{content}
			return result, nil
		}
	}

	errWalk := filepath.WalkDir(sourceBase, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || !acceptFile(path) {
			return nil
		} else if content, err := LoadFile(path); err != nil {
			return err
		} else if m, err := content.Module(); err != nil {
			return err
		} else {
			existingValue := result[m]
			result[m] = append(existingValue, content)
			return nil
		}
	})

	return result, errWalk
}

// Module returns the module of the source file, or error for a malformed file
func (s SourceFile) Module() (string, error) {
	if len(s.Content) < 2 {
		return "", errors.New("no module declaration")
	}

	moduleUnit, nameUnit := s.Content[0], s.Content[1]
	line := moduleUnit.Line
	if moduleUnit.Value != KW_MODULE {
		return "", fmt.Errorf("expecting module at position %d (line %d)", moduleUnit.Position, line)
	} else if nameUnit.Line != line {
		return "", fmt.Errorf("expecting module name at line %d", line)
	}

	return nameUnit.Value, nil
}

// Groups returns the parsed file, but as groups of elements
func (s SourceFile) Groups() [][]ParsingElement {
	if len(s.Content) == 0 {
		return nil
	}

	var currentGroup []ParsingElement
	var result [][]ParsingElement
	for _, value := range s.Content {
		currentPosition := value.Position
		if currentPosition == 0 {
			// store previous group if any
			if currentGroup != nil {
				result = append(result, currentGroup)
			}

			currentGroup = make([]ParsingElement, 0)
		}

		currentGroup = append(currentGroup, value)
	}

	if len(currentGroup) != 0 {
		result = append(result, currentGroup)
	}

	return result
}
