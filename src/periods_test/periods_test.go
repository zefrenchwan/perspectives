package periods_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestPeriodIsEmpty(t *testing.T) {
	if !periods.NewEmptyPeriod().IsEmpty() {
		t.Log("NewEmptyPeriod should be empty")
		t.Fail()
	}
	if periods.NewFullPeriod().IsEmpty() {
		t.Log("NewFullPeriod should not be empty")
		t.Fail()
	}
}

func TestPeriodCopy(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	original := periods.NewPeriodSince(now, true)
	copied := original.Copy()

	if !copied.Equals(original) {
		t.Log("Copied period should equal the original")
		t.Fail()
	}
}

func TestPeriodComplements(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	value := periods.NewPeriodSince(now, true)
	complement := value.Complement()
	expected := periods.NewPeriodUntil(now, false)
	if !expected.Equals(complement) {
		t.Logf("Complement failed, expected %s got %s", expected.AsRawString(), complement.AsRawString())
		t.Fail()
	}

	value = periods.NewPeriodUntil(now, true)
	complement = value.Complement()
	expected = periods.NewPeriodSince(now, false)
	if !expected.Equals(complement) {
		t.Logf("Complement failed, expected %s got %s", expected.AsRawString(), complement.AsRawString())
		t.Fail()
	}

	value = periods.NewEmptyPeriod()
	complement = value.Complement()
	expected = periods.NewFullPeriod()
	if !expected.Equals(complement) {
		t.Logf("Complement failed, expected %s got %s", expected.AsRawString(), complement.AsRawString())
		t.Fail()
	}

	value = periods.NewFullPeriod()
	complement = value.Complement()
	expected = periods.NewEmptyPeriod()
	if !expected.Equals(complement) {
		t.Log("Complement failed to reverse full to empty")
		t.Fail()
	}
}

func TestFinitePeriodComplement(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	before := now.Add(-10 * time.Hour)
	after := now.Add(10 * time.Hour)

	// Period is [before, after]
	value := periods.NewFinitePeriod(before, after, true, true)
	complement := value.Complement()

	// Expected is ]-oo, before[ U ]after, +oo[
	expectedLeft := periods.NewPeriodUntil(before, false)
	expectedRight := periods.NewPeriodSince(after, false)
	expected := expectedLeft.Union(expectedRight)

	if !expected.Equals(complement) {
		t.Logf("Finite complement failed, expected %s got %s", expected.AsRawString(), complement.AsRawString())
		t.Fail()
	}
}

func TestIntersectionWithFull(t *testing.T) {
	value := periods.NewFullPeriod()
	now := time.Now()
	otherValue := periods.NewPeriodSince(now, true)
	result := otherValue.Intersection(value)
	if !result.Equals(otherValue) {
		t.Log("intersection with full failed")
		t.Fail()
	}
}

func TestIntersectionWithEmpty(t *testing.T) {
	value := periods.NewEmptyPeriod()
	now := time.Now()
	otherValue := periods.NewPeriodSince(now, true)
	result := otherValue.Intersection(value)
	if !result.Equals(value) {
		t.Log("intersection with empty failed")
		t.Fail()
	}
}

func TestIntersectionWithOther(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := time.Now().Add(-24 * time.Hour)
	value := periods.NewPeriodSince(before, true)
	otherValue := periods.NewPeriodUntil(now, true)
	expected := periods.NewFinitePeriod(before, now, true, true)
	result := otherValue.Intersection(value)
	if !result.Equals(expected) {
		t.Logf("intersection with other failed: got %s BUT EXPECTED %s", result.AsRawString(), expected.AsRawString())
		t.Fail()
	}
}

func TestIntersectionDisjoint(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	t1 := now.Add(1 * time.Hour)
	t2 := now.Add(2 * time.Hour)
	t3 := now.Add(3 * time.Hour)
	t4 := now.Add(4 * time.Hour)

	// [t1, t2] intersection [t3, t4] -> mathematically empty
	p1 := periods.NewFinitePeriod(t1, t2, true, true)
	p2 := periods.NewFinitePeriod(t3, t4, true, true)

	result := p1.Intersection(p2)
	if !result.IsEmpty() {
		t.Logf("Disjoint intersection should be empty, got %s", result.AsRawString())
		t.Fail()
	}
}

func TestUnionWithEmpty(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	value := periods.NewPeriodSince(now, true)
	expected := value
	result := periods.NewEmptyPeriod().Union(value)
	if !result.Equals(expected) {
		t.Logf("union with empty failed: got %s", result.AsRawString())
		t.Fail()
	}

	result = value.Union(periods.NewEmptyPeriod())
	if !result.Equals(expected) {
		t.Logf("union with empty failed: got %s", result.AsRawString())
		t.Fail()
	}
}

func TestUnionWithSame(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	value := periods.NewPeriodSince(now, true)
	expected := value
	result := value.Union(value)
	if !result.Equals(expected) {
		t.Logf("union with same failed: got %s", result.AsRawString())
		t.Fail()
	}
}

func TestInfiniteUnion(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := time.Now().Add(-24 * time.Hour).Truncate(time.Second)
	first := periods.NewPeriodSince(before, true)
	second := periods.NewPeriodUntil(now, true)
	expected := periods.NewFullPeriod()
	result := first.Union(second)
	if !result.Equals(expected) {
		t.Logf("union to make full failed: got %s", result.AsRawString())
		t.Fail()
	}

	result = second.Union(first)
	if !result.Equals(expected) {
		t.Logf("union to make full failed: got %s", result.AsRawString())
		t.Fail()
	}
}

func TestComplexUnion(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	t1 := now.Add(1 * time.Hour)
	t2 := now.Add(2 * time.Hour)
	t3 := now.Add(3 * time.Hour)
	t4 := now.Add(4 * time.Hour)

	// Overlapping union: [t1, t3] U [t2, t4] -> [t1, t4]
	p1 := periods.NewFinitePeriod(t1, t3, true, true)
	p2 := periods.NewFinitePeriod(t2, t4, true, true)
	res1 := p1.Union(p2)
	exp1 := periods.NewFinitePeriod(t1, t4, true, true)
	if !res1.Equals(exp1) {
		t.Logf("Union with overlap failed: got %s", res1.AsRawString())
		t.Fail()
	}

	// Contiguous union: [t1, t2] U ]t2, t3] -> [t1, t3]
	p3 := periods.NewFinitePeriod(t1, t2, true, true)
	p4 := periods.NewFinitePeriod(t2, t3, false, true)
	res2 := p3.Union(p4)
	exp2 := periods.NewFinitePeriod(t1, t3, true, true)
	if !res2.Equals(exp2) {
		t.Logf("Contiguous union failed: got %s", res2.AsRawString())
		t.Fail()
	}

	// Disjoint union: [t1, t2] U [t3, t4] -> keeps both disjoint intervals
	p5 := periods.NewFinitePeriod(t1, t2, true, true)
	p6 := periods.NewFinitePeriod(t3, t4, true, true)
	res3 := p5.Union(p6)
	if res3.IsEmpty() {
		t.Log("Disjoint union should not be empty")
		t.Fail()
	}
}

func TestPeriodContains(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	before := time.Now().Add(-24 * time.Hour).Truncate(time.Second)
	// period is [before, +oo[
	period := periods.NewPeriodSince(before, true)
	// before is in period
	if !period.Contains(before) {
		t.Log("Failed to test when value is boundary included")
		t.Fail()
	}

	// interval is [before, +oo[ and before < now, it should belong
	if !period.Contains(now) {
		t.Log("Failed to test outside value")
		t.Fail()
	}

	// period is ]before, +oo[
	period = periods.NewPeriodSince(before, false)
	// before is not in interval before it is excluded
	if period.Contains(before) {
		t.Log("Failed to test when value is boundary excluded")
		t.Fail()
	}
	// interval is ]before,+oo[ and before < now
	if !period.Contains(now) {
		t.Log("Failed to test outside value")
		t.Fail()
	}

	// period is ]-oo, now]
	period = periods.NewPeriodUntil(now, true)
	// before < now, so expecting period to contain in
	if !period.Contains(before) {
		t.Log("Failed to test when value is strictly included")
		t.Fail()
	}

	// now is included, so expecting to belong
	if !period.Contains(now) {
		t.Log("Failed to test bound value included")
		t.Fail()
	}

	// period is ]-oo, now[
	period = periods.NewPeriodUntil(now, false)
	// period should not contain now, because now is excluded
	if period.Contains(now) {
		t.Log("Failed to test bound value excluded")
		t.Fail()
	}
}

func TestPeriodIsIncludedIn(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	t1 := now.Add(1 * time.Hour)
	t2 := now.Add(2 * time.Hour)
	t3 := now.Add(3 * time.Hour)
	t4 := now.Add(4 * time.Hour)

	large := periods.NewFinitePeriod(t1, t4, true, true)
	small := periods.NewFinitePeriod(t2, t3, true, true)
	overlapping := periods.NewFinitePeriod(t3, now.Add(5*time.Hour), true, true)

	if !small.IsIncludedIn(large) {
		t.Log("Small should be included in large")
		t.Fail()
	}

	if large.IsIncludedIn(small) {
		t.Log("Large should not be included in small")
		t.Fail()
	}

	if overlapping.IsIncludedIn(large) {
		t.Log("Overlapping interval should not be completely included")
		t.Fail()
	}

	empty := periods.NewEmptyPeriod()
	if !empty.IsIncludedIn(large) {
		t.Log("Empty period should be mathematically included in any period")
		t.Fail()
	}

	full := periods.NewFullPeriod()
	if !large.IsIncludedIn(full) {
		t.Log("Any finite period should be included in full space")
		t.Fail()
	}
}

func TestPeriodRemove(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	before := now.AddDate(-10, 0, 0)
	after := now.AddDate(10, 0, 0)
	period := periods.NewPeriodUntil(after, true)
	toRemove := periods.NewFinitePeriod(before, now, true, false)
	remaining := period.Remove(toRemove)
	expected := periods.NewPeriodUntil(before, false).Union(periods.NewFinitePeriod(now, after, true, true))
	if !remaining.Equals(expected) {
		t.Logf("Failed to remove period, got %s but expected %s", remaining.AsRawString(), expected.AsRawString())
		t.Fail()
	}
}

func TestPeriodRemoveSplitting(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	t1 := now.Add(1 * time.Hour)
	t2 := now.Add(2 * time.Hour)
	t3 := now.Add(3 * time.Hour)
	t4 := now.Add(4 * time.Hour)

	base := periods.NewFinitePeriod(t1, t4, true, true)
	hole := periods.NewFinitePeriod(t2, t3, true, true)

	res := base.Remove(hole)
	// We expect [t1, t2[ U ]t3, t4]
	expLeft := periods.NewFinitePeriod(t1, t2, true, false)
	expRight := periods.NewFinitePeriod(t3, t4, false, true)
	expected := expLeft.Union(expRight)

	if !res.Equals(expected) {
		t.Logf("Failed to split interval via hole removal, got %s but expected %s", res.AsRawString(), expected.AsRawString())
		t.Fail()
	}
}

func TestPeriodRemoveDisjoint(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	t1 := now.Add(1 * time.Hour)
	t2 := now.Add(2 * time.Hour)
	t3 := now.Add(3 * time.Hour)
	t4 := now.Add(4 * time.Hour)

	base := periods.NewFinitePeriod(t1, t2, true, true)
	toRemove := periods.NewFinitePeriod(t3, t4, true, true)

	res := base.Remove(toRemove)
	if !res.Equals(base) {
		t.Logf("Failed to remove disjoint period, got %s but expected %s", res.AsRawString(), base.AsRawString())
		t.Fail()
	}
}

func TestPeriodRemoveLargerPeriod(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	before := now.AddDate(-10, 0, 0)
	after := now.AddDate(10, 0, 0)
	period := periods.NewFinitePeriod(before, now, true, true)
	toRemove := periods.NewPeriodUntil(after, true)
	remaining := period.Remove(toRemove)
	// remaining should be empty because [before, now] is in ]-oo, after]
	if !remaining.IsEmpty() {
		t.Log("remaining should be empty because toRemove contains each point of period")
		t.Fail()
	}

	toRemove = periods.NewFullPeriod()
	remaining = period.Remove(toRemove)
	// remaining should be empty because [before, now] is in the full period
	if !remaining.IsEmpty() {
		t.Log("remaining should be empty because toRemove contains period")
		t.Fail()
	}
}

func TestPeriodRemoveFromEmpty(t *testing.T) {
	p := periods.NewEmptyPeriod()
	res := p.Remove(periods.NewFullPeriod())
	if !res.Equals(p) {
		t.Logf("failed to remove from empty period, got %s", res.AsRawString())
		t.Fail()
	}
}

func TestFinitePeriodEdgeCases(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	before := now.Add(-1 * time.Hour)

	// Math edge case 1: min > max -> should be empty
	inverted := periods.NewFinitePeriod(now, before, true, true)
	if !inverted.IsEmpty() {
		t.Log("Inverted boundaries period (min > max) should be empty")
		t.Fail()
	}

	// Math edge case 2: min == max but boundaries excluded -> should be empty
	pointExcluded := periods.NewFinitePeriod(now, now, false, false)
	if !pointExcluded.IsEmpty() {
		t.Log("Point period with excluded boundaries should be empty")
		t.Fail()
	}

	// Math edge case 3: min == max and boundaries included -> valid point interval
	pointIncluded := periods.NewFinitePeriod(now, now, true, true)
	if pointIncluded.IsEmpty() || !pointIncluded.Contains(now) {
		t.Log("Point period with included boundaries should be valid and contain the exact point")
		t.Fail()
	}
}

func TestPeriodSerde(t *testing.T) {
	tested := periods.NewEmptyPeriod()
	if res, err := periods.PeriodLoad(tested.AsStrings()); err != nil {
		t.Log(err)
		t.Fail()
	} else if !res.Equals(tested) {
		t.Log("Failed to ser + deser empty")
		t.Fail()
	}

	tested = periods.NewFullPeriod()
	if res, err := periods.PeriodLoad(tested.AsStrings()); err != nil {
		t.Log(err)
		t.Fail()
	} else if !res.Equals(tested) {
		t.Log("Failed to ser + deser full")
		t.Fail()
	}

	now := time.Now().Truncate(time.Second)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(10, 0, 0)
	tested = periods.NewFinitePeriod(before, now, true, false)
	tested = tested.Union(periods.NewPeriodSince(after, true))
	if res, err := periods.PeriodLoad(tested.AsStrings()); err != nil {
		t.Log(err)
		t.Fail()
	} else if !res.Equals(tested) {
		t.Log("Failed to ser + deser union of intervals")
		t.Fail()
	}
}

func TestPeriodLoadErrors(t *testing.T) {
	badPartitions := [][]string{
		{"invalid-string"},
		{"[2024-01-01;]"}, // malformed parts
		{"-oo;+oo"},       // missing boundaries
		{"]foo;bar["},     // invalid dates
	}

	for _, parts := range badPartitions {
		if _, err := periods.PeriodLoad(parts); err == nil {
			t.Logf("Expected error for malformed partition %v, but got none", parts)
			t.Fail()
		}
	}
}

func TestPeriodInfiniteBoundaries(t *testing.T) {
	now := time.Now().Truncate(1 * time.Hour)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)
	a := periods.NewPeriodSince(after, true)
	b := periods.NewPeriodUntil(before, true)
	res := a.Union(b).BoundingPeriod()
	expected := periods.NewFullPeriod()
	if !expected.Equals(res) {
		t.Logf("failed to find full as boundaries, got %s", res.AsRawString())
		t.Fail()
	}
}

func TestPeriodFiniteBoundaries(t *testing.T) {
	now := time.Now().Truncate(1 * time.Hour)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)
	evenAfter := after.AddDate(10, 0, 0)
	a := periods.NewFinitePeriod(before, now, false, true)
	b := periods.NewFinitePeriod(after, evenAfter, true, true)
	res := a.Union(b).BoundingPeriod()
	expected := periods.NewFinitePeriod(before, evenAfter, false, true)
	if !expected.Equals(res) {
		t.Logf("failed to find finite intervals as boundaries, got %s", res.AsRawString())
		t.Fail()
	}
}

func TestPeriodUnionWhenIncludedSince(t *testing.T) {
	now := time.Now().Truncate(1 * time.Hour)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)
	container := periods.NewPeriodSince(before, true)
	contained := periods.NewPeriodSince(after, true)
	expected := periods.NewPeriodSince(before, true)

	res1 := container.Union(contained)
	res2 := contained.Union(container)
	if !res1.Equals(res2) {
		t.Logf("union is not commutative, got %s and %s", res1.AsRawString(), res2.AsRawString())
		t.Fail()
	} else if !res1.Equals(expected) {
		t.Logf("union of included periods should be as expected, got %s", res1.AsRawString())
		t.Fail()
	}
}

func TestPeriodUnionWhenIncludedUntil(t *testing.T) {
	now := time.Now().Truncate(1 * time.Hour)
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)
	container := periods.NewPeriodUntil(after, true)
	contained := periods.NewPeriodUntil(before, true)
	expected := periods.NewPeriodUntil(after, true)

	res1 := container.Union(contained)
	res2 := contained.Union(container)
	if !res1.Equals(res2) {
		t.Logf("union is not commutative, got %s and %s", res1.AsRawString(), res2.AsRawString())
		t.Fail()
	} else if !res1.Equals(expected) {
		t.Logf("union of included periods should be as expected, got %s", res1.AsRawString())
		t.Fail()
	}
}

func TestBoundingPeriodEdgeCases(t *testing.T) {
	// Edge case 1: Empty period
	empty := periods.NewEmptyPeriod()
	if !empty.BoundingPeriod().IsEmpty() {
		t.Log("BoundingPeriod of an empty period should be empty")
		t.Fail()
	}

	// Edge case 2: Single interval period
	now := time.Now().Truncate(time.Hour)
	before := now.Add(-1 * time.Hour)
	single := periods.NewFinitePeriod(before, now, true, false)

	res := single.BoundingPeriod()
	if !res.Equals(single) {
		t.Logf("BoundingPeriod of a single interval should be equal to itself, got %s", res.AsRawString())
		t.Fail()
	}
}
