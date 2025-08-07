package models_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
)

func TestAddWithoutCycle(t *testing.T) {
	value := models.NewDAG[string, int]()
	var index int
	if err := value.UpsertNodesLinks("source", "dest", 0); err != nil {
		index++
		t.Logf("Failed at position %d", index)
		t.Fail()
	}
	if err := value.UpsertNodesLinks("dest", "sink", 0); err != nil {
		index++
		t.Logf("Failed at position %d", index)
		t.Fail()
	}
	if err := value.UpsertNodesLinks("source", "other", 0); err != nil {
		index++
		t.Logf("Failed at position %d", index)
		t.Fail()
	}
	if err := value.UpsertNodesLinks("other", "sink", 0); err != nil {
		index++
		t.Logf("Failed at position %d", index)
		t.Fail()
	}
}

func TestAddCycle(t *testing.T) {
	value := models.NewDAG[string, int]()
	var index int
	if err := value.UpsertNodesLinks("source", "dest", 0); err != nil {
		index++
		t.Logf("Failed at position %d", index)
		t.Fail()
	}
	if err := value.UpsertNodesLinks("dest", "sink", 0); err != nil {
		index++
		t.Logf("Failed at position %d", index)
		t.Fail()
	}
	if err := value.UpsertNodesLinks("source", "other", 0); err != nil {
		index++
		t.Logf("Failed at position %d", index)
		t.Fail()
	}
	if err := value.UpsertNodesLinks("sink", "source", 0); err == nil {
		t.Log("Cycle was not detected")
		t.Fail()
	}
}
