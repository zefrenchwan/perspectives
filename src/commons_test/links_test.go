package commons_test

import (
	"slices"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestLinkUsage(t *testing.T) {
	a := DummyIdBasedImplementation{}
	if link, err := commons.NewLink("link", map[string]commons.Linkable{"role": a}); err != nil {
		t.Fail()
	} else if link.Name() != "link" {
		t.Fail()
	} else if link.IsEmpty() {
		t.Fail()
	} else if slices.Compare([]string{"role"}, link.Roles()) != 0 {
		t.Fail()
	} else if link.GetType() != commons.TypeLink {
		t.Fail()
	} else if ops := link.Operands(); len(ops) != 1 {
		t.Fail()
	} else if ops["role"] != a {
		t.Fail()
	} else if _, found := link.Get("no data"); found {
		t.Fail()
	} else if v, found := link.Get("role"); !found {
		t.Fail()
	} else if v != a {
		t.Fail()
	}

	// test errors when create
	noData := make(map[string]commons.Linkable)
	if _, err := commons.NewLink[commons.Linkable]("", nil); err == nil {
		t.Fail()
	} else if _, err := commons.NewLink[commons.Linkable]("a", nil); err == nil {
		t.Fail()
	} else if _, err := commons.NewLink("a", noData); err == nil {
		t.Fail()
	}
}

func TestTemporalLink(t *testing.T) {
	fullPeriod := commons.NewFullPeriod()
	partPeriod := commons.NewPeriodSince(time.Now().Truncate(time.Hour), true)

	if link, err := commons.NewLink("test", map[string]int{"role": 0}); err != nil {
		t.Log(err)
		t.Fail()
	} else if tlink := commons.NewTemporalLink(partPeriod, link); !tlink.ActivePeriod().Equals(partPeriod) {
		t.Fail()
	} else {
		tlink.SetActivePeriod(fullPeriod)
		if !tlink.ActivePeriod().Equals(fullPeriod) {
			t.Fail()
		}
	}
}
