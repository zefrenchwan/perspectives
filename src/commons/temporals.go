package commons

// Temporal is a component with a lifetime (active period)
type Temporal interface {
	// ActivePeriod returns the period the compound is active during
	ActivePeriod() Period
	// SetActivePeriod forces activity for compound
	SetActivePeriod(period Period)
}
