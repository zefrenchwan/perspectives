package commons_test

import (
	"errors"
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
	// LastSource contains the last identifiable that emitted
	LastSource commons.Identifiable
}

func (o *DummyObserver) OnEventProcessing(source commons.Identifiable, in commons.Event, out []commons.Event, e error) {
	o.Incoming += 1
	if e != nil {
		o.Errors++
	}

	o.Outgoing += len(out)
	o.LastSource = source
}

type DummyInterceptor struct {
	Result commons.Event
}

func (i DummyInterceptor) OnRecipientProcessing(event commons.Event, recipient commons.EventProcessor) ([]commons.Event, error) {
	return []commons.Event{i.Result}, nil
}

func TestEventSource(t *testing.T) {
	structure := DummyStructure{id: "structure"}
	event := commons.NewEventLifetimeEnd(structure, time.Now())
	if event.Source() != structure {
		t.Fail()
	} else if !commons.IsEventComingFromStructure(event) {
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
	} else if counter.LastSource.Id() != observer.Id() {
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

func TestEventInterception(t *testing.T) {
	structure := DummyStructure{}
	source := commons.NewEventTick(structure)
	replace := commons.NewEventLifetimeEnd(structure, time.Now())

	mapper := commons.NewEventProcessor(func(e commons.Event) ([]commons.Event, error) { return []commons.Event{source}, nil })
	replacer := commons.NewEventInterceptor(func(e commons.Event, p commons.EventProcessor) ([]commons.Event, error) {
		return []commons.Event{replace}, nil
	})

	interceptor := commons.NewEventInterception(mapper, replacer)
	if result, err := interceptor.Process(nil); err != nil {
		t.Fail()
	} else if len(result) != 1 {
		t.Fail()
	} else if result[0].Id() != replace.Id() {
		t.Fail()
	}
}
