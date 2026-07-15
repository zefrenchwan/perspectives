package values

import (
	"fmt"
	"iter"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

type ValuesDynamicFunction[V Value] interface {
	Datatype() string
	Domain() periods.Period
	IsEmpty() bool
	Range() iter.Seq2[periods.Period, V]
	At(time time.Time) (V, bool)
}

type valuesFunction[V Value] struct {
	valuesMapping[V, periods.DynamicFunction[V]]
}

func (vf *valuesFunction[V]) Equals(other valuesFunction[V]) bool {
	return vf.valuesMapping.Equals(other.valuesMapping)
}

func (vf *valuesFunction[V]) At(time time.Time) (V, bool) {
	return vf.mapper.At(time)
}

func (vf *valuesFunction[V]) Value(time time.Time) (any, bool) {
	result, has := vf.mapper.At(time)
	if !has {
		return nil, false
	}

	return result.Content(), has
}

// PrimitiveTimeFunction represents a time-based primitive function.
type PrimitiveTimeFunction struct {
	valuesFunction[PrimitiveValue]
}

// Equals returns true if the two primitive time functions are equal.
func (pf *PrimitiveTimeFunction) Equals(other PrimitiveTimeFunction) bool {
	return pf.mapper.Equals(other.mapper)
}

// Copy returns a copy of the primitive time function
func (pf *PrimitiveTimeFunction) Copy() PrimitiveTimeFunction {
	base := valuesMapping[PrimitiveValue, periods.DynamicFunction[PrimitiveValue]]{
		mapper: pf.mapper.Copy(),
	}

	function := valuesFunction[PrimitiveValue]{
		valuesMapping: base,
	}

	return PrimitiveTimeFunction{
		valuesFunction: function,
	}
}

// Add adds a value for that given period: it creates the primitive value and adds it to the mapping.
func (pm *PrimitiveTimeFunction) Add(value any, period periods.Period) error {
	expectedType := pm.mapper.DataType()
	newValue, errBuild := BuildPrimitiveValue(value)
	if errBuild != nil {
		return errBuild
	} else if realType := newValue.Datatype(); realType != expectedType {
		return fmt.Errorf("value type %s does not match expected type %s", realType, expectedType)
	}

	pm.mapper.Add(newValue, period)
	return nil
}

// buildPrimitiveTimeFunction is private by design.
// It DOES NOT CHECK if that expectedType is a valid primitive type.
func buildPrimitiveTimeFunction(expectedType string) PrimitiveTimeFunction {
	base := valuesMapping[PrimitiveValue, periods.DynamicFunction[PrimitiveValue]]{
		mapper: periods.NewTimeFunction(expectedType, EqualPrimitiveValue),
	}

	function := valuesFunction[PrimitiveValue]{
		valuesMapping: base,
	}

	return PrimitiveTimeFunction{
		valuesFunction: function,
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
	valuesFunction[ReferenceValue]
}

// Equals checks if two ReferenceTimeFunction instances are equal.
func (rf *ReferenceTimeFunction) Equals(other ReferenceTimeFunction) bool {
	return rf.mapper.Equals(other.mapper)
}

// Add adds a reference (as a string) for that given period: it creates the reference value and adds it to the mapping.
func (rf *ReferenceTimeFunction) Add(reference string, period periods.Period) {
	referenceValue := NewReference(reference)
	rf.mapper.Add(referenceValue, period)
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

	function := valuesFunction[ReferenceValue]{
		valuesMapping: base,
	}

	return ReferenceTimeFunction{
		valuesFunction: function,
	}
}

// NewReferenceTimeFunction creates a new time dependent function of references.
func NewReferenceTimeFunction() ReferenceTimeFunction {
	base := valuesMapping[ReferenceValue, periods.DynamicFunction[ReferenceValue]]{
		mapper: periods.NewTimeFunction(REFERENCE_TYPE, EqualReferences),
	}

	function := valuesFunction[ReferenceValue]{
		valuesMapping: base,
	}

	return ReferenceTimeFunction{
		valuesFunction: function,
	}
}
