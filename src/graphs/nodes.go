package graphs

import "github.com/zefrenchwan/perspectives.git/commons"

type Node interface {
	commons.Identifiable
}

type GenericNode[V any] struct {
	Identifier string
	Value      V
}

func (n *GenericNode[V]) Id() string {
	return n.Identifier
}
