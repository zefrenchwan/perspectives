package commons

import (
	"errors"
	"slices"
)

// RoleSubject specifies the subject role
const RoleSubject = "subject"

// RoleObject specifies the object role
const RoleObject = "object"

// Linkable should be as simple as possible.
type Linkable interface{}

// LinkableEquality defines when two linkables are the same
type LinkableEquality func(a, b Linkable) bool

// LinkableSame is default implementation for same linkables.
// It returns true if a and b may be conidered as equivalent
func LinkableSame(a, b Linkable) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	} else if label, ok := a.(LinkLabel); ok {
		if otherLabel, otherOk := b.(LinkLabel); otherOk {
			return label.Name() == otherLabel.Name()
		} else {
			return false
		}
	} else if variable, ok := a.(LinkVariable); ok {
		if otherVariable, otherOk := b.(LinkVariable); otherOk {
			return variable.Name() == otherVariable.Name()
		} else {
			return false
		}
	} else if idBased, ok := a.(Identifiable); ok {
		if otherIdBased, otherOk := b.(Identifiable); otherOk {
			return idBased.Id() == otherIdBased.Id()
		} else {
			return false
		}
	} else {
		return false
	}
}

// LinkLabel is a label to qualify other elements.
// For instance, human, nice, etc.
type LinkLabel struct {
	// label is the value of the label.
	label string
}

// Name returns the name of that label
func (l LinkLabel) Name() string {
	return l.label
}

// NewLabel builds a new linkable label
func NewLabel(name string) LinkLabel {
	return LinkLabel{label: name}
}

// LinkVariable is a placeholder for values (such as objects and links)
type LinkVariable interface {
	// Name returns the name for that variable
	Name() string
	// Accepts returns true if replacement is valid
	Accepts(Linkable) bool
}

// simpleLinkVariable uses a name and an acceptor and we are done
type simpleLinkVariable struct {
	// name of the variable
	name string
	// acceptance for links to replace (nil means false)
	acceptance func(Linkable) bool
}

// Name returns the name for that variable
func (s simpleLinkVariable) Name() string {
	return s.name
}

// Accepts returns true if l could replace that variable
func (s simpleLinkVariable) Accepts(l Linkable) bool {
	return s.acceptance != nil && s.acceptance(l)
}

// linkableRefuse returns false
func linkableRefuse(link Linkable) bool {
	return false
}

// NewLinkVariable builds a new link variable that accepts based on that predicate
func NewLinkVariable(name string, predicate func(Linkable) bool) LinkVariable {
	if predicate == nil {
		return simpleLinkVariable{name, linkableRefuse}
	}

	return simpleLinkVariable{name, predicate}
}

// NewLinkVariableForObject returns a new variable accepting any NOT NIL object
func NewLinkVariableForObject(name string) LinkVariable {
	return NewLinkVariable(name, func(l Linkable) bool {
		if l == nil {
			return false
		}

		o, ok := l.(ModelObject)
		return ok && o != nil
	})
}

// NewLinkVariableForLink creates a variable that accepts anynon nil link
func NewLinkVariableForLink(name string) LinkVariable {
	return NewLinkVariable(name, func(l Linkable) bool {
		if l == nil {
			return false
		}

		link, ok := l.(Link)
		return ok && link != nil
	})
}

// Link is a constant relation over instances of linkables.
// Link is also Linkable, so it may be used in links.
type Link interface {
	// Links have an id, new each time we build one link
	Identifiable
	// Links are information about elements, they are then part of a model
	Modelable
	// links are composable due to this: a link is linkable
	Linkable
	// IsEmpty returns true if link should be empty and then not being used
	IsEmpty() bool
	// Name returns the name of the link.
	// It is usually a verb or a noun.
	// For instance, knows, couple, etc.
	Name() string
	// Roles returns all the roles set for that link
	Roles() []string
	// Operands returns the roles and related values for that link
	Operands() map[string]Linkable
	// Get returns, if any, the elements for that name.
	// First result is the values if any, second is true if value was found
	Get(role string) (Linkable, bool)
}

// TemporalLink decorates a link over a period.
// It should implement both link and Temporal.
type TemporalLink struct {
	// id (cannot use the same as the orignal one because their properties are different)
	id string
	// decorated link
	Link
	// activity of the link
	period Period
}

// NewTemporalLink decorates a link true for given duration
func NewTemporalLink(duration Period, value Link) *TemporalLink {
	result := new(TemporalLink)
	result.Link = value
	result.id = NewId()
	result.period = duration
	return result
}

// Id() returns the id of the temporal link, not the same as underlying
func (t *TemporalLink) Id() string {
	return t.id
}

// ActivePeriod is the duration which the link is true
func (t *TemporalLink) ActivePeriod() Period {
	if t == nil {
		return NewEmptyPeriod()
	}

	return t.period
}

// SetActivePeriod forces active period
func (t *TemporalLink) SetActivePeriod(period Period) {
	if t != nil {
		t.period = period
	}
}

// simpleLinkNode decorates a value to ensure that it has an unique id (even for similar values)
type simpleLinkNode struct {
	// id should be unique
	id string
	// node is the decorated value
	node Linkable
}

// simpleLink implements links as its canonical implementation
type simpleLink struct {
	// id of the link
	id string
	// name of the link
	name string
	// values map roles to decorated linkable values
	values map[string]simpleLinkNode
}

// Id returns the link unique id
func (s simpleLink) Id() string {
	return s.id
}

// IsEmpty tests if link is empty (no name or no value)
func (s simpleLink) IsEmpty() bool {
	return len(s.values) == 0 || s.name == ""
}

// GetType acts the fact that a link is a model link
func (s simpleLink) GetType() ModelableType {
	return TypeLink
}

// Name returns the name of the link.
func (s simpleLink) Name() string {
	return s.name
}

// Roles returns all the roles set for that link
func (s simpleLink) Roles() []string {
	var result []string
	for role := range s.values {
		result = append(result, role)
	}

	return result
}

// Operands returns the roles and related values for that link.
// To avoid side effects, we return a copy (not direct access to values)
func (s simpleLink) Operands() map[string]Linkable {
	result := make(map[string]Linkable)
	for role, container := range s.values {
		result[role] = container.node
	}

	return result
}

// Get returns, if any, the elements for that name.
func (s simpleLink) Get(role string) (Linkable, bool) {
	if result, found := s.values[role]; found {
		return result.node, true
	}

	return nil, false
}

// NewLink builds a new link, or raises an error if link would be malformed.
// A valid link is not empty: non empty name and at least one value.
// Of course, a "creative" user may create a link with " " name, but it is discouraged
func NewLink(name string, values map[string]Linkable) (Link, error) {
	var empty string
	if len(values) == 0 {
		return nil, errors.New("no value for roles")
	} else if name == empty {
		return nil, errors.New("no name for link")
	}

	var result simpleLink
	result.id = NewId()
	result.name = name
	result.values = make(map[string]simpleLinkNode)
	for role, value := range values {
		decorated := simpleLinkNode{id: NewId(), node: value}
		result.values[role] = decorated
	}

	return result, nil
}

// NewSimpleLink is just creating a link of given name with content equals to "subject" => subject, "object": object
func NewSimpleLink(name string, subject Linkable, object Linkable) (Link, error) {
	return NewLink(name, map[string]Linkable{RoleSubject: subject, RoleObject: object})
}

// linkNodeMapping is a container to put any relevant node data in a link or a leaf (union struct).
// It contains any node information (old values to read, new values to map)
type linkNodeMapping struct {
	// invariantId kept before and after mapping
	invariantId string
	// originalValue is read from the original link
	originalValue Linkable
	// newValue is mapped node (after)
	newValue Linkable
	// changed is true if newValue and originalValue are different
	changed bool
	// role is the role of the value from parent
	role string
}

// initializeLinkNodeMapping gets a linkable and initialiaze mapping node
func initializeLinkNodeMapping(originalValue Linkable) linkNodeMapping {
	var result linkNodeMapping
	result.invariantId = NewId()
	result.originalValue = originalValue
	return result
}

// LinkMapLeafs maps leafs and returns a new link with same structure, just different leafs.
// baseLink is the link to map, mapper maps the leaf to another and returns true to flag a change.
func LinkMapLeafs(baseLink Link, mapper func(Linkable) (Linkable, bool)) (Link, error) {
	if baseLink == nil || baseLink.IsEmpty() {
		return baseLink, nil
	}

	// mapping contains each node to map (old and new values)
	mapping := make(map[string]linkNodeMapping)
	// structure links a parent (link) to all its child links
	structure := make(map[string][]string)
	// parents link a child to its parent (structure does the opposite)
	parents := make(map[string]string)
	// starting point, what is the root
	root := initializeLinkNodeMapping(baseLink)
	// register root, to ensure invariant
	mapping[root.invariantId] = root
	// BFS walkthrough
	queue := []linkNodeMapping{root}

	// for each node
	for len(queue) != 0 {
		// pop node
		current := queue[0]
		queue = queue[1:]
		var currentLink Link
		if l, ok := current.originalValue.(Link); l != nil && ok {
			currentLink = l
		} else {
			// should not happen, but we just ignore a non link
			continue
		}

		currentId := current.invariantId
		// for each child, we already have node in the map, so we add links childs
		for role, child := range currentLink.Operands() {
			// child is an operand, may not be a link
			if child == nil {
				continue
			} else {
				mappedChild := initializeLinkNodeMapping(child)
				mappedChild.role = role
				mapping[mappedChild.invariantId] = mappedChild
				// link current link to that child
				values := structure[currentId]
				values = append(values, mappedChild.invariantId)
				structure[currentId] = SliceDeduplicate(values)
				parents[mappedChild.invariantId] = currentId
				// if child is a link, then add it to processing
				if childLink, ok := child.(Link); childLink != nil && ok {
					queue = append(queue, mappedChild)
				}
			}
		}
	}

	// We now know the full link structure
	// Then, find leafs by taking the difference of all childs and the heads
	var leafs, heads []string
	for head, childs := range structure {
		leafs = append(leafs, childs...)
		heads = append(heads, head)
	}

	// leafs are values that has no child
	leafs = SlicesFilter(SliceDeduplicate(leafs), func(s string) bool { return !slices.Contains(heads, s) })

	// init the bottom up walkthrough : map leafs
	var processed []string
	for _, leafId := range leafs {
		leafContent := mapping[leafId]
		leafMapping, mapped := mapper(leafContent.originalValue)
		leafContent.changed = mapped
		leafContent.newValue = leafMapping
		mapping[leafId] = leafContent
		processed = append(processed, leafId)
	}

	for len(processed) != 0 {
		var nexts []string
		for _, currentId := range processed {
			// current exists for sure, but it may not be a link.
			// Find parent. If not, it means the root
			if parent, found := parents[currentId]; found {
				// parent may have been dealt with already
				if mapping[parent].newValue != nil {
					nexts = append(nexts, parent)
				} else {
					// first time we see the parent.
					// We want to process parent if all its childs are processed
					allSiblingsDone := true
					// siblingChanged is true if at lease one sibling changed (to change parent id)
					siblingChanged := false
					// allSiblings contains the id of all the childs of that parent
					var allSiblings []string
					// for each sibling (then all childs of the parent)
					for _, sibling := range structure[parent] {
						allSiblings = append(allSiblings, sibling)
						siblingValue := mapping[sibling]
						// if a value is not done, then parent may not be processed
						if siblingValue.newValue == nil {
							allSiblingsDone = false
						}
						// if a value changed, then parent id should be different
						if siblingValue.changed {
							siblingChanged = true
						}
					}

					// if all siblings are done, then we may process the parent
					if allSiblingsDone {
						parentValue := mapping[parent]
						// parent is a link for sure
						parentLink, _ := parentValue.originalValue.(Link)
						if siblingChanged {
							// copy the link values from the original link, change its id
							parentValue.changed = true
							roles := make(map[string]Linkable)
							for _, child := range structure[parent] {
								childValue := mapping[child]
								roles[childValue.role] = childValue.newValue
							}

							if newLink, err := NewLink(parentLink.Name(), roles); err != nil {
								return nil, err
							} else {
								parentValue.newValue = newLink
								parentValue.changed = true
								mapping[parent] = parentValue
							}
						} else {
							parentValue.changed = false
							parentValue.newValue = parentValue.originalValue
							mapping[parent] = parentValue
						}

						// parent may be processed on the next run
						nexts = append(nexts, parent)
						// we do not need the siblings values, just clean them
						for _, sibling := range allSiblings {
							delete(mapping, sibling)
						}
					}
				}
			}
		}

		// new processed is current next
		processed = nil
		nexts = SliceDeduplicate(nexts)
		processed = make([]string, len(nexts))
		copy(processed, nexts)
	}

	// at this point, there is no upper level, so we basically read the root mapping
	result := mapping[root.invariantId]
	if result.newValue == nil {
		return nil, errors.New("no mapping for root")
	} else if rootLink, ok := result.newValue.(Link); !ok || rootLink == nil {
		return nil, errors.New("no root link")
	} else {
		return rootLink, nil
	}
}

// LinkAcceptsInstantiation returns a possible variables mapping (if any) to use from pattern to baseLink.
// For instance, Knows(Paul, Jules) is an instantiation of Knows(Paul, X) with X => Jules.
// Signature is:
// the base link (the possible instantiation) to compare to the pattern (may be a variable)
// and the matchingsEqualsFn to test equality on variables matches.
// Then, if there is a mapping that matches, first result is variable name => Linkable to replace.
// Otherwise, it returns nil, false.
func LinkAcceptsInstantiation(baseLink Link, pattern Linkable, matchingsEqualsFn LinkableEquality) (map[string]Linkable, bool) {
	variablesInstantiation := make(map[string]Linkable)
	var patternLink Link
	// either both parameters are links,
	// or second one is a variable (and we test) or we end it as false
	if baseLink == nil && pattern == nil {
		return nil, true
	} else if baseLink == nil || pattern == nil {
		return nil, false
	} else if v, ok := pattern.(LinkVariable); ok {
		if v == nil {
			return nil, false
		} else if v.Accepts(baseLink) {
			variablesInstantiation[v.Name()] = baseLink
			return variablesInstantiation, true
		} else {
			return nil, false
		}
	} else if l, ok := pattern.(Link); !ok {
		return nil, false
	} else {
		patternLink = l
	}

	// go for a walk through each graph, stop as soon as structures differ.
	baseQueue := []Link{baseLink}
	patternQueue := []Link{patternLink}

	// for each node (BFS)
	for len(baseQueue) != 0 {
		if len(baseQueue) != len(patternQueue) {
			return nil, false
		}

		// both are links as an invariant
		currentLink := baseQueue[0]
		baseQueue = baseQueue[1:]
		referenceLink := patternQueue[0]
		patternQueue = patternQueue[1:]

		if currentLink.Name() != referenceLink.Name() {
			return nil, false
		}

		currentOperands := currentLink.Operands()
		referenceOperands := referenceLink.Operands()
		if len(currentOperands) != len(referenceOperands) {
			return nil, false
		}

		// go through each role, and compare each role accordingly.
		for role, baseValue := range currentOperands {
			referenceValue, found := referenceOperands[role]
			if !found {
				return nil, false
			} else if referenceVariable, ok := referenceValue.(LinkVariable); ok {
				if !referenceVariable.Accepts(baseValue) {
					return nil, false
				} else {
					if mappings, found := variablesInstantiation[referenceVariable.Name()]; !found {
						variablesInstantiation[referenceVariable.Name()] = baseValue
					} else if matchingsEqualsFn(baseValue, mappings) {
						variablesInstantiation[referenceVariable.Name()] = baseValue
					} else {
						return nil, false
					}
				}
			} else if !LinkableSame(baseValue, referenceValue) {
				return nil, false
			}

			// if base value is a link (so is reference value), then keep walking through
			if baseLink, ok := baseValue.(Link); ok {
				if otherLink, otherOk := referenceValue.(Link); otherOk {
					// both are links, equivalent (same) with each other
					baseQueue = append(baseQueue, baseLink)
					patternQueue = append(patternQueue, otherLink)
				} else {
					return nil, false
				}
			}
		}
	}

	return variablesInstantiation, len(baseQueue) == len(patternQueue)
}

// LinkSetVariables just sets variables within a link, rest is unchanged.
// In particular, it returns the same link is there is no variable to replace.
func LinkSetVariables(baseLink Link, values map[string]Linkable) (Link, error) {
	if baseLink == nil || len(values) == 0 {
		return nil, nil
	}

	mapper := func(l Linkable) (Linkable, bool) {
		if l == nil {
			return nil, false
		}

		if v, ok := l.(LinkVariable); !ok {
			return l, false
		} else if ok && v != nil {
			if replacer, found := values[v.Name()]; !found {
				return l, false
			} else {
				return replacer, true
			}
		}

		return l, false
	}

	return LinkMapLeafs(baseLink, mapper)
}

// LinkFindAllMatching finds all nodes matching a condition.
// Walkthrough is a BFS, but the map roles implementation does not guarantee a strict deterministic order
func LinkFindAllMatching(base Link, condition func(Linkable) bool) []Linkable {
	if base == nil || condition == nil {
		return nil
	}

	// result contains all linkables matching condition (link, leafs, etc)
	var result []Linkable

	queue := []Link{base}
	for len(queue) != 0 {
		current := queue[0]
		queue = queue[1:]

		// link test
		if condition(current) {
			result = append(result, current)
		}

		for _, operand := range current.Operands() {
			if child, ok := operand.(Link); ok && child != nil {
				// no need to test child, will be done after
				queue = append(queue, child)
			} else if condition(operand) {
				// operand is not a link, so test but not queue
				result = append(result, operand)
			}
		}
	}

	return result
}
