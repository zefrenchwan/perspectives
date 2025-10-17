package commons_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestGroupCreation(t *testing.T) {
	basic := DummyObject{}
	group := commons.NewModelGroup([]commons.ModelEntity{basic})

	if !commons.IsGroup(group) {
		t.Fail()
	} else if group.GetType() != commons.TypeGroup {
		t.Fail()
	} else if !slices.Equal([]commons.ModelEntity{basic}, group.Content()) {
		t.Fail()
	} else if values := slices.Collect(group.Elements()); len(values) != 1 {
		t.Fail()
	} else if values[0] != basic {
		t.Fail()
	}
}
