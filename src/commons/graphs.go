package commons

import (
	"slices"
	"sync"

	"github.com/google/uuid"
)

// Linkable represents an entity that can be linked in a graph.
// A link, especially, is linkable.
// It means that it is possible to use links as operands
type Linkable interface {
	Element
}

// Link represents a link between nodes in a graph.
type Link struct {
	locks    sync.RWMutex          // locks manage the locks for goroutines
	id       string                // id of the link
	name     string                // name of the link
	operands map[string][]Linkable // operands of the link
	validity Period                // moments the link is true
}

// NewLink creates a new empty link with provided name.
func NewLink(name string) *Link {
	return &Link{
		name:     name,
		id:       uuid.NewString(),
		operands: make(map[string][]Linkable),
		validity: NewFullPeriod(),
	}
}

// Id returns the id of the link.
// Name will not be unique, id is for sure
func (l *Link) Id() string {
	return l.id
}

// Name returns the name of the link.
func (l *Link) Name() string {
	if l == nil {
		return ""
	}

	return l.name
}

// Same returns true if the link is the same as the other linkable.
// Same meaning : same structure, same type. If id are equals, they should be the same link
func (l *Link) Same(other Element) bool {
	if l == nil && other == nil {
		return true
	} else if l == nil || other == nil {
		return false
	} else if l.id == other.Id() {
		return true
	}

	if lo, ok := other.(Linkable); !ok {
		return false
	} else {
		return l.id == lo.Id()
	}
}

// DeclaringClasses returns the classes that declare the link, obviously including CLASS_LINK itself
func (l *Link) DeclaringClasses() []Class {
	return []Class{CLASS_LINK}
}

// Validity returns the period the link is active
func (l *Link) Validity() Period {
	return l.validity.Copy()
}

// SetValidity sets the period the link is active
func (l *Link) SetValidity(p Period) {
	if l == nil {
		return
	}

	l.locks.Lock()
	defer l.locks.Unlock()

	l.validity = p
}

// Roles return the roles of the link.
// Result is sorted alphabetically
// For instance, Likes(subject=[Jean],Object=[Pizza]).Roles() will return ["object", "subject"]
func (l *Link) Roles() []string {
	if l == nil {
		return nil
	}

	l.locks.RLock()
	defer l.locks.RUnlock()

	size := len(l.operands)
	result := make([]string, size)
	index := 0
	for role := range l.operands {
		result[index] = role
		index++
	}

	slices.Sort(result)
	return result
}

// Has returns true and related value if the link has an operand with the given name.
// Otherwise, it returns nil, false.
func (l *Link) Has(name string) ([]Linkable, bool) {
	if l == nil {
		return nil, false
	}

	l.locks.RLock()
	defer l.locks.RUnlock()

	operands, ok := l.operands[name]
	return SliceCopy(operands), ok
}

// Add adds a new operand to the link with the given name.
// For instance, Likes(subject=[Jean],Object=[Pizza]).Add("Object", "Tiramisu")
// will make Likes(subject=[Jean],Object=[Pizza, Tiramisu])
func (l *Link) Add(name string, operand Linkable) {
	if l == nil || operand == nil {
		return
	}

	l.locks.Lock()
	defer l.locks.Unlock()

	l.operands[name] = append(l.operands[name], operand)
}

// Remove for a given role, specific values by predicate.
func (l *Link) Remove(name string, op func(linkable Linkable) bool) {
	if l == nil || op == nil {
		return
	}

	l.locks.Lock()
	defer l.locks.Unlock()

	newValues := make([]Linkable, 0)
	for _, element := range l.operands[name] {
		if !op(element) {
			newValues = append(newValues, element)
		}
	}
	l.operands[name] = newValues
}

// Operands return a copy of the role map
func (l *Link) Operands() map[string][]Linkable {
	if l == nil {
		return nil
	}

	l.locks.RLock()
	defer l.locks.RUnlock()

	result := make(map[string][]Linkable)
	for k, v := range l.operands {
		result[k] = SliceCopy(v)
	}

	return result
}
