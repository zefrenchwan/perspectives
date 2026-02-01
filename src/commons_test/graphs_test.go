package commons_test

import (
	"maps"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// StringEvent is a simple implementation of the Event interface for testing
type StringEvent struct {
	Value string
}

func TestGraphPropagationBlackBox(t *testing.T) {
	startTime := time.Now().Truncate(time.Second)
	graph := commons.NewGraph(startTime)

	// Set mappers : A produces, B stores
	mapperA := commons.NewEventMapper(func(events []commons.Event) []commons.Event {
		return []commons.Event{StringEvent{Value: "Pong"}}
	})

	var receivedByB []commons.Event
	mapperB := commons.NewEventMapper(func(events []commons.Event) []commons.Event {
		receivedByB = append(receivedByB, events...)
		return nil
	})

	// Topology: A --(2s)--> B
	graph.Set(mapperA, mapperB, 2*time.Second)

	// Emitting ping at start time
	graph.Emit(mapperA, StringEvent{Value: "Ping"}, startTime)

	// First step: processing ping, sending pong.
	// Current time is one second after
	graph.Step(1 * time.Second)
	if len(receivedByB) != 0 {
		t.Error("event received before delay")
	}

	// Second step: time is 2 seconds later, should read event
	graph.Step(1 * time.Second)
	if len(receivedByB) == 0 {
		t.Error("missing events")
	}
}

func TestGraphNeighbors(t *testing.T) {
	startTime := time.Now().Truncate(time.Second)
	graph := commons.NewGraph(startTime)

	mapperA := commons.NewEventMapper(func(events []commons.Event) []commons.Event { return nil })
	mapperB := commons.NewEventMapper(func(events []commons.Event) []commons.Event { return nil })
	mapperC := commons.NewEventMapper(func(events []commons.Event) []commons.Event { return nil })
	mapperD := commons.NewEventMapper(func(events []commons.Event) []commons.Event { return nil }) // Not in graph

	// Case 1: Element not in graph
	neighbors := maps.Collect(graph.Neighbors(mapperD))
	if len(neighbors) != 0 {
		t.Errorf("Expected 0 neighbors for mapperD (not in graph), got %v", neighbors)
	}

	// Topology: A --(1s)--> B, A --(2s)--> C
	graph.Set(mapperA, mapperB, 1*time.Second)
	graph.Set(mapperA, mapperC, 2*time.Second)
	graph.Set(mapperB, mapperC, 3*time.Second) // B is also a source

	// Case 2: Element is only a destination (mapperC has no outgoing links set)
	neighborsC := maps.Collect(graph.Neighbors(mapperC))
	if neighborsC == nil || len(neighborsC) != 0 {
		t.Errorf("Expected empty neighbors for mapperC (only destination), got %v", neighborsC)
	}

	// Case 3: Element is a source with multiple neighbors (mapperA)
	neighborsA := maps.Collect(graph.Neighbors(mapperA))
	if len(neighborsA) != 2 {
		t.Errorf("Expected 2 neighbors for mapperA, got %d", len(neighborsA))
	}
	if latency, ok := neighborsA[mapperB]; !ok || latency != 1*time.Second {
		t.Errorf("Expected neighbor mapperB with latency 1s, got %v", latency)
	}
	if latency, ok := neighborsA[mapperC]; !ok || latency != 2*time.Second {
		t.Errorf("Expected neighbor mapperC with latency 2s, got %v", latency)
	}

	// Test mapperB as source
	neighborsB := maps.Collect(graph.Neighbors(mapperB))
	if len(neighborsB) != 1 {
		t.Errorf("Expected 1 neighbor for mapperB, got %d", len(neighborsB))
	}
	if latency, ok := neighborsB[mapperC]; !ok || latency != 3*time.Second {
		t.Errorf("Expected neighbor mapperC with latency 3s, got %v", latency)
	}
}
