package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestEventTick(t *testing.T) {
	structure := DummyStructure{id: "structure"}
	now := time.Now().Truncate(commons.TIME_PRECISION)
	tick := commons.NewEventTickTime(structure, now)
	if !tick.ProcessingTime().Equal(now) {
		t.Fail()
	}
}

func TestEventMessage(t *testing.T) {
	structure := DummyStructure{id: "structure"}
	now := time.Now().Truncate(commons.TIME_PRECISION)
	tick := commons.NewEventTickTime(structure, now)

	processor := commons.NewEventProcessor(func(e commons.Event) ([]commons.Event, error) { return nil, nil })
	message := commons.NewMessage(tick, processor)

	if message.Source().Id() != structure.id {
		t.Fail()
	} else if message.Destination().Id() != processor.Id() {
		t.Fail()
	}
}
