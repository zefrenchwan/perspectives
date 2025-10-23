package commons_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestDescribeStateFromStateObject(t *testing.T) {
	desc := commons.NewRequestDescription[string]("x")

	if sign := desc.Signature(); slices.Compare([]string{"x"}, sign.Variables()) != 0 {
		t.Fail()
	}

	obj := commons.NewStateObject[string]()
	obj.SetValue("attr", "test")

	content := commons.NewNamedContent("y", obj)
	if status := desc.Describe(content); status != nil {
		t.Log("no variable selection")
		t.Fail()
	}

	content = commons.NewNamedContent("x", obj)
	if status := desc.Describe(content); status == nil {
		t.Log("no variable selection")
		t.Fail()
	} else if status.Values()["attr"] != "test" {
		t.Fail()
	}

}

func TestDescribeStateFromTemporalObject(t *testing.T) {
	desc := commons.NewRequestDescription[string]("x")
	obj := commons.NewTemporalStateObject[string](commons.NewFullPeriod())
	obj.SetValueDuringPeriod("attr", "test", commons.NewFullPeriod())

	content := commons.NewNamedContent("x", obj)
	if status := desc.Describe(content); status == nil {
		t.Log("no variable selection")
		t.Fail()
	} else if status.Values()["attr"] != "test" {
		t.Fail()
	}
}

func TestDescribeStateFromNonReadable(t *testing.T) {
	desc := commons.NewRequestDescription[string]("x")
	obj := DummyIdBasedImplementation{}

	content := commons.NewNamedContent("x", obj)
	if status := desc.Describe(content); status != nil {
		t.Log("not a state reader")
		t.Fail()
	}
}

func TestTemporalDescribeFromTemporalObject(t *testing.T) {
	desc := commons.NewRequestTemporalDescription[string]("x")
	obj := commons.NewTemporalStateObject[string](commons.NewFullPeriod())
	obj.SetValueDuringPeriod("attr", "test", commons.NewFullPeriod())

	content := commons.NewNamedContent("x", obj)
	if status := desc.Describe(content); status == nil {
		t.Log("no variable selection")
		t.Fail()
	} else if p := status.ActivePeriod(); !p.Equals(obj.ActivePeriod()) {
		t.Fail()
	} else if values := status.Values(); len(values) != 1 {
		t.Fail()
	} else if value := values["attr"]; len(value) != 1 {
		t.Fail()
	} else if !value["test"].Equals(commons.NewFullPeriod()) {
		t.Fail()
	}
}
