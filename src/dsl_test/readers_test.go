package dsl_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/dsl"
)

func TestReadContent(t *testing.T) {
	value := "  I    am      sure"
	expected := []dsl.Element{
		{LineIndex: 0, RowIndex: 2, Value: "I"},
		{LineIndex: 0, RowIndex: 7, Value: "am"},
		{LineIndex: 0, RowIndex: 15, Value: "sure"},
	}
	elements := dsl.Read(value)
	if len(elements) != len(expected) {
		t.Log(elements)
		t.Fail()
	}

	for i := range elements {
		if elements[i] != expected[i] {
			t.Fail()
		}
	}
}
func TestReadContentWithComment(t *testing.T) {
	value := "I am sure // a comment"
	expected := []dsl.Element{
		{LineIndex: 0, RowIndex: 0, Value: "I"},
		{LineIndex: 0, RowIndex: 2, Value: "am"},
		{LineIndex: 0, RowIndex: 5, Value: "sure"},
		{LineIndex: 0, RowIndex: 10, Value: "//"},
	}
	elements := dsl.Read(value)
	if len(elements) != len(expected) {
		t.Log(elements)
		t.Fail()
	}

	for i := range elements {
		if elements[i] != expected[i] {
			t.Fail()
		}
	}
}
