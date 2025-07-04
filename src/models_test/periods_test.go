package models_test

import (
	"testing"
	"time"

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
	} else if !models.PeriodEquals(result, empty) {
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
	} else if !models.PeriodEquals(result, full) {
		t.Log("invalid values for full")
		t.Fail()
	}
}

func TestSerializeLeftFinite(t *testing.T) {
	now := time.Now()
	interval := models.NewPeriodSince(now, true)
	intervalValue := interval.AsString()
	if copyInterval, err := models.PeriodFromString(intervalValue); err != nil {
		t.Logf("failed to read %s", intervalValue)
		t.Fail()
	} else if !models.PeriodEquals(copyInterval, interval) {
		t.Log("Values differ")
		t.Log(interval.PeriodRawValue())
		t.Log(copyInterval.PeriodRawValue())
		t.Fail()
	}
}

func TestSerializeRightFinite(t *testing.T) {
	now := time.Now()
	interval := models.NewPeriodUntil(now, true)
	intervalValue := interval.AsString()
	if copyInterval, err := models.PeriodFromString(intervalValue); err != nil {
		t.Logf("failed to read %s", intervalValue)
		t.Fail()
	} else if !models.PeriodEquals(copyInterval, interval) {
		t.Log("Values differ")
		t.Log(interval.PeriodRawValue())
		t.Log(copyInterval.PeriodRawValue())
		t.Fail()
	}
}

func TestSerializeFinite(t *testing.T) {
	now := time.Now()
	before := now.Add(-1 * time.Hour)
	interval := models.NewFinitePeriod(before, now, true, false)
	intervalValue := interval.AsString()
	if copyInterval, err := models.PeriodFromString(intervalValue); err != nil {
		t.Logf("failed to read %s", intervalValue)
		t.Fail()
	} else if !models.PeriodEquals(copyInterval, interval) {
		t.Log("Values differ")
		t.Log(interval.PeriodRawValue())
		t.Log(copyInterval.PeriodRawValue())
		t.Fail()
	}
}
