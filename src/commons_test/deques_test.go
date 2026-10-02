package commons_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestDequeAdd(t *testing.T) {
	q := commons.NewDeque[int](100)
	q.AddFirst(10)
	q.AddLast(100)
	q.AddFirst(1)

	if q.Size() != 3 {
		t.Errorf("Expected size 3, got %v", q.Size())
		t.Fail()
	}

	result := slices.Collect(q.Elements())
	if !slices.Equal(result, []int{1, 10, 100}) {
		t.Errorf("Expected [1,10,100], got %v", result)
		t.Fail()
	}
}

func TestDequePop(t *testing.T) {
	q := commons.NewDeque[int](100)
	q.AddFirst(10)
	q.AddLast(100)
	q.AddFirst(1)

	if v1, has1 := q.PopFirst(); !has1 || v1 != 1 {
		t.Errorf("Expected 1, got %v", v1)
		t.Fail()
	} else if v2, has2 := q.PopLast(); !has2 || v2 != 100 {
		t.Errorf("Expected 100, got %v", v2)
		t.Fail()
	} else if v3, has3 := q.PopFirst(); !has3 || v3 != 10 {
		t.Errorf("Expected 10, got %v", v3)
		t.Fail()
	} else if _, has := q.PopLast(); has {
		t.Errorf("Expected empty, got %v", has)
		t.Fail()
	}
}

func TestDequeGet(t *testing.T) {
	q := commons.NewDeque[int](10)
	q.AddFirst(2)
	q.AddFirst(1)

	if q.Size() != 2 {
		t.Errorf("Expected size 2, got %v", q.Size())
		t.Fail()
	}

	for i := 0; i < q.Size(); i++ {
		if v, has := q.Get(i); !has {
			t.Errorf("Expected value at index %v, got empty", i)
			t.Fail()
		} else if v != i+1 {
			t.Errorf("Expected value %v at index %v, got %v", i+1, i, v)
			t.Fail()
		}
	}

	q.Clear()
	for i := 0; i < q.Size(); i++ {
		t.Log("expected empty")
		t.Fail()
	}

}
