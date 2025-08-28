package structures_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/structures"
)

func TestGraphAdd(t *testing.T) {
	graph := structures.NewDVGraph[string, int]()
	graph.Link("a", "b", 10)
	if graph.AddNode("a") {
		t.Log("should return false because node exists")
		t.Fail()
	}

	if !graph.AddNode("c") {
		t.Log("should return true because node did not exist")
		t.Fail()
	}

	nodes := graph.Nodes()
	slices.Sort(nodes)

	if slices.Compare(nodes, []string{"a", "b", "c"}) != 0 {
		t.Log("nodes missing")
		t.Fail()
	}
}

func TestCycles(t *testing.T) {
	graph := structures.NewDVGraph[string, int]()
	graph.Link("a", "b", 10)
	graph.Link("b", "c", 10)
	graph.Link("c", "d", 10)

	if graph.HasCycle() {
		t.Log("found cycle whereas there is none")
		t.Fail()
	}

	graph.Link("c", "a", 100)
	if !graph.HasCycle() {
		t.Log("did not find cycle but there is a -> b -> c -> a")
		t.Fail()
	}
}
