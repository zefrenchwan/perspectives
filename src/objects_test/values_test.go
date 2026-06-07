package objects_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/objects"
	"github.com/zefrenchwan/perspectives.git/periods"
)

func TestBuildFromScratch(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Hour * 24)
	if content, err := objects.NewLocalContentBuilder().
		WithActivity(periods.NewPeriodSince(before, true)).
		WithAttributeDuring("name", periods.NewFullPeriod(), "John").
		Build(); err != nil {
		t.Error(err)
	} else if content == nil {
		t.Error("content should not be nil")
	} else if values, exists := content.At(now); !exists {
		t.Error("values should exist for current time")
	} else if values == nil {
		t.Error("values should not be nil")
	} else if len(values) != 1 {
		t.Errorf("expected 1 value, got %d", len(values))
	} else if values["name"] != "John" {
		t.Errorf("expected 'John', got '%s'", values["name"])
	} else if description := content.Description(); len(description) != 1 {
		t.Error("description should not be empty")
	} else if description["name"] != "string" {
		t.Errorf("expected 'string', got '%s'", description["name"])
	}
}
