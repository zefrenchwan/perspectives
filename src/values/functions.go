package values

import (
	"fmt"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

// PrimitiveTimeFunction represents a time-based primitive function.
type PrimitiveTimeFunction struct {
	// primitiveMapping is the mapping reduced to functions (applied to primitive values)
	primitiveMapping[periods.DynamicFunction[PrimitiveValue]]
}

// Equals returns true if the two primitive time functions are equal.
func (pf *PrimitiveTimeFunction) Equals(other PrimitiveTimeFunction) bool {
	return pf.primitiveMapping.Equals(other.primitiveMapping)
}

// At returns the value of the primitive function at the given time.
func (pf *PrimitiveTimeFunction) At(time time.Time) (PrimitiveValue, bool) {
	// instead of "return pf.primitiveMapping.valuesMapping.mapper.At(time)", just use promotion
	return pf.mapper.At(time)
}

// Value returns the underlying value of the primitive function at the given time
func (pf *PrimitiveTimeFunction) Value(time time.Time) (any, bool) {
	result, has := pf.mapper.At(time)
	if !has {
		return nil, false
	}

	return result.value, has
}

// Copy returns a copy of the primitive time function
func (pf *PrimitiveTimeFunction) Copy() PrimitiveTimeFunction {
	base := valuesMapping[PrimitiveValue, periods.DynamicFunction[PrimitiveValue]]{
		mapper: pf.mapper.Copy(),
	}

	primitiveMapper := primitiveMapping[periods.DynamicFunction[PrimitiveValue]]{
		valuesMapping: base,
	}

	return PrimitiveTimeFunction{
		primitiveMapping: primitiveMapper,
	}
}

// buildPrimitiveTimeFunction is private by design.
// It DOES NOT CHECK if that expectedType is a valid primitive type.
func buildPrimitiveTimeFunction(expectedType string) PrimitiveTimeFunction {
	base := valuesMapping[PrimitiveValue, periods.DynamicFunction[PrimitiveValue]]{
		mapper: periods.NewTimeFunction(expectedType, EqualPrimitiveValue),
	}

	primitiveMapper := primitiveMapping[periods.DynamicFunction[PrimitiveValue]]{
		valuesMapping: base,
	}

	return PrimitiveTimeFunction{
		primitiveMapping: primitiveMapper,
	}
}

// NewPrimitiveTimeFunction builds a new primitive time function for a given primitive type.
// If the given name is NOT the name of a primitive type, an error is returned.
func NewPrimitiveTimeFunction(expectedType string) (PrimitiveTimeFunction, error) {
	var empty PrimitiveTimeFunction
	if IsPrimitiveTypeName(expectedType) {
		return buildPrimitiveTimeFunction(expectedType), nil
	}

	return empty, fmt.Errorf("invalid primitive type: %s", expectedType)
}

// NewIntTimeFunction builds a new primitive time function for int values.
func NewIntTimeFunction() PrimitiveTimeFunction {
	return buildPrimitiveTimeFunction(PRIMITIVE_TYPE_INT)
}

// NewStringTimeFunction builds a new primitive time function for string values.
func NewStringTimeFunction() PrimitiveTimeFunction {
	return buildPrimitiveTimeFunction(PRIMITIVE_TYPE_STRING)
}

// ReferenceTimeFunction represents a time function for reference values.
type ReferenceTimeFunction struct {
	// referenceMapping is a function of periods that operate on reference values.
	referenceMapping[periods.DynamicFunction[ReferenceValue]]
}

// Equals checks if two ReferenceTimeFunction instances are equal.
func (rf *ReferenceTimeFunction) Equals(other ReferenceTimeFunction) bool {
	return rf.referenceMapping.Equals(other.referenceMapping)
}

// At returns the reference value at the given time.
func (rf *ReferenceTimeFunction) At(t time.Time) (ReferenceValue, bool) {
	return rf.mapper.At(t)
}

// Value returns the reference id string value at the given time.
func (rf *ReferenceTimeFunction) Value(t time.Time) (string, bool) {
	var empty string
	res, has := rf.mapper.At(t)
	if !has {
		return empty, false
	}

	return res.referenceId, has
}

// Copy returns a copy of the ReferenceTimeFunction.
func (rf *ReferenceTimeFunction) Copy() ReferenceTimeFunction {
	base := valuesMapping[ReferenceValue, periods.DynamicFunction[ReferenceValue]]{
		mapper: rf.mapper.Copy(),
	}

	referenceMapper := referenceMapping[periods.DynamicFunction[ReferenceValue]]{
		valuesMapping: base,
	}

	return ReferenceTimeFunction{
		referenceMapping: referenceMapper,
	}
}
