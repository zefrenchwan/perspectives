package entities_test


import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/entities"
	"github.com/zefrenchwan/perspectives.git/periods"
)


// ============================================
// PUT IN HERE ALL END TO END / DESIGN TESTS ==
// ============================================

func TestEntityBuild(t *testing.T) {
	john, _ := entities.NewLocalEntityBuilder("john").
		WithActivity(periods.NewFullPeriod()).
		Build()
	tiramisu, _ := entities.NewLocalEntityBuilder("tiramisu").
		WithActivity(periods.NewFullPeriod()).
		WithAttributeDuring("calories", periods.NewFullPeriod(), "way too much").
		Build()
	likes, _ := entities.NewLocalEntityBuilder("likes").
		WithActivity(periods.NewFullPeriod()).
		WithOperand("subject", john).
		WithOperand("object", tiramisu).
		Build()

	if subject, hasSubject := likes.Role("subject"); subject == nil || !hasSubject {
		t.Error("Expected subject to be john")
	} else if subject != john {
		t.Error("Expected subject to be john")
	}

	if object, hasObject := likes.Role("object"); object == nil || !hasObject {
		t.Error("Expected object to be tiramisu")
	} else if object != tiramisu {
		t.Error("Expected object to be tiramisu")
	} else if description, has := object.Attribute("calories"); !has {
		t.Error("Expected object to have attribute calories")
	} else if description.AttributeName != "calories" {
		t.Error("Expected attribute name to be calories")
	} else if description.AttributeType != "string" {
		t.Error("Expected attribute type to be string")
	} else if value, hasValue := object.ValueAt("calories", time.Now()); !hasValue {
		t.Error("Expected object to have value for attribute calories")
	} else if value != "way too much" {
		t.Error("Expected value to be way too much")
	}
}
