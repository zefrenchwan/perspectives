package commons

// StateValue is the definition of accepted types
type StateValue interface{ string | int | float64 | bool }

// StateReader declares the ability to deal with attributes that do not depend over time.
type StateReader[T StateValue] interface {
	// Attributes returns the names of the attributes set for that handler
	Attributes() []string
	// GetValue returns the value for that name (if any).
	// It returns the value (if any), true if found or false if not found
	GetValue(name string) (T, bool)
}

// StateHandler is reading and updating a state.
type StateHandler[T StateValue] interface {
	// Handler needs ability to read
	StateReader[T]
	// SetValue sets value for that attribute
	SetValue(name string, value T)
}

// TemporalStateReader reads a state that varies over time and register that state.
type TemporalStateReader[T StateValue] interface {
	// Attributes returns the names of the attributes that are set
	Attributes() []string
	// GetValues returns the values for that attribute (if any) and their period.
	// For instance, a => [now, +oo[ means that a is the value of that attribute during [now, +oo[
	GetValues(name string) map[T]Period
}

// TemporalStateHandler declares a state that varies over time.
// It means the ability to list attributes,
// for each attribute, be able to get values and related periods,
// and change those values during a given period.
type TemporalStateHandler[T StateValue] interface {
	// TemporalStateReader is necessary to change state over time
	TemporalStateReader[T]
	// SetValueDuringPeriod sets value for that attribute during a given period
	SetValueDuringPeriod(name string, value T, period Period)
}
