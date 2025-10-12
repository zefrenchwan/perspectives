package models_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/models"
)

func TestParameterCreation(t *testing.T) {
	durian := models.NewObject([]string{"Fruit"})
	p := commons.NewNamedParameter("x", durian)

	if p.IsEmpty() {
		t.Fail()
	}

	empty := p.Select([]int{})
	if !empty.IsEmpty() {
		t.Fail()
	}
}

func TestParametersGet(t *testing.T) {
	var tanguy, alan commons.ModelElement
	tanguy = models.NewObject([]string{"Human"})
	alan = models.NewObject([]string{"Human"})

	p := commons.NewParameter(tanguy)
	p.Append(alan)

	if p.Get(0) != tanguy {
		t.Fail()
	} else if p.Get(1) != alan {
		t.Fail()
	} else if p.Get(-1) != nil {
		t.Fail()
	}

	if value := p.PositionalParameters(); len(value) != 2 {
		t.Fail()
	} else if !slices.Contains(value, tanguy) {
		t.Fail()
	} else if !slices.Contains(value, alan) {
		t.Fail()
	} else if named := p.NamedParameters(); len(named) != 0 {
		t.Fail()
	}

	q := commons.NewNamedParameter("x", tanguy)
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

	if values := q.NamedParameters(); len(values) != 2 {
		t.Fail()
	} else if v := values["x"]; v != tanguy {
		t.Fail()
	} else if v := values["y"]; v != alan {
		t.Fail()
	} else if len(q.PositionalParameters()) != 0 {
		t.Fail()
	}

	r := commons.NewParameter(tanguy)
	r.AppendAsVariable("x", alan)

	if r.Get(0) != tanguy {
		t.Fail()
	} else if r.Size() != 1 {
		t.Fail()
	} else if slices.Compare([]string{"x"}, r.Variables()) != 0 {
		t.Fail()
	} else if names := r.NamedParameters(); len(names) != 1 {
		t.Fail()
	} else if names["x"] != alan {
		t.Fail()
	} else if positionals := r.PositionalParameters(); len(positionals) != 1 {
		t.Fail()
	} else if positionals[0] != tanguy {
		t.Fail()
	}
}

func TestParameterSelect(t *testing.T) {
	camembert := models.NewObject([]string{"Cheese"})
	brie := models.NewObject([]string{"Cheese"})

	variable := commons.NewNamedParameter("x", &brie)
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
	other := commons.NewParameter(&camembert)
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

func TestParametersUnique(t *testing.T) {
	leila := models.NewObject([]string{"Human"})
	maria := models.NewObject([]string{"Human"})

	p := commons.NewNamedParameter("x", leila)
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

	p = commons.NewParameter(maria)
	if res, matching := p.Unique(); !matching {
		t.Fail()
	} else if res != maria {
		t.Fail()
	}

}
