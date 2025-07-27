package dsl_test

import (
	"strings"
	"testing"

	"github.com/zefrenchwan/perspectives.git/dsl"
)

func TestFileLoading(t *testing.T) {
	if p, err := dsl.LoadFile("sample.dsl"); err != nil {
		t.Log(err)
		t.Fail()
	} else if len(p.Content) == 0 {
		t.Log("could not read content")
		t.Fail()
	} else if !strings.HasSuffix(p.AbsolutePath, "sample.dsl") {
		t.Log("invalid file path")
		t.Fail()
	}
}
