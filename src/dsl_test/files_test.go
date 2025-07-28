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

func TestFileModule(t *testing.T) {
	if p, err := dsl.LoadFile("sample.dsl"); err != nil {
		t.Log(err)
		t.Fail()
	} else if value, err := p.Module(); err != nil {
		t.Log(err)
		t.Fail()
	} else if value != "history" {
		t.Log("failed to read module")
		t.Fail()
	}
}

func TestFileGroups(t *testing.T) {
	var content dsl.SourceFile
	if p, err := dsl.LoadFile("sample.dsl"); err != nil {
		t.Log(err)
		t.Fail()
	} else {
		content = p
	}

	groups := content.Groups()
	if len(groups) != 4 {
		t.Log("failed to regroup in file")
		t.Log(groups)
		t.Fail()
	}

	topicGroup := groups[0]
	importGroup := groups[1]
	classGroup := groups[2]
	linkGroup := groups[3]

	if len(topicGroup) != 2 {
		t.Log("failed to group topic")
		t.Fail()
	} else if len(importGroup) != 2 {
		t.Log("failed to group import")
		t.Fail()
	} else if len(classGroup) != 8 {
		t.Log("failed to group class")
		t.Fail()
	} else if len(linkGroup) != 6 {
		t.Log("failed to group link")
		t.Fail()
	}

	if topicGroup[0].Value != "topic" {
		t.Log("failed to group topic")
		t.Fail()
	} else if importGroup[0].Value != "import" {
		t.Log("failed to group import")
		t.Fail()
	} else if classGroup[0].Value != "class" {
		t.Log("failed to group class")
		t.Fail()
	} else if linkGroup[0].Value != "link" {
		t.Log("failed to group link")
		t.Fail()
	}
}

func TestDirectoryLoading(t *testing.T) {
	res, errLoad := dsl.LoadAllFilesFromBase("samples/")
	if errLoad != nil {
		t.Log(errLoad)
		t.Fail()
	} else if len(res) != 2 {
		t.Log("missing modules when loading dir")
		t.Fail()
	} else if len(res["history"]) != 1 {
		t.Log("failed to load history")
		t.Fail()
	} else if len(res["humans"]) != 1 {
		t.Log("failed to load humans")
		t.Fail()
	}
}
