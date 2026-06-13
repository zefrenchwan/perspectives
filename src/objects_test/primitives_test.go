package objects_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/objects"
)

func TestIsPrimitiveValue(t *testing.T) {
	// Valid primitives
	if !objects.IsPrimitiveValue(42) {
		t.Error("expected int to be a valid primitive value")
	} else if !objects.IsPrimitiveValue(3.14) {
		t.Error("expected float64 to be a valid primitive value")
	} else if !objects.IsPrimitiveValue("hello") {
		t.Error("expected string to be a valid primitive value")
	} else if !objects.IsPrimitiveValue(true) {
		t.Error("expected bool to be a valid primitive value")
	} else if !objects.IsPrimitiveValue(time.Now()) {
		t.Error("expected time.Time to be a valid primitive value")
	}

	// Invalid primitives
	if objects.IsPrimitiveValue([]int{1, 2}) {
		t.Error("expected slice NOT to be a primitive value")
	} else if objects.IsPrimitiveValue(map[string]int{"a": 1}) {
		t.Error("expected map NOT to be a primitive value")
	} else if objects.IsPrimitiveValue(nil) {
		t.Error("expected nil NOT to be a primitive value")
	} else if objects.IsPrimitiveValue(struct{ name string }{name: "test"}) {
		t.Error("expected struct NOT to be a primitive value")
	}
}
