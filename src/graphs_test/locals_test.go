package graphs_test

import (
	"fmt"
	"testing"

	"github.com/zefrenchwan/perspectives.git/graphs"
)

func TestAddLoad(t *testing.T) {
	g := graphs.NewDGraph[graphs.IdentifiableNode]()
	a, b := graphs.IdentifiableNode{Identifer: "a"}, graphs.IdentifiableNode{Identifer: "b"}
	g.SetNode(a)
	g.SetNode(b)
	g.Link(a, b)

	for element := range g.Nodes() {
		fmt.Println(element.Identifer)
	}
}
