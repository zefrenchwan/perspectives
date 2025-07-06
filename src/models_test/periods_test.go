package models_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/models"
)

func TestIntersectionWithFull(t *testing.T) {
	value := models.NewFullPeriod()
	now := time.Now()
	otherValue := models.NewPeriodSince(now, true)
	result := otherValue.Intersection(value)
	if !result.Equals(otherValue) {
		t.Log("intersection with full failed")
		t.Fail()
	}
}

func TestIntersectionWithEmpty(t *testing.T) {
	value := models.NewEmptyPeriod()
	now := time.Now()
	otherValue := models.NewPeriodSince(now, true)
	result := otherValue.Intersection(value)
	if !result.Equals(value) {
		t.Log("intersection with empty failed")
		t.Fail()
	}
}

func TestIntersectionWithOther(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := time.Now().Add(-24 * time.Hour)
	value := models.NewPeriodSince(before, true)
	otherValue := models.NewPeriodUntil(now, true)
	expected := models.NewFinitePeriod(before, now, true, true)
	result := otherValue.Intersection(value)
	if !result.Equals(expected) {
		t.Logf("intersection with other failed: got %s", result.AsRawString())
		t.Fail()
	}
}
