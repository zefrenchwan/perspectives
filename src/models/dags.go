package models

import (
	"errors"
	"slices"
)

// DAG is a directed acyclic graph
type DAG[S interface{ comparable }, L interface{ comparable }] struct {
	successors map[S]map[S]L // successors link source to destination with the value
	roots      map[S]bool    // value is here means that value has no predecessor
}

// NewDAG returns a new empty DAG
func NewDAG[S interface{ comparable }, L interface{ comparable }]() DAG[S, L] {
	var result DAG[S, L]
	result.successors = make(map[S]map[S]L)
	result.roots = make(map[S]bool)
	return result
}

func (d DAG[S, L]) Childs(value S) []S {
	var result []S
	for child := range d.successors[value] {
		result = append(result, child)
	}

	return result
}

// AddNode adds a node with no link if not present
func (d DAG[S, L]) AddNode(node S) {
	if _, found := d.successors[node]; !found {
		d.successors[node] = make(map[S]L)
		d.roots[node] = true
	}
}

// UpsertNodesLinks inserts source and dest if not set already.
// It returns an error if there is a cycle
func (d DAG[S, L]) UpsertNodesLinks(source, dest S, value L) error {
	if source == dest {
		return errors.New("cannot link a node with itself: would make a cycle")
	}
	// to rollback for a check failure
	var rollbackSource, rollbackDest, rollbackValue bool
	var previous L
	if _, keyFound := d.successors[source]; !keyFound {
		rollbackSource = true
	} else if p, valueFound := d.successors[source][dest]; !valueFound {
		rollbackDest = true
	} else {
		previous = p
		rollbackValue = true
	}

	// set values
	if _, found := d.successors[source]; found {
		d.successors[source][dest] = value
	} else {
		d.successors[source] = make(map[S]L)
		d.successors[source][dest] = value
		d.roots[source] = true
	}

	// perform a check to test if graph would be cyclic after deletion
	if d.hasCycle() {
		switch {
		case rollbackValue:
			d.successors[source][dest] = previous
		case rollbackDest:
			delete(d.successors[source], dest)
		case rollbackSource:
			delete(d.successors, source)
		}

		return errors.New("cannot add values because a cycle would exist")
	} else {
		delete(d.roots, dest)
		return nil
	}
}

// hasCycle returns true if there is a cycle within the DAG (AND IT SHOULD NOT)
func (d DAG[S, L]) hasCycle() bool {
	stack := make([]S, 0)
	for root := range d.roots {
		stack = append(stack, root)
	}

	size := len(stack)
	for size != 0 {
		element, stack := stack[size-1], stack[0:size-1]
		size = len(stack)

		for child := range d.successors[element] {
			previous := slices.IndexFunc(stack, func(v S) bool { return v == child })
			if previous >= 0 {
				return true
			}
			stack = append(stack, child)
		}
	}

	return false
}
