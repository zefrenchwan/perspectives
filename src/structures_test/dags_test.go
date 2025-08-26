package structures_test

import (
	"maps"
	"testing"

	"github.com/zefrenchwan/perspectives.git/structures"
)

func TestDAGAddNoCycle(t *testing.T) {
	dag := structures.NewDAG[string, int]()
	if !dag.AddNode("Paris") {
		t.Log("failed to add non existing element")
		t.Fail()
	} else if err := dag.Link("Paris", "Montpellier", 800); err != nil {
		t.Log("failed to link element, found a cycle")
		t.Log(err)
		t.Fail()
	}

	// test neighbors
	if values, found := dag.Neighbors("Paris"); !found {
		t.Log("failed to find neighbors")
		t.Fail()
	} else if len(values) != 1 {
		t.Log("failed to find link to Montpellier")
		t.Fail()
	} else if values["Montpellier"] != 800 {
		t.Log("invalid link value")
		t.Fail()
	}

	// test incoming neighbors
	if values, found := dag.Neighbors("Montpellier"); !found {
		t.Log("failed to find neighbors")
		t.Fail()
	} else if len(values) != 0 {
		t.Log("no link to Paris expected")
		t.Fail()
	}
}

func TestDAGAddWithCycle(t *testing.T) {
	dag := structures.NewDAG[string, int]()
	if err := dag.Link("A", "B", 10); err != nil {
		t.Log("no cycle expected")
		t.Fail()
	} else if err := dag.Link("B", "C", 100); err != nil {
		t.Log("no cycle expected")
		t.Fail()
	} else if err := dag.Link("C", "A", 1000); err == nil {
		t.Log("no cycle detected")
		t.Fail()
	}

	// test rollback
	if len(dag) != 3 {
		t.Log("failed to rollback")
		t.Fail()
	}

	expected := make(structures.DAG[string, int])
	expected["A"] = map[string]int{"B": 10}
	expected["B"] = map[string]int{"C": 100}
	expected["C"] = map[string]int{}

	for k := range expected {
		if !maps.Equal(dag[k], expected[k]) {
			t.Logf("comparison failed for key %s", k)
			t.Log(dag[k])
			t.Fail()
		}
	}
}
