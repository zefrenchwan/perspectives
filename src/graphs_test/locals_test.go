package graphs_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/graphs"
)

func TestLinkPassing(t *testing.T) {
	g := graphs.NewDGraph[graphs.IdentifiableNode]()
	a, b := graphs.IdentifiableNode{Identifer: "a"}, graphs.IdentifiableNode{Identifer: "b"}
	g.SetNode(a)
	g.SetNode(b)
	g.Link(a, b)

	identifiers := make([]string, 0, 2)
	elements := make([]graphs.IdentifiableNode, 0, 2)
	for element := range g.Nodes() {
		identifiers = append(identifiers, element.Identifer)
	}

	slices.Sort(identifiers)
	if !slices.Equal(identifiers, []string{"a", "b"}) {
		t.Log("nodes in graph failure")
		t.Fail()
	}

	elements = slices.Collect(g.Successors(a))
	if len(elements) != 1 || elements[0].Identifer != "b" {
		t.Log("successors of a failure")
		t.Fail()
	}

	elements = slices.Collect(g.Successors(b))
	if len(elements) != 0 {
		t.Log("successors of b failure")
		t.Fail()
	}

	if !g.HasLink(a, b) {
		t.Log("link failure : should have link")
		t.Fail()
	} else if g.HasLink(b, a) {
		t.Log("link failure : no link expected")
		t.Fail()
	}

	c := graphs.IdentifiableNode{Identifer: "c"}
	if g.HasNode(c) {
		t.Log("node failure : should not have node")
		t.Fail()
	} else if !g.HasNode(a) {
		t.Log("node failure : should have node")
		t.Fail()
	}
}

func TestRemoveElements(t *testing.T) {
	g := graphs.NewDGraph[graphs.IdentifiableNode]()
	a, b := graphs.IdentifiableNode{Identifer: "a"}, graphs.IdentifiableNode{Identifer: "b"}
	g.SetNode(a)
	g.SetNode(b)
	g.Link(a, b)
	g.Link(b, a)

	if !g.HasLink(a, b) {
		t.Log("link failure : should have link")
		t.Fail()
	}

	g.Unlink(a, b)
	if g.HasLink(a, b) {
		t.Log("link failure : failed to remove link")
		t.Fail()
	} else if !g.HasLink(b, a) {
		t.Log("link failure : failed to keep other link")
		t.Fail()
	}

	g.RemoveNode(a)
	if g.HasNode(a) {
		t.Log("node failure : should not have node")
		t.Fail()
	} else if !g.HasNode(b) {
		t.Log("node failure : should have node")
		t.Fail()
	} else if g.HasLink(a, b) {
		t.Log("link failure : should not have link")
		t.Fail()
	} else if g.HasLink(b, a) {
		t.Log("link failure : should not have link")
		t.Fail()
	}
}
