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
func TestCloneSimpleLink(t *testing.T) {
	sonia := models.NewObject([]string{"Human"})
	jack := models.NewObject([]string{"Human"})
	married, _ := models.NewSimpleLink("married", sonia, jack)

	// test the run ended
	clone := married.CopyStructure()
	if clone == nil {
		t.Fail()
	}

	if clone.Name() != married.Name() {
		t.Log("bad name for link")
		t.Fail()
	} else if clone.Id() != married.Id() {
		t.Log("copy should keep id")
		t.Fail()
	} else if !clone.Duration().Equals(married.Duration()) {
		t.Log("copy should keep duration")
		t.Fail()
	}

	operands := clone.Operands()
	if len(operands) != 2 {
		t.Log("missing operands")
		t.Fail()
	} else if s, found := operands[models.RoleSubject]; !found {
		t.Log("missing subject")
		t.Fail()
	} else if subject, err := s.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if subject.Id != sonia.Id {
		t.Log("subject not mapped")
		t.Fail()
	} else if o, found := operands[models.RoleObject]; !found {
		t.Log("object not found")
		t.Fail()
	} else if object, err := o.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if object.Id != jack.Id {
		t.Log("wrong object")
		t.Fail()
	}
}

func TestCloneLongLink(t *testing.T) {
	sonia := models.NewObject([]string{"Human"})
	jack := models.NewObject([]string{"Human"})
	married, _ := models.NewSimpleLink("married", sonia, jack)
	marcel := models.NewObject([]string{"Human"})
	ignores, _ := models.NewSimpleLink("ignores", marcel, married)
	knows, _ := models.NewSimpleLink("knows", jack, ignores)
	// Jack knows that Marcel ignores that Sonia married Jack

	clone := knows.CopyStructure()
	if clone == nil {
		t.Log("clone failed")
		t.Fail()
	} else if clone.Name() != knows.Name() {
		t.Log("knows name failed")
		t.Fail()
	} else if clone.Id() != knows.Id() {
		t.Log("knows id failed")
		t.Fail()
	}

	knowsOperands := clone.Operands()
	if len(knowsOperands) != 2 {
		t.Log("root operands failed")
	} else if s, found := knowsOperands[models.RoleSubject]; !found {
		t.Log("root subject failed")
		t.Fail()
	} else if subject, err := s.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if subject.Id != jack.Id {
		t.Log("root subject id failed")
		t.Fail()
	} else if o, found := knowsOperands[models.RoleObject]; !found {
		t.Log("failed to find object in root")
		t.Fail()
	} else if o.GetType() != models.EntityTypeLink {
		t.Log("object should be a link")
		t.Fail()
	}

	ignoresLink, _ := knowsOperands[models.RoleObject].AsLink()
	if ignoresLink.Id() != ignores.Id() {
		t.Log("failed first level id")
		t.Fail()
	} else if ignoresLink.Name() != ignores.Name() {
		t.Log("failed first level name")
		t.Fail()
	} else if operands := ignoresLink.Operands(); len(operands) != 2 {
		t.Log("wrong first level operand")
		t.Fail()
	} else if s, found := operands[models.RoleSubject]; !found {
		t.Log("missing subject at first level")
		t.Fail()
	} else if subject, err := s.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if subject.Id != marcel.Id {
		t.Log("wrong subject at first level")
		t.Fail()
	}

	marriedLink, _ := ignoresLink.Operands()[models.RoleObject].AsLink()
	if marriedLink.Id() != married.Id() {
		t.Log("failed second level id")
		t.Fail()
	} else if marriedLink.Name() != married.Name() {
		t.Log("failed second level name")
		t.Fail()
	} else if operands := marriedLink.Operands(); len(operands) != 2 {
		t.Log("wrong second level operands")
		t.Fail()
	} else if s, found := operands[models.RoleSubject]; !found {
		t.Log("missing subject at level 2")
		t.Fail()
	} else if subject, err := s.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if subject.Id != sonia.Id {
		t.Log("wrong subject at level 2")
		t.Fail()
	} else if o, found := operands[models.RoleObject]; !found {
		t.Log("missing object at level 2")
		t.Fail()
	} else if object, err := o.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if object.Id != jack.Id {
		t.Log("wrong object at level 2")
		t.Fail()
	}

}
