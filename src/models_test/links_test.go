package models

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
	"github.com/zefrenchwan/perspectives.git/structures"
)

func TestLinksCreation(t *testing.T) {
	john := models.NewObject([]string{"Human"})
	cheese := models.NewObject([]string{"cheese"})
	mary := models.NewObject([]string{"Human"})

	if _, err := models.NewLink("likes", map[string]any{"subject": john, "object": "cheese"}, structures.NewFullPeriod()); err == nil {
		t.Log("failed to detect wrong object operand")
		t.Fail()
	}

	if l, err := models.NewLink("likes", map[string]any{"subject": john, "object": cheese}, structures.NewFullPeriod()); err != nil {
		t.Log("failed to use object as operand")
		t.Log(err)
		t.Fail()
	} else if l.Name() != "likes" {
		t.Log("wrong name")
		t.Fail()
	}

	if l, err := models.NewLink("likes", map[string]any{"subject": []models.Object{john, mary}, "object": cheese}, structures.NewFullPeriod()); err != nil {
		t.Log("failed to use group as operand")
		t.Log(err)
		t.Fail()
	} else if l.Name() != "likes" {
		t.Log("wrong name")
		t.Fail()
	} else if k, err := models.NewLink("knows", map[string]any{"subject": mary, "object": l}, structures.NewFullPeriod()); err != nil {
		t.Log("failed to use link as operand")
		t.Fail()
	} else if k.Name() != "knows" {
		t.Log("wrong composite name")
		t.Fail()
	}
}

func TestLinkWalkthrough(t *testing.T) {
	john := models.NewObject([]string{"Human"})
	cheese := models.NewObject([]string{"cheese"})
	mary := models.NewObject([]string{"Human"})

	likes, _ := models.NewLink("likes", map[string]any{"subject": []models.Object{john, mary}, "object": cheese}, structures.NewFullPeriod())
	knows, _ := models.NewLink("knows", map[string]any{"subject": mary, "object": likes}, structures.NewFullPeriod())

	// later, for inner link
	var innerLink models.Link
	// First,
	// test main link
	if childs := knows.ValuesPerRole(); len(childs) != 2 {
		t.Log("wrong link roles")
		t.Fail()
	} else if s, found := childs["subject"]; !found {
		t.Log("missing subject child at top level")
		t.Fail()
	} else if !s.IsObject() {
		t.Log("wrong subject type as top level")
		t.Fail()
	} else if so, err := s.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if so.Id != mary.Id {
		t.Log("wront subject value at top level")
		t.Fail()
	} else if o, found := childs["object"]; !found {
		t.Log("missing object at top level")
		t.Fail()
	} else if !o.IsLink() {
		t.Log("wrong link parameter")
		t.Fail()
	} else if inner, err := o.AsLink(); err != nil {
		t.Log(err)
		t.Fail()
	} else {
		innerLink = inner
	}

	// test inner link
	if childs := innerLink.ValuesPerRole(); len(childs) != 2 {
		t.Log("wrong link roles")
		t.Fail()
	} else if s, found := childs["subject"]; !found {
		t.Log("missing subject child at inner level")
		t.Fail()
	} else if !s.IsGroup() {
		t.Log("wrong subject type as inner level")
		t.Fail()
	} else if so, err := s.AsGroup(); err != nil {
		t.Log(err)
		t.Fail()
	} else if len(so) != 2 {
		t.Log("wront subject values at inner level")
		t.Fail()
	} else if o, found := childs["object"]; !found {
		t.Log("missing object at inner level")
		t.Fail()
	} else if !o.IsObject() {
		t.Log("wrong inner object parameter")
		t.Fail()
	} else if c, err := o.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if c.Id != cheese.Id {
		t.Log("wrong object parameter for inner relation")
		t.Fail()
	}
}
