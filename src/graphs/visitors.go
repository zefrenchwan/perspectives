package graphs

// ValuesVisitor is a visitor over nodes and valued links in a graph.
type ValuesVisitor[N Node, V any] interface {
	// OnNode reacts on a node in the graph.
	OnNode(N)
	// OnLink reacts on a link in the graph with value.
	OnLink(N, N, V)
}

// Visitor is a visitor over nodes in a graph.
type Visitor[N Node] interface {
	// OnNode reacts on a node in the graph.
	OnNode(N)
	// OnLink reacts on a link in the graph.
	OnLink(N, N)
}
