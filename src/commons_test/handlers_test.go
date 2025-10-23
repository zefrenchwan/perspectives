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

func TestActiveObjectHandler(t *testing.T) {
	structure := DummyStructure{id: "test"}
	source := commons.NewEventLifetimeEnd(structure, time.Now())
	event := commons.NewEventTick(structure)
	mapper := func(e commons.Event) ([]commons.Event, error) {
		return []commons.Event{event}, nil
	}

	now := time.Now().Truncate(time.Hour)
	object := commons.NewActiveObjectHandler[int](now, nil, commons.NewEventProcessor(mapper))

	if result, err := object.Process(source); err != nil {
		t.Fail()
	} else if len(result) != 1 {
		t.Fail()
	} else if result[0].Id() != event.Id() {
		t.Fail()
	}

	if !commons.NewPeriodSince(now, true).Equals(object.ActivePeriod()) {
		t.Fail()
	}

	object.SetActivePeriod(commons.NewFullPeriod())
	if !commons.NewFullPeriod().Equals(object.ActivePeriod()) {
		t.Fail()
	}
}

func TestActiveObjectHandlerState(t *testing.T) {
	structure := DummyStructure{id: "test"}
	event := commons.NewEventTick(structure)
	mapper := func(e commons.Event) ([]commons.Event, error) {
		return []commons.Event{event}, nil
	}

	now := time.Now().Truncate(time.Hour)
	values := map[string]int{"a": 10, "b": 100}
	object := commons.NewActiveObjectHandler(now, values, commons.NewEventProcessor(mapper))

	if result := object.Read().Values(); len(result) != 2 {
		t.Fail()
	} else if result["a"] != 10 {
		t.Fail()
	} else if result["b"] != 100 {
		t.Fail()
	}

	object.Remove("a")
	if result := object.Read().Values(); len(result) != 1 {
		t.Fail()
	} else if result["b"] != 100 {
		t.Fail()
	}
}

func TestActiveObjectHandlerObservable(t *testing.T) {
	structure := DummyStructure{id: "test"}
	source := commons.NewEventLifetimeEnd(structure, time.Now())
	event := commons.NewEventTick(structure)
	mapper := func(e commons.Event) ([]commons.Event, error) {
		return []commons.Event{event}, nil
	}

	object := commons.NewActiveObjectHandler[int](time.Now(), nil, commons.NewEventProcessor(mapper))
	observer := DummyObserver{}
	object.AddObserver(&observer)

	object.Process(source)
	if observer.Incoming != 1 {
		t.Fail()
	} else if observer.Outgoing != 1 {
		t.Fail()
	} else if observer.Errors != 0 {
		t.Fail()
	}
}
