package periods_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestPartitionEmpty(t *testing.T) {
	partition := periods.NewDynamicPartition[int]("int", func(a int, b int) bool { return a == b })
	emptySet := periods.NewDynamicSet[int]("int", func(a int, b int) bool { return a == b })
	if !partition.IsEmpty() {
		t.Errorf("Expected partition to be empty")
	} else if partition.DataType() != "int" {
		t.Errorf("Expected partition to have data type int")
	} else if !partition.Domain().Equals(periods.NewEmptyPeriod()) {
		t.Errorf("partition domain should be empty")
	} else if partition.Equals(emptySet) {
		t.Errorf("partition should not be equal to empty set")
	} else if !partition.Equals(partition) {
		t.Errorf("partition should be equal to itself")
	} else if _, has := partition.At(time.Now()); has {
		t.Errorf("Expected partition.At(time.Now()) to be nil")
	}

	partition.Add(1, periods.NewEmptyPeriod())
	if !partition.IsEmpty() {
		t.Errorf("Expected partition to be empty after adding empty period")
	}
}

func TestPartitionAt(t *testing.T) {
	partition := periods.NewDynamicPartition[int]("int", func(a int, b int) bool { return a == b })
	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)
	if _, has := partition.At(now); has {
		t.Errorf("Expected partition.At(time.Now()) to be nil")
	}

	partition.Add(1, periods.NewPeriodSince(now, true))
	partition.Add(2, periods.NewPeriodSince(after, true))

	if matching, has := partition.At(now); !has {
		t.Errorf("Expected partition.At(time.Now()) to be non-nil")
	} else if matching != 1 {
		t.Errorf("Expected partition.At(time.Now()) to be 1")
	}

	if _, has := partition.At(before); has {
		t.Errorf("Expected no element at before")
	}

	if matching, has := partition.At(after); !has {
		t.Errorf("Expected partition.At(after) to be non-nil")
	} else if matching != 2 {
		t.Errorf("Expected partition.At(after)) to be 2")
	}
}

func TestPartitionRemove(t *testing.T) {
	partition := periods.NewDynamicPartition[int]("int", func(a int, b int) bool { return a == b })
	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)

	partition.Add(1, periods.NewPeriodSince(now, true))
	partition.Add(2, periods.NewPeriodSince(after, true))
	partition.Remove(periods.NewPeriodSince(after, true))

	if matching, has := partition.At(now); !has {
		t.Errorf("Expected partition.At(time.Now()) to be non-nil")
	} else if matching != 1 {
		t.Errorf("Expected partition.At(time.Now()) to be 1")
	}

	if _, has := partition.At(before); has {
		t.Errorf("Expected no element at before")
	}

	if _, has := partition.At(after); has {
		t.Errorf("Expected partition.At(after) to be empty due to removal")
	}
}

func TestPartitionDestructiveAdd(t *testing.T) {
	partition := periods.NewDynamicPartition[string]("string", func(a, b string) bool { return a == b })
	now := time.Now().Truncate(time.Hour)
	before := now.Add(-10 * time.Hour)
	after := now.Add(10 * time.Hour)

	// Add Alice for [before, after]
	partition.Add("Alice", periods.NewFinitePeriod(before, after, true, true))

	// Add Bob for [now, after]
	// Because it's a partition, this should truncate Alice's period to [before, now[
	partition.Add("Bob", periods.NewFinitePeriod(now, after, true, true))

	var alicePeriod periods.Period
	var bobPeriod periods.Period
	var aliceCount, bobCount int

	for p, v := range partition.Range() {
		if v == "Alice" {
			alicePeriod = p
			aliceCount++
		} else if v == "Bob" {
			bobPeriod = p
			bobCount++
		}
	}

	if aliceCount != 1 {
		t.Errorf("Expected Alice to have exactly 1 matching period, got %d", aliceCount)
	}
	if bobCount != 1 {
		t.Errorf("Expected Bob to have exactly 1 matching period, got %d", bobCount)
	}

	expectedAlicePeriod := periods.NewFinitePeriod(before, now, true, false)
	if !alicePeriod.Equals(expectedAlicePeriod) {
		t.Errorf("Expected Alice's period to be truncated to %s, got %s", expectedAlicePeriod.AsRawString(), alicePeriod.AsRawString())
	}

	expectedBobPeriod := periods.NewFinitePeriod(now, after, true, true)
	if !bobPeriod.Equals(expectedBobPeriod) {
		t.Errorf("Expected Bob's period to be %s, got %s", expectedBobPeriod.AsRawString(), bobPeriod.AsRawString())
	}
}
