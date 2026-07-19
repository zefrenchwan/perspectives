package values_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

// TestEnsureInvariant proves the point for invariant values break
// if no mechanism is in place to ensure it.
func TestEnsureInvariant(t *testing.T) {
	base := periods.NewTimeFunction[values.PrimitiveValue]("string", values.EqualPrimitiveValue)
	primitiveString := values.NewString("a string")
	primitiveInt := values.NewInt(10)

	now := time.Now().Truncate(time.Second)
	beforePeriod := periods.NewPeriodUntil(now, true)
	afterPeriod := periods.NewPeriodSince(now, false)

	// empty, no invariant break
	if !values.EnsureValuesMappingInvariant(base) {
		t.Errorf("no value, invariant broken unexpected")
	}

	base.Add(primitiveString, beforePeriod)
	// one value, no invariant break
	if !values.EnsureValuesMappingInvariant(base) {
		t.Errorf("one value, invariant broken unexpected")
	}

	// adding int, expected invariant break
	base.Add(primitiveInt, afterPeriod)
	if values.EnsureValuesMappingInvariant(base) {
		t.Errorf("one value int, one string, invariant broken")
	}
}
