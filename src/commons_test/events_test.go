package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestEventFunctionalMapper(t *testing.T) {
	events := []commons.Event{commons.NewEventTick(time.Second)}
	mapper := func(events []commons.Event) []commons.Event {
		return []commons.Event{
			commons.NewMessage("source", []string{"dest"}, "test"),
			commons.NewEventTick(time.Second),
		}
	}

	decorated := commons.NewEventMapper(mapper)
	if result := decorated.OnEvents(events); len(result) != 2 {
		t.Fail()
	}
}

func TestEventIdMapper(t *testing.T) {
	events := []commons.Event{commons.NewEventTick(time.Second)}
	if result := commons.NewEventIdMapper().OnEvents(events); len(result) != 1 {
		t.Fail()
	}
}

func TestEvents(t *testing.T) {
	moment := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	result := commons.NewEventTick(time.Hour).Apply(moment)
	if !moment.Add(time.Hour).Equal(result) {
		t.Fail()
	}
}
