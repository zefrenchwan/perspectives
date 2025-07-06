package dsl_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/zefrenchwan/perspectives.git/dsl"
)

func TestElementsReader(t *testing.T) {
	const content = "Token tests ## comment \n    Value"

	expectedTokens := []dsl.ParsingElement{
		{Value: "Token", Line: 0, Position: 0},
		{Value: "tests", Line: 0, Position: 6},
		{Value: "Value", Line: 1, Position: 4},
	}

	if tokens, err := dsl.Load(strings.NewReader(content)); err != nil {
		t.Log(err)
		t.Fail()
	} else if len(tokens) != 3 {
		t.Log(tokens)
		t.Log("Failed to have all tokens")
		t.Fail()
	} else if !slices.Equal(tokens, expectedTokens) {
		t.Log("Tokens failure")
		t.Log(tokens)
		t.Fail()
	}
}
