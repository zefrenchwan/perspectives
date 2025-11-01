package commons

import (
	"errors"
	"iter"
	"time"
)

// GraphWalker is a walkthrough run (using a BFS)
type GraphWalker[V Identifiable, E any] interface {
	// Position is the current vertex
	Position() V
	// Source is the vertex that walker came from
	Source() V
	// SourceEdge is the value of the edge from the source to the current position
	SourceEdge() E
	// Next moves to next value (if any) and returns true if there was a next value
	Next() bool
	// Stop ends walk immediatly
	Stop()
}

// DynamicGraph is a directed graph with edges active during a given period.
type DynamicGraph[V Identifiable, E any] interface {
	// Neighbors provides neighbors from source at a given time, as vertices and edge values
	Neighbors(source V, moment time.Time) iter.Seq2[V, E]
	// Vertices returns an iterator over the vertices of that graph
	Vertices() iter.Seq[V]
	// Lookup finds (if any) a vertex by id
	Lookup(id string) (V, bool)
}

// DynamicConnectionGraph is a dynamic graph with explicit connections between vertices
type DynamicConnectionGraph[V Identifiable, E any] interface {
	// Features of DynamicGraph are relevant when connections are explicit
	DynamicGraph[V, E]
	// Set adds vertex if not found, or changes its value if already present
	Set(V)
	// Connect source and destination during period, with given value.
	// If vertices are not in graph, add them too.
	Connect(source, destination V, value E, period Period)
	// Remove the edge (if any) from source to destination
	Remove(source, destination V)
}

// localDynamicEdge decorates an edge with the period the edge was active
type localDynamicEdge[E any] struct {
	// value is the edge value
	value E
	// activity is the period the edge is active
	activity Period
}

// localDynamicConnectionGraph implements DynamicConnectionGraph as an adjacency dynamic matrix
type localDynamicConnectionGraph[V Identifiable, E any] struct {
	// elements link vertices with their id
	elements map[string]V
	// edges as sourceId => destinationId => decorated edge
	edges map[string]map[string]localDynamicEdge[E]
}

// Set either adds the vertex, or updates its value
func (g *localDynamicConnectionGraph[V, E]) Set(vertex V) {
	if g == nil {
		return
	}

	if g.elements == nil {
		g.elements = make(map[string]V)
	}

	g.elements[vertex.Id()] = vertex
}

// Connect links source and destination at a given value, since creationTime
func (g *localDynamicConnectionGraph[V, E]) Connect(source, destination V, edge E, period Period) {
	if g == nil {
		return
	}

	if g.elements == nil {
		g.elements = make(map[string]V)
	}

	sourceId := source.Id()
	destId := destination.Id()

	g.elements[sourceId] = source
	g.elements[destId] = destination

	if g.edges == nil {
		g.edges = make(map[string]map[string]localDynamicEdge[E])
	}

	if edges, found := g.edges[sourceId]; !found {
		values := make(map[string]localDynamicEdge[E])
		values[destId] = localDynamicEdge[E]{value: edge, activity: period}
		g.edges[sourceId] = values
	} else if value, found := edges[destination.Id()]; !found {
		edges[destId] = localDynamicEdge[E]{value: edge, activity: period}
		g.edges[sourceId] = edges
	} else {
		value.activity = period
		value.value = edge
		edges[destId] = value
		g.edges[sourceId] = edges
	}
}

// Remove unlinks from source to destination instantly
func (g *localDynamicConnectionGraph[V, E]) Remove(source, destination V) {
	if g == nil {
		return
	} else if len(g.elements) == 0 {
		return
	} else if len(g.edges) == 0 {
		return
	} else if values, foundSource := g.edges[source.Id()]; !foundSource {
		return
	} else {
		delete(values, destination.Id())
		g.edges[source.Id()] = values
	}
}

// Neighbors returns the neighbors of source active at moment.
func (g *localDynamicConnectionGraph[V, E]) Neighbors(source V, moment time.Time) iter.Seq2[V, E] {
	if g == nil {
		return nil
	}

	return func(yield func(V, E) bool) {
		values, found := g.edges[source.Id()]
		if !found {
			return
		}

		for destId, destValue := range values {
			matching := g.elements[destId]
			if destValue.activity.Contains(moment) {
				if !yield(matching, destValue.value) {
					return
				}
			}
		}
	}
}

// Vertices returns an iterator over the vertices of that graph
func (g *localDynamicConnectionGraph[V, E]) Vertices() iter.Seq[V] {
	return func(yield func(V) bool) {
		if g == nil {
			return
		}

		for _, value := range g.elements {
			if !yield(value) {
				return
			}
		}
	}
}

// Lookup finds an element by id (if any).
// It returns said element if any, or zero value with false
func (g *localDynamicConnectionGraph[V, E]) Lookup(id string) (V, bool) {
	var empty V
	if g == nil || g.elements == nil {
		return empty, false
	}

	result, found := g.elements[id]
	return result, found
}

// NewDynamicConnectionGraph returns a new empty graph
func NewDynamicConnectionGraph[V Identifiable, E any]() DynamicConnectionGraph[V, E] {
	result := new(localDynamicConnectionGraph[V, E])
	result.edges = make(map[string]map[string]localDynamicEdge[E])
	result.elements = make(map[string]V)
	return result
}

// localGraphOption is an option, from current position, to visit next nodes
type localGraphOption[V Identifiable, E any] struct {
	// current position
	current V
	// destination to reach (possiple destination)
	destination V
	// value is the edge value from current to destination
	value E
}

// Id returns an id, unique for a given current and destination couple
func (o localGraphOption[V, E]) Id() string {
	return NewCompositeId(o.current, o.destination)
}

// localGraphWalker is a BFS walker
type localGraphWalker[V Identifiable, E any] struct {
	// processingTime
	processingTime time.Time
	// queue to deal with BFS
	elements []localGraphOption[V, E]
	// graph to walk through
	graph DynamicGraph[V, E]
	// current is the last choice made
	current localGraphOption[V, E]
	// seenEdges are all visited edges
	seenEdges map[string]bool
}

// Source returns the vertex walker came from
func (w *localGraphWalker[V, E]) Source() V {
	return w.current.current
}

// SourceEdge returns the edge from source to current position
func (w *localGraphWalker[V, E]) SourceEdge() E {
	return w.current.value
}

// Position returns the current vertex the walker is on
func (w *localGraphWalker[V, E]) Position() V {
	return w.current.destination
}

// Next moves to next vertex (if any) and returns true if we may go on, false otherwise
func (w *localGraphWalker[V, E]) Next() bool {
	if w == nil || len(w.elements) == 0 || w.graph == nil {
		return false
	}

	current := w.elements[0]
	w.seenEdges[current.Id()] = true
	w.elements = w.elements[1:]
	position := current.destination
	w.current = current

	for neighor, value := range w.graph.Neighbors(position, w.processingTime) {
		option := localGraphOption[V, E]{current: position, destination: neighor, value: value}
		if !w.seenEdges[option.Id()] {
			w.elements = append(w.elements, option)
		}
	}

	return true
}

// Stop ends the walk. Next will return false since then
func (w *localGraphWalker[V, E]) Stop() {
	w.elements = nil
	w.elements = make([]localGraphOption[V, E], 0)
}

// NewDynamicGraphWalker walks from startingPoint at current time within the base graph
func NewDynamicGraphWalker[V Identifiable, E any](base DynamicGraph[V, E], startingPoint V, currentTime time.Time) GraphWalker[V, E] {
	if base == nil {
		return nil
	} else if _, found := base.Lookup(startingPoint.Id()); !found {
		return nil
	}

	result := new(localGraphWalker[V, E])
	result.graph = base
	result.processingTime = currentTime
	result.seenEdges = make(map[string]bool)

	for vertex, edge := range base.Neighbors(startingPoint, currentTime) {
		option := localGraphOption[V, E]{current: startingPoint, destination: vertex, value: edge}
		result.elements = append(result.elements, option)
	}

	return result
}

// LocalAction in a graph may apply during a walk.
// For instance, change destination based on source and edge.
type LocalAction[V Identifiable, E any] interface {
	// Apply may choose to apply or not
	Apply(source V, destination V, edge E) error
}

// simpleLocalAction is a simple action, with no state (as a function)
type simpleLocalAction[V Identifiable, E any] func(source V, destination V, edge E) error

// Apply just calls s
func (s simpleLocalAction[V, E]) Apply(source, destination V, edge E) error {
	return s(source, destination, edge)
}

// NewLocalAction builds a local action to apply during a walk
func NewLocalAction[V Identifiable, E any](action func(source V, destination V, edge E) error) LocalAction[V, E] {
	return simpleLocalAction[V, E](action)
}

// DynamicGraphSpreadAction executes an action in a graph based on an iterator
func DynamicGraphSpreadAction[V Identifiable, E any](
	graph DynamicGraph[V, E],
	startingPoint V,
	currentTime time.Time,
	action LocalAction[V, E],
) error {
	var allErrors error
	if graph == nil || action == nil {
		return nil
	}

	iterator := NewDynamicGraphWalker(graph, startingPoint, currentTime)
	for iterator.Next() {
		source := iterator.Source()
		position := iterator.Position()
		edge := iterator.SourceEdge()

		if err := action.Apply(source, position, edge); err != nil {
			allErrors = errors.Join(allErrors, err)
		}
	}

	return allErrors
}

// DynamicGraphAction executes an action from source to its direct neighbors
func DynamicGraphAction[V Identifiable, E any](graph DynamicGraph[V, E], source V, moment time.Time, action LocalAction[V, E]) error {
	var allErrors error
	for neighor, edge := range graph.Neighbors(source, moment) {
		if err := action.Apply(source, neighor, edge); err != nil {
			allErrors = errors.Join(allErrors, err)
		}
	}

	return allErrors
}
