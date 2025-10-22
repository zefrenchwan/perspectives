package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

type DummyObserver struct {
	Incoming int
	Outgoing int
	Errors   int
}

func (o *DummyObserver) OnIncomingEvents(events []commons.Event) {
	o.Incoming += len(events)
}

func (o *DummyObserver) OnProcessingEvents(events []commons.Event, e error) {
	if e != nil {
		o.Errors++
	}

	o.Outgoing += len(events)
}

func TestObjectEventProcessor(t *testing.T) {
	structure := DummyStructure{id: "test"}
	event := commons.NewEventTick(structure)
	mapper := func(events []commons.Event) ([]commons.Event, error) {
		return []commons.Event{event}, nil
	}

	object := commons.NewEventObservableProcessorFromProcessor(commons.NewEventProcessor(mapper))

	if result, err := object.Process(nil); err != nil {
		t.Fail()
	} else if len(result) != 1 {
		t.Fail()
	} else if result[0].Id() != event.Id() {
		t.Fail()
	}

	observer := DummyObserver{}
	object.AddObserver(&observer)

	object.Process([]commons.Event{event})
	if observer.Incoming != 1 {
		t.Fail()
	} else if observer.Outgoing != 1 {
		t.Fail()
	} else if observer.Errors != 0 {
		t.Fail()
	}
}
