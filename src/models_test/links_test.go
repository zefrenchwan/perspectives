package models

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
	"github.com/zefrenchwan/perspectives.git/structures"
)

func TestLinksCreation(t *testing.T) {
	john := models.NewObject([]string{"Human"})
	cheese := models.NewTrait("cheese")
	mary := models.NewObject([]string{"Human"})

	if _, err := models.NewLink("likes", map[string]any{models.RoleSubject: john, models.RoleObject: "cheese"}, structures.NewFullPeriod()); err == nil {
		t.Log("failed to detect wrong object operand")
		t.Fail()
	}

	if l, err := models.NewLink("likes", map[string]any{models.RoleSubject: john, models.RoleObject: cheese}, structures.NewFullPeriod()); err != nil {
		t.Log("failed to use object as operand")
		t.Log(err)
		t.Fail()
	} else if l.Name() != "likes" {
		t.Log("wrong name")
		t.Fail()
	}

	if l, err := models.NewLink("likes", map[string]any{models.RoleSubject: []models.Object{john, mary}, models.RoleObject: cheese}, structures.NewFullPeriod()); err != nil {
		t.Log("failed to use group as operand")
		t.Log(err)
		t.Fail()
	} else if l.Name() != "likes" {
		t.Log("wrong name")
		t.Fail()
	} else if k, err := models.NewLink("knows", map[string]any{models.RoleSubject: mary, models.RoleObject: l}, structures.NewFullPeriod()); err != nil {
		t.Log("failed to use link as operand")
		t.Fail()
	} else if k.Name() != "knows" {
		t.Log("wrong composite name")
		t.Fail()
	}
}

func TestLinkWalkthrough(t *testing.T) {
	john := models.NewObject([]string{"Human"})
	cheese := models.NewTrait("cheese")
	mary := models.NewObject([]string{"Human"})

	likes, _ := models.NewLink("likes", map[string]any{models.RoleSubject: []models.Object{john, mary}, models.RoleObject: cheese}, structures.NewFullPeriod())
	knows, _ := models.NewLink("knows", map[string]any{models.RoleSubject: mary, models.RoleObject: likes}, structures.NewFullPeriod())

	// later, for inner link
	var innerLink models.Link
	// First,
	// test main link
	if childs := knows.ValuesPerRole(); len(childs) != 2 {
		t.Log("wrong link roles")
		t.Fail()
	} else if s, found := childs[models.RoleSubject]; !found {
		t.Log("missing subject child at top level")
		t.Fail()
	} else if s.GetType() != models.LinkValueAsObject {
		t.Log("wrong subject type as top level")
		t.Fail()
	} else if so, err := s.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if so.Id != mary.Id {
		t.Log("wront subject value at top level")
		t.Fail()
	} else if o, found := childs[models.RoleObject]; !found {
		t.Log("missing object at top level")
		t.Fail()
	} else if o.GetType() != models.LinkValueAsLink {
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
	} else if s, found := childs[models.RoleSubject]; !found {
		t.Log("missing subject child at inner level")
		t.Fail()
	} else if s.GetType() != models.LinkValueAsGroup {
		t.Log("wrong subject type as inner level")
		t.Fail()
	} else if so, err := s.AsGroup(); err != nil {
		t.Log(err)
		t.Fail()
	} else if len(so) != 2 {
		t.Log("wront subject values at inner level")
		t.Fail()
	} else if o, found := childs[models.RoleObject]; !found {
		t.Log("missing object at inner level")
		t.Fail()
	} else if o.GetType() != models.LinkValueAsTrait {
		t.Log("wrong inner object parameter")
		t.Fail()
	} else if c, err := o.AsTrait(); err != nil {
		t.Log(err)
		t.Fail()
	} else if c.Name != cheese.Name {
		t.Log("wrong object parameter for inner relation")
		t.Fail()
	}
}

func TestLinkObjectsWalkthrough(t *testing.T) {
	john := models.NewObject([]string{"Human"})
	cheese := models.NewObject([]string{"cheese"})
	mary := models.NewObject([]string{"Human"})

	likes, _ := models.NewLink("likes", map[string]any{models.RoleSubject: []models.Object{john, mary}, models.RoleObject: cheese}, structures.NewFullPeriod())
	knows, _ := models.NewLink("knows", map[string]any{models.RoleSubject: mary, models.RoleObject: likes}, structures.NewFullPeriod())

	if values := likes.AllObjectsOperands(); len(values) != 3 {
		t.Log("failed to find objects")
		t.Log(values)
		t.Fail()
	} else if !slices.ContainsFunc(values, func(o models.Object) bool { return o.Id == john.Id }) {
		t.Log("missing object")
		t.Fail()
	} else if !slices.ContainsFunc(values, func(o models.Object) bool { return o.Id == mary.Id }) {
		t.Log("missing other object")
		t.Fail()
	} else if !slices.ContainsFunc(values, func(o models.Object) bool { return o.Id == cheese.Id }) {
		t.Log("missing other object")
		t.Fail()
	}

	if values := knows.AllObjectsOperands(); len(values) != 3 {
		t.Log("failed to find objects")
		t.Log(values)
		t.Fail()
	} else if !slices.ContainsFunc(values, func(o models.Object) bool { return o.Id == john.Id }) {
		t.Log("missing object")
		t.Fail()
	} else if !slices.ContainsFunc(values, func(o models.Object) bool { return o.Id == mary.Id }) {
		t.Log("missing other object")
		t.Fail()
	} else if !slices.ContainsFunc(values, func(o models.Object) bool { return o.Id == cheese.Id }) {
		t.Log("missing other object")
		t.Fail()
	}
}

func TestLinkWalkthroughFunc(t *testing.T) {
	john := models.NewObject([]string{"Human"})
	cheese := models.NewTrait("cheese")
	if likes, err := models.NewSimpleLink("likes", john, cheese); err != nil {
		t.Log("failed to create simple link")
		t.Log(err)
		t.Fail()
	} else if knows, err := models.NewSimpleLink("knows", john, likes); err != nil {
		t.Log("failed to create parent link")
		t.Log(err)
		t.Fail()
	} else if matches := knows.FindAllMatchingCondition(func(lv models.LinkValue) bool { return lv.GetType() == models.LinkValueAsLink }); len(matches) != 2 {
		t.Log("failed to find links in relation")
		t.Log(matches)
		t.Fail()
	} else {
		var names []string
		for _, value := range matches {
			if value.GetType() != models.LinkValueAsLink {
				t.Log("failed to read link")
				t.Fail()
			} else if l, err := value.AsLink(); err != nil {
				t.Log("failed to read link")
				t.Fail()
			} else {
				names = append(names, l.Name())
			}
		}

		slices.Sort(names)
		if slices.Compare(names, []string{"knows", "likes"}) != 0 {
			t.Log("wrong content")
			t.Fail()
		}
	}

	if likes, err := models.NewSimpleLink("likes", john, cheese); err != nil {
		t.Log("failed to create simple link")
		t.Log(err)
		t.Fail()
	} else if knows, err := models.NewSimpleLink("knows", john, likes); err != nil {
		t.Log("failed to create parent link")
		t.Log(err)
		t.Fail()
	} else if matches := knows.FindAllMatchingCondition(func(lv models.LinkValue) bool { return lv.GetType() == models.LinkValueAsTrait }); len(matches) != 1 {
		t.Log("failed to find traits in relation")
		t.Log(matches)
		t.Fail()
	} else if val := matches[0]; val.GetType() != models.LinkValueAsTrait {
		t.Log("wrong type")
		t.Fail()
	} else if trait, err := val.AsTrait(); err != nil {
		t.Log(err)
		t.Fail()
	} else if !trait.Equals(cheese) {
		t.Log("wrong read trait")
		t.Fail()
	}
}
