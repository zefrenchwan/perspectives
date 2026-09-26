package graphs

import "github.com/zefrenchwan/perspectives.git/commons"

type Node interface {
	commons.Identifiable
}

type IdentifiableNode struct {
	Identifer string
}

func (n IdentifiableNode) Id() string {
	return n.Identifer
}
