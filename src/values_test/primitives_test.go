package values_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/values"
)

func TestIsPrimitiveValue(t *testing.T) {
	// Valid primitives
	if !values.IsPrimitiveValue(42) {
		t.Error("expected int to be a valid primitive value")
	} else if !values.IsPrimitiveValue(3.14) {
		t.Error("expected float64 to be a valid primitive value")
	} else if !values.IsPrimitiveValue("hello") {
		t.Error("expected string to be a valid primitive value")
	} else if !values.IsPrimitiveValue(true) {
		t.Error("expected bool to be a valid primitive value")
	} else if !values.IsPrimitiveValue(time.Now()) {
		t.Error("expected time.Time to be a valid primitive value")
	}

	// Invalid primitives
	if values.IsPrimitiveValue([]int{1, 2}) {
		t.Error("expected slice NOT to be a primitive value")
	} else if values.IsPrimitiveValue(map[string]int{"a": 1}) {
		t.Error("expected map NOT to be a primitive value")
	} else if values.IsPrimitiveValue(nil) {
		t.Error("expected nil NOT to be a primitive value")
	} else if values.IsPrimitiveValue(struct{ name string }{name: "test"}) {
		t.Error("expected struct NOT to be a primitive value")
	}
}

func TestIsPrimitiveTypeName(t *testing.T) {
	if !values.IsPrimitiveTypeName("int") {
		t.Error("expected int to be a valid primitive type name")
	} else if !values.IsPrimitiveTypeName("float64") {
		t.Error("expected float64 to be a valid primitive type name")
	} else if !values.IsPrimitiveTypeName("string") {
		t.Error("expected string to be a valid primitive type name")
	} else if !values.IsPrimitiveTypeName("bool") {
		t.Error("expected bool to be a valid primitive type name")
	} else if !values.IsPrimitiveTypeName("time.Time") {
		t.Error("expected time.Time to be a valid primitive type name")
	}

	if values.IsPrimitiveTypeName("[]int") {
		t.Error("expected slice NOT to be a primitive type name")
	} else if values.IsPrimitiveTypeName("map[string]int") {
		t.Error("expected map[string]int NOT to be a primitive type name")
	} else if values.IsPrimitiveTypeName("nil") {
		t.Error("expected nil NOT to be a primitive type name")
	} else if values.IsPrimitiveTypeName("struct{ name string }") {
		t.Error("expected other struct NOT to be a primitive type name")
	}
}

func TestPrimitiveInt(t *testing.T) {
	a := values.NewInt(0)
	b := values.NewInt(1)
	s := values.NewString("test")

	if a.Datatype() != values.PRIMITIVE_TYPE_INT {
		t.Error("expected datatype to be int")
	} else if s.Datatype() != values.PRIMITIVE_TYPE_STRING {
		t.Error("expected datatype to be string")
	} else if a.Equals(b) {
		t.Error("expected int values NOT to be equal")
	} else if a.Equals(s) {
		t.Error("expected int and string values NOT to be equal")
	} else if !a.Equals(values.NewInt(0)) {
		t.Error("expected int values to be equal")
	} else if b.Content() != 1 {
		t.Errorf("Content should  be 1, got %v", b.Content())
	}
}

func TestPrimitiveString(t *testing.T) {
	a := values.NewString("a")
	b := values.NewString("b")
	s := values.NewString("a")

	if a.Datatype() != values.PRIMITIVE_TYPE_STRING {
		t.Error("expected datatype to be string")
	} else if a.Equals(b) {
		t.Error("expected string values NOT to be equal")
	} else if !a.Equals(s) {
		t.Error("expected string values to be equal")
	} else if a.Content() != "a" {
		t.Error("same content, should be true")
	}
}

func TestPrimitiveTime(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().Truncate(time.Second).In(loc)
	standardNow := now.UTC()
	if !standardNow.Equal(now) {
		t.Error("Wrong assumption when building code : timezones differ but time is the same")
	}

	value := values.NewTime(now)
	other := values.NewTime(standardNow)
	if !value.Equals(other) {
		t.Error("expected time values to be equal despite timezone")
	}

	value = values.NewTime(now.UTC())
	other = values.NewTime(standardNow.UTC())
	if !value.Equals(other) {
		t.Error("expected time values to be equals (same value, same timezone)")
	}

	value = values.NewTime(now.Add(time.Hour))
	other = values.NewTime(standardNow)
	if value.Equals(other) {
		t.Error("expected time values to be different (values differ)")
	}

	if value.Datatype() != values.PRIMITIVE_TYPE_TIME {
		t.Error("expected datatype to be time")
	}
}
