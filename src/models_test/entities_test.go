package models

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/models"
)

func TestEntityLifetime(t *testing.T) {
	now := time.Now()
	before := now.AddDate(-10, 0, 0)
	base := models.NewEntity("id", before)

	expected := models.NewPeriodSince(now, true)
	if !base.LifetimeDuringPeriod(expected).Equals(expected) {
		t.Log("lifetime should be [before, +oo[ and then intersection with [now, +oo[ should be itself")
		t.Fail()
	}

	base.End(now)
	expected = models.NewFinitePeriod(before, now, true, false)
	if !base.LifetimeDuringPeriod(models.NewFullPeriod()).Equals(expected) {
		t.Log("entity should live during [before, now[")
		t.Fail()
	}
}
