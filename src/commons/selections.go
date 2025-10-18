package commons

import "slices"

// unarySelector is the defintion of an unary selection over a Modelable.
// It reads a variable from content (by name), then performs a test over reference and a modelable.
// For instance:
// to test an id,
// reference would be expected id,
// acceptance reads a source and compares its id wih reference.
type unarySelector[T any] struct {
	// variable to read from content
	variable string
	// reference value for that condition.
	// We accept the source based on that reference.
	reference T
	// acceptance is the match condition:
	// we read a part of the source and compare with reference.
	// It returns true or false if acceptance, an error if any.
	// Returning true, nil means accepting the source.
	acceptance func(source Modelable, reference T) (bool, error)
}

// Signature returns formal parameters for its variable
func (s unarySelector[T]) Signature() FormalParameters {
	return NewNamedFormalParameters([]string{s.variable})
}

// Matches uses the acceptance function as the inner condition
func (s unarySelector[T]) Matches(c Content) (bool, error) {
	if value, found := c.GetByName(s.variable); !found {
		return false, nil
	} else if value == nil {
		return false, nil
	} else if s.acceptance == nil {
		return false, nil
	} else if res, err := s.acceptance(value, s.reference); err != nil {
		return false, err
	} else {
		return res, nil
	}
}

// localOperatorSelector defines  a condition to deal with operators.
// Algorithm is:
// to pick variable from the content,
// extract value from content (using extractor) if any,
// compare value with reference using the operator
type localOperatorSelector[T any] struct {
	// variable defines which name to pick
	variable string
	// operator is the operator per se
	operator LocalOperator[T]
	// reference is the reference value to compare extracted value
	reference T
	// extractor picks the value from the modelable and returns the value if any.
	// Result is the extracted value, true if extractor managed to extract the value, an error if any
	extractor func(Modelable) (T, bool, error)
}

// Signature returns formal parameters accepting variable
func (l localOperatorSelector[T]) Signature() FormalParameters {
	return NewNamedFormalParameters([]string{l.variable})
}

// Matches returns true if operator applied to content matches reference
func (l localOperatorSelector[T]) Matches(c Content) (bool, error) {
	if l.extractor == nil {
		return false, nil
	} else if l.operator == nil {
		return false, nil
	} else if c == nil {
		return false, nil
	} else if value, found := c.GetByName(l.variable); !found {
		return false, nil
	} else if res, matches, err := l.extractor(value); err != nil {
		return false, err
	} else if !matches {
		return false, nil
	} else {
		return l.operator.Accepts(res, l.reference), nil
	}
}

// acceptsModelableByid returns true if modelable is identifiable with that id, false otherwise
func acceptsModelableByid(source Modelable, reference string) (bool, error) {
	if value, ok := source.(Identifiable); !ok {
		return false, nil
	} else if value == nil {
		return false, nil
	} else {
		return value.Id() == reference, nil
	}
}

// NewFilterById returns a new condition for a variable to have a given id.
// If variable = x and expected id is "id", then condition is x.id == "id".
func NewFilterById(variable string, expectedId string) Condition {
	var result unarySelector[string]
	result.variable = variable
	result.reference = expectedId
	result.acceptance = acceptsModelableByid
	return result
}

// acceptModelableByTypes accepts source if it is not nil and its type is in reference (as a set)
func acceptModelableByTypes(source Modelable, reference []ModelableType) (bool, error) {
	if source == nil || len(reference) == 0 {
		return false, nil
	} else {
		currentType := source.GetType()
		return slices.ContainsFunc(reference, func(a ModelableType) bool { return a == currentType }), nil
	}
}

// NewFilterByTypes creates a new condition for a variable to have its type in expected types
func NewFilterByTypes(variable string, expectedTypes []ModelableType) Condition {
	var result unarySelector[[]ModelableType]
	result.variable = variable
	result.reference = SliceDeduplicate(expectedTypes)
	result.acceptance = acceptModelableByTypes
	return result
}

// activePeriodExtractor extracts the active period of m, if any.
// It returns the active period (if any), true if m is temporal, an error if any
func activePeriodExtractor(m Modelable) (Period, bool, error) {
	var empty Period
	if m == nil {
		return empty, false, nil
	} else if value, ok := m.(TemporalReader); !ok {
		return empty, false, nil
	} else if value == nil {
		return empty, false, nil
	} else {
		return value.ActivePeriod(), true, nil
	}
}

// NewFilterActivePeriod returns a new condition to read active period from variable and compare it to period.
// Parameters are:
// name of the variable to pick once a context is given,
// the temporal operator,
// period as the reference period
func NewFilterActivePeriod(variable string, operator TemporalOperator, period Period) Condition {
	var result localOperatorSelector[Period]
	result.operator = operator
	result.reference = period
	result.variable = variable
	result.extractor = activePeriodExtractor
	return result
}
