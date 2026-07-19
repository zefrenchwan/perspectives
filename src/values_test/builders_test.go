package values_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

func TestReferenceBuilder(t *testing.T) {
	base := periods.NewTimeRelation(values.REFERENCE_TYPE, values.EqualReferences)
	builder := values.NewReferenceMappingBuilder(base)
	full := periods.NewFullPeriod()

	value1 := "one id"
	value2 := "another id"

	builder.Add(value1, periods.NewFullPeriod())
	builder.Add(value2, periods.NewFullPeriod())
	if mapping, err := builder.Build(); err != nil {
		t.Error("failed to build mapping", err)
	} else if mapping == nil {
		t.Error("mapping is nil")
	} else {
		found := make(map[string]bool)
		for period, value := range mapping.Range() {
			if !full.Equals(period) {
				t.Error("unexpected period", period)
			}

			found[value.Content().(string)] = true
		}

		if len(found) != 2 {
			t.Error("missing value")
		} else if !found[value1] || !found[value2] {
			t.Error("missing value")
		}
	}
}

func TestPrimitiveBuilder(t *testing.T) {
	base := periods.NewTimeFunction(values.PRIMITIVE_TYPE_INT, values.EqualPrimitiveValue)
	builder := values.NewPrimitiveMappingBuilder(base)
	full := periods.NewFullPeriod()

	builder.Add(50, periods.NewFullPeriod())
	builder.Add(10, periods.NewFullPeriod())

	if mapping, err := builder.Build(); err != nil {
		t.Error("failed to build mapping", err)
	} else if mapping == nil {
		t.Error("mapping is nil")
	} else {
		counter := 0
		for period, value := range mapping.Range() {
			counter++
			if value.Content() != 10 {
				t.Error("unexpected value", value)
			} else if !full.Equals(period) {
				t.Error("unexpected period", period)
			}
		}

		if counter != 1 {
			t.Error("unexpected number of values", counter)
		}
	}
}
