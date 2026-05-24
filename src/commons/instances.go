package commons

import "time"

type TemporalValues interface {
	At(t time.Time) (any, bool)
	Range(yield func(p Period, v any) bool)
	IsEmpty() bool
}

type Instance interface {
	Lifetime() Period
	SetLifetime(Period)
	
	Description() map[string]string

	SetAttribute(name string, p Period, value any) error
	Attribute(name string) TemporalValues
}
