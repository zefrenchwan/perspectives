package models_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/models"
)

func TestIdBasedCondition(t *testing.T) {
	helene := models.NewObject([]string{"Human"})
	id := helene.Id()
	daniel := models.NewObject([]string{"Human"})
	condition := commons.IdBasedCondition{Id: id}

	p := commons.NewNamedParameter("x", helene)
	p.AppendAsVariable("y", daniel)

	if condition.Matches(p) {
		t.Log("multiple values should not match")
		t.Fail()
	}

	p = commons.NewNamedParameter("y", daniel)
	if condition.Matches(p) {
		t.Log("bad id matching")
		t.Fail()
	}

	p = commons.NewParameter(helene)
	if !condition.Matches(p) {
		t.Log("id should match")
		t.Fail()
	}

	// makes no sense, but prove that a non identifiable object would be rejected
	variable := models.NewVariableForTrait(id)
	p = commons.NewNamedParameter("x", variable)
	if condition.Matches(p) {
		t.Fail()
	}

}
