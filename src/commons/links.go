package commons

type Linkable any

// Link is a constant relation over instances of linkables.
// Link is also Linkable, so it may be used in links.
type Link[T Linkable] interface {
	// links are composable due to this: a link is linkable
	Linkable
	// Name returns the name of the link.
	// It is usually a verb or a noun.
	// For instance, knows, couple, etc.
	Name() string
	// Roles returns all the roles set for that link
	Roles() []string
	// Size returns the number of elements per role.
	// For no value, it returns 0
	Size(role string) int
	// Operands returns the roles and related values for that link
	Operands() map[string]T
	// Get returns, if any, the elements for that name.
	// First result is the values if any, second is true if value was found
	Get(role string) (T, bool)
}

// TemporalLink decorates a link over a period.
// It should implement both link and Temporal.
type TemporalLink[T Linkable] struct {
	// activity of the link
	period Period
	// value is link decoration
	value Link[T]
}

// NewTemporalLink decorates a link true for given duration
func NewTemporalLink[T Linkable](duration Period, value Link[T]) *TemporalLink[T] {
	result := new(TemporalLink[T])
	result.period = duration
	result.value = value
	return result
}

// ActivePeriod is the duration which the link is true
func (t *TemporalLink[T]) ActivePeriod() Period {
	if t == nil {
		return NewEmptyPeriod()
	}

	return t.period
}

// SetActivePeriod forces active period
func (t *TemporalLink[T]) SetActivePeriod(period Period) {
	if t != nil {
		t.period = period
	}
}

// Name returns the name of the link.
func (t *TemporalLink[T]) Name() string {
	var empty string
	if t == nil {
		return empty
	}

	return t.Name()
}

// Roles returns all the roles set for that link
func (t *TemporalLink[T]) Roles() []string {
	if t == nil {
		return nil
	}

	return t.Roles()
}

// Size returns the number of elements per role.
func (t *TemporalLink[T]) Size(role string) int {
	if t == nil {
		return 0
	}

	return t.Size(role)
}

// Operands returns the roles and related values for that link
func (t *TemporalLink[T]) Operands() map[string]T {
	if t == nil {
		return nil
	}

	return t.Operands()
}

// Get returns, if any, the element for that name.
func (t *TemporalLink[T]) Get(role string) (T, bool) {
	var empty T
	if t == nil {
		return empty, false
	}

	return t.Get(role)
}

// First returns the first value, if any, for that role
func (t *TemporalLink[T]) First(role string) (T, bool) {
	var empty T
	if t == nil {
		return empty, false
	}

	return t.First(role)
}
