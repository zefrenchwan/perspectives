package objects

import (
	"math"
	"reflect"
	"time"

	"github.com/zefrenchwan/perspectives.git/configuration"
)

// allowedPrimitives is the part to change to add / remove primitive types.
// As a single source of truth, it associates the name of the primitive type with the corresponding equality function.
var allowedPrimitives = map[string]func(any, any) bool{
	"int": defaultEquals,
	//"int8":      defaultEquals,
	//"int16":     defaultEquals,
	"int32": defaultEquals,
	"int64": defaultEquals,
	//"uint":      defaultEquals,
	//"uint8":     defaultEquals,
	//"uint16":    defaultEquals,
	//"uint32":    defaultEquals,
	//"uint64":    defaultEquals,
	"float32":   equalsFloat32,
	"float64":   equalsFloat64,
	"string":    defaultEquals,
	"bool":      defaultEquals,
	"time.Time": equalsTime,
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
		return fn
	}
	return defaultEquals // Fallback par défaut ou nil selon tes préférences
}

// IsPrimitiveValue checks whether any is a primitive type instance or not.
func IsPrimitiveValue(v any) bool {
	return primitiveTypeName(v) != ""
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
