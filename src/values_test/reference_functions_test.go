package values

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

func TestReferenceFunctions(t *testing.T) {
	referenceVal := "id to some ref"
	later := time.Now().Truncate(time.Second).AddDate(1, 0, 0)
	ref := values.NewReferenceTimeFunction()
	ref.Add("id to some ref", periods.NewFullPeriod())
	period := periods.NewPeriodSince(later, true)
	ref.Remove(period)

	if v, has := ref.Value(time.Now()); !has {
		t.Errorf("Expected value to be present")
	} else if v != referenceVal {
		t.Errorf("Expected value to be %s, got %s", referenceVal, v)
	}

	if _, has := ref.Value(later); has {
		t.Errorf("Expected value to be absent")
	}
}

func TestReferenceCopy(t *testing.T) {
	referenceVal := "id to some ref"
	later := time.Now().Truncate(time.Second).AddDate(1, 0, 0)
	ref := values.NewReferenceTimeFunction()
	ref.Add("id to some ref", periods.NewFullPeriod())

	refCopy := ref.Copy()

	// remove from ref
	period := periods.NewPeriodSince(later, true)
	ref.Remove(period)

	if v, has := refCopy.Value(time.Now()); !has {
		t.Errorf("Expected value to be present")
	} else if v != referenceVal {
		t.Errorf("Expected value to be %s, got %s", referenceVal, v)
	} else if otherVal, has := refCopy.Value(later); !has {
		t.Errorf("Expected value to be present ")
	} else if otherVal != referenceVal {
		t.Errorf("Expected value to be %s, got %s", referenceVal, otherVal)
	}

	// test equality
	if !ref.Equals(ref) {
		t.Errorf("ref should equal itself")
	} else if ref.Equals(refCopy) {
		t.Errorf("ref should not equal copy (they differ)")
	}
}
