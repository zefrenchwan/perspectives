package values_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/values"
)

func TestReferences(t *testing.T) {
	ref := values.NewReference("id")
	s := values.NewString("id")
	if ref.Datatype() != values.REFERENCE_TYPE {
		t.Errorf("Expected REFERENCE_TYPE, got %s", ref.Datatype())
	} else if ref.Equals(s) {
		t.Errorf("Expected false because a reference is not primitive, got true")
	} else if ref.Content() != s.Content() {
		t.Errorf("same content, should be true")
	}
}
