package maths_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/maths"
)

// TestNewSerie validates the initialization of a new series with different types and sizes.
func TestNewSerie(t *testing.T) {
	// Test with float64
	s64 := maths.NewSerie(10, 0.0)
	if s64.Size() != 10 {
		t.Errorf("Expected size 10, got %d", s64.Size())
	}

	// Test with float32
	s32 := maths.NewSerie(5, float32(1.1))
	if s32.Size() != 5 {
		t.Errorf("Expected size 5, got %d", s32.Size())
	}

	// Test with zero size
	szero := maths.NewSerie(0, 0.0)
	if szero.Size() != 0 {
		t.Errorf("Expected size 0, got %d", szero.Size())
	}
}

// TestSerie_GetSet tests basic data access and the "auto-grow" behavior of the Set method.
func TestSerie_GetSet(t *testing.T) {
	defaultValue := 10.5
	s := maths.NewSerie(3, defaultValue)

	// Test initial state (default values)
	val, ok := s.Get(1)
	if !ok || val != defaultValue {
		t.Errorf("Expected default value %f, got %f", defaultValue, val)
	}

	// Test Set within bounds
	newValue := 99.9
	err := s.Set(1, newValue)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	val, ok = s.Get(1)
	if !ok || val != newValue {
		t.Errorf("Expected updated value %f, got %f", newValue, val)
	}

	// Test Set out of bounds (should grow the series)
	growIndex := 5
	growValue := 123.0
	err = s.Set(growIndex, growValue)
	if err != nil {
		t.Errorf("Set failed on grow: %v", err)
	}
	if s.Size() != growIndex+1 {
		t.Errorf("Expected size %d after grow, got %d", growIndex+1, s.Size())
	}

	// Test invalid negative index
	err = s.Set(-1, 1.0)
	if err == nil {
		t.Error("Expected error for negative index in Set, got nil")
	}

	// Test Get out of bounds
	_, ok = s.Get(100)
	if ok {
		t.Error("Expected Get to return ok=false for out of bounds index")
	}
}

// TestSerie_Append verifies that appending values correctly increases the size.
func TestSerie_Append(t *testing.T) {
	s := maths.NewSerie(2, 0.0)
	s.Append(1.0)
	s.Append(2.0)

	if s.Size() != 4 {
		t.Errorf("Expected size 4 after appends, got %d", s.Size())
	}

	val, _ := s.Get(3)
	if val != 2.0 {
		t.Errorf("Expected last element to be 2.0, got %f", val)
	}
}

// TestSerie_Values checks the materialization of the series into a standard slice.
func TestSerie_Values(t *testing.T) {
	s := maths.NewSerie(3, 0.0)
	s.Set(0, 1.1)
	s.Set(2, 3.3)

	values := s.Values()
	expected := []float64{1.1, 0.0, 3.3}

	if len(values) != len(expected) {
		t.Fatalf("Expected slice length %d, got %d", len(expected), len(values))
	}

	for i, v := range values {
		if v != expected[i] {
			t.Errorf("At index %d: expected %f, got %f", i, expected[i], v)
		}
	}
}

// TestSerie_Equals validates the comparison logic, including floating point precision.
func TestSerie_Equals(t *testing.T) {
	// Test cases for float64 (using LONG_EPSILON)
	s1 := maths.NewSerie(2, 0.0)
	s1.Set(0, 1.0)

	s2 := maths.NewSerie(2, 0.0)
	s2.Set(0, 1.000000000001) // Difference < 1e-9 (LONG_EPSILON)

	if !s1.Equals(s2) {
		t.Error("Series should be equal within float64 epsilon")
	}

	s3 := maths.NewSerie(2, 0.0)
	s3.Set(0, 1.1)
	if s1.Equals(s3) {
		t.Error("Series should not be equal")
	}

	// Test different sizes
	s4 := maths.NewSerie(3, 0.0)
	if s1.Equals(s4) {
		t.Error("Series with different sizes should not be equal")
	}
}

// TestSerie_Cut tests the sub-series creation and index validation.
func TestSerie_Cut(t *testing.T) {
	s := maths.NewSerie(10, 0.0)
	for i := 0; i < 10; i++ {
		s.Set(i, float64(i))
	}

	// Valid cut [2, 4] -> size 3
	sub, err := s.Cut(2, 4)
	if err != nil {
		t.Fatalf("Cut failed: %v", err)
	}

	if sub.Size() != 3 {
		t.Errorf("Expected sub-series size 3, got %d", sub.Size())
	}

	val, _ := sub.Get(0) // Should be original index 2
	if val != 2.0 {
		t.Errorf("Expected value 2.0 at sub-index 0, got %f", val)
	}

	// Test invalid indices
	_, err = s.Cut(-1, 5)
	if err == nil {
		t.Error("Expected error for negative start index")
	}

	_, err = s.Cut(5, 2)
	if err == nil {
		t.Error("Expected error for start > to")
	}

	_, err = s.Cut(0, 20)
	if err == nil {
		t.Error("Expected error for out of bounds end index")
	}
}

// TestSerie_SparseMemoryEfficiency is a logical test ensuring that setting
// values back to default doesn't break the series behavior.
func TestSerie_SparseMemoryEfficiency(t *testing.T) {
	defaultValue := 5.0
	s := maths.NewSerie(100, defaultValue)

	s.Set(50, 10.0)
	s.Set(50, defaultValue) // Reset to default

	val, ok := s.Get(50)
	if !ok || val != defaultValue {
		t.Errorf("Expected value to return to default %f, got %f", defaultValue, val)
	}

	// Ensure Size is still correct
	if s.Size() != 100 {
		t.Errorf("Size should remain 100, got %d", s.Size())
	}
}
