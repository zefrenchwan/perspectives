package commons

import (
	"errors"
	"maps"
	"slices"
)

// Condition is the most abstract defintion of a condition to match.
// A condition is said to be a conditional operator if it is a "not", "and" or "or".
// Otherwise, it is a leaf in a logical expression.
type Condition interface {
	// Signature returns the expected parameters
	Signature() FormalParameters
	// Matches returns true if this condition accepts content.
	// It may return an error during its evaluation
	Matches(Content) (bool, error)
}

// ConditionConstant keeps returning the same value
type ConditionConstant struct {
	// value to return
	value bool
}

// Matches returns the constant value, no error
func (c ConditionConstant) Matches(p Content) (bool, error) {
	return c.value, nil
}

// Signature for constant parameters is to accept anything
func (c ConditionConstant) Signature() FormalParameters {
	return NewMostPermissiveFormalParameters()
}

// NewConditionConstant returns a constant condition with given value
func NewConditionConstant(value bool) Condition {
	return ConditionConstant{value: value}
}

// conditionChilds returns childs of a condition if it is a conditional operator, or nil, false.
// First result is the childs, second result is true only for conditional operators
func conditionChilds(condition Condition) ([]Condition, bool) {
	if v, ok := condition.(ConditionOr); ok {
		return v.operands, true
	} else if v, ok := condition.(ConditionNot); ok {
		return []Condition{v.operand}, true
	} else if v, ok := condition.(ConditionAnd); ok {
		return v.operands, true
	}

	return nil, false
}

// conditionalEvaluationNode is a technical class to use in walkthrough and deal with lazy users reusing condition.
// If we use conditions directly in walkthroughs, duplicates would create cycles preventing the walkthrough.
// If, each time we see a condition, we create a new conditionalEvaluationNode, no more cycle issue
type conditionalEvaluationNode struct {
	// id should be unique (globally) to prevent two conditions to have same id no matter if they are equals as conditions
	id string
	//  decorated condition (two nodes may then have a common condition )
	condition Condition
}

// Matches decorates en inner condition
func (en conditionalEvaluationNode) Matches(p Content) (bool, error) {
	if p == nil || en.condition == nil {
		return false, nil
	}

	return en.condition.Matches(p)
}

// Signature just decorates en inner condition
func (en conditionalEvaluationNode) Signature(p Content) FormalParameters {
	return en.condition.Signature()
}

// conditionOperatorEvaluate evaluates a conditional operator on given operands.
// If condition is not a conditional operator, then it returns false.
// Note that if operands contain no value, we return false
func conditionOperatorEvaluate(condition Condition, operands []bool) bool {
	if _, ok := condition.(ConditionOr); ok {
		return slices.Contains(operands, true)
	} else if _, ok := condition.(ConditionAnd); ok {
		if len(operands) == 0 {
			return false
		}

		return !slices.Contains(operands, false)
	} else if _, ok := condition.(ConditionNot); ok {
		if len(operands) != 1 {
			return false

		}

		return !operands[0]
	} else {
		return false
	}
}

// conditionTreeMapReduce is a tree map reduce with nodes as conditions.
// It reads the tree from top to bottom, mapping each leaf using a mappe.
// Then, it goes back and reduces each node using its childs values.
// Parameters are the type of the result as T,
// condition as the root, parameters as leafs parameters, mapper and reducer.
func conditionTreeMapReduce[T any](condition Condition, parameters Content,
	mapper func(condition Condition, parameters Content) (T, error),
	reducer func(condition Condition, values []T) T,
) (T, error) {
	var empty T
	if condition == nil {
		return empty, nil
	} else if parameters == nil {
		return empty, nil
	}

	// structure links a condition to its id
	structure := make(map[string]conditionalEvaluationNode)
	// parents are necessary to find unique parent of a node and then keep going up
	parents := make(map[string]conditionalEvaluationNode)
	// childs are useful to list all brothers of current nodes in the walkthrough via parents and then childs.
	// It would be possible to go up, find the node, and then find condition, and then rebuild childs.
	// Way too long, we allocate a childs map and deal with a map instead of a super memory efficient algorithm
	childs := make(map[string][]string)
	// mapping links resolved nodes with their values.
	// A node is said to be resolved if we evaluated itself and all its descendants
	mapping := make(map[string]T)
	// rootId is the id of the root, and then mapping[rootId] is the result
	rootId := NewId()

	// first, initialize structure with current node.
	// Then, once popping a node, invariant is that its id is known
	rootValue := conditionalEvaluationNode{id: rootId, condition: condition}
	structure[rootId] = rootValue

	// queue for BFS
	queue := []conditionalEvaluationNode{rootValue}

	// STEP ONE:
	// LINK EACH NODE TO AN ID
	// REGISTER TREE STRUCTURE within parents and childs.
	// WHEN NODE IS A LEAF, EVALUATE IT and then find leafs
	for len(queue) != 0 {
		element := queue[0]
		queue = queue[1:]
		if element.condition == nil {
			mapping[element.id] = empty
		}

		// currentId of the element.
		// Because root was set and invariant to include it for childs, always set
		currentId := element.id

		// find childs if any
		currentChilds, isComposed := conditionChilds(element.condition)
		switch {
		case !isComposed:
			if value, err := mapper(element.condition, parameters); err != nil {
				return empty, err
			} else {
				mapping[currentId] = value
			}
		case len(currentChilds) == 0:
			mapping[currentId] = empty
		case len(currentChilds) != 0:
			// for each child, register child <-> current link and add it for later processing
			for _, child := range currentChilds {
				// allocate id and structure link
				childId := NewId()
				childNode := conditionalEvaluationNode{id: childId, condition: child}
				structure[childId] = childNode
				// register link to parent
				parents[childId] = element
				// register link to childs
				existingChilds := childs[currentId]
				existingChilds = SliceDeduplicate(append(existingChilds, childId))
				childs[currentId] = existingChilds
				// keep going
				queue = append(queue, childNode)
			}
		}
	}

	// STEP TWO:
	// STARTING FROM LEAFS, TRY EACH TIME TO FILL PARENTS OF FILLED VALUES
	// Starting from the leafs:
	// until root is known
	// * for each parent, if childs are all set, then set parent
	// So each outer loop sets one parent at least

	// currents are currently explored nodes
	currents := make(map[string]T)
	// nexts are the current elements to process (next run nodes are currents parents)
	nexts := make(map[string]T)
	// initially, they are leafs of the condition graph
	maps.Copy(currents, mapping)

	// as long as we ignore the value of the root
	for _, found := mapping[rootId]; !found; {
		for currentId, currentValue := range currents {
			parent, hasParent := parents[currentId]
			if !hasParent {
				// found root, we know the value, so just return
				return currentValue, nil
			}

			// if we already know the parent, skip
			if value, found := mapping[parent.id]; found {
				nexts[parent.id] = value
				continue
			}

			// childs[parent.id] contains the childs of the parent, that is the brothers and sisters (siblings)
			// If all siblings are evaluated (or evaluable), then fill parent value
			// values of the siblings
			var values []T
			// number of elements
			counter := 0
			// true if all siblings are nil (including current)
			allNilConditions := true
			for _, siblingId := range childs[parent.id] {
				sibling := structure[siblingId]
				// ensure that we set allNilConditions the first time we see a non nil condition
				if allNilConditions && sibling.condition != nil {
					// we have a non nil condition, so update allNilConditions
					allNilConditions = false
				}

				// add the value if we already evaluated the sibling
				if value, found := mapping[siblingId]; found {
					values = append(values, value)
				}

				// counter is changed no matter the node condition
				counter++
			}

			// if we know all the values, then we may evaluate the parent.
			// Otherwise, it is too early.
			// We know all the values when:
			// EITHER all childs conditions of parent are nil => parent is evaluated to default parent behavior
			// OR all values for siblings are known, and then we may evaluate the parent
			if allNilConditions {
				mapping[parent.id] = empty
				nexts[parent.id] = empty
			} else if counter == len(values) {
				// we have everything to evaluate the parent
				result := reducer(parent.condition, values)
				mapping[parent.id] = result
				nexts[parent.id] = result
			}
		}

		// we then explored all current nodes
		if len(nexts) == 0 {
			// should be impossible
			return empty, errors.New("no progression in walkthrough")
		} else if result, found := nexts[rootId]; found {
			return result, nil
		}

		// clear currents
		for k := range currents {
			delete(currents, k)
		}

		// currents is the new nexts.
		// It means we try to go higher
		maps.Copy(currents, nexts)
		// clean nexts to ensure only new values will be added
		for k := range nexts {
			delete(nexts, k)
		}
	}

	if value, found := mapping[rootId]; !found {
		return empty, errors.New("algorithm did not end")
	} else {
		return value, nil
	}
}

// conditionLeafEvaluate evaluates a leaf, returns false for nil conten or no condition
func conditionLeafEvaluate(condition Condition, content Content) (bool, error) {
	if content == nil {
		return false, nil
	} else if condition == nil {
		return false, nil
	} else if _, isOperator := conditionChilds(condition); isOperator {
		return false, nil
	} else {
		return condition.Matches(content)
	}
}

// conditionTreeEvaluate evaluates a logical tree starting at condition with those parameters (as content).
// If condition or parameters is nil, then it returns false.
// Implementation is based on an iterative walkthrough.
// Reason is a recursive walkthrough would put pressure on the stack.
// It would be easier, of course, but let us preserve the running environment.
func conditionTreeEvaluate(condition Condition, parameters Content) (bool, error) {
	return conditionTreeMapReduce(condition, parameters, conditionLeafEvaluate, conditionOperatorEvaluate)
}

// conditionTreeCombineSignatures parses the conditions tree and get the parameters it accepts as a whole.
// For instance, if a needs variable x and b needs variable y, then result should be expected variables x and y.
// It returns the necessary formal parameters or an error.
func conditionTreeCombineSignatures(condition Condition) FormalParameters {
	conditionLeafSignature := func(condition Condition, content Content) (FormalParameters, error) {
		if condition == nil {
			return FormalParameters{}, nil
		} else if _, isOperator := conditionChilds(condition); isOperator {
			return FormalParameters{}, nil
		} else {
			return condition.Signature(), nil
		}
	}

	conditionParametersReduce := func(condition Condition, values []FormalParameters) FormalParameters {
		return parametersCombine(values)
	}

	dummyValue := NewContent[Modelable](nil)
	result, _ := conditionTreeMapReduce(condition, dummyValue, conditionLeafSignature, conditionParametersReduce)
	return result
}

// cleanConditionsOperands returns the conditions with nil excluded.
// If there is no element or conditions are only nil, return nil
func cleanConditionsOperands(conditions []Condition) []Condition {
	if len(conditions) == 0 {
		return nil
	}

	var result []Condition
	for _, condition := range conditions {
		if condition != nil {
			result = append(result, condition)
		}
	}

	return result
}

// ConditionAnd is true if it has at least one operand and if all its operands matches their parameters
type ConditionAnd struct {
	operands []Condition
}

// Matches returns true if there is at least one operand condition and all operands matches the parameters
func (a ConditionAnd) Matches(p Content) (bool, error) {
	return conditionTreeEvaluate(a, p)
}

// Signature returns the formal parameters inferred from its descendants
func (a ConditionAnd) Signature() FormalParameters {
	return conditionTreeCombineSignatures(a)
}

// ConditionOr is true if it has at least one operand and if any of them matches its parameters
type ConditionOr struct {
	operands []Condition
}

// Matches returns true if at least one condition in operands matches the content
func (o ConditionOr) Matches(p Content) (bool, error) {
	return conditionTreeEvaluate(o, p)
}

// Signature returns the formal parameters inferred from its descendants
func (o ConditionOr) Signature() FormalParameters {
	return conditionTreeCombineSignatures(o)
}

// ConditionNot negates a condition
type ConditionNot struct {
	// operand is the condition to negate
	operand Condition
}

// Matches returns false for no operand, and not (the result of operand applied to parameters) otherwise
func (n ConditionNot) Matches(p Content) (bool, error) {
	return conditionTreeEvaluate(n, p)
}

// Signature returns the formal parameters inferred from its descendants
func (n ConditionNot) Signature() FormalParameters {
	return conditionTreeCombineSignatures(n)
}

// NewConditionOr returns OR(conditions) as a condition
func NewConditionOr(conditions []Condition) Condition {
	return ConditionOr{operands: cleanConditionsOperands(conditions)}
}

// NewConditionAnd returns AND(conditions) as a condition
func NewConditionAnd(conditions []Condition) Condition {
	return ConditionAnd{operands: cleanConditionsOperands(conditions)}
}

// NewConditionNot returns NOT(condition) as a condition
func NewConditionNot(condition Condition) Condition {
	return ConditionNot{operand: condition}
}
