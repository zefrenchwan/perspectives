package commons_test

import (
	"slices"
	"strings"
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
	if _, err := commons.NewLink("", nil); err == nil {
		t.Fail()
	} else if _, err := commons.NewLink("a", nil); err == nil {
		t.Fail()
	} else if _, err := commons.NewLink("a", noData); err == nil {
		t.Fail()
	}
}

func TestTemporalLink(t *testing.T) {
	fullPeriod := commons.NewFullPeriod()
	partPeriod := commons.NewPeriodSince(time.Now().Truncate(time.Hour), true)

	if link, err := commons.NewLink("test", map[string]commons.Linkable{"role": DummyObject{}}); err != nil {
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

	if link, err := commons.NewLink("test", map[string]commons.Linkable{"role": DummyIdBasedImplementation{id: "test"}}); err != nil {
		t.Log(err)
		t.Fail()
	} else if tlink := commons.NewTemporalLink(partPeriod, link); !tlink.ActivePeriod().Equals(partPeriod) {
		t.Fail()
	} else if link.Name() != tlink.Name() {
		t.Fail()
	} else if tlink.Id() == link.Id() {
		t.Fail()
	} else if len(tlink.Operands()) != len(link.Operands()) {
		t.Fail()
	} else {
		other := tlink.Operands()
		for k, v := range link.Operands() {
			if other[k] != v {
				t.Fail()
			}
		}
	}

}

func TestLinkComposition(t *testing.T) {
	marie := DummyObject{id: "marie"}
	paul := DummyObject{id: "paul"}

	if marie == paul {
		t.Fail()
	}

	link, errLink := commons.NewLink("loves", map[string]commons.Linkable{"subject": marie, "object": paul})
	if errLink != nil {
		t.Fail()
	}

	knows, errKnows := commons.NewSimpleLink("knows", marie, link)
	if errKnows != nil {
		t.Fail()
	}

	if opKnows := knows.Operands(); opKnows["subject"] != marie {
		t.Fail()
	} else if value := opKnows["object"]; value == nil {
		t.Fail()
	} else if l := value.(commons.Link); l == nil {
		t.Fail()
	} else if l.Name() != link.Name() {
		t.Fail()
	} else if l.Id() != link.Id() {
		t.Fail()
	} else if opLink := l.Operands(); opLink["subject"] != marie {
		t.Fail()
	} else if opLink["object"] != paul {
		t.Fail()
	}
}

func TestLeafMappingId(t *testing.T) {
	link, _ := commons.NewSimpleLink("extends", commons.NewLabel("dessert"), commons.NewLabel("food"))
	if result, err := commons.LinkMapLeafs(link, func(l commons.Linkable) (commons.Linkable, bool) { return l, false }); err != nil {
		t.Log(err)
		t.Fail()
	} else if result.Id() != link.Id() {
		t.Fail()
	} else if ops := result.Operands(); len(ops) != 2 {
		t.Fail()
	} else if s, found := result.Get(commons.RoleSubject); !found {
		t.Fail()
	} else if l, ok := s.(commons.LinkLabel); !ok {
		t.Fail()
	} else if l.Name() != "dessert" {
		t.Fail()
	} else if o, found := result.Get(commons.RoleObject); !found {
		t.Fail()
	} else if l, ok := o.(commons.LinkLabel); !ok {
		t.Fail()
	} else if l.Name() != "food" {
		t.Fail()
	}
}

func TestLeafMappingIdValuesReuse(t *testing.T) {
	value := commons.NewModelObject()
	link, _ := commons.NewSimpleLink("is", value, value)
	if result, err := commons.LinkMapLeafs(link, func(l commons.Linkable) (commons.Linkable, bool) { return l, false }); err != nil {
		t.Log(err)
		t.Fail()
	} else if result.Id() != link.Id() {
		t.Fail()
	} else if ops := result.Operands(); len(ops) != 2 {
		t.Fail()
	} else if s, found := result.Get(commons.RoleSubject); !found {
		t.Fail()
	} else if s != value {
		t.Fail()
	} else if o, found := result.Get(commons.RoleObject); !found {
		t.Fail()
	} else if o != value {
		t.Fail()
	}
}

func TestLeafMappingChangedLeaf(t *testing.T) {
	link, _ := commons.NewSimpleLink("extends", commons.NewLabel("dessert"), commons.NewLabel("food"))
	mapper := func(l commons.Linkable) (commons.Linkable, bool) {
		if s, ok := l.(commons.LinkLabel); !ok {
			return l, false
		} else {
			return commons.NewLabel(strings.ToUpper(s.Name())), true
		}
	}

	if result, err := commons.LinkMapLeafs(link, mapper); err != nil {
		t.Log(err)
		t.Fail()
	} else if result.Id() == link.Id() {
		t.Log("should have changed id")
		t.Fail()
	} else if ops := result.Operands(); len(ops) != 2 {
		t.Fail()
	} else if s, found := result.Get(commons.RoleSubject); !found {
		t.Fail()
	} else if l, ok := s.(commons.LinkLabel); !ok {
		t.Fail()
	} else if l.Name() != "DESSERT" {
		t.Fail()
	} else if o, found := result.Get(commons.RoleObject); !found {
		t.Fail()
	} else if l, ok := o.(commons.LinkLabel); !ok {
		t.Fail()
	} else if l.Name() != "FOOD" {
		t.Fail()
	}
}

func TestLeafMappingUnbalanced(t *testing.T) {
	tiramisu, _ := commons.NewSimpleLink("is", commons.NewLabel("tiramisu"), commons.NewLabel("amazing"))
	marie := commons.NewModelObject()
	julie := commons.NewModelObject()
	thinks, _ := commons.NewSimpleLink("thinks", marie, tiramisu)
	knows, _ := commons.NewSimpleLink("knows", julie, thinks)

	mapper := func(l commons.Linkable) (commons.Linkable, bool) {
		if s, ok := l.(commons.LinkLabel); !ok {
			return l, false
		} else {
			return commons.NewLabel(strings.ToUpper(s.Name())), true
		}
	}

	if result, err := commons.LinkMapLeafs(knows, mapper); err != nil {
		t.Log(err)
		t.Fail()
	} else if result.Id() == knows.Id() {
		t.Log("should have changed id")
		t.Fail()
	} else if knowsOps := result.Operands(); len(knowsOps) != 2 {
		t.Fail()
	} else if julieOp := knowsOps[commons.RoleSubject]; julieOp != julie {
		t.Fail()
	} else if thinksOp, found := knowsOps[commons.RoleObject]; !found {
		t.Fail()
	} else if thinksLink, ok := thinksOp.(commons.Link); thinksLink == nil || !ok {
		t.Fail()
	} else if thinksLink.Name() != "thinks" {
		t.Fail()
	} else if thinkOps := thinksLink.Operands(); len(thinkOps) != 2 {
		t.Fail()
	} else if thinkOps[commons.RoleSubject] != marie {
		t.Fail()
	} else if tiramisuOp, found := thinkOps[commons.RoleObject]; !found {
		t.Fail()
	} else if tiramisuLink, ok := tiramisuOp.(commons.Link); !ok || tiramisuLink == nil {
		t.Fail()
	} else if tiramisuOps := tiramisuLink.Operands(); len(tiramisuOps) != 2 {
		t.Fail()
	} else if s, ok := tiramisuOps[commons.RoleSubject].(commons.LinkLabel); !ok {
		t.Fail()
	} else if s.Name() != "TIRAMISU" {
		t.Fail()
	} else if o, ok := tiramisuOps[commons.RoleObject].(commons.LinkLabel); !ok {
		t.Fail()
	} else if o.Name() != "AMAZING" {
		t.Fail()
	}
}
