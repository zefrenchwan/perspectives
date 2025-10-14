package commons_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
)

func TestContentCreation(t *testing.T) {
	durian := DummyIdBasedImplementation{id: "durian"}
	p := commons.NewNamedContent("x", durian)

	if p.IsEmpty() {
		t.Fail()
	}

	empty := p.Select([]int{})
	if !empty.IsEmpty() {
		t.Fail()
	}
}

func TestContentGet(t *testing.T) {
	var tanguy, alan commons.ModelComponent
	tanguy = DummyIdBasedImplementation{id: "tanguy"}
	alan = DummyIdBasedImplementation{id: "alan"}

	p := commons.NewContent(tanguy)
	p.Append(alan)

	if p.Get(0) != tanguy {
		t.Fail()
	} else if p.Get(1) != alan {
		t.Fail()
	} else if p.Get(-1) != nil {
		t.Fail()
	}

	if value := p.PositionalContent(); len(value) != 2 {
		t.Fail()
	} else if !slices.Contains(value, tanguy) {
		t.Fail()
	} else if !slices.Contains(value, alan) {
		t.Fail()
	} else if named := p.NamedContent(); len(named) != 0 {
		t.Fail()
	}

	q := commons.NewNamedContent("x", tanguy)
	q.AppendAsVariable("y", alan)

	variables := q.Variables()
	slices.Sort(variables)
	if slices.Compare(variables, []string{"x", "y"}) != 0 {
		t.Fail()
	}

	if q.GetVariable("x") != tanguy {
		t.Fail()
	} else if q.GetVariable("y") != alan {
		t.Fail()
	} else if q.GetVariable("z") != nil {
		t.Fail()
	} else if q.Get(0) != nil {
		t.Fail()
	}

	if values := q.NamedContent(); len(values) != 2 {
		t.Fail()
	} else if v := values["x"]; v != tanguy {
		t.Fail()
	} else if v := values["y"]; v != alan {
		t.Fail()
	} else if len(q.PositionalContent()) != 0 {
		t.Fail()
	}

	r := commons.NewContent(tanguy)
	r.AppendAsVariable("x", alan)

	if r.Get(0) != tanguy {
		t.Fail()
	} else if r.Size() != 1 {
		t.Fail()
	} else if slices.Compare([]string{"x"}, r.Variables()) != 0 {
		t.Fail()
	} else if names := r.NamedContent(); len(names) != 1 {
		t.Fail()
	} else if names["x"] != alan {
		t.Fail()
	} else if positionals := r.PositionalContent(); len(positionals) != 1 {
		t.Fail()
	} else if positionals[0] != tanguy {
		t.Fail()
	}
}

func TestAddExtraContent(t *testing.T) {
	// test separated variables
	base := commons.NewNamedContent("x", DummyIdBasedImplementation{id: "base"})
	other := commons.NewNamedContent("y", DummyIdBasedImplementation{id: "other"})
	if result, extra := base.AddExtraContent(other); !extra {
		t.Log("should have added value")
		t.Fail()
	} else if len(result.PositionalContent()) != 0 {
		t.Log("no more positional content")
		t.Fail()
	} else if values := base.NamedContent(); len(values) != 2 {
		t.Fail()
	} else if x := values["x"]; x == nil {
		t.Fail()
	} else if y := values["y"]; y == nil {
		t.Fail()
	}

	// test common variables mean no change
	base = commons.NewNamedContent("x", DummyIdBasedImplementation{id: "base"})
	other = commons.NewNamedContent("x", DummyIdBasedImplementation{id: "other"})
	if result, extra := base.AddExtraContent(other); extra {
		t.Fail()
	} else if len(result.PositionalContent()) != 0 {
		t.Fail()
	} else if variables := result.NamedContent(); len(variables) != 1 {
		t.Fail()
	} else if value, found := variables["x"]; !found || value == nil {
		t.Fail()
	} else if v, ok := value.(commons.Identifiable); !ok {
		t.Fail()
	} else if v.Id() != "base" {
		t.Fail()
	}

	// add extra positionals from scratch
	base = commons.NewNamedContent("x", DummyIdBasedImplementation{id: "base"})
	other = commons.NewContent(DummyIdBasedImplementation{id: "other"})
	if result, extra := base.AddExtraContent(other); !extra {
		t.Fail()
	} else if positionals := result.PositionalContent(); len(positionals) != 1 {
		t.Fail()
	} else if p := positionals[0]; p == nil {
		t.Fail()
	} else if o, ok := p.(commons.Identifiable); !ok {
		t.Fail()
	} else if o.Id() != "other" {
		t.Fail()
	} else if variables := result.NamedContent(); len(variables) != 1 {
		t.Fail()
	} else if value, found := variables["x"]; !found || value == nil {
		t.Fail()
	} else if v, ok := value.(commons.Identifiable); !ok {
		t.Fail()
	} else if v.Id() != "base" {
		t.Fail()
	}

	// add extra positionals from existing
	base = commons.NewContent(DummyIdBasedImplementation{id: "base"})
	other = commons.NewContent(DummyIdBasedImplementation{id: "other"})
	other.Append(DummyIdBasedImplementation{id: "new"})

	if result, extra := base.AddExtraContent(other); !extra {
		t.Fail()
	} else if positionals := result.PositionalContent(); len(positionals) != 2 {
		t.Fail()
	} else if p := positionals[0]; p == nil {
		t.Fail()
	} else if o, ok := p.(commons.Identifiable); !ok {
		t.Fail()
	} else if o.Id() != "base" {
		t.Log("existing value was erased")
		t.Fail()
	} else if p := positionals[1]; p == nil {
		t.Fail()
	} else if o, ok := p.(commons.Identifiable); !ok {
		t.Fail()
	} else if o.Id() != "new" {
		t.Fail()
	}

	// test larger positions should not change values
	base = commons.NewContent(DummyIdBasedImplementation{id: "base"})
	base.Append(DummyIdBasedImplementation{id: "other"})
	other = commons.NewContent(DummyIdBasedImplementation{id: "none"})

	if result, extra := base.AddExtraContent(other); extra {
		t.Fail()
	} else if positionals := result.PositionalContent(); len(positionals) != 2 {
		t.Fail()
	} else if p := positionals[0]; p == nil {
		t.Fail()
	} else if o, ok := p.(commons.Identifiable); !ok {
		t.Fail()
	} else if o.Id() != "base" {
		t.Fail()
	} else if p := positionals[1]; p == nil {
		t.Fail()
	} else if o, ok := p.(commons.Identifiable); !ok {
		t.Fail()
	} else if o.Id() != "other" {
		t.Fail()
	}

}

func TestContentSelect(t *testing.T) {
	camembert := DummyIdBasedImplementation{id: "camembert"}
	brie := DummyIdBasedImplementation{id: "brie"}

	variable := commons.NewNamedContent("x", &brie)
	variable.AppendAsVariable("y", &camembert)

	// test empty
	result := variable.Select([]int{0, 1})
	if !result.IsEmpty() {
		t.Fail()
	}

	result = variable.SelectVariables([]string{"a", "b"})
	if !result.IsEmpty() {
		t.Fail()
	}

	// test select
	result = variable.SelectVariables([]string{"x"})
	if value := result.GetVariable("x"); value != &brie {
		t.Fail()
	}

	// test select ints
	other := commons.NewContent(&camembert)
	other.Append(&brie)

	if !other.Select(nil).IsEmpty() {
		t.Fail()
	} else if !other.Select([]int{-1, 2, 3, 4}).IsEmpty() {
		t.Fail()
	}

	reduced := other.Select([]int{0})
	if reduced.Size() != 1 {
		t.Fail()
	} else if reduced.Get(0) != &camembert {
		t.Fail()
	}

}

func TestContentUnique(t *testing.T) {
	leila := DummyIdBasedImplementation{id: "leila"}
	maria := DummyIdBasedImplementation{id: "maria"}

	p := commons.NewNamedContent("x", leila)
	if res, matching := p.Unique(); !matching {
		t.Fail()
	} else if res != leila {
		t.Fail()
	}

	p.Append(maria)
	if res, matching := p.Unique(); matching {
		t.Fail()
	} else if res != nil {
		t.Fail()
	}

	p = commons.NewContent(maria)
	if res, matching := p.Unique(); !matching {
		t.Fail()
	} else if res != maria {
		t.Fail()
	}

}
