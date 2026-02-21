package commons

import (
	"container/heap"
	"errors"
	"iter"
	"slices"
	"sync"
	"time"
)

// ============================================================================
// PART 1: Priority Queue Implementation
// ============================================================================
// Architecture Note:
// Go's container/heap requires us to implement the heap.Interface.
// This section defines the low-level structure that manages the "Time Horizon".
// It acts as the engine allowing O(1) access to the next event and O(log N) insertions.
// This structure is NOT thread-safe by itself; it relies on the Graph's mutex.
// ============================================================================

// activeNode represents a graph element scheduled for processing.
// It is a lightweight wrapper used only within the priority queue.
type activeNode struct {
	// id of the target graph element.
	id string

	// wakeUp is the timestamp when this node needs processing.
	// This is the priority key (lower time = higher priority).
	wakeUp time.Time

	// index is the position of this node in the heap slice.
	// It is maintained by the heap interface methods.
	// We need this for the heap.Fix() method, which allows us to update
	// a node's priority in O(log N) if a new event arrives earlier than expected.
	index int
}

// priorityQueue implements heap.Interface and holds activeNodes.
type priorityQueue []*activeNode

func (pq priorityQueue) Len() int { return len(pq) }

// Less dictates the sorting logic of the heap.
// We want a Min-Heap based on time (earliest event first).
func (pq priorityQueue) Less(i, j int) bool {
	// Primary Sort Key: Time.
	if !pq[i].wakeUp.Equal(pq[j].wakeUp) {
		return pq[i].wakeUp.Before(pq[j].wakeUp)
	}
	// Secondary Sort Key: ID (Deterministic tie-breaking).
	// This ensures that two simulations with the exact same inputs
	// always yield the exact same execution order.
	return pq[i].id < pq[j].id
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	// Critical: Maintain the index map for O(1) lookups during updates.
	pq[i].index = i
	pq[j].index = j
}

// Push adds an item. logic is for the slice backend, heap.Push handles the sifting.
func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*activeNode)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes the last item. heap.Pop handles the sifting before calling this.
func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // Avoid memory leak
	item.index = -1 // Safety
	*pq = old[0 : n-1]
	return item
}

// ============================================================================
// PART 2: Graph Data Structures
// ============================================================================

// graphElement represents a node in the graph state space.
// It wraps an EventMapper (the logic/object) and manages its interactions in time.
type graphElement struct {
	// id is the unique identifier of the element.
	id string

	// mapper is the logic unit that processes events.
	mapper EventMapper

	// successors represents the outgoing links from this node.
	// Key: Destination ID, Value: Latency (transport time).
	successors map[string]time.Duration

	// events is the mailbox of the node.
	// It stores future or current events indexed by their scheduled arrival time.
	// Implementation Note: We use a map here because the number of distinct timestamps
	// per node is usually small. If a single node queues thousands of distinct
	// timestamps, this could be replaced by a local per-node Min-Heap.
	events map[time.Time][]Event
}

func (g *graphElement) Id() string {
	return g.id
}

// Graph represents the dynamic system of local events mappers.
// It uses a Discrete Event Simulation (DES) kernel powered by a Min-Heap.
type Graph interface {
	// Set creates or updates a directed link between a source and a destination.
	Set(source, destination EventMapper, latency time.Duration) error

	// Neighbors returns an iterator over the direct successors of a source node.
	Neighbors(source EventMapper) iter.Seq2[EventMapper, time.Duration]

	// Step advances the simulation by a fixed duration dt.
	// It processes all events scheduled up to (current_time + dt).
	//
	// Complexity:
	// O(E * log A), where:
	// - E is the number of events processed in this step.
	// - A is the number of ACTIVE nodes (nodes present in the heap).
	// This is vastly superior to O(N) for large, sparse graphs.
	Step(dt time.Duration) error

	// Emit schedules an external event to be processed by a specific mapper at a given time.
	Emit(target EventMapper, event Event, at time.Time) error
}

// localGraph is the optimized implementation of the Graph interface.
// Architecture Note (Thread-Safety):
// This structure is designed to be concurrently accessed. Multiple goroutines
// can inject events (Emit) or read topology (Neighbors) safely while Step is running.
type localGraph struct {
	// mu is the global lock securing the state of the graph.
	// We use an RWMutex to allow concurrent reads (like Neighbors snapshotting)
	// when no mutation is occurring.
	mu sync.RWMutex

	// sharedTime is the current global clock of the simulation.
	sharedTime time.Time

	// elements holds all the nodes in the graph, indexed by their ID.
	elements map[string]*graphElement

	// queue is the central engine of the simulation.
	// It contains only the nodes that have pending work.
	// Invariant: If a node is in this queue, it has at least one event in its mailbox.
	queue priorityQueue

	// lookup provides O(1) access to items inside the priority queue.
	// This allows us to update a node's wake-up time (DecreaseKey operation) efficiently
	// without scanning the entire heap.
	lookup map[string]*activeNode
}

// NewGraph creates a new empty graph initialized at a specific start time.
func NewGraph(startTime time.Time) Graph {
	// We initialize the heap with capacity 0, it will grow dynamically.
	pq := make(priorityQueue, 0)
	heap.Init(&pq)

	return &localGraph{
		sharedTime: startTime,
		elements:   make(map[string]*graphElement),
		queue:      pq,
		lookup:     make(map[string]*activeNode),
	}
}

// scheduleNode manages the presence of a node in the priority queue.
// Architecture Note:
// This method implements the "Decrease-Key" operation common in Dijkstra-like algorithms.
// If a node receives an event earlier than its currently scheduled wake-up time,
// we update its position in the heap to ensure it is processed sooner.
// WARNING: The caller MUST hold the g.mu lock (Write Lock) before invoking this function.
func (g *localGraph) scheduleNode(nodeId string, at time.Time) {
	item, exists := g.lookup[nodeId]

	if !exists {
		// Case 1: Node is not currently scheduled.
		// Create a new entry and push it to the heap.
		item = &activeNode{
			id:     nodeId,
			wakeUp: at,
		}
		heap.Push(&g.queue, item)
		g.lookup[nodeId] = item
	} else {
		// Case 2: Node is already scheduled.
		// We only update if the new event arrives BEFORE the currently known wake-up time.
		// (We maintain the "Earliest Deadline First" invariant).
		if at.Before(item.wakeUp) {
			item.wakeUp = at
			// heap.Fix re-establishes the heap invariant in O(log N)
			heap.Fix(&g.queue, item.index)
		}
	}
}

// Set creates or updates a directed link between a source and a destination.
func (g *localGraph) Set(source, destination EventMapper, latency time.Duration) error {
	if latency <= 0 {
		return errors.New("latency must be positive to respect causality")
	}

	// Acquire write lock to safely mutate the topology
	g.mu.Lock()
	defer g.mu.Unlock()

	// Lazy initialization check (safety)
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

	// Create the topological link
	srcElem.successors[destination.Id()] = latency
	return nil
}

// Neighbors returns the current neighbors of the source.
// Architecture Note (Deadlock Prevention):
// We cannot hold a Read Lock during the 'yield' call. If the consumer of this iterator
// attempts to call Emit() or Set() within the loop, it would request a Write Lock while
// the Read Lock is still held, resulting in a classic deadlock.
// Therefore, we take a rapid "snapshot" of the state, release the lock, and then yield.
func (g *localGraph) Neighbors(source EventMapper) iter.Seq2[EventMapper, time.Duration] {
	// 1. Acquire Read Lock for a safe snapshot
	g.mu.RLock()

	type neighbor struct {
		mapper  EventMapper
		latency time.Duration
	}
	var snapshot []neighbor

	if link := g.elements[source.Id()]; link != nil {
		for k, t := range link.successors {
			// We assume destination exists if it is in successors map
			if dest, ok := g.elements[k]; ok {
				snapshot = append(snapshot, neighbor{mapper: dest.mapper, latency: t})
			}
		}
	}

	// 2. Release Lock immediately
	g.mu.RUnlock()

	// 3. Iterate over the isolated snapshot safely
	return func(yield func(EventMapper, time.Duration) bool) {
		for _, n := range snapshot {
			if !yield(n.mapper, n.latency) {
				return
			}
		}
	}
}

// Emit schedules an external event.
// It acts as the "Interrupt Controller", injecting stimuli into the system.
func (g *localGraph) Emit(target EventMapper, event Event, at time.Time) error {
	// Acquire write lock as we are mutating node mailboxes and the heap
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.elements == nil {
		return errors.New("graph is empty")
	}

	targetId := target.Id()
	element, exists := g.elements[targetId]
	if !exists {
		return errors.New("target mapper not found in graph")
	}

	// 1. Data Persistence: Store the event in the node's local mailbox.
	element.events[at] = append(element.events[at], event)

	// 2. Scheduler Notification: Inform the heap that this node needs CPU time at 'at'.
	// Lock is safely held here.
	g.scheduleNode(targetId, at)

	return nil
}

// Step advances the simulation by a duration dt.
// Unlike a fixed time-step loop, this method "jumps" between active nodes
// using the priority queue, ensuring we only burn cycles on nodes with actual work.
//
// Architecture Note (Concurrency):
// The entire Step is executed under a Write Lock to preserve strict causal ordering.
// This means EventMappers (OnEvents) are executed synchronously within the lock.
// WARNING: An EventMapper MUST NOT call graph.Emit() or graph.Set() internally,
// as Go mutexes are non-reentrant and this would cause an immediate deadlock.
func (g *localGraph) Step(dt time.Duration) error {
	if dt <= 0 {
		return errors.New("step duration must be positive")
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	if g == nil || len(g.elements) == 0 {
		return nil
	}

	// The horizon defines the limit of our "Look-Ahead".
	// We will process all events occurring before or at this time.
	processingHorizon := g.sharedTime.Add(dt)

	// Processing Loop:
	// We drain the priority queue as long as the next scheduled event
	// is within our time window.
	for g.queue.Len() > 0 {

		// 1. Peek at the next candidate (Optimistic check)
		// We look at queue[0] without removing it to check the timestamp.
		nextNodesItem := g.queue[0]

		if nextNodesItem.wakeUp.After(processingHorizon) {
			// The next event is in the future (beyond this Step).
			// We can safely stop processing for this tick.
			break
		}

		// 2. Pop the node from the scheduler (Commit)
		// We remove the node from the heap because we are about to process it.
		// If it has remaining work afterwards, we will re-insert it.
		item := heap.Pop(&g.queue).(*activeNode)
		nodeId := item.id

		// Remove from lookup map since it's no longer in the heap queue.
		delete(g.lookup, nodeId)

		element := g.elements[nodeId]

		// 3. Identify Relevant Events
		// A node might have events at T=1, T=5, T=100. If Horizon=10,
		// we only want to process T=1 and T=5.
		var timeSlots []time.Time
		for t := range element.events {
			if !t.After(processingHorizon) {
				timeSlots = append(timeSlots, t)
			}
		}

		// Optimization: If for some reason the node woke up but has no valid events
		// (e.g., event was cancelled or already processed), verify reschedule and continue.
		if len(timeSlots) == 0 {
			goto RescheduleCheck
		}

		// 4. Sort Events Chronologically
		// Events within a node must be processed in strict causal order.
		slices.SortFunc(timeSlots, func(a, b time.Time) int {
			return a.Compare(b)
		})

		// 5. Execution Loop (The "Mapper" Logic)
		for _, t := range timeSlots {
			inputEvents := element.events[t]

			// Execute the business logic
			// (Assuming OnEvents is a pure function or state mutator that does NOT call the Graph API)
			outputEvents := element.mapper.OnEvents(inputEvents)

			// Clean up processed state
			delete(element.events, t)

			// 6. Propagation (The "Network" Logic)
			if len(outputEvents) > 0 {
				for successorId, latency := range element.successors {
					successor, exists := g.elements[successorId]
					if !exists {
						continue
					}

					// Calculate when the event arrives at the destination
					arrivalTime := t.Add(latency)

					// Push data to destination mailbox
					successor.events[arrivalTime] = append(successor.events[arrivalTime], outputEvents...)

					// Schedule the destination node in the Heap
					// (If it's already scheduled, this might update its priority)
					g.scheduleNode(successorId, arrivalTime)
				}
			}
		}

	RescheduleCheck:
		// 7. Re-scheduling (Context Switch)
		// The node may still have events scheduled AFTER the current horizon.
		// We must find the earliest future event and put the node back in the heap.
		// Example: Processed T=5, but mailbox has T=15. We re-queue for T=15.
		var nextWakeUp time.Time
		foundFutureEvent := false

		for t := range element.events {
			// Find the minimum time 't' remaining in the map
			if !foundFutureEvent || t.Before(nextWakeUp) {
				nextWakeUp = t
				foundFutureEvent = true
			}
		}

		if foundFutureEvent {
			g.scheduleNode(nodeId, nextWakeUp)
		}
	}

	// 8. Advance Global Clock
	g.sharedTime = processingHorizon
	return nil
}
