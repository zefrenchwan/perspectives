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
