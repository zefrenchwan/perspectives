package commons_test

import (
	"slices"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

type DummyProcessor struct {
	id      string
	counter int
}

func (dp *DummyProcessor) Id() string {
	return dp.id
}

func (dp *DummyProcessor) OnEvent(event []commons.Event) []commons.Event {
	dp.counter += len(event)
	return event
}

func TestLocalContainer(t *testing.T) {
	a := &DummyProcessor{id: "a"}
	b := &DummyProcessor{id: "b"}

	container := commons.NewProcessorsGroup()
	container.Register(a, time.Now())
	container.Register(b, time.Now())

	elements := slices.Collect(container.Processors())
	if len(elements) != 2 {
		t.Fail()
	}

	container.Append(a, commons.NewMessage("b", []string{"a"}, "Hi"))
	if result := container.Launch(a); len(result) != 1 {
		t.Fail()
	} else if a.counter != 1 {
		t.Fail()
	}

	container.Remove(a)
	container.Append(a, commons.NewMessage("b", []string{"a"}, "Hi"))
	if result := container.Launch(a); result != nil {
		t.Fail()
	} else if a.counter != 1 {
		t.Fail()
	}

	elements = slices.Collect(container.Processors())
	if len(elements) != 1 {
		t.Fail()
	}

	container.Remove(b)
	elements = slices.Collect(container.Processors())
	if len(elements) != 0 {
		t.Fail()
	}
}
