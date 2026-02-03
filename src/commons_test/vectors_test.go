package commons_test

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestColumnMatrixAdd(t *testing.T) {
	a := commons.NewColumnMatrix([]float64{1, 2, 3})
	b := commons.NewColumnMatrix([]float64{3, 2, 1})
	noMatch := commons.NewColumnMatrix([]float64{1, 2, 3, 4})

	if _, err := noMatch.Add(a); err == nil {
		t.Log("size mismatch not seen")
		t.Fail()
	}

	expected := commons.NewColumnMatrix([]float64{4, 4, 4})
	if result, err := a.Add(b); err != nil {
		t.Log(err)
		t.Fail()
	} else if !result.Equals(expected) {
		t.Log(result)
		t.Fail()
	}
}

func TestColumnMatrixExternalProduct(t *testing.T) {
	a := commons.NewColumnMatrix([]float64{1, 2})
	b := commons.NewColumnMatrix([]float64{3, 4})
	noMatch := commons.NewColumnMatrix([]float64{1, 2, 3})

	if _, err := a.ExternalProduct(noMatch); err == nil {
		t.Log("size mismatch not seen")
		t.Fail()
	}
	expected, _ := commons.NewSquareMatrix(2, [][]float64{{3, 4}, {6, 8}})
	if result, err := a.ExternalProduct(b); err != nil {
		t.Log(err)
		t.Fail()
	} else if !result.Equals(expected) {
		t.Log(result)
		t.Fail()
	}
}

func TestColumnMatrixEquals(t *testing.T) {
	a := commons.NewColumnMatrix([]float64{1, 2, 3})
	b := commons.NewColumnMatrix([]float64{1, 2, 3})
	c := commons.NewColumnMatrix([]float64{1, 2, 4})
	d := commons.NewColumnMatrix([]float64{1, 2})

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
	if _, err := commons.NewSquareMatrix(2, [][]float64{{1, 2}, {3, 4}}); err != nil {
		t.Log("Valid matrix creation failed")
		t.Fail()
	}
	if _, err := commons.NewSquareMatrix(2, [][]float64{{1, 2}}); err == nil {
		t.Log("Invalid matrix (missing row) creation succeeded")
		t.Fail()
	}
	if _, err := commons.NewSquareMatrix(2, [][]float64{{1, 2}, {3}}); err == nil {
		t.Log("Invalid matrix (short row) creation succeeded")
		t.Fail()
	}
}

func TestSquareMatrixEquals(t *testing.T) {
	m1, _ := commons.NewSquareMatrix(2, [][]float64{{1, 2}, {3, 4}})
	m2, _ := commons.NewSquareMatrix(2, [][]float64{{1, 2}, {3, 4}})
	m3, _ := commons.NewSquareMatrix(2, [][]float64{{1, 2}, {3, 5}})

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
	id, _ := commons.NewSquareMatrix(2, [][]float64{{1, 0}, {0, 1}})
	vec := commons.NewColumnMatrix([]float64{2, 3})

	if res, err := id.Multiply(vec); err != nil {
		t.Log(err)
		t.Fail()
	} else if !res.Equals(vec) {
		t.Logf("Identity multiplication failed. Got %v, expected %v", res, vec)
		t.Fail()
	}

	// Permutation
	perm, _ := commons.NewSquareMatrix(2, [][]float64{{0, 1}, {1, 0}})
	expected := commons.NewColumnMatrix([]float64{3, 2})

	if res, err := perm.Multiply(vec); err != nil {
		t.Log(err)
		t.Fail()
	} else if !res.Equals(expected) {
		t.Logf("Permutation multiplication failed. Got %v, expected %v", res, expected)
		t.Fail()
	}

	// Size mismatch
	wrongVec := commons.NewColumnMatrix([]float64{1, 2, 3})
	if _, err := id.Multiply(wrongVec); err == nil {
		t.Log("Size mismatch not detected")
		t.Fail()
	}
}
