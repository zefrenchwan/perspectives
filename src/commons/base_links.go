package commons

import (
	"reflect"
	"slices"
)

// safePointer extracts the memory address of an interface strictly if it is a pointer.
// It returns the pointer address and true, or 0 and false if it's a value type.
// Explanation: In Go, infinite memory cycles can ONLY be formed using pointers.
// Value types are copied, thus inherently acyclic.
func safePointer(i any) (uintptr, bool) {
	if i == nil {
		return 0, false
	}
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		return v.Pointer(), true
	}
	return 0, false
}

// =============================================
// STARTING IMPLEMENTATION
// =============================================

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

// Same checks for deep equality between two Links.
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

		if IsElementDeclaredInstance(a, CLASS_LINK) {
			linkA, okA := a.(Link)
			linkB, okB := b.(Link)
			if !okA || !okB {
				return false
			}

			// --- SECURE CYCLE DETECTION ---
			ptrA, isPtrA := safePointer(linkA)
			ptrB, isPtrB := safePointer(linkB)

			if isPtrA && isPtrB {
				memPair := [2]uintptr{ptrA, ptrB}
				if visited[memPair] {
					continue // Déjà comparés, on casse le cycle
				}
				visited[memPair] = true
			}

			// 1. Fast checks
			if linkA.Name() != linkB.Name() || !linkA.Validity().Equals(linkB.Validity()) {
				return false
			}

			// 2. Structural checks
			opsA := linkA.Operands()
			opsB := linkB.Operands()
			if len(opsA) != len(opsB) {
				return false
			}

			// 3. Queue children
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
			if !a.Same(b) {
				return false
			}
		}
	}
	return true
}
