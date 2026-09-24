package graphs

import "iter"

// Graph defines the generic graph interface: manage nodes.
// Each node is actually an implementation of the Node interface as N.
// This way, one may add any type of node to the graph as long as it implements the Node interface.
type Graph[N Node] interface {
	// SetNode upserts a node.
	// If node already exists, it changes its value.
	SetNode(N) error
	// RemoveNode a node if isolated (not linked to any other node)
	RemoveNode(N) error
	// Nodes return the nodes within the graph
	Nodes() iter.Seq[N]
	// HasNode returns true if the node is within the graph
	HasNode(N) bool
}

// DGraph defines the generic directed graph interface: manage nodes and links.
// This is the current implementation of a directed graph with no values linked to the edges.
type DGraph[N Node] interface {
	// Graph defines the generic graph interface: manage nodes.
	Graph[N]
	// Link a node to another node (if exists).
	// It does nothing if the link already exists.
	Link(N, N) error
	// Unlink a node from another node (if exists).
	// It does nothing if the link does not exist.
	Unlink(N, N) error
	// Successors return the successors of a node (if exists)
	Successors(N) iter.Seq[N]
	// HasLink returns true if the link exists
	HasLink(N, N) bool
}

// DWGraph is the implementation for a graph with values linked to the edges.
// Values may be anything (a reducer will manage int or float values).
type DWGraph[N Node, V any] interface {
	// Graph defines the generic graph interface: manage nodes.
	Graph[N]
	// Link adds a value for a couple of nodes.
	// It changes the value if the link already exists.
	Link(N, N, V) error
	// Unlink nodes (if exists).
	// No change if the link does not exist.
	Unlink(N, N) error
	// Successors returns the weights and successors of a node (if exists)
	Successors(N) iter.Seq2[N, V]
	// HasLink returns the value and true if the link exists
	HasLink(N, N) (V, bool)
}
