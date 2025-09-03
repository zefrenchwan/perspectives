package models_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
)

func TestClassBuild(t *testing.T) {
	c := models.NewFormalClass("humans")
	c.SetAttribute("account", []string{"social network"})

	if attr, found := c.Attributes["account"]; !found {
		t.Log("missing attribute")
		t.Fail()
	} else if attr.Name != "account" {
		t.Fail()
	} else if slices.Compare([]string{"social network"}, attr.Semantics) != 0 {
		t.Log("missing semantics")
		t.Fail()
	}

	c = models.NewFormalCompleteClass("humans", map[string][]string{"account": {"social network"}})

	if attr, found := c.Attributes["account"]; !found {
		t.Log("missing attribute")
		t.Fail()
	} else if attr.Name != "account" {
		t.Fail()
	} else if slices.Compare([]string{"social network"}, attr.Semantics) != 0 {
		t.Log("missing semantics")
		t.Fail()
	}

}
