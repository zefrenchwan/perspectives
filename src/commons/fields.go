package commons

import "time"

// Field is the tool to manage its elements
type Field interface {
	EventProcessor
	// Register the agent at a given time
	Register(element EventProcessor, moment time.Time)
}
