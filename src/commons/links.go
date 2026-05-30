package commons

import (
	"reflect"
	"slices"
)

// Link is a relation between elements.
// Think of it as a node in a graph that connects to other elements (Traits, Variables, or nested Links).
type Link interface {
	Element // Inherits Same(Element) and DeclaringClass()

	// --- READ OPERATIONS : State queries ---

	// Name returns the name of the link. For instance, in Loves(subject=John, object=Pizza) => "Loves"
	Name() string
	// Validity returns the time period during which this link is considered true/active.
	Validity() Period
	// Operands returns all operand keys (the names of the relationships, e.g., "subject", "object") deterministically sorted.
	Operands() []string
	// Operand returns the slice of Elements associated with the given operand name.
	// It returns a boolean indicating if the operand key exists.
	Operand(name string) ([]Element, bool)

	// --- FUNCTIONAL MUTATIONS : Copy-on-write operations ---
	// Since Links are immutable, these methods never modify the current instance.
	// They return a newly allocated Link with the requested changes.

	// WithValidity returns a copy of the link with the new validity period.
	WithValidity(p Period) Link

	// WithOperand returns a copy of the link with the given operand forced (overwrites previous values).
	WithOperand(name string, operands []Element) Link

	// WithAppended returns a copy of the link with a new element added to the specified operand.
	WithAppended(name string, operand Element) Link

	// Without returns a copy of the link, filtering out operand values that match the condition.
	// For instance, Loves(John, [Pizza, Salad]) => Without("object", isPizza) => Loves(John, [Salad])
	Without(name string, op func(linkable Element) bool) Link

	// ReplaceVariable traverses the graph and replaces occurrences of a specific Variable with a new Element.
	ReplaceVariable(variable Variable, value Element) Link
}

// baseLink is the immutable implementation of the Link interface.
// Explanation:  Immutability provides native thread-safety. We don't need sync.RWMutex
// because once created, a baseLink's internal state can never be altered.
type baseLink struct {
	name     string
	validity Period
	operands map[string][]Element
}

// NewLink creates a new immutable Link instance.
func NewLink(name string, validity Period) Link {
	return &baseLink{
		name:     name,
		validity: validity,
		operands: make(map[string][]Element),
	}
}

// ============================================================================
// READ METHODS
// ============================================================================

func (l *baseLink) DeclaringClass() Class {
	return CLASS_LINK
}

func (l *baseLink) Name() string {
	if l == nil {
		return ""
	}
	return l.name
}

func (l *baseLink) Validity() Period {
	if l == nil {
		return NewEmptyPeriod()
	}
	return l.validity
}

// Operands returns a deterministically sorted list of operand keys.
// Explanation:  Sorting is crucial here. If we need to serialize the graph or compare
// two links, iterating over map keys without sorting leads to random orders in Go,
// which would break equality algorithms.
func (l *baseLink) Operands() []string {
	if l == nil || l.operands == nil {
		return nil
	}
	keys := make([]string, 0, len(l.operands))
	for k := range l.operands {
		keys = append(keys, k)
	}
	slices.Sort(keys) // Guarantee stable order
	return keys
}

// Operand returns the elements associated with a specific operand name.
// Explanation:  It returns a defensive copy (SliceCopy) to prevent external code from
// mutating the internal slice. If we returned the raw slice, a user could do:
// `ops, _ := myLink.Operand("key"); ops[0] = somethingElse`, destroying immutability.
func (l *baseLink) Operand(name string) ([]Element, bool) {
	if l == nil || l.operands == nil {
		return nil, false
	}
	vals, exists := l.operands[name]
	if !exists {
		return nil, false
	}
	return SliceCopy(vals), true // Defensive copy
}

// ============================================================================
// MUTATION METHODS (FUNCTIONAL / COPY-ON-WRITE)
// ============================================================================

// copyMap is an internal utility to shallow-copy the operands map.
// Explanation:  This implements "Structural Sharing". We duplicate the map itself,
// but the underlying slices (and their elements) are shared in memory.
// This saves significant CPU and RAM when mutating large graphs.
func (l *baseLink) copyMap() map[string][]Element {
	res := make(map[string][]Element, len(l.operands))
	for k, v := range l.operands {
		res[k] = v // Slices are referenced, not deep-copied here
	}
	return res
}

func (l *baseLink) WithValidity(p Period) Link {
	if l == nil {
		return nil
	}
	return &baseLink{
		name:     l.name,
		validity: p,
		operands: l.operands, // Safe to share the map (read-only usage in new instance)
	}
}

func (l *baseLink) WithOperand(name string, operands []Element) Link {
	if l == nil {
		return nil
	}
	newOps := l.copyMap()
	if len(operands) == 0 {
		delete(newOps, name) // Clean up empty keys to keep the map tidy
	} else {
		newOps[name] = SliceCopy(operands) // Prevent external slice mutation
	}

	return &baseLink{
		name:     l.name,
		validity: l.validity,
		operands: newOps,
	}
}

func (l *baseLink) WithAppended(name string, operand Element) Link {
	if l == nil || operand == nil {
		return l
	}
	newOps := l.copyMap()
	oldSlice := newOps[name]

	// Allocate a new slice with exactly the required capacity to avoid hidden re-allocations.
	newSlice := make([]Element, len(oldSlice), len(oldSlice)+1)
	copy(newSlice, oldSlice)
	newSlice = append(newSlice, operand)

	newOps[name] = newSlice

	return &baseLink{
		name:     l.name,
		validity: l.validity,
		operands: newOps,
	}
}

func (l *baseLink) Without(name string, op func(linkable Element) bool) Link {
	if l == nil || l.operands == nil {
		return l
	}
	oldSlice, exists := l.operands[name]
	if !exists {
		return l // Fast exit: no changes needed, reuse current instance
	}

	var newSlice []Element
	for _, el := range oldSlice {
		// If op() returns false, it means we want to KEEP the element
		if !op(el) {
			newSlice = append(newSlice, el)
		}
	}

	// Optimization: If nothing was removed, return the current instance.
	// This avoids useless memory allocations.
	if len(newSlice) == len(oldSlice) {
		return l
	}

	newOps := l.copyMap()
	if len(newSlice) == 0 {
		delete(newOps, name)
	} else {
		newOps[name] = newSlice
	}

	return &baseLink{
		name:     l.name,
		validity: l.validity,
		operands: newOps,
	}
}

// ============================================================================
// NON-RECURSIVE GRAPH ALGORITHMS
// ============================================================================

// Same checks for deep equality between two Links.
// Explanation:  This uses an iterative Breadth-First Search (BFS) using a Queue.
// Why BFS? Because if the top-level validity or names differ, we exit immediately
// (early return) without wasting time exploring deep sub-graphs.
// We also track visited nodes to prevent Infinite Loops if the user created a cycle.
func (l *baseLink) Same(other Element) bool {
	if l == nil && other == nil {
		return true
	}
	if l == nil || other == nil {
		return false
	}

	type pair struct {
		a Element
		b Element
	}

	queue := []pair{{l, other}}

	// Cycle detection: track pairs of memory addresses we have already compared.
	// If we see the same pair again, we can assume they are equal in this path.
	visited := make(map[[2]uintptr]bool)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		a, b := curr.a, curr.b

		if a == nil && b == nil {
			continue
		}
		if a == nil || b == nil {
			return false
		}
		if a.DeclaringClass() != b.DeclaringClass() {
			return false
		}

		// Check if we are comparing two links (which can contain cycles)
		if IsElementDeclaredInstance(a, CLASS_LINK) {
			linkA, _ := a.(Link)
			linkB, _ := b.(Link)

			// Get pointer addresses for cycle detection
			ptrA := reflect.ValueOf(linkA).Pointer()
			ptrB := reflect.ValueOf(linkB).Pointer()
			memPair := [2]uintptr{ptrA, ptrB}

			if visited[memPair] {
				continue // We have already verified this exact pair, break the cycle.
			}
			visited[memPair] = true

			// 1. Fast checks: Name and Validity
			if linkA.Name() != linkB.Name() {
				return false
			}
			if !linkA.Validity().Equals(linkB.Validity()) {
				return false
			}

			// 2. Structural checks: Operand keys
			opsA := linkA.Operands()
			opsB := linkB.Operands()
			if len(opsA) != len(opsB) {
				return false
			}

			// 3. Queue children for the next BFS iterations
			for _, opName := range opsA {
				elsA, _ := linkA.Operand(opName)
				elsB, okB := linkB.Operand(opName)

				if !okB || len(elsA) != len(elsB) {
					return false
				}
				for i := range elsA {
					queue = append(queue, pair{elsA[i], elsB[i]})
				}
			}
		} else {
			// Base case: For instances, traits, or variables, defer to their native Same()
			if !a.Same(b) {
				return false
			}
		}
	}
	return true
}

// ReplaceVariable finds a target Variable and replaces it with a new Element.
// Explanation:  Since Links are immutable, substituting a deep leaf requires rebuilding all
// parent nodes up to the root (Bottom-Up reconstruction).
// We use an iterative Depth-First Search (DFS) with a custom Call Stack (LIFO)
// to prevent Stack Overflows on huge graphs.
// It also implements Structural Sharing and Cycle Prevention.
func (l *baseLink) ReplaceVariable(variable Variable, value Element) Link {
	if l == nil {
		return nil
	}

	// Frame simulates a function call context in the Go runtime.
	type frame struct {
		link       Link                 // The Link currently being processed
		keys       []string             // All operand keys of this link
		keyIdx     int                  // Current position in the 'keys' slice
		valIdx     int                  // Current position in the current key's elements
		newOps     map[string][]Element // The reconstructed map for the new link
		currVals   []Element            // The reconstructed elements for the current key
		hasChanges bool                 // Flags if ANY child actually changed (for Structural Sharing)
	}

	stack := []*frame{{
		link:   l,
		keys:   l.Operands(),
		newOps: make(map[string][]Element),
	}}

	// Cycle Prevention and DAG Optimization
	memo := make(map[uintptr]Link)    // Maps old Link pointer -> new processed Link
	inStack := make(map[uintptr]bool) // Tracks links currently in the processing stack

	// Mark root as currently processing
	inStack[reflect.ValueOf(l).Pointer()] = true

	var result Element

	for len(stack) > 0 {
		curr := stack[len(stack)-1]

		// Condition A: Have we finished iterating over all keys in this Link?
		if curr.keyIdx >= len(curr.keys) {
			var resultLink Link

			// OPTIMIZATION (Structural Sharing): If no variables were replaced in this branch,
			// do NOT allocate a new map/link. Just reuse the original pointer.
			if !curr.hasChanges {
				resultLink = curr.link
			} else {
				resultLink = &baseLink{
					name:     curr.link.Name(),
					validity: curr.link.Validity(),
					operands: curr.newOps,
				}
			}

			// Pop the current frame
			stack = stack[:len(stack)-1]
			ptr := reflect.ValueOf(curr.link).Pointer()

			// Update tracking maps
			delete(inStack, ptr)
			memo[ptr] = resultLink

			// If there is a parent Link waiting above us, append our newly built link to it
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.currVals = append(parent.currVals, resultLink)

				// If the child was reconstructed (changed), the parent is also considered changed.
				if resultLink != curr.link {
					parent.hasChanges = true
				}
			} else {
				// We reached the root
				result = resultLink
			}
			continue
		}

		key := curr.keys[curr.keyIdx]
		vals, _ := curr.link.Operand(key)

		// Condition B: Have we finished iterating over all elements of the current key?
		if curr.valIdx >= len(vals) {
			if len(curr.currVals) > 0 {
				curr.newOps[key] = curr.currVals
			}
			curr.currVals = nil
			curr.keyIdx++
			curr.valIdx = 0
			continue
		}

		// Condition C: Process the current Element
		elem := vals[curr.valIdx]
		curr.valIdx++

		// Scenario 1: The element is a Variable
		if elem != nil && IsElementDeclaredInstance(elem, CLASS_VARIABLE) {
			if elem.Same(&variable) {
				// BUG FIX: We MUST verify if the value respects the AllowedTypes of the Variable.
				if variable.CanBeReplacedBy(value) {
					curr.currVals = append(curr.currVals, value)
					curr.hasChanges = true
				} else {
					// Explanation:  If the constraint fails, we safely ignore the replacement
					// and keep the original variable to prevent corrupting the graph logic.
					curr.currVals = append(curr.currVals, elem)
				}
				continue
			}
		}

		// Scenario 2: The element is a nested Link (Dive deeper)
		if elem != nil && IsElementDeclaredInstance(elem, CLASS_LINK) {
			if childLink, ok := elem.(Link); ok {
				ptr := reflect.ValueOf(childLink).Pointer()

				// Check for Cycles
				if inStack[ptr] {
					// Explanation:  You cannot deep-copy a cyclic graph immutably bottom-up without infinite recursion.
					// If we detect a cycle, we break it by keeping the original reference.
					curr.currVals = append(curr.currVals, childLink)
					continue
				}

				// Check for DAG Memoization (already processed sub-graph)
				if cachedLink, exists := memo[ptr]; exists {
					curr.currVals = append(curr.currVals, cachedLink)
					if cachedLink != childLink {
						curr.hasChanges = true
					}
					continue
				}

				// Push a new frame to process this child
				inStack[ptr] = true
				stack = append(stack, &frame{
					link:   childLink,
					keys:   childLink.Operands(),
					newOps: make(map[string][]Element),
				})
				continue
			}
		}

		// Scenario 3: Terminal element (Instance, Trait, or un-matching Variable)
		curr.currVals = append(curr.currVals, elem)
	}

	return result.(Link)
}
