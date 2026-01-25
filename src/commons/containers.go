package commons

import (
	"iter"
	"sync"
	"time"
)

// ProcessorsGroup is a group of events processors that interact together.
// For any field use, there is, technically, this kind of code to manage events
type ProcessorsGroup interface {
	// Register the agent at a given time
	Register(element EventProcessor, moment time.Time)
	// Processors lists all processors in that group
	Processors() iter.Seq[EventProcessor]
	// Append adds events to an element (but do not fire it)
	Append(element EventProcessor, events ...Event)
	// Launch all events to process for a given element.
	// Result is the events that processor generated
	Launch(element EventProcessor) []Event
	// Remove that element of that group
	Remove(element EventProcessor) bool
}

// containerNode represents an event processor, its events to process, and its local time
type containerNode struct {
	// element is the processor per se
	element EventProcessor
	// localTime is the time for that processor
	localTime time.Time
	// events are the events that processor has to manage
	events []Event
}

// LocalContainer manages a group of processors
type LocalContainer struct {
	// mutex deals with concurrency operations
	mutex sync.RWMutex
	// elements are indexed by the id of the processor
	elements map[string]*containerNode
}

// Processors returns the event processors in that container
func (c *LocalContainer) Processors() iter.Seq[EventProcessor] {
	if c == nil {
		return func(yield func(EventProcessor) bool) {}
	}

	// copy them to manage the lock
	c.mutex.RLock()
	processors := make([]EventProcessor, 0, len(c.elements))
	for _, node := range c.elements {
		processors = append(processors, node.element)
	}
	c.mutex.RUnlock()

	// and then we let others deal with them at their own pace
	return func(yield func(EventProcessor) bool) {
		for _, p := range processors {
			if !yield(p) {
				return
			}
		}
	}
}

// Register adds a processor at that time
func (c *LocalContainer) Register(element EventProcessor, moment time.Time) {
	if c != nil && element != nil {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		previousNode := c.elements[element.Id()]
		if previousNode == nil {
			previousNode = new(containerNode)
			previousNode.element = element
		}

		previousNode.events = nil
		previousNode.localTime = moment
		c.elements[element.Id()] = previousNode
	}
}

// Append adds events for that element to process
func (c *LocalContainer) Append(element EventProcessor, events ...Event) {
	if c != nil && c.elements != nil && element != nil {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		previousNode := c.elements[element.Id()]
		if previousNode != nil {
			previousNode.events = append(previousNode.events, events...)
			c.elements[element.Id()] = previousNode
		}
	}
}

// Launch triggers that element (if in container) and processes all its registered events
func (c *LocalContainer) Launch(element EventProcessor) []Event {
	if c != nil && c.elements != nil && element != nil {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		previousNode := c.elements[element.Id()]
		if previousNode != nil {
			result := previousNode.element.OnEvent(previousNode.events)
			previousNode.events = nil
			c.elements[element.Id()] = previousNode
			return result
		}
	}
	return nil
}

// Remove element if it is in this container, return true if it did
func (c *LocalContainer) Remove(element EventProcessor) bool {
	if c != nil && c.elements != nil && element != nil {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		previousNode := c.elements[element.Id()]
		delete(c.elements, element.Id())
		return previousNode != nil
	}
	return false
}

// NewProcessorsGroup returns the most basic implementation, with local storage
func NewProcessorsGroup() ProcessorsGroup {
	result := new(LocalContainer)
	result.mutex = sync.RWMutex{}
	result.elements = make(map[string]*containerNode)
	return result
}
