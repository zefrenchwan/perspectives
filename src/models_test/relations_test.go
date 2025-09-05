package models_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
	"github.com/zefrenchwan/perspectives.git/structures"
)

func TestRelationsBuild(t *testing.T) {
	john := models.NewObjectTerm(models.NewObject([]string{"human"}))
	cheese := models.NewObjectTerm(models.NewObject([]string{"cheese"}))
	loves := models.NewRelationTerm("loves", map[string]models.RelationTerm{"subject": john, "object": cheese}, structures.NewFullPeriod())

	if relation, err := loves.Build(); err != nil {
		t.Log(err)
		t.Fail()
	} else if relation == nil {
		t.Log("nil relation")
		t.Fail()
	} else if relation.Link != "loves" {
		t.Log("wrong verb")
		t.Fail()
	} else if !relation.Lifetime.Equals(structures.NewFullPeriod()) {
		t.Log("wrong period")
		t.Fail()
	} else if len(relation.Parameters) != 2 {
		t.Log("wrong roles")
		t.Fail()
	} else if v, found := relation.Parameters["subject"]; !found {
		t.Log("wrong subject")
		t.Fail()
	} else if objects := v.AsObjects(); len(objects) != 1 {
		t.Log("wrong subject slice")
		t.Fail()
	} else if slices.Compare([]string{"human"}, objects[0].DeclaringTraits()) != 0 {
		t.Log("wrong subject value")
		t.Fail()
	} else if v, found := relation.Parameters["object"]; !found {
		t.Log("wrong object")
		t.Fail()
	} else if objects := v.AsObjects(); len(objects) != 1 {
		t.Log("wrong object slice")
		t.Fail()
	} else if slices.Compare([]string{"cheese"}, objects[0].DeclaringTraits()) != 0 {
		t.Log("wrong object value")
		t.Fail()
	}
}

func TestRelationsComposeBuild(t *testing.T) {
	john := models.NewObjectTerm(models.NewObject([]string{"human", "man"}))
	cheese := models.NewObjectTerm(models.NewObject([]string{"cheese"}))
	loves := models.NewRelationTerm("loves", map[string]models.RelationTerm{"subject": john, "object": cheese}, structures.NewFullPeriod())
	marie := models.NewObjectTerm(models.NewObject([]string{"human", "woman"}))
	knows := models.NewRelationTerm("knows", map[string]models.RelationTerm{"subject": marie, "object": loves}, structures.NewFullPeriod())

	if relation, err := knows.Build(); err != nil {
		t.Log(err)
		t.Fail()
	} else if relation == nil {
		t.Log("nil relation")
		t.Fail()
	} else if relation.Link != "knows" {
		t.Log("wrong verb")
		t.Fail()
	} else if !relation.Lifetime.Equals(structures.NewFullPeriod()) {
		t.Log("wrong period")
		t.Fail()
	} else if len(relation.Parameters) != 2 {
		t.Log("wrong roles")
		t.Fail()
	} else if v, found := relation.Parameters["subject"]; !found {
		t.Log("wrong subject")
		t.Fail()
	} else if objects := v.AsObjects(); len(objects) != 1 {
		t.Log("wrong subject slice")
		t.Fail()
	} else if slices.Compare([]string{"human", "woman"}, objects[0].DeclaringTraits()) != 0 {
		t.Log("wrong subject value")
		t.Fail()
	} else if v, found := relation.Parameters["object"]; !found {
		t.Log("wrong object")
		t.Fail()
	} else if child, err := v.Build(); err != nil {
		t.Log("wrong composition")
		t.Fail()
	} else if child.Link != "loves" {
		t.Log("wrong verb for child")
		t.Fail()
	}
}

func TestRelationsWalkthrough(t *testing.T) {
	source := models.NewObject([]string{"Human"})
	dest := models.NewObject([]string{"Dessert"})
	likes := models.NewRelationTerm(
		"likes",
		map[string]models.RelationTerm{
			"subject": models.NewObjectTerm(source),
			"object":  models.NewObjectTerm(dest),
		},
		structures.NewFullPeriod(),
	)

	knows := models.NewRelationTerm(
		"knows",
		map[string]models.RelationTerm{
			"subject": models.NewObjectTerm(source),
			"object":  likes,
		},
		structures.NewFullPeriod(),
	)

	operands := knows.ObjectsOperands()
	values := make([]string, 0)
	for _, obj := range operands {
		values = append(values, obj.Id)
	}

	if len(values) != 2 {
		t.Log("missing operands")
		t.Fail()
	} else if !slices.Contains(values, source.Id) {
		t.Log("missing source")
		t.Fail()
	} else if !slices.Contains(values, dest.Id) {
		t.Log("missing dest")
		t.Fail()
	}

}
