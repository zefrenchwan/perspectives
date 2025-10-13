package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

type DummyComponentImplementation struct {
}

type DummyIdBasedImplementation struct {
	id string
}

func (d DummyIdBasedImplementation) Id() string {
	return d.id
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
