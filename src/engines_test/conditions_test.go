package models

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/engines"
	"github.com/zefrenchwan/perspectives.git/models"
	"github.com/zefrenchwan/perspectives.git/structures"
)

func TestTypeCondition(t *testing.T) {
	object := models.NewObject([]string{"Human"})
	link, _ := models.NewQualifier(object, "good", structures.NewFullPeriod())

	condition := engines.NewTypeCondition(models.EntityTypeLink)
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
	condition := engines.LocalMatchingAttributeCondition{
		AttributeName:     "gender",
		AttributeValue:    "M",
		AttributeOperator: engines.ValuesEqual,
		ReferencePeriod:   structures.NewFullPeriod(),
		PeriodOperator:    engines.AcceptsAllPeriods,
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
	condition = engines.LocalMatchingAttributeCondition{
		AttributeName:     "gender",
		AttributeValue:    "F",
		AttributeOperator: engines.ValuesEqual,
		ReferencePeriod:   structures.NewPeriodSince(date, true),
		PeriodOperator:    engines.NonDisjoinPeriods,
	}

	mary := models.NewObjectSince([]string{"Human"}, birthdate)
	mary.SetValue("gender", "F")
	if !condition.Matches(mary) {
		t.Fail()
	}
}

func TestObjectAttributeRegexpCondition(t *testing.T) {
	objectValueMatch := models.NewObject([]string{"Human"})
	objectValueMatch.SetValue("attr", "abc")

	if condition, err := engines.NewAttributeRegexpCondition("attr", "\\w"); err != nil {
		t.Log("valid regexp should not fail")
		t.Fail()
	} else if !condition.Matches(objectValueMatch) {
		t.Log("regexp accepting any word should match attribute")
		t.Fail()
	}

	objectValueMatch.SetValue("attr", "a")
	if condition, err := engines.NewAttributeRegexpCondition("attr", "\\d+"); err != nil {
		t.Log("valid regexp should not fail")
		t.Fail()
	} else if condition.Matches(objectValueMatch) {
		t.Log("regexp does not match value")
		t.Fail()
	}
}

func TestActiveConditionForTemporalEntities(t *testing.T) {
	birthdate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now().Truncate(time.Second)
	after := now.AddDate(10, 0, 0)
	valid := models.NewObjectSince([]string{"Human"}, birthdate)
	condition := engines.NewActiveCondition(structures.NewPeriodSince(now, true))
	invalid := models.NewObjectSince([]string{"Human"}, after)

	if !condition.Matches(valid) {
		t.Fail()
	} else if condition.Matches(invalid) {
		t.Fail()
	}

	// and of course
	object := models.NewObject([]string{"Human"})
	if !condition.Matches(object) {
		t.Fail()
	}

	// test forever condition
	condition = engines.NewActiveCondition(structures.NewFullPeriod())
	if !condition.Matches(object) {
		t.Fail()
	}

	// apply to link
	other := models.NewObject([]string{"Human"})
	likes, _ := models.NewTimedSimpleLink("likes", structures.NewPeriodSince(now, true), object, other)
	if condition.Matches(likes) {
		t.Log("likes cannot match because was false before now")
		t.Fail()
	}

	// test when equals
	condition = engines.NewActiveCondition(structures.NewPeriodSince(now, true))
	if !condition.Matches(likes) {
		t.Fail()
	}
}

func TestCombineConditions(t *testing.T) {
	object := models.NewObject([]string{"Human"})
	object.SetValue("name", "Doe")

	// testing NOT condition
	condition := engines.NewActiveCondition(structures.NewPeriodSince(time.Now(), true))
	if !condition.Matches(object) {
		t.Fail()
	} else if notCondition := engines.NotCondition(condition); notCondition.Matches(object) {
		t.Fail()
	}

	otherCondition := engines.NewAttributeValueCondition("name", "Danny", engines.ValuesEqualIgnoreCase)
	if otherCondition.Matches(object) {
		t.Fail()
	} else if !engines.NotCondition(otherCondition).Matches(object) {
		t.Fail()
	}

	// condition should accept and otherCondition should not
	// So test OR and AND.
	// And should be false because one is false,
	// Or should be true because one is true
	orCondition := engines.OrConditions([]engines.LocalCondition{otherCondition, condition})
	andCondition := engines.AndConditions([]engines.LocalCondition{otherCondition, condition})

	if andCondition.Matches(object) {
		t.Fail()
	} else if !orCondition.Matches(object) {
		t.Fail()
	}
}

func TestLinkValueCondition(t *testing.T) {
	homer := models.NewObject([]string{"Human"})
	donuts := models.NewObject([]string{"Food"})
	likes, _ := models.NewSimpleLink("likes", homer, donuts)

	// test when condition matches
	condition := engines.NewLinkNameInValuesCondition([]string{"likes", "loves", "adores"}, engines.ValuesEqualIgnoreCase)
	if !condition.Matches(likes) {
		t.Fail()
	}

	// test for no match
	condition = engines.NewLinkNameInValuesCondition([]string{"hates", "do not like", "ignores"}, engines.ValuesEqualIgnoreCase)
	if condition.Matches(likes) {
		t.Fail()
	}

}
