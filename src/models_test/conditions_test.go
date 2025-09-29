package models

import (
	"testing"
	"time"

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

func TestObjectAttributeCondition(t *testing.T) {
	// gender = M no matter the period
	condition := models.LocalMatchingAttributeCondition{
		AttributeName:     "gender",
		AttributeValue:    "M",
		AttributeOperator: models.ValuesEqual,
		ReferencePeriod:   structures.NewFullPeriod(),
		PeriodOoperator:   models.AcceptsAllOperator,
	}

	objectNoMatch := models.NewObject([]string{"Human"})
	objectNoMatch.SetValue("no match", "popo")
	if condition.Matches(objectNoMatch) {
		t.Fail()
	}

	objectValueMismatch := models.NewObject([]string{"Human"})
	objectValueMismatch.SetValue("gender", "F")
	if condition.Matches(objectValueMismatch) {
		t.Fail()
	}

	objectValueMatch := models.NewObject([]string{"Human"})
	objectValueMatch.SetValue("gender", "M")
	if !condition.Matches(objectValueMatch) {
		t.Fail()
	}

	date := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	birthdate := time.Now().Add(-24 * time.Hour).Truncate(time.Minute)
	// gender = F since a date
	condition = models.LocalMatchingAttributeCondition{
		AttributeName:     "gender",
		AttributeValue:    "F",
		AttributeOperator: models.ValuesEqual,
		ReferencePeriod:   structures.NewPeriodSince(date, true),
		PeriodOoperator:   models.NonDisjoinPeriods,
	}

	mary := models.NewObjectSince([]string{"Human"}, birthdate)
	mary.SetValue("gender", "F")
	if !condition.Matches(mary) {
		t.Fail()
	}

}
