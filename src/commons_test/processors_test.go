package commons_test

import (
	"errors"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

type DummyObserver struct {
	// to implement identifiable
	id string
	// counter of incoming events
	Incoming int
	// counter of outgoing events
	Outgoing int
	// counter of errors
	Errors int
	// LastSource contains the last identifiable that emitted
	LastSource commons.Identifiable
}

func (o *DummyObserver) Id() string {
	return o.id
}

func (o *DummyObserver) OnEventProcessing(source commons.EventProcessor, in commons.Event, out []commons.Event, e error) {
	o.Incoming += 1
	if e != nil {
		o.Errors++
	}

	o.Outgoing += len(out)
	o.LastSource = source
}

func TestEventProcessorObservation(t *testing.T) {

	structure := DummyStructure{id: "structure"}
	source := commons.NewEventLifetimeEnd(structure, time.Now())
	dest := commons.NewEventTick(structure)

	mapper := func(e commons.Event) ([]commons.Event, error) {
		return []commons.Event{dest}, nil
	}

	processor := commons.NewEventProcessor(mapper)

	counter := &DummyObserver{id: "observer"}
	processor.AddObserver(counter)

	if values, err := processor.Process(source); err != nil {
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
	} else if counter.LastSource.Id() != processor.Id() {
		t.Fail()
	}

	if obs := processor.Observers(); len(obs) != 1 {
		t.Fail()
	} else if obs[0].Id() != counter.Id() {
		t.Fail()
	}

	processor.RemoveObservers(func(eo commons.EventObserver) bool { return eo.Id() != counter.id })
	if obs := processor.Observers(); len(obs) != 1 {
		t.Log("should keep observers intact because we EXCLUDE counter")
		t.Fail()
	} else if obs[0].Id() != counter.Id() {
		t.Fail()
	}

	processor.RemoveObservers(func(eo commons.EventObserver) bool { return eo.Id() == counter.id })
	if obs := processor.Observers(); len(obs) != 0 {
		t.Log("should remove counter because of the predicate")
		t.Fail()
	}
}

func TestEventRedirection(t *testing.T) {
	structure := DummyStructure{id: "structure"}
	event := commons.NewEventLifetimeEnd(structure, time.Now())
	response := commons.NewEventTick(structure)

	processor := commons.NewEventProcessor(func(e commons.Event) ([]commons.Event, error) {
		return nil, errors.ErrUnsupported
	})

	catcher := commons.NewEventProcessor(func(e commons.Event) ([]commons.Event, error) {
		return []commons.Event{response}, nil
	})

	mapper := commons.NewEventRedirection(catcher, processor, func(e commons.Event) bool { return true })

	if result, err := mapper.Process(event); err != nil {
		t.Fail()
	} else if len(result) != 1 {
		t.Fail()
	} else if result[0].Id() != response.Id() {
		t.Fail()
	}

	mapper = commons.NewEventRedirection(catcher, processor, func(e commons.Event) bool { return false })

	if _, err := mapper.Process(event); err == nil {
		t.Fail()
	}
}
