package graphs

import (
	"errors"
	"iter"
	"slices"
	"sync"

	"github.com/zefrenchwan/perspectives.git/commons"
)

type localEdge[V any] struct {
	destination string
	value       V
}

type readOnlyEdge[N Node, V any] struct {
	destination N
	value       V
}

type localGraph[N Node, V any] struct {
	synchronizer sync.RWMutex
	edges        map[string][]localEdge[V]
	nodes        map[string]N
}

func newLocalGraph[N Node, V any]() localGraph[N, V] {
	return localGraph[N, V]{
		edges: make(map[string][]localEdge[V]),
		nodes: make(map[string]N),
	}
}

func (g *localGraph[N, V]) hasNode(nodeId string) bool {
	g.synchronizer.RLock()
	defer g.synchronizer.RUnlock()
	_, has := g.nodes[nodeId]
	return has
}

func (g *localGraph[N, V]) upsertNode(node N) error {
	g.synchronizer.Lock()
	defer g.synchronizer.Unlock()

	g.nodes[node.Id()] = node
	return nil
}

func (g *localGraph[N, V]) removeNode(nodeId string) error {
	g.synchronizer.Lock()
	defer g.synchronizer.Unlock()

	if _, has := g.nodes[nodeId]; !has {
		return errors.New("node does not exist")
	}

	delete(g.nodes, nodeId)
	delete(g.edges, nodeId)

	filterFunc := func(e localEdge[V]) bool {
		return e.destination != nodeId
	}

	for key, values := range g.edges {
		filteredValues := commons.SlicesFilter(values, filterFunc)
		g.edges[key] = filteredValues
	}

	return nil
}

func (g *localGraph[N, V]) nodesIterator() iter.Seq[N] {
	g.synchronizer.RLock()
	defer g.synchronizer.RUnlock()

	localCopy := make([]N, len(g.nodes))
	i := 0
	for _, node := range g.nodes {
		localCopy[i] = node
		i++
	}

	return slices.Values(localCopy)
}

func (g *localGraph[N, V]) linkNodes(source, destination string, value V) error {
	g.synchronizer.Lock()
	defer g.synchronizer.Unlock()

	if _, hasSource := g.nodes[source]; !hasSource {
		return errors.New("source node does not exist")
	} else if _, hasDestination := g.nodes[destination]; !hasDestination {
		return errors.New("destination node does not exist")
	}

	found := false
	sourceEdges := g.edges[source]
	for index, currentEdge := range sourceEdges {
		if currentEdge.destination == destination {
			found = true
			g.edges[source][index].value = value
			break
		}
	}

	if !found {
		g.edges[source] = append(g.edges[source], localEdge[V]{destination, value})
	}

	return nil
}

func (g *localGraph[N, V]) unlinkNodes(source, destination string) error {
	g.synchronizer.Lock()
	defer g.synchronizer.Unlock()

	if _, hasSource := g.nodes[source]; !hasSource {
		return errors.New("source node does not exist")
	} else if _, hasDestination := g.nodes[destination]; !hasDestination {
		return errors.New("destination node does not exist")
	}

	found := false
	matchingIndex := 0
	sourceEdges := g.edges[source]
	for index, currentEdge := range sourceEdges {
		if currentEdge.destination == destination {
			found = true
			matchingIndex = index
			break
		}
	}

	if found {
		size := len(sourceEdges)
		sourceEdges[matchingIndex] = sourceEdges[size-1]
		g.edges[source] = sourceEdges[:size-1]
	}

	if len(g.edges[source]) == 0 {
		delete(g.edges, source)
	}

	return nil
}

func (g *localGraph[N, V]) hasLink(source, destination string) (V, bool) {
	g.synchronizer.RLock()
	defer g.synchronizer.RUnlock()

	var empty V
	if _, hasSource := g.nodes[source]; !hasSource {
		return empty, false
	} else if _, hasDest := g.nodes[destination]; !hasDest {
		return empty, false
	}

	for _, currentEdge := range g.edges[source] {
		if currentEdge.destination == destination {
			return currentEdge.value, true
		}
	}

	return empty, false
}

func (g *localGraph[N, V]) nodeNeighbors(source string) iter.Seq[N] {
	g.synchronizer.RLock()
	defer g.synchronizer.RUnlock()

	if _, hasSource := g.nodes[source]; !hasSource {
		return nil
	}

	currentEdges := g.edges[source]
	copyValues := make([]N, len(currentEdges))
	index := 0
	for _, edge := range currentEdges {
		node := g.nodes[edge.destination]
		copyValues[index] = node
		index++
	}

	return slices.Values(copyValues)
}

func (g *localGraph[N, V]) valuedNeighbors(source string) iter.Seq2[N, V] {
	return func(yield func(N, V) bool) {
		g.synchronizer.RLock()
		// unlock dealt manually, before end of function
		if _, hasSource := g.nodes[source]; !hasSource {
			g.synchronizer.RUnlock()
			return
		}

		// make a copy of neighbors as a snapshot
		currentEdges := g.edges[source]
		snapshot := make([]readOnlyEdge[N, V], 0, len(currentEdges))

		for _, edge := range currentEdges {
			if destNode, exists := g.nodes[edge.destination]; exists {
				snapshot = append(snapshot, readOnlyEdge[N, V]{
					destination: destNode,
					value:       edge.value,
				})
			}
		}

		// copy done, unlock
		g.synchronizer.RUnlock()

		// iterator on snapshot, no race issue
		for _, roEdge := range snapshot {
			if !yield(roEdge.destination, roEdge.value) {
				return
			}
		}
	}
}

/////////////////////////////////////
// LINK TO GENERAL IMPLEMENTATIONS //
/////////////////////////////////////

type directedUnweightedGraph[N Node] struct {
	adapter localGraph[N, struct{}]
}

func NewDGraph[N Node]() DGraph[N] {
	return &directedUnweightedGraph[N]{
		adapter: newLocalGraph[N, struct{}](),
	}
}

func (g *directedUnweightedGraph[N]) SetNode(node N) error {
	return g.adapter.upsertNode(node)
}

func (g *directedUnweightedGraph[N]) RemoveNode(node N) error {
	return g.adapter.removeNode(node.Id())
}

func (g *directedUnweightedGraph[N]) Nodes() iter.Seq[N] {
	return g.adapter.nodesIterator()
}

func (g *directedUnweightedGraph[N]) HasNode(node N) bool {
	return g.adapter.hasNode(node.Id())
}

func (g *directedUnweightedGraph[N]) Link(source, destination N) error {
	empty := struct{}{}
	return g.adapter.linkNodes(source.Id(), destination.Id(), empty)
}

func (g *directedUnweightedGraph[N]) Unlink(source, destination N) error {
	return g.adapter.unlinkNodes(source.Id(), destination.Id())
}

func (g *directedUnweightedGraph[N]) Successors(source N) iter.Seq[N] {
	return g.adapter.nodeNeighbors(source.Id())
}

func (g *directedUnweightedGraph[N]) HasLink(source, destination N) bool {
	_, has := g.adapter.hasLink(source.Id(), destination.Id())
	return has
}
