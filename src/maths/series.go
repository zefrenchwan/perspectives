package maths

import (
	"errors"
	"math"
)

// Serie is a generic interface representing a sequence of floating-point numbers.
// It supports basic operations for data manipulation, comparison, and slicing.
type Serie[F FloatNumber] interface {
	// Equals checks if this series is identical to another series.
	// Returns true if both have the same size and all elements are equal
	// based on the floating-point precision logic.
	Equals(other Serie[F]) bool
	// Size returns the total number of elements in the series.
	Size() int
	// Values returns the full sequence of values as a slice.
	Values() []F
	// Set assigns a value at the specified index.
	// If the index is beyond the current size, the series automatically grows.
	Set(index int, value F) error
	// Get retrieves the value at the specified index.
	// Returns the value and a boolean indicating if the index is valid.
	Get(index int) (F, bool)
	// Append adds a new value at the end of the series, increasing its size by one.
	Append(value F)
	// Cut creates a sub-series from the 'from' index to the 'to' index (inclusive).
	// Returns an error if indices are out of bounds.
	Cut(from, to int) (Serie[F], error)
	// Indicators returns the mean and standard deviation of the series.
	Indicators() (mean, stddev float64)
}

// localSerie is a memory-efficient implementation of the Serie interface.
// Implementation choice: It uses a map to store values that differ from a defaultValue.
// This "sparse" approach is highly efficient for large series containing many repeated values.
type localSerie[F FloatNumber] struct {
	// defaultValue is what we return if no other value was set
	defaultValue F
	// values contains the index based value
	values map[int]F
	// size is the current size of the serie
	size int
	// equality is the way to compare elements in it
	equality func(F, F) bool
}

// Equals compares two series for equality.
// Complexity: O(N) where N is the size of the series, as it iterates through every index.
// It ensures that even values not explicitly stored in the map (default values) are compared correctly.
func (l *localSerie[F]) Equals(other Serie[F]) bool {
	if l == nil && other == nil {
		return true
	} else if l == nil || other == nil {
		return false
	}

	size := l.Size()
	if other.Size() != size {
		return false
	}

	for i := 0; i < size; i++ {
		valA, _ := l.Get(i)
		valB, _ := other.Get(i)
		if !l.equality(valA, valB) {
			return false
		}
	}

	return true
}

// Size returns the current logical length of the series.
// Complexity: O(1).
func (l *localSerie[F]) Size() int {
	return l.size
}

// Values materializes the series into a slice of type F.
// Complexity: O(N) where N is the size of the series.
// Implementation choice: It pre-allocates the slice to avoid multiple reallocations during the loop.
func (l *localSerie[F]) Values() []F {
	if l == nil {
		return nil
	}

	result := make([]F, l.size)
	for i := 0; i < l.size; i++ {
		if value, found := l.values[i]; found {
			result[i] = value
		} else {
			result[i] = l.defaultValue
		}
	}

	return result
}

// Set updates the value at a specific index.
// Complexity: O(1) average for map insertion.
// Implementation choice: Only values different from the defaultValue are stored in the map to save memory.
// If the index is greater than the current size, the size is updated to index + 1.
func (l *localSerie[F]) Set(index int, value F) error {
	if index < 0 {
		return errors.New("invalid index")
	} else if index >= l.size {
		l.size = index + 1
	}

	if !l.equality(value, l.defaultValue) {
		l.values[index] = value
	} else {
		// Clean up the map if the value is changed back to the default
		delete(l.values, index)
	}

	return nil
}

// Get retrieves a value from the series.
// Complexity: O(1) average for map lookup.
// If the index exists but is not in the map, it returns the defaultValue.
func (l *localSerie[F]) Get(index int) (F, bool) {
	if index < 0 {
		return l.defaultValue, false
	} else if index < l.size {
		if value, found := l.values[index]; found {
			return value, true
		} else {
			return l.defaultValue, true
		}
	}

	return l.defaultValue, false
}

// Append adds a value to the end of the series.
// Complexity: O(1) average.
// This is a specialized case of Set(l.size, value).
func (l *localSerie[F]) Append(value F) {
	if l != nil {
		if !l.equality(value, l.defaultValue) {
			l.values[l.size] = value
		}

		l.size = l.size + 1
	}
}

// Cut returns a new Serie containing a subset of the original elements.
// Complexity: O(V) where V is the number of non-default values stored in the map.
// Implementation choice: It iterates over the internal map to copy only relevant stored values
// within the requested range, maintaining the sparse efficiency.
func (l *localSerie[F]) Cut(from, to int) (Serie[F], error) {
	if from < 0 || to < 0 || from >= l.size || to >= l.size || from > to {
		return nil, errors.New("invalid index")
	}

	// The new size of the cut series is determined by the input range
	result := newLocalSerie[F](to-from+1, l.defaultValue)
	for k, v := range l.values {
		if k >= from && k <= to {
			result.values[k-from] = v
		}
	}

	return result, nil
}

// Indicators calculates and returns the population mean and standard deviation of the series.
//
// It expects that none of the values are NaN. To ensure numerical stability,
// especially with large numbers or small variances, this implementation utilizes
// Welford's online algorithm.
//
// Furthermore, to maintain the memory and CPU efficiency of the sparse series representation,
// the algorithm is optimized to run in O(V) time, where V is the number of explicitly stored
// values in the map. It achieves this by running the standard Welford update on the map values
// first, followed by a single mathematical "batched" update (based on Welford's parallel formula)
// to account for all the remaining implicit default values at once.
func (l *localSerie[F]) Indicators() (mean, stddev float64) {
	if l == nil || l.size == 0 {
		return math.NaN(), math.NaN()
	}

	count := 0
	mean = 0.0
	M2 := 0.0 // Sum of squares of differences from the current mean

	// 1. Standard Welford's algorithm for explicitly defined values in the sparse map.
	for _, value := range l.values {
		count++
		v := float64(value)
		delta := v - mean
		mean += delta / float64(count)
		delta2 := v - mean
		M2 += delta * delta2
	}

	// 2. Batched Welford update for the remaining implicit default values.
	// This avoids looping over potentially millions of default values,
	// preserving the O(V) performance characteristic of the sparse series.
	remaining := l.size - len(l.values)
	if remaining > 0 {
		v := float64(l.defaultValue)
		if count == 0 {
			// Fast path: if the series entirely consists of default values,
			// the mean is exactly the default value and the variance is 0.
			mean = v
			// M2 remains 0.0
		} else {
			// Welford's parallel/merge formula:
			// Safely merging a batch of 'k' identical elements (all equal to 'v')
			// into an already processed distribution of size 'nA'.
			k := float64(remaining)
			nA := float64(count)
			nNew := nA + k

			delta := v - mean

			// Update the overall mean combining the existing set and the new batch
			mean += delta * (k / nNew)

			// Update the sum of squared differences
			M2 += delta * delta * (nA * k / nNew)
		}
	}

	// Calculate the population variance (M2 / N).
	// Note: If sample variance were needed, the divisor would be (N - 1).
	variance := M2 / float64(l.size)

	// Safeguard against floating-point inaccuracies that could rarely produce
	// an infinitesimally small negative variance (e.g., -1e-16).
	if variance < 0 {
		variance = 0
	}

	stddev = math.Sqrt(variance)
	return mean, stddev
}

// newLocalSerie is a private constructor that initializes the internal state.
// Implementation choice: It automatically selects the appropriate epsilon-based
// equality function based on the underlying type (float32 vs float64).
func newLocalSerie[F FloatNumber](size int, defaultValue F) *localSerie[F] {
	if size < 0 {
		return nil
	}

	result := new(localSerie[F])
	result.size = size
	result.defaultValue = defaultValue

	// Determine which comparison precision to use
	if isFloat64(defaultValue) {
		result.equality = equalsFloat64
	} else {
		result.equality = equalsFloat32
	}

	result.values = make(map[int]F)
	return result
}

// NewSerie creates and returns a new Serie interface instance.
func NewSerie[F FloatNumber](size int, defaultValue F) Serie[F] {
	return newLocalSerie(size, defaultValue)
}

// NewEmptySerie returns a new empty serie with the default value to set
func NewEmptySerie[F FloatNumber](defaultValue F) Serie[F] {
	return NewSerie(0, defaultValue)
}
