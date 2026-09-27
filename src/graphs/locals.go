package graphs

import (
	"errors"
	"iter"
	"slices"
	"sync"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// localEdge is inner structure for edges : destination and value
type localEdge[V any] struct {
	// destination of the edge (source is stored in map)
	destination string
	// value of the edge
	value V
}

// readOnlyEdge is a struct used for read-only access to edges when building sequences
type readOnlyEdge[N Node, V any] struct {
	destination N
	value       V
}

// localGraph is a graph implementation that stores edges and nodes locally as a map of edges and nodes.
// Implementation is in memory, thread-safe, and defensive (copies on read).
// Its purpose is more to validate the graph interface, and should not be used in production for large graphs.
type localGraph[N Node, V any] struct {
	// synchronizer is used as a read-write lock to ensure thread safety when accessing the graph
	synchronizer sync.RWMutex
	// edges of the graph by source node id
	edges map[string][]localEdge[V]
	// nodes of the graph by node id (to have content once)
	nodes map[string]N
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

////////////////////////////////
// DIRECTED UNWEIGHTED GRAPHS //
////////////////////////////////

type directedUnweightedGraph[N Node] struct {
	adapter localGraph[N, struct{}]
}

func NewDGraph[N Node]() DGraph[N] {
	return &directedUnweightedGraph[N]{
		adapter: newLocalGraph[N, struct{}](),
	}
}

// SetNode adds a node to the graph or change its status.
// It does not raise any error.
func (g *directedUnweightedGraph[N]) SetNode(node N) error {
	return g.adapter.upsertNode(node)
}

// RemoveNode removes the node by id.
// It raises an error if the node does not exist.
func (g *directedUnweightedGraph[N]) RemoveNode(node N) error {
	return g.adapter.removeNode(node.Id())
}

// Nodes iterates over all the nodes.
// It guarantees that the nodes are returned once only.
func (g *directedUnweightedGraph[N]) Nodes() iter.Seq[N] {
	return g.adapter.nodesIterator()
}

// HasNode returns true if the node exists in the graph (based on id).
func (g *directedUnweightedGraph[N]) HasNode(node N) bool {
	return g.adapter.hasNode(node.Id())
}

// Link adds a link (if no one existed).
// It returns an error if the nodes do not exist.
// It has no effect if the link already exists.
func (g *directedUnweightedGraph[N]) Link(source, destination N) error {
	empty := struct{}{}
	return g.adapter.linkNodes(source.Id(), destination.Id(), empty)
}

// Unlink removes a link between two nodes if any.
// It returns an error if the nodes do not exist.
func (g *directedUnweightedGraph[N]) Unlink(source, destination N) error {
	return g.adapter.unlinkNodes(source.Id(), destination.Id())
}

// Successors returns the successors of the node.
// It performs a defensive copy to avoid locking while iterating.
func (g *directedUnweightedGraph[N]) Successors(source N) iter.Seq[N] {
	return g.adapter.nodeNeighbors(source.Id())
}

// HasLink returns true if the link exists in the graph (hence, nodes exist too).
func (g *directedUnweightedGraph[N]) HasLink(source, destination N) bool {
	_, has := g.adapter.hasLink(source.Id(), destination.Id())
	return has
}

////////////////////////////////
// DIRECTED UNWEIGHTED GRAPHS //
////////////////////////////////

type directedWeightedGraph[N Node, V any] struct {
	adapter localGraph[N, V]
}

func NewDWGraph[N Node, V any]() DWGraph[N, V] {
	return &directedWeightedGraph[N, V]{
		adapter: newLocalGraph[N, V](),
	}
}

// SetNode adds a node to the graph or change its status.
// It does not raise any error.
func (g *directedWeightedGraph[N, V]) SetNode(node N) error {
	return g.adapter.upsertNode(node)
}

// RemoveNode removes the node by id.
// It raises an error if the node does not exist.
func (g *directedWeightedGraph[N, V]) RemoveNode(node N) error {
	return g.adapter.removeNode(node.Id())
}

// Nodes iterates over all the nodes.
// It guarantees that the nodes are returned once only.
func (g *directedWeightedGraph[N, V]) Nodes() iter.Seq[N] {
	return g.adapter.nodesIterator()
}

// HasNode returns true if the node exists in the graph (based on id).
func (g *directedWeightedGraph[N, V]) HasNode(node N) bool {
	return g.adapter.hasNode(node.Id())
}

// Link adds a link (if no one existed) between two nodes with given value.
// It returns an error if the nodes do not exist.
// It has no effect if the link already exists.
func (g *directedWeightedGraph[N, V]) Link(source N, destination N, value V) error {
	return g.adapter.linkNodes(source.Id(), destination.Id(), value)
}

// Unlink removes a link between two nodes if any.
// It returns an error if the nodes do not exist.
func (g *directedWeightedGraph[N, V]) Unlink(source, destination N) error {
	return g.adapter.unlinkNodes(source.Id(), destination.Id())
}

// Successors returns the successors of the node and edges values.
// It performs a defensive copy to avoid locking while iterating.
func (g *directedWeightedGraph[N, V]) Successors(source N) iter.Seq2[N, V] {
	return g.adapter.valuedNeighbors(source.Id())
}

// HasLink returns true if the link exists in the graph (hence, nodes exist too).
func (g *directedWeightedGraph[N, V]) HasLink(source, destination N) (V, bool) {
	return g.adapter.hasLink(source.Id(), destination.Id())
}
