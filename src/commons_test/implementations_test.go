package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

const DummyTestingType commons.ModelableType = -1

type DummyComponentImplementation struct {
}

func (d DummyComponentImplementation) GetType() commons.ModelableType {
	return DummyTestingType
}

type DummyIdBasedImplementation struct {
	id string
}

func (d DummyIdBasedImplementation) GetType() commons.ModelableType {
	return DummyTestingType
}

func (d DummyIdBasedImplementation) Id() string {
	return d.id
}

type DummyObject struct {
	id string
}

func (d DummyObject) Id() string {
	return d.id
}

func (d DummyObject) GetType() commons.ModelableType {
	return commons.TypeObject
}

func TestDummyIdBasedImplementation(t *testing.T) {
	// create dummy value, but forces it to any
	value := any(DummyIdBasedImplementation{id: "id"})
	// confirms it implements IdentifiableElement
	if v, ok := value.(commons.Identifiable); !ok {
		t.Fail()
	} else if v.Id() != "id" {
		t.Fail()
	}
}

func TestDummyObjectImplementation(t *testing.T) {
	value := any(DummyObject{id: "test"})
	if v, ok := value.(commons.Identifiable); !ok {
		t.Fail()
	} else if v.Id() != "test" {
		t.Fail()
	} else if v, ok := value.(commons.ModelObject); !ok {
		t.Fail()
	} else if v.GetType() != commons.TypeObject {
		t.Fail()
	} else if v.Id() != "test" {
		t.Fail()
	}
}
