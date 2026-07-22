package values_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

// TestMappingLocalConcept tests the concept of local mapping : easy to read ?
func TestMappingLocalConcept(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	beforePeriod := periods.NewPeriodUntil(now, true)
	afterPeriod := periods.NewPeriodSince(now, false)
	result := values.NewStringLocalMapping(map[string]periods.Period{
		"before": beforePeriod,
		"after":  afterPeriod,
	})

	for period, value := range result.Range() {
		if value.Content() == "before" {
			if !period.Equals(beforePeriod) {
				t.Errorf("BFORE : Expected period to be %v, but got %v", beforePeriod, period)
			}
		} else if value.Content() == "after" {
			if !period.Equals(afterPeriod) {
				t.Errorf("AFTER : Expected period to be %v, but got %v", afterPeriod, period)
			}
		} else {
			t.Errorf("Unexpected value: %v", value)
		}
	}
}

func TestMappingLocalHash(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	beforePeriod := periods.NewPeriodUntil(now, true)
	afterPeriod := periods.NewPeriodSince(now, false)

	resultBeforeAfter := values.NewStringLocalMapping(map[string]periods.Period{
		"before": beforePeriod,
		"after":  afterPeriod,
	})

	resultAfter := values.NewStringLocalMapping(map[string]periods.Period{
		"after": afterPeriod,
	})

	resultBefore := values.NewStringLocalMapping(map[string]periods.Period{
		"before": beforePeriod,
	})

	resultAfterBefore := values.NewStringLocalMapping(map[string]periods.Period{
		"after":  afterPeriod,
		"before": beforePeriod,
	})

	if resultBeforeAfter.ToHashString() != resultAfterBefore.ToHashString() {
		t.Errorf("Expected hash of resultBeforeAfter to be equal to resultAfterBefore")
	} else if resultAfter.ToHashString() == resultBefore.ToHashString() {
		t.Errorf("Expected hash of resultAfter to be different to resultBefore")
	} else if resultAfter.ToHashString() == resultAfterBefore.ToHashString() {
		t.Errorf("Expected hash of resultAfter to be different to resultAfterBefore")
	}
}
