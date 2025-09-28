package models

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
	"github.com/zefrenchwan/perspectives.git/structures"
)

func TestTypeCondition(t *testing.T) {
	object := models.NewObject([]string{"Human"})
	link, _ := models.NewQualifier(object, "good", structures.NewFullPeriod())

	condition := models.NewTypeCondition(models.EntityTypeLink)
	if condition.Matches(object) {
		t.Log("should accept link only")
		t.Fail()
	} else if !condition.Matches(link) {
		t.Log("should accept link")
		t.Fail()
	}
}
