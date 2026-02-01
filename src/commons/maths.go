package commons

import (
	"errors"
	"fmt"
)

// ColumnMatrix is a multi lines, single column, matrix
type ColumnMatrix []float64

// Add two columns matrixs and return the result.
// May raise an error if sizes will not match
func (c ColumnMatrix) Add(d ColumnMatrix) (ColumnMatrix, error) {
	if len(c) != len(d) {
		return nil, errors.New("dimensions are not equal")
	}

	result := make(ColumnMatrix, len(c))
	for index := 0; index < len(d); index++ {
		result[index] = c[index] + d[index]
	}

	return result, nil
}

// ExternalProduct returns the external product if sizes match, false otherwise
func (c ColumnMatrix) ExternalProduct(d ColumnMatrix) (SquareMatrix, error) {
	if len(c) != len(d) {
		return nil, errors.New("dimensions are not equal")
	} else if len(c) == 0 || len(d) == 0 {
		return nil, nil
	}

	size := len(c)
	result := make(SquareMatrix, size)
	for cIndex := 0; cIndex < size; cIndex++ {
		result[cIndex] = make([]float64, size)
		for dIndex := 0; dIndex < size; dIndex++ {
			result[cIndex][dIndex] = c[cIndex] * d[dIndex]
		}
	}

	return result, nil
}

// Equals tests equality (same size, same values)
func (c ColumnMatrix) Equals(d ColumnMatrix) bool {
	if len(c) != len(d) {
		return false
	} else if len(c) == 0 || len(d) == 0 {
		return true
	}

	size := len(c)
	for cIndex := 0; cIndex < size; cIndex++ {
		if c[cIndex] != d[cIndex] {
			return false
		}
	}

	return true
}

// NewColumnMatrix returns a column matrix (nil for 0 value)
func NewColumnMatrix(values []float64) ColumnMatrix {
	if len(values) == 0 {
		return nil
	}

	result := make(ColumnMatrix, len(values))
	copy(result, values)
	return result
}

// SquareMatrix is a square matrix with a given size (same by definition for lines and columns)
type SquareMatrix [][]float64

// Dimensions returns the size of the matrix (same size, returned as a couple)
func (s SquareMatrix) Dimensions() (int, int) {
	return len(s), len(s)
}

// Multiply calculated $s \times c$, that is a column matrix with len(c) lines and a single column
func (s SquareMatrix) Multiply(c ColumnMatrix) (ColumnMatrix, error) {
	size := len(c)
	if size == 0 {
		return ColumnMatrix{}, errors.New("empty matrix")
	}

	result := make(ColumnMatrix, size)
	for index, values := range s {
		if len(values) != size {
			return ColumnMatrix{}, errors.New("invalid matrix size")
		}

		resultValue := 0.0
		for colIndex, value := range values {
			resultValue += value * c[colIndex]
		}

		result[index] = resultValue
	}

	return result, nil
}

// Equals returns true if sizes and values are equals from each other
func (s SquareMatrix) Equals(d SquareMatrix) bool {
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
			if s[rowIndex][colIndex] != d[rowIndex][colIndex] {
				return false
			}
		}
	}

	return true
}

// NewSquareMatrix returns a complete square matrix with provided content, or an error if sizes mismatch
func NewSquareMatrix(size int, elements [][]float64) (SquareMatrix, error) {
	result := make(SquareMatrix, size)
	if len(elements) != size {
		return nil, errors.New("invalid matrix size")
	}

	for index, values := range elements {
		if len(values) != size {
			return nil, fmt.Errorf("invalid matrix size at index %d", index)
		}
		result[index] = make([]float64, size)
		copy(result[index], values)
	}

	return result, nil
}
