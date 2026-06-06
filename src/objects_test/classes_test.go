package objects_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/objects"
)

func TestIsDeclaredClasses(t *testing.T) {
	trait := objects.NewTrait("Person")
	if trait == nil {
		t.Errorf("Expected trait to be non-nil, got nil")
	} else if !objects.IsElementDeclaredInstance(trait, objects.CLASS_TRAIT) {
		t.Errorf("Expected trait to be declared as CLASS_TRAIT, got undeclared")
	} else if objects.IsElementDeclaredInstance(trait, objects.CLASS_LINK) {
		t.Errorf("Expected trait to not be declared as CLASS_LINK, got declared. Bad typing")
	}
}
