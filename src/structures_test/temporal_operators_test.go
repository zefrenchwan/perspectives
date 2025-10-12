package structures_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/structures"
)

func TestTemporalOperatorEquals(t *testing.T) {
	reference := structures.NewFullPeriod()
	other := structures.NewPeriodSince(time.Now(), true)

	if structures.TemporalEquals.Accepts(other, reference) {
		t.Fail()
	} else if !structures.TemporalEquals.Accepts(reference, reference) {
		t.Fail()
	}
}

func TestTemporalOperatorConstants(t *testing.T) {
	reference := structures.NewFullPeriod()
	other := structures.NewPeriodSince(time.Now(), true)

	if structures.TemporalAlwaysRefuse.Accepts(reference, reference) {
		t.Fail()
	} else if !structures.TemporalAlwaysAccept.Accepts(other, reference) {
		t.Fail()
	}
}

func TestTemporalOperatorCommonPoint(t *testing.T) {
	all := structures.NewFullPeriod()
	before := time.Now().AddDate(-1, 0, 0)
	after := time.Now().AddDate(1, 0, 0)
	base := structures.NewPeriodSince(after, true)
	other := structures.NewPeriodUntil(before, true)

	if structures.TemporalCommonPoint.Accepts(base, other) {
		t.Fail()
	} else if !structures.TemporalCommonPoint.Accepts(all, base) {
		t.Fail()
	}
}

func TestTemporalOperatorInclusion(t *testing.T) {
	all := structures.NewFullPeriod()
	before := time.Now().AddDate(-1, 0, 0)
	after := time.Now().AddDate(1, 0, 0)
	base := structures.NewPeriodSince(after, true)
	other := structures.NewPeriodUntil(before, true)

	if structures.TemporalReferenceContains.Accepts(base, other) {
		t.Fail()
	} else if !structures.TemporalCommonPoint.Accepts(base, all) {
		t.Fail()
	}
}
