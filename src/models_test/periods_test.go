package models_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
)

func TestSerializeEmpty(t *testing.T) {
	empty := models.NewEmptyPeriod()
	serVal := empty.AsString()
	if serVal != "][" {
		t.Log("failed to serialize empty interval")
		t.Fail()
	}

	if result, err := models.PeriodFromString(serVal); err != nil {
		t.Log("Failed to deserialize empty", err)
		t.Fail()
	} else if result != empty {
		t.Log("invalid values for empty")
		t.Fail()
	}
}

func TestSerializeFull(t *testing.T) {
	full := models.NewFullPeriod()
	serVal := full.AsString()
	if serVal != "]-oo;+oo[" {
		t.Log("failed to serialize full interval")
		t.Fail()
	}

	if result, err := models.PeriodFromString(serVal); err != nil {
		t.Log("Failed to deserialize full", err)
		t.Fail()
	} else if result != full {
		t.Log("invalid values for full")
		t.Fail()
	}
}
