package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

type DummyEvent struct{}

func TestEventFunctionalMapper(t *testing.T) {
	events := []commons.Event{DummyEvent{}}
	mapper := func(e []commons.Event) []commons.Event {
		return []commons.Event{
			DummyEvent{},
			DummyEvent{},
		}
	}

	decorated := commons.NewEventMapper(mapper)
	if result := decorated.OnEvents(events); len(result) != 2 {
		t.Fail()
	}
}

func TestEventIdMapper(t *testing.T) {
	events := []commons.Event{DummyEvent{}}
	if result := commons.NewEventIdMapper().OnEvents(events); len(result) != 1 {
		t.Fail()
	}
}
