package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestIsDeclaredClasses(t *testing.T) {
	link := commons.NewLink("FRIEND_OF")
	if link == nil {
		t.Errorf("Expected link to be non-nil, got nil")
	} else if !commons.IsElementDeclaredInstance(link, commons.CLASS_LINK) {
		t.Errorf("Expected link to be declared as CLASS_LINK, got undeclared")
	} else if commons.IsElementDeclaredInstance(link, commons.CLASS_TRAIT) {
		t.Errorf("Expected link to not be declared as CLASS_TRAIT, got declared. Bad typing")
	}

	trait := commons.NewTrait("Person")
	if !commons.IsElementDeclaredInstance(trait, commons.CLASS_TRAIT) {
		t.Errorf("Expected trait to be declared as CLASS_TRAIT, got undeclared")
	} else if commons.IsElementDeclaredInstance(trait, commons.CLASS_LINK) {
		t.Errorf("Expected trait to not be declared as CLASS_LINK, got declared. Bad typing")
	}
}
