package commons

// Element is an element within a field.
// It may be an actor or an event
type Element interface {
	// Identifiable to distinguish from each other
	Identifiable
	// OnElements applies when an element interacts with others
	OnElements([]Element) []Element
}

// Topology defines how elements synchronize via their position.
// Position is any identifiable (or composite of identifiables)
type Topology[Position Identifiable] interface {
	// Set an element at that position within the field
	Set(Element, Position)
	// OnEmission is what happens when elements interact and they emit events coming from there
	OnEmission(source Position, emitted []Element)
}
