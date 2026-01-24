package commons

import (
	"time"

	"github.com/google/uuid"
)

// Event defines any impulse that would change a state
type Event interface {
}

// EventProcessor processes events.
// Special cases: event tick to act, no event to emit (just change state)
type EventProcessor interface {
	Identifiable
	// OnEvent forces the processor to act, it may emit events
	OnEvent(event []Event) []Event
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
	Payload      string            `json:"payload"`       // Payload is the content to send
	Metadata     map[string]string `json:"metadata"`      // Metadata defines any metadata for the message
}

// NewMessage returns a new message
func NewMessage(source string, target []string, payload string) Message {
	return Message{
		Id:           uuid.NewString(),
		CreationTime: time.Now(),
		Source:       source,
		Target:       target,
		Payload:      payload,
		Metadata:     make(map[string]string),
	}
}
