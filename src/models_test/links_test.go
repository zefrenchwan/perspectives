package models

import (
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
