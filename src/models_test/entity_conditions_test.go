package models_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
)

func TestAcceptTypesCondition(t *testing.T) {
	jane := models.NewObject([]string{"Human"})
	desk := models.NewObject([]string{"Desk"})
	link, _ := models.NewSimpleLink("works on", jane, desk)

	values := []models.EntityType{models.EntityTypeObject}

	condition := models.ConditionOnEntityType{MatchingTypes: values}
	if condition.Triggers(models.NewParameter(link)) {
		t.Fail()
	} else if !condition.Triggers(models.NewParameter(jane)) {
		t.Fail()
	} else if condition.Triggers(nil) {
		t.Fail()
	}

	// test for nil array
	condition = models.ConditionOnEntityType{MatchingTypes: nil}
	if condition.Triggers(models.NewParameter(link)) {
		t.Fail()
	} else if condition.Triggers(models.NewParameter(jane)) {
		t.Fail()
	} else if condition.Triggers(nil) {
		t.Fail()
	}

}
