package commons

import "slices"

// objectPredicate defines a predicate condition
type objectPredicate[T any] struct {
	// predicate to accept or reject an element
	predicate func(T) bool
	// variable name to pick in the content
	variable string
}

// Signature returns the expected parameters
func (p objectPredicate[T]) Signature() FormalParameters {
	return NewNamedFormalParameters([]string{p.variable})
}

// Matches returns true if this condition accepts content, or an error if any.
func (p objectPredicate[T]) Matches(content Content) (bool, error) {
	if content == nil {
		return false, nil
	} else if p.predicate == nil {
		return false, nil
	} else if value, found := content.GetByName(p.variable); !found {
		return false, nil
	} else if obj, ok := value.(T); !ok {
		return false, nil
	} else {
		return p.predicate(obj), nil
	}
}

// NewFilterPredicate is a named condition for a given predicate.
// It accepts a predicate on any type (hence the type parameter).
// It applies on variables only based on a given predicate.
func NewFilterPredicate[T any](variable string, predicate func(T) bool) Condition {
	if predicate == nil {
		return nil
	}

	return objectPredicate[T]{variable: variable, predicate: predicate}
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

// NewFilterById returns a new condition for a variable to have a given id.
// If variable = x and expected id is "id", then condition is x.id == "id".
func NewFilterById(variable string, expectedId string) Condition {
	return NewFilterPredicate(variable, func(e Modelable) bool {
		if e == nil {
			return false
		} else if value, ok := e.(Identifiable); !ok {
			return false
		} else if value == nil {
			return false
		} else {
			return value.Id() == expectedId
		}
	})
}

// NewFilterByTypes creates a new condition for a variable to have its type in expected types
func NewFilterByTypes(variable string, expectedTypes []ModelableType) Condition {
	return NewFilterPredicate(variable, func(e Modelable) bool {
		return e != nil && len(expectedTypes) != 0 && slices.Contains(expectedTypes, e.GetType())
	})
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

// localStateValueSelector defines a condition variable.attribute operator reference.
// For instance x.price < 0, or x.firstname = "Heinrich"
type localStateValueSelector[T StateValue] struct {
	// variable is the name of the variable
	variable string
	// attribute is the attribute to pick in the state
	attribute string
	// setOperator is the operator to go through references
	setOperator LocalSetOperator[T]
	// reference is the constant value to compare to loaded data
	references []T
}

// Signature defines expected variable
func (l localStateValueSelector[T]) Signature() FormalParameters {
	return NewNamedFormalParameters([]string{l.variable})
}

// Matches uses the operators as the inner condition:
// it runs through references via the setOperator and applies to each element its operator
func (l localStateValueSelector[T]) Matches(c Content) (bool, error) {
	if l.setOperator == nil {
		return false, nil
	} else if value, found := c.GetByName(l.variable); !found {
		return false, nil
	} else if value == nil {
		return false, nil
	} else if v, ok := value.(StateReader[T]); !ok {
		return false, nil
	} else if v == nil {
		return false, nil
	} else if state := v.Read(); state == nil {
		return false, nil
	} else if values := state.Values(); len(values) == 0 {
		return false, nil
	} else if operand, found := values[l.attribute]; !found {
		return false, nil
	} else {
		return l.setOperator.Accepts(operand, l.references), nil
	}
}

// NewFilterByStateOperator returns a new condition to compare an attribute value to a reference for that operator.
// For instance object.x >= 0 (x coordinate for object should be positive).
func NewFilterByStateOperator[T StateValue](variable, attribute string, operator LocalOperator[T], reference T) Condition {
	if operator == nil {
		return nil
	}

	return localStateValueSelector[T]{
		variable:    variable,
		attribute:   attribute,
		setOperator: NewLocalSetOperator(MatchesOneInSetOperator, operator),
		references:  []T{reference},
	}
}

// NewFilterByStateSetOperator reads a value from an attribute and compares it to the reference set based on the operator.
// For instance, to test if variable.attribute is in a given set of int values,
// use NewLocalSetOperator(MatchesOneInSetOperator, IntEquals) for the operator.
func NewFilterByStateSetOperator[T StateValue](variable, attribute string, operator LocalSetOperator[T], reference []T) Condition {
	if operator == nil {
		return nil
	}

	return localStateValueSelector[T]{
		variable:    variable,
		attribute:   attribute,
		setOperator: operator,
		references:  reference,
	}
}
