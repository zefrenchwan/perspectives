package periods

// TimeBounded defines the activity of an element : outside, it is inactive / dead.
type TimeBounded interface {
	// Activity returns the period during which the element is active.
	Activity() Period
}
