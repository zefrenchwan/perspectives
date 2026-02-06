package commons

import (
	"errors"
	"fmt"
)

// Vectorizer is a function type that converts any input into a ColumnMatrix (vector).
// It returns the resulting ColumnMatrix or an error if the conversion fails.
type Vectorizer func(any) (ColumnMatrix, error)

// ColumnMatrix represents a mathematical column vector (a matrix with a single column).
type ColumnMatrix interface {
	// Add adds another ColumnMatrix of the same size to this one.
	// Returns the resulting ColumnMatrix or an error if dimensions do not match.
	Add(ColumnMatrix) (ColumnMatrix, error)
	// ExternalProduct computes the outer product of this vector and another ColumnMatrix.
	// Returns the resulting SquareMatrix or an error if dimensions do not match.
	ExternalProduct(ColumnMatrix) (SquareMatrix, error)
	// Equals checks if this matrix is equal to another ColumnMatrix.
	// Returns true if both have the same size and contain the same values.
	Equals(ColumnMatrix) bool
	// Export returns the matrix content as a slice of float64.
	Export() []float64
	// Size returns the number of rows in the column matrix.
	Size() int
}

// SquareMatrix represents a square matrix (same number of rows and columns).
type SquareMatrix interface {
	// Equals checks if this matrix is equal to another SquareMatrix.
	// Returns true if both have the same size and contain the same values.
	Equals(SquareMatrix) bool
	// Multiply multiplies this square matrix by a ColumnMatrix (vector).
	// Returns the resulting ColumnMatrix or an error if dimensions do not match.
	Multiply(ColumnMatrix) (ColumnMatrix, error)
	// Export returns the matrix content as a 2D slice of float64.
	Export() [][]float64
	// Row returns the row at the specified index as a slice of float64.
	// Returns an error if the index is out of bounds.
	Row(index int) ([]float64, error)
	// Column returns the column at the specified index as a slice of float64.
	// Returns an error if the index is out of bounds.
	Column(index int) ([]float64, error)
	// Size returns the dimension of the matrix (number of rows or columns).
	Size() int
}

// denseColumnMatrix is an implementation of ColumnMatrix using a slice of float64.
type denseColumnMatrix []float64

// Export returns the raw content of the matrix as a slice of float64.
func (c denseColumnMatrix) Export() []float64 {
	return c
}

// Size returns the number of elements (rows) in the column matrix.
func (c denseColumnMatrix) Size() int {
	return len(c)
}

// Add adds another ColumnMatrix to this one element-wise.
// It returns a new ColumnMatrix containing the sum.
// Returns an error if the sizes of the two matrices do not match.
func (c denseColumnMatrix) Add(d ColumnMatrix) (ColumnMatrix, error) {
	other := d.Export()
	if len(c) != len(other) {
		return nil, errors.New("dimensions are not equal")
	}

	result := make(denseColumnMatrix, len(c))
	for index := 0; index < len(other); index++ {
		result[index] = c[index] + other[index]
	}

	return result, nil
}

// ExternalProduct computes the outer product of this vector and another vector.
// The result is a SquareMatrix where result[i][j] = c[i] * other[j].
// Returns an error if the sizes do not match.
func (c denseColumnMatrix) ExternalProduct(other ColumnMatrix) (SquareMatrix, error) {
	d := other.Export()
	if len(c) != len(d) {
		return nil, errors.New("dimensions are not equal")
	} else if len(c) == 0 {
		return denseSquareMatrix{}, nil
	}

	size := len(c)
	// Optimization: Allocate a single contiguous block for better cache locality
	// and avoid double allocation/copy by constructing denseSquareMatrix directly.
	result := make(denseSquareMatrix, size)
	data := make([]float64, size*size)

	for cIndex := 0; cIndex < size; cIndex++ {
		result[cIndex] = data[cIndex*size : (cIndex+1)*size]
		for dIndex := 0; dIndex < size; dIndex++ {
			result[cIndex][dIndex] = c[cIndex] * d[dIndex]
		}
	}

	return result, nil
}

// Equals compares this matrix with another ColumnMatrix for equality.
// Two matrices are equal if they have the same size and all corresponding elements are equal.
func (c denseColumnMatrix) Equals(other ColumnMatrix) bool {
	d := other.Export()
	if len(c) != len(d) {
		return false
	} else if len(c) == 0 || len(d) == 0 {
		return true
	}

	size := len(c)
	for cIndex := 0; cIndex < size; cIndex++ {
		if !equalsFloats(c[cIndex], d[cIndex]) {
			return false
		}
	}

	return true
}

// NewColumnMatrix creates a new ColumnMatrix from a slice of float64 values.
// If the input slice is empty, it returns an empty matrix.
func NewColumnMatrix(values []float64) ColumnMatrix {
	if len(values) == 0 {
		return denseColumnMatrix{}
	}

	result := make(denseColumnMatrix, len(values))
	copy(result, values)
	return result
}

// denseSquareMatrix is an implementation of SquareMatrix using a 2D slice of float64.
type denseSquareMatrix [][]float64

// Size returns the dimension (number of rows/columns) of the square matrix.
func (s denseSquareMatrix) Size() int {
	return len(s)
}

// Multiply computes the product of this square matrix and a column vector.
// The result is a ColumnMatrix representing the matrix-vector product.
// Returns an error if the matrix is empty or if dimensions are incompatible.
func (s denseSquareMatrix) Multiply(other ColumnMatrix) (ColumnMatrix, error) {
	c := other.Export()
	size := len(c)

	if len(s) != size {
		return nil, errors.New("dimensions are not equal")
	}

	result := make([]float64, size)
	for index, values := range s {
		if len(values) != size {
			return nil, errors.New("invalid matrix size")
		}

		resultValue := 0.0
		for colIndex, value := range values {
			resultValue += value * c[colIndex]
		}

		result[index] = resultValue
	}

	return denseColumnMatrix(result), nil
}

// Equals compares this matrix with another SquareMatrix for equality.
// Two matrices are equal if they have the same size and all corresponding elements are equal.
func (s denseSquareMatrix) Equals(other SquareMatrix) bool {
	d := other.Export()
	if len(s) != len(d) {
		return false
	} else if len(s) == 0 || len(d) == 0 {
		return true
	}

	rows := len(s)
	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		if len(s[rowIndex]) != len(d[rowIndex]) {
			return false
		}

		cols := len(s[rowIndex])
		for colIndex := 0; colIndex < cols; colIndex++ {
			if !equalsFloats(s[rowIndex][colIndex], d[rowIndex][colIndex]) {
				return false
			}
		}
	}

	return true
}

// Export returns the raw content of the matrix as a 2D slice of float64.
func (s denseSquareMatrix) Export() [][]float64 {
	return s
}

// Row returns a copy of the row at the specified index.
// Returns an error if the index is out of bounds.
func (s denseSquareMatrix) Row(index int) ([]float64, error) {
	length := len(s)
	if index < 0 || index >= length {
		return nil, errors.New("invalid matrix index")
	}

	result := make([]float64, length)
	copy(result, s[index])
	return result, nil
}

// Column returns a copy of the column at the specified index.
// Returns an error if the index is out of bounds.
func (s denseSquareMatrix) Column(index int) ([]float64, error) {
	length := len(s)
	if index < 0 || index >= length {
		return nil, errors.New("invalid matrix index")
	}

	result := make([]float64, length)
	for row := 0; row < length; row++ {
		result[row] = s[row][index]
	}

	return result, nil
}

// NewSquareMatrix creates a new SquareMatrix from a 2D slice of float64 values.
// It validates that the input is a square matrix of the specified size.
// Returns an error if the input dimensions do not match the specified size.
func NewSquareMatrix(size int, elements [][]float64) (SquareMatrix, error) {
	if len(elements) != size {
		return nil, errors.New("invalid matrix size")
	}

	// Validate all rows before allocation
	for index, values := range elements {
		if len(values) != size {
			return nil, fmt.Errorf("invalid matrix size at index %d", index)
		}
	}

	// Optimization: Allocate a single contiguous block for better cache locality
	result := make(denseSquareMatrix, size)
	data := make([]float64, size*size)

	for index, values := range elements {
		result[index] = data[index*size : (index+1)*size]
		copy(result[index], values)
	}

	return result, nil
}
