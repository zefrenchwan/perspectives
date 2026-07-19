package values_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

func TestPrimitiveBuilderIndependence(t *testing.T) {
	base := periods.NewTimeFunction(values.PRIMITIVE_TYPE_INT, values.EqualPrimitiveValue)

	if valuesMapping, err := values.NewPrimitiveMappingBuilder(base).Build(); err != nil {
		t.Error("failed to build mapping", err)
	} else if !valuesMapping.IsEmpty() {
		t.Error("values mapping should not be empty (independent from base)")
	} else {
		base.Add(values.NewInt(50), periods.NewFullPeriod())
		if !valuesMapping.IsEmpty() {
			t.Error("values mapping should be empty (independent from base)")
		}
	}
}

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
	} else if mapping.ValuesType() != values.REFERENCE_TYPE {
		t.Error("unexpected values type")
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

	if err := builder.Add(50, periods.NewFullPeriod()); err != nil {
		t.Error("failed to add value", err)
	} else if err := builder.Add(10, periods.NewFullPeriod()); err != nil {
		t.Error("failed to add value", err)
	} else if err := builder.Add("not the right type", periods.NewFullPeriod()); err == nil {
		t.Error("expected error : expected int, got string")
	}

	if mapping, err := builder.Build(); err != nil {
		t.Error("failed to build mapping", err)
	} else if mapping == nil {
		t.Error("mapping is nil")
	} else if mapping.ValuesType() != values.PRIMITIVE_TYPE_INT {
		t.Error("unexpected values type")
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

func TestPrimitiveBuilderLoad(t *testing.T) {
	base := periods.NewTimeFunction(values.PRIMITIVE_TYPE_INT, values.EqualPrimitiveValue)
	base.Add(values.NewInt(50), periods.NewFullPeriod())
	valuesBuilder := values.NewPrimitiveMappingBuilder(base)
	valuesMapping, _ := valuesBuilder.Build()

	other := periods.NewTimeFunction(values.PRIMITIVE_TYPE_INT, values.EqualPrimitiveValue)
	other.Add(values.NewInt(10), periods.NewFullPeriod())
	otherBuilder := values.NewPrimitiveMappingBuilder(other)

	// Other contains 10; we add 50, we expect 50
	if err := otherBuilder.Load(valuesMapping); err != nil {
		t.Error("failed to load mapping", err)
	}

	otherMapping, _ := otherBuilder.Build()
	for period, value := range otherMapping.Range() {
		if value.Content() != 50 {
			t.Error("unexpected value", value)
		} else if !periods.NewFullPeriod().Equals(period) {
			t.Error("unexpected period", period)
		}
	}
}

func TestPrimitiveBuilderLoadInconsistentType(t *testing.T) {
	base := periods.NewTimeFunction(values.PRIMITIVE_TYPE_INT, values.EqualPrimitiveValue)
	base.Add(values.NewInt(50), periods.NewFullPeriod())
	valuesBuilder := values.NewPrimitiveMappingBuilder(base)
	valuesMapping, _ := valuesBuilder.Build()

	other := periods.NewTimeFunction(values.PRIMITIVE_TYPE_STRING, values.EqualPrimitiveValue)
	other.Add(values.NewString("10"), periods.NewFullPeriod())
	otherBuilder := values.NewPrimitiveMappingBuilder(other)

	// Other contains 10; we add 50, we expect 50
	if err := otherBuilder.Load(valuesMapping); err == nil {
		t.Error("mixed int and string should fail ")
	}
}
