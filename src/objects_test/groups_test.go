package objects_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/objects"
	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestGroupsSets(t *testing.T) {
	instance, _ := objects.
		NewLocalInstanceBuilder("1").
		WithActivity(periods.NewFullPeriod()).
		Build()

	other, _ := objects.
		NewLocalInstanceBuilder("2").
		WithActivity(periods.NewFullPeriod()).
		Build()

	setTwo := objects.NewSetOfInstances([]objects.Instance{instance, other})
	setOne := objects.NewSetOfInstances([]objects.Instance{instance})
	setZero := objects.NewSetOfInstances([]objects.Instance{})

	if setZero.Size() != 0 {
		t.Errorf("Sets should be equal")
	} else if setOne.Size() != 1 {
		t.Errorf("Sets should be equal")
	} else if setTwo.Size() != 2 {
		t.Errorf("Sets should be equal")
	}

	if len(setZero.SortedInstances()) != 0 {
		t.Errorf("Sorted instances should be equal for 0 element")
	} else if setOneValues := setOne.SortedInstances(); len(setOneValues) != 1 {
		t.Errorf("Sorted instances should be equal for 1 element")
	} else if !setOneValues[0].Same(instance) {
		t.Errorf("Sorted values should be equal")
	} else if setTwoValues := setTwo.SortedInstances(); len(setTwoValues) != 2 {
		t.Errorf("Sorted instances should be equal for 2 element")
	} else if !setTwoValues[0].Same(instance) {
		t.Errorf("Sorted values should be equal due to order")
	} else if !setTwoValues[1].Same(other) {
		t.Errorf("Sorted values should be equal due to order")
	}

	for element := range setOne.Range {
		if !element.Same(instance) {
			t.Errorf("Element should be equal to instance")
		}
	}

	if !setOne.Contains(instance) {
		t.Errorf("Set should contain instance")
	} else if setOne.Contains(other) {
		t.Errorf("Set should not contain other")
	} else if setZero.Contains(instance) {
		t.Errorf("Set should not contain instance")
	} else if setZero.Contains(other) {
		t.Errorf("Set should not contain other")
	}
}

func TestGroupsSetsSame(t *testing.T) {
	instance, _ := objects.
		NewLocalInstanceBuilder("1").
		WithActivity(periods.NewFullPeriod()).
		Build()

	other, _ := objects.
		NewLocalInstanceBuilder("2").
		WithActivity(periods.NewFullPeriod()).
		Build()

	setTwo := objects.NewSetOfInstances([]objects.Instance{instance, other})
	setTwoSame := objects.NewSetOfInstances([]objects.Instance{other, instance})

	setOne := objects.NewSetOfInstances([]objects.Instance{instance})

	setZero := objects.NewSetOfInstances([]objects.Instance{})

	if !setTwo.Same(setTwoSame) {
		t.Errorf("Sets should be equal")
	} else if !setOne.Same(setOne) {
		t.Errorf("Sets should be equal")
	} else if !setZero.Same(setZero) {
		t.Errorf("Sets should be equal")
	}

	if setOne.Same(setTwo) {
		t.Errorf("Sets should not be equal")
	} else if setTwo.Same(setOne) {
		t.Errorf("Sets should not be equal")
	} else if setOne.Same(setZero) {
		t.Errorf("Sets should not be equal")
	} else if setZero.Same(setOne) {
		t.Errorf("Sets should not be equal")
	}
}
