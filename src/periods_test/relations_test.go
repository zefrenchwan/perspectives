package periods_test

import (
	"slices"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestTimeRelationEmpty(t *testing.T) {
	currentSet := periods.NewTimeRelation[int]("int", func(a, b int) bool { return a == b })
	if !currentSet.IsEmpty() {
		t.Errorf("Expected set to be empty")
	} else if currentSet.DataType() != "int" {
		t.Errorf("Expected data type to be int")
	} else if !currentSet.Domain().IsEmpty() {
		t.Errorf("Expected set domain to be empty")
	} else if _, has := currentSet.At(time.Now()); has {
		t.Errorf("Expected set to not contain value at time.Now()")
	}

	counter := 0
	for p, v := range currentSet.Range() {
		_ = p
		_ = v
		counter++
	}

	if counter != 0 {
		t.Errorf("Expected set to be empty")
	}
}

func TestTimeRelationAt(t *testing.T) {
	currentSet := periods.NewTimeRelation[int]("int", func(a, b int) bool { return a == b })
	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-20, 0, 0)
	after := now.AddDate(20, 0, 0)

	currentSet.Add(1, periods.NewPeriodSince(before, true))
	currentSet.Add(2, periods.NewPeriodSince(after, true))

	if !currentSet.Domain().Equals(periods.NewPeriodSince(before, true)) {
		t.Errorf("set domain should be [before,+oo[, it is %s", currentSet.Domain().AsRawString())
	}

	// at after, we expect 2 values
	if iterator, has := currentSet.At(after); !has {
		t.Errorf("Expected set to contain 2 values at after")
	} else if values := slices.Collect(iterator); len(values) != 2 {
		t.Errorf("Expected set to contain 2 value at after")
	} else if !slices.Contains(values, 2) {
		t.Errorf("Expected set to contain value 2 at after")
	} else if !slices.Contains(values, 1) {
		t.Errorf("Expected set to contain value 1 at after")
	}

	// at now, we expect 1 value because we added at after
	if iterator, has := currentSet.At(now); !has {
		t.Errorf("Expected set to contain 1 value at now")
	} else if values := slices.Collect(iterator); len(values) != 1 {
		t.Errorf("Expected set to contain 1 value at now")
	} else if !slices.Contains(values, 1) {
		t.Errorf("Expected set to contain value 1 at now")
	}
}

func TestTimeRelationRemove(t *testing.T) {
	currentSet := periods.NewTimeRelation[int]("int", func(a, b int) bool { return a == b })
	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-20, 0, 0)
	after := now.AddDate(20, 0, 0)

	currentSet.Add(1, periods.NewPeriodSince(before, true))
	currentSet.Add(2, periods.NewPeriodSince(after, true))

	// let us remove everything since now
	currentSet.Remove(periods.NewPeriodSince(now, true))

	// at after, no value
	if _, has := currentSet.At(after); has {
		t.Errorf("Expected set to contain nothing at after")
	}

	// Now, remove it all
	currentSet.Remove(periods.NewFullPeriod())
	if !currentSet.IsEmpty() {
		t.Errorf("Expected set to be empty")
	}
}

func TestTimeRelationEquals(t *testing.T) {
	currentSet := periods.NewTimeRelation[int]("int", func(a, b int) bool { return a == b })
	if !currentSet.Equals(currentSet) {
		t.Errorf("Expected set to be equal to itself")
	}

	newSet := currentSet.Copy()
	if !currentSet.Equals(newSet) {
		t.Errorf("Expected set to be equal to its copy")
	}

	newSet.Add(1, periods.NewEmptyPeriod())
	if !currentSet.Equals(newSet) {
		t.Errorf("Expected set to be equal to its copy after adding an empty period")
	}

	newSet.Remove(periods.NewEmptyPeriod())
	if !currentSet.Equals(newSet) {
		t.Errorf("Expected set to be equal to its copy after removing an empty period")
	}

	newSet.Add(1, periods.NewFullPeriod())
	if currentSet.Equals(newSet) {
		t.Errorf("Expected set to NOT be equal to its copy after adding a non empty period")
	} else if !newSet.Equals(newSet) {
		t.Errorf("Expected set to be equal to itself")
	}
}

func TestTimeRelationRange(t *testing.T) {
	currentSet := periods.NewTimeRelation[int]("int", func(a, b int) bool { return a == b })
	now := time.Now().Truncate(time.Second)

	// Add multiple values at the exact same period
	currentSet.Add(1, periods.NewPeriodSince(now, true))
	currentSet.Add(2, periods.NewPeriodSince(now, true))
	currentSet.Add(3, periods.NewPeriodSince(now, true))

	counterAt := 0
	if iterator, has := currentSet.At(now); !has {
		t.Errorf("Expected to have elements at now")
	} else {
		for v := range iterator {
			_ = v
			counterAt++
		}
	}

	if counterAt != 3 {
		t.Errorf("Expected At iterator to stop after 3 iteration, but ran %d times", counterAt)
	}
}

func TestTimeRelationCopy(t *testing.T) {
	r := periods.NewTimeRelation[int]("int", func(a, b int) bool { return a == b })
	r.Add(10, periods.NewFullPeriod())
	r.Add(20, periods.NewFullPeriod())

	if v, has := r.At(time.Now()); !has {
		t.Errorf("Expected to have elements at now")
	} else if values := slices.Collect(v); values == nil {
		t.Errorf("Expected iterator to not be nil")
	} else if len(values) != 2 {
		t.Errorf("Expected iterator to have 2 values, but had %d", len(values))
	}

	rCopy := r.Copy()
	if !r.Equals(rCopy) {
		t.Errorf("Expected r to be equal to its copy")
	} else if v, has := rCopy.At(time.Now()); !has {
		t.Errorf("Expected to have elements at now")
	} else if values := slices.Collect(v); values == nil {
		t.Errorf("Expected iterator to not be nil")
	} else if len(values) != 2 {
		t.Errorf("Expected iterator to have 2 values, but had %d", len(values))
	}

	r.Remove(periods.NewFullPeriod())
	if v, has := rCopy.At(time.Now()); !has {
		t.Errorf("Expected to have elements at now")
	} else if values := slices.Collect(v); values == nil {
		t.Errorf("Expected iterator to not be nil")
	} else if len(values) != 2 {
		t.Errorf("Expected iterator to have 2 values, but had %d", len(values))
	}
}
