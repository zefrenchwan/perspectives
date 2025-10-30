package commons

import (
	"iter"
	"time"
)

// DynamicGraph is a directed graph with edges active during a given period.
type DynamicGraph[V Identifiable, E any] interface {
	// Set adds vertex if not found, or changes its value if already present
	Set(V)
	// Connect source and destination during period, with given value.
	// If vertices are not in graph, add them too.
	Relate(source, destination V, value E, period Period)
	// Remove the edge (if any) from source to destination
	Remove(source, destination V)
	// Neighbors provides neighbors from source at a given time, as vertices and edge values
	Neighbors(source V, moment time.Time) iter.Seq2[V, E]
	// Vertices returns an iterator over the vertices of that graph
	Vertices() iter.Seq[V]
	// Lookup finds (if any) a vertex by id
	Lookup(id string) (V, bool)
}

// localDynamicEdge decorates an edge with the period the edge was active
type localDynamicEdge[E any] struct {
	// value is the edge value
	value E
	// activity is the period the edge is active
	activity Period
}

// localDynamicGraph implements DynamicGraph as an adjacency dynamic matrix
type localDynamicGraph[V Identifiable, E any] struct {
	// elements link vertices with their id
	elements map[string]V
	// edges as sourceId => destinationId => decorated edge
	edges map[string]map[string]localDynamicEdge[E]
}

// Set either adds the vertex, or updates its value
func (g *localDynamicGraph[V, E]) Set(vertex V) {
	if g == nil {
		return
	}

	if g.elements == nil {
		g.elements = make(map[string]V)
	}

	g.elements[vertex.Id()] = vertex
}

// Connect links source and destination at a given value, since creationTime
func (g *localDynamicGraph[V, E]) Relate(source, destination V, edge E, period Period) {
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
func (g *localDynamicGraph[V, E]) Remove(source, destination V) {
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
func (g *localDynamicGraph[V, E]) Neighbors(source V, moment time.Time) iter.Seq2[V, E] {
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
func (g *localDynamicGraph[V, E]) Vertices() iter.Seq[V] {
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
func (g *localDynamicGraph[V, E]) Lookup(id string) (V, bool) {
	var empty V
	if g == nil || g.elements == nil {
		return empty, false
	}

	result, found := g.elements[id]
	return result, found
}

// NewDynamicGraph returns a new empty graph
func NewDynamicGraph[V Identifiable, E any]() DynamicGraph[V, E] {
	result := new(localDynamicGraph[V, E])
	result.edges = make(map[string]map[string]localDynamicEdge[E])
	result.elements = make(map[string]V)
	return result
}
