package dsl

import (
	"strings"
	"unicode"
)

// Element defines a token (unclassified)
type Element struct {
	// LineIndex of that token (lines start at 0)
	LineIndex int
	// RowIndex of that token (rows start at 0)
	RowIndex int
	// Value of the token
	Value string
}

// Read maps a content (of a file) to all matching elements
func Read(content string) []Element {
	var elements []Element
	for lineNumber, line := range strings.Split(content, "\n") {
		remaining := line
		indexComment := -1
		if index := strings.Index(line, KW_COMMENT); index >= 0 {
			remaining = remaining[0:index]
			indexComment = index
		}

		currentBuffer := ""
		isSpaceBefore := true
		for index, value := range remaining {
			if unicode.IsSpace(value) {
				if !isSpaceBefore {
					startingPoint := index - len(currentBuffer)
					currentValue := Element{LineIndex: lineNumber, RowIndex: startingPoint, Value: currentBuffer}
					elements = append(elements, currentValue)
					currentBuffer = ""
				}

				isSpaceBefore = true
			} else {
				currentBuffer = currentBuffer + string(value)
				isSpaceBefore = false
			}
		}

		if len(currentBuffer) > 0 {
			startingPoint := len(line) - len(currentBuffer)
			if indexComment >= 0 {
				startingPoint = indexComment - len(currentBuffer)
			}
			currentValue := Element{LineIndex: lineNumber, RowIndex: startingPoint, Value: currentBuffer}
			elements = append(elements, currentValue)
		}

		if indexComment >= 0 {
			commentElement := Element{LineIndex: lineNumber, RowIndex: indexComment, Value: "//"}
			elements = append(elements, commentElement)
		}

	}

	return elements
}
