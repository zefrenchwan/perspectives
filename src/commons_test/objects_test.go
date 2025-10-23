package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestStateObjects(t *testing.T) {
	obj := commons.NewStateObject[int]()
	if obj.GetType() != commons.TypeObject {
		t.Fail()
	} else if !obj.ActivePeriod().Equals(commons.NewFullPeriod()) {
		t.Fail()
	}

	p := commons.NewPeriodSince(time.Now().Truncate(time.Second), true)
	obj.SetActivePeriod(p)
	if !obj.ActivePeriod().Equals(p) {
		t.Fail()
	}

	obj.SetValue("attr", 10)
	obj.SetValue("other", 30)

	if v, found := obj.GetValue("attr"); !found {
		t.Fail()
	} else if v != 10 {
		t.Fail()
	} else if _, found := obj.GetValue("not exist"); found {
		t.Fail()
	} else if obj.Remove("not here") {
		t.Fail()
	} else if !obj.Remove("other") {
		t.Fail()
	} else if r := obj.Read(); r == nil {
		t.Fail()
	} else if values := r.Values(); len(values) != 1 {
		t.Fail()
	} else if values["attr"] != 10 {
		t.Fail()
	}
}

func TestTemporalStateObject(t *testing.T) {
	activity := commons.NewPeriodSince(time.Now().Truncate(time.Second), true)
	obj := commons.NewTemporalStateObject[int](activity)
	if !obj.ActivePeriod().Equals(activity) {
		t.Fail()
	} else if obj.GetType() != commons.TypeObject {
		t.Fail()
	}

	before := time.Now().AddDate(-10, 0, 0)
	beforePeriod := commons.NewPeriodSince(before, true)
	obj.SetValueDuringPeriod("attr", 10, beforePeriod)

	if r := obj.ReadAtTime(time.Now().AddDate(10, 0, 0)); r == nil {
		t.Fail()
	} else if values := r.Values(); len(values) != 1 {
		t.Fail()
	} else if values["attr"] != 10 {
		t.Fail()
	} else if r := obj.ReadAtTime(before); r != nil {
		t.Fail()
	} else if r := obj.Read(); !r.ActivePeriod().Equals(activity) {
		t.Fail()
	} else if attrs := r.Values(); len(attrs) != 1 {
		t.Fail()
	} else if values := attrs["attr"]; len(values) != 1 {
		t.Fail()
	} else if !values[10].Equals(beforePeriod) {
		t.Fail()
	}
}
