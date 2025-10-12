package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestTemporalOperatorEquals(t *testing.T) {
	reference := commons.NewFullPeriod()
	other := commons.NewPeriodSince(time.Now(), true)

	if commons.TemporalEquals.Accepts(other, reference) {
		t.Fail()
	} else if !commons.TemporalEquals.Accepts(reference, reference) {
		t.Fail()
	}
}

func TestTemporalOperatorConstants(t *testing.T) {
	reference := commons.NewFullPeriod()
	other := commons.NewPeriodSince(time.Now(), true)

	if commons.TemporalAlwaysRefuse.Accepts(reference, reference) {
		t.Fail()
	} else if !commons.TemporalAlwaysAccept.Accepts(other, reference) {
		t.Fail()
	}
}

func TestTemporalOperatorCommonPoint(t *testing.T) {
	all := commons.NewFullPeriod()
	before := time.Now().AddDate(-1, 0, 0)
	after := time.Now().AddDate(1, 0, 0)
	base := commons.NewPeriodSince(after, true)
	other := commons.NewPeriodUntil(before, true)

	if commons.TemporalCommonPoint.Accepts(base, other) {
		t.Fail()
	} else if !commons.TemporalCommonPoint.Accepts(all, base) {
		t.Fail()
	}
}

func TestTemporalOperatorInclusion(t *testing.T) {
	all := commons.NewFullPeriod()
	before := time.Now().AddDate(-1, 0, 0)
	after := time.Now().AddDate(1, 0, 0)
	base := commons.NewPeriodSince(after, true)
	other := commons.NewPeriodUntil(before, true)

	if commons.TemporalReferenceContains.Accepts(base, other) {
		t.Fail()
	} else if !commons.TemporalCommonPoint.Accepts(base, all) {
		t.Fail()
	}
}
