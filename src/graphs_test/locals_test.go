package graphs_test

import "testing"

// Dummy test to ensure maps work as expected : clean in inner map does not mean
// pushing it back to the outer map by key. Change is visible in outer map.
func TestMapDeleteElement(t *testing.T) {
	v := make(map[string]map[string]int)
	v["test"] = make(map[string]int)
	v["test"]["other"] = 1

	for _, currentMap := range v {
		delete(currentMap, "other")
	}

	res := len(v["test"])
	if res != 0 {
		t.Errorf("Expected map to be empty, but it has %d elements", res)
	}
}
