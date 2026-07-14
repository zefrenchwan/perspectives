package values

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// PrimitiveValue decorates primitives types allowed on entities.
// It is a mandatory gate keeper to use and then ensure entities have good properties (serde for instance).
type PrimitiveValue struct {
	// dataType as listed in values/values.go
	dataType string
	// value is the actual value of the primitive.
	// We used bool, ints, time.Time, string, float64, etc.
	// So common type is any.
	value any
	// serialized is the string representation of the primitive value.
	// It is calculated once and will be used extensively in hashing and serialization.
	serialized string
	// hashString is hashed value of serialized
	hashString string
}

// EqualPrimitiveValue compares two PrimitiveValue for equality.
func EqualPrimitiveValue(a, b PrimitiveValue) bool {
	return a.Equals(b)
}

// isReference forces sealed interface
func (p PrimitiveValue) isReference() bool {
	return false
}

// serializeContent serializes a PrimitiveValue into a string representation.
// For inner purpose only, as soon as an element is built
func serializeContent(p PrimitiveValue) string {
	var buffer strings.Builder
	buffer.WriteString("PRIMITIVE VALUE|")
	buffer.WriteString(p.dataType)
	buffer.WriteString("|")
	var content string
	if asTime, ok := p.value.(time.Time); ok {
		content = asTime.Format(time.RFC3339Nano)
	} else {
		content = fmt.Sprintf("%v", p.value)
	}
	buffer.WriteString(strconv.Itoa(len(content)))
	buffer.WriteString("|")
	buffer.WriteString(content)
	return buffer.String()
}

// Serialize returns the serialized representation of the PrimitiveValue.
func (p PrimitiveValue) Serialize() string {
	return p.serialized
}

// ToHashString returns the hash of the serialized representation of the PrimitiveValue.
func (p PrimitiveValue) ToHashString() string {
	return p.hashString
}

// Datatype returns the data type of the PrimitiveValue.
func (p PrimitiveValue) Datatype() string {
	return p.dataType
}

// Equals checks if two PrimitiveValues are equal.
// For anything that is not a PrimitiveValue, it returns false.
func (p PrimitiveValue) Equals(other any) bool {
	if other == nil {
		return false
	} else if otherValue, ok := other.(PrimitiveValue); !ok {
		return false
	} else if p.dataType != otherValue.dataType {
		return false
	} else {
		pTime, aOk := p.value.(time.Time)
		oTime, bOk := otherValue.value.(time.Time)
		if aOk && bOk {
			return pTime.Equal(oTime)
		}

		return p.value == otherValue.value
	}
}

// Content returns the actual value of the PrimitiveValue.
func (p PrimitiveValue) Content() any {
	return p.value
}

// NewBool creates a new PrimitiveValue of type bool.
func NewBool(value bool) PrimitiveValue {
	result := PrimitiveValue{
		dataType: PRIMITIVE_TYPE_BOOL,
		value:    value,
	}

	result.serialized = serializeContent(result)
	result.hashString = commons.HashString(result.serialized)
	return result
}

// NewInt creates a new PrimitiveValue of type int.
func NewInt(value int) PrimitiveValue {
	result := PrimitiveValue{
		dataType: PRIMITIVE_TYPE_INT,
		value:    value,
	}

	result.serialized = serializeContent(result)
	result.hashString = commons.HashString(result.serialized)
	return result
}

// NewString returns a new PrimitiveValue of type string.
func NewString(value string) PrimitiveValue {
	result := PrimitiveValue{
		dataType: PRIMITIVE_TYPE_STRING,
		value:    value,
	}

	result.serialized = serializeContent(result)
	result.hashString = commons.HashString(result.serialized)
	return result
}

// NewTime returns a new PrimitiveValue of type time.Time.
func NewTime(value time.Time) PrimitiveValue {
	result := PrimitiveValue{
		dataType: PRIMITIVE_TYPE_TIME,
		value:    value,
	}

	result.serialized = serializeContent(result)
	result.hashString = commons.HashString(result.serialized)
	return result
}

// NewFloat creates a new PrimitiveValue of type float64.
func NewFloat(value float64) PrimitiveValue {
	result := PrimitiveValue{
		dataType: PRIMITIVE_TYPE_FLOAT,
		value:    value,
	}

	result.serialized = serializeContent(result)
	result.hashString = commons.HashString(result.serialized)
	return result
}

// BuildPrimitiveValue creates a new PrimitiveValue from the given value.
// It supports bool, int, string, time.Time, and float64 types.
// Otherwise, it returns an error.
func BuildPrimitiveValue(v any) (PrimitiveValue, error) {
	var empty PrimitiveValue
	if v == nil {
		return empty, errors.New("nil value, not a primitive type")
	}
	switch v.(type) {
	case bool:
		b, _ := v.(bool)
		return NewBool(b), nil
	case int:
		i, _ := v.(int)
		return NewInt(i), nil
	case string:
		s, _ := v.(string)
		return NewString(s), nil
	case time.Time:
		t, _ := v.(time.Time)
		return NewTime(t), nil
	case float64:
		f, _ := v.(float64)
		return NewFloat(f), nil
	default:
		return empty, errors.New("unsupported value, not a primitive type")
	}
}

// IsPrimitiveValue checks if the given value is a PrimitiveValue.
func IsPrimitiveValue(v any) bool {
	_, ok := GetPrimitiveType(v)
	return ok
}

// GetPrimitiveType returns the type name of the given value if it is a PrimitiveValue, otherwise, empty string and false.
func GetPrimitiveType(v any) (string, bool) {
	var empty string
	if v == nil {
		return empty, false
	}
	switch v.(type) {
	case int:
		return PRIMITIVE_TYPE_INT, true
	case float64:
		return PRIMITIVE_TYPE_FLOAT, true
	case string:
		return PRIMITIVE_TYPE_STRING, true
	case bool:
		return PRIMITIVE_TYPE_BOOL, true
	case time.Time:
		return PRIMITIVE_TYPE_TIME, true
	default:
		return empty, false
	}
}

// IsPrimitiveTypeName checks if the given name is a valid PrimitiveValue type name.
func IsPrimitiveTypeName(name string) bool {
	switch name {
	case PRIMITIVE_TYPE_INT:
		return true
	case PRIMITIVE_TYPE_FLOAT:
		return true
	case PRIMITIVE_TYPE_STRING:
		return true
	case PRIMITIVE_TYPE_BOOL:
		return true
	case PRIMITIVE_TYPE_TIME:
		return true
	default:
		return false
	}
}
