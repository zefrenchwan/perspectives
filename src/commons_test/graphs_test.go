package commons_test

import (
	"maps"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestTemporalGraph(t *testing.T) {
	source := commons.NewModelObject()
	dest := commons.NewModelObject()
	sink := commons.NewModelObject()
	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(10, 0, 0)
	graph := commons.NewDynamicGraph[commons.ModelObject, int]()

	// test empty => iterate but nothing (test if ends)
	if result := maps.Collect(graph.Neighbors(source, now)); len(result) != 0 {
		t.Fail()
	} else if values := slices.Collect(graph.Vertices()); len(values) != 0 {
		t.Fail()
	}

	// test one element at different times
	graph.Relate(source, dest, 10, commons.NewPeriodSince(now, true))
	if result := maps.Collect(graph.Neighbors(source, now)); len(result) != 1 {
		t.Fail()
	} else if result[dest] != 10 {
		t.Fail()
	} else if result := maps.Collect(graph.Neighbors(dest, now)); len(result) != 0 {
		t.Fail()
	} else if result := maps.Collect(graph.Neighbors(source, before)); len(result) != 0 {
		t.Fail()
	} else if values := slices.Collect(graph.Vertices()); len(values) != 2 {
		t.Fail()
	} else if !slices.Contains(values, source) {
		t.Fail()
	} else if !slices.Contains(values, dest) {
		t.Fail()
	}

	// delete the unique edge
	graph.Remove(source, dest)
	if result := maps.Collect(graph.Neighbors(source, now)); len(result) != 0 {
		t.Fail()
	} else if values := slices.Collect(graph.Vertices()); len(values) != 2 {
		t.Fail()
	} else if !slices.Contains(values, source) {
		t.Fail()
	} else if !slices.Contains(values, dest) {
		t.Fail()
	}

	// test with two elements
	graph.Relate(source, dest, 10, commons.NewPeriodSince(before, true))
	graph.Relate(source, sink, 100, commons.NewPeriodSince(now, true))
	if result := maps.Collect(graph.Neighbors(source, now)); len(result) != 2 {
		t.Fail()
	} else if result[dest] != 10 {
		t.Fail()
	} else if result[sink] != 100 {
		t.Fail()
	} else if values := slices.Collect(graph.Vertices()); len(values) != 3 {
		t.Fail()
	} else if !slices.Contains(values, source) {
		t.Fail()
	} else if !slices.Contains(values, dest) {
		t.Fail()
	} else if !slices.Contains(values, sink) {
		t.Fail()
	}

	// test find
	if _, found := graph.Lookup("not there"); found {
		t.Fail()
	} else if value, found := graph.Lookup(source.Id()); !found {
		t.Fail()
	} else if value != source {
		t.Fail()
	}

	// test time management
	graph = commons.NewDynamicGraph[commons.ModelObject, int]()
	graph.Relate(source, dest, 10, commons.NewFinitePeriod(before, now, true, false))
	if result := maps.Collect(graph.Neighbors(source, before)); len(result) != 1 {
		t.Fail()
	} else if result[dest] != 10 {
		t.Fail()
	} else if result := maps.Collect(graph.Neighbors(source, now)); len(result) != 0 {
		t.Fail()
	} else if result := maps.Collect(graph.Neighbors(source, after)); len(result) != 0 {
		t.Fail()
	} else if values := slices.Collect(graph.Vertices()); len(values) != 2 {
		t.Fail()
	} else if !slices.Contains(values, source) {
		t.Fail()
	} else if !slices.Contains(values, dest) {
		t.Fail()
	}

	// Test set
	graph = commons.NewDynamicGraph[commons.ModelObject, int]()
	graph.Set(source)
	if result := maps.Collect(graph.Neighbors(source, before)); len(result) != 0 {
		t.Fail()
	} else if v, found := graph.Lookup(source.Id()); !found {
		t.Fail()
	} else if v.Id() != source.Id() {
		t.Fail()
	}
}

func TestTemporalWalker(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-10, 0, 0)
	period := commons.NewPeriodSince(now, true)
	source := DummyIdBasedImplementation{id: "0"}
	dest := DummyIdBasedImplementation{id: "1"}
	other := DummyIdBasedImplementation{id: "2"}
	sink := DummyIdBasedImplementation{id: "3"}
	graph := commons.NewDynamicGraph[DummyIdBasedImplementation, int]()
	graph.Relate(source, dest, 5, period)
	graph.Relate(source, other, 50, period)
	graph.Relate(dest, sink, 500, period)
	graph.Relate(other, sink, 5000, period)

	// test if neighbors work NOW
	if values := maps.Collect(graph.Neighbors(source, now)); len(values) != 2 {
		t.Log("failed neighbors at now")
		t.Fail()
	} else if values[dest] != 5 {
		t.Log("failed neighbors at now")
		t.Fail()
	} else if values[other] != 50 {
		t.Log("failed neighbors at now")
		t.Fail()
	} else if values := maps.Collect(graph.Neighbors(dest, now)); len(values) != 1 {
		t.Log("failed neighbors at now")
		t.Fail()
	} else if values[sink] != 500 {
		t.Log("failed neighbors at now")
		t.Fail()
	} else if values := maps.Collect(graph.Neighbors(other, now)); len(values) != 1 {
		t.Log("failed neighbors at now")
		t.Fail()
	} else if values[sink] != 5000 {
		t.Log("failed neighbors at now")
		t.Fail()
	} else if values := maps.Collect(graph.Neighbors(sink, now)); len(values) != 0 {
		t.Log("failed neighbors at now")
		t.Fail()
	}

	// test if neighbors work before (should have nothing)
	if values := maps.Collect(graph.Neighbors(source, before)); len(values) != 0 {
		t.Log("failed neighbors at before")
		t.Fail()
	} else if values := maps.Collect(graph.Neighbors(dest, before)); len(values) != 0 {
		t.Log("failed neighbors at before")
		t.Fail()
	} else if values := maps.Collect(graph.Neighbors(other, before)); len(values) != 0 {
		t.Log("failed neighbors at before")
		t.Fail()
	} else if values := maps.Collect(graph.Neighbors(sink, before)); len(values) != 0 {
		t.Log("failed neighbors at before")
		t.Fail()
	}

	// test full walk
	var values []string
	walker := commons.NewDynamicGraphWalker(graph, source, now)
	for walker.Next() {
		composite := walker.Source().Id() + ";" + walker.Position().Id() + ";" + strconv.Itoa(walker.SourceEdge())
		values = append(values, composite)
	}

	expected := []string{"0;2;50", "0;1;5", "2;3;5000", "1;3;500"}
	slices.Sort(expected)
	slices.Sort(values)

	if len(values) != len(expected) {
		t.Log("failed full walk")
		t.Fail()
	} else if slices.Compare(expected, values) != 0 {
		t.Log("failed full walk")
		t.Fail()
	}

	// stop at second step
	values = nil
	walker = commons.NewDynamicGraphWalker(graph, source, now)
	for walker.Next() {
		composite := walker.Source().Id() + ";" + walker.Position().Id() + ";" + strconv.Itoa(walker.SourceEdge())
		values = append(values, composite)
		if len(values) >= 2 {
			walker.Stop()
		}
	}

	expected = []string{"0;2;50", "0;1;5"}
	slices.Sort(expected)
	slices.Sort(values)

	if len(values) != len(expected) {
		t.Log("failed stopped walk")
		t.Fail()
	} else if slices.Compare(expected, values) != 0 {
		t.Log("failed stopped walk")
		t.Fail()
	}

	// Walk before (should not walk)
	values = nil
	walker = commons.NewDynamicGraphWalker(graph, source, before)
	for walker.Next() {
		t.Log("NO active edge, should fail")
		t.Fail()
	}
}

func TestTemporalQuery(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	period := commons.NewPeriodSince(now, true)
	source := DummyIdBasedImplementation{id: "0"}
	dest := DummyIdBasedImplementation{id: "1"}
	other := DummyIdBasedImplementation{id: "2"}
	sink := DummyIdBasedImplementation{id: "3"}
	graph := commons.NewDynamicGraph[DummyIdBasedImplementation, int]()
	graph.Relate(source, dest, 5, period)
	graph.Relate(source, other, 50, period)
	graph.Relate(dest, sink, 500, period)
	graph.Relate(other, sink, 5000, period)

	acceptor5000 := func(a, b DummyIdBasedImplementation, edge int) bool {
		return edge >= 5000
	}

	// test on last part of the graph
	if value := maps.Collect(commons.DynamicGraphLocalQuery(graph, source, now, acceptor5000)); len(value) != 1 {
		t.Log("failed edge predicate")
		t.Fail()
	} else if value[other] != sink {
		t.Log("failed edge predicate")
		t.Fail()
	}

	// test on false
	if value := maps.Collect(commons.DynamicGraphLocalQuery(graph, source, now, func(a, b DummyIdBasedImplementation, v int) bool { return false })); len(value) != 0 {
		t.Log("failed false predicate")
		t.Fail()
	}
}

func TestTemporalAction(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	period := commons.NewPeriodSince(now, true)
	source := commons.NewStateObject[int]()
	dest := commons.NewStateObject[int]()
	other := commons.NewStateObject[int]()
	sink := commons.NewStateObject[int]()
	graph := commons.NewDynamicGraph[*commons.StateObject[int], int]()
	graph.Relate(source, dest, 5, period)
	graph.Relate(source, other, 50, period)
	graph.Relate(dest, sink, 500, period)
	graph.Relate(other, sink, 5000, period)

	acceptor5000 := func(a, b *commons.StateObject[int], edge int) bool {
		return edge >= 5000
	}

	processor := func(o *commons.StateObject[int]) error {
		o.SetValue("validated", 100)
		return nil
	}

	if err := commons.DynamicGraphLocalAction(graph, source, now, acceptor5000, processor); err != nil {
		t.Log(err)
		t.Fail()
	} else if value, found := sink.GetValue("validated"); !found {
		t.Fail()
	} else if value != 100 {
		t.Fail()
	}

}
