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

func TestGraphRemove(t *testing.T) {
	graph := structures.NewDVGraph[string, int]()
	graph.Link("a", "b", 10)
	graph.AddNode("d")

	if !graph.RemoveNode("d") {
		t.Log("node should be here")
		t.Fail()
	}

	if graph.RemoveNode("d") {
		t.Log("node should NOT be here (was removed)")
		t.Fail()
	}

	if !graph.RemoveNode("b") {
		t.Log("b should be here")
		t.Fail()
	} else if !graph.Has("a") {
		t.Log("a was removed due to link")
		t.Fail()
	}

}

func TestNeighbors(t *testing.T) {
	graph := structures.NewDVGraph[string, int]()
	graph.Link("a", "b", 10)
	graph.Link("a", "c", 100)
	graph.Link("c", "d", 1000)
	graph.AddNode("z")

	// not exists
	if values, found := graph.Neighbors("x"); found || len(values) != 0 {
		t.Log("x not in graph")
		t.Fail()
	}

	if values, found := graph.Edges("x"); found || len(values) != 0 {
		t.Log("x not in graph")
		t.Fail()
	}

	// exists
	if values, found := graph.Neighbors("a"); !found || len(values) != 2 {
		t.Log("a in graph with two childs")
		t.Fail()
	} else if values["b"] != 10 {
		t.Fail()
	} else if values["c"] != 100 {
		t.Fail()
	}

	if values, found := graph.Edges("c"); !found || len(values) != 1 {
		t.Log("c in graph with one child")
		t.Fail()
	} else {
		value := values[0]
		if value.Source != "c" {
			t.Fail()
		} else if value.Destination != "d" {
			t.Fail()
		} else if value.Value != 1000 {
			t.Fail()
		}
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

func TestWalkNoCycle(t *testing.T) {
	graph := structures.NewDVGraph[string, int]()
	graph.AddNode("a")
	graph.AddNode("b")
	graph.Link("b", "c", 10)
	graph.Link("c", "d", 10)

	var collector []string
	graph.Walk("a", func(source string) {
		collector = append(collector, source)
	})

	if slices.Compare([]string{"a"}, collector) != 0 {
		t.Log("failed to read single element")
		t.Fail()
	}

	collector = nil
	graph.Walk("b", func(source string) {
		collector = append(collector, source)
	})

	slices.Sort(collector)
	if slices.Compare([]string{"b", "c", "d"}, collector) != 0 {
		t.Log("walk failed for path")
		t.Log(collector)
		t.Log(graph)
		t.Fail()
	}

	collector = nil
	graph.Walk("c", func(source string) {
		collector = append(collector, source)
	})

	slices.Sort(collector)
	if slices.Compare([]string{"c", "d"}, collector) != 0 {
		t.Log("walk failed for path")
		t.Log(collector)
		t.Log(graph)
		t.Fail()
	}
}

func TestWalkCycle(t *testing.T) {
	graph := structures.NewDVGraph[string, int]()
	graph.Link("b", "c", 10)
	graph.Link("c", "d", 10)
	graph.Link("d", "a", 10)
	graph.Link("d", "b", 10)

	var collector []string
	graph.Walk("b", func(source string) {
		collector = append(collector, source)
	})

	slices.Sort(collector)
	if slices.Compare([]string{"a", "b", "c", "d"}, collector) != 0 {
		t.Log("failed to read cycle")
		t.Fail()
	}
}

func TestWalkReverse(t *testing.T) {
	graph := structures.NewDVGraph[string, int]()
	graph.Link("b", "c", 10)
	graph.Link("a", "c", 10)
	graph.Link("c", "d", 10)

	var collector []string
	graph.ReverseWalk("c", func(current string) { collector = append(collector, current) })
	slices.Sort(collector)
	if slices.Compare(collector, []string{"a", "b", "c"}) != 0 {
		t.Log("failed to reverse graph")
		t.Fail()
	}
}

func TestAddWithoutCycle(t *testing.T) {
	graph := structures.NewDVGraph[string, int]()
	if !graph.LinkWithoutCycle("a", "b", 10) {
		t.Log("cycle detected, while there is none")
		t.Fail()
	}

	if !graph.LinkWithoutCycle("a", "c", 10) {
		t.Log("cycle detected, while there is none")
		t.Fail()
	}

	if !graph.LinkWithoutCycle("b", "c", 10) {
		t.Log("cycle detected, while there is none")
		t.Fail()
	}

	// now make the cycle
	if graph.LinkWithoutCycle("c", "b", 10) {
		t.Log("cycle undetected, while there is b -> c -> b")
		t.Fail()
	}

	// test the rollback
	if !graph.Has("a") || !graph.Has("b") || !graph.Has("c") {
		t.Log("missing node")
		t.Fail()
	}

	// expecting :
	// a -> c, a -> b, b -> c
	if value, found := graph.Neighbors("c"); !found || len(value) != 0 {
		t.Log("neighbor failed for c")
		t.Fail()
	}

	if value, found := graph.Neighbors("b"); !found || len(value) != 1 {
		t.Log("neighbor failed for b")
		t.Fail()
	}

	if value, found := graph.Neighbors("a"); !found || len(value) != 2 {
		t.Log("neighbor failed for a")
		t.Fail()
	}
}

func TestEdgesWalk(t *testing.T) {
	graph := structures.NewDVGraph[string, int]()
	graph.AddNode("a")
	graph.Link("b", "c", 10)
	graph.Link("c", "d", 10)
	graph.Link("d", "e", 10)
	graph.Link("d", "f", 10)

	expected := []structures.GraphEdge[string, int]{
		{Source: "b", Destination: "c", Value: 10},
		{Source: "c", Destination: "d", Value: 10},
		{Source: "d", Destination: "e", Value: 10},
		{Source: "d", Destination: "f", Value: 10},
	}

	got := graph.EdgesFrom("b")

	if len(got) != len(expected) {
		t.Log("missing edges")
		t.Log(got)
		t.Fail()
	}

	for _, link := range expected {
		if !slices.Contains(got, link) {
			t.Log("missing link")
			t.Log(link)
			t.Fail()
		}
	}
}
