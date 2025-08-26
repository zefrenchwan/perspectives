package structures

import (
	"errors"
	"maps"
)

// DAG represents a directed acyclic graph as map of nodes and links.
// Given a node in the graph, it appears once (no outgoing link) or twice (ingoing and outgoing links)
type DAG[S comparable, L comparable] map[S]map[S]L

// NewDAG returns a new empty DAG
func NewDAG[S comparable, L comparable]() DAG[S, L] {
	return make(DAG[S, L])
}

// AddNode adds a node in the graph and returns true if graph did not already contained it.
// If a node is destination of a link, it also appears in the map as primay value
func (d DAG[S, L]) AddNode(source S) bool {
	if _, found := d[source]; found {
		return false
	}

	d[source] = make(map[S]L)
	return true
}

// Link adds source, destination and the link in between.
// It returns an error if a cycle appears
func (d DAG[S, L]) Link(source, destination S, link L) error {
	if source == destination {
		return errors.New("adding same source and destination makes an obvious cycle")
	}

	var sourceExisted, destinationExisted, linkExisted bool
	var previousLink L

	// get previous values for a potential rollback
	_, sourceExisted = d[source]
	_, destinationExisted = d[destination]

	if !sourceExisted {
		d[source] = make(map[S]L)
	} else if destinationExisted {
		previousLink, linkExisted = d[source][destination]
	}

	// perform the action
	d[source][destination] = link
	if !destinationExisted {
		d[destination] = make(map[S]L)
	}

	// test if adding element would create a cycle.
	// If so, perform a rollback !
	if d.hasCycle() {
		// Rollback
		if !linkExisted {
			delete(d[source], destination)
		} else {
			d[source][destination] = previousLink
		}

		if !destinationExisted {
			delete(d[source], destination)
			delete(d, destination)
		}

		if !sourceExisted {
			delete(d, source)
		}

		return errors.New("adding this link would create a cycle")
	}

	return nil
}

// Neighbors returns a copy of the neighborhood of a node, false for not found in the graph
func (d DAG[S, L]) Neighbors(source S) (map[S]L, bool) {
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

// hasCycle returns true if graph contains a cycle
func (d DAG[S, L]) hasCycle() bool {
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
