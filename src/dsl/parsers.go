package dsl

import (
	"io"
	"regexp"
	"strings"
)

// ParsingElement is a string that contains something to parse.
// Note that it is a consecutive suite of non empty chars.
// It may contain multiple tokens such as end; ("end" + ";")
type ParsingElement struct {
	Value    string // Value of the element
	Line     int    // line number in the original source (starts at 0)
	Position int    // position in line (starts at 0)
}

// Load reads a content and isolates all the parsed elements
func Load(reader io.Reader) ([]ParsingElement, error) {
	var elements []ParsingElement

	var payload string
	if p, err := io.ReadAll(reader); err != nil {
		return nil, err
	} else {
		payload = string(p)
	}

	spaceTester := regexp.MustCompile(`\A\s+\z`)
	var lineNumber int
	var positionNumber int
	for rawLine := range strings.Lines(payload) {
		line := strings.Split(rawLine, COMMENT)[0]
		var space bool
		var previousSpace bool
		var buffer []string
		positionNumber = 0
		for value := range strings.SplitSeq(line, "") {
			space = spaceTester.MatchString(value)
			switch {
			case !space:
				buffer = append(buffer, value)
			case space && !previousSpace:
				if len(buffer) != 0 {
					content := strings.Join(buffer, "")
					size := len(content)
					elements = append(elements, ParsingElement{Value: content, Line: lineNumber, Position: positionNumber - size})
					buffer = nil
				}
			}

			previousSpace = space
			positionNumber = positionNumber + 1
		}

		if !space {
			if len(buffer) != 0 {
				content := strings.Join(buffer, "")
				size := len(content)
				elements = append(elements, ParsingElement{Value: content, Line: lineNumber, Position: positionNumber - size})
				buffer = nil
			}
		}

		lineNumber = lineNumber + 1
	}

	return elements, nil
}
