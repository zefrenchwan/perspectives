package entities_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/entities"
	"github.com/zefrenchwan/perspectives.git/periods"
)

// More a doc of empty behavior
func TestEmptyValues(t *testing.T) {
	base := periods.NewDynamicPartition[any]("int", func(a, b any) bool { return a == b })
	values := entities.NewDynamicValuesFromPartition(base)

	if !values.IsEmpty() {
		t.Errorf("Expected empty values from base")
	} else if values.DataType() != "int" {
		t.Errorf("Expected data type to be int")
	} else if !values.Validity().IsEmpty() {
		t.Errorf("Expected validity to be empty")
	}

	counter := 0
	for p := range values.Range() {
		_ = p
		counter++
	}

	if counter != 0 {
		t.Errorf("Expected counter to be 0")
	}
}
