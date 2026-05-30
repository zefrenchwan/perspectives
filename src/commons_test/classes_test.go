package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestIsDeclaredClasses(t *testing.T) {
	trait := commons.NewTrait("Person")
	if trait == nil {
		t.Errorf("Expected trait to be non-nil, got nil")
	} else if !commons.IsElementDeclaredInstance(trait, commons.CLASS_TRAIT) {
		t.Errorf("Expected trait to be declared as CLASS_TRAIT, got undeclared")
	} else if commons.IsElementDeclaredInstance(trait, commons.CLASS_LINK) {
		t.Errorf("Expected trait to not be declared as CLASS_LINK, got declared. Bad typing")
	}
}
