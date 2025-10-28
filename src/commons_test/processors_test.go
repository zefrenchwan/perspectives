package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestEventProcessorObservation(t *testing.T) {

	structure := DummyStructure{id: "structure"}
	source := commons.NewEventLifetimeEnd(structure, time.Now())
	dest := commons.NewEventTick(structure)

	mapper := func(e commons.Event) ([]commons.Event, error) {
		return []commons.Event{dest}, nil
	}

	processor := commons.NewEventProcessor(mapper)

	if values, err := processor.Process(source); err != nil {
		t.Fail()
	} else if len(values) != 1 {
		t.Fail()
	} else if values[0] != dest {
		t.Fail()
	}
}
