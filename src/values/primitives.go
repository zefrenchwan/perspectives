package values

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
	return result
}

// NewInt creates a new PrimitiveValue of type int.
func NewInt(value int) PrimitiveValue {
	result := PrimitiveValue{
		dataType: PRIMITIVE_TYPE_INT,
		value:    value,
	}

	result.serialized = serializeContent(result)
	return result
}

// NewString returns a new PrimitiveValue of type string.
func NewString(value string) PrimitiveValue {
	result := PrimitiveValue{
		dataType: PRIMITIVE_TYPE_STRING,
		value:    value,
	}

	result.serialized = serializeContent(result)
	return result
}

// NewTime returns a new PrimitiveValue of type time.Time.
func NewTime(value time.Time) PrimitiveValue {
	result := PrimitiveValue{
		dataType: PRIMITIVE_TYPE_TIME,
		value:    value,
	}

	result.serialized = serializeContent(result)
	return result
}

// NewFloat creates a new PrimitiveValue of type float64.
func NewFloat(value float64) PrimitiveValue {
	result := PrimitiveValue{
		dataType: PRIMITIVE_TYPE_FLOAT,
		value:    value,
	}

	result.serialized = serializeContent(result)
	return result
}

// IsPrimitiveValue checks if the given value is a PrimitiveValue.
func IsPrimitiveValue(v any) bool {
	if v == nil {
		return false
	}
	switch v.(type) {
	case int:
		return true
	case float64:
		return true
	case string:
		return true
	case bool:
		return true
	case time.Time:
		return true
	default:
		return false
	}
}

// IsPrimitiveTypeName checks if the given name is a valid PrimitiveValue type name.
func IsPrimitiveTypeName(name string) bool {
	switch name {
	case "int":
		return true
	case "float64":
		return true
	case "string":
		return true
	case "bool":
		return true
	case "time.Time":
		return true
	default:
		return false
	}
}
