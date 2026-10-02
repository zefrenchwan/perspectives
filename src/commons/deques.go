package commons

import (
	"iter"
	"slices"
	"sync"
)

// Deque is a double-ended queue.
// To use for LIFO, best performance option would be add last, pop last
// To use for FIFO, best performance option would be add last, pop first
type Deque[E any] interface {
	// AddFirst adds an element to the front of the deque.
	AddFirst(E)
	// AddLast adds an element to the back of the deque.
	AddLast(E)
	// PopFirst removes and returns the first element of the deque.
	// If no element, then it returns empty and false
	PopFirst() (E, bool)
	// PopLast removes and returns the last element of the deque.
	// If no element, then it returns empty and false
	PopLast() (E, bool)
	// Size returns the actual size of the deque.
	Size() int
	// Get returns the element at the given index, if any.
	// For no element, the second result is false
	Get(int) (E, bool)
	// Elements returns an iterator over the elements of the deque.
	Elements() iter.Seq[E]
	// Clear removes all elements from the deque.
	Clear()
}

// localDeque is a local implementation of a deque, with a slice as the underlying data structure.
// Use it for small deques.
type localDeque[E any] struct {
	mu      sync.RWMutex // Protects the deque from concurrent access
	content []E
}

func (ll *localDeque[E]) AddFirst(e E) {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	ll.content = slices.Insert(ll.content, 0, e)
}

func (ll *localDeque[E]) AddLast(e E) {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	ll.content = append(ll.content, e)
}

func (ll *localDeque[E]) PopFirst() (E, bool) {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	var empty E
	if len(ll.content) == 0 {
		return empty, false
	}

	result := ll.content[0]
	ll.content[0] = empty
	ll.content = ll.content[1:]

	return result, true
}

func (ll *localDeque[E]) PopLast() (E, bool) {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	var empty E
	if len(ll.content) == 0 {
		return empty, false
	}

	lastIndex := len(ll.content) - 1
	result := ll.content[lastIndex]
	ll.content[lastIndex] = empty
	ll.content = ll.content[:lastIndex]

	return result, true
}

func (ll *localDeque[E]) Clear() {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	clear(ll.content)
	ll.content = ll.content[:0]
}

func (ll *localDeque[E]) Elements() iter.Seq[E] {
	ll.mu.RLock()
	defer ll.mu.RUnlock()

	// Since iterators are evaluated lazily, we clone the underlying slice
	// before returning. This prevents deadlocks if the caller mutates the deque
	// during iteration, and prevents race conditions. This aligns well with
	// the "Use it for small deques" heuristic.
	return slices.Values(slices.Clone(ll.content))
}

func (ll *localDeque[E]) Size() int {
	ll.mu.RLock()
	defer ll.mu.RUnlock()

	return len(ll.content)
}

func (ll *localDeque[E]) Get(index int) (E, bool) {
	ll.mu.RLock()
	defer ll.mu.RUnlock()

	if index < 0 || index >= len(ll.content) {
		var empty E
		return empty, false
	}

	return ll.content[index], true
}

// NewDeque creates a new local deque with the given capacity.
// Implementation is thread safe.
func NewDeque[E any](capacity int) Deque[E] {
	return &localDeque[E]{
		content: make([]E, 0, capacity),
	}
}
