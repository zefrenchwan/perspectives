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

func TestGraphsErrors(t *testing.T) {
	node := graphs.IdentifiableNode{Identifer: "node"}
	other := graphs.IdentifiableNode{Identifer: "other"}
	last := graphs.IdentifiableNode{Identifer: "last"}

	gr := graphs.NewDGraph[graphs.IdentifiableNode]()
	gr.SetNode(node)
	if err := gr.Link(node, other); err == nil {
		t.Log("node failure : failed to add link with missing node")
		t.Fail()
	} else if err := gr.Link(last, node); err == nil {
		t.Log("node failure : failed to add link with missing node")
		t.Fail()
	} else if err := gr.Unlink(node, last); err == nil {
		t.Log("node failure : failed to remove link with missing node")
		t.Fail()
	} else if err := gr.Unlink(other, node); err == nil {
		t.Log("node failure : failed to remove link with missing node")
		t.Fail()
	}
}

func TestValuedGraph(t *testing.T) {
	node := graphs.IdentifiableNode{Identifer: "node"}
	other := graphs.IdentifiableNode{Identifer: "other"}
	last := graphs.IdentifiableNode{Identifer: "last"}

	g := graphs.NewDWGraph[graphs.IdentifiableNode, int]()
	g.SetNode(node)
	g.SetNode(other)
	g.Link(node, other, 1)
	if v, has := g.HasLink(node, other); !has {
		t.Log("link failure : failed to add link")
		t.Fail()
	} else if v != 1 {
		t.Log("link failure : failed to add link with correct value")
		t.Fail()
	} else if _, has := g.HasLink(other, node); has {
		t.Log("link failure : should not have link")
		t.Fail()
	}

	// go for walkthrough from an existing node
	for v, e := range g.Successors(node) {
		if v.Id() != other.Id() {
			t.Log("walkthrough failure : failed to walk through successors")
			t.Fail()
		} else if e != 1 {
			t.Log("walkthrough failure : failed to walk through successors with correct value")
			t.Fail()
		}
	}

	// walkt through successors from a unlinked node
	for v, e := range g.Successors(other) {
		t.Log("walkthrough failure : other has no successors", v.Id(), e)
		t.Fail()
	}

	// walk from node not in graph
	for v, e := range g.Successors(last) {
		t.Log("walkthrough failure : not in graph", v.Id(), e)
		t.Fail()
	}
}
