package dsl

import (
	"os"
	"path/filepath"
)

// SourceFile is a file that contains source code
type SourceFile struct {
	AbsolutePath string
	Content      []ParsingElement
}

//func WalkDir(root string, fn fs.WalkDirFunc) error

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
