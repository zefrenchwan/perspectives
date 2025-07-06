package dsl

import (
	"io"
	"regexp"
	"strings"
)

const COMMENT = "##"

type ParsingElement struct {
	Value    string
	Line     int
	Position int
}

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
		for _, value := range strings.Split(line, "") {
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
