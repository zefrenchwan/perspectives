package commons_test

import (
	"errors"
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

func (d DummyObject) AsGroup() (commons.ModelGroup, error) {
	return nil, errors.ErrUnsupported
}

func (d DummyObject) GetType() commons.ModelableType {
	return commons.TypeObject
}

func (d DummyObject) AsObject() (commons.ModelObject, error) {
	return d, nil
}

func (d DummyObject) IsGroup() bool {
	return false
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
