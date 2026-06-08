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

func TestBuildFromOther(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Hour * 24)
	if content, err := objects.NewLocalContentBuilder().
		WithActivity(periods.NewPeriodSince(before, true)).
		WithAttributeDuring("name", periods.NewFullPeriod(), "John").
		Build(); err != nil {
		t.Error(err)
	} else if content == nil {
		t.Error("content should not be nil")
	} else if other, errOther := objects.LocalContentBuilderLoad(content).Build(); errOther != nil {
		t.Error(errOther)
	} else if other == nil {
		t.Error("other should not be nil")
	} else if values, exists := other.At(now); !exists {
		t.Error("values should exist for current time")
	} else if values == nil {
		t.Error("values should not be nil")
	} else if len(values) != 1 {
		t.Errorf("expected 1 value, got %d", len(values))
	} else if values["name"] != "John" {
		t.Errorf("expected 'John', got '%s'", values["name"])
	} else if description := other.Description(); len(description) != 1 {
		t.Error("description should not be empty")
	} else if description["name"] != "string" {
		t.Errorf("expected 'string', got '%s'", description["name"])
	}
}

func TestBuildError(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Hour * 24)
	if _, err := objects.NewLocalContentBuilder().
		WithActivity(periods.NewPeriodSince(before, true)).
		WithAttributeDuring("name", periods.NewFullPeriod(), "John").
		WithAttributeDuring("name", periods.NewFullPeriod(), 10).
		Build(); err == nil {
		t.Error("expected error for invalid attribute that changed its type")
	}
}
