package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestStateObject(t *testing.T) {
	obj := commons.NewStateObject[int]()

	if len(obj.Attributes()) != 0 {
		t.Fail()
	} else if _, found := obj.GetValue("test"); found {
		t.Fail()
	}

	obj.SetValue("key", 10)
	if value, found := obj.GetValue("key"); !found {
		t.Fail()
	} else if value != 10 {
		t.Fail()
	} else if attr := obj.Attributes(); len(attr) != 1 {
		t.Fail()
	} else if attr[0] != "key" {
		t.Fail()
	}
}
