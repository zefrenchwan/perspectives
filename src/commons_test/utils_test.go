package commons_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestSliceReduce(t *testing.T) {
	values := []int{0, 5, 10, 4, 15, 10, 10, 10}
	expected := []int{0, 4, 5, 10, 15}
	if result := commons.SliceReduce(values); slices.Compare(result, expected) != 0 {
		t.Fail()
	}
}

func TestSliceDeduplicate(t *testing.T) {
	values := []int{5, 10, 5, 10, 10, 10}
	if result := commons.SliceDeduplicate(values); len(result) != 2 {
		t.Fail()
	} else if !slices.Contains(result, 5) {
		t.Fail()
	} else if !slices.Contains(result, 10) {
		t.Fail()
	}
}

func TestSliceDeduplicateFunc(t *testing.T) {
	values := []int{5, 10, 5, 10, 10, 10}
	if result := commons.SliceDeduplicateFunc(values, func(a, b int) bool { return a == b }); len(result) != 2 {
		t.Fail()
	} else if !slices.Contains(result, 5) {
		t.Fail()
	} else if !slices.Contains(result, 10) {
		t.Fail()
	}
}

func TestSliceCommonElement(t *testing.T) {
	values := []int{0, 2, 4}
	noMatch := []int{1, 3, 5}
	match := []int{1, 2, 5}
	if commons.SliceCommonElement(values, noMatch) {
		t.Fail()
	} else if !commons.SliceCommonElement(values, match) {
		t.Fail()
	}
}

func TestSliceCommonElementFunc(t *testing.T) {
	values := []int{0, 2, 4}
	noMatch := []int{1, 3, 5}
	match := []int{1, 2, 5}
	if commons.SliceCommonElementFunc(values, noMatch, func(a, b int) bool { return a == b }) {
		t.Fail()
	} else if !commons.SliceCommonElementFunc(values, match, func(a, b int) bool { return a == b }) {
		t.Fail()
	}
}

func TestSlicesEqualsAsSetsFunc(t *testing.T) {
	values := []int{0, 2, 4}
	noMatch := []int{0, 2, 5}
	match := []int{4, 2, 0}
	if commons.SlicesEqualsAsSetsFunc(values, noMatch, func(a, b int) bool { return a == b }) {
		t.Fail()
	} else if !commons.SlicesEqualsAsSetsFunc(values, match, func(a, b int) bool { return a == b }) {
		t.Fail()
	}
}

func TestMapsReverseFind(t *testing.T) {
	values := make(map[string][]string)
	values["noMatch"] = []string{"a", "b"}
	values["match"] = []string{"matching"}
	if result := commons.MapsReverseFind(values, "matching"); len(result) != 1 {
		t.Fail()
	} else if result[0] != "match" {
		t.Fail()
	}
}

func TestSlicesContainsAllFunc(t *testing.T) {
	values := []int{0, 2, 4}
	noMatch := []int{0, 5}
	match := []int{2, 0}

	if !commons.SlicesContainsAllFunc(values, match, func(a, b int) bool { return a == b }) {
		t.Fail()
	} else if commons.SlicesContainsAllFunc(values, noMatch, func(a, b int) bool { return a == b }) {
		t.Fail()
	}
}

func TestSlicesFilter(t *testing.T) {
	expected := []int{4}
	if result := commons.SlicesFilter([]int{0, 2, 4}, func(a int) bool { return a >= 3 }); slices.Compare(expected, result) != 0 {
		t.Fail()
	}
}
