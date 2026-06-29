package entities

import (
	"fmt"
	"time"
)

// ===================================================
// DESIGN INFORMATION
// This code was optimized for performance.
// It means that if you change primitive types, you should check the full code file.
// Many times here, code contains switch statements with duplicated code.
// It means that adding or removing a primitive type might fail in the switch statements.
// ===================================================

// primitiveActions store operands to apply on primitive types.
type primitiveActions struct {
	// equals operator on a specified type
	equals func(any, any) bool
	// toString operator on a specified type
	toString func(any) string
}

// defaultValues returns default primitive actions (int, string, bool)
func defaultValues() primitiveActions {
	return primitiveActions{
		equals:   defaultEquals,
		toString: defaultString,
	}
}

// float32Actions returns primitive actions for float32 type
func float32Actions() primitiveActions {
	return primitiveActions{
		equals:   equalsFloat32,
		toString: defaultString,
	}
}

// float64Actions returns primitive actions for float64 type
func float64Actions() primitiveActions {
	return primitiveActions{
		equals:   equalsFloat64,
		toString: defaultString,
	}
}

// timeActions returns primitive actions for time.Time type
func timeActions() primitiveActions {
	return primitiveActions{
		equals:   equalsTime,
		toString: timeString,
	}
}

// allowedPrimitives associates the name of the primitive type with the corresponding equality function.
// It is NOT the unique source of truth, code was optimized for performance.
var allowedPrimitives = map[string]primitiveActions{
	"int":       defaultValues(),
	"int32":     defaultValues(),
	"int64":     defaultValues(),
	"float32":   float32Actions(),
	"float64":   float64Actions(),
	"string":    defaultValues(),
	"bool":      defaultValues(),
	"time.Time": timeActions(),
}

// IsPrimitiveTypeName checks if the given name is a primitive type name.
func IsPrimitiveTypeName(name string) bool {
	_, ok := allowedPrimitives[name]
	return ok
}

// primitiveTypeName returns the name of the primitive type if it is a primitive type, otherwise an empty string.
func primitiveTypeName(v any) string {
	if v == nil {
		return ""
	}
	switch v.(type) {
	case int:
		return "int"
	case int32:
		return "int32"
	case int64:
		return "int64"
	case float32:
		return "float32"
	case float64:
		return "float64"
	case string:
		return "string"
	case bool:
		return "bool"
	case time.Time:
		return "time.Time"
	default:
		return ""
	}
}

// primitiveTypeEqualsFunc returns the function to use for comparing primitive types.
func primitiveTypeEqualsFunc(typeName string) func(any, any) bool {
	if fn, ok := allowedPrimitives[typeName]; ok {
		return fn.equals
	}
	return defaultEquals
}

// IsPrimitiveValue checks whether any is a primitive type instance or not.
func IsPrimitiveValue(v any) bool {
	return primitiveTypeName(v) != ""
}

// primitiveValueToString returns the string representation of a primitive value.
// It uses the right serializer for that given value based on its type.
func primitiveValueToString(v any) string {
	if value, ok := v.(time.Time); ok {
		return timeString(value)
	} else {
		return fmt.Sprintf("%v", v)
	}
}

// ===========================================================================
// STRING FUNCTIONS FOR DEDICATED TYPES
// ===========================================================================

// defaultString is how to convert a primitive value to a string by default.
func defaultString(v any) string {
	if v == nil {
		return ""
	}

	return fmt.Sprintf("%v", v)
}

// timeString is dedicated to time.Time values.
// It returns either an empty string (not a time.Time instance) or a formatted time string.
// It is NOT based on the configuration : we want to serialize time.Time values in a standard way.
func timeString(v any) string {
	if v == nil {
		return ""
	}
	t, ok := v.(time.Time)
	if !ok {
		return ""
	}
	// DO NOT USE CONFIGURATION HERE : CHANGE OF CONF => NO NEED TO RECALCULATE ALL HASHES
	return t.UTC().Format(time.RFC3339Nano)
}

// ===========================================================================
// EQUALS FUNCTIONS FOR DEDICATED TYPES
// ===========================================================================

// defaultEquals is what we use for primitive types that don't have a dedicated equals function.
func defaultEquals(a, b any) bool {
	return a == b
}

// equalsTime tests equality between two time.Time values.
func equalsTime(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	t1, ok1 := a.(time.Time)
	t2, ok2 := b.(time.Time)

	if !ok1 || !ok2 {
		return false
	}
	return t1.Equal(t2)
}

// equalsFloat32 tests equality between two float32 values
func equalsFloat32(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	v1, ok1 := a.(float32)
	v2, ok2 := b.(float32)

	if !ok1 || !ok2 {
		return false
	}
	return v1 == v2
}

// equalsFloat64 tests equality between two float64 values
func equalsFloat64(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	v1, ok1 := a.(float64)
	v2, ok2 := b.(float64)

	if !ok1 || !ok2 {
		return false
	}
	return v1 == v2
}
