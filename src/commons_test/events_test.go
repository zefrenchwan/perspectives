package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

type DummyObserver struct {
	Incoming int
	Outgoing int
	Errors   int
}

func (o *DummyObserver) OnIncomingEvent(event commons.Event) {
	o.Incoming += 1
}

func (o *DummyObserver) OnProcessingEvents(events []commons.Event, e error) {
	if e != nil {
		o.Errors++
	}

	o.Outgoing += len(events)
}

func TestEventObservableProcessor(t *testing.T) {

	structure := DummyStructure{id: "structure"}
	source := commons.NewEventLifetimeEnd(structure, time.Now())
	dest := commons.NewEventTick(structure)

	mapper := func(e commons.Event) ([]commons.Event, error) {
		return []commons.Event{dest}, nil
	}

	processor := commons.NewEventProcessor(mapper)
	observer := commons.NewEventObservableProcessor(processor)

	counter := &DummyObserver{}
	observer.AddObserver(counter)

	if values, err := observer.Process(source); err != nil {
		t.Fail()
	} else if len(values) != 1 {
		t.Fail()
	} else if values[0] != dest {
		t.Fail()
	} else if counter.Errors != 0 {
		t.Fail()
	} else if counter.Incoming != 1 {
		t.Fail()
	} else if counter.Outgoing != 1 {
		t.Fail()
	}
}
