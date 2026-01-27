package commons

import (
	"time"

	"github.com/google/uuid"
)

// Event defines any impulse that would change a state
type Event interface {
}

// EventMapper processes events.
// Special cases: event tick to act, no event to emit (just change state)
type EventMapper interface {
	// OnEvents forces the processor to act, it may emit events
	OnEvents(events []Event) []Event
}

// functionaEventMapper implements EventMapper as a basic function
type functionaEventMapper func(events []Event) []Event

// OnEvents is indeed a function call
func (f functionaEventMapper) OnEvents(events []Event) []Event {
	if f == nil {
		return nil
	}

	return f(events)
}

// NewEventMapper return a new EventMapper decoring a function
func NewEventMapper(mapper func([]Event) []Event) EventMapper {
	return functionaEventMapper(mapper)
}

// NewEventIdMapper returns the id mapper
func NewEventIdMapper() EventMapper {
	return functionaEventMapper(func(events []Event) []Event {
		return events
	})
}

// EventTick is a special event to force action even for no other event
type EventTick struct {
	// dt is the local time to increase
	dt time.Duration
}

// Apply gets the new moment when we apply dt
func (e EventTick) Apply(moment time.Time) time.Time {
	return moment.Add(e.dt)
}

// NewEventTick builds a new tick
func NewEventTick(dt time.Duration) EventTick {
	return EventTick{dt: dt}
}

// Message defines a generic message, directed from a source to targets
type Message struct {
	Id           string            `json:"id"`            // Id of the message (for later log)
	CreationTime time.Time         `json:"creation_time"` // CreationTime defines the moment the message was created
	Source       string            `json:"source"`        // Source is the agent that sent the message
	Target       []string          `json:"targets"`       // Targets are the agents to send message to
	Payload      []byte            `json:"payload"`       // Payload is the content to send
	Metadata     map[string]string `json:"metadata"`      // Metadata defines any metadata for the message
}

// NewMessage returns a new message
func NewMessage(source string, targets []string, payload string) Message {
	return Message{
		Id:           uuid.NewString(),
		CreationTime: time.Now(),
		Source:       source,
		Target:       targets,
		Payload:      []byte(payload),
		Metadata:     make(map[string]string),
	}
}

// DecoratedEvent decorates events coming from a single source
type DecoratedEvent struct {
	// Source is the id of the common emitter
	Source string
	// OriginalEvents are the events coming from that source
	OriginalEvents []Event
}
