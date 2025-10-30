package commons_test

import (
	"maps"
	"slices"
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
