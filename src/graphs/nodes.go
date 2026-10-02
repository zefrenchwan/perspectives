package graphs

import "github.com/zefrenchwan/perspectives.git/commons"

// Node is the general contract for a vertex in a graph
type Node interface {
	commons.Identifiable
}

type IdentifiableNode struct {
	Identifier string
}

func (n IdentifiableNode) Id() string {
	return n.Identifier
}
