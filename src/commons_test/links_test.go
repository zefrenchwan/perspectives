package commons_test

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestLinkableSame(t *testing.T) {
	labelA := commons.NewLabel("a")
	labelB := commons.NewLabel("b")
	objectA := commons.NewModelObject()
	objectB := commons.NewModelObject()
	varA := commons.NewLinkVariable("a", func(l commons.Linkable) bool { return true })
	varB := commons.NewLinkVariable("b", func(l commons.Linkable) bool { return true })

	if commons.LinkableSame(labelA, labelB) {
		t.Fail()
	} else if commons.LinkableSame(labelA, objectA) {
		t.Fail()
	} else if !commons.LinkableSame(labelA, labelA) {
		t.Fail()
	} else if commons.LinkableSame(objectA, objectB) {
		t.Fail()
	} else if !commons.LinkableSame(objectA, objectA) {
		t.Fail()
	} else if commons.LinkableSame(varA, varB) {
		t.Fail()
	} else if !commons.LinkableSame(varA, varA) {
		t.Fail()
	} else if commons.LinkableSame(varA, objectA) {
		t.Fail()
	}
}

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

func TestLinkAcceptsInstantiationSameLink(t *testing.T) {
	link, _ := commons.NewSimpleLink("extends", commons.NewLabel("dessert"), commons.NewLabel("food"))
	otherLink, _ := commons.NewSimpleLink("extends", commons.NewLabel("cakes"), commons.NewLabel("food"))

	if _, ok := commons.LinkAcceptsInstantiation(link, link, commons.LinkableSame); !ok {
		t.Log("same link")
		t.Fail()
	}

	if _, ok := commons.LinkAcceptsInstantiation(link, otherLink, commons.LinkableSame); ok {
		t.Log("different structure")
		t.Fail()
	}
}

func TestLinkAcceptsInstantiationRoot(t *testing.T) {
	linkAcceptor := func(l commons.Linkable) bool {
		if l == nil {
			return false
		}

		link, ok := l.(commons.Link)
		return link != nil && ok
	}

	labelAcceptor := func(l commons.Linkable) bool {
		if l == nil {
			return false
		}

		_, ok := l.(commons.LinkLabel)
		return ok
	}

	link, _ := commons.NewSimpleLink("extends", commons.NewLabel("dessert"), commons.NewLabel("food"))
	other := commons.NewLinkVariable("x", linkAcceptor)

	if assignations, ok := commons.LinkAcceptsInstantiation(link, other, commons.LinkableSame); !ok {
		t.Log("match")
		t.Fail()
	} else if len(assignations) != 1 {
		t.Fail()
	} else if value, found := assignations["x"]; !found {
		t.Fail()
	} else if l, ok := value.(commons.Link); !ok || l == nil {
		t.Fail()
	} else if l.Id() != link.Id() {
		t.Fail()
	}

	otherLabel := commons.NewLinkVariable("x", labelAcceptor)
	if _, ok := commons.LinkAcceptsInstantiation(link, otherLabel, commons.LinkableSame); ok {
		t.Log("no match")
		t.Fail()
	}
}

func TestLinkAcceptsInstantiationStructure(t *testing.T) {
	linkAcceptor := func(l commons.Linkable) bool {
		if l == nil {
			return false
		}

		link, ok := l.(commons.Link)
		return link != nil && ok
	}

	labelAcceptor := func(l commons.Linkable) bool {
		if l == nil {
			return false
		}

		_, ok := l.(commons.LinkLabel)
		return ok
	}

	x := commons.NewLinkVariable("x", linkAcceptor)
	y := commons.NewLinkVariable("y", labelAcceptor)

	kate := commons.NewModelObject()
	chocolate := commons.NewLabel("chocolate")
	baseLink, _ := commons.NewSimpleLink("likes", kate, chocolate)
	differentLink, _ := commons.NewSimpleLink("loves", kate, chocolate)
	patternMatching, _ := commons.NewSimpleLink("likes", kate, y)
	patternNOMatch, _ := commons.NewSimpleLink("likes", kate, x)

	if matching, ok := commons.LinkAcceptsInstantiation(baseLink, patternMatching, commons.LinkableSame); !ok {
		t.Log("y => chocolate is a match")
		t.Fail()
	} else if len(matching) != 1 {
		t.Fail()
	} else if m, found := matching["y"]; !found {
		t.Fail()
	} else if label, ok := m.(commons.LinkLabel); !ok {
		t.Fail()
	} else if label.Name() != chocolate.Name() {
		t.Fail()
	}

	if _, ok := commons.LinkAcceptsInstantiation(differentLink, patternMatching, commons.LinkableSame); ok {
		t.Log("different link")
		t.Fail()
	}

	if _, ok := commons.LinkAcceptsInstantiation(baseLink, patternNOMatch, commons.LinkableSame); ok {
		t.Log("x accepts links only")
		t.Fail()
	}
}

func TestLinkAcceptsInstantiationMultipleMatches(t *testing.T) {

	labelAcceptor := func(l commons.Linkable) bool {
		if l == nil {
			return false
		}

		_, ok := l.(commons.LinkLabel)
		return ok
	}

	x := commons.NewLinkVariable("x", labelAcceptor)
	y := commons.NewLinkVariable("y", labelAcceptor)

	women := commons.NewLabel("women")
	humans := commons.NewLabel("humans")
	baseLink, _ := commons.NewSimpleLink("extends", women, humans)
	patternMatching, _ := commons.NewSimpleLink("extends", x, y)
	patternNOMatch, _ := commons.NewSimpleLink("extends", x, x)

	if matching, ok := commons.LinkAcceptsInstantiation(baseLink, patternMatching, commons.LinkableSame); !ok {
		t.Log("y => chocolate is a match")
		t.Fail()
	} else if len(matching) != 2 {
		t.Fail()
	} else if mx, found := matching["x"]; !found {
		t.Fail()
	} else if lx, ok := mx.(commons.LinkLabel); !ok {
		t.Fail()
	} else if lx.Name() != women.Name() {
		t.Fail()
	} else if my, found := matching["y"]; !found {
		t.Fail()
	} else if ly, ok := my.(commons.LinkLabel); !ok {
		t.Fail()
	} else if ly.Name() != humans.Name() {
		t.Fail()
	}

	if _, ok := commons.LinkAcceptsInstantiation(baseLink, patternNOMatch, commons.LinkableSame); ok {
		t.Log("x assigned with different values")
		t.Fail()
	}
}

func TestLinkFindAll(t *testing.T) {
	books, _ := commons.NewSimpleLink("is", commons.NewLabel("reading"), commons.NewLabel("awesome"))
	hugh := commons.NewModelObject()
	timothy := commons.NewModelObject()
	thinks, _ := commons.NewSimpleLink("thinks", timothy, books)
	likes, _ := commons.NewSimpleLink("likes", hugh, thinks)

	findLinks := func(l commons.Linkable) bool {
		if l == nil {
			return false
		}

		link, ok := l.(commons.Link)

		return ok && link != nil
	}

	findAwesome := func(l commons.Linkable) bool {
		if l == nil {
			return false
		}

		label, ok := l.(commons.LinkLabel)

		return ok && label.Name() == "awesome"
	}

	if result := commons.LinkFindAllMatching(thinks, findLinks); len(result) != 2 {
		t.Fail()
	} else if l, ok := result[0].(commons.Link); !ok || l == nil {
		t.Fail()
	} else if l.Id() != thinks.Id() {
		t.Fail()
	} else if l, ok := result[1].(commons.Link); !ok || l == nil {
		t.Fail()
	} else if l.Id() != books.Id() {
		t.Fail()
	}

	if result := commons.LinkFindAllMatching(likes, findAwesome); len(result) != 1 {
		t.Fail()
	} else if l, ok := result[0].(commons.LinkLabel); !ok {
		t.Fail()
	} else if l.Name() != "awesome" {
		t.Fail()
	}
}

func TestLinkUseVariables(t *testing.T) {
	dad := commons.NewModelObject()
	kid := commons.NewModelObject()
	variable := commons.NewLinkVariableForObject("x")
	parent, _ := commons.NewSimpleLink("parent", dad, kid)
	generic, _ := commons.NewSimpleLink("parent", dad, variable)

	if inst, matches := commons.LinkAcceptsInstantiation(parent, generic, commons.LinkableSame); !matches {
		t.Fail()
	} else if len(inst) != 1 {
		t.Fail()
	} else if value, found := inst["x"]; !found {
		t.Fail()
	} else if k, ok := value.(commons.ModelObject); !ok {
		t.Fail()
	} else if k.Id() != kid.Id() {
		t.Fail()
	} else if result, err := commons.LinkSetVariables(generic, inst); err != nil {
		t.Log(err)
		t.Fail()
	} else if result.Name() != parent.Name() {
		t.Fail()
	} else if ops := result.Operands(); len(ops) != 2 {
		t.Fail()
	} else if s, found := ops[commons.RoleSubject]; !found {
		t.Fail()
	} else if subject, ok := s.(commons.ModelObject); !ok {
		t.Fail()
	} else if subject.Id() != dad.Id() {
		t.Fail()
	} else if o, found := ops[commons.RoleObject]; !found {
		t.Fail()
	} else if object, ok := o.(commons.ModelObject); !ok {
		t.Fail()
	} else if object.Id() != kid.Id() {
		t.Fail()
	}

	// test no change when no variable
	if result, err := commons.LinkSetVariables(parent, map[string]commons.Linkable{"x": kid}); err != nil {
		t.Log(err)
		t.Fail()
	} else if result.Id() != parent.Id() {
		t.Fail()
	} else if result.Name() != parent.Name() {
		t.Fail()
	} else if ops := result.Operands(); len(ops) != 2 {
		t.Fail()
	} else if s, found := ops[commons.RoleSubject]; !found {
		t.Fail()
	} else if subject, ok := s.(commons.ModelObject); !ok {
		t.Fail()
	} else if subject.Id() != dad.Id() {
		t.Fail()
	} else if o, found := ops[commons.RoleObject]; !found {
		t.Fail()
	} else if object, ok := o.(commons.ModelObject); !ok {
		t.Fail()
	} else if object.Id() != kid.Id() {
		t.Fail()
	}
}
