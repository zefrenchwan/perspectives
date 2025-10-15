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

	empty := p.PickByIndexes([]int{})
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

	if value, found := p.Get(0); !found || value != tanguy {
		t.Fail()
	} else if value, found := p.Get(1); !found || value != alan {
		t.Fail()
	} else if value, found := p.Get(-1); found || value != nil {
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
	q.AppendAs("y", alan)

	variables := q.Names()
	slices.Sort(variables)
	if slices.Compare(variables, []string{"x", "y"}) != 0 {
		t.Fail()
	}

	if value, found := q.GetByName("x"); !found || value != tanguy {
		t.Fail()
	} else if value, found := q.GetByName("y"); !found || value != alan {
		t.Fail()
	} else if value, found := q.GetByName("z"); found || value != nil {
		t.Fail()
	} else if value, found := q.Get(0); found || value != nil {
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
	r.AppendAs("x", alan)

	if value, found := r.Get(0); !found || value != tanguy {
		t.Fail()
	} else if r.Size() != 1 {
		t.Fail()
	} else if slices.Compare([]string{"x"}, r.Names()) != 0 {
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

func TestContentSelect(t *testing.T) {
	camembert := DummyIdBasedImplementation{id: "camembert"}
	brie := DummyIdBasedImplementation{id: "brie"}

	variable := commons.NewNamedContent("x", &brie)
	variable.AppendAs("y", &camembert)

	// test empty
	result := variable.PickByIndexes([]int{0, 1})
	if !result.IsEmpty() {
		t.Fail()
	}

	result = variable.PickByNames([]string{"a", "b"})
	if !result.IsEmpty() {
		t.Fail()
	}

	// test select
	result = variable.PickByNames([]string{"x"})
	if value, found := result.GetByName("x"); !found || value != &brie {
		t.Fail()
	}

	// test select ints
	other := commons.NewContent(&camembert)
	other.Append(&brie)

	if !other.PickByIndexes(nil).IsEmpty() {
		t.Fail()
	} else if !other.PickByIndexes([]int{-1, 2, 3, 4}).IsEmpty() {
		t.Fail()
	}

	reduced := other.PickByIndexes([]int{0})
	if reduced.Size() != 1 {
		t.Fail()
	} else if value, found := reduced.Get(0); !found || value != &camembert {
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

func TestContentDisjoin(t *testing.T) {
	var a, b commons.Content
	a = commons.NewContent(DummyComponentImplementation{})
	if !a.Disjoin(b) {
		t.Fail()
	}

	b = commons.NewNamedContent("x", DummyComponentImplementation{})
	if !a.Disjoin(b) {
		t.Fail()
	} else if !b.Disjoin(a) {
		t.Fail()
	}

	a.AppendAs("y", DummyComponentImplementation{})
	if !a.Disjoin(b) {
		t.Fail()
	} else if !b.Disjoin(a) {
		t.Fail()
	}

	b.Append(DummyComponentImplementation{})
	if a.Disjoin(b) {
		t.Fail()
	} else if b.Disjoin(a) {
		t.Fail()
	}

	a = commons.NewNamedContent("x", DummyComponentImplementation{})
	b = commons.NewNamedContent("y", DummyComponentImplementation{})
	if !a.Disjoin(b) {
		t.Fail()
	} else if !b.Disjoin(a) {
		t.Fail()
	}

	a.AppendAs("y", DummyComponentImplementation{})
	if a.Disjoin(b) {
		t.Fail()
	} else if b.Disjoin(a) {
		t.Fail()
	}
}
