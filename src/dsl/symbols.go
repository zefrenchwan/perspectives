package dsl

import (
	"slices"
	"strings"
)

// COMMENT is the string that opens a comment.
// Rest of the line is ignored
const COMMENT = "##"

// ESCAPE_CHAR defines the escape character
const ESCAPE_CHAR = "\\"

// /////////
// MARKS //
// /////////
const MARK_POINT = "."
const MARK_TRIPLE_POINTS = "..."
const MARK_EXCLAMATION = "!"
const MARK_COMMA = ","
const MARK_QUESTION = "?"
const MARK_SEMICOLON = ";"
const MARK_PARENTHESIS_LEFT = "("
const MARK_PARENTHESIS_RIGHT = ")"
const MARK_DOUBLE_QUOTES = `"`
const MARK_SINGLE_QUOTES = "'"

// MARK_SYMBOLS contains all symbols to consider as marks
var MARK_SYMBOLS = []string{
	MARK_POINT, MARK_EXCLAMATION, MARK_COMMA,
	MARK_QUESTION, MARK_SEMICOLON, MARK_TRIPLE_POINTS,
	MARK_DOUBLE_QUOTES, MARK_SINGLE_QUOTES,
	MARK_PARENTHESIS_LEFT, MARK_PARENTHESIS_RIGHT,
}

// IsMarkSymbol tests if all the values in the string are marks (and there is at least one char)
func IsMarkSymbol(value string) bool {
	for _, val := range strings.Split(value, "") {
		if !slices.Contains(MARK_SYMBOLS, val) {
			return false
		}
	}

	return len(value) != 0
}

// //////////////////////////////////
// KEYWORDS FOR SOURCE MANAGEMENT //
// //////////////////////////////////

// module declaration
const KW_MODULE = "topic"

// import content from another module
const KW_IMPORT = "import"
