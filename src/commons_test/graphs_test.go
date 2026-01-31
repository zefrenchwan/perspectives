package commons_test

import (
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
