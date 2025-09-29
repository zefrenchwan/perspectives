package engines

import (
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/zefrenchwan/perspectives.git/models"
	"github.com/zefrenchwan/perspectives.git/structures"
)

// LocalCondition is a boolean condition to accept a single entity
type LocalCondition interface {
	// Matches returns true if the entity matches the condition
	Matches(models.ModelEntity) bool
}

// ValueOperator applies to values (such as string) for conditions definition
type ValueOperator int

// ValuesEqual test if values are the same (including case)
const ValuesEqual ValueOperator = 1

// ValuesEqualIgnoreCase tests if values are the same no matter the case
const ValuesEqualIgnoreCase ValueOperator = 2

// ValueMatchesRegexp tests if a value matches a given regexp as a reference
const ValueMatchesRegexp ValueOperator = 3

// valuesOperatorRun applies the ValueOperator on source (parameter) and reference (from condition)
func valuesOperatorRun(source, reference string, operator ValueOperator) bool {
	switch operator {
	case ValuesEqual:
		return source == reference
	case ValuesEqualIgnoreCase:
		return strings.EqualFold(source, reference)
	case ValueMatchesRegexp:
		if validator := regexp.MustCompile(reference); validator == nil {
			return false
		} else {
			return validator.MatchString(source)
		}
	default:
		return false
	}
}

// PeriodOperator is the type to define a condition on periods
type PeriodOperator int

// NonDisjoinPeriods tests if periods have a common point
const NonDisjoinPeriods PeriodOperator = 1

// SamePeriods tests if periods are equals
const SamePeriods PeriodOperator = 2

// AcceptsAllPeriods returns true no matter the condition
const AcceptsAllPeriods PeriodOperator = 3

// IsIncludedInReferencePeriod tests if current period is in reference period
const IsIncludedInReferencePeriod PeriodOperator = 4

// periodOperatorRun returns true if operator applied to current and reference matches
func periodOperatorRun(current, reference structures.Period, operator PeriodOperator) bool {
	switch operator {
	case IsIncludedInReferencePeriod:
		return reference.Intersection(current).Equals(reference)
	case AcceptsAllPeriods:
		return true
	case NonDisjoinPeriods:
		return !current.Intersection(reference).IsEmpty()
	case SamePeriods:
		return current.Equals(reference)
	default:
		return false
	}
}

// LogicalCombinationOperator defines basic logic operators: not, and, or
type LogicalCombinationOperator int

// NegateCombineConditions defines not as a logical operator
const NegateCombineConditions LogicalCombinationOperator = 1

// OrCombineConditions defines or as a logical operator
const OrCombineConditions LogicalCombinationOperator = 2

// AndCombineConditions defines and as a logical operator
const AndCombineConditions LogicalCombinationOperator = 3

// AlwaysTrueCombiner defines a constant logical operator as true
const AlwaysTrueCombiner LogicalCombinationOperator = 100

// AlwaysFalseCombiner defines a constant logical operator as false
const AlwaysFalseCombiner LogicalCombinationOperator = 101

// LocalCompositeCondition defines a local condition as a combination of local other conditions.
// For instance, if a and b are conditions, then a AND b is a condition
type LocalCompositeCondition struct {
	id       string                     // id of the node
	operator LogicalCombinationOperator // operator to combine conditions
	operands []LocalCondition           // the conditions to combine
}

// evaluate returns the result of operation applied to booleans.
// For instance, if composite is AND, then it returns true if all values are true.
// By default, for no value, it returns false no matter the operator.
// Otherwise, it just applies:
// "and" and "or" to all values
// "not" only to first value
func (l LocalCompositeCondition) evaluate(values []bool) bool {
	if len(values) == 0 {
		return false
	}

	switch l.operator {
	case NegateCombineConditions:
		return !values[0]
	case OrCombineConditions:
		for _, value := range values {
			if value {
				return true
			}
		}

		return false
	case AndCombineConditions:
		for _, value := range values {
			if !value {
				return false
			}
		}

		return true
	default:
		return false
	}
}

// Matches returns the combined results of its operands.
// To avoid stack forcing with a recursive solution, it is a tree walkthrough.
// We look for leafs to evaluate and then walk back until the root
func (c LocalCompositeCondition) Matches(m models.ModelEntity) bool {
	// childs contain composite childs of a composite node.
	childs := make(map[string][]string)
	// content gets the value of a composition node
	content := make(map[string]LocalCompositeCondition)
	// resolved contains all composite conditions evaluated.
	// At the end of the algorithm, we know root
	resolved := make(map[string]bool)
	// parents maps a composite node to its parent
	parents := make(map[string]string)

	// BFS to get composite conditions structure
	queue := []LocalCompositeCondition{c}
	// for each node
	for len(queue) != 0 {
		current := queue[0]
		queue = queue[1:]
		// map it to its actual value
		content[current.id] = current

		// true of all childs are leafs (and then node may be evaluated)
		allLeafs := true
		for _, child := range current.operands {
			if childComposite, ok := child.(LocalCompositeCondition); ok {
				// register structure
				currentChilds := childs[current.id]
				currentChilds = append(currentChilds, childComposite.id)
				childs[current.id] = currentChilds
				parents[childComposite.id] = current.id
				// add the child node for a later processing
				queue = append(queue, childComposite)
				// node is not a leaf
				allLeafs = false
			}
		}

		// if allLeafs, we may evaluate current node
		if allLeafs {
			var operands []bool
			for _, child := range current.operands {
				if child != nil {
					operands = append(operands, child.Matches(m))
				}
			}

			localResult := current.evaluate(operands)
			// flag the node as resolved
			resolved[current.id] = localResult
		}
	}

	// deal with special case: root is already solved
	if result, found := resolved[c.id]; found {
		return result
	}

	// it means that root is not resolved.
	// We dealt with leafs, we know the graph structure.
	// So, starting from resolved nodes, go a level higher if possible.
	currents := make(map[string]bool)
	maps.Copy(currents, resolved)

	// until we solve the root, go from childs to parent
	for foundRoot := currents[c.id]; !foundRoot; {
		// we know currents, we feed nexts with parents of nodes to process if possible
		nexts := make(map[string]bool)
		// for each node we get solution for
		for node, value := range currents {
			// special case: root is already solved
			if node == c.id {
				return value
			}

			// Root is not resolved, so each node in there should have a parent.
			// Find parents and try to solve them
			parent := parents[node]
			// true if the brothers of the current node are ALL solved.
			// It means that we may evaluate the parent and move up
			allSolved := true
			// in case of all brothers are set, collect them right now
			var allBrothersSolutions []bool
			// for each child of the parent, that is, for each brother of the node
			for _, brother := range childs[parent] {
				// if we already solved that brother
				if value, solved := resolved[brother]; !solved {
					allSolved = false
				} else {
					allBrothersSolutions = append(allBrothersSolutions, value)
				}
			}

			// if allSolved is true, it means that we may evaluate the parent.
			// All operands to evaluate parent are in allBrothersSolutions
			if allSolved {
				parentOperation := content[parent]
				solution := parentOperation.evaluate(allBrothersSolutions)
				resolved[parent] = solution
				nexts[parent] = solution
			}
		}

		if len(nexts) == 0 {
			// cannot happen but stop loop
			break
		}

		// update currents
		currents = nexts
	}

	// at this point for sure, tree is completely solved, so find root solution
	return currents[c.id]
}

// NotCondition returns the not condition as a local condition.
// Special case: if condition is nil, result will return false no matter the value
func NotCondition(condition LocalCondition) LocalCondition {
	if condition == nil {
		return LocalCompositeCondition{id: uuid.NewString(), operator: AlwaysFalseCombiner}
	}

	return LocalCompositeCondition{
		id:       uuid.NewString(),
		operator: NegateCombineConditions,
		operands: []LocalCondition{condition},
	}
}

// buildLocalConditionsCombiner picks only not null conditions and build the combined local condition.
// Special case: only nil values make the "always false" condition
func buildLocalConditionsCombiner(conditions []LocalCondition, operation LogicalCombinationOperator) LocalCondition {
	var operands []LocalCondition
	for _, operand := range conditions {
		if operand != nil {
			operands = append(operands, operand)
		}
	}

	if len(operands) == 0 {
		return LocalCompositeCondition{
			id:       uuid.NewString(),
			operator: AlwaysFalseCombiner,
		}
	}

	return LocalCompositeCondition{
		id:       uuid.NewString(),
		operator: operation,
		operands: operands,
	}
}

// OrConditions builds an OR applied to non nil conditions parameter.
// Special case: if conditions contains no nil condition, then result is the "always false" condition
func OrConditions(conditions []LocalCondition) LocalCondition {
	return buildLocalConditionsCombiner(conditions, OrCombineConditions)
}

// AndConditions builds an AND applied to non nil conditions parameter.
// Special case: if conditions contains no nil condition, then result is the "always false" condition
func AndConditions(conditions []LocalCondition) LocalCondition {
	return buildLocalConditionsCombiner(conditions, AndCombineConditions)
}

// LocalTypeCondition accepts entities if they match a type.
// For instance, a condition may be "only links" or "links or objects".
type LocalTypeCondition struct {
	// values are the types to accept.
	values []models.EntityType
}

// NewTypeCondition returns a condition that accepts only a set of types.
func NewTypeCondition(types ...models.EntityType) LocalCondition {
	return LocalTypeCondition{values: types}
}

// Matches for a LocalTypeCondition is true if the entity's type is in the list of accepted types.
// Note that nil value is rejected (because we cannot tell its type for sure)
func (l LocalTypeCondition) Matches(e models.ModelEntity) bool {
	// No type accepted => false
	if len(l.values) == 0 {
		return false
	} else if e == nil {
		// Nil => false
		return false
	} else {
		return slices.Contains(l.values, e.GetType())
	}
}

// LocalActiveCondition means that a temporal entity is active during that reference period
type LocalActiveCondition struct {
	ReferencePeriod structures.Period
}

// Matches returns true if m is a temporal entity with a period included in reference period
func (l LocalActiveCondition) Matches(m models.ModelEntity) bool {
	if m == nil {
		return false
	} else if t, ok := m.(models.TemporalEntity); !ok {
		return false
	} else {
		return periodOperatorRun(t.ActivePeriod(), l.ReferencePeriod, IsIncludedInReferencePeriod)
	}
}

// NewActiveCondition builds a condition to be active during a provided period.
// Condition is true if object has an active period that contains reference (or is equals to)
func NewActiveCondition(reference structures.Period) LocalCondition {
	return LocalActiveCondition{ReferencePeriod: reference}
}

// LocalMatchingAttributeCondition is a condition for an object attribute on a given period.
// For instance, nationality = "french" during a period
type LocalMatchingAttributeCondition struct {
	AttributeName     string            // Name of the attribute to find in the object
	AttributeValue    string            // Value of the attribute to compare with
	AttributeOperator ValueOperator     // Operator for value (such as equals)
	ReferencePeriod   structures.Period // Period to match attribute during
	PeriodOperator    PeriodOperator    // Operator for the period (such as with at least a common point)
}

// Matches returns true if all of those conditions apply :
// The parameter is indeed an object and has that attribute
// The condition on attribute compared to value matches
// The period of matching is acceptable regarding the period condition
func (l LocalMatchingAttributeCondition) Matches(e models.ModelEntity) bool {
	if e == nil {
		return false
	} else if e.GetType() != models.EntityTypeObject {
		return false
	}

	object, _ := e.AsObject()
	if matchingValues, found := object.GetValue(l.AttributeName); !found {
		return false
	} else {
		for value, period := range matchingValues {
			if valuesOperatorRun(value, l.AttributeValue, l.AttributeOperator) {
				if periodOperatorRun(period, l.ReferencePeriod, l.PeriodOperator) {
					return true
				}
			}
		}
	}

	return false
}

// NewAttributeValueCondition builds a time independant condition on attribute and value.
// For instance, attribute = value no matter the time.
func NewAttributeValueCondition(attribute, value string, operator ValueOperator) LocalCondition {
	return LocalMatchingAttributeCondition{
		AttributeName:     attribute,
		AttributeValue:    value,
		AttributeOperator: operator,
		ReferencePeriod:   structures.NewFullPeriod(),
		PeriodOperator:    AcceptsAllPeriods,
	}
}

// NewAttributeRegexpCondition builds a condition for an attribute value matching a regexp no matter the period.
// The regexp definition may not be valid, so we return an error if so.
func NewAttributeRegexpCondition(attribute, regexpDefinition string) (LocalCondition, error) {
	if _, err := regexp.Compile(regexpDefinition); err != nil {
		return nil, err
	}

	return LocalMatchingAttributeCondition{
		AttributeName:     attribute,
		AttributeValue:    regexpDefinition,
		AttributeOperator: ValueMatchesRegexp,
		ReferencePeriod:   structures.NewFullPeriod(),
		PeriodOperator:    AcceptsAllPeriods,
	}, nil
}

// LocalLinkNameValuesCondition tests if the name of a link matches a condition.
// For instance, test if link's name equals (ignore case) a specific value in a provided list.
// Note that the operator applies to each element in LinkValues, not the values as a whole
type LocalLinkNameValueCondition struct {
	LinkValues    []string      // operand to match
	ValueOperator ValueOperator // operation to test link on (for instance equals as string)
}

// Matches returns true if the condition applies to that link's name.
// Algorithm is to go through all the values in the condition and to get a match
func (l LocalLinkNameValueCondition) Matches(m models.ModelEntity) bool {
	if m == nil {
		return false
	} else if m.GetType() != models.EntityTypeLink {
		return false
	}

	link, _ := m.AsLink()
	value := link.Name()
	for _, option := range l.LinkValues {
		if valuesOperatorRun(value, option, l.ValueOperator) {
			return true
		}
	}

	return false
}

// NewLinkNameInValuesCondition builds a condition for a link to have a name in options for that operator.
// For instance, operator may be equals ignore case, options a set of values.
// An entity would pass that condition if it is a link and if its name is included in the set of LinkValues
func NewLinkNameInValuesCondition(options []string, operator ValueOperator) LocalCondition {
	values := structures.SliceDeduplicate(options)
	return LocalLinkNameValueCondition{LinkValues: values, ValueOperator: operator}
}
