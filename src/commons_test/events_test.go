package commons_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestEventTick(t *testing.T) {
	now := time.Now().Truncate(commons.TIME_PRECISION)
	tick := commons.NewEventTickTime(now)
	if !tick.ProcessingTime().Equal(now) {
		t.Fail()
	}
}
