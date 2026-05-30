package commons

import (
	"errors"

	"github.com/google/uuid"
)

// InstanceBuilder manages building instances from database (serde).
// Typical usage is to create a new instance builder, set its properties, and then build the instance.
// Id especially is necessary when building instances from a storage system or file
type InstanceBuilder struct {
	// instance is the temporal instance being built
	instance *temporalInstance
	// globalError accumulates errors during instance building
	globalError error
}

// SetId sets the identifier for the temporal instance being built.
func (ib *InstanceBuilder) SetId(id string) *InstanceBuilder {
	ib.instance.id = id
	return ib
}

// SetValidityPeriod sets the validity period for the temporal instance being built.
func (ib *InstanceBuilder) SetValidityPeriod(period Period) *InstanceBuilder {
	if period.IsEmpty() {
		ib.globalError = errors.Join(ib.globalError, errors.New("period cannot be nil"))
		return ib
	}

	ib.instance.lifetime = period
	return ib
}

// SetValue sets an attribute with a given name, period, and value for the temporal instance being built.
func (ib *InstanceBuilder) SetValue(name string, period Period, value any) *InstanceBuilder {
	if err := ib.instance.SetAttribute(name, period, value); err != nil {
		ib.globalError = errors.Join(ib.globalError, err)
	}

	return ib
}

// Build constructs the temporal instance from the builder, returning it along with any accumulated errors.
func (ib *InstanceBuilder) Build() (Instance, error) {
	if ib.globalError != nil {
		return nil, ib.globalError
	}

	return ib.instance, ib.globalError
}

// NewInstanceBuilder creates a new instance builder with a default temporal instance initialized with a new UUID.
func NewInstanceBuilder() *InstanceBuilder {
	return &InstanceBuilder{
		instance: initTemporalInstance(uuid.NewString()),
	}
}
