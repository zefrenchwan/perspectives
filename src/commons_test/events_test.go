package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

type DummyObserver struct {
	// counter of incoming events
	Incoming int
	// counter of outgoing events
	Outgoing int
	// counter of errors
	Errors int
	// last executed processing
	LastReceived commons.EventProcessing
}

func (o *DummyObserver) OnEventProcessing(p commons.EventProcessing) {
	o.Incoming += 1
	if p.Error != nil {
		o.Errors++
	}

	o.Outgoing += len(p.Outgoings)
	o.LastReceived = p
}

func TestEventSource(t *testing.T) {
	structure := DummyStructure{id: "structure"}
	sevent := commons.NewEventLifetimeEnd(structure, time.Now())
	if sevent.Source() != structure {
		t.Fail()
	} else if !commons.IsEventComingFromStructure(sevent) {
		t.Fail()
	}
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
	} else if counter.LastReceived.Source.Id() != observer.Id() {
		t.Fail()
	}
}
