package commons

import (
	"errors"
	"slices"
	"time"
)

// graphElement represents a node in the graph state space.
// It wraps an EventMapper (the logic/object) and manages its interactions in time.
type graphElement struct {
	// id is the unique identifier of the element (same as the underlying mapper).
	id string
	// mapper is the logic unit that processes events.
	mapper EventMapper
	// successors represents the outgoing links from this node.
	// Key: Destination ID
	// Value: Latency (transport time) to reach that destination.
	successors map[string]time.Duration
	// events is the mailbox of the node.
	// It stores future or current events indexed by their scheduled arrival time.
	events map[time.Time][]Event
}

// Id returns the unique identifier of the graph element.
func (g *graphElement) Id() string {
	return g.id
}

// Graph represents the dynamic system of local events mappers.
// It manages time, topology (connections), and event propagation.
type Graph struct {
	// sharedTime is the current global clock of the simulation.
	// It advances via the Step method.
	sharedTime time.Time
	// elements holds all the nodes in the graph, indexed by their ID.
	elements map[string]*graphElement
}

// NewGraph creates a new empty graph initialized at a specific start time.
func NewGraph(startTime time.Time) *Graph {
	return &Graph{
		sharedTime: startTime,
		elements:   make(map[string]*graphElement),
	}
}

// Set creates or updates a directed link between a source and a destination.
// If the nodes do not exist in the graph, they are automatically added.
//
// source: The node emitting events.
// destination: The node receiving events.
// latency: The time it takes for an event to travel from source to destination (must be > 0).
func (g *Graph) Set(source, destination EventMapper, latency time.Duration) error {
	if latency <= 0 {
		return errors.New("latency must be positive to respect causality")
	}

	// 1. Ensure the graph map is initialized
	if g.elements == nil {
		g.elements = make(map[string]*graphElement)
	}

	sourceId := source.Id()
	destId := destination.Id()

	// 2. Ensure destination exists in the graph
	if _, ok := g.elements[destId]; !ok {
		g.elements[destId] = &graphElement{
			id:         destId,
			mapper:     destination,
			successors: make(map[string]time.Duration),
			events:     make(map[time.Time][]Event),
		}
	}

	// 3. Ensure source exists in the graph
	if _, ok := g.elements[sourceId]; !ok {
		g.elements[sourceId] = &graphElement{
			id:         sourceId,
			mapper:     source,
			successors: make(map[string]time.Duration),
			events:     make(map[time.Time][]Event),
		}
	}

	// 4. Create or update the link (the edge)
	// We set the latency on the source's outgoing map.
	g.elements[sourceId].successors[destId] = latency

	return nil
}

// Step advances the simulation by a duration dt.
// It processes events chronologically to respect causality and propagation delays.
func (g *Graph) Step(dt time.Duration) error {
	if g == nil {
		return nil
	}
	if dt <= 0 {
		return errors.New("step duration must be positive")
	}

	// 1. Define the horizon: we process everything up to this new time (inclusive)
	processingHorizon := g.sharedTime.Add(dt)

	// true as soon as exists one event before the processing horizon
	var changes = true
	// Run graph walk as long as previous run made a change
	for changes {
		// ensure no change first !
		changes = false
		// for each node, find events to map and propagate
		for _, element := range g.elements {
			// 2. Snapshot: Identify all time slots to process in this step
			// We use a slice to store keys because we need to sort them
			var timeSlots []time.Time
			for t := range element.events {
				if !t.After(processingHorizon) {
					timeSlots = append(timeSlots, t)
				}
			}

			// 3. Sort times chronologically
			// Essential for stateful mappers (memory) and correct causal chain
			slices.SortFunc(timeSlots, func(a, b time.Time) int {
				return a.Compare(b)
			})

			// 4. Process loop
			for _, t := range timeSlots {
				events := element.events[t]

				// A. The Node reacts (transformation)
				outputEvents := element.mapper.OnEvents(events)

				// B. Propagation (displacement)
				// We use 't' (the event occurrence time) as the base, not the current simulation clock.
				if len(outputEvents) > 0 {
					// at least a change
					changes = true
					for successorId, transportLatency := range element.successors {
						if successor, exists := g.elements[successorId]; exists {
							// Arrival is relative to when the cause happened (t), plus the travel time.
							arrivalTime := t.Add(transportLatency)
							successor.events[arrivalTime] = append(successor.events[arrivalTime], outputEvents...)
						}
					}
				}

				// C. Cleanup
				// Now that we have processed and propagated, we can remove these events.
				delete(element.events, t)
			}
		}
	}

	// 5. Advance global time
	g.sharedTime = processingHorizon

	return nil
}

// Emit schedules an external event to be processed by a specific mapper at a given time.
func (g *Graph) Emit(target EventMapper, event Event, at time.Time) error {
	if g.elements == nil {
		return errors.New("graph is empty")
	}

	targetId := target.Id()
	element, exists := g.elements[targetId]
	if !exists {
		return errors.New("target mapper not found in graph")
	}

	// Schedule the event
	element.events[at] = append(element.events[at], event)
	return nil
}
