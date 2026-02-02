package commons

import (
	"errors"
	"iter"
	"slices"
	"time"
)

// graphElement represents a node in the graph state space.
// It wraps an EventMapper (the logic/object) and manages its interactions in time.
//
// In the Observer pattern, this element acts as the Subject that holds state (events).
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
//
// Architecture Note:
// This Graph implementation uses an "Activity-Driven" approach (Observer Pattern).
// Instead of scanning every node at every tick (Polling), it maintains a calendar
// of active nodes. A node is only processed if it has pending events.
type Graph interface {
	// Set creates or updates a directed link between a source and a destination.
	// If the nodes do not exist in the graph, they are automatically added.
	//
	// Parameters:
	// - source: The node emitting events.
	// - destination: The node receiving events.
	// - latency: The propagation delay (must be > 0).
	Set(source, destination EventMapper, latency time.Duration) error

	// Neighbors returns an iterator over the direct successors of a source node.
	// It yields the destination mapper and the latency of the link.
	Neighbors(source EventMapper) iter.Seq2[EventMapper, time.Duration]

	// Step advances the simulation by a fixed duration dt.
	// It processes all events scheduled up to (current_time + dt).
	//
	// Complexity:
	// The complexity is roughly O(A * E), where:
	// - A is the number of ACTIVE nodes in this time window (not total nodes).
	// - E is the average number of events processed per node.
	// This is significantly more efficient than O(N) for sparse activity graphs.
	Step(dt time.Duration) error

	// Emit schedules an external event to be processed by a specific mapper at a given time.
	// This is the entry point to inject stimuli into the system.
	Emit(target EventMapper, event Event, at time.Time) error
}

// localGraph is the default implementation of the Graph interface.
type localGraph struct {
	// sharedTime is the current global clock of the simulation.
	// It advances via the Step method.
	sharedTime time.Time

	// elements holds all the nodes in the graph, indexed by their ID.
	elements map[string]*graphElement

	// calendar acts as the "Observer" state.
	// It tracks the earliest known wake-up time for each node.
	// If a node ID is present in this map, it means the graph knows it has work to do.
	//
	// Logical Invariant: If a node has an event at time T, and T <= Horizon,
	// the node MUST be either in this calendar or currently in the processing queue.
	calendar map[string]time.Time
}

// NewGraph creates a new empty graph initialized at a specific start time.
func NewGraph(startTime time.Time) Graph {
	return &localGraph{
		sharedTime: startTime,
		elements:   make(map[string]*graphElement),
		calendar:   make(map[string]time.Time),
	}
}

// notifyActivity updates the calendar (Observer) for a node.
// It ensures we track the *earliest* wake-up time for the node.
// This is called whenever an event is scheduled (Emit or Propagation).
func (g *localGraph) notifyActivity(nodeId string, at time.Time) {
	existingTime, known := g.calendar[nodeId]
	// We only update if the node was unknown OR if the new time is earlier.
	// This ensures priority is given to the soonest event.
	if !known || at.Before(existingTime) {
		g.calendar[nodeId] = at
	}
}

// Set creates or updates a directed link between a source and a destination.
func (g *localGraph) Set(source, destination EventMapper, latency time.Duration) error {
	if latency <= 0 {
		return errors.New("latency must be positive to respect causality")
	}

	// Lazy initialization of the elements map
	if g.elements == nil {
		g.elements = make(map[string]*graphElement)
	}

	// Helper to ensure node existence in the graph structure
	ensureNode := func(em EventMapper) *graphElement {
		id := em.Id()
		if _, ok := g.elements[id]; !ok {
			g.elements[id] = &graphElement{
				id:         id,
				mapper:     em,
				successors: make(map[string]time.Duration),
				events:     make(map[time.Time][]Event),
			}
		}
		return g.elements[id]
	}

	srcElem := ensureNode(source)
	ensureNode(destination)

	// Create the edge
	srcElem.successors[destination.Id()] = latency
	return nil
}

// Neighbors returns the current neighbors of the source.
func (g *localGraph) Neighbors(source EventMapper) iter.Seq2[EventMapper, time.Duration] {
	return func(yield func(EventMapper, time.Duration) bool) {
		if link := g.elements[source.Id()]; link != nil {
			for k, t := range link.successors {
				// We assume destination exists if it is in successors map
				if dest, ok := g.elements[k]; ok {
					if !yield(dest.mapper, t) {
						return
					}
				}
			}
		}
	}
}

// Emit schedules an external event.
// It notifies the graph that the target node needs attention at 'at'.
func (g *localGraph) Emit(target EventMapper, event Event, at time.Time) error {
	if g.elements == nil {
		return errors.New("graph is empty")
	}

	targetId := target.Id()
	element, exists := g.elements[targetId]
	if !exists {
		return errors.New("target mapper not found in graph")
	}

	// 1. Store the event in the node's mailbox
	element.events[at] = append(element.events[at], event)

	// 2. Observer Notification:
	// We flag this node as "active" in the calendar.
	// The Step method will check this calendar to know whom to visit.
	g.notifyActivity(targetId, at)

	return nil
}

// Step advances the simulation by a duration dt using an activity queue.
// Instead of polling all nodes, it only processes nodes known to have events.
func (g *localGraph) Step(dt time.Duration) error {
	if g == nil || len(g.elements) == 0 {
		return nil
	}
	if dt <= 0 {
		return errors.New("step duration must be positive")
	}

	// Calculate the time horizon for this step.
	// We will process all events where t <= processingHorizon.
	processingHorizon := g.sharedTime.Add(dt)

	// 1. Initialization: Wake up nodes from the Calendar.
	// 'inQueue' tracks which nodes are currently in the processing line.
	// This prevents adding the same node multiple times redundantly,
	// but allows re-adding it if it receives new work after being processed.
	inQueue := make(map[string]bool)
	var workQueue []string

	for id, wakeUpTime := range g.calendar {
		if !wakeUpTime.After(processingHorizon) {
			inQueue[id] = true
			workQueue = append(workQueue, id)
			// Remove from calendar: the node is now "in flight".
			// If it has remaining future events, it must be re-scheduled later.
			delete(g.calendar, id)
		}
	}

	// Sort initial queue for deterministic behavior (optional but good for debugging/replay)
	slices.Sort(workQueue)

	// 2. Processing Loop (Breadth-First-Like propagation)
	// The queue grows dynamically: if Node A activates Node B (within horizon),
	// Node B is appended to the queue.
	for i := 0; i < len(workQueue); i++ {
		nodeId := workQueue[i]
		element := g.elements[nodeId]

		// Mark node as "not in queue anymore" (it is being processed).
		// This allows it to be re-added later in the loop if a feedback loop occurs.
		inQueue[nodeId] = false

		// A. Collect relevant events for this node (Timeline check)
		var timeSlots []time.Time
		for t := range element.events {
			if !t.After(processingHorizon) {
				timeSlots = append(timeSlots, t)
			}
		}

		// Skip if false alarm (optimization)
		if len(timeSlots) == 0 {
			// Even if no work was done, check if we need to reschedule future events
			// (See step E below)
			goto RescheduleCheck
		}

		// Sort events chronologically to respect causality inside the node
		slices.SortFunc(timeSlots, func(a, b time.Time) int {
			return a.Compare(b)
		})

		// B. Process events
		for _, t := range timeSlots {
			inputEvents := element.events[t]

			// The Mapper Logic (The "Brain" of the node)
			outputEvents := element.mapper.OnEvents(inputEvents)

			// Clean up processed events from memory
			delete(element.events, t)

			// C. Propagate results (The "Network" effect)
			if len(outputEvents) > 0 {
				for successorId, latency := range element.successors {
					successor, exists := g.elements[successorId]
					if !exists {
						continue
					}

					// Calculate arrival time relative to the event occurrence 't'
					arrivalTime := t.Add(latency)

					// Push events to destination
					successor.events[arrivalTime] = append(successor.events[arrivalTime], outputEvents...)

					// D. Observer Logic: Schedule the successor
					if !arrivalTime.After(processingHorizon) {
						// Case 1: Immediate Reaction.
						// The event arrives WITHIN the current step window.
						// We must ensure the successor processes it in this loop.
						if !inQueue[successorId] {
							inQueue[successorId] = true
							workQueue = append(workQueue, successorId)
						}
					} else {
						// Case 2: Future Reaction.
						// The event arrives AFTER the current step.
						// We store it in the calendar for the next Step() call.
						g.notifyActivity(successorId, arrivalTime)
					}
				}
			}
		}

	RescheduleCheck:
		// E. Reschedule remaining future events.
		// The node might have events scheduled after the horizon.
		// We must ensure they remain tracked in the calendar.
		// (Efficiency: we only need to find the earliest one).
		for t := range element.events {
			if t.After(processingHorizon) {
				g.notifyActivity(nodeId, t)
				// We can break here if we trust notifyActivity handles min comparison,
				// but since iteration order is random, we should strictly check all or rely on repeated calls.
				// Given map size is usually small, simple iteration is fine.
			}
		}
	}

	// 3. Advance global time
	g.sharedTime = processingHorizon
	return nil
}
