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

func TestGetPrimitiveType(t *testing.T) {
	now := time.Now()
	if _, has := values.GetPrimitiveType(nil); has {
		t.Error("expected nil NOT to be a valid primitive type")
	} else if s, has := values.GetPrimitiveType("popo"); !has {
		t.Error("expected string to be a valid primitive type")
	} else if s != values.PRIMITIVE_TYPE_STRING {
		t.Errorf("expected string to be a valid primitive type, got %v", s)
	} else if b, has := values.GetPrimitiveType(true); !has {
		t.Error("expected bool to be a valid primitive type")
	} else if b != values.PRIMITIVE_TYPE_BOOL {
		t.Errorf("expected bool to be a valid primitive type, got %v", b)
	} else if i, has := values.GetPrimitiveType(1); !has {
		t.Error("expected int to be a valid primitive type")
	} else if i != values.PRIMITIVE_TYPE_INT {
		t.Errorf("expected int to be a valid primitive type, got %v", i)
	} else if f, has := values.GetPrimitiveType(1.0); !has {
		t.Error("expected float to be a valid primitive type")
	} else if f != values.PRIMITIVE_TYPE_FLOAT {
		t.Errorf("expected float to be a valid primitive type, got %v", f)
	} else if tt, has := values.GetPrimitiveType(now); !has {
		t.Error("expected time to be a valid primitive type")
	} else if tt != values.PRIMITIVE_TYPE_TIME {
		t.Errorf("expected time to be a valid primitive type, got %v", tt)
	}
}

func TestPrimitiveBool(t *testing.T) {
	b := values.NewBool(true)
	if b.Datatype() != values.PRIMITIVE_TYPE_BOOL {
		t.Error("expected datatype to be bool")
	} else if b.Content() != true {
		t.Errorf("Content should  be true, got %v", b.Content())
	}

	a := values.NewBool(false)
	if a.Equals(b) {
		t.Error("expected bool values NOT to be equal")
	} else if !a.Equals(a) {
		t.Error("expected bool values to be equal")
	} else if a.ToHashString() == b.ToHashString() {
		t.Error("expected different bool values NOT to have the same hash")
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
	} else if a.ToHashString() == b.ToHashString() {
		t.Error("expected different int values NOT to have the same hash")
	} else if a.Equals(s) {
		t.Error("expected int and string values NOT to be equal")
	} else if !a.Equals(values.NewInt(0)) {
		t.Error("expected int values to be equal")
	} else if b.Content() != 1 {
		t.Errorf("Content should  be 1, got %v", b.Content())
	}
}

func TestPrimitiveFloat(t *testing.T) {
	f := values.NewFloat(1.0)
	if f.Datatype() != values.PRIMITIVE_TYPE_FLOAT {
		t.Error("expected datatype to be float")
	} else if f.Content() != 1.0 {
		t.Errorf("Content should  be 1.0, got %v", f.Content())
	}

	same := values.NewFloat(1.0)
	if !f.Equals(same) {
		t.Error("expected float values to be equal")
	} else if f.ToHashString() != same.ToHashString() {
		t.Error("expected equal float values to have the same hash")
	}

	different := values.NewFloat(2.0)
	if f.Equals(different) {
		t.Error("expected float values NOT to be equal")
	} else if f.ToHashString() == different.ToHashString() {
		t.Error("expected different float values NOT to have the same hash")
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
	} else if a.ToHashString() == b.ToHashString() {
		t.Error("expected different string values NOT to have the same hash")
	} else if !a.Equals(s) {
		t.Error("expected string values to be equal")
	} else if a.ToHashString() != s.ToHashString() {
		t.Error("expected equal string values to have the same hash")
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

func TestBuilder(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	if _, err := values.BuildPrimitiveValue(nil); err == nil {
		t.Error("expected error when building primitive value from nil")
	} else if i, err := values.BuildPrimitiveValue(10); err != nil {
		t.Error("expected no error when building primitive value from int")
	} else if i.Datatype() != values.PRIMITIVE_TYPE_INT {
		t.Error("expected datatype to be int")
	} else if i.Content() != 10 {
		t.Error("expected content to be 10")
	} else if b, err := values.BuildPrimitiveValue(true); err != nil {
		t.Error("expected no error when building primitive value from bool")
	} else if b.Datatype() != values.PRIMITIVE_TYPE_BOOL {
		t.Error("expected datatype to be bool")
	} else if b.Content() != true {
		t.Error("expected content to be true")
	} else if s, err := values.BuildPrimitiveValue("hello"); err != nil {
		t.Error("expected no error when building primitive value from string")
	} else if s.Datatype() != values.PRIMITIVE_TYPE_STRING {
		t.Error("expected datatype to be string")
	} else if s.Content() != "hello" {
		t.Error("expected content to be hello")
	} else if f, err := values.BuildPrimitiveValue(10.5); err != nil {
		t.Error("expected no error when building primitive value from float")
	} else if f.Datatype() != values.PRIMITIVE_TYPE_FLOAT {
		t.Error("expected datatype to be float")
	} else if f.Content() != 10.5 {
		t.Error("expected content to be 10.5")
	} else if tt, err := values.BuildPrimitiveValue(now); err != nil {
		t.Error("expected no error when building primitive value from time")
	} else if tt.Datatype() != values.PRIMITIVE_TYPE_TIME {
		t.Error("expected datatype to be time")
	} else if tt.Content() != now {
		t.Error("expected content to be time")
	}
}
