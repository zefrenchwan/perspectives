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
