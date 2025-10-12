package models_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
)

func TestIdBasedCondition(t *testing.T) {
	helene := models.NewObject([]string{"Human"})
	id := helene.Id()
	daniel := models.NewObject([]string{"Human"})
	condition := models.IdBasedCondition{Id: id}

	p := models.NewNamedParameter("x", helene)
	p.AppendAsVariable("y", daniel)

	if condition.Matches(p) {
		t.Log("multiple values should not match")
		t.Fail()
	}

	p = models.NewNamedParameter("y", daniel)
	if condition.Matches(p) {
		t.Log("bad id matching")
		t.Fail()
	}

	p = models.NewParameter(helene)
	if !condition.Matches(p) {
		t.Log("id should match")
		t.Fail()
	}

}
