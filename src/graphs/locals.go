package graphs

import (
	"errors"
	"iter"
	"maps"
	"sync"
)

// localGraph defines a weighted graph that stores nodes and links locally.
// Each node has a value of type V.
// Implementation stores the values per identifier, then manages the ids only.
type localGraph[N Node, V any] struct {
	// mutex to prevent concurrent access to the nodes and links maps.
	mutex sync.RWMutex
	// nodes stores the nodes of the graph by id
	nodes map[string]N
	// links stores the links of the graph by source id and target id
	links map[string]map[string]V
}

// newLocalGraph creates a new localGraph instance.
func newLocalGraph[N Node, V any]() *localGraph[N, V] {
	return &localGraph[N, V]{
		nodes: make(map[string]N),
		links: make(map[string]map[string]V),
	}
}

// upsertNode adds a node to the graph.
func (g *localGraph[N, V]) upsertNode(node N) {
	if node != nil {
		g.mutex.Lock()
		defer g.mutex.Unlock()
		id := node.Id()
		g.nodes[id] = node
	}
}

// removeNode removes a node from the graph.
func (g *localGraph[N, V]) removeNode(node N) {
	if node != nil {
		id := node.Id()
		g.mutex.Lock()
		defer g.mutex.Unlock()
		delete(g.nodes, id)
		for _, linksMap := range g.links {
			delete(linksMap, id)
		}
		delete(g.links, id)
	}
}

// hasNode checks if a node exists in the graph
func (g *localGraph[N, V]) hasNode(node N) bool {
	if node != nil {
		g.mutex.RLock()
		defer g.mutex.RUnlock()
		id := node.Id()
		_, ok := g.nodes[id]
		return ok
	}

	return false
}

// allNodes returns all nodes in the graph as a sequence
func (g *localGraph[N, V]) allNodes() iter.Seq[N] {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	return maps.Values(g.nodes)
}

// linkNodesWithValue links two nodes in the graph with a value
func (g *localGraph[N, V]) linkNodesWithValue(source, destination N, value V) error {
	if source == nil || destination == nil {
		return errors.New("source and destination nodes cannot be nil")
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()
	if _, hasSource := g.nodes[source.Id()]; !hasSource {
		return errors.New("source should exist in graph")
	} else if _, hasDest := g.nodes[destination.Id()]; !hasDest {
		return errors.New("destination should exist in graph")
	}

	if _, has := g.links[source.Id()]; !has {
		g.links[source.Id()] = make(map[string]V)
	}

	g.links[source.Id()][destination.Id()] = value
	return nil
}

// linkNodes links two nodes in the graph with an empty value so that it works as an unweighted link
func (g *localGraph[N, V]) linkNodes(source, destination N) error {
	var empty V
	return g.linkNodesWithValue(source, destination, empty)
}
