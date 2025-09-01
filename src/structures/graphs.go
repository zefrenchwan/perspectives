package structures

import (
	"maps"
)

// DVGraph represents a directed valued graph as map of nodes and links.
// Given a node in the graph, it appears once (no outgoing link) or twice (ingoing and outgoing links)
type DVGraph[S comparable, L comparable] map[S]map[S]L

// GraphEdge defines an edge of a graph
type GraphEdge[S comparable, L comparable] struct {
	// Source of the edge
	Source S
	// Destination of the edge
	Destination S
	// Value of the edge
	Value L
}

// NewDVGraph returns a new empty DAG
func NewDVGraph[S comparable, L comparable]() DVGraph[S, L] {
	return make(DVGraph[S, L])
}

// AddNode adds a node in the graph and returns true if graph did not already contained it.
// If a node is destination of a link, it also appears in the map as primay value
func (d DVGraph[S, L]) AddNode(source S) bool {
	if _, found := d[source]; found {
		return false
	}

	d[source] = make(map[S]L)
	return true
}

// RemoveNode removes that node (as destination and source).
// It returns true if the node was in the graph, false otherwise
func (d DVGraph[S, L]) RemoveNode(node S) bool {
	for _, values := range d {
		delete(values, node)
	}

	_, found := d[node]
	delete(d, node)
	return found
}

// ClearLinks deletes all existing links from a node.
// If there were links return true, false otherwise
func (d DVGraph[S, L]) ClearLinks(node S) bool {
	if d == nil {
		return false
	} else if values, found := d[node]; !found {
		return false
	} else if len(values) == 0 {
		return false
	} else {
		delete(d, node)
		return true
	}
}

// Has returns true if node is in the graph
func (d DVGraph[S, L]) Has(node S) bool {
	_, found := d[node]
	return found
}

// Link adds source, destination and the link in between.
func (d DVGraph[S, L]) Link(source, destination S, link L) {
	_, sourceExists := d[source]
	_, destinationExists := d[destination]

	if !destinationExists {
		d[destination] = make(map[S]L)
	}

	if !sourceExists {
		d[source] = make(map[S]L)
	}

	d[source][destination] = link
}

// LinkWithoutCycle adds a link if it makes no cycle.
// If it makes a cycle, then rollback this link.
// Result is true if link was added, false otherwise
func (d DVGraph[S, L]) LinkWithoutCycle(source, destination S, link L) bool {
	if source == destination {
		// obvious cycle
		return false
	}

	// measure previous state
	_, sourceExists := d[source]
	_, destinationExists := d[destination]
	var previousLink bool
	var previousValue L

	if sourceExists {
		if l, found := d[source][destination]; found {
			previousLink = found
			previousValue = l
		} else {
			previousLink = found
		}
	}

	// perform the action
	d.Link(source, destination, link)

	// test if it would make a cycle, and then rollback
	if !d.HasCycle() {
		return true
	}

	// rollback because of the cycle
	if previousLink {
		d[source][destination] = previousValue
		return false
	} else {
		// delete the link for sure
		delete(d[source], destination)
	}

	if !destinationExists {
		delete(d, destination)
	}

	if !sourceExists {
		delete(d, source)
	}

	return false
}

// Neighbors returns a copy of the neighborhood of a node, false for not found in the graph
func (d DVGraph[S, L]) Neighbors(source S) (map[S]L, bool) {
	if values, found := d[source]; !found {
		return nil, false
	} else if len(values) != 0 {
		// clone values
		result := make(map[S]L)
		maps.Copy(result, values)

		return result, true
	} else {
		return nil, true
	}
}

// Edges return the edges from a node, if any (true) or empty, false
func (d DVGraph[S, L]) Edges(source S) ([]GraphEdge[S, L], bool) {
	values, found := d[source]
	if !found {
		return nil, false
	}

	var result []GraphEdge[S, L]
	for dest, link := range values {
		result = append(result, GraphEdge[S, L]{Source: source, Destination: dest, Value: link})
	}

	return result, true
}

// HasCycle returns true if graph contains a cycle
func (d DVGraph[S, L]) HasCycle() bool {
	// DFS needs a stack
	var stack []S

	// nodeExplorationStatus defines the exploration status for a given node
	// 0: never seen, 1: processing, 2: node and childs were explored
	nodeExplorationStatus := make(map[S]int)

	// for each unexplored node
	for startNode := range d {
		if nodeExplorationStatus[startNode] == 0 {
			// start a new DFS walkthrough from startNode
			stack = append(stack, startNode)
			nodeExplorationStatus[startNode] = 1

			for len(stack) > 0 {
				// Get last element
				node := stack[len(stack)-1]

				// test if all the childs of startNode were explored
				allChildsExplored := true
				// for each neighbor
				for neighbor := range d[node] {
					switch nodeExplorationStatus[neighbor] {
					case 1: // we already saw that node. Hence, there is a cycle
						return true
					case 0: // first time we see that node, explore it later
						allChildsExplored = false
						// neighbor is then about to be processed
						nodeExplorationStatus[neighbor] = 1
						stack = append(stack, neighbor)
					}
				}

				// If all neighbors (if any) were visited, node is then visited
				if allChildsExplored {
					stack = stack[:len(stack)-1]
					nodeExplorationStatus[node] = 2 // Node is completely visited
				}
			}
		}
	}

	return false
}

// Nodes returns the nodes of the graph
func (d DVGraph[S, L]) Nodes() []S {
	var result []S
	for key := range d {
		result = append(result, key)
	}

	return result
}

// Walk goes through a graph and reading once each element.
// Processor is a function to apply to each node (for instance, get its neighbors and so something)
func (d DVGraph[S, L]) Walk(starting S, processor func(source S)) {
	seen := make(map[S]bool)
	fifo := []S{starting}

	_, found := d[starting]
	if !found {
		return
	}

	for len(fifo) != 0 {
		element := fifo[0]
		fifo = fifo[1:]

		if seen[element] {
			continue
		}

		processor(element)
		seen[element] = true

		for other := range d[element] {
			if !seen[other] {
				fifo = append(fifo, other)
			}
		}
	}
}

// ReverseWalk walks through a graph, going backward (from a node to predecessors)
func (d DVGraph[S, L]) ReverseWalk(starting S, processor func(current S)) {
	_, found := d[starting]
	if !found {
		return
	}

	reverseGraph := make(map[S][]S)
	for key, values := range d {
		for value := range values {
			if existing, found := reverseGraph[value]; !found {
				reverseGraph[value] = []S{key}
			} else {
				existing = append(existing, key)
				reverseGraph[value] = existing
			}
		}
	}

	seen := make(map[S]bool)
	fifo := []S{starting}

	for len(fifo) != 0 {
		element := fifo[0]
		fifo = fifo[1:]

		if seen[element] {
			continue
		}

		processor(element)
		seen[element] = true

		for _, other := range reverseGraph[element] {
			if !seen[other] {
				fifo = append(fifo, other)
			}
		}
	}
}

// EdgesFrom gets the edges since a starting node
func (d DVGraph[S, L]) EdgesFrom(starting S) []GraphEdge[S, L] {
	_, found := d[starting]
	if !found {
		return nil
	}

	elements := []S{starting}
	var result []GraphEdge[S, L]
	seen := make(map[S]bool)

	for len(elements) != 0 {
		element := elements[0]
		elements = elements[1:]

		if seen[element] {
			continue
		}

		for destination, link := range d[element] {
			edge := GraphEdge[S, L]{Source: element, Destination: destination, Value: link}
			result = append(result, edge)

			if !seen[destination] {
				elements = append(elements, destination)
			}
		}

		seen[element] = true
	}

	return result
}
