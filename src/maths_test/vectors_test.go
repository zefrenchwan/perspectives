package maths_test

import (
	"math"
	"testing"

	"github.com/zefrenchwan/perspectives.git/maths"
)

func TestColumnMatrixAdd(t *testing.T) {
	a := maths.NewColumnMatrix([]float64{1, 2, 3})
	b := maths.NewColumnMatrix([]float64{3, 2, 1})
	noMatch := maths.NewColumnMatrix([]float64{1, 2, 3, 4})

	if _, err := noMatch.Add(a); err == nil {
		t.Log("size mismatch not seen")
		t.Fail()
	}

	expected := maths.NewColumnMatrix([]float64{4, 4, 4})
	if result, err := a.Add(b); err != nil {
		t.Log(err)
		t.Fail()
	} else if !result.Equals(expected) {
		t.Log(result)
		t.Fail()
	}
}

func TestColumnMatrixExternalProduct(t *testing.T) {
	a := maths.NewColumnMatrix([]float64{1, 2})
	b := maths.NewColumnMatrix([]float64{3, 4})
	noMatch := maths.NewColumnMatrix([]float64{1, 2, 3})

	if _, err := a.ExternalProduct(noMatch); err == nil {
		t.Log("size mismatch not seen")
		t.Fail()
	}
	expected, _ := maths.NewSquareMatrix(2, [][]float64{{3, 4}, {6, 8}})
	if result, err := a.ExternalProduct(b); err != nil {
		t.Log(err)
		t.Fail()
	} else if !result.Equals(expected) {
		t.Log(result)
		t.Fail()
	}
}

func TestColumnMatrixDotProduct(t *testing.T) {
	a := maths.NewColumnMatrix([]float64{1, 2})
	b := maths.NewColumnMatrix([]float64{3, 4})
	noMatch := maths.NewColumnMatrix([]float64{1, 2, 3})

	if _, err := a.DotProduct(noMatch); err == nil {
		t.Log("size mismatch not seen")
		t.Fail()
	}
	expected := 11.0
	if result, err := a.DotProduct(b); err != nil {
		t.Log(err)
		t.Fail()
	} else if math.Abs(expected-result) > 0.0001 {
		t.Log(result)
		t.Fail()
	}
}

func TestColumnMatrixEquals(t *testing.T) {
	a := maths.NewColumnMatrix([]float64{1, 2, 3})
	b := maths.NewColumnMatrix([]float64{1, 2, 3})
	c := maths.NewColumnMatrix([]float64{1, 2, 4})
	d := maths.NewColumnMatrix([]float64{1, 2})

	if !a.Equals(b) {
		t.Log("Equal matrices considered different")
		t.Fail()
	}
	if a.Equals(c) {
		t.Log("Different matrices considered equal")
		t.Fail()
	}
	if a.Equals(d) {
		t.Log("Matrices of different sizes considered equal")
		t.Fail()
	}
}

func TestNewSquareMatrix(t *testing.T) {
	if _, err := maths.NewSquareMatrix(2, [][]float64{{1, 2}, {3, 4}}); err != nil {
		t.Log("Valid matrix creation failed")
		t.Fail()
	}
	if _, err := maths.NewSquareMatrix(2, [][]float64{{1, 2}}); err == nil {
		t.Log("Invalid matrix (missing row) creation succeeded")
		t.Fail()
	}
	if _, err := maths.NewSquareMatrix(2, [][]float64{{1, 2}, {3}}); err == nil {
		t.Log("Invalid matrix (short row) creation succeeded")
		t.Fail()
	}
}

func TestSquareMatrixEquals(t *testing.T) {
	m1, _ := maths.NewSquareMatrix(2, [][]float64{{1, 2}, {3, 4}})
	m2, _ := maths.NewSquareMatrix(2, [][]float64{{1, 2}, {3, 4}})
	m3, _ := maths.NewSquareMatrix(2, [][]float64{{1, 2}, {3, 5}})

	if !m1.Equals(m2) {
		t.Log("Equal matrices considered different")
		t.Fail()
	}
	if m1.Equals(m3) {
		t.Log("Different matrices considered equal")
		t.Fail()
	}
}

func TestSquareMatrixMultiply(t *testing.T) {
	// Identity matrix
	id, _ := maths.NewSquareMatrix(2, [][]float64{{1, 0}, {0, 1}})
	vec := maths.NewColumnMatrix([]float64{2, 3})

	if res, err := id.Multiply(vec); err != nil {
		t.Log(err)
		t.Fail()
	} else if !res.Equals(vec) {
		t.Logf("Identity multiplication failed. Got %v, expected %v", res, vec)
		t.Fail()
	}

	// Permutation
	perm, _ := maths.NewSquareMatrix(2, [][]float64{{0, 1}, {1, 0}})
	expected := maths.NewColumnMatrix([]float64{3, 2})

	if res, err := perm.Multiply(vec); err != nil {
		t.Log(err)
		t.Fail()
	} else if !res.Equals(expected) {
		t.Logf("Permutation multiplication failed. Got %v, expected %v", res, expected)
		t.Fail()
	}

	// Size mismatch
	wrongVec := maths.NewColumnMatrix([]float64{1, 2, 3})
	if _, err := id.Multiply(wrongVec); err == nil {
		t.Log("Size mismatch not detected")
		t.Fail()
	}
}

func TestSquareMatrixAdd(t *testing.T) {
	m1, _ := maths.NewSquareMatrix(2, [][]float64{{1, 2}, {3, 4}})
	m2, _ := maths.NewSquareMatrix(2, [][]float64{{5, 6}, {7, 8}})
	expected, _ := maths.NewSquareMatrix(2, [][]float64{{6, 8}, {10, 12}})

	if res, err := m1.Add(m2); err != nil {
		t.Log(err)
		t.Fail()
	} else if !res.Equals(expected) {
		t.Logf("Addition failed. Got %v, expected %v", res, expected)
		t.Fail()
	}

	// Size mismatch
	wrongSize, _ := maths.NewSquareMatrix(3, [][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}})
	if _, err := m1.Add(wrongSize); err == nil {
		t.Log("Size mismatch not detected")
		t.Fail()
	}
}

func TestSquareMatrixRow(t *testing.T) {
	m, _ := maths.NewSquareMatrix(2, [][]float64{{1, 2}, {3, 4}})

	// Valid row
	row0, err := m.Row(0)
	if err != nil {
		t.Logf("Failed to get row 0: %v", err)
		t.Fail()
	}
	if len(row0) != 2 || row0[0] != 1 || row0[1] != 2 {
		t.Logf("Row 0 incorrect. Got %v", row0)
		t.Fail()
	}

	row1, err := m.Row(1)
	if err != nil {
		t.Logf("Failed to get row 1: %v", err)
		t.Fail()
	}
	if len(row1) != 2 || row1[0] != 3 || row1[1] != 4 {
		t.Logf("Row 1 incorrect. Got %v", row1)
		t.Fail()
	}

	// Invalid row
	if _, err := m.Row(-1); err == nil {
		t.Log("Negative row index not detected")
		t.Fail()
	}
	if _, err := m.Row(2); err == nil {
		t.Log("Out of bounds row index not detected")
		t.Fail()
	}
}

func TestSquareMatrixColumn(t *testing.T) {
	m, _ := maths.NewSquareMatrix(2, [][]float64{{1, 2}, {3, 4}})

	// Valid column
	col0, err := m.Column(0)
	if err != nil {
		t.Logf("Failed to get column 0: %v", err)
		t.Fail()
	}
	if len(col0) != 2 || col0[0] != 1 || col0[1] != 3 {
		t.Logf("Column 0 incorrect. Got %v", col0)
		t.Fail()
	}

	col1, err := m.Column(1)
	if err != nil {
		t.Logf("Failed to get column 1: %v", err)
		t.Fail()
	}
	if len(col1) != 2 || col1[0] != 2 || col1[1] != 4 {
		t.Logf("Column 1 incorrect. Got %v", col1)
		t.Fail()
	}

	// Invalid column
	if _, err := m.Column(-1); err == nil {
		t.Log("Negative column index not detected")
		t.Fail()
	}
	if _, err := m.Column(2); err == nil {
		t.Log("Out of bounds column index not detected")
		t.Fail()
	}
}

func TestColumnMatrixSize(t *testing.T) {
	c := maths.NewColumnMatrix([]float64{1, 2, 3})
	if c.Size() != 3 {
		t.Logf("Incorrect size for ColumnMatrix. Expected 3, got %d", c.Size())
		t.Fail()
	}

	empty := maths.NewColumnMatrix([]float64{})
	if empty.Size() != 0 {
		t.Logf("Incorrect size for empty ColumnMatrix. Expected 0, got %d", empty.Size())
		t.Fail()
	}
}

func TestSquareMatrixSize(t *testing.T) {
	m, _ := maths.NewSquareMatrix(2, [][]float64{{1, 2}, {3, 4}})
	if m.Size() != 2 {
		t.Logf("Incorrect size for SquareMatrix. Expected 2, got %d", m.Size())
		t.Fail()
	}
}

func TestSquareMatrixExport(t *testing.T) {
	data := [][]float64{{1, 2}, {3, 4}}
	m, _ := maths.NewSquareMatrix(2, data)
	exported := m.Export()

	if len(exported) != 2 {
		t.Log("Exported matrix has wrong number of rows")
		t.Fail()
	}
	for i := range data {
		for j := range data[i] {
			if exported[i][j] != data[i][j] {
				t.Logf("Exported data mismatch at %d,%d", i, j)
				t.Fail()
			}
		}
	}
}

func TestColumnMatrixExport(t *testing.T) {
	data := []float64{1, 2, 3}
	c := maths.NewColumnMatrix(data)
	exported := c.Export()

	if len(exported) != 3 {
		t.Log("Exported vector has wrong length")
		t.Fail()
	}
	for i := range data {
		if exported[i] != data[i] {
			t.Logf("Exported data mismatch at %d", i)
			t.Fail()
		}
	}
}
