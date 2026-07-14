package values_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

func TestIntPrimitiveFunction(t *testing.T) {
	function := values.NewIntTimeFunction()
	function.Add(10, periods.NewFullPeriod())
	if v, has := function.Value(time.Now()); !has {
		t.Errorf("Expected value to be present")
	} else if v != 10 {
		t.Errorf("Expected value to be 10")
	}

}
