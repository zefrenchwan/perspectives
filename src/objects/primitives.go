package objects

import (
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/zefrenchwan/perspectives.git/configuration"
)

type primitiveActions struct {
	equals   func(any, any) bool
	toString func(any) string
}

func defaultValues() primitiveActions {
	return primitiveActions{
		equals:   defaultEquals,
		toString: defaultString,
	}
}

func float32Actions() primitiveActions {
	return primitiveActions{
		equals:   equalsFloat32,
		toString: defaultString,
	}
}

func float64Actions() primitiveActions {
	return primitiveActions{
		equals:   equalsFloat32,
		toString: defaultString,
	}
}

func timeActions() primitiveActions {
	return primitiveActions{
		equals:   equalsTime,
		toString: timeString,
	}
}

// allowedPrimitives is the part to change to add / remove primitive types.
// As a single source of truth, it associates the name of the primitive type with the corresponding equality function.
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

	// reflect.TypeOf().String() is slower than case switch for sure,
	// BUT it reduces readability.
	name := reflect.TypeOf(v).String()
	if IsPrimitiveTypeName(name) {
		return name
	}
	return ""
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

func primitiveValueToString(v any) string {
	if v == nil {
		return ""
	}
	name := reflect.TypeOf(v).String()
	if fn, ok := allowedPrimitives[name]; ok {
		return fn.toString(v)
	}
	return fmt.Sprintf("%v", v)
}

// ===========================================================================
// STRING FUNCTIONS FOR DEDICATED TYPES
// ===========================================================================

func defaultString(v any) string {
	if v == nil {
		return ""
	}

	return fmt.Sprintf("%v", v)
}

func timeString(v any) string {
	if v == nil {
		return ""
	}
	t, ok := v.(time.Time)
	if !ok {
		return ""
	}
	return t.Format(configuration.TIME_FORMAT)
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

// equalsFloat32 tests equality between two float32 values using the SHORT_EPSILON constant.
func equalsFloat32(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	v1, ok1 := a.(float32)
	v2, ok2 := b.(float32)

	if !ok1 || !ok2 {
		return false
	}
	return math.Abs(float64(v1-v2)) < configuration.SHORT_EPSILON
}

// equalsFloat64 tests equality between two float64 values using the LONG_EPSILON constant.
func equalsFloat64(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	v1, ok1 := a.(float64)
	v2, ok2 := b.(float64)

	if !ok1 || !ok2 {
		return false
	}
	return math.Abs(v1-v2) < configuration.LONG_EPSILON
}
